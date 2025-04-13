package response

import resUser "chg/internal/model/response/user"

//图片响应编辑消息
type PictureEditResponseMessage struct {
	Type       string         `json:"type"`       // 消息类型，例如 "INFO", "ERROR", "ENTER_EDIT", "EXIT_EDIT", "EDIT_ACTION"
	Message    string         `json:"message"`    // 信息
	EditAction string         `json:"editAction"` // 执行的编辑动作
	User       resUser.UserVO `json:"user"`       // 用户信息
}
