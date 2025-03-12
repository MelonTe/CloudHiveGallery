package user

import (
	"chg/internal/model/entity"
	"time"
)

// 脱敏后的用户信息
type UserVO struct {
	ID          uint64    `json:"id,string" swaggertype:"string"`
	UserAccount string    `json:"userAccount"`
	UserName    string    `json:"userName"`
	UserAvatar  string    `json:"userAvatar"`
	UserProfile string    `json:"userProfile"`
	UserRole    string    `json:"userRole"`
	CreateTime  time.Time `json:"createTime"`
}

func GetUserVO(user entity.User) UserVO {
	return UserVO{
		ID:          user.ID,
		UserAccount: user.UserAccount,
		UserName:    user.UserName,
		UserAvatar:  user.UserAvatar,
		UserProfile: user.UserProfile,
		UserRole:    user.UserRole,
		CreateTime:  user.CreateTime,
	}
}

func GetUserVOList(users []entity.User) []UserVO {
	userVOList := make([]UserVO, 0)
	for _, user := range users {
		userVOList = append(userVOList, GetUserVO(user))
	}
	return userVOList
}
