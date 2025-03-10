package controller

import (
	"chg/internal/common"
	"chg/internal/ecode"
	"chg/internal/model/request"
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
	}
	if id, err := s.UserRegister(uReg.UserAccount, uReg.UserPassword, uReg.CheckPassword); err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
	} else {
		common.Success(c, id)
	}
}
