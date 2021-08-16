package feishu

type aTagMsg struct {
	Tag      string `json:"tag"`
	Text     string `json:"text"`
	Href     string `json:"href"`
	UnEscape bool   `json:"un_escape,omitempty"`
}

func newATagMsg(text, href string) aTagMsg {
	msg := aTagMsg{
		Tag:  "a",
		Text: text,
		Href: href,
	}
	return msg
}
