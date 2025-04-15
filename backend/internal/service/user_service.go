package service

import (
	"chg/internal/common"
	"chg/internal/consts"
	"chg/internal/ecode"
	"chg/internal/manager"
	"chg/internal/model/entity"
	reqUser "chg/internal/model/request/user"
	resUser "chg/internal/model/response/user"
	"chg/internal/repository"
	"chg/pkg/argon2"
	"chg/pkg/casbin"
	"chg/pkg/db"
	"chg/pkg/session"
	"fmt"
	"mime/multipart"
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

// 根据id获取用户，不存在返回错误
func (s *UserService) GetUserById(id uint64) (*entity.User, *ecode.ErrorWithCode) {
	user, err := s.UserRepo.FindById(nil, id)
	if err != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库查询错误")
	}
	if user == nil {
		return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "用户不存在")
	}
	return user, nil
}

// 根据id获取用户视图
func (s *UserService) GetUserVOById(id uint64) (*resUser.UserVO, *ecode.ErrorWithCode) {
	user, err := s.GetUserById(id)
	if err != nil {
		return nil, err
	}
	userVO := resUser.GetUserVO(*user)
	return &userVO, nil
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
	if cnt, err = s.UserRepo.CountByAccount(nil, userAccount); err != nil {
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
	if err = s.UserRepo.CreateUser(nil, user); err != nil {
		return 0, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库错误，注册失败")
	}
	//5.添加RBAC权限，新用户为在public域下的viewer
	casClient := casbin.LoadCasbinMethod()
	casbin.UpdateUserRoleInDomain(casClient, user.ID, consts.SPACEROLE_VIEWER, consts.DOM_PUBLIC)
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
	user, err := s.UserRepo.FindByAccountAndPassword(nil, userAccount, hashPsw)
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
// 未获取到用户信息时，返回nil和错误
func (s *UserService) GetLoginUser(c *gin.Context) (*entity.User, *ecode.ErrorWithCode) {
	//从session中提取用户信息
	currentUser, ok := session.GetSession(c, consts.USER_LOGIN_STATE).(entity.User)
	if !ok {
		//对应的用户不存在
		return nil, ecode.GetErrWithDetail(ecode.NOT_LOGIN_ERROR, "用户未登录")
	}
	//数据库进行ID查询，避免数据不一致。追求性能可以不查询。
	curUser, err := s.UserRepo.FindById(nil, currentUser.ID)
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
	if suc, err := s.UserRepo.RemoveById(nil, id); err != nil {
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
	if suc, err := s.UserRepo.UpdateUser(nil, u); err != nil {
		return ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库错误")
	} else {
		if !suc {
			return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "用户不存在")
		}
		return nil
	}
}

// 更新用户信息的特定字段，不存在则返回错误
func (s *UserService) UpdateUserByMap(req *reqUser.UserEditRequest) *ecode.ErrorWithCode {
	updateMap := map[string]interface{}{
		"user_name":    req.UserName,
		"user_profile": req.UserProfile,
	}
	if suc, err := s.UserRepo.UpdateUserByMap(nil, req.ID, updateMap); err != nil {
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
	total, _ := s.UserRepo.GetQueryUsersNum(nil, query)
	//拼接分页
	if queryReq.Current == 0 {
		queryReq.Current = 1
	}
	//重置query
	query, _ = s.GetQueryWrapper(db.LoadDB(), queryReq)
	query = query.Offset((queryReq.Current - 1) * queryReq.PageSize).Limit(queryReq.PageSize)
	users, errr := s.UserRepo.ListUserByPage(nil, query)
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

// 上传头像接口
func (s *UserService) UploadAvatar(file *multipart.FileHeader, userId uint64) (bool, *ecode.ErrorWithCode) {
	//1.校验
	if file == nil {
		return false, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "文件为空")
	}
	if file.Size > 5*1024*1024 {
		return false, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "文件过大")
	}
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		return false, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "文件格式错误")
	}
	//校验文件大小
	fileSize := file.Size
	ONE_MB := int64(1024 * 1024)
	if fileSize > 5*ONE_MB {
		return false, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "文件过大，不能超过2MB")
	}
	//2.获取文件路径
	//定义前缀
	uploadPrefix := fmt.Sprintf("avatar/%d", userId)
	result, err := manager.UploadPicture(file, uploadPrefix)
	if err != nil {
		return false, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "文件上传失败")
	}
	//3.更新数据库
	updateMap := map[string]interface{}{"user_avatar": fmt.Sprintf("%s", result.URL)}
	query := db.LoadDB()
	query.Model(&entity.User{}).Where("id = ?", userId).Updates(updateMap)
	if err := query.Error; err != nil {
		return false, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库更新失败")
	}
	return true, nil
}
