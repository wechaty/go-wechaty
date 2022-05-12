package puppetservice

import (
	"bytes"
	"errors"
	pbwechaty "github.com/wechaty/go-grpc/wechaty"
	pbwechatypuppet "github.com/wechaty/go-grpc/wechaty/puppet"
	"github.com/wechaty/go-wechaty/wechaty-puppet/filebox"
	"github.com/wechaty/go-wechaty/wechaty-puppet/helper"
	"io"
)

// ErrNoName err no name
var ErrNoName = errors.New("no name")

// NewFileBoxFromMessageFileStream ...
func NewFileBoxFromMessageFileStream(client pbwechaty.Puppet_MessageFileStreamClient) (*filebox.FileBox, error) {
	recv, err := client.Recv()
	if err != nil {
		return nil, err
	}
	name := recv.FileBoxChunk.GetName()
	if name == "" {
		return nil, ErrNoName
	}

	return filebox.FromStream(NewMessageFile(client), filebox.WithName(name)), nil
}

// MessageFile 把 grpc 流包装到 io.Reader 接口
type MessageFile struct {
	client pbwechaty.Puppet_MessageFileStreamClient
	buffer bytes.Buffer
	done   bool
}

// Read 把 grpc 流包装到 io.Reader 接口
func (m *MessageFile) Read(p []byte) (n int, err error) {
	if m.done {
		return m.buffer.Read(p)
	}

	for {
		if m.buffer.Len() >= len(p) {
			break
		}
		recv, err := m.client.Recv()
		if err == io.EOF {
			m.done = true
			break
		}
		if err != nil {
			return 0, err
		}
		_, err = m.buffer.Write(recv.FileBoxChunk.GetData())
		if err != nil {
			return 0, err
		}
	}
	return m.buffer.Read(p)
}

// NewMessageFile ...
func NewMessageFile(client pbwechaty.Puppet_MessageFileStreamClient) *MessageFile {
	return &MessageFile{
		client: client,
		buffer: bytes.Buffer{},
		done:   false,
	}
}

// MessageSendFile 把 grpc 流包装到 io.Writer 接口
type MessageSendFile struct {
	client  pbwechaty.Puppet_MessageSendFileStreamClient
	fileBox *filebox.FileBox
	count   int
}

// Write 把 grpc 流包装到 io.Writer 接口
func (m *MessageSendFile) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	fileDataRequest := &pbwechatypuppet.MessageSendFileStreamRequest{
		FileBoxChunk: &pbwechatypuppet.FileBoxChunk{
			Data: p,
			Name: nil,
		},
	}
	m.count++
	if err := m.client.Send(fileDataRequest); err != nil {
		return 0, err
	}
	return len(p), nil
}

// ToMessageSendFileWriter 把 grpc 流包装到 io.Writer 接口
func ToMessageSendFileWriter(client pbwechaty.Puppet_MessageSendFileStreamClient, conversationID string, fileBox *filebox.FileBox) (io.Writer, error) {
	// 发送 conversationID
	{
		idRequest := &pbwechatypuppet.MessageSendFileStreamRequest{
			ConversationId: &conversationID,
		}
		if err := client.Send(idRequest); err != nil {
			return nil, err
		}
	}

	// 发送 fileName
	{
		fileNameRequest := &pbwechatypuppet.MessageSendFileStreamRequest{
			FileBoxChunk: &pbwechatypuppet.FileBoxChunk{
				Name: &fileBox.Name,
			},
		}
		if err := client.Send(fileNameRequest); err != nil {
			return nil, err
		}
	}
	return &MessageSendFile{
		client:  client,
		fileBox: fileBox,
	}, nil
}

// DownloadFile 把 grpc download 流包装到 io.Reader 接口
type DownloadFile struct {
	client pbwechaty.Puppet_DownloadClient
	buffer bytes.Buffer
	done   bool
}

// Read 把 grpc download 流包装到 io.Reader 接口
func (m *DownloadFile) Read(p []byte) (n int, err error) {
	if m.done {
		return m.buffer.Read(p)
	}

	for {
		if m.buffer.Len() >= len(p) {
			break
		}
		recv, err := m.client.Recv()
		if err == io.EOF {
			m.done = true
			break
		}
		if err != nil {
			return 0, err
		}
		_, err = m.buffer.Write(recv.Chunk)
		if err != nil {
			return 0, err
		}
	}
	return m.buffer.Read(p)
}

// NewDownloadFile 把 grpc download 流包装到 io.Reader 接口
func NewDownloadFile(client pbwechaty.Puppet_DownloadClient) *DownloadFile {
	return &DownloadFile{
		client: client,
		buffer: bytes.Buffer{},
		done:   false,
	}
}

/**
 *  for testing propose, use 20KB as the threshold
 *  after stable we should use a value between 64KB to 256KB as the threshold
 */
const passThroughThresholdBytes = 20 * 1024 //nolint:unused, deadcode, varcheck  // TODO 未来会被用到

/**
 * 1. Green:
 *  Can be serialized directly
 */
var greenFileBoxTypes = helper.ArrayInt{
	filebox.TypeUrl,
	filebox.TypeUuid,
	filebox.TypeQRCode,
}

/**
 * 2. Yellow:
 *  Can be serialized directly, if the size is less than a threshold
 *  if it's bigger than the threshold,
 *  then it should be convert to a UUID file box before send out
 */
var yellowFileBoxTypes = helper.ArrayInt{
	filebox.TypeBase64,
}

func serializeFileBox(box *filebox.FileBox) (*filebox.FileBox, error) {
	if canPassthrough(box) {
		return box, nil
	}
	reader, err := box.ToReader()
	if err != nil {
		return nil, err
	}
	uuid, err := filebox.FromStream(reader).ToUuid()
	if err != nil {
		return nil, err
	}
	return filebox.FromUuid(uuid, filebox.WithName(box.Name)), nil
}

func canPassthrough(box *filebox.FileBox) bool {
	if greenFileBoxTypes.InArray(int(box.Type())) {
		return true
	}

	if !yellowFileBoxTypes.InArray(int(box.Type())) {
		return false
	}

	// checksize
	return true
}
