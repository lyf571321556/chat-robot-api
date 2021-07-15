package bot

import (
	"encoding/json"
	"fmt"
	"github.com/lyf571321556/qiye-wechat-bot-api/api"
	"github.com/lyf571321556/qiye-wechat-bot-api/file"
	"github.com/lyf571321556/qiye-wechat-bot-api/image"
	"github.com/lyf571321556/qiye-wechat-bot-api/markdown"
	"github.com/lyf571321556/qiye-wechat-bot-api/news"
	"github.com/lyf571321556/qiye-wechat-bot-api/text"
	"io/ioutil"
	"log"
	"net/http"
)

type QiyeWechatBot struct {
	webhook string
	key     string
}

func NewQiyeWechatBot(key string) *QiyeWechatBot {
	bot := new(QiyeWechatBot)
	bot.webhook = fmt.Sprintf(api.GroupBotSendUrl, key)

	bot.key = key
	return bot
}

func (b *QiyeWechatBot) PushTextMessage(content string, opts ...text.TextMsgOption) (err error) {
	msg := text.NewTextMsg(content, opts...)
	return b.pushMsg(msg)
}

func (b *QiyeWechatBot) PushMarkdownMessage(content string) (err error) {
	msg := markdown.NewMarkdownMsg(content)
	return b.pushMsg(msg)
}

func (b *QiyeWechatBot) PushImageMessage(img []byte) (err error) {
	msg := image.NewImageMsg(img)
	return b.pushMsg(msg)
}

func (b *QiyeWechatBot) PushNewsMessage(art news.Article, articles ...news.Article) (err error) {
	msg := news.NewNewsMsg(art, articles...)
	return b.pushMsg(msg)
}

func (b *QiyeWechatBot) PushFileMessage(media file.Media) error {
	msg := file.NewFileMsg(media.Id)
	return b.pushMsg(msg)
}

func handleResp(resp *http.Response) (err error) {
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

func (b *QiyeWechatBot) pushMsg(msg interface{}) (err error) {
	var bsJSON []byte
	if bsJSON, err = json.Marshal(msg); err != nil {
		return err
	}
	var req *http.Request
	if req, err = api.NewRequest(http.MethodPost, b.webhook, bsJSON); err != nil {
		return err
	}

	rawResp, err := api.ExecuteHTTP(req, handleResp)
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

func (b *QiyeWechatBot) UploadFile(filename string) (media file.Media, err error) {
	var req *http.Request
	if req, err = api.NewUploadRequest(http.MethodPost, fmt.Sprintf(api.UploadMediaUrl, b.key), filename); err != nil {
		return file.Media{}, err
	}
	var rawResp []byte = nil
	if rawResp, err = api.ExecuteHTTP(req, handleResp); err != nil {
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
