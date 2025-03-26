package controller

import (
	"chg/internal/common"
	"chg/internal/ecode"
	reqSpace "chg/internal/model/request/space"
	resSpace "chg/internal/model/response/space"
	"chg/internal/service"

	"github.com/gin-gonic/gin"
)

func dumb() {
	temp := resSpace.ListSpaceResponse{}
	_ = temp
}

// 获取一个spaceService单例
var sSpace *service.SpaceService = service.NewSpaceService()

// UpdateSpace godoc
// @Summary      更新空间「管理员」
// @Description  若空间不存在，则返回false
// @Tags         space
// @Accept       json
// @Produce      json
// @Param		request body reqSpace.SpaceUpdateRequest true "需要更新的空间信息"
// @Success      200  {object}  common.Response{data=bool} "更新成功"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /v1/space/update [POST]
func UpdateSpace(c *gin.Context) {
	updateReq := reqSpace.SpaceUpdateRequest{}
	if err := c.ShouldBind(&updateReq); err != nil {
		common.BaseResponse(c, false, "参数绑定失败", ecode.PARAMS_ERROR)
		return
	}
	if updateReq.ID <= 0 {
		common.BaseResponse(c, false, "参数错误", ecode.PARAMS_ERROR)
		return
	}
	//更新操作，参数校验，权限校验等在service层完成
	loginUser, _ := sUser.GetLoginUser(c)
	if err := sSpace.UpdateSpace(&updateReq, loginUser); err != nil {
		common.BaseResponse(c, false, err.Msg, err.Code)
		return
	}
	common.Success(c, true)
}

// EditSpace godoc
// @Summary      编辑空间
// @Description  若空间不存在，则返回false
// @Tags         space
// @Accept       json
// @Produce      json
// @Param		request body reqSpace.SpaceUpdateRequest true "需要更新的空间信息"
// @Success      200  {object}  common.Response{data=bool} "更新成功"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /v1/space/edit [POST]
func EditSpace(c *gin.Context) {
	updateReq := reqSpace.SpaceUpdateRequest{}
	if err := c.ShouldBind(&updateReq); err != nil {
		common.BaseResponse(c, false, "参数绑定失败", ecode.PARAMS_ERROR)
		return
	}
	if updateReq.ID <= 0 {
		common.BaseResponse(c, false, "参数错误", ecode.PARAMS_ERROR)
		return
	}
	//更新操作，参数校验，权限校验等在service层完成
	loginUser, _ := sUser.GetLoginUser(c)
	if loginUser == nil {
		common.BaseResponse(c, false, "未登录", ecode.NOT_LOGIN_ERROR)
		return
	}
	if err := sSpace.UpdateSpace(&updateReq, loginUser); err != nil {
		common.BaseResponse(c, false, err.Msg, err.Code)
		return
	}
	common.Success(c, true)
}

// ListSpaceByPage godoc
// @Summary      分页获取一系列空间信息「管理员」
// @Tags         space
// @Accept       json
// @Produce      json
// @Param		request body reqSpace.SpaceQueryRequest true "需要查询的空间信息字段"
// @Success      200  {object}  common.Response{data=resSpace.ListSpaceResponse} "查询成功"
// @Failure      400  {object}  common.Response "查询失败，详情见响应中的code"
// @Router       /v1/space/list/page [POST]
func ListSpaceByPage(c *gin.Context) {
	queryReq := reqSpace.SpaceQueryRequest{}
	if err := c.ShouldBind(&queryReq); err != nil {
		common.BaseResponse(c, nil, "参数绑定失败", ecode.PARAMS_ERROR)
		return
	}
	//获取分页查询对象
	pics, err := sSpace.ListSpaceByPage(&queryReq)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	common.Success(c, *pics)
}

// ListSpaceVOByPage godoc
// @Summary      分页获取一系列空间视图信息
// @Tags         space
// @Accept       json
// @Produce      json
// @Param		request body reqSpace.SpaceQueryRequest true "需要查询的空间信息字段"
// @Success      200  {object}  common.Response{data=resSpace.ListSpaceVOResponse} "查询成功"
// @Failure      400  {object}  common.Response "查询失败，详情见响应中的code"
// @Router       /v1/space/list/page/vo [POST]
func ListSpaceVOByPage(c *gin.Context) {
	queryReq := reqSpace.SpaceQueryRequest{}
	if err := c.ShouldBind(&queryReq); err != nil {
		common.BaseResponse(c, nil, "参数绑定失败", ecode.PARAMS_ERROR)
		return
	}
	//获取分页查询对象
	pics, err := sSpace.ListSpaceVOByPage(&queryReq)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	common.Success(c, *pics)
}

// AddSpace godoc
// @Summary      增加空间「需要登录」
// @Tags         space
// @Accept       json
// @Produce      json
// @Param		request body reqSpace.SpaceAddRequest true "需要增加的空间信息字段"
// @Success      200  {object}  common.Response{data=string} "返回空间ID，字符串格式"
// @Failure      400  {object}  common.Response "查询失败，详情见响应中的code"
// @Router       /v1/space/add [POST]
func AddSpace(c *gin.Context) {
	queryReq := reqSpace.SpaceAddRequest{}
	if err := c.ShouldBind(&queryReq); err != nil {
		common.BaseResponse(c, nil, "参数绑定失败", ecode.PARAMS_ERROR)
		return
	}
	loginUser, _ := sUser.GetLoginUser(c)
	//获取分页查询对象
	spaceId, err := sSpace.AddSpace(&queryReq, loginUser)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	common.Success(c, spaceId)
}
