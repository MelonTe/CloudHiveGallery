package repository

import (
	"chg/internal/model/entity"
	"chg/pkg/db"

	"gorm.io/gorm"
)

// 数据库操作层
type SpaceRepository struct {
	db *gorm.DB
}

// 开启事务
func (r *SpaceRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func NewSpaceRepository() *SpaceRepository {
	return &SpaceRepository{db.LoadDB()}
}

// 若数据库查询失败，返回err，若不存在记录，space为nil
func (r *SpaceRepository) GetSpaceById(tx *gorm.DB, id uint64) (*entity.Space, error) {
	if tx == nil {
		tx = r.db
	}
	space := &entity.Space{}
	err := tx.First(space, id).Error
	// 区分记录不存在和查询异常
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return space, nil
}

// 根据空间id判断空间是否存在
func (r *SpaceRepository) IsExistById(tx *gorm.DB, id uint64) (bool, error) {
	if tx == nil {
		tx = r.db
	}
	var count int64
	err := tx.Model(&entity.Space{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// 空间的保存操作
func (r *SpaceRepository) SaveSpace(tx *gorm.DB, space *entity.Space) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Save(space).Error
}

// 根据表的字段更新空间信息
func (r *SpaceRepository) UpdateSpaceById(tx *gorm.DB, id uint64, updateMap map[string]interface{}) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Model(&entity.Space{ID: id}).Updates(updateMap).Error
}

// 根据用户ID判断空间是否存在
func (r *SpaceRepository) IsExistByUserId(tx *gorm.DB, userId uint64) bool {
	if tx == nil {
		tx = r.db
	}
	var exists bool
	tx.Raw("select exists(select 1 from spaces where user_id = ?)", userId).Scan(&exists)
	return exists
}
