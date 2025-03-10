package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// 初始化 Session 中间件
func InitSession(r *gin.Engine) {
	// 使用 Cookie 作为 Session 存储
	store := cookie.NewStore([]byte("cloud"))     // 对信息的加密密钥
	r.Use(sessions.Sessions("GSESSIONID", store)) // "GSESSION" 是 Cookie 的名称
}

// 设置 Session 数据
func SetSession(c *gin.Context, key string, value interface{}) error {
	session := sessions.Default(c)
	session.Set(key, value)
	return session.Save()
}

// 获取 Session 数据
func GetSession(c *gin.Context, key string) interface{} {
	session := sessions.Default(c)
	return session.Get(key)
}

// 删除 Session 数据
func DeleteSession(c *gin.Context, key string) error {
	session := sessions.Default(c)
	session.Delete(key)
	return session.Save()
}

// 清空 Session 数据
func ClearSession(c *gin.Context) error {
	session := sessions.Default(c)
	session.Clear()
	return session.Save()
}
