package service

import (
	"chg/internal/common"
	"chg/internal/consts"
	"chg/internal/ecode"
	"chg/internal/model/entity"
	reqSpace "chg/internal/model/request/space"
	resSpace "chg/internal/model/response/space"
	resUser "chg/internal/model/response/user"
	"chg/internal/repository"
	"chg/pkg/casbin"
	"chg/pkg/db"
	"chg/pkg/redlock"
	"fmt"
	"log"
	"strconv"
	"time"
	"unicode/utf8"

	"gorm.io/gorm"
)

type SpaceService struct {
	SpaceRepo *repository.SpaceRepository
}

func NewSpaceService() *SpaceService {
	return &SpaceService{
		SpaceRepo: repository.NewSpaceRepository(),
	}
}

// 校验空间更新数据是否正常，包括昵称、级别
func (s *SpaceService) ValidSpace(space *reqSpace.SpaceUpdateRequest, add bool) *ecode.ErrorWithCode {
	spaceName := space.SpaceName
	spaceLevel := consts.GetSpaceLevelByValue(space.SpaceLevel)
	// 要创建
	if add {
		if spaceName == "" {
			return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "空间名称不能为空")
		}
		if spaceLevel == nil {
			return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "空间级别错误")
		}
	}
	// 修改数据时，如果要改空间级别
	if spaceLevel == nil {
		return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "空间级别不存在")
	}
	if space.SpaceName != "" && utf8.RuneCountInString(spaceName) > 30 {
		return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "空间名称过长")
	}
	return nil
}

// 自动填充空间的等级额度到map中
func (s *SpaceService) FillSpaceByLevelInMap(spaceLevel int, updateMap map[string]interface{}) {
	spaceLevelEnum := consts.GetSpaceLevelByValue(spaceLevel)
	if spaceLevelEnum != nil {
		updateMap["max_count"] = spaceLevelEnum.MaxCount
		updateMap["max_size"] = spaceLevelEnum.MaxSize
	}
}

// 自动填充空间的等级额度
func (s *SpaceService) FillSpaceByLevel(space *entity.Space) {
	spaceLevelEnum := consts.GetSpaceLevelByValue(space.SpaceLevel)
	if spaceLevelEnum != nil {
		space.MaxCount = spaceLevelEnum.MaxCount
		space.MaxSize = spaceLevelEnum.MaxSize
	}
}

// 更新空间
func (s *SpaceService) UpdateSpace(space *reqSpace.SpaceUpdateRequest, loginUser *entity.User) *ecode.ErrorWithCode {
	//查找旧的空间，校验是否存在
	oldSpace, err := s.SpaceRepo.GetSpaceById(nil, space.ID)
	if err != nil {
		return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "数据库查询失败")
	}
	if oldSpace == nil {
		return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "空间不存在")
	}
	//权限校验
	if !(oldSpace.UserID == loginUser.ID) {
		return ecode.GetErrWithDetail(ecode.NO_AUTH_ERROR, "无权限")
	}
	//校验数据
	if err := s.ValidSpace(space, false); err != nil {
		return err
	}
	updateMap := make(map[string]interface{}, 8)
	//填充数据
	updateMap["space_name"] = space.SpaceName
	updateMap["space_level"] = space.SpaceLevel
	s.FillSpaceByLevelInMap(space.SpaceLevel, updateMap)
	//更新数据库数据
	if err := s.SpaceRepo.UpdateSpaceById(nil, space.ID, updateMap); err != nil {
		return ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "更新失败")
	}
	return nil
}

// 编辑空间
func (s *SpaceService) EditSpace(space *reqSpace.SpaceEditRequest, loginUser *entity.User) *ecode.ErrorWithCode {
	//查找旧的空间，校验是否存在
	oldSpace, err := s.SpaceRepo.GetSpaceById(nil, space.ID)
	if err != nil {
		return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "数据库查询失败")
	}
	if oldSpace == nil {
		return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "空间不存在")
	}
	//权限校验
	if !(oldSpace.UserID == loginUser.ID) {
		return ecode.GetErrWithDetail(ecode.NO_AUTH_ERROR, "无权限")
	}
	updateMap := make(map[string]interface{}, 8)
	//填充数据
	updateMap["space_name"] = space.SpaceName
	updateMap["edit_time"] = time.Now()
	//更新数据库数据
	if err := s.SpaceRepo.UpdateSpaceById(nil, space.ID, updateMap); err != nil {
		return ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "更新失败")
	}
	return nil
}

// 分页查询空间视图
func (s *SpaceService) ListSpaceVOByPage(req *reqSpace.SpaceQueryRequest) (*resSpace.ListSpaceVOResponse, *ecode.ErrorWithCode) {
	// 调用 SpaceList
	list, err := s.ListSpaceByPage(req)
	if err != nil {
		return nil, err
	}
	//获取用户对象
	userService := NewUserService()
	userVO, err := userService.GetUserVOById(req.UserID)
	// 获取 VO 对象
	listVO := &resSpace.ListSpaceVOResponse{
		PageResponse: list.PageResponse,
		Records:      s.GetSpaceVOList(list.Records, userVO),
	}
	return listVO, nil
}

// Spaces 数组转换为 SpaceVO 数组
func (s *SpaceService) GetSpaceVOList(records []entity.Space, userVO *resUser.UserVO) []resSpace.SpaceVO {
	var vos []resSpace.SpaceVO
	for _, record := range records {
		vos = append(vos, resSpace.EntityToVO(record, *userVO))
	}
	return vos
}

// 分页查询空间
func (s *SpaceService) ListSpaceByPage(req *reqSpace.SpaceQueryRequest) (*resSpace.ListSpaceResponse, *ecode.ErrorWithCode) {
	// 参数校验
	if req.Current <= 0 || req.PageSize <= 0 {
		return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "页数或者页面大小不能小于0")
	}
	// 获取查询对象
	query, err := s.GetQueryWrapper(db.LoadDB(), req)
	if err != nil {
		return nil, err
	}
	// 查询总数
	var total int64
	query.Model(&entity.Space{}).Count(&total)
	to := int(total)
	// 分页查询
	var spaces []entity.Space
	// 重置 query
	query, _ = s.GetQueryWrapper(db.LoadDB(), req)
	query = query.Offset((req.Current - 1) * req.PageSize).Limit(req.PageSize)
	query.Find(&spaces)
	p := (to + req.PageSize - 1) / req.PageSize
	//封装返回
	list := &resSpace.ListSpaceResponse{
		Records: spaces,
		PageResponse: common.PageResponse{
			Total:   to,
			Size:    req.PageSize,
			Pages:   p,
			Current: req.Current,
		},
	}
	return list, nil
}

// 获取一个链式查询对象
func (s *SpaceService) GetQueryWrapper(db *gorm.DB, req *reqSpace.SpaceQueryRequest) (*gorm.DB, *ecode.ErrorWithCode) {
	query := db.Session(&gorm.Session{})
	if req.ID != 0 {
		query = query.Where("id = ?", req.ID)
	}
	if req.UserID != 0 {
		query = query.Where("user_id = ?", req.UserID)
	}
	if req.SpaceName != "" {
		query = query.Where("space_name LIKE ?", "%"+req.SpaceName+"%")
	}
	if req.SpaceLevel != nil {
		query = query.Where("space_level = ?", *req.SpaceLevel)
	}
	if req.SpaceType != nil {
		query = query.Where("space_type = ?", *req.SpaceType)
	}
	if req.SortField != "" {
		sortOrder := "ASC"
		if req.SortOrder == "descend" {
			sortOrder = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", req.SortField, sortOrder))
	}
	return query, nil
}

// 添加创建一个空间
func (s *SpaceService) AddSpace(addRequest *reqSpace.SpaceAddRequest, loginUser *entity.User) (uint64, *ecode.ErrorWithCode) {
	//1.校验数据
	if addRequest.SpaceName == "" {
		addRequest.SpaceName = loginUser.UserName + "的空间"
	}
	spaceLevel := consts.GetSpaceLevelByValue(addRequest.SpaceLevel)
	if spaceLevel == nil {
		spaceLevel = consts.GetSpaceLevelByValue(consts.COMMON.Value) //默认为0级别空间
	}
	if spaceTypeValid := consts.IsSpaceTypeValid(addRequest.SpaceType); !spaceTypeValid {
		return 0, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "空间类型错误")
	}
	//实体类创建
	space := &entity.Space{
		SpaceName:  addRequest.SpaceName,
		SpaceLevel: addRequest.SpaceLevel,
		SpaceType:  addRequest.SpaceType,
	}
	//参数填充
	s.FillSpaceByLevel(space)
	space.UserID = loginUser.ID
	//2.校验权限，只允许管理员创建指定级别的空间
	if consts.COMMON.Value != addRequest.SpaceLevel && loginUser.UserRole != consts.ADMIN_ROLE {
		return 0, ecode.GetErrWithDetail(ecode.NO_AUTH_ERROR, "无权限创建指定级别的空间")
	}
	//3.锁+事务，保证同一用户只能创建一个空间，以及只能创建一个团队空间
	rs := redlock.GetRedSync()
	lock := rs.NewMutex(strconv.FormatUint(loginUser.ID, 10))
	//加锁，超时时间默认8s
	lock.Lock()
	defer lock.Unlock()
	//开启事务
	tx := s.SpaceRepo.BeginTransaction()
	//进行数据库校验，查看是否存在数据
	exist := s.SpaceRepo.IsExistByUserId(tx, loginUser.ID, addRequest.SpaceType)
	if exist {
		return 0, ecode.GetErrWithDetail(ecode.NO_AUTH_ERROR, "每个用户每类空间只允许创建一个")
	}
	//写入数据库
	err := s.SpaceRepo.SaveSpace(tx, space)
	if err != nil {
		return 0, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库错误")
	}
	//若是团队空间，需要在space_user表插入当前创建人的信息
	if space.SpaceType == consts.SPACE_TEAM {
		tx.Model(&entity.SpaceUser{}).Save(&entity.SpaceUser{
			SpaceID:   space.ID,
			UserID:    loginUser.ID,
			SpaceRole: consts.SPACEROLE_ADMIN,
		})
	}
	//记录创建者的权限
	casClient := casbin.LoadCasbinMethod()
	domain := fmt.Sprintf("space_%d", space.ID)
	originErr := casbin.UpdateUserRoleInDomain(casClient, loginUser.ID, consts.ADMIN_ROLE, domain)
	if originErr != nil {
		log.Println("更新权限失败", originErr)
		return 0, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "权限更新失败")
	}
	//提交事务
	err = tx.Commit().Error
	if err != nil {
		log.Println("事务提交失败", err)
		return 0, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库错误")
	}
	return space.ID, nil
}

// 根据ID获取空间，若不存在返回错误
func (s *SpaceService) GetSpaceById(id uint64) (*entity.Space, *ecode.ErrorWithCode) {
	space, err := s.SpaceRepo.GetSpaceById(nil, id)
	if err != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库错误")
	}
	if space == nil {
		return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "空间不存在")
	}
	return space, nil
}
