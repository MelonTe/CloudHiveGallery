package repository

import (
	"chg/internal/model/entity"
	"chg/pkg/db"

	"gorm.io/gorm"
)

// 数据库操作层
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db.LoadDB()}
}

// 根据账号查找用户
func (r *UserRepository) FindByAccount(userAccount string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("user_account = ?", userAccount).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser 创建新用户
func (r *UserRepository) CreateUser(user *entity.User) error {
	return r.db.Create(user).Error
}

// CountByAccount 统计账号数量（用于判断账号是否重复）
func (r *UserRepository) CountByAccount(userAccount string) (int64, error) {
	var count int64
	err := r.db.Model(&entity.User{}).Where("user_account = ?", userAccount).Count(&count).Error
	return count, err
}
