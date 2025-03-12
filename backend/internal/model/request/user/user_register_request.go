package user

// 注册请求参数
type UserRegsiterRequest struct {
	UserAccount   string `json:"userAccount" binding:"required"`
	UserPassword  string `json:"userPassword" binding:"required"`
	CheckPassword string `json:"checkPassword" binding:"required"`
}
