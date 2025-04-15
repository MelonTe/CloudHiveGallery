package v1

import (
	"chg/internal/consts"
	"chg/internal/controller"
	"chg/internal/manager/websocket"
	"chg/internal/middleware"

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
			userAPI.POST("/list/page/vo", middleware.AuthCheck(consts.ADMIN_ROLE), controller.ListUserVOByPage)
			userAPI.POST("/update", middleware.AuthCheck(consts.ADMIN_ROLE), controller.UpdateUser)
			userAPI.POST("/delete", middleware.AuthCheck(consts.ADMIN_ROLE), controller.DeleteUser)
			userAPI.POST("/add", middleware.AuthCheck(consts.ADMIN_ROLE), controller.AddUser)
			userAPI.GET("/get", middleware.AuthCheck(consts.ADMIN_ROLE), controller.GetUserById)
			userAPI.POST("/avatar", middleware.LoginCheck(), controller.UploadAvatar)
			userAPI.POST("/edit", middleware.LoginCheck(), controller.EditUser)
		}
		fileAPI := apiV1.Group("/file")
		{
			fileAPI.POST("/test/upload", middleware.AuthCheck(consts.ADMIN_ROLE), controller.TestUploadFile)
			fileAPI.GET("/test/download", middleware.AuthCheck(consts.ADMIN_ROLE), controller.TestDownloadFile)
		}
		pictureAPI := apiV1.Group("/picture")
		{
			pictureAPI.POST("/upload", middleware.LoginCheck(), controller.UploadPicture)
			pictureAPI.POST("/upload/url", middleware.LoginCheck(), controller.UploadPictureByUrl)
			pictureAPI.POST("/upload/batch", middleware.AuthCheck(consts.ADMIN_ROLE), controller.UploadPictureByBatch)
			pictureAPI.POST("/delete", middleware.LoginCheck(), controller.DeletePicture)
			pictureAPI.POST("/update", middleware.LoginCheck(), controller.UpdatePicture)
			pictureAPI.POST("/edit", middleware.LoginCheck(), controller.UpdatePicture)
			pictureAPI.GET("/get", middleware.AuthCheck(consts.ADMIN_ROLE), controller.GetPictureById)
			pictureAPI.GET("/get/vo", controller.GetPictureVOById)
			pictureAPI.POST("/list/page", middleware.AuthCheck(consts.ADMIN_ROLE), controller.ListPictureByPage)
			pictureAPI.POST("/list/page/vo", controller.ListPictureVOByPage)
			pictureAPI.POST("/list/page/vo/cache", controller.ListPictureVOByPageWithCache)
			pictureAPI.GET("/tag_category", controller.ListPictureTagCategory)
			pictureAPI.POST("/review", middleware.AuthCheck(consts.ADMIN_ROLE), controller.DoPictureReview)
			pictureAPI.POST("/search/picture", controller.SearchPictureByPicture)
			pictureAPI.POST("/search/color", middleware.LoginCheck(), controller.SearchPictureByColor)
			pictureAPI.POST("/edit/batch", middleware.LoginCheck(), controller.PictureEditByBatch)
			pictureAPI.POST("/out_painting/create_task", middleware.LoginCheck(), controller.CreatePictureOutPaintingTask)
			pictureAPI.GET("/out_painting/create_task", middleware.LoginCheck(), controller.GetOutPaintingTaskResponse)
		}
		spaceAPI := apiV1.Group("/space")
		{
			spaceAPI.POST("/update", middleware.AuthCheck(consts.ADMIN_ROLE), controller.UpdateSpace)
			spaceAPI.POST("/edit", controller.EditPicture)
			spaceAPI.POST("/list/page", middleware.AuthCheck(consts.ADMIN_ROLE), controller.ListSpaceByPage)
			spaceAPI.POST("/list/page/vo", controller.ListSpaceVOByPage)
			spaceAPI.POST("/add", middleware.LoginCheck(), controller.AddSpace)
			spaceAPI.GET("/list/level", controller.ListSpaceLevel)
			spaceAPI.GET("/get/vo", middleware.LoginCheck(), controller.GetSpaceVOById)
		}
		spaceAnalyzeAPI := apiV1.Group("/space/analyze")
		{
			spaceAnalyzeAPI.POST("/usage", middleware.LoginCheck(), controller.GetSpaceUsageAnalyze)
			spaceAnalyzeAPI.POST("/category", middleware.LoginCheck(), controller.GetSpaceCategoryAnalyze)
			spaceAnalyzeAPI.POST("/tag", middleware.LoginCheck(), controller.GetSpaceTagAnalyze)
			spaceAnalyzeAPI.POST("/size", middleware.LoginCheck(), controller.GetSpaceSizeAnalyze)
			spaceAnalyzeAPI.POST("/user", middleware.LoginCheck(), controller.GetSpaceUserAnalyze)
			spaceAnalyzeAPI.POST("/rank", middleware.AuthCheck(consts.ADMIN_ROLE), controller.GetSpaceRankAnalyze)
		}
		spaceUserAPI := apiV1.Group("/spaceUser")
		{
			spaceUserAPI.POST("/add", middleware.CasbinAuthCheck(consts.DOM_SPACE, consts.OBJ_SPACEUSER, consts.ACT_SPACEUSER_MANAGE), controller.AddSpaceUser)
			spaceUserAPI.POST("/delete", middleware.CasbinAuthCheck(consts.DOM_SPACE, consts.OBJ_SPACEUSER, consts.ACT_SPACEUSER_MANAGE), controller.DeleteSpaceUser)
			spaceUserAPI.POST("/get", controller.GetSpaceUser)
			spaceUserAPI.POST("/list", controller.ListSpaceUser)
			spaceUserAPI.POST("/edit", middleware.CasbinAuthCheck(consts.DOM_SPACE, consts.OBJ_SPACEUSER, consts.ACT_SPACEUSER_MANAGE), controller.EditSpaceUser)
			spaceUserAPI.POST("/list/my", controller.ListMyTeamSpace)
		}
		apiV1.GET("/ws/picture/edit", websocket.PictureEditHandShake)
	}
}
