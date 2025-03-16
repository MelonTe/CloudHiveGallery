package service

import (
	"chg/internal/common"
	"chg/internal/consts"
	"chg/internal/ecode"
	"chg/internal/model/entity"
	reqUser "chg/internal/model/request/user"
	resUser "chg/internal/model/response/user"
	"chg/internal/repository"
	"chg/pkg/argon2"
	"chg/pkg/db"
	"chg/pkg/session"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserService struct {
	UserRepo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		UserRepo: repository.NewUserRepository(),
	}
}

// 执行用户注册服务，用户默认权限为user，昵称为无名
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
	if cnt, err = s.UserRepo.CountByAccount(userAccount); err != nil {
		return 0, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库查询错误")
	}
	if cnt > 0 {
		return 0, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "账号重复")
	}
	//3.加密
	encryptPassword := GetEncryptPassword(userPassword)
	//4.插入数据
	user := &entity.User{
		UserAccount:  userAccount,
		UserPassword: encryptPassword,
		UserName:     "无名",
		UserRole:     "user",
	}
	if err = s.UserRepo.CreateUser(user); err != nil {
		return 0, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库错误，注册失败")
	}
	return user.ID, nil
}

func GetEncryptPassword(userPassword string) string {
	//前四位充当盐值
	return argon2.GetEncryptString(userPassword, userPassword[:5])
}

// 用户登录服务，返回脱敏后的用户信息
func (s *UserService) UserLogin(c *gin.Context, userAccount, userPassword string) (*resUser.UserLoginVO, *ecode.ErrorWithCode) {
	//1.校验
	if userAccount == "" || userPassword == "" {
		return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "账号或密码为空")
	}
	if len(userAccount) < 4 || len(userPassword) < 8 {
		return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "账号或密码过短")
	}
	//2.加密、查询用户是否存在
	hashPsw := argon2.GetEncryptString(userPassword, userPassword[:5])
	user, err := s.UserRepo.FindByAccountAndPassword(userAccount, hashPsw)
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
	session.SetSession(c, consts.USER_LOGIN_STATE, userCopy)

	return resUser.GetUserLoginVO(userCopy), nil
}

// 获取当前登录用户，是数据库实体，用于内部可以复用
func (s *UserService) GetLoginUser(c *gin.Context) (*entity.User, *ecode.ErrorWithCode) {
	//从session中提取用户信息
	currentUser, ok := session.GetSession(c, consts.USER_LOGIN_STATE).(entity.User)
	if !ok {
		//对应的用户不存在
		return nil, ecode.GetErrWithDetail(ecode.NOT_LOGIN_ERROR, "用户未登录")
	}
	//数据库进行ID查询，避免数据不一致。追求性能可以不查询。
	curUser, err := s.UserRepo.FindById(currentUser.ID)
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

// 判断当前登录的用户是否是管理员
func (s *UserService) IsAdmin(c *gin.Context) bool {
	user, _ := s.GetLoginUser(c)
	if user != nil && user.UserRole == consts.ADMIN_ROLE {
		return true
	}
	return false
}

// 用户注销
func (s *UserService) UserLogout(c *gin.Context) (bool, *ecode.ErrorWithCode) {
	//从session中提取用户信息
	_, ok := session.GetSession(c, consts.USER_LOGIN_STATE).(entity.User)
	if !ok {
		//用户未登录
		return false, ecode.GetErrWithDetail(ecode.OPERATION_ERROR, "未登录")
	}
	//移除登录态
	session.DeleteSession(c, consts.USER_LOGIN_STATE)
	return true, nil
}

// 获取一个链式查询对象
func (s *UserService) GetQueryWrapper(db *gorm.DB, req *reqUser.UserQueryRequest) (*gorm.DB, *ecode.ErrorWithCode) {
	query := db.Session(&gorm.Session{})
	if req == nil {
		return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "参数为空")
	}
	if req.ID != 0 {
		query = query.Where("id = ?", req.ID)
	}
	if req.UserRole != "" {
		query = query.Where("user_role = ?", req.UserRole)
	}
	//模糊查询
	if req.UserAccount != "" {
		query = query.Where("user_account LIKE ?", "%"+req.UserAccount+"%")
	}
	if req.UserName != "" {
		query = query.Where("user_name LIKE ?", "%"+req.UserName+"%")
	}
	if req.UserProfile != "" {
		query = query.Where("user_profile LIKE ?", "%"+req.UserProfile+"%")
	}
	if req.SortField != "" {
		order := "ASC"
		if strings.ToLower(req.SortOrder) == "descend" {
			order = "DESC"
		}
		query = query.Order(req.SortField + " " + order)
	}
	return query, nil
}

// 根据ID软删除用户
func (s *UserService) RemoveById(id uint64) (bool, *ecode.ErrorWithCode) {
	if suc, err := s.UserRepo.RemoveById(id); err != nil {
		return false, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库错误")
	} else {
		if !suc {
			return false, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "用户不存在")
		}
		return true, nil
	}
}

// 更新用户信息，不存在则返回错误
func (s *UserService) UpdateUser(u *entity.User) *ecode.ErrorWithCode {
	if suc, err := s.UserRepo.UpdateUser(u); err != nil {
		return ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库错误")
	} else {
		if !suc {
			return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "用户不存在")
		}
		return nil
	}
}

// 获取用户列表
func (s *UserService) ListUserByPage(queryReq *reqUser.UserQueryRequest) (*resUser.ListUserVOResponse, *ecode.ErrorWithCode) {
	query, err := s.GetQueryWrapper(db.LoadDB(), queryReq)
	if err != nil {
		return nil, err
	}
	total, _ := s.UserRepo.GetQueryUsersNum(query)
	//拼接分页
	if queryReq.Current == 0 {
		queryReq.Current = 1
	}
	//重置query
	query, _ = s.GetQueryWrapper(db.LoadDB(), queryReq)
	query = query.Offset((queryReq.Current - 1) * queryReq.PageSize).Limit(queryReq.PageSize)
	users, errr := s.UserRepo.ListUserByPage(query)
	if errr != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库错误")
	}
	usersVO := resUser.GetUserVOList(users)
	p := (total + queryReq.PageSize - 1) / queryReq.PageSize
	return &resUser.ListUserVOResponse{
		Records: usersVO,
		PageResponse: common.PageResponse{
			Total:   total,
			Size:    queryReq.PageSize,
			Pages:   p,
			Current: queryReq.Current,
		},
	}, nil
}
