package controller

import (
	"chg/internal/common"
	"chg/internal/consts"
	"chg/internal/ecode"
	reqSpace "chg/internal/model/request/space"
	resSpace "chg/internal/model/response/space"
	"chg/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func dumb1() {
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
// @Param		request body reqSpace.SpaceEditRequest true "需要更新的空间信息"
// @Success      200  {object}  common.Response{data=bool} "更新成功"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /v1/space/edit [POST]
func EditSpace(c *gin.Context) {
	updateReq := reqSpace.SpaceEditRequest{}
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
	if err := sSpace.EditSpace(&updateReq, loginUser); err != nil {
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

// GetSpaceVOById godoc
// @Summary      获取当个空间的视图信息
// @Tags         space
// @Accept       json
// @Produce      json
// @Param		id query string true "空间的ID"
// @Success      200  {object}  common.Response{data=resSpace.SpaceVO} "获取成功"
// @Failure      400  {object}  common.Response "查询失败，详情见响应中的code"
// @Router       /v1/space/get/vo [GET]
func GetSpaceVOById(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	if id <= 0 {
		common.BaseResponse(c, nil, "参数错误", ecode.PARAMS_ERROR)
		return
	}
	space, err := sSpace.GetSpaceById(id)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	common.Success(c, *space)
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

// ListSpaceLevel godoc
// @Summary      获取所有的空间等级信息
// @Tags         space
// @Produce      json
// @Success      200  {object}  common.Response{data=[]resSpace.SpaceLevelResponse} "返回所有空间等级信息数组"
// @Failure      400  {object}  common.Response "查询失败，详情见响应中的code"
// @Router       /v1/space/list/level [GET]
func ListSpaceLevel(c *gin.Context) {
	res := []resSpace.SpaceLevelResponse{}
	for i := consts.FirstSpaceLevel; i <= consts.LastSpaceLevel; i++ {
		spaceLevel := consts.GetSpaceLevelByValue(i)
		res = append(res, resSpace.SpaceLevelResponse{
			Value:    i,
			Text:     spaceLevel.Text,
			MaxCount: spaceLevel.MaxCount,
			MaxSize:  spaceLevel.MaxSize,
		})
	}
	common.Success(c, res)
}
