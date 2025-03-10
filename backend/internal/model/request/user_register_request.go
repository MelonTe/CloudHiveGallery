package request

//登录请求参数
type UserRegsiterRequest struct {
	UserAccount   string `json:"userAccount" binding:"required"`
	UserPassword  string `json:"userPassword" binding:"required"`
	CheckPassword string `json:"checkPassword" binding:"required"`
}
