package feishu

type TextMsg struct {
	// 消息类型, 固定为 `text`
	Tag      string `json:"tag"`
	Text     string `json:"text"`
	UnEscape bool   `json:"un_escape,omitempty"`
}

func NewTextMsg(content string) TextMsg {
	msg := TextMsg{
		Tag:  "text",
		Text: content,
	}
	return msg
}
