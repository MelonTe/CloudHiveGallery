package consts

//定义使用websocket进行消协作编辑用到的常量
const (
	WS_PICTURE_EDIT_MESSAGE_INFO        = "INFO"        //发送通知信息
	WS_PICTURE_EDIT_MESSAGE_ERROR       = "ERROR"       //发送错误信息
	WS_PICTURE_EDIT_MESSAGE_ENTER_EDIT  = "ENTER_EDIT"  //进入编辑状态
	WS_PICTURE_EDIT_MESSAGE_EXIT_EDIT   = "EXIT_EDIT"   //退出编辑状态
	WS_PICTURE_EDIT_MESSAGE_EDIT_ACTION = "EDIT_ACTION" //执行编辑动作，如放大或缩小
)

//定义图片编辑动作常量
const (
	WS_PICTURE_EDIT_ACTION_ZOOM_IN      = "ZOOM_IN"      //放大
	WS_PICTURE_EDIT_ACTION_ZOOM_OUT     = "ZOOM_OUT"     //缩小
	WS_PICTURE_EDIT_ACTION_ROTATE_LEFT  = "ROTATE_LEFT"  //左旋转
	WS_PICTURE_EDIT_ACTION_ROTATE_RIGHT = "ROTATE_RIGHT" //右旋转
)

// IsEditAction 判断是否是编辑动作
func IsEditAction(action string) bool {
	switch action {
	case WS_PICTURE_EDIT_ACTION_ZOOM_IN,
		WS_PICTURE_EDIT_ACTION_ZOOM_OUT,
		WS_PICTURE_EDIT_ACTION_ROTATE_LEFT,
		WS_PICTURE_EDIT_ACTION_ROTATE_RIGHT:
		return true
	default:
		return false
	}
}

// GetActionName 获取编辑动作名称
func GetActionName(action string) string {
	switch action {
	case WS_PICTURE_EDIT_ACTION_ZOOM_IN:
		return "放大"
	case WS_PICTURE_EDIT_ACTION_ZOOM_OUT:
		return "缩小"
	case WS_PICTURE_EDIT_ACTION_ROTATE_LEFT:
		return "左旋转"
	case WS_PICTURE_EDIT_ACTION_ROTATE_RIGHT:
		return "右旋转"
	default:
		return "未知操作"
	}
}
