package entity

import (
	"github.com/sony/sonyflake"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           uint64         `gorm:"primaryKey;comment:id" json:"id,string" swaggertype:"string"`
	UserAccount  string         `gorm:"type:varchar(256);uniqueIndex;not null;comment:账号"`
	UserPassword string         `gorm:"type:varchar(512);not null;comment:密码"`
	UserName     string         `gorm:"type:varchar(256);index;comment:用户昵称"`
	UserAvatar   string         `gorm:"type:varchar(1024);comment:用户头像"`
	UserProfile  string         `gorm:"type:varchar(512);comment:用户简介"`
	UserRole     string         `gorm:"type:varchar(256);default:user;not null;comment:用户角色：user/admin"`
	EditTime     time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP;not null;comment:编辑时间"`
	CreateTime   time.Time      `gorm:"autoCreateTime;comment:创建时间"`
	UpdateTime   time.Time      `gorm:"autoUpdateTime;comment:更新时间"`
	IsDelete     gorm.DeletedAt `gorm:"comment:是否删除" swaggerignore:"true"`
}

// AutoMigrateUser 执行数据库迁移
func AutoMigrateUser(db *gorm.DB) {
	err := db.AutoMigrate(&User{})
	if err != nil {
		panic("⚠️ 用户表迁移失败: " + err.Error())
	}
}

var sf = sonyflake.NewSonyflake(sonyflake.Settings{})

// 钩子，使用sonyflake生成ID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == 0 {
		id, err := sf.NextID()
		if err != nil {
			return err
		}
		u.ID = id
	}
	return nil
}
