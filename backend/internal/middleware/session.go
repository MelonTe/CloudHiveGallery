package middleware

import (
	"chg/config"
	"chg/internal/model/entity"
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"log"
)

// 需要提前注册数据结构，否则无法存储
func init() {
	gob.Register(entity.User{})
}

// 初始化 Session 中间件，使用redis存储
func InitSession(r *gin.Engine) {
	cfg := config.LoadConfig()
	store, err := redis.NewStore(10, "tcp", fmt.Sprintf("%s:%d", cfg.Rds.Host, cfg.Rds.Port), cfg.Rds.Password, []byte("cloud"))
	if err != nil {
		log.Fatalf("初始化 Redis Session 失败: %v", err)
	}
	// 设置 Session 选项
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   3600,  // Session 过期时间（秒），这里是 1 小时
		HttpOnly: true,  // 保护 Session，不让 JS 访问
		Secure:   false, // 生产环境应设为 true（HTTPS）
	})
	r.Use(sessions.Sessions("GSESSIONID", store)) // "GSESSION" 是 Cookie 的名称
}
