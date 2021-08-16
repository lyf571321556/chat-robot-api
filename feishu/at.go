package feishu

type atMsg struct {
	Tag    string `json:"tag"`
	UserID string `json:"user_id"`
}

func newAtMsg(userId string) atMsg {
	msg := atMsg{
		Tag:    "at",
		UserID: userId,
	}
	return msg
}
