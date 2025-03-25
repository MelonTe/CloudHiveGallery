package entity

import (
	"gorm.io/gorm"
	"time"
)

type Picture struct {
	ID            uint64         `gorm:"primaryKey;comment:id" json:"id,string" swaggertype:"string"`
	URL           string         `gorm:"type:varchar(512);not null;comment:图片 url" json:"url"`
	ThumbnailURL  string         `gorm:"type:varchar(512);comment:缩略图 url;default:null" json:"thumbnailUrl"`
	Name          string         `gorm:"type:varchar(128);not null;index:idx_name;comment:图片名称" json:"name"`
	Introduction  string         `gorm:"type:varchar(512);index:idx_introduction;comment:简介" json:"introduction"`
	Category      string         `gorm:"type:varchar(64);index:idx_category;comment:分类" json:"category"`
	Tags          string         `gorm:"type:varchar(512);index:idx_tags;comment:标签（JSON 数组）" json:"tags"` //存储的格式：["golang","java","c++"]
	PicSize       int64          `gorm:"comment:图片体积" json:"picSize"`
	PicWidth      int            `gorm:"comment:图片宽度" json:"picWidth"`
	PicHeight     int            `gorm:"comment:图片高度" json:"picHeight"`
	PicScale      float64        `gorm:"comment:图片宽高比例" json:"picScale"`
	PicFormat     string         `gorm:"type:varchar(32);comment:图片格式" json:"picFormat"`
	UserID        uint64         `gorm:"not null;index:idx_userId;comment:创建用户 id" json:"userId,string" swaggertype:"string"`
	EditTime      time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP;not null;comment:编辑时间" json:"editTime"`
	CreateTime    time.Time      `gorm:"autoCreateTime;comment:创建时间" json:"createTime"`
	UpdateTime    time.Time      `gorm:"autoUpdateTime;comment:更新时间" json:"updateTime"`
	IsDelete      gorm.DeletedAt `gorm:"comment:是否删除" json:"isDelete" swaggerignore:"true"`
	ReviewStatus  int            `gorm:"default:0;comment:审核状态：0-待审核；1-通过；2-拒绝;not null;index:idx_reviewStatus" json:"reviewStatus"`
	ReviewMessage string         `gorm:"type:varchar(512);comment:审核信息" json:"reviewMessage"`
	ReviewerID    uint64         `gorm:"comment:审核人 ID" json:"reviewerId,string" swaggertype:"string"`
	ReviewTime    time.Time      `gorm:"type:datetime;comment:审核时间" json:"reviewTime"`
}

// AutoMigrateUser 执行数据库迁移
func AutoMigratePicture(db *gorm.DB) {
	err := db.AutoMigrate(&Picture{})
	if err != nil {
		panic("⚠️ 用户表迁移失败: " + err.Error())
	}
}

// 钩子，使用sonyflake生成ID
func (p *Picture) BeforeCreate(tx *gorm.DB) error {
	if p.ID == 0 {
		id, err := sf.NextID()
		if err != nil {
			return err
		}
		p.ID = id
	}
	return nil
}
