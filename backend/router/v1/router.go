package v1

import (
	"chg/internal/controller"

	"github.com/gin-gonic/gin"
)

// RegisterV1Routes 注册 v1 版本的路由
func RegisterV1Routes(r *gin.Engine) {
	apiV1 := r.Group("/v1")
	apiV1.GET("/hello", controller.HelloWorld)
}
