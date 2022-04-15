/*
README:
-----------------------
config.yaml >>>
bot:
  adminid: wxid_xxxxxx
  chat: "on"
  name: BotName # 自动获取
ding: # 钉钉 webhock
  keyword: Wechaty # 关键字
  token: xxxxxxxxxxxxxxx
  url: https://oapi.dingtalk.com/robot/send?access_token=
tuling:
  token: xxxxxxxxxxxxxxxxxxxx&info= # 后面这段 `&info=` 属于懒人操作
  url: http://www.tuling123.com/openapi/api?key=
wechaty:
  wechaty_puppet_endpoint: 127.0.0.1:25000
  wechaty_puppet_service_token: insecure_xxxxxxxxxxx # 现在需要使用这种 token
config.yaml <<<
-----------------------
puppet-xp
目测 好友和群聊操作不了，但是确实最稳定的 消息管理 工具
-----------------------
Viper 组件随缘，如果你不想把配置都写进源码里的话，或者使用更大内存的 数据库
-----------------------
dingding 消息推送，主要是为了提示你有人@过你 ，以及 微信账号 退出通知，十有八九封号了
-----------------------
tuling
http://www.turingapi.com/ 免费版 有100条消息每天，也可以接入微信的 对话开放平台
-----------------------
如果有兴趣可以与我一起管理 : https://github.com/XRSec/Go-wechaty-Bot

谢谢！
*/
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty"
	wp "github.com/wechaty/go-wechaty/wechaty-puppet"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	_interface "github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
)

type (
	TulingBotResult struct {
		Code int    `json:"code"`
		Text string `json:"text"`
	}
	DingBotResult struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}
)

var (
	err             error
	resp            *http.Response
	dingBotResult   DingBotResult
	tulingBotResult TulingBotResult
)

func init() {
	// 设置日志格式
	log.SetPrefix("[xrsec] [\033[01;33m➜\033[0m] ") // 设置日志前缀
	log.SetFlags(log.Ltime | log.Lshortfile)

	// 初始化配置文件
	rootPath, _ := os.Getwd()     // 当前用户路劲
	exePath, _ := os.Executable() // 当前 程序路径
	log.Printf("rootPath: %s\nexePath: %s", rootPath, filepath.Dir(exePath))

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Dir(exePath))
	viper.AddConfigPath(rootPath)
	viper.Set("rootPath", rootPath)
	viper.Set("exePath", exePath)
}

func viperReload() { // 重新读取 viper 配置文件，方便修改
	if err = viper.ReadInConfig(); err != nil { // 永远不能相信用户的输入，所以要做好判断和提示
		log.Printf("Viper Read Config Error, Error: %s", err)
		if _, err = os.Stat(viper.GetString("rootPath") + "/config.yaml"); err != nil {
			if _, err = os.Stat(viper.GetString("exePath") + "/config.yaml"); err != nil {
				log.Println("配置文件放在当前路劲即可, 注意检测配置是否正确")
			}
			log.Printf("config.yaml not found, Error: %s", err)
			viper.Set("wechaty.wechaty_puppet_endpoint", "Please Fill In Your Server Address")
			viper.Set("wechaty.wechaty_puppet_service_token", "Please Fill In Your Token")
			var f *os.File
			if f, err = os.Create(viper.GetString("exePath") + "/config.yaml"); err != nil {
				log.Printf("Create Config File, Error: %s", err)
			} else {
				log.Printf("请修改你的配置文件: %s/config.yaml", viper.GetString("rootPath"))
			}
			defer func(f *os.File) {
				if err = f.Close(); err != nil {
					log.Printf("Close Config File, Error: %s", err)
				}
			}(f)
		}
	}
}

func ViperWrite() {
	if err = viper.WriteConfigAs(viper.ConfigFileUsed()); err != nil {
		log.Printf("Viper Write file Error: %s", err)
	} else {
		log.Println("Viper Write file Success")
	}
}

func NightMode() bool {
	//当前时间
	startTimeStr := "00:00:00"
	endTimeStr := "06:00:00"
	now := time.Now()
	//当前时间转换为"年-月-日"的格式
	format := now.Format("2006-01-02")
	//转换为time类型需要的格式
	layout := "2006-01-02 15:04:05"
	//将开始时间拼接“年-月-日 ”转换为time类型
	timeStart, _ := time.ParseInLocation(layout, format+" "+startTimeStr, time.Local)
	//将结束时间拼接“年-月-日 ”转换为time类型
	timeEnd, _ := time.ParseInLocation(layout, format+" "+endTimeStr, time.Local)
	//使用time的Before和After方法，判断当前时间是否在参数的时间范围
	return now.Before(timeEnd) && now.After(timeStart)
}

func DingBotCheck() {
	if viper.GetString("Ding.URl") == "" {
		log.Printf("DingDing, Error: %s", errors.New("机器人URL为空"))
	} else {
		// 单个机器人有单位时间内消息条数的限制，如果有需要可以初始化多个token，发消息时随机发给其中一个机器人。
		log.Println("已设置钉钉提醒")
	}
}

func DingMessage(message string) {
	if NightMode() {
		log.Println("现在处于夜间模式，请在白天使用")
		return
	} else {
		dingWebHook := viper.GetString("Ding.URL") + viper.GetString("Ding.TOKEN")
		content := fmt.Sprintf(" {\"msgtype\": \"text\",\"text\": {\"content\": \"%s %s\"}}", viper.GetString("Ding.KEYWORD"), message)
		// 发送请求
		if resp, err = http.Post(dingWebHook, "application/json; charset=utf-8", strings.NewReader(content)); err != nil {
			log.Printf("机器人请求错误, Error: %s", err)
		} else {
			if err = json.NewDecoder(resp.Body).Decode(&dingBotResult); err != nil {
				log.Printf("机器人请求错误, Error: %s", err)
			} else {
				if dingBotResult.Errcode == 0 {
					log.Println("消息发送成功!")
				} else {
					log.Printf("消息发送失败, Error: %s", err)
				}
			}
		}
		// 关闭请求
		defer func(Body io.ReadCloser) {
			if err = Body.Close(); err != nil {
				log.Printf("关闭请求错误, Error: %s", err)
			}
		}(resp.Body)
	}
}

func TulingMessage(msg string) string {
	if NightMode() {
		log.Println("现在处于夜间模式，请在白天使用")
		return ""
	} else {
		// 发送请求
		tulingWebhook := viper.GetString("Tuling.URL") + viper.GetString("Tuling.TOKEN")
		if resp, err = http.Get(tulingWebhook + msg); err != nil {
			log.Printf("图灵机器人请求错误, Error: %s", err)
		} else {
			if err = json.NewDecoder(resp.Body).Decode(&tulingBotResult); err != nil {
				return ""
			} else {
				if tulingBotResult.Code != 100000 {
					return ""
				} else {
					log.Printf("图灵机器人 回复信息: %+v", tulingBotResult.Text)
					return tulingBotResult.Text
				}
			}
		}
		// 关闭请求
		defer func(Body io.ReadCloser) {
			if err = Body.Close(); err != nil {
				log.Println("Close body error:", err)
			}
		}(resp.Body)
		return ""
	}
}

/*
	--------------------------------------------------------
*/

func onScan(context *wechaty.Context, qrCode string, status schemas.ScanStatus, data string) {
	log.Printf("%s[Scan] %s %s %s\n", viper.GetString("info"), qrCode, status, data)
}

/*
	@method onLogin 当机器人成功登陆后，会触发事件，并会在事件中传递当前登陆机器人的信息
	@param {*} user
*/
func onLogin(context *wechaty.Context, user *user.ContactSelf) {
	log.Printf(`
                           //
               \\         //
                \\       //
        ##DDDDDDDDDDDDDDDDDDDDDD##
        ## DDDDDDDDDDDDDDDDDDDD ##      
        ## DDDDDDDDDDDDDDDDDDDD ##      
        ## hh                hh ##      ##         ## ## ## ##   ## ## ## ###   ##    ####     ##     
        ## hh    //    \\    hh ##      ##         ##       ##   ##             ##    ## ##    ##
        ## hh   //      \\   hh ##      ##         ##       ##   ##             ##    ##   ##  ##
        ## hh                hh ##      ##         ##       ##   ##     ##      ##    ##    ## ##
        ## hh      wwww      hh ##      ##         ##       ##   ##       ##    ##    ##     ####
        ## hh                hh ##      ## ## ##   ## ## ## ##   ## ## ## ###   ##    ##      ###
        ## MMMMMMMMMMMMMMMMMMMM ##    
        ##MMMMMMMMMMMMMMMMMMMMMM##      微信机器人: [%s] 已经登录成功了。
        `, user.Name())
	viper.Set("bot.name", user.Name())
}

/**
@method onLogout 当机器人检测到登出的时候，会触发事件，并会在事件中传递机器人的信息。
@param {*} user
*/
func onLogout(context *wechaty.Context, user *user.ContactSelf, reason string) {
	log.Println("========================onLogout👇========================")
	DingMessage(user.Name() + "账号已退出登录, 请检查账号!" + reason)
}

/*
  @method onRoomInvite 当收到群邀请的时候，会触发这个事件。
  @param {*} user
*/
func onRoomInvite(context *wechaty.Context, roomInvitation *user.RoomInvitation) {
	log.Println("========================onRoomInvite👇========================")
	if err = roomInvitation.Accept(); err != nil {
		log.Printf("Accept Room Invitation, Error: %s", err)
		//	好像有点问题，群聊设置了邀请确认就用不了
	}
	log.Println(roomInvitation.String())
}

/*
	@method onRoomTopic 当有人修改群名称的时候会触发这个事件。
	@param {*} user
*/
func onRoomTopic(context *wechaty.Context, room *user.Room, newTopic string, oldTopic string, changer _interface.IContact, date time.Time) {
	log.Println("========================onRoomTopic👇========================")
}

/*
	进入房间监听回调 room-群聊 inviteeList-受邀者名单 inviter-邀请者
	判断配置项群组id数组中是否存在该群聊id
*/
func onRoomJoin(context *wechaty.Context, room *user.Room, inviteeList []_interface.IContact, inviter _interface.IContact, date time.Time) {
}

/*
	@method onRoomleave 当机器人把群里某个用户移出群聊的时候会触发这个时间。用户主动退群是无法检测到的。
	@param {*} user
*/
func onRoomleave(context *wechaty.Context, room *user.Room, leaverList []_interface.IContact, remover _interface.IContact, date time.Time) {
	log.Println("========================onRoomleave👇========================")
	log.Printf("用户[%s]被踢出去聊", remover.Name())
}

func onFriendship(context *wechaty.Context, friendship *user.Friendship) {
	log.Println("========================onFriendship👇========================")
	switch friendship.Type() {
	case 1:
	//FriendshipTypeUnknown
	case 2:
		//FriendshipTypeConfirm
		/**
		 * 2. 友谊确认
		 */
		log.Printf("friend ship confirmed with%s", friendship.Contact().Name())
	case 3:
		//FriendshipTypeReceive
		/*
			1. 新的好友请求
			设置请求后，我们可以从request.hello中获得验证消息,
			并通过`request.accept（）`接受此请求
		*/
		if friendship.Hello() == viper.GetString("addFriendKeywords") {
			if err = friendship.Accept(); err != nil {
				log.Printf("添加好友失败, Error: %s", err)
			}
		} else {
			log.Printf("%s未能自动通过好友申请, 因为验证消息是%s", friendship.Contact().Name(), friendship.Hello())
		}
	case 4:
	//FriendshipTypeVerify
	default:
	}
	log.Printf("%s好友关系是: %s", friendship.Contact().Name(), friendship.Type())
}

/*
	@method onHeartbeat 获取机器人的心跳。
	@param {*} user
*/
func onHeartbeat(context *wechaty.Context, data string) {
	log.Println("========================onHeartbeat👇========================")
	log.Printf("获取机器人的心跳: %s", data)
}

func OnMessage(context *wechaty.Context, message *user.Message) {
	if message.Self() {
		return
	}
	if message.Age() > 2*60*time.Second {
		log.Println("消息已丢弃，因为它太旧（超过2分钟）")
	}
	if message.Type() == schemas.MessageTypeText {
		if message.Room() != nil { // 群聊
			if message.MentionSelf() || strings.Contains(message.Text(), "@"+viper.GetString("bot.name")) { // @我 的我操作
				if reply := TulingMessage(strings.Replace(message.Text(), "@"+viper.GetString("bot.name"), "", 1)); reply != "" { // 获取图灵接口信息
					if _, err = message.Say(reply); err == nil { // 回复图灵返回的内容
						DingMessage(fmt.Sprintf("群聊名称: %s 用户名: %s 消息内容: %s", message.Room().String(), message.From().Name(), message.Text())) // 钉钉 推送
					}
				}
			} else {
				// 没有 @我 就老老实实的
			}
		} else { // 私聊
			if strings.Contains("加群", message.Text()) {
				// 邀请进群
			}
		}
	} else {
		//	 其他类型的消息
	}
	// 每条消息都在终端输出
	log.Printf("用户: [%s] 聊天内容:[%s]", message.From().Name(), message.Text())
}

func onError(context *wechaty.Context, err error) {
	log.Printf("机器人错误, Error: %s", err)
}

func main() {
	i := 0
	// 重试次数 10
	for i <= 10 {
		i++
		// 读取配置文件
		viperReload()
		// 钉钉推送
		DingBotCheck()
		var bot = wechaty.NewWechaty(wechaty.WithPuppetOption(wp.Option{
			Token:    viper.GetString("wechaty.wechaty_puppet_service_token"),
			Endpoint: viper.GetString("wechaty.wechaty_puppet_endpoint"),
		}))
		log.Printf("Token:%s", viper.GetString("wechaty.wechaty_puppet_service_token"))
		log.Printf("Endpoint: %s", viper.GetString("wechaty.wechaty_puppet_endpoint"))

		bot.OnScan(onScan).
			OnLogin(onLogin).
			OnLogout(onLogout).
			OnMessage(OnMessage).
			OnRoomInvite(onRoomInvite).
			OnRoomTopic(onRoomTopic).
			OnRoomJoin(onRoomJoin).
			OnRoomLeave(onRoomleave).
			OnFriendship(onFriendship).
			OnHeartbeat(onHeartbeat).
			OnError(onError)

		if err = bot.Start(); err != nil {
			// TODO 重启Bot，当 RPC error 的时候就不能成功重启
			log.Printf("Bot 错误, Error: %s", err)
			if i > 10 {
				os.Exit(0)
			}
			log.Printf("正在重新启动程序, 当前重试次数: 第%v次", i)
			time.Sleep(10 * time.Second)
		} else {
			i = 0
			// Bot 守护程序
			var quitSig = make(chan os.Signal)
			signal.Notify(quitSig, os.Interrupt, os.Kill)
			select {
			case <-quitSig:
				ViperWrite()
				log.Fatal("程序退出!")
			}
		}
	}
}
