package puppethostie

import (
  "context"
  "encoding/json"
  "errors"
  "fmt"
  "github.com/golang/protobuf/ptypes/wrappers"
  "github.com/gorilla/websocket"
  pbwechaty "github.com/wechaty/go-grpc/wechaty"
  file_box "github.com/wechaty/go-wechaty/wechaty-puppet/file-box"
  option2 "github.com/wechaty/go-wechaty/wechaty-puppet/option"
  "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
  "google.golang.org/grpc"
  "io"
  "log"
)

// ErrNoEndpoint err no endpoint
var ErrNoEndpoint = errors.New("no endpoint")

var pbEventType2PuppetEventName = map[pbwechaty.EventType]schemas.PuppetEventName{
  pbwechaty.EventType_EVENT_TYPE_DONG:        schemas.PuppetEventNameDong,
  pbwechaty.EventType_EVENT_TYPE_ERROR:       schemas.PuppetEventNameError,
  pbwechaty.EventType_EVENT_TYPE_HEARTBEAT:   schemas.PuppetEventNameHeartbeat,
  pbwechaty.EventType_EVENT_TYPE_FRIENDSHIP:  schemas.PuppetEventNameFriendShip,
  pbwechaty.EventType_EVENT_TYPE_LOGIN:       schemas.PuppetEventNameLogin,
  pbwechaty.EventType_EVENT_TYPE_LOGOUT:      schemas.PuppetEventNameLogout,
  pbwechaty.EventType_EVENT_TYPE_MESSAGE:     schemas.PuppetEventNameMessage,
  pbwechaty.EventType_EVENT_TYPE_READY:       schemas.PuppetEventNameReady,
  pbwechaty.EventType_EVENT_TYPE_ROOM_INVITE: schemas.PuppetEventNameRoomInvite,
  pbwechaty.EventType_EVENT_TYPE_ROOM_JOIN:   schemas.PuppetEventNameRoomJoin,
  pbwechaty.EventType_EVENT_TYPE_ROOM_LEAVE:  schemas.PuppetEventNameRoomLeave,
  pbwechaty.EventType_EVENT_TYPE_ROOM_TOPIC:  schemas.PuppetEventNameRoomTopic,
  pbwechaty.EventType_EVENT_TYPE_SCAN:        schemas.PuppetEventNameScan,
  pbwechaty.EventType_EVENT_TYPE_RESET:       schemas.PuppetEventNameReset,
  pbwechaty.EventType_EVENT_TYPE_UNSPECIFIED: schemas.PuppetEventNameUnknown,
}

var pbEventType2GeneratePayloadFunc = map[pbwechaty.EventType]func() interface{}{
  pbwechaty.EventType_EVENT_TYPE_DONG:        func() interface{} { return &schemas.EventDongPayload{} },
  pbwechaty.EventType_EVENT_TYPE_ERROR:       func() interface{} { return &schemas.EventErrorPayload{} },
  pbwechaty.EventType_EVENT_TYPE_HEARTBEAT:   func() interface{} { return &schemas.EventHeartbeatPayload{} },
  pbwechaty.EventType_EVENT_TYPE_FRIENDSHIP:  func() interface{} { return &schemas.EventFriendshipPayload{} },
  pbwechaty.EventType_EVENT_TYPE_LOGIN:       func() interface{} { return &schemas.EventLoginPayload{} },
  pbwechaty.EventType_EVENT_TYPE_LOGOUT:      func() interface{} { return &schemas.EventLogoutPayload{} },
  pbwechaty.EventType_EVENT_TYPE_MESSAGE:     func() interface{} { return &schemas.EventMessagePayload{} },
  pbwechaty.EventType_EVENT_TYPE_READY:       func() interface{} { return &schemas.EventReadyPayload{} },
  pbwechaty.EventType_EVENT_TYPE_ROOM_INVITE: func() interface{} { return &schemas.EventRoomInvitePayload{} },
  pbwechaty.EventType_EVENT_TYPE_ROOM_JOIN:   func() interface{} { return &schemas.EventRoomJoinPayload{} },
  pbwechaty.EventType_EVENT_TYPE_ROOM_LEAVE:  func() interface{} { return &schemas.EventRoomLeavePayload{} },
  pbwechaty.EventType_EVENT_TYPE_ROOM_TOPIC:  func() interface{} { return &schemas.EventRoomTopicPayload{} },
  pbwechaty.EventType_EVENT_TYPE_SCAN:        func() interface{} { return &schemas.EventScanPayload{} },
  pbwechaty.EventType_EVENT_TYPE_RESET:       func() interface{} { return &schemas.EventResetPayload{} },
}

// PuppetHostie struct
type PuppetHostie struct {
  option      *option2.Option
  grpcConn    *grpc.ClientConn
  grpcClient  pbwechaty.PuppetClient
  eventStream pbwechaty.Puppet_EventClient
  id          string
}

// NewPuppetHostie new PuppetHostie struct
func NewPuppetHostie(o *option2.Option) *PuppetHostie {
  return &PuppetHostie{
    option: o,
  }
}

// MessageImage ...
func (p *PuppetHostie) MessageImage(messageID string, imageType schemas.ImageType) (*file_box.FileBox, error) {
  log.Printf("PuppetHostie MessageImage(%s, %s)\n", messageID, imageType)
  response, err := p.grpcClient.MessageImage(context.Background(), &pbwechaty.MessageImageRequest{
    Id:   messageID,
    Type: pbwechaty.ImageType(imageType),
  })
  if err != nil {
    return nil, err
  }
  return file_box.NewFileBoxFromJSONString(response.Filebox)
}

// Start ...
func (p *PuppetHostie) Start() (err error) {
  log.Println("PuppetHostie Start()")
  defer func() {
    if err != nil {
      err = fmt.Errorf("PuppetHostie Start() rejection: %w", err)
    }
  }()

  err = p.startGrpcClient()
  if err != nil {
    return err
  }
  err = p.startGrpcStream()
  if err != nil {
    return err
  }
  _, err = p.grpcClient.Start(context.Background(), &pbwechaty.StartRequest{})
  if err != nil {
    return err
  }
  return nil
}

// Stop ...
func (p *PuppetHostie) Stop() {
  var err error
  defer func() {
    if err != nil {
      log.Printf("PuppetHostie Stop err: %s\n", err)
    }
  }()
  if p.logonoff() {
    p.option.Emit(schemas.EventLogoutPayload{
      ContactId: p.id,
      Data:      "PuppetHostie Stop()",
    })
    p.id = ""
  }

  if err = p.stopGrpcStream(); err != nil {
    return
  }

  if p.grpcClient != nil {
    if _, err = p.grpcClient.Stop(context.Background(), &pbwechaty.StopRequest{}); err != nil {
      return
    }
  }

  if err = p.stopGrpcClient(); err != nil {
    return
  }
}

func (p *PuppetHostie) stopGrpcClient() error {
  if p.grpcConn == nil {
    return errors.New("puppetClient had not inited")
  }
  p.grpcConn.Close()
  p.grpcConn = nil
  p.grpcClient = nil
  return nil
}

func (p *PuppetHostie) stopGrpcStream() error {
  log.Println("PuppetHostie stopGrpcStream()")

  if p.eventStream == nil {
    return errors.New("no event stream")
  }

  if err := p.eventStream.CloseSend(); err != nil {
    log.Printf("PuppetHostie stopGrpcStream() err: %s\n", err)
  }
  p.eventStream = nil
  return nil
}

func (p *PuppetHostie) logonoff() bool {
  return p.id != ""
}

func (p *PuppetHostie) startGrpcClient() error {
  endpoint := p.option.Endpoint
  if endpoint == "" {
    ip, err := p.discoverHostieIP()
    if err != nil {
      return err
    }
    if ip == "" || ip == "0.0.0.0" {
      return ErrNoEndpoint
    }
    endpoint = ip + ":8788"
  }
  conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
  if err != nil {
    return err
  }
  p.grpcConn = conn
  p.grpcClient = pbwechaty.NewPuppetClient(conn)
  return nil
}

func (p *PuppetHostie) discoverHostieIP() (s string, err error) {
  defer func() {
    if err != nil {
      err = fmt.Errorf("discoverHostieIP() err:%w", err)
    }
  }()
  if p.option.Token == "" {
    return "", errors.New("wechaty-puppet-hostie: token not found. See: <https://github.com/wechaty/wechaty-puppet-hostie#1-wechaty_puppet_hostie_token>")
  }
  const chatieEndpoint = "wss://api.chatie.io/v0/websocket/token/%s"
  const protocol = "puppet-hostie|0.0.1"

  dialer := websocket.Dialer{}
  dialer.Subprotocols = append(dialer.Subprotocols, protocol)
  conn, _, err := dialer.Dial(fmt.Sprintf(chatieEndpoint, p.option.Token), nil)
  if err != nil {
    return "", err
  }
  defer conn.Close()
  err = conn.WriteJSON(map[string]interface{}{
    "name": "hostie",
  })
  if err != nil {
    return "", err
  }

  result := make(chan string)
  errChan := make(chan error)
  go func() {
    for {
      _, message, err := conn.ReadMessage()
      if err != nil {
        errChan <- err
      }
      resultMap := make(map[string]string)
      err = json.Unmarshal(message, &resultMap)
      if err != nil {
        errChan <- err
      }
      if resultMap["name"] == "hostie" {
        result <- resultMap["payload"]
        return
      }
    }
  }()

  select {
  case ip := <-result:
    return ip, nil
  case err := <-errChan:
    return "", err
  }
}

func (p *PuppetHostie) startGrpcStream() (err error) {
  defer func() {
    if err != nil {
      err = fmt.Errorf("startGrpcStream err:%w", err)
    }
  }()
  if p.eventStream != nil {
    return errors.New("event stream exists")
  }
  p.eventStream, err = p.grpcClient.Event(context.Background(), &pbwechaty.EventRequest{})
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
        log.Printf("PuppetHostie startGrpcStream() eventStream err %s", err)
        reason := "startGrpcStream() eventStream err: " + err.Error()
        p.option.Emit(schemas.PuppetEventNameReset, schemas.EventResetPayload{Data: reason})
        break
      }
      go p.onGrpcStreamEvent(reply)
    }
  }()
  return nil
}

func (p *PuppetHostie) onGrpcStreamEvent(event *pbwechaty.EventResponse) {
  log.Printf("PuppetHostie onGrpcStreamEvent({type:%s payload:%s})", event.Type, event.Payload)

  if event.Type != pbwechaty.EventType_EVENT_TYPE_HEARTBEAT {
    p.option.Emit(schemas.PuppetEventNameHeartbeat, &schemas.EventHeartbeatPayload{
      Data: fmt.Sprintf("onGrpcStreamEvent(%s)", event.Type),
    })
  }
  if event.Type == pbwechaty.EventType_EVENT_TYPE_UNSPECIFIED {
    log.Println("PuppetHostie onGrpcStreamEvent() got an EventType.EVENT_TYPE_UNSPECIFIED ")
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
  case pbwechaty.EventType_EVENT_TYPE_RESET:
    log.Println("PuppetHostie onGrpcStreamEvent() got an EventType.EVENT_TYPE_RESET ?")
    // the `reset` event should be dealed not send out
    return
  case pbwechaty.EventType_EVENT_TYPE_LOGIN:
    p.id = payload.(*schemas.EventLoginPayload).ContactId
  case pbwechaty.EventType_EVENT_TYPE_LOGOUT:
    p.id = ""
  }
  p.option.Emit(eventName, payload)
}

func (p *PuppetHostie) unMarshal(data string, v interface{}) {
  err := json.Unmarshal([]byte(data), v)
  if err != nil {
    log.Printf("PuppetHostie unMarshal err: %s\n", err)
  }
}

// Logout ...
func (p *PuppetHostie) Logout() error {
  log.Println("PuppetHostie Logout()")
  if !p.logonoff() {
    return errors.New("logout before login? ")
  }
  _, err := p.grpcClient.Logout(context.Background(), &pbwechaty.LogoutRequest{})
  if err != nil {
    return fmt.Errorf("PuppetHostie Logout() err: %w", err)
  }
  go p.option.Emit(schemas.PuppetEventNameLogout, schemas.EventLogoutPayload{
    ContactId: p.id,
  })
  p.id = ""
  return nil
}

// Ding ...
func (p *PuppetHostie) Ding(data string) {
  log.Printf("PuppetHostie Ding(%s)\n", data)
  _, err := p.grpcClient.Ding(context.Background(), &pbwechaty.DingRequest{
    Data: data,
  })
  if err != nil {
    log.Printf("PuppetHostie Ding() err: %s\n", err)
  }
}

// SetContactAlias ...
func (p *PuppetHostie) SetContactAlias(contactID string, alias string) error {
  log.Printf("PuppetHostie, 'SetContactAlias(%s, %s)\n", contactID, alias)
  _, err := p.grpcClient.ContactAlias(context.Background(), &pbwechaty.ContactAliasRequest{
    Id: contactID,
    Alias: &wrappers.StringValue{
      Value: alias,
    },
  })
  if err != nil {
    return fmt.Errorf("PuppetHostie SetContactAlias err: %w", err)
  }
  return nil
}

// GetContactAlias ...
func (p *PuppetHostie) GetContactAlias(contactID string) (string, error) {
  log.Printf("PuppetHostie, 'GetContactAlias(%s)\n", contactID)
  response, err := p.grpcClient.ContactAlias(context.Background(), &pbwechaty.ContactAliasRequest{
    Id: contactID,
  })
  if err != nil {
    return "", fmt.Errorf("PuppetHostie GetContactAlias err: %w", err)
  }
  if response.Alias == nil {
    return "", fmt.Errorf("can not get aliasWrapper")
  }
  return response.Alias.Value, nil
}

// ContactList ...
func (p *PuppetHostie) ContactList() ([]string, error) {
  log.Println("PuppetHostie ContactList()")
  response, err := p.grpcClient.ContactList(context.Background(), &pbwechaty.ContactListRequest{})
  if err != nil {
    return nil, fmt.Errorf("PuppetHostie ContactList err: %w", err)
  }
  return response.Ids, nil
}

// ContactQRCode ...
func (p *PuppetHostie) ContactQRCode(contactID string) (string, error) {
  log.Printf("PuppetHostie ContactQRCode(%s)\n", contactID)
  if contactID != p.id {
    return "", errors.New("can not set avatar for others")
  }
  response, err := p.grpcClient.ContactSelfQRCode(context.Background(), &pbwechaty.ContactSelfQRCodeRequest{})
  if err != nil {
    return "", err
  }
  return response.Qrcode, nil
}

// SetContactAvatar ...
func (p *PuppetHostie) SetContactAvatar(contactID string, fileBox *file_box.FileBox) error {
  log.Printf("PuppetHostie SetContactAvatar(%s)\n", contactID)
  jsonString, err := fileBox.ToJSONString()
  if err != nil {
    return err
  }
  _, err = p.grpcClient.ContactAvatar(context.Background(), &pbwechaty.ContactAvatarRequest{
    Id: contactID,
    Filebox: &wrappers.StringValue{
      Value: jsonString,
    },
  })
  if err != nil {
    return nil
  }
  return nil
}

// GetContactAvatar ...
func (p *PuppetHostie) GetContactAvatar(contactID string) (*file_box.FileBox, error) {
  log.Printf("PuppetHostie GetContactAvatar(%s)\n", contactID)
  response, err := p.grpcClient.ContactAvatar(context.Background(), &pbwechaty.ContactAvatarRequest{
    Id: contactID,
  })
  if err != nil {
    return nil, err
  }
  return file_box.NewFileBoxFromJSONString(response.Filebox.Value)
}

// ContactPayload ...
func (p *PuppetHostie) ContactPayload(contactID string) (*schemas.ContactPayload, error) {
  log.Printf("PuppetHostie ContactPayload(%s)\n", contactID)
  response, err := p.grpcClient.ContactPayload(context.Background(), &pbwechaty.ContactPayloadRequest{
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
    Start:     response.Star,
    WeiXin:    response.Weixin,
  }, nil
}

// SetContactSelfName ...
func (p *PuppetHostie) SetContactSelfName(name string) error {
  log.Printf("PuppetHostie SetContactSelfName(%s)\n", name)
  _, err := p.grpcClient.ContactSelfName(context.Background(), &pbwechaty.ContactSelfNameRequest{
    Name: name,
  })
  return err
}

// ContactSelfQRCode ...
func (p *PuppetHostie) ContactSelfQRCode() (string, error) {
  log.Println("PuppetHostie ContactSelfQRCode()")
  response, err := p.grpcClient.ContactSelfQRCode(context.Background(), &pbwechaty.ContactSelfQRCodeRequest{})
  if err != nil {
    return "", err
  }
  return response.Qrcode, nil
}

// SetContactSelfSignature ...
func (p *PuppetHostie) SetContactSelfSignature(signature string) error {
  log.Printf("PuppetHostie SetContactSelfSignature(%s)\n", signature)
  _, err := p.grpcClient.ContactSelfSignature(context.Background(), &pbwechaty.ContactSelfSignatureRequest{
    Signature: signature,
  })
  return err
}

// MessageMiniProgram ...
func (p *PuppetHostie) MessageMiniProgram(messageID string) (*schemas.MiniProgramPayload, error) {
  log.Printf("PuppetHostie MessageMiniProgram(%s)\n", messageID)
  response, err := p.grpcClient.MessageMiniProgram(context.Background(), &pbwechaty.MessageMiniProgramRequest{
    Id: messageID,
  })
  if err != nil {
    return nil, err
  }
  payload := &schemas.MiniProgramPayload{}
  p.unMarshal(response.MiniProgram, payload)
  return payload, nil
}

// MessageContact ...
func (p *PuppetHostie) MessageContact(messageID string) (string, error) {
  log.Printf("PuppetHostie MessageContact(%s)\n", messageID)
  response, err := p.grpcClient.MessageContact(context.Background(), &pbwechaty.MessageContactRequest{
    Id: messageID,
  })
  if err != nil {
    return "", err
  }
  return response.Id, nil
}

// MessageSendMiniProgram ...
func (p *PuppetHostie) MessageSendMiniProgram(conversationID string, miniProgramPayload *schemas.MiniProgramPayload) (string, error) {
  log.Printf("PuppetHostie MessageSendMiniProgram(%s,%s)\n", conversationID, miniProgramPayload)
  response, err := p.grpcClient.MessageSendMiniProgram(context.Background(), &pbwechaty.MessageSendMiniProgramRequest{
    ConversationId: conversationID,
    MiniProgram:    miniProgramPayload.ToJson(),
  })
  if err != nil {
    return "", err
  }
  return response.Id.Value, nil
}

// MessageRecall ...
func (p *PuppetHostie) MessageRecall(messageID string) (bool, error) {
  log.Printf("PuppetHostie MessageRecall(%s)\n", messageID)
  response, err := p.grpcClient.MessageRecall(context.Background(), &pbwechaty.MessageRecallRequest{
    Id: messageID,
  })
  if err != nil {
    return false, err
  }
  return response.Success, nil
}

// MessageFile ...
func (p *PuppetHostie) MessageFile(id string) (*file_box.FileBox, error) {
  log.Printf("PuppetHostie MessageFile(%s)\n", id)
  response, err := p.grpcClient.MessageFile(context.Background(), &pbwechaty.MessageFileRequest{
    Id: id,
  })
  if err != nil {
    return nil, err
  }
  return file_box.NewFileBoxFromJSONString(response.Filebox)
}

// MessagePayload ...
func (p *PuppetHostie) MessagePayload(id string) (*schemas.MessagePayload, error) {
  log.Printf("PuppetHostie MessagePayload(%s)\n", id)
  response, err := p.grpcClient.MessagePayload(context.Background(), &pbwechaty.MessagePayloadRequest{
    Id: id,
  })
  if err != nil {
    return nil, err
  }
  return &schemas.MessagePayload{
    MessagePayloadBase: schemas.MessagePayloadBase{
      Id:            response.Id,
      MentionIdList: response.MentionIds,
      FileName:      response.Filename,
      Text:          response.Text,
      Timestamp:     response.Timestamp,
      Type:          schemas.MessageType(response.Type),
    },
    MessagePayloadRoom: schemas.MessagePayloadRoom{
      FromId: response.FromId,
      RoomId: response.RoomId,
      ToId:   response.ToId,
    },
  }, nil
}

// MessageSendText ...
func (p *PuppetHostie) MessageSendText(conversationID string, text string) (string, error) {
  log.Printf("PuppetHostie messageSendText(%s, %s)\n", conversationID, text)
  response, err := p.grpcClient.MessageSendText(context.Background(), &pbwechaty.MessageSendTextRequest{
    ConversationId: conversationID,
    Text:           text,
    MentonalIds:    nil,
  })
  if err != nil {
    return "", err
  }
  return response.Id.Value, nil
}

// MessageSendFile ...
func (p *PuppetHostie) MessageSendFile(conversationID string, fileBox *file_box.FileBox) (string, error) {
  log.Printf("PuppetHostie MessageSendFile(%s)\n", conversationID)
  jsonString, err := fileBox.ToJSONString()
  if err != nil {
    return "", err
  }
  response, err := p.grpcClient.MessageSendFile(context.Background(), &pbwechaty.MessageSendFileRequest{
    ConversationId: conversationID,
    Filebox:        jsonString,
  })
  if err != nil {
    return "", err
  }
  return response.Id.Value, nil
}

// MessageSendContact ...
func (p *PuppetHostie) MessageSendContact(conversationID string, contactID string) (string, error) {
  log.Printf("PuppetHostie MessageSendContact(%s, %s)\n", conversationID, contactID)
  response, err := p.grpcClient.MessageSendContact(context.Background(), &pbwechaty.MessageSendContactRequest{
    ConversationId: conversationID,
    ContactId:      contactID,
  })
  if err != nil {
    return "", err
  }
  return response.Id.Value, nil
}

// MessageSendURL ...
func (p *PuppetHostie) MessageSendURL(conversationID string, urlLinkPayload *schemas.UrlLinkPayload) (string, error) {
  log.Printf("PuppetHostie MessageSendURL(%s, %s)\n", conversationID, urlLinkPayload)
  response, err := p.grpcClient.MessageSendUrl(context.Background(), &pbwechaty.MessageSendUrlRequest{
    ConversationId: conversationID,
    UrlLink:        urlLinkPayload.ToJson(),
  })
  if err != nil {
    return "", err
  }
  return response.Id.Value, nil
}

// MessageURL ...
func (p *PuppetHostie) MessageURL(messageID string) (*schemas.UrlLinkPayload, error) {
  log.Printf("PuppetHostie MessageURL(%s)\n", messageID)
  response, err := p.grpcClient.MessageUrl(context.Background(), &pbwechaty.MessageUrlRequest{
    Id: messageID,
  })
  if err != nil {
    return nil, err
  }
  payload := &schemas.UrlLinkPayload{}
  p.unMarshal(response.UrlLink, payload)
  return payload, nil
}

// RoomPayload ...
func (p *PuppetHostie) RoomPayload(id string) (*schemas.RoomPayload, error) {
  log.Printf("PuppetHostie RoomPayload(%s)\n", id)
  response, err := p.grpcClient.RoomPayload(context.Background(), &pbwechaty.RoomPayloadRequest{
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
func (p *PuppetHostie) RoomList() ([]string, error) {
  log.Printf("PuppetHostie RoomList()\n")
  response, err := p.grpcClient.RoomList(context.Background(), &pbwechaty.RoomListRequest{})
  if err != nil {
    return nil, err
  }
  return response.Ids, nil
}

// RoomDel ...
func (p *PuppetHostie) RoomDel(roomID, contactID string) error {
  log.Printf("PuppetHostie roomDel(%s, %s)\n", roomID, contactID)
  _, err := p.grpcClient.RoomDel(context.Background(), &pbwechaty.RoomDelRequest{
    Id:        roomID,
    ContactId: contactID,
  })
  if err != nil {
    return err
  }
  return nil
}

// RoomAvatar ...
func (p *PuppetHostie) RoomAvatar(roomID string) (*file_box.FileBox, error) {
  log.Printf("PuppetHostie RoomAvatar(%s)\n", roomID)
  response, err := p.grpcClient.RoomAvatar(context.Background(), &pbwechaty.RoomAvatarRequest{
    Id: roomID,
  })
  if err != nil {
    return nil, err
  }
  return file_box.NewFileBoxFromJSONString(response.Filebox)
}

// RoomAdd ...
func (p *PuppetHostie) RoomAdd(roomID, contactID string) error {
  log.Printf("PuppetHostie RoomAdd(%s, %s)\n", roomID, contactID)
  _, err := p.grpcClient.RoomAdd(context.Background(), &pbwechaty.RoomAddRequest{
    Id:        roomID,
    ContactId: contactID,
  })
  if err != nil {
    return err
  }
  return nil
}

// SetRoomTopic ...
func (p *PuppetHostie) SetRoomTopic(roomID string, topic string) error {
  log.Printf("PuppetHostie setRoomTopic(%s, %s)\n", roomID, topic)
  _, err := p.grpcClient.RoomTopic(context.Background(), &pbwechaty.RoomTopicRequest{
    Id: roomID,
    Topic: &wrappers.StringValue{
      Value: topic,
    },
  })
  if err != nil {
    return err
  }
  return nil
}

// GetRoomTopic ...
func (p *PuppetHostie) GetRoomTopic(roomID string) (string, error) {
  log.Printf("PuppetHostie GetRoomTopic(%s)\n", roomID)
  response, err := p.grpcClient.RoomTopic(context.Background(), &pbwechaty.RoomTopicRequest{
    Id: roomID,
  })
  if err != nil {
    return "", err
  }
  return response.Topic.Value, nil
}

// RoomCreate ...
func (p *PuppetHostie) RoomCreate(contactIDList []string, topic string) (string, error) {
  log.Printf("PuppetHostie roomCreate(%s, %s)\n", contactIDList, topic)
  response, err := p.grpcClient.RoomCreate(context.Background(), &pbwechaty.RoomCreateRequest{
    ContactIds: contactIDList,
    Topic:      topic,
  })
  if err != nil {
    return "", err
  }
  return response.Id, nil
}

// RoomQuit ...
func (p *PuppetHostie) RoomQuit(roomID string) error {
  log.Printf("PuppetHostie RoomQuit(%s)\n", roomID)
  _, err := p.grpcClient.RoomQuit(context.Background(), &pbwechaty.RoomQuitRequest{
    Id: roomID,
  })
  if err != nil {
    return err
  }
  return nil
}

// RoomQRCode ...
func (p *PuppetHostie) RoomQRCode(roomID string) (string, error) {
  log.Printf("PuppetHostie RoomQRCode(%s)\n", roomID)
  response, err := p.grpcClient.RoomQRCode(context.Background(), &pbwechaty.RoomQRCodeRequest{
    Id: roomID,
  })
  if err != nil {
    return "", err
  }
  return response.Qrcode, nil
}

// RoomMemberList ...
func (p *PuppetHostie) RoomMemberList(roomID string) ([]string, error) {
  log.Printf("PuppetHostie RoomMemberList(%s)\n", roomID)
  response, err := p.grpcClient.RoomMemberList(context.Background(), &pbwechaty.RoomMemberListRequest{
    Id: roomID,
  })
  if err != nil {
    return nil, err
  }
  return response.MemberIds, nil
}

// RoomMemberPayload ...
func (p *PuppetHostie) RoomMemberPayload(roomID string, contactID string) (*schemas.RoomMemberPayload, error) {
  log.Printf("PuppetHostie RoomMemberPayload(%s, %s)\n", roomID, contactID)
  response, err := p.grpcClient.RoomMemberPayload(context.Background(), &pbwechaty.RoomMemberPayloadRequest{
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
func (p *PuppetHostie) SetRoomAnnounce(roomID, text string) error {
  log.Printf("PuppetHostie SetRoomAnnounce(%s, %s)\n", roomID, text)
  _, err := p.grpcClient.RoomAnnounce(context.Background(), &pbwechaty.RoomAnnounceRequest{
    Id:   roomID,
    Text: &wrappers.StringValue{Value: text},
  })
  if err != nil {
    return err
  }
  return nil
}

// GetRoomAnnounce ...
func (p *PuppetHostie) GetRoomAnnounce(roomID string) (string, error) {
  log.Printf("PuppetHostie GetRoomAnnounce(%s)\n", roomID)
  response, err := p.grpcClient.RoomAnnounce(context.Background(), &pbwechaty.RoomAnnounceRequest{
    Id: roomID,
  })
  if err != nil {
    return "", err
  }
  return response.Text.Value, nil
}

// RoomInvitationAccept ...
func (p *PuppetHostie) RoomInvitationAccept(roomInvitationID string) error {
  log.Printf("PuppetHostie RoomInvitationAccept(%s)\n", roomInvitationID)
  _, err := p.grpcClient.RoomInvitationAccept(context.Background(), &pbwechaty.RoomInvitationAcceptRequest{
    Id: roomInvitationID,
  })
  if err != nil {
    return err
  }
  return nil
}

// RoomInvitationPayload ...
func (p *PuppetHostie) RoomInvitationPayload(id string) (*schemas.RoomInvitationPayload, error) {
  log.Printf("PuppetHostie RoomInvitationPayload(%s)\n", id)
  response, err := p.grpcClient.RoomInvitationPayload(context.Background(), &pbwechaty.RoomInvitationPayloadRequest{
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
    Timestamp:    int64(response.Timestamp),
    ReceiverId:   response.ReceiverId,
  }, nil
}

// FriendshipSearchPhone ...
func (p *PuppetHostie) FriendshipSearchPhone(phone string) (string, error) {
  log.Printf("PuppetHostie FriendshipSearchPhone(%s)\n", phone)
  response, err := p.grpcClient.FriendshipSearchPhone(context.Background(), &pbwechaty.FriendshipSearchPhoneRequest{
    Phone: phone,
  })
  if err != nil {
    return "", err
  }
  return response.ContactId.Value, nil
}

// FriendshipSearchWeixin ...
func (p *PuppetHostie) FriendshipSearchWeixin(weixin string) (string, error) {
  log.Printf("PuppetHostie FriendshipSearchWeixin(%s)\n", weixin)
  response, err := p.grpcClient.FriendshipSearchWeixin(context.Background(), &pbwechaty.FriendshipSearchWeixinRequest{
    Weixin: weixin,
  })
  if err != nil {
    return "", err
  }
  return response.ContactId.Value, nil
}

// FriendshipPayload ...
func (p *PuppetHostie) FriendshipPayload(id string) (*schemas.FriendshipPayload, error) {
  log.Printf("PuppetHostie FriendshipPayload(%s)\n", id)
  response, err := p.grpcClient.FriendshipPayload(context.Background(), &pbwechaty.FriendshipPayloadRequest{
    Id: id,
  })
  if err != nil {
    return nil, err
  }
  return &schemas.FriendshipPayload{
    FriendshipPayloadReceive: schemas.FriendshipPayloadReceive{
      FriendshipPayloadBase: schemas.FriendshipPayloadBase{},
      Type:                  schemas.FriendshipType(response.Type),
      Scene:                 schemas.FriendshipSceneType(response.Scene),
      Stranger:              response.Stranger,
      Ticket:                response.Ticket,
    },
  }, nil
}

// FriendshipAdd ...
func (p *PuppetHostie) FriendshipAdd(contactID, hello string) (err error) {
  log.Printf("PuppetHostie FriendshipAdd(%s, %s)\n", contactID, hello)
  _, err = p.grpcClient.FriendshipAdd(context.Background(), &pbwechaty.FriendshipAddRequest{
    ContactId: contactID,
    Hello:     hello,
  })
  return err
}

// FriendshipAccept ...
func (p *PuppetHostie) FriendshipAccept(friendshipID string) (err error) {
  log.Printf("PuppetHostie FriendshipAccept(%s)\n", friendshipID)
  _, err = p.grpcClient.FrendshipAccept(context.Background(), &pbwechaty.FriendshipAcceptRequest{
    Id: friendshipID,
  })
  return err
}

// TagContactAdd ...
func (p *PuppetHostie) TagContactAdd(id, contactID string) (err error) {
  log.Printf("PuppetHostie TagContactAdd(%s, %s)\n", id, contactID)
  _, err = p.grpcClient.TagContactAdd(context.Background(), &pbwechaty.TagContactAddRequest{
    Id:        id,
    ContactId: id,
  })
  return err
}

// TagContactRemove ...
func (p *PuppetHostie) TagContactRemove(id, contactID string) (err error) {
  log.Printf("PuppetHostie TagContactRemove(%s, %s)\n", id, contactID)
  _, err = p.grpcClient.TagContactRemove(context.Background(), &pbwechaty.TagContactRemoveRequest{
    Id:        id,
    ContactId: contactID,
  })
  return err
}

// TagContactDelete ...
func (p *PuppetHostie) TagContactDelete(id string) (err error) {
  log.Printf("PuppetHostie TagContactDelete(%s)\n", id)
  _, err = p.grpcClient.TagContactDelete(context.Background(), &pbwechaty.TagContactDeleteRequest{
    Id: id,
  })
  return err
}

// TagContactList ...
func (p *PuppetHostie) TagContactList(contactID string) ([]string, error) {
  log.Printf("PuppetHostie TagContactList(%s)\n", contactID)
  response, err := p.grpcClient.TagContactList(context.Background(), &pbwechaty.TagContactListRequest{
    ContactId: &wrappers.StringValue{Value: contactID},
  })
  if err != nil {
    return nil, err
  }
  return response.Ids, nil
}