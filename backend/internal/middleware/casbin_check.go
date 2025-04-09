package middleware

import (
	"bytes"
	"chg/internal/common"
	"chg/internal/ecode"
	"chg/internal/service"
	"chg/pkg/casbin"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
)

//引入casbin的RBAC鉴权

// 中间件鉴权函数，必须需要登录
// Dom: 访问的资源域，对于公共图库是public，对于特定的空间是space，具体的空间ID会从请求中提取
// Obj: 访问的资源对象，目前有picture和spaceUser两种
// Act: 访问的行为，对于picture有upload/delete/view/edit，对于spaceUser有manage
func CasbinAuthCheck(Dom, Obj, Act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取一个用户Service
		userService := service.NewUserService()
		//获取当前登录对象
		loginUser, err := userService.GetLoginUser(c)
		if err != nil {
			//未登录或者出错
			common.BaseResponse(c, nil, err.Msg, err.Code)
			c.Abort()
			return
		}
		//获取sub，用户ID
		sub := "user_" + fmt.Sprintf("%d", loginUser.ID)

		//若Dom是space，则需要从请求中获取空间ID
		if Dom == "space" {
			//复制一个请求体
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			// 恢复请求体内容，供后续绑定等操作使用
			c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			var OriginBodyMap map[string]interface{}
			if OriginErr := json.Unmarshal(bodyBytes, &OriginBodyMap); OriginErr != nil {
				common.BaseResponse(c, nil, "请求体解析失败", ecode.SYSTEM_ERROR)
				c.Abort()
				return
			}
			//将OriginBodyMap中，为string类型的ID转化为uint64，存入bodymap
			bodyMap := make(map[string]interface{})
			for k, v := range OriginBodyMap {
				if k == "Id" || k == "spaceId" {
					//将ID转化为uint64类型
					if id, ok := v.(string); ok {
						idUint64, err := strconv.ParseUint(id, 10, 64)
						if err == nil {
							bodyMap[k] = idUint64
						}
					}
				}
			}
			//尝试直接获取spaceId
			spaceId, ok := bodyMap["spaceId"]
			if ok {
				Dom = fmt.Sprintf("%s_%v", Dom, spaceId)
			} else {
				//若Obj是spaceUser，那么需要从数据库根据Id找到对应的记录获取spaceId
				if Obj == "spaceUser" {
					//获取空间成员ID
					spaceUserId, ok := bodyMap["Id"]
					if !ok {
						common.BaseResponse(c, nil, "请求体中缺少spaceUser ID", ecode.SYSTEM_ERROR)
						c.Abort()
						return
					}
					//获取空间成员Service
					spaceUserService := service.NewSpaceUserService()
					//根据ID查询空间成员信息
					spaceUserInfo, err := spaceUserService.GetSpaceUserById(spaceUserId.(uint64))
					if err != nil {
						common.BaseResponse(c, nil, "获取空间成员信息失败", ecode.SYSTEM_ERROR)
						c.Abort()
						return
					}
					//获取空间ID
					Dom = fmt.Sprintf("%s_%d", Dom, spaceUserInfo.SpaceID)
				} else {
					common.BaseResponse(c, nil, "请求体中缺少space ID", ecode.SYSTEM_ERROR)
					c.Abort()
					return
				}
			}
		}
		//获取casbin鉴权中间件
		casMethod := casbin.LoadCasbinMethod()
		//判断是否有权限
		ok, originErr := casMethod.Enforcer.Enforce(sub, Dom, Obj, Act)
		if originErr != nil {
			//权限校验出错
			common.BaseResponse(c, nil, "权限校验出错", ecode.SYSTEM_ERROR)
			c.Abort()
			return
		}
		if !ok {
			//没有权限
			common.BaseResponse(c, nil, "没有权限", ecode.NO_AUTH_ERROR)
			c.Abort()
			return
		}
		//权限通过，放行
		c.Next()
	}
}
