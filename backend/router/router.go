package router

//全局路由注册
import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"shg/router/v1" // 导入 v1 路由
)

// SetupRoutes 全局路由设置
func SetupRoutes(r *gin.Engine) {
	//注册swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 注册 v1 路由
	v1.RegisterV1Routes(r)
}
