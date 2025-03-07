package v1

import (
	"shg/internal/controller"

	"github.com/gin-gonic/gin"
)

// RegisterV1Routes 注册 v1 版本的路由
func RegisterV1Routes(r *gin.Engine) {
	r.GET("/v1/hello", controller.HelloWorld)
}
