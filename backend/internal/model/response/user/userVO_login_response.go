package user

import (
	"chg/internal/model/entity"
	"time"
)

// 创建用户VO
type UserLoginVO struct {
	ID          uint64    `json:"id,string" swaggertype:"string"`
	UserAccount string    `json:"userAccount"`
	UserName    string    `json:"userName"`
	UserAvatar  string    `json:"userAvatar"`
	UserProfile string    `json:"userProfile"`
	UserRole    string    `json:"userRole"`
	EditTime    time.Time `json:"editTime"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
}

// 获取脱敏后的用户视图
func GetUserLoginVO(user entity.User) *UserLoginVO {
	return &UserLoginVO{
		ID:          user.ID,
		UserAccount: user.UserAccount,
		UserName:    user.UserName,
		UserAvatar:  user.UserAvatar,
		UserProfile: user.UserProfile,
		UserRole:    user.UserRole,
		EditTime:    user.EditTime,
		CreateTime:  user.CreateTime,
		UpdateTime:  user.UpdateTime,
	}
}
