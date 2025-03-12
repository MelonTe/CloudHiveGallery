package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//封装session函数

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
