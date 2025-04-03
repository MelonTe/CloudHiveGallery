package repository

import (
	"chg/pkg/db"
	"gorm.io/gorm"
)

// 数据库操作层
type SpaceAnalyzeRepository struct {
	db *gorm.DB
}

// 开启事务
func (r *SpaceAnalyzeRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func NewSpaceAnalyzeRepository() *SpaceAnalyzeRepository {
	return &SpaceAnalyzeRepository{db.LoadDB()}
}
