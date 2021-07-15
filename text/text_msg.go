package text

type TextMsg struct {
	// 消息类型, 固定为 `text`
	MsgType     string       `json:"msgtype"`
	ContentData *ContentData `json:"text"`
}

type ContentData struct {
	// **必填**
	//
	// 文本内容, 最长不超过 2048 个字节, 必须是 `utf8` 编码
	Content string `json:"content"`

	// **选填**
	//
	// `userid` 列表, 提醒群中的指定成员 (@某个成员), `@all` 表示提醒所有人,
	//
	// 如果获取不到 `userid`, 可以使用 `mentioned_mobile_list`
	MentionedList []string `json:"mentioned_list,omitempty"`

	// **选填**
	//
	// 手机号列表, 提醒手机号对应的群成员 (@某个成员), `@all`表示提醒所有人
	MentionedMobileList []string `json:"mentioned_mobile_list,omitempty"`
}

type TextMsgOption func(d *ContentData)

func NewTextMsg(content string, options ...TextMsgOption) TextMsg {
	msg := TextMsg{
		MsgType: "text",
		ContentData: &ContentData{
			Content:             content,
			MentionedList:       make([]string, 0),
			MentionedMobileList: make([]string, 0),
		},
	}

	for _, option := range options {
		option(msg.ContentData)
	}
	return msg
}

// MentionByUserid 通过 `userid` @某个成员
func MentionByUserid(userid string) TextMsgOption {
	return func(contentData *ContentData) {
		contentData.MentionedList = append(contentData.MentionedList, userid)
	}
}

// MentionAllByUserid `@all` 提醒所有人, 等同于 MentionAllByMobile
func MentionAllByUserid() TextMsgOption {
	return func(contentData *ContentData) {
		contentData.MentionedList = append(contentData.MentionedList, "@all")
	}
}

// MentionByMobile 通过 `手机号码` @某个成员
func MentionByMobile(mobile string) TextMsgOption {
	return func(contentData *ContentData) {
		contentData.MentionedMobileList = append(contentData.MentionedMobileList, mobile)
	}
}

// MentionAllByMobile `@all` 提醒所有人, 等同于 MentionAllByUserid
func MentionAllByMobile() TextMsgOption {
	return func(contentData *ContentData) {
		contentData.MentionedMobileList = append(contentData.MentionedMobileList, "@all")
	}
}
