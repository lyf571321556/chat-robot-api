package test

import (
	"github.com/lyf571321556/chat-robot-api/api"
	"github.com/lyf571321556/chat-robot-api/markdown"
	"github.com/lyf571321556/chat-robot-api/text"
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
		"test", text.MentionByUserid(userid), text.MentionByUserid("ChuYang"),
	)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_PushMarkdownMessage(t *testing.T) {
	setup()

	err := wechatbot.PushMarkdownMessage(
		markdown.Heading(1, "警告⚠️") + "工时登记提醒," + markdown.WarningText("") + "，请相关同事注意。\n" +
			markdown.QuoteText("@ChuYang :"+markdown.CommentText("test001")) +
			markdown.QuoteText("@WuXingJuan :"+markdown.CommentText("test002")) +
			markdown.QuoteText("@all :"+markdown.CommentText("test003")),
	)
	if err != nil {
		t.Fatal(err)
	}
}
