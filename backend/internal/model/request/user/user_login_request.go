package user

// 登录请求参数
type UserLoginRequest struct {
	UserAccount  string `json:"userAccount" binding:"required"`
	UserPassword string `json:"userPassword" binding:"required"`
}
