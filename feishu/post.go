package feishu

/**
{
        "chat_id":chat_id,
        "msg_type":"post",
        "content":{
            "post":{
                "zh_cn":{
                    "title":subject,
                    "content":[
                        [
                        {
                            "tag": "text",
                            "un_escape": True,
                            "text": messages
                        },
                        {
                            "tag": "at",
                            "user_id": user_id

                        }
                    ],
                    [
                        {
                            "tag": "img",
                            "image_key": image_key,
                            "width": 700,
                            "height": 400
                        }
                    ]
                ]
            }
        }
    }
    }
*/
const (
	MsgTypeText        string = "text"
	MsgTypePost        string = "post"
	MsgTypeImage       string = "image"
	MsgTypeShareChat   string = "share_chat"
	MsgTypeInteractive string = "interactive"
)

type PostItem interface{}

type PostUnit struct {
	Title   string       `json:"title,omitempty"`
	Content [][]PostItem `json:"content"`
}

type PostBody struct {
	ZH PostUnit `json:"zh_cn,omitempty"`
	JA PostUnit `json:"ja_jp,omitempty"`
	EN PostUnit `json:"en_us,omitempty"`
}

type PostContent struct {
	Post PostBody `json:"post"`
}

type PostMessage struct {
	MsgType string      `json:"msg_type"`
	Content PostContent `json:"content"`
}

func NewPostMessage() *PostMessage {
	return &PostMessage{
		MsgType: MsgTypePost,
	}
}

//func (m *PostMessage) Body() PostMessage {
//	return structToMap(m)
//}

func (m *PostMessage) SetZH(unit PostUnit) *PostMessage {
	m.Content.Post.ZH = unit
	return m
}

func (m *PostMessage) SetZHTitle(title string) *PostMessage {
	m.Content.Post.ZH.Title = title
	return m
}

func (m *PostMessage) AppendZHContent(items []PostItem) *PostMessage {
	m.Content.Post.ZH.Content = append(m.Content.Post.ZH.Content, items)
	return m
}

func (m *PostMessage) SetJA(unit PostUnit) *PostMessage {
	m.Content.Post.JA = unit
	return m
}

func (m *PostMessage) SetJATitle(title string) *PostMessage {
	m.Content.Post.JA.Title = title
	return m
}

func (m *PostMessage) AppendJAContent(items []PostItem) *PostMessage {
	m.Content.Post.JA.Content = append(m.Content.Post.JA.Content, items)
	return m
}

func (m *PostMessage) SetEN(unit PostUnit) *PostMessage {
	m.Content.Post.EN = unit
	return m
}

func (m *PostMessage) SetENTitle(title string) *PostMessage {
	m.Content.Post.EN.Title = title
	return m
}

func (m *PostMessage) AppendENContent(items []PostItem) *PostMessage {
	m.Content.Post.EN.Content = append(m.Content.Post.EN.Content, items)
	return m
}
