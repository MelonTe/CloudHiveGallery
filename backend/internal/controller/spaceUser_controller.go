package controller

import (
	"chg/internal/common"
	"chg/internal/ecode"
	"chg/internal/model/entity"
	reqSpaceUser "chg/internal/model/request/spaceuser"
	resSpaceUser "chg/internal/model/response/spaceuser"
	"chg/internal/service"
	"chg/pkg/db"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func dump3() {
	temp := resSpaceUser.SpaceUserVO{}
	_ = temp
}

//接口的权限校验会使用Sa-token来实现

// 获取一个spaceService单例
var sSpaceUser *service.SpaceUserService = service.NewSpaceUserService()

// AddSpaceUser godoc
// @Summary      增加成员到空间
// @Tags         spaceUser
// @Accept       json
// @Produce      json
// @Param		request body reqSpaceUser.SpaceUserAddRequest true "成员的ID和空间ID，以及添加的成员角色"
// @Success      200  {object}  common.Response{data=string} "返回空间成员表的数据ID，字符串格式"
// @Failure      400  {object}  common.Response "查询失败，详情见响应中的code"
// @Router       /v1/spaceUser/add [POST]
func AddSpaceUser(c *gin.Context) {
	Req := reqSpaceUser.SpaceUserAddRequest{}
	if err := c.ShouldBind(&Req); err != nil {
		common.BaseResponse(c, nil, "参数绑定失败", ecode.PARAMS_ERROR)
		return
	}
	id, err := sSpaceUser.AddSpaceUser(&Req)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	common.Success(c, fmt.Sprintf("%d", id))
}

// DeleteSpaceUser godoc
// @Summary      从空间移除成员
// @Tags         spaceUser
// @Accept       json
// @Produce      json
// @Param		request body reqSpaceUser.SpaceUserRemoveRequest true "空间成员表中数据的ID"
// @Success      200  {object}  common.Response{data=bool} "返回成功与否"
// @Failure      400  {object}  common.Response "查询失败，详情见响应中的code"
// @Router       /v1/spaceUser/delete [POST]
func DeleteSpaceUser(c *gin.Context) {
	req := reqSpaceUser.SpaceUserRemoveRequest{}
	if err := c.ShouldBind(&req); err != nil {
		common.BaseResponse(c, nil, "参数绑定失败", ecode.PARAMS_ERROR)
		return
	}
	err := sSpaceUser.RemoveSpaceUserById(req.ID)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	common.Success(c, true)
}

// GetSpaceUser godoc
// @Summary      查询某个成员在某个空间的信息
// @Tags         spaceUser
// @Accept       json
// @Produce      json
// @Param		request body reqSpaceUser.SpaceUserQueryRequest true "必须携带spaceID和userID"
// @Success      200  {object}  common.Response{data=entity.SpaceUser} "返回成功与否"
// @Failure      400  {object}  common.Response "查询失败，详情见响应中的code"
// @Router       /v1/spaceUser/get [POST]
func GetSpaceUser(c *gin.Context) {
	req := reqSpaceUser.SpaceUserQueryRequest{}
	if err := c.ShouldBind(&req); err != nil {
		common.BaseResponse(c, nil, "参数绑定失败", ecode.PARAMS_ERROR)
		return
	}
	query, err := sSpaceUser.GetQueryWrapper(db.LoadDB(), &req)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	spaceUser := &entity.SpaceUser{}
	originErr := query.First(spaceUser).Error
	if originErr != nil {
		if originErr == gorm.ErrRecordNotFound {
			common.BaseResponse(c, nil, "没有找到该空间成员", ecode.PARAMS_ERROR)
			return
		}
		common.BaseResponse(c, nil, "查询失败", ecode.SYSTEM_ERROR)
		return
	}
	common.Success(c, *spaceUser)
}

// ListSpaceUser godoc
// @Summary      查询成员信息列表
// @Tags         spaceUser
// @Accept       json
// @Produce      json
// @Param		request body reqSpaceUser.SpaceUserQueryRequest true "可以携带的参数"
// @Success      200  {object}  common.Response{data=[]resSpaceUser.SpaceUserVO} "返回详细数据"
// @Failure      400  {object}  common.Response "查询失败，详情见响应中的code"
// @Router       /v1/spaceUser/list [POST]
func ListSpaceUser(c *gin.Context) {
	req := reqSpaceUser.SpaceUserQueryRequest{}
	if err := c.ShouldBind(&req); err != nil {
		common.BaseResponse(c, nil, "参数绑定失败", ecode.PARAMS_ERROR)
		return
	}
	err, spaceVOList := sSpaceUser.ListSpaceUserVO(&req)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	common.Success(c, spaceVOList)
}

// EditSpaceUser godoc
// @Summary      编辑成员信息
// @Tags         spaceUser
// @Accept       json
// @Produce      json
// @Param		request body reqSpaceUser.SpaceUserEditRequest true "记录的ID和需要调整的权限"
// @Success      200  {object}  common.Response{data=bool} "编辑成功与否"
// @Failure      400  {object}  common.Response "查询失败，详情见响应中的code"
// @Router       /v1/spaceUser/edit [POST]
func EditSpaceUser(c *gin.Context) {
	req := reqSpaceUser.SpaceUserEditRequest{}
	if err := c.ShouldBind(&req); err != nil {
		common.BaseResponse(c, nil, "参数绑定失败", ecode.PARAMS_ERROR)
		return
	}
	suc, err := sSpaceUser.EditSpaceUser(&req)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	common.Success(c, suc)
}

// ListMyTeamSpace godoc
// @Summary      查询我加入的团队空间列表
// @Tags         spaceUser
// @Accept       json
// @Produce      json
// @Success      200  {object}  common.Response{data=[]resSpaceUser.SpaceUserVO} "返回详细数据"
// @Failure      400  {object}  common.Response "查询失败，详情见响应中的code"
// @Router       /v1/spaceUser/list/my [POST]
func ListMyTeamSpace(c *gin.Context) {
	loginUser, err := sUser.GetLoginUser(c)
	if err != nil {
		common.BaseResponse(c, nil, "获取登录用户失败", ecode.NOT_LOGIN_ERROR)
		return
	}
	req := &reqSpaceUser.SpaceUserQueryRequest{
		UserID: loginUser.ID,
	}
	err, spaceVOList := sSpaceUser.ListSpaceUserVO(req)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	common.Success(c, spaceVOList)
}
