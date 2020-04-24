package wechaty_puppet_padplus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"

	pd "github.com/wechaty/go-wechaty/wechaty-puppet-padplus/proto"
	padschemas "github.com/wechaty/go-wechaty/wechaty-puppet-padplus/schemas"
	file_box "github.com/wechaty/go-wechaty/wechaty-puppet/file-box"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"

	option2 "github.com/wechaty/go-wechaty/wechaty-puppet/option"
)

const Endpoint = "padplus.juzibot.com:50051"

// PuppetPadPlus struct
type PuppetPadPlus struct {
	option      *option2.Option
	grpcConn    *grpc.ClientConn
	grpcClient  pd.PadPlusServerClient
	eventStream pd.PadPlusServer_InitClient
	Uin         string
}

// NewPuppetPadPlus new PuppetHostie struct
func NewPuppetPadPlus(o *option2.Option) *PuppetPadPlus {
	return &PuppetPadPlus{
		option: o,
	}
}

// FriendshipPayload ...
func (p *PuppetPadPlus) FriendshipPayload(id string) (*schemas.FriendshipPayload, error) {
	return nil, nil
}

// MessageImage ...
func (p *PuppetPadPlus) MessageImage(messageID string, imageType schemas.ImageType) (*file_box.FileBox, error) {
	return nil, nil
}

// FriendshipAccept ...
func (p *PuppetPadPlus) FriendshipAccept(friendshipID string) (err error) {
	return err
}

// RoomInvitationPayload ...
func (p *PuppetPadPlus) RoomInvitationPayload(id string) (*schemas.RoomInvitationPayload, error) {
	return nil, nil
}

// RoomInvitationAccept ...
func (p *PuppetPadPlus) RoomInvitationAccept(roomInvitationID string) error {
	return nil
}

// MessageSendText ...
func (p *PuppetPadPlus) MessageSendText(conversationID string, text string) (string, error) {
	return "", nil
}

// MessageSendContact ...
func (p *PuppetPadPlus) MessageSendContact(conversationID string, contactID string) (string, error) {
	return "", nil
}

// MessageSendURL ...
func (p *PuppetPadPlus) MessageSendURL(conversationID string, urlLinkPayload *schemas.UrlLinkPayload) (string, error) {
	return "", nil
}

// MessageSendFile ...
func (p *PuppetPadPlus) MessageSendFile(conversationID string, fileBox *file_box.FileBox) (string, error) {
	return "", nil
}

// MessageSendMiniProgram ...
func (p *PuppetPadPlus) MessageSendMiniProgram(conversationID string, miniProgramPayload *schemas.MiniProgramPayload) (string, error) {
	return "", nil
}

// Start ...
func (p *PuppetPadPlus) Start() (err error) {
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

	if p.isLogin() {
		err = p.AutoLogin()
		if err != nil {
			return err
		}
	} else {
		err = p.Login()
		if err != nil {
			return err
		}
	}
	return nil
}

// Login 登录
func (p *PuppetPadPlus) Login() (err error) {
	err, _ = p.Request(pd.ApiType_GET_QRCODE, nil)
	return
}

// AutoLogin 自动登录
func (p *PuppetPadPlus) AutoLogin() (err error) {
	err, _ = p.Request(pd.ApiType_INIT, nil)
	return
}

// Stop ...
func (p *PuppetPadPlus) Stop() {
	var err error
	defer func() {
		if err != nil {
			log.Printf("PuppetHostie Stop err: %s\n", err)
		}
	}()
	if p.isLogin() {
		p.option.Emit(schemas.EventLogoutPayload{
			ContactId: p.Uin,
			Data:      "PuppetHostie Stop()",
		})
		p.Uin = ""
	}

	if err = p.stopGrpcStream(); err != nil {
		return
	}
	if err = p.stopGrpcClient(); err != nil {
		return
	}
}

// startGrpcClient start GRPC Client
func (p *PuppetPadPlus) startGrpcClient() error {
	endpoint := p.option.Endpoint
	if endpoint == "" {
		endpoint = Endpoint
	}
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	p.grpcConn = conn
	p.grpcClient = pd.NewPadPlusServerClient(conn)
	return nil
}

// startGrpcStream start GRPC Stream
func (p *PuppetPadPlus) startGrpcStream() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("startGrpcStream err:%w", err)
		}
	}()
	if p.eventStream != nil {
		return errors.New("event stream exists")
	}
	p.eventStream, err = p.grpcClient.Init(context.Background(), &pd.InitConfig{
		Token: &p.option.Token,
	})
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

func (p *PuppetPadPlus) stopGrpcClient() error {
	if p.grpcConn == nil {
		return errors.New("puppetClient had not inited")
	}
	p.grpcConn.Close()
	p.grpcConn = nil
	p.grpcClient = nil
	return nil
}

// stopGrpcStream stop GRPC Stream
func (p *PuppetPadPlus) stopGrpcStream() error {
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

// Request 发送请求
// 有些事件应该在一定时间后得到响应，否则则超市
// SEND_MESSAGE,SEND_FILE 3 * 60 * 1000
// GET_MESSAGE_MEDIA 5 * 60 * 1000
// SEARCH_CONTACT,ADD_CONTACT,CREATE_ROOM,GET_ROOM_QRCODE 1 * 60 * 1000
func (p *PuppetPadPlus) Request(apiType pd.ApiType, data interface{}) (err error, response string) {
	us := uuid.NewV4().String()
	trace := uuid.NewV4().String()
	paramsByte, err := json.Marshal(data)
	params := string(paramsByte)

	resp, err := p.grpcClient.Request(context.Background(), &pd.RequestObject{
		RequestId: &us, // uuid
		Token:     &p.option.Token,
		ApiType:   &apiType,
		TraceId:   &trace, // uuid
		Params:    &params,
		Uin:       &p.Uin,
	})
	if err != nil {
		log.Printf("Type: %s, err: %v", apiType, err)
		return err, ""
	}
	log.Printf("Type: %s, Result: %s", apiType, *resp.Result)
	return err, *resp.Result
}

// isLogin is login
func (p *PuppetPadPlus) isLogin() bool {
	return p.Uin != ""
}

func (p *PuppetPadPlus) onGrpcStreamEvent(resp *pd.StreamResponse) {
	//log.Printf("resp: uin: %s, data: %s, requestId: %s, responseType: %s, traceId: %s", *resp.Uin, *resp.Data, *resp.RequestId, *resp.ResponseType, *resp.TraceId)
	if *resp.Data == "EXPIRED_TOKEN" || *resp.Data == "INVALID_TOKEN" {

	}

	//log.Println("resp: ", resp.String())
	switch *resp.ResponseType {
	case pd.ResponseType_LOGIN_QRCODE: // 登录二维码
		var data padschemas.PadPlusQrCode
		err := json.Unmarshal([]byte(*resp.Data), &data)
		if err != nil {
			return
		}
		scan := schemas.EventScanPayload{
			BaseEventPayload: schemas.BaseEventPayload{Data: *resp.Data},
			QrCode:           "http://weixin.qq.com/x/" + data.QrCodeId,
		}
		p.option.Emit(schemas.PuppetEventNameScan, &scan)
		break
	case pd.ResponseType_ACCOUNT_LOGOUT: // 退出登录
		var data padschemas.LogoutGRPCResponse
		err := json.Unmarshal([]byte(*resp.Data), &data)
		if err != nil {
			return
		}
		break
	case pd.ResponseType_MESSAGE_RECEIVE: // 收到消息
		//err = w.onMessage(*resp.Data)
		break
	case pd.ResponseType_MESSAGE_MEDIA_SRC: // 收到媒资源信息
		break
	case pd.ResponseType_QRCODE_SCAN: // 扫描二维码
	case pd.ResponseType_QRCODE_LOGIN: // 登录二维码
		var data padschemas.GRPCQrCodeLogin
		err := json.Unmarshal([]byte(*resp.Data), &data)
		if err != nil {
			return
		}
		p.option.Emit(schemas.PuppetEventNameScan, data)
		break
	case pd.ResponseType_CONTACT_LIST, pd.ResponseType_CONTACT_MODIFY: // 通讯录列表
		// 提取UserName
		var re = regexp.MustCompile(`(?U)"UserName":"(.*)"`)
		userName := re.FindString(*resp.Data)
		if len(userName) == 0 {
			var us map[string]interface{}
			err := json.Unmarshal([]byte(*resp.Data), &us)
			if err == nil {
				userName = us["UserName"].(string)
			}
		}
		if isRoomId(userName) {
			var contact padschemas.GRPCContactPayload
			err := json.Unmarshal([]byte(*resp.Data), &contact)
			if err != nil {
				return
			}
		} else {
			var contact padschemas.GRPCRoomPayload
			err := json.Unmarshal([]byte(*resp.Data), &contact)
			if err != nil {
				return
			}
		}
	}
	return
}

func (p *PuppetPadPlus) onMessage(data string) (err error) {
	var payload padschemas.GRPCMessagePayload
	err = json.Unmarshal([]byte(data), &payload)
	if err != nil {
		return
	}
	switch payload.MsgType {
	case padschemas.WechatMessageTypeImage:
		log.Println("获取图片媒资信息")
		err, response := p.Request(pd.ApiType_GET_MESSAGE_MEDIA, padschemas.PadPlusRichMediaData{
			Content:      "",
			MsgType:      payload.MsgType,
			ContentType:  "img",
			Src:          payload.Url2,
			AppMsgType:   0,
			FileName:     payload.FileName2,
			MsgId:        payload.MsgId,
			CreateTime:   payload.CreateTime,
			FromUserName: payload.FromUserName,
			ToUserName:   payload.ToUserName,
		})
		if err != nil {
			return err
		}
		var resp padschemas.PadPlusMediaData
		_ = json.Unmarshal([]byte(response), &resp)
		log.Println(resp)
	}
	return
}

func isRoomId(s string) bool {
	return strings.HasSuffix(s, "@chatroom")
}
