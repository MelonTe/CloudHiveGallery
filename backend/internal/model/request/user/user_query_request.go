package user

import "chg/internal/common"

type UserQueryRequest struct {
	common.PageRequest
	ID          uint64 `json:"id,string" swaggertype:"string"` //用户ID
	UserAccount string `json:"userAccount"`                    //用户账号
	UserName    string `json:"userName"`                       //用户昵称
	UserProfile string `json:"userProfile"`                    //用户简介
	UserRole    string `json:"userRole"`                       //用户权限
}
