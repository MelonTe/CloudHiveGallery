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

// 开启事务
func (r *UserRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

// 根据账号查找用户
func (r *UserRepository) FindByAccount(tx *gorm.DB, userAccount string) (*entity.User, error) {
	if tx == nil {
		tx = r.db
	}
	var user entity.User
	if err := tx.Where("user_account = ?", userAccount).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //无记录
		}
		return nil, err //数据库查询异常
	}
	return &user, nil
}

// 根据ID查找用户
func (r *UserRepository) FindById(tx *gorm.DB, id uint64) (*entity.User, error) {
	if tx == nil {
		tx = r.db
	}
	var user entity.User
	if err := tx.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //无记录
		}
		return nil, err //数据库查询异常
	}
	return &user, nil
}

// 根据账号和密码查找用户
func (r *UserRepository) FindByAccountAndPassword(tx *gorm.DB, userAccount string, userPassword string) (*entity.User, error) {
	if tx == nil {
		tx = r.db
	}
	var user entity.User
	if err := tx.Where("user_account = ? AND user_password = ?", userAccount, userPassword).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //无记录
		}
		return nil, err //数据库查询异常
	}
	return &user, nil
}

// CreateUser 创建新用户
func (r *UserRepository) CreateUser(tx *gorm.DB, user *entity.User) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(user).Error
}

// CountByAccount 统计账号数量（用于判断账号是否重复）
func (r *UserRepository) CountByAccount(tx *gorm.DB, userAccount string) (int64, error) {
	if tx == nil {
		tx = r.db
	}
	var count int64
	err := tx.Model(&entity.User{}).Where("user_account = ?", userAccount).Count(&count).Error
	return count, err
}

func (r *UserRepository) RemoveById(tx *gorm.DB, id uint64) (bool, error) {
	if tx == nil {
		tx = r.db
	}
	result := tx.Delete(&entity.User{}, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, result.Error
		}
	}
	return true, nil
}

func (r *UserRepository) UpdateUser(tx *gorm.DB, user *entity.User) (bool, error) {
	if tx == nil {
		tx = r.db
	}
	result := tx.Model(&entity.User{}).Where("id = ?", user.ID).Updates(user)
	err := result.Error
	if err != nil {
		return false, err
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (r *UserRepository) ListUserByPage(tx *gorm.DB, query *gorm.DB) ([]entity.User, error) {
	if tx == nil {
		tx = r.db
	}
	var users []entity.User
	err := query.Find(&users).Error
	return users, err
}

// 获取query查询到的user数量
func (r *UserRepository) GetQueryUsersNum(tx *gorm.DB, query *gorm.DB) (int, error) {
	if tx == nil {
		tx = r.db
	}
	total := int64(0)
	query.Find(&[]entity.User{}).Count(&total)
	return int(total), nil
}
