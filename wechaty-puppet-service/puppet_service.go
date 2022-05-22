package puppetservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	pbwechaty "github.com/wechaty/go-grpc/wechaty"
	pbwechatypuppet "github.com/wechaty/go-grpc/wechaty/puppet"
	wechatyPuppet "github.com/wechaty/go-wechaty/wechaty-puppet"
	"github.com/wechaty/go-wechaty/wechaty-puppet/filebox"
	"github.com/wechaty/go-wechaty/wechaty-puppet/helper"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"log"
	"time"
)

var (
	// ErrNoEndpoint err no endpoint
	ErrNoEndpoint = errors.New("no endpoint")
	// ErrURLLinkPayloadNotFound ...
	ErrURLLinkPayloadNotFound = errors.New("UrlLinkPayloadNotFound")
)

var pbEventType2PuppetEventName = schemas.PbEventType2PuppetEventName()

var pbEventType2GeneratePayloadFunc = schemas.PbEventType2GeneratePayloadFunc()

var _ wechatyPuppet.IPuppetAbstract = &PuppetService{}

// PuppetService struct
type PuppetService struct {
	*wechatyPuppet.Puppet
	grpcConn    *grpc.ClientConn
	grpcClient  pbwechaty.PuppetClient
	eventStream pbwechaty.Puppet_EventClient

	stop chan struct{}
}

// NewPuppetService new PuppetService struct
func NewPuppetService(o wechatyPuppet.Option) (*PuppetService, error) {
	if o.Token == "" {
		o.Token = getPuppetServiceTokenFromEnv()
	}
	if o.Endpoint == "" {
		o.Endpoint = getPuppetServiceEndpointFromEnv()
	}
	puppetAbstract, err := wechatyPuppet.NewPuppet(o)
	if err != nil {
		return nil, err
	}
	puppetService := &PuppetService{
		Puppet: puppetAbstract,
		stop:   make(chan struct{}, 1),
	}
	puppetAbstract.SetPuppetImplementation(puppetService)
	return puppetService, nil
}

func getPuppetServiceTokenFromEnv() string {
	if WechatyPuppetServiceToken != "" {
		return WechatyPuppetServiceToken
	}
	if WechatyPuppetHostieToken != "" {
		log.Println(`warn: WECHATY_PUPPET_HOSTIE_TOKEN environment be deprecated
please use new environment name<WECHATY_PUPPET_SERVICE_TOKEN> to avoid unnecessary bugs`)
		return WechatyPuppetHostieToken
	}
	return ""
}

func getPuppetServiceEndpointFromEnv() string {
	if WechatyPuppetServiceEndpoint != "" {
		return WechatyPuppetServiceEndpoint
	}
	if WechatyPuppetServiceEndpoint != "" {
		log.Println(`warn: WECHATY_PUPPET_HOSTIE_ENDPOINT environment be deprecated
please use new environment name<WECHATY_PUPPET_SERVICE_ENDPOINT> to avoid unnecessary bugs`)
		return WechatyPuppetServiceEndpoint
	}
	return ""
}

// MessageImage ...
func (p *PuppetService) MessageImage(messageID string, imageType schemas.ImageType) (*filebox.FileBox, error) {
	log.Printf("PuppetService MessageImage(%s, %s)\n", messageID, imageType)
	response, err := p.grpcClient.MessageImage(context.Background(), &pbwechatypuppet.MessageImageRequest{
		Id:   messageID,
		Type: pbwechatypuppet.ImageType(imageType),
	})
	if err != nil {
		return nil, err
	}
	return filebox.FromJSON(response.FileBox), nil
}

// Start ...
func (p *PuppetService) Start() (err error) {
	log.Println("PuppetService Start()")
	defer func() {
		if err != nil {
			err = fmt.Errorf("PuppetService Start() rejection: %w", err)
		}
	}()

	err = p.startGrpcClient()
	if err != nil {
		return err
	}

	filebox.SetUuidLoader(p.uuidLoader)
	filebox.SetUuidSaver(p.uuidSaver)

	err = p.startGrpcStream()
	if err != nil {
		return err
	}
	_, err = p.grpcClient.Start(context.Background(), &pbwechatypuppet.StartRequest{})
	if err != nil {
		return err
	}
	return nil
}

func (p *PuppetService) uuidSaver(reader io.Reader) (uuid string, err error) {
	client, err := p.grpcClient.Upload(context.Background())
	if err != nil {
		return "", err
	}

	b := make([]byte, 4000000)
	for {
		l, err := reader.Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		err = client.Send(&pbwechatypuppet.UploadRequest{Chunk: b[0:l]})
		if err != nil {
			return "", err
		}
	}

	response, err := client.CloseAndRecv()
	if err != nil {
		return "", err
	}
	return response.Id, nil
}

func (p *PuppetService) uuidLoader(uuid string) (io.Reader, error) {
	client, err := p.grpcClient.Download(context.Background(), &pbwechatypuppet.DownloadRequest{
		Id: uuid,
	})
	if err != nil {
		return nil, err
	}
	return NewDownloadFile(client), nil
}

// Stop ...
func (p *PuppetService) Stop() {
	p.stop <- struct{}{}
	var err error
	defer func() {
		if err != nil {
			log.Printf("PuppetService Stop err: %s\n", err)
		}
	}()
	if p.logonoff() {
		p.Emit(schemas.EventLogoutPayload{
			ContactId: p.SelfID(),
			Data:      "PuppetService Stop()",
		})
		p.SetID("")
	}

	if err = p.stopGrpcStream(); err != nil {
		return
	}

	if p.grpcClient != nil {
		if _, err = p.grpcClient.Stop(context.Background(), &pbwechatypuppet.StopRequest{}); err != nil {
			return
		}
	}

	if err = p.stopGrpcClient(); err != nil {
		return
	}
}

func (p *PuppetService) stopGrpcClient() error {
	if p.grpcConn == nil {
		return errors.New("puppetClient had not inited")
	}
	p.grpcConn.Close()
	p.grpcConn = nil
	p.grpcClient = nil
	return nil
}

func (p *PuppetService) stopGrpcStream() error {
	log.Println("PuppetService stopGrpcStream()")

	if p.eventStream == nil {
		return errors.New("no event stream")
	}

	if err := p.eventStream.CloseSend(); err != nil {
		log.Printf("PuppetService stopGrpcStream() err: %s\n", err)
	}
	p.eventStream = nil
	return nil
}

func (p *PuppetService) logonoff() bool {
	return p.SelfID() != ""
}

func (p *PuppetService) startGrpcClient() error {
	endpoint := p.Endpoint
	if endpoint == "" {
		serviceEndPoint, err := p.discoverServiceEndPoint()
		if err != nil {
			return err
		}
		if !serviceEndPoint.IsValid() {
			return ErrNoEndpoint
		}
		endpoint = serviceEndPoint.Target()
	}

	// TODO 支持 tls
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithAuthority(p.Token))
	if err != nil {
		return err
	}
	p.grpcConn = conn

	go p.autoReconnectGrpcConn()

	p.grpcClient = pbwechaty.NewPuppetClient(conn)
	return nil
}

func (p *PuppetService) autoReconnectGrpcConn() {
	interval := 2 * time.Second
	if p.Option.GrpcReconnectInterval > 0 {
		interval = p.Option.GrpcReconnectInterval
	}
	ticker := time.NewTicker(interval)
	isClose := false
	for {
		select {
		case <-ticker.C:
			connState := p.grpcConn.GetState()
			// 重新连接成功
			if isClose && connectivity.Ready == connState {
				isClose = false
				log.Printf("PuppetService.autoReconnectGrpcConn grpc reconnection successful")
				if err := p.startGrpcStream(); err != nil {
					log.Printf("PuppetService.autoReconnectGrpcConn startGrpcStream err:%s", err.Error())
				}
			}

			if p.grpcConn.GetState() == connectivity.Idle {
				isClose = true
				p.grpcConn.Connect()
				log.Printf("PuppetService.autoReconnectGrpcConn grpc reconnection...")
			}
		case <-p.stop:
			return
		}
	}
}

func (p *PuppetService) startGrpcStream() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("startGrpcStream err:%w", err)
		}
	}()
	if p.eventStream != nil {
		return errors.New("event stream exists")
	}
	p.eventStream, err = p.grpcClient.Event(context.Background(), &pbwechatypuppet.EventRequest{})
	if err != nil {
		return err
	}

	go func() {
		for {
			reply, err := p.eventStream.Recv()
			if err == io.EOF {
				log.Println("eventStream.Recv EOF")
				break
			}
			if err != nil {
				log.Printf("PuppetService startGrpcStream() eventStream err %s", err)
				reason := "startGrpcStream() eventStream err: " + err.Error()
				p.Emit(schemas.PuppetEventNameReset, schemas.EventResetPayload{Data: reason})
				p.eventStream = nil
				break
			}
			go p.onGrpcStreamEvent(reply)
		}
	}()
	return nil
}

func (p *PuppetService) onGrpcStreamEvent(event *pbwechatypuppet.EventResponse) {
	log.Printf("PuppetService onGrpcStreamEvent({type:%s payload:%s})", event.Type, event.Payload)

	if event.Type != pbwechatypuppet.EventType_EVENT_TYPE_HEARTBEAT {
		p.Emit(schemas.PuppetEventNameHeartbeat, &schemas.EventHeartbeatPayload{
			Data: fmt.Sprintf("onGrpcStreamEvent(%s)", event.Type),
		})
	}
	if event.Type == pbwechatypuppet.EventType_EVENT_TYPE_UNSPECIFIED {
		log.Println("PuppetService onGrpcStreamEvent() got an EventType.EVENT_TYPE_UNSPECIFIED ")
		return
	}
	eventName, ok := pbEventType2PuppetEventName[event.Type]
	if !ok {
		log.Printf("'eventType %s unsupported! (code should not reach here)\n", event.Type)
		return
	}
	payload := pbEventType2GeneratePayloadFunc[event.Type]()
	p.unMarshal(event.Payload, payload)
	switch event.Type {
	case pbwechatypuppet.EventType_EVENT_TYPE_RESET:
		log.Println("PuppetService onGrpcStreamEvent() got an EventType.EVENT_TYPE_RESET ?")
		// the `reset` event should be dealed not send out
		return
	case pbwechatypuppet.EventType_EVENT_TYPE_LOGIN:
		p.SetID(payload.(*schemas.EventLoginPayload).ContactId)
	case pbwechatypuppet.EventType_EVENT_TYPE_LOGOUT:
		p.SetID("")
	}
	p.Emit(eventName, payload)
}

func (p *PuppetService) unMarshal(data string, v interface{}) {
	err := json.Unmarshal([]byte(data), v)
	if err != nil {
		log.Printf("PuppetService unMarshal err: %s\n", err)
	}
}

// Logout ...
func (p *PuppetService) Logout() error {
	log.Println("PuppetService Logout()")
	if !p.logonoff() {
		return errors.New("logout before login? ")
	}
	_, err := p.grpcClient.Logout(context.Background(), &pbwechatypuppet.LogoutRequest{})
	if err != nil {
		return fmt.Errorf("PuppetService Logout() err: %w", err)
	}
	go p.Emit(schemas.PuppetEventNameLogout, &schemas.EventLogoutPayload{
		ContactId: p.SelfID(),
	})
	p.SetID("")
	return nil
}

// Ding ...
func (p *PuppetService) Ding(data string) {
	log.Printf("PuppetService Ding(%s)\n", data)
	_, err := p.grpcClient.Ding(context.Background(), &pbwechatypuppet.DingRequest{
		Data: data,
	})
	if err != nil {
		log.Printf("PuppetService Ding() err: %s\n", err)
	}
}

// SetContactAlias ...
func (p *PuppetService) SetContactAlias(contactID string, alias string) error {
	log.Printf("PuppetService, SetContactAlias(%s, %s)\n", contactID, alias)
	_, err := p.grpcClient.ContactAlias(context.Background(), &pbwechatypuppet.ContactAliasRequest{
		Id:    contactID,
		Alias: &alias,
	})
	if err != nil {
		return fmt.Errorf("PuppetService SetContactAlias err: %w", err)
	}
	return nil
}

// ContactAlias ...
func (p *PuppetService) ContactAlias(contactID string) (string, error) {
	log.Printf("PuppetService, 'ContactAlias(%s)\n", contactID)
	response, err := p.grpcClient.ContactAlias(context.Background(), &pbwechatypuppet.ContactAliasRequest{
		Id: contactID,
	})
	if err != nil {
		return "", fmt.Errorf("PuppetService ContactAlias err: %w", err)
	}
	return response.Alias, nil
}

// ContactList ...
func (p *PuppetService) ContactList() ([]string, error) {
	log.Println("PuppetService ContactList()")
	response, err := p.grpcClient.ContactList(context.Background(), &pbwechatypuppet.ContactListRequest{})
	if err != nil {
		return nil, fmt.Errorf("PuppetService ContactList err: %w", err)
	}
	return response.Ids, nil
}

// ContactQRCode ...
func (p *PuppetService) ContactQRCode(contactID string) (string, error) {
	log.Printf("PuppetService ContactQRCode(%s)\n", contactID)
	if contactID != p.SelfID() {
		return "", errors.New("can not set avatar for others")
	}
	response, err := p.grpcClient.ContactSelfQRCode(context.Background(), &pbwechatypuppet.ContactSelfQRCodeRequest{})
	if err != nil {
		return "", err
	}
	return response.Qrcode, nil
}

// SetContactAvatar ...
func (p *PuppetService) SetContactAvatar(contactID string, fileBox *filebox.FileBox) error {
	log.Printf("PuppetService SetContactAvatar(%s)\n", contactID)

	var err error
	fileBox, err = serializeFileBox(fileBox)
	if err != nil {
		return fmt.Errorf("serializeFileBox %w", err)
	}
	jsonString, err := fileBox.ToJSON()
	if err != nil {
		return fmt.Errorf("fileBox.ToJSON() %w", err)
	}
	_, err = p.grpcClient.ContactAvatar(context.Background(), &pbwechatypuppet.ContactAvatarRequest{
		Id:      contactID,
		FileBox: &jsonString,
	})
	if err != nil {
		return err
	}
	return nil
}

// ContactAvatar ...
func (p *PuppetService) ContactAvatar(contactID string) (*filebox.FileBox, error) {
	log.Printf("PuppetService ContactAvatar(%s)\n", contactID)
	response, err := p.grpcClient.ContactAvatar(context.Background(), &pbwechatypuppet.ContactAvatarRequest{
		Id: contactID,
	})
	if err != nil {
		return nil, err
	}
	return filebox.FromJSON(response.FileBox), nil
}

// ContactRawPayload ...
func (p *PuppetService) ContactRawPayload(contactID string) (*schemas.ContactPayload, error) {
	log.Printf("PuppetService ContactRawPayload(%s)\n", contactID)
	response, err := p.grpcClient.ContactPayload(context.Background(), &pbwechatypuppet.ContactPayloadRequest{
		Id: contactID,
	})
	if err != nil {
		return nil, err
	}
	return &schemas.ContactPayload{
		Id:        response.Id,
		Gender:    schemas.ContactGender(response.Gender),
		Type:      schemas.ContactType(response.Type),
		Name:      response.Name,
		Avatar:    response.Avatar,
		Address:   response.Address,
		Alias:     response.Alias,
		City:      response.City,
		Friend:    response.Friend,
		Province:  response.Province,
		Signature: response.Signature,
		Star:      response.Star,
		WeiXin:    response.Weixin,
	}, nil
}

// SetContactSelfName ...
func (p *PuppetService) SetContactSelfName(name string) error {
	log.Printf("PuppetService SetContactSelfName(%s)\n", name)
	_, err := p.grpcClient.ContactSelfName(context.Background(), &pbwechatypuppet.ContactSelfNameRequest{
		Name: name,
	})
	return err
}

// ContactSelfQRCode ...
func (p *PuppetService) ContactSelfQRCode() (string, error) {
	log.Println("PuppetService ContactSelfQRCode()")
	response, err := p.grpcClient.ContactSelfQRCode(context.Background(), &pbwechatypuppet.ContactSelfQRCodeRequest{})
	if err != nil {
		return "", err
	}
	return response.Qrcode, nil
}

// SetContactSelfSignature ...
func (p *PuppetService) SetContactSelfSignature(signature string) error {
	log.Printf("PuppetService SetContactSelfSignature(%s)\n", signature)
	_, err := p.grpcClient.ContactSelfSignature(context.Background(), &pbwechatypuppet.ContactSelfSignatureRequest{
		Signature: signature,
	})
	return err
}

// MessageRawMiniProgramPayload ...
func (p *PuppetService) MessageRawMiniProgramPayload(messageID string) (*schemas.MiniProgramPayload, error) {
	log.Printf("PuppetService MessageMiniProgram(%s)\n", messageID)
	response, err := p.grpcClient.MessageMiniProgram(context.Background(), &pbwechatypuppet.MessageMiniProgramRequest{
		Id: messageID,
	})
	if err != nil {
		return nil, err
	}

	// Deprecated: will be removed after Dec 31, 2022
	//nolint:staticcheck
	if response.MiniProgram == nil {
		payload := &schemas.MiniProgramPayload{}
		p.unMarshal(response.MiniProgramDeprecated, payload)
		return payload, nil
	}

	payload := &schemas.MiniProgramPayload{
		Appid:       response.MiniProgram.Appid,
		Description: response.MiniProgram.Description,
		PagePath:    response.MiniProgram.PagePath,
		ThumbUrl:    response.MiniProgram.ThumbUrl,
		Title:       response.MiniProgram.Title,
		Username:    response.MiniProgram.Username,
		ThumbKey:    response.MiniProgram.ThumbKey,
		ShareId:     response.MiniProgram.ShareId,
		IconUrl:     response.MiniProgram.IconUrl,
	}
	return payload, nil
}

// MessageContact ...
func (p *PuppetService) MessageContact(messageID string) (string, error) {
	log.Printf("PuppetService MessageContact(%s)\n", messageID)
	response, err := p.grpcClient.MessageContact(context.Background(), &pbwechatypuppet.MessageContactRequest{
		Id: messageID,
	})
	if err != nil {
		return "", err
	}
	return response.Id, nil
}

// MessageSendMiniProgram ...
func (p *PuppetService) MessageSendMiniProgram(conversationID string, miniProgramPayload *schemas.MiniProgramPayload) (string, error) {
	log.Printf("PuppetService MessageSendMiniProgram(%s,%#v)\n", conversationID, miniProgramPayload)
	response, err := p.grpcClient.MessageSendMiniProgram(context.Background(), &pbwechatypuppet.MessageSendMiniProgramRequest{
		ConversationId: conversationID,
		MiniProgram: &pbwechatypuppet.MiniProgramPayload{
			Appid:       miniProgramPayload.Appid,
			Description: miniProgramPayload.Description,
			PagePath:    miniProgramPayload.PagePath,
			IconUrl:     miniProgramPayload.IconUrl,
			ShareId:     miniProgramPayload.ShareId,
			ThumbUrl:    miniProgramPayload.ThumbUrl,
			Title:       miniProgramPayload.Title,
			Username:    miniProgramPayload.Username,
			ThumbKey:    miniProgramPayload.ThumbKey,
		},
		// Deprecated: will be removed after Dec 31, 2022
		MiniProgramDeprecated: miniProgramPayload.ToJson(),
	})
	if err != nil {
		return "", err
	}
	return response.Id, nil
}

// MessageRecall ...
func (p *PuppetService) MessageRecall(messageID string) (bool, error) {
	log.Printf("PuppetService MessageRecall(%s)\n", messageID)
	response, err := p.grpcClient.MessageRecall(context.Background(), &pbwechatypuppet.MessageRecallRequest{
		Id: messageID,
	})
	if err != nil {
		return false, err
	}
	return response.Success, nil
}

// MessageFile ...
func (p *PuppetService) MessageFile(id string) (*filebox.FileBox, error) {
	log.Printf("PuppetService MessageFile(%s)\n", id)
	response, err := p.grpcClient.MessageFileStream(context.Background(), &pbwechatypuppet.MessageFileStreamRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return NewFileBoxFromMessageFileStream(response)
}

// MessageRawPayload ...
func (p *PuppetService) MessageRawPayload(id string) (*schemas.MessagePayload, error) {
	log.Printf("PuppetService MessagePayload(%s)\n", id)
	response, err := p.grpcClient.MessagePayload(context.Background(), &pbwechatypuppet.MessagePayloadRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	payload := &schemas.MessagePayload{
		MessagePayloadBase: schemas.MessagePayloadBase{
			Id:            response.Id,
			MentionIdList: response.MentionIds,
			FileName:      response.Filename,
			Text:          response.Text,
			Type:          schemas.MessageType(response.Type),
		},
		MessagePayloadRoom: schemas.MessagePayloadRoom{
			TalkerId:   response.TalkerId,
			RoomId:     response.RoomId,
			ListenerId: response.ListenerId,
		},
	}

	if response.ReceiveTime != nil {
		payload.Timestamp = grpcTimestampToGoTime(response.ReceiveTime)
	} else {
		payload.Timestamp = time.Unix(int64(response.TimestampDeprecated), 0) //nolint:staticcheck
	}
	return payload, nil
}

func grpcTimestampToGoTime(t *timestamppb.Timestamp) time.Time {
	second := t.Seconds*1000 + int64(t.Nanos)/1000000
	return time.Unix(second, 0)
}

// MessageSendText ...
func (p *PuppetService) MessageSendText(conversationID string, text string, mentionIDList ...string) (string, error) {
	log.Printf("PuppetService messageSendText(%s, %s)\n", conversationID, text)
	response, err := p.grpcClient.MessageSendText(context.Background(), &pbwechatypuppet.MessageSendTextRequest{
		ConversationId: conversationID,
		Text:           text,
		MentionalIds:   mentionIDList,
	})
	if err != nil {
		return "", err
	}
	return response.Id, nil
}

// MessageSendFile ...
func (p *PuppetService) MessageSendFile(conversationID string, fileBox *filebox.FileBox) (string, error) {
	log.Printf("PuppetService MessageSendFile(%s)\n", conversationID)
	if msgID, err := p.messageSendFileNonStream(conversationID, fileBox); err == nil {
		return msgID, nil
	}

	return p.messageSendFileStream(conversationID, fileBox)
}

func (p *PuppetService) messageSendFileStream(conversationID string, fileBox *filebox.FileBox) (string, error) {
	stream, err := p.grpcClient.MessageSendFileStream(context.Background())
	if err != nil {
		return "", err
	}

	writer, err := ToMessageSendFileWriter(stream, conversationID, fileBox)
	if err != nil {
		return "", err
	}

	reader, err := fileBox.ToReader()
	if err != nil {
		return "", err
	}

	b := make([]byte, 4000000)
	for {
		l, err := reader.Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		_, err = writer.Write(b[0:l])
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		return "", err
	}
	return response.Id, nil
}

var serializableFileBoxTypes = helper.ArrayInt{
	filebox.TypeBase64,
	filebox.TypeUrl,
	filebox.TypeQRCode,
}

func (p *PuppetService) messageSendFileNonStream(conversationID string, fileBox *filebox.FileBox) (string, error) {
	log.Printf("PuppetService MessageSendFile(%s)\n", conversationID)
	var err error

	jsonText := ""
	if serializableFileBoxTypes.InArray(int(fileBox.Type())) {
		jsonText, err = fileBox.ToJSON()
		if err != nil {
			return "", err
		}
	} else {
		base64, err := fileBox.ToBase64()
		if err != nil {
			return "", err
		}
		jsonText, err = filebox.FromBase64(base64, filebox.WithName(fileBox.Name)).ToJSON()
		if err != nil {
			return "", err
		}
	}
	response, err := p.grpcClient.MessageSendFile(context.Background(), &pbwechatypuppet.MessageSendFileRequest{
		ConversationId: conversationID,
		FileBox:        jsonText,
	})
	if err != nil {
		return "", err
	}
	return response.Id, nil
}

// MessageSendContact ...
func (p *PuppetService) MessageSendContact(conversationID string, contactID string) (string, error) {
	log.Printf("PuppetService MessageSendContact(%s, %s)\n", conversationID, contactID)
	response, err := p.grpcClient.MessageSendContact(context.Background(), &pbwechatypuppet.MessageSendContactRequest{
		ConversationId: conversationID,
		ContactId:      contactID,
	})
	if err != nil {
		return "", err
	}
	return response.Id, nil
}

// MessageSendURL ...
func (p *PuppetService) MessageSendURL(conversationID string, urlLinkPayload *schemas.UrlLinkPayload) (string, error) {
	log.Printf("PuppetService MessageSendURL(%s, %+v)\n", conversationID, urlLinkPayload)
	response, err := p.grpcClient.MessageSendUrl(context.Background(), &pbwechatypuppet.MessageSendUrlRequest{
		ConversationId: conversationID,
		UrlLink: &pbwechatypuppet.UrlLinkPayload{
			Description:  urlLinkPayload.Description,
			ThumbnailUrl: urlLinkPayload.ThumbnailUrl,
			Title:        urlLinkPayload.Title,
			Url:          urlLinkPayload.Url,
		},

		// Deprecated: will be removed after Dec 31, 2022
		UrlLinkDeprecated: urlLinkPayload.ToJson(),
	})
	if err != nil {
		return "", err
	}
	return response.Id, nil
}

// MessageURL ...
func (p *PuppetService) MessageURL(messageID string) (*schemas.UrlLinkPayload, error) {
	log.Printf("PuppetService MessageURL(%s)\n", messageID)
	response, err := p.grpcClient.MessageUrl(context.Background(), &pbwechatypuppet.MessageUrlRequest{
		Id: messageID,
	})
	if err != nil {
		return nil, err
	}

	if response.UrlLink == nil {
		// Deprecated: will be removed after Dec 31, 2022
		payload := &schemas.UrlLinkPayload{}
		p.unMarshal(response.UrlLinkDeprecated, payload) //nolint:staticcheck
		return payload, nil
	}

	payload := &schemas.UrlLinkPayload{
		Description:  response.UrlLink.Description,
		ThumbnailUrl: response.UrlLink.ThumbnailUrl,
		Title:        response.UrlLink.Title,
		Url:          response.UrlLink.Url,
	}
	return payload, nil
}

// RoomRawPayload ...
func (p *PuppetService) RoomRawPayload(id string) (*schemas.RoomPayload, error) {
	log.Printf("PuppetService RoomRawPayload(%s)\n", id)
	response, err := p.grpcClient.RoomPayload(context.Background(), &pbwechatypuppet.RoomPayloadRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return &schemas.RoomPayload{
		Id:           response.Id,
		Topic:        response.Topic,
		Avatar:       response.Avatar,
		MemberIdList: response.MemberIds,
		OwnerId:      response.OwnerId,
		AdminIdList:  response.AdminIds,
	}, nil
}

// RoomList ...
func (p *PuppetService) RoomList() ([]string, error) {
	log.Printf("PuppetService RoomList()\n")
	response, err := p.grpcClient.RoomList(context.Background(), &pbwechatypuppet.RoomListRequest{})
	if err != nil {
		return nil, err
	}
	return response.Ids, nil
}

// RoomDel ...
func (p *PuppetService) RoomDel(roomID, contactID string) error {
	log.Printf("PuppetService roomDel(%s, %s)\n", roomID, contactID)
	_, err := p.grpcClient.RoomDel(context.Background(), &pbwechatypuppet.RoomDelRequest{
		Id:        roomID,
		ContactId: contactID,
	})
	if err != nil {
		return err
	}
	return nil
}

// RoomAvatar ...
func (p *PuppetService) RoomAvatar(roomID string) (*filebox.FileBox, error) {
	log.Printf("PuppetService RoomAvatar(%s)\n", roomID)
	response, err := p.grpcClient.RoomAvatar(context.Background(), &pbwechatypuppet.RoomAvatarRequest{
		Id: roomID,
	})
	if err != nil {
		return nil, err
	}
	return filebox.FromJSON(response.FileBox), nil
}

// RoomAdd ...
func (p *PuppetService) RoomAdd(roomID, contactID string) error {
	log.Printf("PuppetService RoomAdd(%s, %s)\n", roomID, contactID)
	_, err := p.grpcClient.RoomAdd(context.Background(), &pbwechatypuppet.RoomAddRequest{
		Id:        roomID,
		ContactId: contactID,
	})
	if err != nil {
		return err
	}
	return nil
}

// SetRoomTopic ...
func (p *PuppetService) SetRoomTopic(roomID string, topic string) error {
	log.Printf("PuppetService setRoomTopic(%s, %s)\n", roomID, topic)
	_, err := p.grpcClient.RoomTopic(context.Background(), &pbwechatypuppet.RoomTopicRequest{
		Id:    roomID,
		Topic: &topic,
	})
	return err
}

// RoomTopic ...
func (p *PuppetService) RoomTopic(roomID string) (string, error) {
	log.Printf("PuppetService RoomTopic(%s)\n", roomID)
	response, err := p.grpcClient.RoomTopic(context.Background(), &pbwechatypuppet.RoomTopicRequest{
		Id: roomID,
	})
	if err != nil {
		return "", err
	}
	return response.Topic, nil
}

// RoomCreate ...
func (p *PuppetService) RoomCreate(contactIDList []string, topic string) (string, error) {
	log.Printf("PuppetService roomCreate(%s, %s)\n", contactIDList, topic)
	response, err := p.grpcClient.RoomCreate(context.Background(), &pbwechatypuppet.RoomCreateRequest{
		ContactIds: contactIDList,
		Topic:      topic,
	})
	if err != nil {
		return "", err
	}
	return response.Id, nil
}

// RoomQuit ...
func (p *PuppetService) RoomQuit(roomID string) error {
	log.Printf("PuppetService RoomQuit(%s)\n", roomID)
	_, err := p.grpcClient.RoomQuit(context.Background(), &pbwechatypuppet.RoomQuitRequest{
		Id: roomID,
	})
	if err != nil {
		return err
	}
	return nil
}

// RoomQRCode ...
func (p *PuppetService) RoomQRCode(roomID string) (string, error) {
	log.Printf("PuppetService RoomQRCode(%s)\n", roomID)
	response, err := p.grpcClient.RoomQRCode(context.Background(), &pbwechatypuppet.RoomQRCodeRequest{
		Id: roomID,
	})
	if err != nil {
		return "", err
	}
	return response.Qrcode, nil
}

// RoomMemberList ...
func (p *PuppetService) RoomMemberList(roomID string) ([]string, error) {
	log.Printf("PuppetService RoomMemberList(%s)\n", roomID)
	response, err := p.grpcClient.RoomMemberList(context.Background(), &pbwechatypuppet.RoomMemberListRequest{
		Id: roomID,
	})
	if err != nil {
		return nil, err
	}
	return response.MemberIds, nil
}

// RoomMemberRawPayload ...
func (p *PuppetService) RoomMemberRawPayload(roomID string, contactID string) (*schemas.RoomMemberPayload, error) {
	log.Printf("PuppetService RoomMemberRawPayload(%s, %s)\n", roomID, contactID)
	response, err := p.grpcClient.RoomMemberPayload(context.Background(), &pbwechatypuppet.RoomMemberPayloadRequest{
		Id:       roomID,
		MemberId: contactID,
	})
	if err != nil {
		return nil, err
	}
	return &schemas.RoomMemberPayload{
		Id:        response.Id,
		RoomAlias: response.RoomAlias,
		InviterId: response.InviterId,
		Avatar:    response.Avatar,
		Name:      response.Name,
	}, nil
}

// SetRoomAnnounce ...
func (p *PuppetService) SetRoomAnnounce(roomID, text string) error {
	log.Printf("PuppetService SetRoomAnnounce(%s, %s)\n", roomID, text)
	_, err := p.grpcClient.RoomAnnounce(context.Background(), &pbwechatypuppet.RoomAnnounceRequest{
		Id:   roomID,
		Text: &text,
	})
	if err != nil {
		return err
	}
	return nil
}

// RoomAnnounce ...
func (p *PuppetService) RoomAnnounce(roomID string) (string, error) {
	log.Printf("PuppetService RoomAnnounce(%s)\n", roomID)
	response, err := p.grpcClient.RoomAnnounce(context.Background(), &pbwechatypuppet.RoomAnnounceRequest{
		Id: roomID,
	})
	if err != nil {
		return "", err
	}
	return response.Text, nil
}

// RoomInvitationAccept ...
func (p *PuppetService) RoomInvitationAccept(roomInvitationID string) error {
	log.Printf("PuppetService RoomInvitationAccept(%s)\n", roomInvitationID)
	_, err := p.grpcClient.RoomInvitationAccept(context.Background(), &pbwechatypuppet.RoomInvitationAcceptRequest{
		Id: roomInvitationID,
	})
	return err
}

// RoomInvitationRawPayload ...
func (p *PuppetService) RoomInvitationRawPayload(id string) (*schemas.RoomInvitationPayload, error) {
	log.Printf("PuppetService RoomInvitationRawPayload(%s)\n", id)
	response, err := p.grpcClient.RoomInvitationPayload(context.Background(), &pbwechatypuppet.RoomInvitationPayloadRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return &schemas.RoomInvitationPayload{
		Id:           response.Id,
		InviterId:    response.InviterId,
		Topic:        response.Topic,
		Avatar:       response.Avatar,
		Invitation:   response.Invitation,
		MemberCount:  int(response.MemberCount),
		MemberIdList: response.MemberIds,
		Timestamp:    grpcTimestampToGoTime(response.ReceiveTime),
		ReceiverId:   response.ReceiverId,
	}, nil
}

// FriendshipSearchPhone ...
func (p *PuppetService) FriendshipSearchPhone(phone string) (string, error) {
	log.Printf("PuppetService FriendshipSearchPhone(%s)\n", phone)
	response, err := p.grpcClient.FriendshipSearchPhone(context.Background(), &pbwechatypuppet.FriendshipSearchPhoneRequest{
		Phone: phone,
	})
	if err != nil {
		return "", err
	}
	return response.ContactId, nil
}

// FriendshipSearchWeixin ...
func (p *PuppetService) FriendshipSearchWeixin(weixin string) (string, error) {
	log.Printf("PuppetService FriendshipSearchWeixin(%s)\n", weixin)
	response, err := p.grpcClient.FriendshipSearchWeixin(context.Background(), &pbwechatypuppet.FriendshipSearchHandleRequest{
		Weixin: weixin,
	})
	if err != nil {
		return "", err
	}
	return response.ContactId, nil
}

// FriendshipRawPayload ...
func (p *PuppetService) FriendshipRawPayload(id string) (*schemas.FriendshipPayload, error) {
	log.Printf("PuppetService FriendshipRawPayload(%s)\n", id)
	response, err := p.grpcClient.FriendshipPayload(context.Background(), &pbwechatypuppet.FriendshipPayloadRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return &schemas.FriendshipPayload{
		FriendshipPayloadReceive: schemas.FriendshipPayloadReceive{
			FriendshipPayloadBase: schemas.FriendshipPayloadBase{
				ContactId: response.ContactId,
				Id:        response.Id,
				Hello:     response.Hello,
			},
			Type:     schemas.FriendshipType(response.Type),
			Scene:    schemas.FriendshipSceneType(response.Scene),
			Stranger: response.Stranger,
			Ticket:   response.Ticket,
		},
	}, nil
}

// FriendshipAdd ...
func (p *PuppetService) FriendshipAdd(contactID, hello string) (err error) {
	log.Printf("PuppetService FriendshipAdd(%s, %s)\n", contactID, hello)
	_, err = p.grpcClient.FriendshipAdd(context.Background(), &pbwechatypuppet.FriendshipAddRequest{
		ContactId: contactID,
		Hello:     hello,
	})
	return err
}

// FriendshipAccept ...
func (p *PuppetService) FriendshipAccept(friendshipID string) (err error) {
	log.Printf("PuppetService FriendshipAccept(%s)\n", friendshipID)
	_, err = p.grpcClient.FriendshipAccept(context.Background(), &pbwechatypuppet.FriendshipAcceptRequest{
		Id: friendshipID,
	})
	return err
}

// TagContactAdd ...
func (p *PuppetService) TagContactAdd(id, contactID string) (err error) {
	log.Printf("PuppetService TagContactAdd(%s, %s)\n", id, contactID)
	_, err = p.grpcClient.TagContactAdd(context.Background(), &pbwechatypuppet.TagContactAddRequest{
		Id:        id,
		ContactId: id,
	})
	return err
}

// TagContactRemove ...
func (p *PuppetService) TagContactRemove(id, contactID string) (err error) {
	log.Printf("PuppetService TagContactRemove(%s, %s)\n", id, contactID)
	_, err = p.grpcClient.TagContactRemove(context.Background(), &pbwechatypuppet.TagContactRemoveRequest{
		Id:        id,
		ContactId: contactID,
	})
	return err
}

// TagContactDelete ...
func (p *PuppetService) TagContactDelete(id string) (err error) {
	log.Printf("PuppetService TagContactDelete(%s)\n", id)
	_, err = p.grpcClient.TagContactDelete(context.Background(), &pbwechatypuppet.TagContactDeleteRequest{
		Id: id,
	})
	return err
}

// TagContactList ...
func (p *PuppetService) TagContactList(contactID string) ([]string, error) {
	log.Printf("PuppetService TagContactList(%s)\n", contactID)
	request := &pbwechatypuppet.TagContactListRequest{}
	if contactID != "" {
		request.ContactId = contactID
	}
	response, err := p.grpcClient.TagContactList(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return response.Ids, nil
}

// DirtyPayload ...
func (p *PuppetService) DirtyPayload(payloadType schemas.PayloadType, id string) error {
	log.Printf("PuppetService DirtyPayload(%v, %v)\n", payloadType, id)
	err := p.Puppet.OnDirty(payloadType, id)
	if err != nil {
		return err
	}
	request := &pbwechatypuppet.DirtyPayloadRequest{
		Type: pbwechatypuppet.PayloadType(payloadType),
		Id:   id,
	}
	_, err = p.grpcClient.DirtyPayload(context.Background(), request)
	if err != nil {
		return err
	}
	return nil
}

// MessageForward message forward
func (p *PuppetService) MessageForward(conversationID string, messageID string) (string, error) {
	log.Printf("PuppetService MessageForward(%v, %v)\n", conversationID, messageID)
	request := &pbwechatypuppet.MessageForwardRequest{
		MessageId:      messageID,
		ConversationId: conversationID,
	}
	response, err := p.grpcClient.MessageForward(context.Background(), request)
	if err != nil {
		return "", err
	}
	return response.Id, nil
}
