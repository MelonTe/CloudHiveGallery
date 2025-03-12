package repository

import (
	"chg/internal/model/entity"
	"chg/pkg/db"
	"errors"
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //无记录
		}
		return nil, err //数据库查询异常
	}
	return &user, nil
}

// 根据ID查找用户
func (r *UserRepository) FindById(id uint64) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //无记录
		}
		return nil, err //数据库查询异常
	}
	return &user, nil
}

// 根据账号和密码查找用户
func (r *UserRepository) FindByAccountAndPassword(userAccount string, userPassword string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("user_account = ? AND user_password = ?", userAccount, userPassword).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //无记录
		}
		return nil, err //数据库查询异常
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

func (r *UserRepository) RemoveById(id uint64) (bool, error) {
	result := r.db.Delete(&entity.User{}, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, result.Error
		}
	}
	return true, nil
}

func (r *UserRepository) UpdateUser(user *entity.User) (bool, error) {
	result := r.db.Model(&entity.User{}).Where("id = ?", user.ID).Updates(user)
	err := result.Error
	if err != nil {
		return false, err
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (r *UserRepository) ListUserByPage(query *gorm.DB) ([]entity.User, error) {
	var users []entity.User
	err := query.Find(&users).Error
	return users, err
}

// 获取query查询到的user数量
func (r *UserRepository) GetQueryUsersNum(query *gorm.DB) (int, error) {
	total := int64(0)
	query.Find(&[]entity.User{}).Count(&total)
	return int(total), nil
}
