package v1

import (
	"chg/internal/consts"
	"chg/internal/controller"
	"chg/internal/middleware"
	"chg/internal/service"

	"github.com/gin-gonic/gin"
)

// RegisterV1Routes 注册 v1 版本的路由
func RegisterV1Routes(r *gin.Engine) {
	apiV1 := r.Group("/v1")
	{
		userAPI := apiV1.Group("/user")
		{
			userAPI.POST("/register", controller.UserRegister)
			userAPI.POST("/login", controller.UserLogin)
			userAPI.GET("/get/login", controller.GetLoginUser)
			userAPI.POST("/logout", controller.UserLogout)
			userAPI.GET("/get/vo", controller.GetUserVOById)
			//以下需要权限
			userAPI.POST("/list/page/vo", middleware.AuthCheck(service.NewUserService(), consts.ADMIN_ROLE), controller.ListUserVOByPage)
			userAPI.POST("/update", middleware.AuthCheck(service.NewUserService(), consts.ADMIN_ROLE), controller.UpdateUser)
			userAPI.POST("/delete", middleware.AuthCheck(service.NewUserService(), consts.ADMIN_ROLE), controller.DeleteUser)
			userAPI.POST("/add", middleware.AuthCheck(service.NewUserService(), consts.ADMIN_ROLE), controller.AddUser)
			userAPI.GET("/get", middleware.AuthCheck(service.NewUserService(), consts.ADMIN_ROLE), controller.GetUserById)
		}
		fileAPI := apiV1.Group("/file")
		{
			fileAPI.POST("/test/upload", middleware.AuthCheck(service.NewUserService(), consts.ADMIN_ROLE), controller.TestUploadFile)
			fileAPI.GET("/test/download", middleware.AuthCheck(service.NewUserService(), consts.ADMIN_ROLE), controller.TestDownloadFile)
		}
		pictureAPI := apiV1.Group("/picture")
		{
			pictureAPI.POST("/upload", middleware.AuthCheck(service.NewUserService(), consts.ADMIN_ROLE), controller.UploadPicture)
		}
	}

}
