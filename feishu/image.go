package feishu

type imageMsg struct {
	// 消息类型, 固定为 `image`
	Tag      string `json:"tag"`
	ImageKey string `json:"image_key"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}

func NewImageMsg(imageKey string, height, width int) imageMsg {
	image := imageMsg{
		Tag:      "image",
		ImageKey: imageKey,
		Height:   height,
		Width:    width,
	}
	return image
}
