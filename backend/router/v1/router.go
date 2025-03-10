package v1

import (
	"chg/internal/controller"

	"github.com/gin-gonic/gin"
)

// RegisterV1Routes 注册 v1 版本的路由
func RegisterV1Routes(r *gin.Engine) {
	apiV1 := r.Group("/v1")
	{
		userAPI := apiV1.Group("/user")
		userAPI.POST("/register", controller.UserRegister)
		userAPI.POST("/login", controller.UserLogin)
		userAPI.GET("/get/login", controller.GetLoginUser)
		userAPI.POST("/logout", controller.UserLogout)
	}

}
