package user

type UserUploadAvatarRequest struct {
	ID          uint64 `json:"id,string" swaggertype:"string"` //用户ID
	UserAccount string `json:"userAccount" binding:"required"` //用户账号
	UserName    string `json:"userName"`                       //用户昵称
	UserAvatar  string `json:"userAvatar"`                     //用户头像
	UserProfile string `json:"userProfile"`                    //用户简介
	UserRole    string `json:"userRole"`                       //用户权限
}
