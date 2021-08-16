package robot

import (
	"encoding/json"
	"fmt"
	"github.com/lyf571321556/chat-robot-api/api"
	"github.com/lyf571321556/chat-robot-api/file"
	"github.com/lyf571321556/chat-robot-api/image"
	"github.com/lyf571321556/chat-robot-api/markdown"
	"github.com/lyf571321556/chat-robot-api/news"
	"github.com/lyf571321556/chat-robot-api/text"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	// FeishuGroupBotSendUrl 企业微信群机器人 webhook
	FeishuGroupRobotSendurl = "https://open.feishu.cn/open-apis/bot/v2/hook/?key=%s"
	FeishuAppRobotSendURL   = "https://qyapi.weixin.qq.com/cgi-bin/message/send"
	// UploadMediaUrl 企业微信上传文件接口 `url`, `type` 固定传 `file`
	FeishuUploadMediaUrl = "https://qyapi.weixin.qq.com/cgi-bin/webhook/upload_media?key=%s&type=file"
)

type FeishuRobot struct {
	Name    string
	webhook string
	key     string
}

func NewFeishuRobot(name string, key string) *FeishuRobot {
	bot := new(FeishuRobot)
	bot.webhook = fmt.Sprintf(QiyeWeinxingGroupRobotSendurl, key)
	bot.Name = name
	bot.key = key
	return bot
}

func (b *FeishuRobot) PushTextMessage(content string, opts ...text.TextMsgOption) (err error) {
	msg := text.NewTextMsg(content, opts...)
	return b.pushMsg(msg)
}

func (b *FeishuRobot) PushMarkdownMessage(content string) (err error) {
	msg := markdown.NewMarkdownMsg(content)
	return b.pushMsg(msg)
}

func (b *FeishuRobot) PushImageMessage(img []byte) (err error) {
	msg := image.NewImageMsg(img)
	return b.pushMsg(msg)
}

func (b *FeishuRobot) PushNewsMessage(art news.Article, articles ...news.Article) (err error) {
	msg := news.NewNewsMsg(art, articles...)
	return b.pushMsg(msg)
}

func (b *FeishuRobot) PushFileMessage(media file.Media) error {
	msg := file.NewFileMsg(media.Id)
	return b.pushMsg(msg)
}

func handleFeishuResp(resp *http.Response) (err error) {
	var reply = new(struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	})

	rawResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(rawResp, reply); err != nil {
		return fmt.Errorf("unknown response: %w\nraw response: %s", err, rawResp)
	}
	if reply.ErrMsg != "ok" {
		return fmt.Errorf("unknown response: %s", rawResp)
	}
	return nil
}

func (b *FeishuRobot) pushMsg(msg interface{}) (err error) {
	var bsJSON []byte
	if bsJSON, err = json.Marshal(msg); err != nil {
		return err
	}
	var req *http.Request
	if req, err = api.NewRequest(http.MethodPost, b.webhook, bsJSON); err != nil {
		return err
	}

	rawResp, err := api.ExecuteHTTP(req, handleFeishuResp)
	if err != nil {
		return err
	}
	if rawResp != nil {
		log.Printf("result:%s", string(rawResp))
	}
	//var reply = new(struct {
	//	ErrCode   int    `json:"errcode"`
	//	ErrMsg    string `json:"errmsg"`
	//	Type      string `json:"type"`
	//	CreatedAt string `json:"created_at"`
	//})
	//if err = json.Unmarshal(rawResp, reply); err != nil {
	//	return fmt.Errorf("unknown response: %w\nraw response: %s", err, rawResp)
	//}
	//print(reply)
	return
}

func (b *FeishuRobot) UploadFile(filename string) (media file.Media, err error) {
	var req *http.Request
	if req, err = api.NewUploadRequest(http.MethodPost, fmt.Sprintf(QiyeWeinxingUploadMediaUrl, b.key), filename); err != nil {
		return file.Media{}, err
	}
	var rawResp []byte = nil
	if rawResp, err = api.ExecuteHTTP(req, handleFeishuResp); err != nil {
		return file.Media{}, err
	}

	var reply = new(struct {
		ErrCode   int    `json:"errcode"`
		ErrMsg    string `json:"errmsg"`
		Type      string `json:"type"`
		MediaId   string `json:"media_id"`
		CreatedAt string `json:"created_at"`
	})
	if err = json.Unmarshal(rawResp, reply); err != nil {
		return file.Media{}, fmt.Errorf("unknown response: %w\nraw response: %s", err, rawResp)
	}
	media = file.Media{Id: reply.MediaId}
	return
}
