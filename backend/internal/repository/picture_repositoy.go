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

// 根据ID查找图片
func (r *PictureRepository) FindById(id uint64) (*entity.Picture, error) {
	var picture entity.Picture
	if err := r.db.Where("id = ?", id).First(&picture).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //无记录
		}
		return nil, err //数据库查询异常
	}
	return &picture, nil
}

// save图片
func (r *PictureRepository) SavePicture(picture *entity.Picture) error {
	return r.db.Save(picture).Error
}

// 删除图片
func (r *PictureRepository) DeleteById(id uint64) error {
	err := r.db.Where("id = ?", id).Delete(&entity.Picture{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil //无记录
		}
		return err
	}
	return nil
}

// 通过map更新图片，只更新map中含有的字段
func (r *PictureRepository) UpdateById(id uint64, updateMap map[string]interface{}) error {
	return r.db.Model(&entity.Picture{ID: id}).Updates(updateMap).Error
}
