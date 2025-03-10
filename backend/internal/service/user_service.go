package service

import (
	"chg/internal/consts"
	"chg/internal/ecode"
	"chg/internal/middleware"
	"chg/internal/model/entity"
	"chg/internal/model/response"
	"chg/internal/repository"
	"chg/pkg/argon2"

	"github.com/gin-gonic/gin"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}
func (s *UserService) UserRegister(userAccount, userPassword, checkPassword string) (uint64, *ecode.ErrorWithCode) {
	//1.校验
	if userAccount == "" || userPassword == "" || checkPassword == "" {
		return 0, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "参数为空")
	}
	if len(userAccount) < 4 {
		return 0, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "用户账号过短")
	}
	if len(userPassword) < 8 || len(checkPassword) < 8 {
		return 0, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "用户密码过短")
	}
	if userPassword != checkPassword {
		return 0, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "两次输入的密码不一致")
	}

	//2.检查是否重复
	var cnt int64
	var err error
	if cnt, err = s.userRepo.CountByAccount(userAccount); err != nil {
		return 0, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库查询错误")
	}
	if cnt > 0 {
		return 0, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "账号重复")
	}
	//3.加密
	encryptPassword := getEncryptPassword(userPassword)
	//4.插入数据
	user := &entity.User{
		UserAccount:  userAccount,
		UserPassword: encryptPassword,
		UserName:     "无名",
		UserRole:     "user",
	}
	if err = s.userRepo.CreateUser(user); err != nil {
		return 0, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库错误，注册失败")
	}
	return user.ID, nil
}

func getEncryptPassword(userPassword string) string {
	//前四位充当盐值
	return argon2.GetEncryptString(userPassword, userPassword[:5])
}

// 用户登录服务，返回脱敏后的用户信息
func (s *UserService) UserLogin(c *gin.Context, userAccount, userPassword string) (*response.UserLoginVO, *ecode.ErrorWithCode) {
	//1.校验
	if userAccount == "" || userPassword == "" {
		return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "账号或密码为空")
	}
	if len(userAccount) < 4 || len(userPassword) < 8 {
		return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "账号或密码过短")
	}
	//2.加密、查询用户是否存在
	hashPsw := argon2.GetEncryptString(userPassword, userPassword[:5])
	user, err := s.userRepo.FindByAccountAndPassword(userAccount, hashPsw)
	if err != nil {
		//数据库异常
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库查询异常")
	}
	if user == nil {
		//用户不存在
		return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "用户不存在或密码错误")
	}
	//3.存储用户的登录态信息
	userCopy := *user //存储结构体，避免指针悬空
	middleware.SetSession(c, consts.USER_LOGIN_STATE, userCopy)

	return response.GetUserLoginVO(userCopy), nil
}

// 获取当前登录用户，是数据库实体，用于内部可以复用
func (s *UserService) GetLoginUser(c *gin.Context) (*entity.User, *ecode.ErrorWithCode) {
	//从session中提取用户信息
	currentUser, ok := middleware.GetSession(c, consts.USER_LOGIN_STATE).(entity.User)
	if !ok {
		//对应的用户不存在
		return nil, ecode.GetErrWithDetail(ecode.NOT_LOGIN_ERROR, "")
	}
	//数据库进行ID查询，避免数据不一致。追求性能可以不查询。
	curUser, err := s.userRepo.FindByAId(currentUser.ID)
	if err != nil {
		//数据库异常
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库查询失败")
	}
	if curUser == nil {
		//用户不存在
		return nil, ecode.GetErrWithDetail(ecode.NOT_LOGIN_ERROR, "用户不存在")
	}
	return curUser, nil
}

// 用户注销
func (s *UserService) UserLogout(c *gin.Context) (bool, *ecode.ErrorWithCode) {
	//从session中提取用户信息
	_, ok := middleware.GetSession(c, consts.USER_LOGIN_STATE).(entity.User)
	if !ok {
		//用户未登录
		return false, ecode.GetErrWithDetail(ecode.OPERATION_ERROR, "未登录")
	}
	//移除登录态
	middleware.DeleteSession(c, consts.USER_LOGIN_STATE)
	return true, nil
}
