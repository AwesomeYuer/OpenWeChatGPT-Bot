package handlers

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/869413421/wechatbot/config"
	"github.com/869413421/wechatbot/service"
	"github.com/eatmoreapple/openwechat"
	"github.com/skip2/go-qrcode"
)

// MessageHandlerInterface 消息处理接口
type MessageHandlerInterface interface {
	handle(*openwechat.Message) error
	ReplyText(*openwechat.Message) error
}

type HandlerType string

const (
	GroupHandler = "group"
	UserHandler  = "user"
)

// QrCodeCallBack 登录扫码回调，
func QrCodeCallBack(uuid string) {
	if 1 == 0 && runtime.GOOS == "windows" {
		// 运行在Windows系统上
		openwechat.PrintlnQrcodeUrl(uuid)
	} else {
		log.Println("login in linux/windows")
		q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
		fmt.Println(q.ToString(true))
	}
}

// handlers 所有消息类型类型的处理器
var handlers map[HandlerType]MessageHandlerInterface
var UserService service.UserServiceInterface

func init() {
	handlers = make(map[HandlerType]MessageHandlerInterface)
	handlers[GroupHandler] = NewGroupMessageHandler()
	handlers[UserHandler] = NewUserMessageHandler()
	UserService = service.NewUserService()
}

// Handler 全局处理入口
func Handler(msg *openwechat.Message) {
	log.Printf("hadler Received msg : %v", msg.Content)
	// 处理群消息
	if msg.IsSendByGroup() {
		handlers[GroupHandler].handle(msg)
		return
	}

	// 好友申请
	if msg.IsFriendAdd() {
		log.Printf("hadler Received msg IsFriendAdd request:\r\n%v,\r\n%v", msg.FromUserName, msg.Content)
		if strings.Contains(msg.Content, " content=\"awesomeyuer\" ") {
			if config.LoadConfig().AutoPass {
				_, err := msg.Agree("你好我是基于chatGPT引擎开发的微信机器人，你可以向我提问任何问题。")
				if err != nil {
					log.Fatalf("add friend agree error : %v", err)
					return
				} else {
					log.Printf("accept new friend : %v", msg.FromUserName)
				}
			}
		} else {
			log.Printf("deny new friend : %v", msg.FromUserName)
		}
	}

	// 私聊
	handlers[UserHandler].handle(msg)
}
