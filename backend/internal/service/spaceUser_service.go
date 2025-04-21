package service

import (
	"chg/internal/consts"
	"chg/internal/ecode"
	"chg/internal/model/entity"
	reqSpaceUser "chg/internal/model/request/spaceuser"
	resSpace "chg/internal/model/response/space"
	resSpaceUser "chg/internal/model/response/spaceuser"
	resUser "chg/internal/model/response/user"
	"chg/internal/repository"
	"chg/pkg/casbin"
	"chg/pkg/db"
	"fmt"

	"gorm.io/gorm"
)

type SpaceUserService struct {
	SpaceUserRepo *repository.SpaceUserRepository
}

func NewSpaceUserService() *SpaceUserService {
	return &SpaceUserService{
		SpaceUserRepo: repository.NewSpaceUserRepository(),
	}
}

// 添加空间成员方法
func (s *SpaceUserService) AddSpaceUser(req *reqSpaceUser.SpaceUserAddRequest) (uint64, *ecode.ErrorWithCode) {
	//参数校验
	spaceUser := &entity.SpaceUser{
		SpaceID:   req.SpaceID,
		UserID:    req.UserID,
		SpaceRole: req.SpaceRole,
	}
	if req.SpaceRole == "" {
		//默认为浏览者
		spaceUser.SpaceRole = consts.SPACEROLE_VIEWER
	}
	if err := ValidSpaceUser(spaceUser, true); err != nil {
		return 0, err
	}
	//数据库添加新成员
	query := db.LoadDB()
	originErr := query.Save(spaceUser).Error
	if originErr != nil {
		return 0, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库操作失败")
	}
	//更新RBAC权限
	dom := fmt.Sprintf("space_%d", req.SpaceID)
	casbin.UpdateUserRoleInDomain(casbin.Casbin, req.UserID, consts.SPACEROLE_VIEWER, dom)
	return spaceUser.ID, nil
}

// 校验空间成员对象，区分是编辑校验还是增加成员校验
func ValidSpaceUser(spaceUser *entity.SpaceUser, add bool) *ecode.ErrorWithCode {
	//若创建，需校验是否填写了空间ID和用户ID
	if add {
		_, err := NewSpaceService().GetSpaceById(spaceUser.SpaceID)
		if err != nil {
			return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "空间不存在")
		}
		_, err = NewUserService().GetUserById(spaceUser.UserID)
		if err != nil {
			return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "用户不存在")
		}
	}
	//校验空间角色
	if exist := consts.IsSpaceUserRoleExist(spaceUser.SpaceRole); !exist {
		return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "空间角色不存在")
	}
	return nil
}

// 封装链式查询对象
func (s *SpaceUserService) GetQueryWrapper(db *gorm.DB, req *reqSpaceUser.SpaceUserQueryRequest) (*gorm.DB, *ecode.ErrorWithCode) {
	query := db.Session(&gorm.Session{})
	if req.ID > 0 {
		query = query.Where("id = ?", req.ID)
	}
	if req.SpaceID > 0 {
		query = query.Where("space_id = ?", req.SpaceID)
	}
	if req.UserID > 0 {
		query = query.Where("user_id = ?", req.UserID)
	}
	if req.SpaceRole != "" {
		query = query.Where("space_role = ?", req.SpaceRole)
	}
	return query, nil
}

// 获取空间成员视图的上层封装
func (s *SpaceUserService) GetSpaceUserVO(spaceUser *entity.SpaceUser) *resSpaceUser.SpaceUserVO {
	//主要获取UserVO和SpaceVO
	vo := &resSpaceUser.SpaceUserVO{
		ID:        spaceUser.ID,
		SpaceID:   spaceUser.SpaceID,
		UserID:    spaceUser.UserID,
		SpaceRole: spaceUser.SpaceRole,
	}
	if spaceUser.UserID > 0 {
		user, _ := NewUserService().GetUserById(spaceUser.UserID)
		if user != nil {
			vo.User = resUser.GetUserVO(*user)
		}
	}
	if spaceUser.SpaceID > 0 {
		space, _ := NewSpaceService().GetSpaceById(spaceUser.SpaceID)
		if space != nil {
			vo.SpaceVO = resSpace.EntityToVO(*space, vo.User)
		}
	}
	return vo
}

// 根据空间成员实体列表获取空间成员视图列表
func (s *SpaceUserService) GetSpaceUserVOList(spaceUsers []entity.SpaceUser) []resSpaceUser.SpaceUserVO {
	recordUserVO := make(map[uint64]resUser.UserVO)
	recordSpaceVO := make(map[uint64]resSpace.SpaceVO)
	for _, spaceUser := range spaceUsers {
		if _, ok := recordUserVO[spaceUser.UserID]; !ok {
			//该用户没有被查询过，进行查询
			user, _ := NewUserService().GetUserById(spaceUser.UserID)
			//保证用户的存在
			userVO := resUser.GetUserVO(*user)
			recordUserVO[spaceUser.UserID] = userVO
		}
		if _, ok := recordSpaceVO[spaceUser.SpaceID]; !ok {
			//该空间没有被查询过，进行查询
			space, _ := NewSpaceService().GetSpaceById(spaceUser.SpaceID)
			//保证空间的存在
			spaceVO := resSpace.EntityToVO(*space, recordUserVO[spaceUser.UserID])
			recordSpaceVO[spaceUser.SpaceID] = spaceVO
		}
	}
	//封装返回
	voList := make([]resSpaceUser.SpaceUserVO, 0, len(spaceUsers))
	for _, spaceUser := range spaceUsers {
		vo := resSpaceUser.SpaceUserVO{
			ID:        spaceUser.ID,
			SpaceID:   spaceUser.SpaceID,
			UserID:    spaceUser.UserID,
			SpaceRole: spaceUser.SpaceRole,
		}
		vo.User = recordUserVO[spaceUser.UserID]
		vo.SpaceVO = recordSpaceVO[spaceUser.SpaceID]
		voList = append(voList, vo)
	}
	return voList
}

// 根据ID移除空间成员
func (s *SpaceUserService) RemoveSpaceUserById(id uint64) *ecode.ErrorWithCode {
	//校验ID是否存在
	if id <= 0 {
		return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "ID不能为空")
	}
	//删除空间成员
	if err := db.LoadDB().Where("id = ?", id).Delete(&entity.SpaceUser{}).Error; err != nil {
		return ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库操作失败")
	}
	return nil
}

// 获取空间成员视图列表
func (s *SpaceUserService) ListSpaceUserVO(req *reqSpaceUser.SpaceUserQueryRequest) (*ecode.ErrorWithCode, []resSpaceUser.SpaceUserVO) {
	if exists := consts.IsSpaceUserRoleExist(req.SpaceRole); !exists && req.SpaceRole != "" {
		return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "空间角色不存在"), nil
	}
	query, err := s.GetQueryWrapper(db.LoadDB(), req)
	if err != nil {
		return err, nil
	}
	var spaceUsers []entity.SpaceUser
	query.Model(&entity.SpaceUser{}).Find(&spaceUsers)
	//获取空间成员视图列表
	voList := s.GetSpaceUserVOList(spaceUsers)
	return nil, voList
}

func (s *SpaceUserService) EditSpaceUser(req *reqSpaceUser.SpaceUserEditRequest) (bool, *ecode.ErrorWithCode) {
	//参数校验
	if req.ID <= 0 {
		return false, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "ID不能为空")
	}
	if req.SpaceRole != "" && !consts.IsSpaceUserRoleExist(req.SpaceRole) {
		return false, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "空间角色不存在")
	}
	//记录校验
	oldSpaceUser := &entity.SpaceUser{}
	query := db.LoadDB()
	originErr := query.Model(&entity.SpaceUser{}).Where("id = ?", req.ID).First(oldSpaceUser).Error
	if originErr != nil {
		return false, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "没有找到该空间成员")
	}
	if oldSpaceUser.SpaceRole == req.SpaceRole {
		return true, nil
	}
	if err := ValidSpaceUser(oldSpaceUser, false); err != nil {
		return false, err
	}
	//更新空间成员
	query = db.LoadDB()
	query.Model(&entity.SpaceUser{}).Where("id = ?", req.ID).Updates(map[string]interface{}{
		"space_role": req.SpaceRole,
	})
	//更新这个空间成员的权限
	casClient := casbin.LoadCasbinMethod()
	domain := fmt.Sprintf("space_%d", oldSpaceUser.SpaceID)
	casbin.UpdateUserRoleInDomain(casClient, oldSpaceUser.UserID, req.SpaceRole, domain)
	return true, nil
}

// 根据记录的Id获取空间成员记录
func (s *SpaceUserService) GetSpaceUserById(id uint64) (*entity.SpaceUser, *ecode.ErrorWithCode) {
	if id <= 0 {
		return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "ID不能为空")
	}
	//查询空间成员记录
	query := db.LoadDB()
	spaceUser := &entity.SpaceUser{}
	originErr := query.Model(&entity.SpaceUser{}).Where("id = ?", id).First(spaceUser).Error
	if originErr != nil {
		return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "没有找到该空间成员")
	}
	return spaceUser, nil
}

// 根据空间Id和用户Id获取空间成员记录
func (s *SpaceUserService) GetSpaceUserBySpaceIdAndUserId(spaceId uint64, userId uint64) (*entity.SpaceUser, *ecode.ErrorWithCode) {
	if spaceId <= 0 {
		return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "ID不能为空")
	}
	//查询空间成员记录
	query := db.LoadDB()
	spaceUser := &entity.SpaceUser{}
	originErr := query.Model(&entity.SpaceUser{}).Where("space_id = ? and user_id = ?", spaceId, userId).First(spaceUser).Error
	if originErr != nil {
		return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "没有找到该空间成员")
	}
	return spaceUser, nil
}

// 获取成员的权限列表，权限有picture:edit、picture:view、picture:upload、picture:delete、spaceUser:manage
// 若是公共图库，space请传入nil
func GetPermissionList(space *entity.Space, loginUser *entity.User) []string {
	permissionList := []string{}
	if loginUser == nil {
		return permissionList
	}
	adminPermission := []string{"picture:view", "picture:edit", "picture:delete", "picture:upload", "spaceUser:manage"}
	editorPermission := []string{"picture:view", "picture:edit", "picture:delete", "picture:upload"}
	viewerPermission := []string{"picture:view"}
	//公共图库
	if space == nil {
		if loginUser.UserRole == consts.ADMIN_ROLE {
			return adminPermission
		}
		//默认都拥有只读权限
		return viewerPermission
	}
	switch space.SpaceType {
	case consts.SPACE_PRIVATE:
		//私人空间，本人或管理员有所有权限
		if space.UserID == loginUser.ID || loginUser.UserRole == consts.ADMIN_ROLE {
			return adminPermission
		} else {
			return permissionList
		}
	case consts.SPACE_TEAM:
		//团队空间，查询loginUser的身份
		spaceUserInfo, err := NewSpaceUserService().GetSpaceUserBySpaceIdAndUserId(space.ID, loginUser.ID)
		if err != nil {
			return permissionList
		}
		if spaceUserInfo != nil {
			if spaceUserInfo.SpaceRole == consts.SPACEROLE_ADMIN {
				return adminPermission
			} else if spaceUserInfo.SpaceRole == consts.SPACEROLE_EDITOR {
				return editorPermission
			} else {
				return viewerPermission
			}
		}
	}
	return permissionList
}
