package controller

import (
	"chg/internal/common"
	reqPicture "chg/internal/model/request/picture"
	resPicture "chg/internal/model/response/picture"
	"chg/internal/service"
	"github.com/gin-gonic/gin"
)

// 获取一个pictureService单例
var sPicture *service.PictureService = service.NewPictureService()
var sUser *service.UserService = service.NewUserService()

func init() {
	//防止引用错误
	_ = resPicture.PictureVO{}
}

// UploadPicture godoc
// @Summary      上传图片接口「管理员」
// @Description  根据是否存在ID来上传图片或者修改图片信息，返回图片信息视图
// @Tags         picture
// @Accept       mpfd
// @Produce      json
// @Param        file formData file true "图片"
// @Param        id formData string false "图片的ID，非必需"
// @Success      200  {object}  common.Response{data=resPicture.PictureVO} "上传成功，返回图片信息视图"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /v1/picture/upload [POST]
func UploadPicture(c *gin.Context) {
	file, _ := c.FormFile("file")
	picReq := &reqPicture.PictureUploadRequest{}
	c.ShouldBind(picReq)
	//因为存在中间件，所以一定会成功获取
	loginUser, _ := sUser.GetLoginUser(c)
	picVO, err := sPicture.UploadPicture(file, picReq, loginUser)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	common.Success(c, *picVO)
}
