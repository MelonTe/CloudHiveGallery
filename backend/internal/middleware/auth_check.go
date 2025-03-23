package middleware

import (
	"chg/internal/common"
	"chg/internal/consts"
	"chg/internal/ecode"
	"chg/internal/service"

	"github.com/gin-gonic/gin"
)

// 权限校验中间件，只有拥有权限的用户才允许执行接下来的服务
func AuthCheck(userService *service.UserService, mustRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取当前登录对象
		loginUser, err := userService.GetLoginUser(c)
		if err != nil {
			//未登录或者出错
			common.BaseResponse(c, nil, err.Msg, err.Code)
			c.Abort()
			return
		}
		if mustRole == "" {
			//不需要权限，放行
			c.Next()
			return
		}
		//需要权限：
		//校验用户角色，管理员权限只能管理员用，普通权限管理员和用户可用
		if (mustRole == consts.ADMIN_ROLE && loginUser.UserRole != consts.ADMIN_ROLE) || loginUser.UserRole == "" {
			common.BaseResponse(c, nil, "无权限", ecode.NO_AUTH_ERROR)
			c.Abort()
			return
		}
		//权限通过，放行
		c.Next()
	}
}

// 登录校验中间件，只有登录的用户才允许执行接下来的服务
func LoginCheck(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取当前登录对象
		_, err := userService.GetLoginUser(c)
		if err != nil {
			//未登录或者出错
			common.BaseResponse(c, nil, err.Msg, err.Code)
			c.Abort()
			return
		}
		//放行
		c.Next()
	}
}
