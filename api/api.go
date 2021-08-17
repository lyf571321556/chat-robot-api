package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path"
	"strings"
	"time"
)

var HTTPClient = http.DefaultClient

type MessageSentDelegate interface {
	SendMessage(msg interface{}, handleResp HandRespFunc) (err error)
	UploadFile(filePath string, handleResp HandRespFunc) (err error)
	WebhookUrl() (webhookUrl string, err error)
	UploadUrl() (uploadUrl string, err error)
}

type Robot struct {
	MessageSentDelegate *MessageSentDelegate
	//PushTextMessage(content string, opts ...qiye_wechat.TextMsgOption) error
	//PushMarkdownMessage(content string) error
	//PushImageMessage(img []byte) error
	//PushNewsMessage(art qiye_wechat.Article, articles ...qiye_wechat.Article) error
	//PushFileMessage(media qiye_wechat.Media) error
	//UploadFile(filename string) (media qiye_wechat.Media, err error)
}

func (r *Robot) SendMessage(msg interface{}, handleResp HandRespFunc) (err error) {
	var bsJSON []byte
	if bsJSON, err = json.Marshal(msg); err != nil {
		return err
	}
	var req *http.Request
	var webhookUrl string
	if webhookUrl, err = r.WebhookUrl(); err != nil {
		return err
	}
	if req, err = NewRequest(http.MethodPost, webhookUrl, bsJSON); err != nil {
		return err
	}

	rawResp, err := ExecuteHTTP(req, handleResp)
	if err != nil {
		return err
	}
	if rawResp != nil {
		log.Printf("result:%s", string(rawResp))
	}
	return
}

func (r *Robot) UploadFile(filePath string, handleResp HandRespFunc) (err error) {
	var req *http.Request
	var uploadUrl string
	if uploadUrl, err = r.UploadUrl(); err != nil {
		return err
	}
	if req, err = NewUploadRequest(http.MethodPost, uploadUrl, filePath); err != nil {
		return err
	}
	rawResp, err := ExecuteHTTP(req, handleResp)
	if err != nil {
		return err
	}
	if rawResp != nil {
		log.Printf("result:%s", string(rawResp))
	}
	return
}

func (r *Robot) WebhookUrl() (webhookUrl string, err error) {
	panic("must implement method WebhookUrl()")
}

func (r *Robot) UploadUrl() (uploadUrl string, err error) {
	panic("must implement method UploadUrl()")
}

func NewRequest(method string, rawUrl string, rawBody []byte) (request *http.Request, err error) {
	debugLog(fmt.Sprintf("start request --> %s %s<<<%s>>>", method, rawUrl, rawBody))
	if request, err = http.NewRequest(method, rawUrl, bytes.NewBuffer(rawBody)); err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	return
}

func NewUploadRequest(method string, rawUrl string, filename string) (request *http.Request, err error) {
	bodyBuffer := bytes.NewBufferString("")
	writer := multipart.NewWriter(bodyBuffer)

	mediaFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() { _ = mediaFile.Close() }()

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			"media", escapeQuotes(path.Base(filename))))
	h.Set("Content-Type", "application/octet-stream")
	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, err
	}
	_, _ = io.Copy(part, mediaFile)

	_ = writer.Close()

	debugLog(fmt.Sprintf("--> %s %s [MEDIA_DATA: Don't display]", method, rawUrl))

	if request, err = http.NewRequest(method, rawUrl, bodyBuffer); err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	return
}

type HandRespFunc func(resp *http.Response) error

func ExecuteHTTP(req *http.Request, handRespFunc HandRespFunc) (rawResp []byte, err error) {
	var resp *http.Response
	if resp, err = HTTPClient.Do(req); err != nil {
		return nil, err
	}
	end := time.Now()
	log.Println(fmt.Sprintf("<-- %s %s %d %s", req.Method, req.URL.String(), resp.StatusCode, time.Since(end)))

	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}()
	if err = handRespFunc(resp); err != nil {
		return nil, fmt.Errorf("error: %w\n", err)
	}
	return
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

var debugFlag = false

func SetDebug(debug bool) {
	debugFlag = debug
}

func debugLog(msg string) {
	if !debugFlag {
		return
	}
	log.Println("[DEBUG-WeCom-Bot-API] " + msg)
}
