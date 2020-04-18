package wechaty_puppet_hostie

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
  pbwechaty.EventType_EVENT_TYPE_UNSPECIFIED: func() interface{} { return nil },
}

type PuppetHostie struct {
  option      *option2.Option
  grpcConn    *grpc.ClientConn
  grpcClient  pbwechaty.PuppetClient
  eventStream pbwechaty.Puppet_EventClient
  id          string
}

func NewPuppetHostie(optFns ...option2.OptionFn) *PuppetHostie {
  ph := &PuppetHostie{
    option: &option2.Option{},
  }

  for _, fn := range optFns {
    fn(ph.option)
  }
  return ph
}

func (p *PuppetHostie) MessageImage(ctx context.Context, messageID string, imageType schemas.ImageType) (*file_box.FileBox, error) {
  log.Printf("PuppetHostie MessageImage(%s, %s)\n", messageID, imageType)
  response, err := p.grpcClient.MessageImage(ctx, &pbwechaty.MessageImageRequest{
    Id:   messageID,
    Type: pbwechaty.ImageType(imageType),
  })
  if err != nil {
    return nil, err
  }
  return file_box.NewFileBoxFromJSONString(response.Filebox)
}

func (p *PuppetHostie) FriendshipPayloadReceive(friendshipID string) schemas.FriendshipPayloadReceive {
  panic("implement me")
}

func (p *PuppetHostie) FriendshipPayloadConfirm(friendshipID string) schemas.FriendshipPayloadConfirm {
  panic("implement me")
}

func (p *PuppetHostie) FriendshipPayloadVerify(friendshipID string) schemas.FriendshipPayloadVerify {
  panic("implement me")
}

func (p *PuppetHostie) Start(ctx context.Context) (err error) {
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
  err = p.startGrpcStream(ctx)
  if err != nil {
    return err
  }
  _, err = p.grpcClient.Start(ctx, &pbwechaty.StartRequest{})
  if err != nil {
    return err
  }
  return nil
}

func (p *PuppetHostie) Stop(ctx context.Context) {
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
    if _, err = p.grpcClient.Stop(ctx, &pbwechaty.StopRequest{}); err != nil {
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
    ip, err := p.discoverHostieIp()
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

func (p *PuppetHostie) discoverHostieIp() (s string, err error) {
  defer func() {
    if err != nil {
      err = fmt.Errorf("discoverHostieIp() err:%w", err)
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
      log.Println(string(message))
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

func (p *PuppetHostie) startGrpcStream(ctx context.Context) (err error) {
  defer func() {
    if err != nil {
      err = fmt.Errorf("startGrpcStream err:%w", err)
    }
  }()
  if p.eventStream != nil {
    return errors.New("event stream exists")
  }
  p.eventStream, err = p.grpcClient.Event(ctx, &pbwechaty.EventRequest{})
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
      }
      go p.onGrpcStreamEvent(reply)
    }
  }()
  return nil
}

func (p *PuppetHostie) onGrpcStreamEvent(event *pbwechaty.EventResponse) {
  log.Printf("PuppetHostie onGrpcStreamEvent({type:%s payload:%s})", event.Type, event.Payload)

  if event.Type != pbwechaty.EventType_EVENT_TYPE_HEARTBEAT {
    p.option.Emit(schemas.PuppetEventNameHeartbeat, schemas.EventHeartbeatPayload{
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

func (p *PuppetHostie) Logout(ctx context.Context) error {
  log.Println("PuppetHostie Logout()")
  if !p.logonoff() {
    return errors.New("logout before login? ")
  }
  _, err := p.grpcClient.Logout(ctx, &pbwechaty.LogoutRequest{})
  if err != nil {
    return fmt.Errorf("PuppetHostie Logout() err: %w", err)
  }
  go p.option.Emit(schemas.PuppetEventNameLogout, schemas.EventLogoutPayload{
    ContactId: p.id,
  })
  p.id = ""
  return nil
}

func (p *PuppetHostie) Ding(ctx context.Context, data string) {
  log.Printf("PuppetHostie Ding(%s)\n", data)
  _, err := p.grpcClient.Ding(ctx, &pbwechaty.DingRequest{
    Data: data,
  })
  if err != nil {
    log.Printf("PuppetHostie Ding() err: %s\n", err)
  }
}

func (p *PuppetHostie) SetContactAlias(ctx context.Context, contactId string, alias string) error {
  log.Printf("PuppetHostie, 'SetContactAlias(%s, %s)\n", contactId, alias)
  _, err := p.grpcClient.ContactAlias(ctx, &pbwechaty.ContactAliasRequest{
    Id: contactId,
    Alias: &wrappers.StringValue{
      Value: alias,
    },
  })
  if err != nil {
    return fmt.Errorf("PuppetHostie SetContactAlias err: %w", err)
  }
  return nil
}

func (p *PuppetHostie) GetContactAlias(ctx context.Context, contactId string) (string, error) {
  log.Printf("PuppetHostie, 'GetContactAlias(%s)\n", contactId)
  response, err := p.grpcClient.ContactAlias(ctx, &pbwechaty.ContactAliasRequest{
    Id: contactId,
  })
  if err != nil {
    return "", fmt.Errorf("PuppetHostie GetContactAlias err: %w", err)
  }
  if response.Alias == nil {
    return "", fmt.Errorf("can not get aliasWrapper")
  }
  return response.Alias.Value, nil
}

func (p *PuppetHostie) ContactList(ctx context.Context) ([]string, error) {
  log.Println("PuppetHostie ContactList()")
  response, err := p.grpcClient.ContactList(ctx, &pbwechaty.ContactListRequest{})
  if err != nil {
    return nil, fmt.Errorf("PuppetHostie ContactList err: %w", err)
  }
  return response.Ids, nil
}

func (p *PuppetHostie) ContactQRCode(ctx context.Context, contactId string) (string, error) {
  log.Printf("PuppetHostie ContactQRCode(%s)\n", contactId)
  if contactId != p.id {
    return "", errors.New("can not set avatar for others")
  }
  response, err := p.grpcClient.ContactSelfQRCode(ctx, &pbwechaty.ContactSelfQRCodeRequest{})
  if err != nil {
    return "", err
  }
  return response.Qrcode, nil
}

func (p *PuppetHostie) SetContactAvatar(ctx context.Context, contactId string, fileBox *file_box.FileBox) error {
  log.Printf("PuppetHostie SetContactAvatar(%s)\n", contactId)
  jsonString, err := fileBox.ToJSONString()
  if err != nil {
    return err
  }
  _, err = p.grpcClient.ContactAvatar(ctx, &pbwechaty.ContactAvatarRequest{
    Id: contactId,
    Filebox: &wrappers.StringValue{
      Value: jsonString,
    },
  })
  if err != nil {
    return nil
  }
  return nil
}

func (p *PuppetHostie) GetContactAvatar(ctx context.Context, contactId string) (*file_box.FileBox, error) {
  log.Printf("PuppetHostie GetContactAvatar(%s)\n", contactId)
  response, err := p.grpcClient.ContactAvatar(ctx, &pbwechaty.ContactAvatarRequest{
    Id: contactId,
  })
  if err != nil {
    return nil, err
  }
  return file_box.NewFileBoxFromJSONString(response.Filebox.Value)
}

func (p *PuppetHostie) ContactPayload(ctx context.Context, contactId string) (*schemas.ContactPayload, error) {
  log.Printf("PuppetHostie ContactPayload(%s)\n", contactId)
  response, err := p.grpcClient.ContactPayload(ctx, &pbwechaty.ContactPayloadRequest{
    Id: contactId,
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

func (p *PuppetHostie) SetContactSelfName(ctx context.Context, name string) error {
  log.Printf("PuppetHostie SetContactSelfName(%s)\n", name)
  _, err := p.grpcClient.ContactSelfName(ctx, &pbwechaty.ContactSelfNameRequest{
    Name: name,
  })
  return err
}

func (p *PuppetHostie) ContactSelfQRCode(ctx context.Context) (string, error) {
  log.Println("PuppetHostie ContactSelfQRCode()")
  response, err := p.grpcClient.ContactSelfQRCode(ctx, &pbwechaty.ContactSelfQRCodeRequest{})
  if err != nil {
    return "", err
  }
  return response.Qrcode, nil
}

func (p *PuppetHostie) SetContactSelfSignature(ctx context.Context, signature string) error {
  log.Printf("PuppetHostie SetContactSelfSignature(%s)\n", signature)
  _, err := p.grpcClient.ContactSelfSignature(ctx, &pbwechaty.ContactSelfSignatureRequest{
    Signature: signature,
  })
  return err
}

func (p *PuppetHostie) MessageMiniProgram(ctx context.Context, messageId string) (*schemas.MiniProgramPayload, error) {
  log.Printf("PuppetHostie MessageMiniProgram(%s)\n", messageId)
  response, err := p.grpcClient.MessageMiniProgram(ctx, &pbwechaty.MessageMiniProgramRequest{
    Id: messageId,
  })
  if err != nil {
    return nil, err
  }
  payload := &schemas.MiniProgramPayload{}
  p.unMarshal(response.MiniProgram, payload)
  return payload, nil
}

func (p *PuppetHostie) MessageContact(ctx context.Context, messageID string) (string, error) {
  log.Printf("PuppetHostie MessageContact(%s)\n", messageID)
  response, err := p.grpcClient.MessageContact(ctx, &pbwechaty.MessageContactRequest{
    Id: messageID,
  })
  if err != nil {
    return "", err
  }
  return response.Id, nil
}

func (p *PuppetHostie) MessageSendMiniProgram(ctx context.Context, conversationId string, miniProgramPayload *schemas.MiniProgramPayload) (string, error) {
  log.Printf("PuppetHostie MessageSendMiniProgram(%s,%s)\n", conversationId, miniProgramPayload)
  response, err := p.grpcClient.MessageSendMiniProgram(ctx, &pbwechaty.MessageSendMiniProgramRequest{
    ConversationId: conversationId,
    MiniProgram:    miniProgramPayload.ToJson(),
  })
  if err != nil {
    return "", err
  }
  return response.Id.Value, nil
}

func (p *PuppetHostie) MessageRecall(ctx context.Context, messageId string) (bool, error) {
  log.Printf("PuppetHostie MessageRecall(%s)\n", messageId)
  response, err := p.grpcClient.MessageRecall(ctx, &pbwechaty.MessageRecallRequest{
    Id: messageId,
  })
  if err != nil {
    return false, err
  }
  return response.Success, nil
}

func (p *PuppetHostie) MessageFile(ctx context.Context, id string) (*file_box.FileBox, error) {
  log.Printf("PuppetHostie MessageFile(%s)\n", id)
  response, err := p.grpcClient.MessageFile(ctx, &pbwechaty.MessageFileRequest{
    Id: id,
  })
  if err != nil {
    return nil, err
  }
  return file_box.NewFileBoxFromJSONString(response.Filebox)
}

func (p *PuppetHostie) MessagePayload(ctx context.Context, id string) (*schemas.MessagePayload, error) {
  log.Printf("PuppetHostie MessagePayload(%s)\n", id)
  response, err := p.grpcClient.MessagePayload(ctx, &pbwechaty.MessagePayloadRequest{
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

func (p *PuppetHostie) MessageSendText(ctx context.Context, conversationId string, text string) (string, error) {
  log.Printf("PuppetHostie messageSendText(%s, %s)\n", conversationId, text)
  response, err := p.grpcClient.MessageSendText(ctx, &pbwechaty.MessageSendTextRequest{
    ConversationId: conversationId,
    Text:           text,
    MentonalIds:    nil,
  })
  if err != nil {
    return "", err
  }
  return response.Id.Value, nil
}

func (p *PuppetHostie) MessageSendFile(ctx context.Context, conversationId string, fileBox *file_box.FileBox) (string, error) {
  log.Printf("PuppetHostie MessageSendFile(%s)\n", conversationId)
  jsonString, err := fileBox.ToJSONString()
  if err != nil {
    return "", err
  }
  response, err := p.grpcClient.MessageSendFile(ctx, &pbwechaty.MessageSendFileRequest{
    ConversationId: conversationId,
    Filebox:        jsonString,
  })
  if err != nil {
    return "", err
  }
  return response.Id.Value, nil
}

func (p *PuppetHostie) MessageSendContact(ctx context.Context, conversationId string, contactId string) (string, error) {
  log.Printf("PuppetHostie MessageSendContact(%s, %s)\n", conversationId, contactId)
  response, err := p.grpcClient.MessageSendContact(ctx, &pbwechaty.MessageSendContactRequest{
    ConversationId: conversationId,
    ContactId:      contactId,
  })
  if err != nil {
    return "", err
  }
  return response.Id.Value, nil
}

func (p *PuppetHostie) MessageSendUrl(ctx context.Context, conversationId string, urlLinkPayload *schemas.UrlLinkPayload) (string, error) {
  log.Printf("PuppetHostie MessageSendUrl(%s, %s)\n", conversationId, urlLinkPayload)
  response, err := p.grpcClient.MessageSendUrl(ctx, &pbwechaty.MessageSendUrlRequest{
    ConversationId: conversationId,
    UrlLink:        urlLinkPayload.ToJson(),
  })
  if err != nil {
    return "", err
  }
  return response.Id.Value, nil
}

func (p *PuppetHostie) MessageUrl(ctx context.Context, messageId string) (*schemas.UrlLinkPayload, error) {
  log.Printf("PuppetHostie MessageUrl(%s)\n", messageId)
  response, err := p.grpcClient.MessageUrl(ctx, &pbwechaty.MessageUrlRequest{
    Id: messageId,
  })
  if err != nil {
    return nil, err
  }
  payload := &schemas.UrlLinkPayload{}
  p.unMarshal(response.UrlLink, payload)
  return payload, nil
}

func (p *PuppetHostie) RoomPayload(ctx context.Context, id string) (*schemas.RoomPayload, error) {
  log.Printf("PuppetHostie RoomPayload(%s)\n", id)
  response, err := p.grpcClient.RoomPayload(ctx, &pbwechaty.RoomPayloadRequest{
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

func (p *PuppetHostie) RoomList(ctx context.Context) ([]string, error) {
  log.Printf("PuppetHostie RoomList()\n")
  response, err := p.grpcClient.RoomList(ctx, &pbwechaty.RoomListRequest{})
  if err != nil {
    return nil, err
  }
  return response.Ids, nil
}

func (p *PuppetHostie) RoomDel(ctx context.Context, roomId, contactId string) error {
  log.Printf("PuppetHostie roomDel(%s, %s)\n", roomId, contactId)
  _, err := p.grpcClient.RoomDel(ctx, &pbwechaty.RoomDelRequest{
    Id:        roomId,
    ContactId: contactId,
  })
  if err != nil {
    return err
  }
  return nil
}

func (p *PuppetHostie) RoomAvatar(ctx context.Context, roomId string) (*file_box.FileBox, error) {
  log.Printf("PuppetHostie RoomAvatar(%s)\n", roomId)
  response, err := p.grpcClient.RoomAvatar(ctx, &pbwechaty.RoomAvatarRequest{
    Id: roomId,
  })
  if err != nil {
    return nil, err
  }
  return file_box.NewFileBoxFromJSONString(response.Filebox)
}

func (p *PuppetHostie) RoomAdd(ctx context.Context, roomId, contactId string) error {
  log.Printf("PuppetHostie RoomAdd(%s, %s)\n", roomId, contactId)
  _, err := p.grpcClient.RoomAdd(ctx, &pbwechaty.RoomAddRequest{
    Id:        roomId,
    ContactId: contactId,
  })
  if err != nil {
    return err
  }
  return nil
}

func (p *PuppetHostie) setRoomTopic(ctx context.Context, roomId string, topic string) error {
  log.Printf("PuppetHostie setRoomTopic(%s, %s)\n", roomId, topic)
  _, err := p.grpcClient.RoomTopic(ctx, &pbwechaty.RoomTopicRequest{
    Id: roomId,
    Topic: &wrappers.StringValue{
      Value: topic,
    },
  })
  if err != nil {
    return err
  }
  return nil
}

func (p *PuppetHostie) GetRoomTopic(ctx context.Context, roomId string) (string, error) {
  log.Printf("PuppetHostie GetRoomTopic(%s)\n", roomId)
  response, err := p.grpcClient.RoomTopic(ctx, &pbwechaty.RoomTopicRequest{
    Id: roomId,
  })
  if err != nil {
    return "", err
  }
  return response.Topic.Value, nil
}

func (p *PuppetHostie) RoomCreate(ctx context.Context, contactIdList []string, topic string) (string, error) {
  log.Printf("PuppetHostie roomCreate(%s, %s)\n", contactIdList, topic)
  response, err := p.grpcClient.RoomCreate(ctx, &pbwechaty.RoomCreateRequest{
    ContactIds: contactIdList,
    Topic:      topic,
  })
  if err != nil {
    return "", err
  }
  return response.Id, nil
}

func (p *PuppetHostie) RoomQuit(ctx context.Context, roomId string) error {
  log.Printf("PuppetHostie RoomQuit(%s)\n", roomId)
  _, err := p.grpcClient.RoomQuit(ctx, &pbwechaty.RoomQuitRequest{
    Id: roomId,
  })
  if err != nil {
    return err
  }
  return nil
}

func (p *PuppetHostie) RoomQRCode(ctx context.Context, roomId string) (string, error) {
  log.Printf("PuppetHostie RoomQRCode(%s)\n", roomId)
  response, err := p.grpcClient.RoomQRCode(ctx, &pbwechaty.RoomQRCodeRequest{
    Id: roomId,
  })
  if err != nil {
    return "", err
  }
  return response.Qrcode, nil
}

func (p *PuppetHostie) RoomMemberList(ctx context.Context, roomId string) ([]string, error) {
  log.Printf("PuppetHostie RoomMemberList(%s)\n", roomId)
  response, err := p.grpcClient.RoomMemberList(ctx, &pbwechaty.RoomMemberListRequest{
    Id: roomId,
  })
  if err != nil {
    return nil, err
  }
  return response.MemberIds, nil
}

func (p *PuppetHostie) RoomMemberPayload(ctx context.Context, roomId string, contactId string) (*schemas.RoomMemberPayload, error) {
  log.Printf("PuppetHostie RoomMemberPayload(%s, %s)\n", roomId, contactId)
  response, err := p.grpcClient.RoomMemberPayload(ctx, &pbwechaty.RoomMemberPayloadRequest{
    Id:       roomId,
    MemberId: contactId,
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

func (p *PuppetHostie) SetRoomAnnounce(ctx context.Context, roomId, text string) error {
  log.Printf("PuppetHostie SetRoomAnnounce(%s, %s)\n", roomId, text)
  _, err := p.grpcClient.RoomAnnounce(ctx, &pbwechaty.RoomAnnounceRequest{
    Id:   roomId,
    Text: &wrappers.StringValue{Value: text},
  })
  if err != nil {
    return err
  }
  return nil
}

func (p *PuppetHostie) GetRoomAnnounce(ctx context.Context, roomId string) (string, error) {
  log.Printf("PuppetHostie GetRoomAnnounce(%s)\n", roomId)
  response, err := p.grpcClient.RoomAnnounce(ctx, &pbwechaty.RoomAnnounceRequest{
    Id: roomId,
  })
  if err != nil {
    return "", err
  }
  return response.Text.Value, nil
}

func (p *PuppetHostie) RoomInvitationAccept(ctx context.Context, roomInvitationId string) error {
  log.Printf("PuppetHostie RoomInvitationAccept(%s)\n", roomInvitationId)
  _, err := p.grpcClient.RoomInvitationAccept(ctx, &pbwechaty.RoomInvitationAcceptRequest{
    Id: roomInvitationId,
  })
  if err != nil {
    return err
  }
  return nil
}

func (p *PuppetHostie) RoomInvitationPayload(ctx context.Context, id string) (*schemas.RoomInvitationPayload, error) {
  log.Printf("PuppetHostie RoomInvitationPayload(%s)\n", id)
  response, err := p.grpcClient.RoomInvitationPayload(ctx, &pbwechaty.RoomInvitationPayloadRequest{
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

func (p *PuppetHostie) FriendshipSearchPhone(ctx context.Context, phone string) (string, error) {
  log.Printf("PuppetHostie FriendshipSearchPhone(%s)\n", phone)
  response, err := p.grpcClient.FriendshipSearchPhone(ctx, &pbwechaty.FriendshipSearchPhoneRequest{
    Phone: phone,
  })
  if err != nil {
    return "", err
  }
  return response.ContactId.Value, nil
}

func (p *PuppetHostie) FriendshipSearchWeixin(ctx context.Context, weixin string) (string, error) {
  log.Printf("PuppetHostie FriendshipSearchWeixin(%s)\n", weixin)
  response, err := p.grpcClient.FriendshipSearchWeixin(ctx, &pbwechaty.FriendshipSearchWeixinRequest{
    Weixin: weixin,
  })
  if err != nil {
    return "", err
  }
  return response.ContactId.Value, nil
}

func (p *PuppetHostie) FriendshipPayload(ctx context.Context, id string) (*schemas.FriendShipPayload, error) {
  log.Printf("PuppetHostie FriendshipPayload(%s)\n", id)
  response, err := p.grpcClient.FriendshipPayload(ctx, &pbwechaty.FriendshipPayloadRequest{
    Id: id,
  })
  if err != nil {
    return nil, err
  }
  return &schemas.FriendShipPayload{
    FriendshipPayloadReceive: schemas.FriendshipPayloadReceive{
      FriendshipPayloadBase: schemas.FriendshipPayloadBase{},
      Type:                  schemas.FriendshipType(response.Type),
      Scene:                 schemas.FriendshipSceneType(response.Scene),
      Stranger:              response.Stranger,
      Ticket:                response.Ticket,
    },
  }, nil
}

func (p *PuppetHostie) FriendshipAdd(ctx context.Context, contactId, hello string) (err error) {
  log.Printf("PuppetHostie FriendshipAdd(%s, %s)\n", contactId, hello)
  _, err = p.grpcClient.FriendshipAdd(ctx, &pbwechaty.FriendshipAddRequest{
    ContactId: contactId,
    Hello:     hello,
  })
  return err
}

func (p *PuppetHostie) FriendshipAccept(ctx context.Context, friendshipId string) (err error) {
  log.Printf("PuppetHostie FriendshipAccept(%s)\n", friendshipId)
  _, err = p.grpcClient.FrendshipAccept(ctx, &pbwechaty.FriendshipAcceptRequest{
    Id: friendshipId,
  })
  return err
}

func (p *PuppetHostie) TagContactAdd(ctx context.Context, id, contactId string) (err error) {
  log.Printf("PuppetHostie TagContactAdd(%s, %s)\n", id, contactId)
  _, err = p.grpcClient.TagContactAdd(ctx, &pbwechaty.TagContactAddRequest{
    Id:        id,
    ContactId: id,
  })
  return err
}

func (p *PuppetHostie) TagContactRemove(ctx context.Context, id, contactId string) (err error) {
  log.Printf("PuppetHostie TagContactRemove(%s, %s)\n", id, contactId)
  _, err = p.grpcClient.TagContactRemove(ctx, &pbwechaty.TagContactRemoveRequest{
    Id:        id,
    ContactId: id,
  })
  return err
}

func (p *PuppetHostie) TagContactDelete(ctx context.Context, id string) (err error) {
  log.Printf("PuppetHostie TagContactDelete(%s)\n", id)
  _, err = p.grpcClient.TagContactDelete(ctx, &pbwechaty.TagContactDeleteRequest{
    Id: id,
  })
  return err
}

func (p *PuppetHostie) TagContactList(ctx context.Context, contactId string) ([]string, error) {
  log.Printf("PuppetHostie TagContactList(%s)\n", contactId)
  response, err := p.grpcClient.TagContactList(ctx, &pbwechaty.TagContactListRequest{
    ContactId: &wrappers.StringValue{Value: contactId},
  })
  if err != nil {
    return nil, err
  }
  return response.Ids, nil
}
