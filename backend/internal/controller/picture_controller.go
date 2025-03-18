package controller

import (
	"chg/internal/common"
	"chg/internal/ecode"
	reqPicture "chg/internal/model/request/picture"
	resPicture "chg/internal/model/response/picture"
	"chg/internal/service"
	"strconv"

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

// DeletePicture godoc
// @Summary      根据ID软删除图片
// @Tags         picture
// @Accept       json
// @Produce      json
// @Param		id body common.DeleteRequest true "图片的ID"
// @Success      200  {object}  common.Response{data=bool} "删除成功"
// @Failure      400  {object}  common.Response "删除失败，详情见响应中的code"
// @Router       /v1/picture/delete [POST]
func DeletePicture(c *gin.Context) {
	deleReq := common.DeleteRequest{}
	c.ShouldBind(&deleReq)
	if deleReq.Id <= 0 {
		common.BaseResponse(c, false, "删除失败，参数错误", ecode.PARAMS_ERROR)
		return
	}
	//判断id图片是否存在
	oldPic, err := sPicture.GetPictureById(deleReq.Id)
	if err != nil {
		common.BaseResponse(c, false, err.Msg, err.Code)
		return
	}
	//仅本人或管理员允许删除图片
	user, _ := sUser.GetLoginUser(c)
	if user == nil || !(oldPic.UserID == user.ID) && !(user.UserRole == "admin") {
		common.BaseResponse(c, false, "无权限", ecode.NO_AUTH_ERROR)
		return
	}
	if err := sPicture.DeletePictureById(deleReq.Id); err != nil {
		common.BaseResponse(c, false, err.Msg, err.Code)
		return
	}
	common.Success(c, true)
}

// UpdatePicture godoc
// @Summary      更新图片「管理员」
// @Description  若图片不存在，则返回false
// @Tags         picture
// @Accept       json
// @Produce      json
// @Param		request body reqPicture.PictureUpdateRequest true "需要更新的图片信息"
// @Success      200  {object}  common.Response{data=bool} "更新成功"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /v1/picture/update [POST]
func UpdatePicture(c *gin.Context) {
	updateReq := reqPicture.PictureUpdateRequest{}
	c.ShouldBind(&updateReq)
	if updateReq.ID <= 0 {
		common.BaseResponse(c, false, "参数错误", ecode.PARAMS_ERROR)
		return
	}
	//更新操作，参数校验等在service层完成
	if err := sPicture.UpdatePicture(&updateReq); err != nil {
		common.BaseResponse(c, false, err.Msg, err.Code)
		return
	}
	common.Success(c, true)
}

// GetPictureById godoc
// @Summary      根据ID获取图片「管理员」
// @Tags         picture
// @Accept		json
// @Produce      json
// @Param		id query string true "图片的ID"
// @Success      200  {object}  common.Response{data=entity.Picture} "获取成功"
// @Failure      400  {object}  common.Response "获取失败，详情见响应中的code"
// @Router       /v1/picture/get [GET]
func GetPictureById(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	if id <= 0 {
		common.BaseResponse(c, nil, "参数错误", ecode.PARAMS_ERROR)
		return
	}
	pic, err := sPicture.GetPictureById(id)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	common.Success(c, *pic)
}

// GetPictureVOById godoc
// @Summary      根据ID获取脱敏的图片
// @Tags         picture
// @Accept		json
// @Produce      json
// @Param		id query string true "图片的ID"
// @Success      200  {object}  common.Response{data=resPicture.PictureVO} "获取成功"
// @Failure      400  {object}  common.Response "获取失败，详情见响应中的code"
// @Router       /v1/picture/get/vo [GET]
func GetPictureVOById(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	if id <= 0 {
		common.BaseResponse(c, nil, "参数错误", ecode.PARAMS_ERROR)
		return
	}
	pic, err := sPicture.GetPictureById(id)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	picVO := sPicture.GetPictureVO(pic)
	common.Success(c, *picVO)
}

// ListPictureByPage godoc
// @Summary      分页获取一系列图片信息「管理员」
// @Tags         picture
// @Accept       json
// @Produce      json
// @Param		request body reqPicture.PictureQueryRequest true "需要查询的页数、以及图片关键信息"
// @Success      200  {object}  common.Response{data=resPicture.ListPictureResponse} "查询成功"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /v1/picture/list/page [POST]
func ListPictureByPage(c *gin.Context) {
	queryReq := reqPicture.PictureQueryRequest{}
	c.ShouldBind(&queryReq)
	//获取分页查询对象
	pics, err := sPicture.ListPictureByPage(&queryReq)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	common.Success(c, *pics)
}

// ListPictureVOByPage godoc
// @Summary      分页获取一系列图片信息
// @Tags         picture
// @Accept       json
// @Produce      json
// @Param		request body reqPicture.PictureQueryRequest true "需要查询的页数、以及图片关键信息"
// @Success      200  {object}  common.Response{data=resPicture.ListPictureVOResponse} "查询成功"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /v1/picture/list/page/vo [POST]
func ListPictureVOByPage(c *gin.Context) {
	queryReq := reqPicture.PictureQueryRequest{}
	c.ShouldBind(&queryReq)
	//限制爬虫
	if queryReq.PageSize > 20 {
		common.BaseResponse(c, nil, "参数错误", ecode.PARAMS_ERROR)
		return
	}
	//获取分页查询对象
	pics, err := sPicture.ListPictureVOByPage(&queryReq)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	common.Success(c, *pics)
}

// EditPicture godoc
// @Summary      更新图片
// @Description  若图片不存在，则返回false
// @Tags         picture
// @Accept       json
// @Produce      json
// @Param		request body reqPicture.PictureEditRequest true "需要更新的图片信息"
// @Success      200  {object}  common.Response{data=bool} "更新成功"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /v1/picture/edit [POST]
func EditPicture(c *gin.Context) {
	updateReq := reqPicture.PictureUpdateRequest{}
	c.ShouldBind(&updateReq)
	if updateReq.ID <= 0 {
		common.BaseResponse(c, false, "参数错误", ecode.PARAMS_ERROR)
		return
	}
	//校验是否本人或管理员操作
	user, _ := sUser.GetLoginUser(c)
	if user == nil || !(updateReq.ID == user.ID) && !(user.UserRole == "admin") {
		common.BaseResponse(c, false, "无权限", ecode.NO_AUTH_ERROR)
		return
	}
	//更新操作，参数校验等在service层完成
	if err := sPicture.UpdatePicture(&updateReq); err != nil {
		common.BaseResponse(c, false, err.Msg, err.Code)
		return
	}
	common.Success(c, true)
}

// ListPictureTagCategory godoc
// @Summary      获取图片的标签和分类（固定）
// @Tags         picture
// @Accept       json
// @Produce      json
// @Success      200  {object}  common.Response{data=resPicture.PictureTagCategory} "获取成功"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /v1/picture/tag_category [GET]
func ListPictureTagCategory(c *gin.Context) {
	tagCate := resPicture.PictureTagCategory{
		TagList:      []string{"热门", "搞笑", "生活", "高清", "艺术", "校园", "背景", "简历", "创意"},
		CategoryList: []string{"模板", "电商", "表情包", "素材", "海报"},
	}
	common.Success(c, tagCate)
}
