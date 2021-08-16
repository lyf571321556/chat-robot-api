package api

import (
	"bytes"
	"fmt"
	"github.com/lyf571321556/chat-robot-api/file"
	"github.com/lyf571321556/chat-robot-api/news"
	"github.com/lyf571321556/chat-robot-api/text"
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

type Robot interface {
	PushTextMessage(content string, opts ...text.TextMsgOption) error
	PushMarkdownMessage(content string) error
	PushImageMessage(img []byte) error
	PushNewsMessage(art news.Article, articles ...news.Article) error
	PushFileMessage(media file.Media) error
	UploadFile(filename string) (media file.Media, err error)
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
