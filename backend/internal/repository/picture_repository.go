package repository

import (
	"chg/internal/model/entity"
	"chg/pkg/db"
	"errors"
	"gorm.io/gorm"
)

// 数据库操作层
type PictureRepository struct {
	db *gorm.DB
}

func NewPictureRepository() *PictureRepository {
	return &PictureRepository{db.LoadDB()}
}

// 开启事务
func (r *PictureRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

// 根据ID查找图片
func (r *PictureRepository) FindById(tx *gorm.DB, id uint64) (*entity.Picture, error) {
	if tx == nil {
		tx = r.db
	}
	var picture entity.Picture
	if err := tx.Where("id = ?", id).First(&picture).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 无记录
		}
		return nil, err // 数据库查询异常
	}
	return &picture, nil
}

// save图片
func (r *PictureRepository) SavePicture(tx *gorm.DB, picture *entity.Picture) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Save(picture).Error
}

// 删除图片
func (r *PictureRepository) DeleteById(tx *gorm.DB, id uint64) error {
	if tx == nil {
		tx = r.db
	}
	err := tx.Where("id = ?", id).Delete(&entity.Picture{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // 无记录
		}
		return err
	}
	return nil
}

// 通过map更新图片，只更新map中含有的字段
func (r *PictureRepository) UpdateById(tx *gorm.DB, id uint64, updateMap map[string]interface{}) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Model(&entity.Picture{ID: id}).Updates(updateMap).Error
}
