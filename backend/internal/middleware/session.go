package middleware

import (
	"chg/internal/model/entity"
	"encoding/gob"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// 需要提前注册数据结构，否则无法存储
func init() {
	gob.Register(entity.User{})
}

// 初始化 Session 中间件
func InitSession(r *gin.Engine) {
	// 使用 Cookie 作为 Session 存储
	store := cookie.NewStore([]byte("cloud"))     // 对信息的加密密钥
	r.Use(sessions.Sessions("GSESSIONID", store)) // "GSESSION" 是 Cookie 的名称
}
