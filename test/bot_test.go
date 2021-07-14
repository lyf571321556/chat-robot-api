package test

import (
	"github.com/lyf571321556/qiye-wechat-bot-api/api"
	"github.com/lyf571321556/qiye-wechat-bot-api/text"
	"testing"
)
import "github.com/lyf571321556/qiye-wechat-bot-api/bot"

var botKey = "3889d32e-5511-423c-a18b-12feeab07bab"
var phoneNumber = ""
var userid = "WuXingJuan||ChuYang"
var wechatbot api.QiyeWechatBot

func setup() {
	wechatbot = bot.NewQiyeWechatBot(botKey)
}

func Test_PushTextMessage(t *testing.T) {
	setup()

	err := wechatbot.PushTextMessage(
		"test", text.MentionByUserid(userid),
	)
	if err != nil {
		t.Fatal(err)
	}
}
