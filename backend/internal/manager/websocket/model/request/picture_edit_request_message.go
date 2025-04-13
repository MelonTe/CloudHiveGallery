package request

//图片编辑请求消息
type PictureEditRequestMessage struct {
	Type       string `json:"type"`       // 消息类型，如"ENTER_EDIT","EXIT_EDIT","EDIT_ACTION"
	EditAction string `json:"editAction"` // 编辑动作，如放大或缩小
}
