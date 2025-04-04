package db

import (
	"chg/config"
	"chg/internal/model/entity"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var db *gorm.DB

func init() {
	config := config.LoadConfig()
	//初始化数据库
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Database.User,
		config.Database.Password,
		config.Database.Port,
		config.Database.Name)
	var err error
	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 设置日志级别为 Info
	}); err != nil {
		log.Fatalf("Failed to connect DB, %s", err)
	}
	//自动迁移model
	entity.AutoMigrateUser(db)
	entity.AutoMigratePicture(db)
	entity.AutoMigrateSpace(db)
	entity.AutoMigrateSpaceUser(db)
}
func LoadDB() *gorm.DB {
	return db.Session(&gorm.Session{})
}
