package puppetservice

import (
	"bytes"
	"errors"
	pbwechaty "github.com/wechaty/go-grpc/wechaty"
	pbwechatypuppet "github.com/wechaty/go-grpc/wechaty/puppet"
	"github.com/wechaty/go-wechaty/wechaty-puppet/filebox"
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
			err = nil
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
			err = nil
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
