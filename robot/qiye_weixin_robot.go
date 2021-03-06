package robot

import (
	"encoding/json"
	"fmt"
	"github.com/lyf571321556/chat-robot-api/api"
	"github.com/lyf571321556/chat-robot-api/qiye_wechat"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	// GroupBotSendUrl 企业微信群机器人 webhook
	QiyeWeinxingGroupRobotSendurl = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s"
	QiyeWeixingAppRobotSendURL    = "https://qyapi.weixin.qq.com/cgi-bin/message/send"
	// UploadMediaUrl 企业微信上传文件接口 `url`, `type` 固定传 `file`
	QiyeWeinxingUploadMediaUrl = "https://qyapi.weixin.qq.com/cgi-bin/webhook/upload_media?key=%s&type=file"
)

type QiyeWechatRobot struct {
	Name    string
	webhook string
	key     string
	*api.Robot
}

func NewQiyeWechatRobot(name string, key string) *QiyeWechatRobot {
	bot := new(QiyeWechatRobot)
	bot.webhook = fmt.Sprintf(QiyeWeinxingGroupRobotSendurl, key)
	bot.Name = name
	bot.key = key
	return bot
}

func (b *QiyeWechatRobot) PushTextMessage(content string, opts ...qiye_wechat.TextMsgOption) (err error) {
	msg := qiye_wechat.NewTextMsg(content, opts...)
	return b.pushMsg(msg)
}

func (b *QiyeWechatRobot) PushMarkdownMessage(content string) (err error) {
	msg := qiye_wechat.NewMarkdownMsg(content)
	return b.pushMsg(msg)
}

func (b *QiyeWechatRobot) PushImageMessage(img []byte) (err error) {
	msg := qiye_wechat.NewImageMsg(img)
	return b.pushMsg(msg)
}

func (b *QiyeWechatRobot) PushNewsMessage(art qiye_wechat.Article, articles ...qiye_wechat.Article) (err error) {
	msg := qiye_wechat.NewNewsMsg(art, articles...)
	return b.pushMsg(msg)
}

func (b *QiyeWechatRobot) PushFileMessage(media qiye_wechat.Media) error {
	msg := qiye_wechat.NewFileMsg(media.Id)
	return b.pushMsg(msg)
}

func handleQiyeWeixinResp(resp *http.Response) (err error) {
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

func (b *QiyeWechatRobot) pushMsg(msg interface{}) (err error) {
	var bsJSON []byte
	if bsJSON, err = json.Marshal(msg); err != nil {
		return err
	}
	var req *http.Request
	if req, err = api.NewRequest(http.MethodPost, b.webhook, bsJSON); err != nil {
		return err
	}

	rawResp, err := api.ExecuteHTTP(req, handleQiyeWeixinResp)
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

func (b *QiyeWechatRobot) UploadFile(filename string) (media qiye_wechat.Media, err error) {
	var req *http.Request
	if req, err = api.NewUploadRequest(http.MethodPost, fmt.Sprintf(QiyeWeinxingUploadMediaUrl, b.key), filename); err != nil {
		return qiye_wechat.Media{}, err
	}
	var rawResp []byte = nil
	if rawResp, err = api.ExecuteHTTP(req, handleQiyeWeixinResp); err != nil {
		return qiye_wechat.Media{}, err
	}

	var reply = new(struct {
		ErrCode   int    `json:"errcode"`
		ErrMsg    string `json:"errmsg"`
		Type      string `json:"type"`
		MediaId   string `json:"media_id"`
		CreatedAt string `json:"created_at"`
	})
	if err = json.Unmarshal(rawResp, reply); err != nil {
		return qiye_wechat.Media{}, fmt.Errorf("unknown response: %w\nraw response: %s", err, rawResp)
	}
	media = qiye_wechat.Media{Id: reply.MediaId}
	return
}
