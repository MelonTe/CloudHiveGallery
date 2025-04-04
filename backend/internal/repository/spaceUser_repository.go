package repository

import (
	"chg/pkg/db"
	"gorm.io/gorm"
)

// 数据库操作层
type SpaceUserRepository struct {
	db *gorm.DB
}

// 开启事务
func (r *SpaceUserRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func NewSpaceUserRepository() *SpaceUserRepository {
	return &SpaceUserRepository{db.LoadDB()}
}
