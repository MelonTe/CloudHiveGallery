package entity

import (
	"gorm.io/gorm"
	"time"
)

type SpaceUser struct {
	ID         uint64    `gorm:"primaryKey;comment:id" json:"id,string" swaggertype:"string"`
	SpaceID    uint64    `gorm:"not null;index:idx_spaceId;comment:空间 id;uniqueIndex:uk_space_user" json:"spaceId,string" swaggertype:"string"`
	UserID     uint64    `gorm:"not null;index:idx_userId;comment:用户 id;uniqueIndex:uk_space_user" json:"userId,string" swaggertype:"string"`
	SpaceRole  string    `gorm:"type:varchar(128);default:'viewer';comment:空间角色：viewer/editor/admin" json:"spaceRole"`
	CreateTime time.Time `gorm:"autoCreateTime;comment:创建时间" json:"createTime"`
	UpdateTime time.Time `gorm:"autoUpdateTime;comment:更新时间" json:"updateTime"`
}

// AutoMigratePicture 执行数据库迁移
func AutoMigrateSpaceUser(db *gorm.DB) {
	err := db.AutoMigrate(&SpaceUser{})
	if err != nil {
		panic("⚠️ 用户表迁移失败: " + err.Error())
	}
}

// 钩子，使用sonyflake生成ID
func (su *SpaceUser) BeforeCreate(tx *gorm.DB) error {
	if su.ID == 0 {
		id, err := sf.NextID()
		if err != nil {
			return err
		}
		su.ID = id
	}
	return nil
}
