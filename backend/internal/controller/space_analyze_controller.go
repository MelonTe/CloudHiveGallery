package controller

import (
	"chg/internal/common"
	"chg/internal/ecode"
	reqSpaceAnalyze "chg/internal/model/request/space/analyze"
	resSpaceAnalyze "chg/internal/model/response/space/analyze"
	"chg/internal/service"

	"github.com/gin-gonic/gin"
)

var sSpaceAnalyze *service.SpaceAnalyzeService = service.NewSpaceAnalyzeService()

func dump2() {
	temp := resSpaceAnalyze.SpaceUsageAnalyzeResponse{}
	_ = temp
}

// GetSpaceUsageAnalyze godoc
// @Summary      获取空间使用分析「登录校验」
// @Tags         space/analyze
// @Accept       json
// @Produce      json
// @Param		request body reqSpaceAnalyze.SpaceUsageAnalyzeRequest true "查询条件"
// @Success      200  {object}  common.Response{data=resSpaceAnalyze.SpaceUsageAnalyzeResponse} "查询成功"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /v1/space/analyze/usage [POST]
func GetSpaceUsageAnalyze(c *gin.Context) {
	req := reqSpaceAnalyze.SpaceUsageAnalyzeRequest{}
	if err := c.ShouldBind(&req); err != nil {
		common.BaseResponse(c, false, "参数绑定失败", ecode.PARAMS_ERROR)
		return
	}
	loginUser, _ := sUser.GetLoginUser(c)
	res, err := sSpaceAnalyze.GetSpaceUsageAnalyze(&req, loginUser)
	if err != nil {
		common.BaseResponse(c, false, err.Msg, err.Code)
		return
	}
	common.Success(c, *res)
}

// GetSpaceCategoryAnalyze godoc
// @Summary      获取空间图片分类分析「登录校验」
// @Tags         space/analyze
// @Accept       json
// @Produce      json
// @Param		request body reqSpaceAnalyze.SpaceCategoryAnalyzeRequest true "查询条件"
// @Success      200  {object}  common.Response{data=[]resSpaceAnalyze.SpaceCategoryAnalyzeResponse} "查询成功"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /v1/space/analyze/category [POST]
func GetSpaceCategoryAnalyze(c *gin.Context) {
	req := reqSpaceAnalyze.SpaceCategoryAnalyzeRequest{}
	if err := c.ShouldBind(&req); err != nil {
		common.BaseResponse(c, false, "参数绑定失败", ecode.PARAMS_ERROR)
		return
	}
	loginUser, _ := sUser.GetLoginUser(c)
	res, err := sSpaceAnalyze.GetSpaceCategoryAnalyze(&req, loginUser)
	if err != nil {
		common.BaseResponse(c, false, err.Msg, err.Code)
		return
	}
	common.Success(c, res)
}

// GetSpaceTagAnalyze godoc
// @Summary      获取空间标签出现量分析「登录校验」
// @Tags         space/analyze
// @Accept       json
// @Produce      json
// @Param		request body reqSpaceAnalyze.SpaceTagAnalyzeRequest true "查询条件"
// @Success      200  {object}  common.Response{data=[]resSpaceAnalyze.SpaceTagAnalyzeResponse} "查询成功"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /v1/space/analyze/tag [POST]
func GetSpaceTagAnalyze(c *gin.Context) {
	req := reqSpaceAnalyze.SpaceTagAnalyzeRequest{}
	if err := c.ShouldBind(&req); err != nil {
		common.BaseResponse(c, false, "参数绑定失败", ecode.PARAMS_ERROR)
		return
	}
	loginUser, _ := sUser.GetLoginUser(c)
	res, err := sSpaceAnalyze.GetSpaceTagAnalyze(&req, loginUser)
	if err != nil {
		common.BaseResponse(c, false, err.Msg, err.Code)
		return
	}
	common.Success(c, res)
}

// GetSpaceSizeAnalyze godoc
// @Summary      获取空间图片大小范围统计分析「登录校验」
// @Tags         space/analyze
// @Accept       json
// @Produce      json
// @Param		request body reqSpaceAnalyze.SpaceSizeAnalyzeRequest true "查询条件"
// @Success      200  {object}  common.Response{data=[]resSpaceAnalyze.SpaceSizeAnalyzeResponse} "查询成功"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /v1/space/analyze/size [POST]
func GetSpaceSizeAnalyze(c *gin.Context) {
	req := reqSpaceAnalyze.SpaceSizeAnalyzeRequest{}
	if err := c.ShouldBind(&req); err != nil {
		common.BaseResponse(c, false, "参数绑定失败", ecode.PARAMS_ERROR)
		return
	}
	loginUser, _ := sUser.GetLoginUser(c)
	res, err := sSpaceAnalyze.GetSpaceSizeAnalyze(&req, loginUser)
	if err != nil {
		common.BaseResponse(c, false, err.Msg, err.Code)
		return
	}
	common.Success(c, res)
}

// GetSpaceUserAnalyze godoc
// @Summary      获取用户上传图片统计分析，支持分析特定用户「登录校验」
// @Tags         space/analyze
// @Accept       json
// @Produce      json
// @Param		request body reqSpaceAnalyze.SpaceUserAnalyzeRequest true "查询条件"
// @Success      200  {object}  common.Response{data=[]resSpaceAnalyze.SpaceUserAnalyzeResponse} "查询成功"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /v1/space/analyze/user [POST]
func GetSpaceUserAnalyze(c *gin.Context) {
	req := reqSpaceAnalyze.SpaceUserAnalyzeRequest{}
	if err := c.ShouldBind(&req); err != nil {
		common.BaseResponse(c, false, "参数绑定失败", ecode.PARAMS_ERROR)
		return
	}
	loginUser, _ := sUser.GetLoginUser(c)
	res, err := sSpaceAnalyze.GetSpaceUserAnalyze(&req, loginUser)
	if err != nil {
		common.BaseResponse(c, false, err.Msg, err.Code)
		return
	}
	common.Success(c, res)
}

// GetSpaceRankAnalyze godoc
// @Summary      获取空间使用情况排名「管理员」
// @Tags         space/analyze
// @Accept       json
// @Produce      json
// @Param		request body reqSpaceAnalyze.SpaceRankAnalyzeRequest true "查询条件"
// @Success      200  {object}  common.Response{data=[]entity.Space} "查询成功"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /v1/space/analyze/rank [POST]
func GetSpaceRankAnalyze(c *gin.Context) {
	req := reqSpaceAnalyze.SpaceRankAnalyzeRequest{}
	if err := c.ShouldBind(&req); err != nil {
		common.BaseResponse(c, false, "参数绑定失败", ecode.PARAMS_ERROR)
		return
	}
	loginUser, _ := sUser.GetLoginUser(c)
	res, err := sSpaceAnalyze.GetSpaceRankAnalyze(&req, loginUser)
	if err != nil {
		common.BaseResponse(c, false, err.Msg, err.Code)
		return
	}
	common.Success(c, res)
}
