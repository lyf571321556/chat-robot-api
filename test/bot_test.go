package test

import (
	"github.com/lyf571321556/chat-robot-api/api"
	"github.com/lyf571321556/chat-robot-api/qiye_wechat"
	"testing"
)
import "github.com/lyf571321556/chat-robot-api/robot"

var botKey = "8c8c33c5-2ce2-4fae-b610-be137741a06d"
var phoneNumber = ""
var userid = "WuXingJuan"
var wechatbot api.Robot

func setup() {
	wechatbot = robot.NewQiyeWechatRobot("", botKey)
	api.SetDebug(true)
}

func Test_PushTextMessage(t *testing.T) {
	setup()

	err := wechatbot.PushTextMessage(
		"test", qiye_wechat.MentionByUserid(userid), qiye_wechat.MentionByUserid("ChuYang"),
	)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_PushMarkdownMessage(t *testing.T) {
	setup()

	err := wechatbot.PushMarkdownMessage(
		qiye_wechat.Heading(1, "警告⚠️") + "工时登记提醒," + qiye_wechat.WarningText("") + "，请相关同事注意。\n" +
			qiye_wechat.QuoteText("@ChuYang :"+qiye_wechat.CommentText("test001")) +
			qiye_wechat.QuoteText("@WuXingJuan :"+qiye_wechat.CommentText("test002")) +
			qiye_wechat.QuoteText("@all :"+qiye_wechat.CommentText("test003")),
	)
	if err != nil {
		t.Fatal(err)
	}
}
