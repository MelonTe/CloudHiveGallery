package controller

import (
	"chg/internal/common"
	"chg/internal/consts"
	"chg/internal/ecode"
	reqPicture "chg/internal/model/request/picture"
	resPicture "chg/internal/model/response/picture"
	"chg/internal/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// 获取一个pictureService单例
var sPicture *service.PictureService = service.NewPictureService()

// UploadPicture godoc
// @Summary      上传图片接口「需要登录校验」
// @Description  根据是否存在ID来上传图片或者修改图片信息，返回图片信息视图
// @Tags         picture
// @Accept       mpfd
// @Produce      json
// @Param        file formData file true "图片"
// @Param        id formData string false "图片的ID，非必需"
// @Param        spaceId formData string false "图片的上传空间ID，非必需"
// @Success      200  {object}  common.Response{data=resPicture.PictureVO} "上传成功，返回图片信息视图"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /v1/picture/upload [POST]
func UploadPicture(c *gin.Context) {
	file, _ := c.FormFile("file")
	// 手动解析表单参数
	id, _ := strconv.ParseUint(c.PostForm("id"), 10, 64)
	spaceId, _ := strconv.ParseUint(c.PostForm("spaceId"), 10, 64)
	picReq := &reqPicture.PictureUploadRequest{
		ID:      id,      // 获取 id
		SpaceID: spaceId, // 获取 spaceId
	}
	loginUser, err := sUser.GetLoginUser(c)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	picVO, err := sPicture.UploadPicture(file, picReq, loginUser)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	common.Success(c, *picVO)
}

// UploadPictureByUrl godoc
// @Summary      根据URL上传图片接口「需要登录校验」
// @Tags         picture
// @Accept       json
// @Produce      json
// @Param        request body reqPicture.PictureUploadRequest true "图片URL"
// @Success      200  {object}  common.Response{data=resPicture.PictureVO} "上传成功，返回图片信息视图"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /v1/picture/upload/url [POST]
func UploadPictureByUrl(c *gin.Context) {
	picReq := &reqPicture.PictureUploadRequest{}
	c.ShouldBind(picReq)
	loginUser, err := sUser.GetLoginUser(c)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	//若picUrl包含了?解析参数，需要去掉
	if idx := strings.LastIndex(picReq.FileUrl, "?"); idx != -1 {
		picReq.FileUrl = picReq.FileUrl[:idx]
	}
	picVO, err := sPicture.UploadPicture(picReq.FileUrl, picReq, loginUser)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	common.Success(c, *picVO)
}

// UploadPictureByBatch godoc
// @Summary      批量抓取图片「管理员」
// @Tags         picture
// @Accept       json
// @Produce      json
// @Param        request body reqPicture.PictureUploadByBatchRequest true "图片的关键词"
// @Success      200  {object}  common.Response{data=int} "返回抓取图片数量"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /v1/picture/upload/batch [POST]
func UploadPictureByBatch(c *gin.Context) {
	picReq := &reqPicture.PictureUploadByBatchRequest{}
	c.ShouldBind(picReq)
	//一定能获取
	loginUser, _ := sUser.GetLoginUser(c)
	cnt, err := sPicture.UploadPictureByBatch(picReq, loginUser)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	common.Success(c, cnt)
}

// DeletePicture godoc
// @Summary      根据ID软删除图片「登录校验」
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
	user, _ := sUser.GetLoginUser(c)
	if err := sPicture.DeletePicture(user, &deleReq); err != nil {
		common.BaseResponse(c, false, err.Msg, err.Code)
		return
	}
	common.Success(c, true)
}

// UpdatePicture godoc
// @Summary      更新图片「登录校验」
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
	//获取登录用户，使用中间件保证可以获取到用户
	loginUser, _ := sUser.GetLoginUser(c)
	//更新操作，参数校验等在service层完成
	if err := sPicture.UpdatePicture(&updateReq, loginUser); err != nil {
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
	//空间权限校验
	if pic.SpaceID != 0 {
		loginUser, _ := sUser.GetLoginUser(c)
		if loginUser == nil || sPicture.CheckPictureAuth(loginUser, pic) != nil {
			common.BaseResponse(c, nil, "无权限", ecode.NO_AUTH_ERROR)
			return
		}
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
		common.BaseResponse(c, nil, "最多只允许获取20张/页", ecode.PARAMS_ERROR)
		return
	}
	//空间权限校验
	if queryReq.SpaceID != 0 {
		//私有空间
		loginUser, err := sUser.GetLoginUser(c)
		if err != nil {
			common.BaseResponse(c, nil, err.Msg, err.Code)
			return
		}
		space, err := sSpace.GetSpaceById(queryReq.SpaceID)
		if err != nil {
			common.BaseResponse(c, nil, err.Msg, err.Code)
			return
		}
		if space.UserID != loginUser.ID {
			common.BaseResponse(c, nil, "无权限", ecode.NO_AUTH_ERROR)
			return
		}
	} else {
		//公开图库
		//普通用户默认只允许查询过审图片
		if queryReq.ReviewStatus == nil {
			queryReq.ReviewStatus = new(int) //创建指针
		}
		*queryReq.ReviewStatus = consts.PASS
		queryReq.IsNullSpaceID = true
	}
	//获取分页查询对象
	pics, err := sPicture.ListPictureVOByPage(&queryReq)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	common.Success(c, *pics)
}

// ListPictureVOByPageWithCache godoc
// @Summary      带有缓存的分页获取一系列图片信息
// @Tags         picture
// @Accept       json
// @Produce      json
// @Param		request body reqPicture.PictureQueryRequest true "需要查询的页数、以及图片关键信息"
// @Success      200  {object}  common.Response{data=resPicture.ListPictureVOResponse} "查询成功"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /v1/picture/list/page/vo/cache [POST]
func ListPictureVOByPageWithCache(c *gin.Context) {
	queryReq := reqPicture.PictureQueryRequest{}
	c.ShouldBind(&queryReq)
	//限制爬虫
	if queryReq.PageSize > 20 {
		common.BaseResponse(c, nil, "最多只允许获取20张/页", ecode.PARAMS_ERROR)
		return
	}
	//空间权限校验
	if queryReq.SpaceID != 0 {
		//私有空间
		loginUser, err := sUser.GetLoginUser(c)
		if err != nil {
			common.BaseResponse(c, nil, err.Msg, err.Code)
			return
		}
		space, err := sSpace.GetSpaceById(queryReq.SpaceID)
		if err != nil {
			common.BaseResponse(c, nil, err.Msg, err.Code)
			return
		}
		if space.UserID != loginUser.ID {
			common.BaseResponse(c, nil, "无权限", ecode.NO_AUTH_ERROR)
			return
		}
	} else {
		//公开图库
		//普通用户默认只允许查询过审图片
		if queryReq.ReviewStatus == nil {
			queryReq.ReviewStatus = new(int) //创建指针
		}
		*queryReq.ReviewStatus = consts.PASS
		queryReq.IsNullSpaceID = true
	}
	//获取分页查询对象
	pics, err := sPicture.ListPictureVOByPageWithCache(&queryReq)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	common.Success(c, *pics)
}

// EditPicture godoc
// @Summary      编辑图片
// @Description  若图片不存在，则返回false
// @Tags         picture
// @Accept       json
// @Produce      json
// @Param		request body reqPicture.PictureEditRequest true "需要更新的图片信息"
// @Success      200  {object}  common.Response{data=bool} "更新成功"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /v1/picture/edit [POST]
func EditPicture(c *gin.Context) {
	//update和edit复用了同一个请求
	updateReq := reqPicture.PictureUpdateRequest{}
	c.ShouldBind(&updateReq)
	if updateReq.ID <= 0 {
		common.BaseResponse(c, false, "参数错误", ecode.PARAMS_ERROR)
		return
	}
	//校验是否本人或管理员操作
	user, _ := sUser.GetLoginUser(c)
	if user == nil {
		common.BaseResponse(c, false, "未登录", ecode.NOT_LOGIN_ERROR)
		return
	}
	//更新操作，参数校验和权限等在service层完成
	if err := sPicture.UpdatePicture(&updateReq, user); err != nil {
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

// DoPictureReview godoc
// @Summary      执行图片审核「管理员」
// @Tags         picture
// @Accept       json
// @Produce      json
// @Param		request body reqPicture.PictureReviewRequest true "审核图片所需信息"
// @Success      200  {object}  common.Response{data=bool} "审核更新成功"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /v1/picture/review [POST]
func DoPictureReview(c *gin.Context) {
	var req reqPicture.PictureReviewRequest
	c.ShouldBind(&req)
	//获取当前登录用户
	user, _ := sUser.GetLoginUser(c)
	if err := sPicture.DoPictureReview(&req, user); err != nil {
		common.BaseResponse(c, false, err.Msg, err.Code)
		return
	}
	common.Success(c, true)
}
