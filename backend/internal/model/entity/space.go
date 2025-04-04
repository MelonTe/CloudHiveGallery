package entity

import (
	"gorm.io/gorm"
	"time"
)

type Space struct {
	ID         uint64         `gorm:"primaryKey;comment:id" json:"id,string" swaggertype:"string"`
	SpaceName  string         `gorm:"type:varchar(128);comment:空间名称;index:idx_spaceName" json:"spaceName"`
	SpaceLevel int            `gorm:"default:0;comment:空间级别：0-普通版 1-专业版 2-旗舰版;index:idx_spaceLevel" json:"spaceLevel"`
	MaxSize    int64          `gorm:"default:0;comment:空间图片的最大总大小" json:"maxSize"`
	MaxCount   int64          `gorm:"default:0;comment:空间图片的最大数量" json:"maxCount"`
	TotalSize  int64          `gorm:"default:0;comment:当前空间下图片的总大小" json:"totalSize"`
	TotalCount int64          `gorm:"default:0;comment:当前空间下的图片数量" json:"totalCount"`
	UserID     uint64         `gorm:"not null;index:idx_userId;comment:创建用户 id" json:"userId,string" swaggertype:"string"`
	CreateTime time.Time      `gorm:"autoCreateTime;comment:创建时间" json:"createTime"`
	EditTime   time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP;not null;comment:编辑时间" json:"editTime"`
	UpdateTime time.Time      `gorm:"autoUpdateTime;comment:更新时间" json:"updateTime"`
	IsDelete   gorm.DeletedAt `gorm:"comment:是否删除" json:"isDelete" swaggerignore:"true"`
	SpaceType  int            `gorm:"default:0;comment:空间类型：0-个人空间 1-团队空间;index:idx_spaceType" json:"spaceType"`
}

// AutoMigrateSpace 执行数据库迁移
func AutoMigrateSpace(db *gorm.DB) {
	err := db.AutoMigrate(&Space{})
	if err != nil {
		panic("⚠️ 用户表迁移失败: " + err.Error())
	}
}

// 钩子，使用sonyflake生成ID
func (p *Space) BeforeCreate(tx *gorm.DB) error {
	if p.ID == 0 {
		id, err := sf.NextID()
		if err != nil {
			return err
		}
		p.ID = id
	}
	return nil
}
