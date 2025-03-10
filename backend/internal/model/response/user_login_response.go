package response

import (
	"chg/internal/model/entity"
	"time"
)

// 创建用户VO
type UserLoginVO struct {
	ID          uint64
	UserAccount string
	UserName    string
	UserAvatar  string
	UserProfile string
	UserRole    string
	EditTime    time.Time
	CreateTime  time.Time
	UpdateTime  time.Time
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
