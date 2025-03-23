package controller

import (
	"chg/internal/common"
	"chg/internal/ecode"
	"chg/internal/model/entity"
	reqUser "chg/internal/model/request/user"
	resUser "chg/internal/model/response/user"
	"chg/internal/service"
	"strconv"

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
// @Param		request body reqUser.UserRegsiterRequest true "用户请求注册参数"
// @Success      200  {object}  common.Response{data=string} "注册成功，返回注册用户的ID"
// @Failure      400  {object}  common.Response "注册失败，详情见响应中的code"
// @Router       /v1/user/register [POST]
func UserRegister(c *gin.Context) {
	//使用shouldbind绑定参数，参数不可复用
	var uReg reqUser.UserRegsiterRequest
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
// @Param		request body reqUser.UserLoginRequest true "用户登录请求参数"
// @Success      200  {object}  common.Response{data=resUser.UserLoginVO} "登录成功，返回用户视图"
// @Failure      400  {object}  common.Response "登录失败，详情见响应中的code"
// @Router       /v1/user/login [POST]
func UserLogin(c *gin.Context) {
	var uLog reqUser.UserLoginRequest
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
// @Success      200  {object}  common.Response{data=resUser.UserLoginVO} "获取用户视图成功"
// @Failure      400  {object}  common.Response "获取失败，详情见响应中的code"
// @Router       /v1/user/get/login [GET]
func GetLoginUser(c *gin.Context) {
	user, err := s.GetLoginUser(c)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	//返回用户视图
	common.Success(c, *resUser.GetUserLoginVO(*user))
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

// AddUser godoc
// @Summary      创建一个用户「管理员」
// @Description  默认密码为12345678
// @Tags         user
// @Accept       json
// @Produce      json
// @Param		request body reqUser.UserAddRequest true "用户添加申请参数"
// @Success      200  {object}  common.Response{data=string} "添加成功，返回添加用户的ID"
// @Failure      400  {object}  common.Response "添加失败，详情见响应中的code"
// @Router       /v1/user/add [POST]
func AddUser(c *gin.Context) {
	//使用shouldbind绑定参数，参数不可复用
	var uReg reqUser.UserAddRequest
	if err := c.ShouldBind(&uReg); err != nil {
		common.BaseResponse(c, nil, "参数绑定错误", ecode.PARAMS_ERROR)
		return
	}
	//简单逻辑，不放在服务里面写了
	user := &entity.User{
		UserAccount:  uReg.UserAccount,
		UserName:     uReg.UserName,
		UserRole:     uReg.UserRole,
		UserPassword: service.GetEncryptPassword("12345678"),
		UserProfile:  uReg.UserProfile,
		UserAvatar:   uReg.UserAvatar,
	}
	if err := s.UserRepo.CreateUser(user); err != nil {
		common.BaseResponse(c, nil, "数据库错误，注册失败", ecode.SYSTEM_ERROR)
		return
	}
	common.Success(c, user.ID)
}

// GetUserById godoc
// @Summary      根据ID获取用户「管理员」
// @Tags         user
// @Accept       json
// @Produce      json
// @Param		id query string true "用户的ID"
// @Success      200  {object}  common.Response{data=entity.User} "查询成功，返回用户的所有信息"
// @Failure      400  {object}  common.Response "查询失败，详情见响应中的code"
// @Router       /v1/user/get [GET]
func GetUserById(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	if id <= 0 {
		common.BaseResponse(c, nil, "参数错误", ecode.PARAMS_ERROR)
		return
	}
	user, err := s.UserRepo.FindById(id)
	if err != nil {
		common.BaseResponse(c, nil, "数据库错误，查询失败", ecode.SYSTEM_ERROR)
		return
	}
	if user == nil {
		common.BaseResponse(c, nil, "用户不存在", ecode.NOT_FOUND_ERROR)
		return
	}
	common.Success(c, *user)
}

// GetUserVOById godoc
// @Summary      根据ID获取简略信息用户
// @Tags         user
// @Accept       json
// @Produce      json
// @Param		id query string true "用户的ID"
// @Success      200  {object}  common.Response{data=resUser.UserVO} "查询成功，返回用户的脱敏信息"
// @Failure      400  {object}  common.Response "查询失败，详情见响应中的code"
// @Router       /v1/user/get/vo [GET]
func GetUserVOById(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	if id <= 0 {
		common.BaseResponse(c, nil, "参数错误", ecode.PARAMS_ERROR)
		return
	}
	user, err := s.UserRepo.FindById(id)
	if err != nil {
		common.BaseResponse(c, nil, "数据库错误，查询失败", ecode.SYSTEM_ERROR)
		return
	}
	if user == nil {
		common.BaseResponse(c, nil, "用户不存在", ecode.NOT_FOUND_ERROR)
		return
	}
	u := resUser.GetUserVO(*user)
	common.Success(c, u)
}

// DeleteUser godoc
// @Summary      根据ID软删除用户「管理员」
// @Tags         user
// @Accept       json
// @Produce      json
// @Param		request body common.DeleteRequest true "用户的ID"
// @Success      200  {object}  common.Response{data=bool} "删除成功"
// @Failure      400  {object}  common.Response "删除失败，详情见响应中的code"
// @Router       /v1/user/delete [POST]
func DeleteUser(c *gin.Context) {
	deleReq := common.DeleteRequest{}
	c.ShouldBind(&deleReq)
	if deleReq.Id <= 0 {
		common.BaseResponse(c, false, "删除失败，参数错误", ecode.PARAMS_ERROR)
		return
	}
	if suc, err := s.RemoveById(deleReq.Id); err != nil {
		common.BaseResponse(c, suc, err.Msg, err.Code)
		return
	}
	common.Success(c, true)
}

// UpdateUser godoc
// @Summary      更新用户信息「管理员」
// @Description  若用户不存在，则返回失败
// @Tags         user
// @Accept       json
// @Produce      json
// @Param		request body reqUser.UserUpdateRequest true "需要更新的用户信息"
// @Success      200  {object}  common.Response{data=bool} "更新成功"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /v1/user/update [POST]
func UpdateUser(c *gin.Context) {
	updateReq := reqUser.UserUpdateRequest{}
	c.ShouldBind(&updateReq)
	if updateReq.ID <= 0 {
		common.BaseResponse(c, false, "参数错误", ecode.PARAMS_ERROR)
		return
	}
	u := entity.User{
		ID:          updateReq.ID,
		UserName:    updateReq.UserName,
		UserAvatar:  updateReq.UserAvatar,
		UserProfile: updateReq.UserProfile,
		UserRole:    updateReq.UserRole,
	}
	if err := s.UpdateUser(&u); err != nil {
		common.BaseResponse(c, false, err.Msg, err.Code)
		return
	}
	common.Success(c, true)
}

// ListUserVOByPage godoc
// @Summary      分页获取一系列用户信息「管理员」
// @Description  根据用户关键信息进行模糊查询
// @Tags         user
// @Accept       json
// @Produce      json
// @Param		request body reqUser.UserQueryRequest true "需要查询的页数、以及用户关键信息"
// @Success      200  {object}  common.Response{data=resUser.ListUserVOResponse} "查询成功"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /v1/user/list/page/vo [POST]
func ListUserVOByPage(c *gin.Context) {
	queryReq := reqUser.UserQueryRequest{}
	c.ShouldBind(&queryReq)
	users, err := s.ListUserByPage(&queryReq)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	common.Success(c, *users)
}
