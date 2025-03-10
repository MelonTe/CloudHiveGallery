package controller

import (
	"chg/internal/common"
	"chg/internal/ecode"
	"chg/internal/model/request"
	"chg/internal/model/response"
	"chg/internal/service"

	"github.com/gin-gonic/gin"
)

// 接口前缀为/user
// param	用空格分隔的参数。param name,param type,data type,is mandatory?,comment attribute(optional)
// 获取一个userservice单例
var s *service.UserService = service.NewUserService()

// UserRegister godoc
// @Summary      注册用户
// @Description  根据账号密码进行注册
// @Tags         user
// @Accept       json
// @Produce      json
// @Param		request body request.UserRegsiterRequest true "用户请求注册参数"
// @Success      200  {object}  common.Response{data=string} "注册成功，返回注册用户的ID"
// @Failure      400  {object}  common.Response "注册失败，详情见响应中的code"
// @Router       /v1/user/register [POST]
func UserRegister(c *gin.Context) {
	//使用shouldbind绑定参数，参数不可复用
	var uReg request.UserRegsiterRequest
	if err := c.ShouldBind(&uReg); err != nil {
		common.BaseResponse(c, nil, "参数绑定错误", ecode.PARAMS_ERROR)
		return
	}
	if id, err := s.UserRegister(uReg.UserAccount, uReg.UserPassword, uReg.CheckPassword); err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	} else {
		common.Success(c, id)
		return
	}
}

// UserLogin godoc
// @Summary      用户登录
// @Description  根据账号密码进行登录
// @Tags         user
// @Accept       json
// @Produce      json
// @Param		request body request.UserLoginRequest true "用户登录请求参数"
// @Success      200  {object}  common.Response{data=response.UserLoginVO} "登录成功，返回用户视图"
// @Failure      400  {object}  common.Response "登录失败，详情见响应中的code"
// @Router       /v1/user/login [POST]
func UserLogin(c *gin.Context) {
	var uLog request.UserLoginRequest
	if err := c.ShouldBind(&uLog); err != nil {
		common.BaseResponse(c, nil, "参数绑定错误", ecode.PARAMS_ERROR)
		return
	}
	if userVO, err := s.UserLogin(c, uLog.UserAccount, uLog.UserPassword); err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	} else {
		common.Success(c, *userVO)
		return
	}
}

// GetLoginUser godoc
// @Summary      获取登录的用户信息
// @Tags         user
// @Produce      json
// @Success      200  {object}  common.Response{data=response.UserLoginVO} "获取用户视图成功"
// @Failure      400  {object}  common.Response "获取失败，详情见响应中的code"
// @Router       /v1/user/get/login [GET]
func GetLoginUser(c *gin.Context) {
	user, err := s.GetLoginUser(c)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	//返回用户视图
	common.Success(c, *response.GetUserLoginVO(*user))
}

// UserLogout godoc
// @Summary      执行用户注销（退出）
// @Tags         user
// @Produce      json
// @Success      200  {object}  common.Response{data=bool} "退出成功"
// @Failure      400  {object}  common.Response "注册失败，详情见响应中的code"
// @Router       /v1/user/logout [POST]
func UserLogout(c *gin.Context) {
	suc, err := s.UserLogout(c)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	common.Success(c, suc)
}
