package service

import (
	"chg/internal/common"
	"chg/internal/ecode"
	"chg/internal/manager"
	"chg/internal/model/entity"
	reqPicture "chg/internal/model/request/picture"
	resPicture "chg/internal/model/response/picture"
	resUser "chg/internal/model/response/user"
	"chg/internal/repository"
	"chg/pkg/db"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"time"
	"unicode/utf8"

	"gorm.io/gorm"
)

type PictureService struct {
	PictureRepo *repository.PictureRepository
}

func NewPictureService() *PictureService {
	return &PictureService{
		PictureRepo: repository.NewPictureRepository(),
	}
}

// 修改或插入图片数据到服务器中
func (s *PictureService) UploadPicture(multipartFile *multipart.FileHeader, PictureUploadRequest *reqPicture.PictureUploadRequest, loginUser *entity.User) (*resPicture.PictureVO, *ecode.ErrorWithCode) {
	//判断图片是需要新增还是需要更新
	picId := uint64(0)
	if PictureUploadRequest.ID != 0 {
		picId = PictureUploadRequest.ID
	}
	//若更新图片，则需要校验图片是否存在
	if picId != 0 {
		_, err := s.PictureRepo.FindById(picId)
		if err != nil {
			return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "图片不存在")
		}
	}
	//上传图片，得到信息
	uploadPathPrefix := fmt.Sprintf("public/%d", loginUser.ID)
	info, err := manager.UploadPicture(multipartFile, uploadPathPrefix)
	if err != nil {
		return nil, err
	}
	//构造插入数据库的实体
	pic := &entity.Picture{
		URL:       info.URL,
		Name:      info.PicName,
		PicSize:   info.PicSize,
		PicWidth:  info.PicWidth,
		PicHeight: info.PicHeight,
		PicScale:  info.PicScale,
		PicFormat: info.PicFormat,
		UserID:    loginUser.ID,
		EditTime:  time.Now(),
	}
	//若是更新，则需要更新ID
	if picId != 0 {
		pic.ID = picId
	}
	//进行插入或者更新操作，即save
	errr := s.PictureRepo.SavePicture(pic)
	if errr != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库错误")
	}
	userVO := resUser.GetUserVO(*loginUser)
	picVO := resPicture.EntityToVO(*pic, userVO)
	return &picVO, nil
}

// 获取一个链式查询对象
func (s *PictureService) GetQueryWrapper(db *gorm.DB, req *reqPicture.PictureQueryRequest) (*gorm.DB, *ecode.ErrorWithCode) {
	query := db.Session(&gorm.Session{})
	if req.SearchText != "" {
		query = query.Where("name LIKE ? OR introduction LIKE ?", "%"+req.SearchText+"%", "%"+req.SearchText+"%")
	}
	if req.ID != 0 {
		query = query.Where("id = ?", req.ID)
	}
	if req.UserID != 0 {
		query = query.Where("user_id = ?", req.UserID)
	}
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Introduction != "" {
		query = query.Where("introduction LIKE ?", "%"+req.Introduction+"%")
	}
	if req.PicFormat != "" {
		query = query.Where("pic_format LIKE ?", "%"+req.PicFormat+"%")
	}
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}
	if req.PicWidth != 0 {
		query = query.Where("pic_width = ?", req.PicWidth)
	}
	if req.PicHeight != 0 {
		query = query.Where("pic_height = ?", req.PicHeight)
	}
	if req.PicSize != 0 {
		query = query.Where("pic_size = ?", req.PicSize)
	}
	if req.PicScale != 0 {
		query = query.Where("pic_scale = ?", req.PicScale)
	}
	//tags在数据库中的存储格式：["golang","java","c++"]
	if len(req.Tags) > 0 {
		//and (tags LIKE %"commic" and tags LIKE %"manga"% ...)
		for _, tag := range req.Tags {
			query = query.Where("tags LIKE ?", "%\""+tag+"\"%")
		}
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

// 获取一个PicVO对象
func (s *PictureService) GetPictureVO(Picture *entity.Picture) *resPicture.PictureVO {
	user, err := repository.NewUserRepository().FindById(Picture.UserID)
	if err != nil {
		return nil
	}
	//防止user为空
	var picVO resPicture.PictureVO
	if user != nil {
		userVO := resUser.GetUserVO(*user)
		picVO = resPicture.EntityToVO(*Picture, userVO)
	} else {
		picVO = resPicture.EntityToVO(*Picture, resUser.UserVO{})
	}
	return &picVO
}

// 获取PictureVO列表
func (s *PictureService) GetPictureVOList(Pictures []entity.Picture) []resPicture.PictureVO {
	var picVOList []resPicture.PictureVO
	//保存所有需要的user对象
	userMap := make(map[uint64]resUser.UserVO)
	for _, Picture := range Pictures {
		//user还没被查询，那么就查询
		if _, ok := userMap[Picture.UserID]; !ok {
			user, err := repository.NewUserRepository().FindById(Picture.UserID)
			if err != nil {
				log.Println("GetPictureVOList: 查询用户失败，错误为", err)
				//跳过
				continue
			}
			userVO := resUser.GetUserVO(*user)
			userMap[Picture.UserID] = userVO
		}
	}
	for _, v := range Pictures {
		picVOList = append(picVOList, resPicture.EntityToVO(v, userMap[v.UserID]))
	}
	return picVOList
}

// 图片参数校验，在更新和修改图片前进行判断
func (s *PictureService) ValidPicture(Picture *entity.Picture) *ecode.ErrorWithCode {
	if Picture.ID == 0 {
		return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "图片ID不能为空")
	}
	if len(Picture.URL) > 1024 {
		return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "图片URL过长")
	}
	if len(Picture.Introduction) > 800 {
		return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "图片简介过长")
	}
	if Picture.Name == "" || utf8.RuneCountInString(Picture.Name) > 20 {
		fmt.Println(Picture.Name)
		fmt.Println(len(Picture.Name))
		return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "图片名不能为空或不能超过20个字符")
	}
	return nil
}

// 根据ID获取图片，若图片不存在则返回错误
func (s *PictureService) GetPictureById(id uint64) (*entity.Picture, *ecode.ErrorWithCode) {
	Picture, err := s.PictureRepo.FindById(id)
	if err != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库错误")
	}
	if Picture == nil {
		return nil, ecode.GetErrWithDetail(ecode.NOT_FOUND_ERROR, "图片不存在")
	}
	return Picture, nil
}

// 根据ID删除图片
func (s *PictureService) DeletePictureById(id uint64) *ecode.ErrorWithCode {
	err := s.PictureRepo.DeleteById(id)
	if err != nil {
		return ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库错误")
	}
	return nil
}

// 更新图片
func (s *PictureService) UpdatePicture(updateReq *reqPicture.PictureUpdateRequest) *ecode.ErrorWithCode {
	//查询图片是否存在
	oldPic, err := s.GetPictureById(updateReq.ID)
	if err != nil {
		return err
	}
	//校验图片参数
	oldPic.Name = updateReq.Name
	oldPic.Introduction = updateReq.Introduction
	oldPic.Category = updateReq.Category
	if err := s.ValidPicture(oldPic); err != nil {
		return err
	}
	//保留更新字段
	updateMap := make(map[string]interface{}, 8)
	updateMap["name"] = oldPic.Name
	updateMap["introduction"] = oldPic.Introduction
	updateMap["category"] = oldPic.Category
	tags, _ := json.Marshal(updateReq.Tags)
	updateMap["tags"] = string(tags)
	updateMap["edit_time"] = time.Now()
	//更新
	if err := s.PictureRepo.UpdateById(updateReq.ID, updateMap); err != nil {
		return ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库错误")
	}
	return nil
}

// 分页查询图片
func (s *PictureService) ListPictureByPage(req *reqPicture.PictureQueryRequest) (*resPicture.ListPictureResponse, *ecode.ErrorWithCode) {
	//参数校验
	if req.Current <= 0 || req.PageSize <= 0 {
		return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "页数或者页面大小不能小于0")
	}
	//获取查询对象
	query, err := s.GetQueryWrapper(db.LoadDB(), req)
	if err != nil {
		return nil, err
	}
	//查询总数
	var total int64
	query.Model(&entity.Picture{}).Count(&total)
	to := int(total)
	//分页查询
	var Pictures []entity.Picture
	//重置query
	query, _ = s.GetQueryWrapper(db.LoadDB(), req)
	query = query.Offset((req.Current - 1) * req.PageSize).Limit(req.PageSize)
	query.Find(&Pictures)
	p := (to + req.PageSize - 1) / req.PageSize
	//获取VO对象
	list := &resPicture.ListPictureResponse{
		Records: Pictures,
		PageResponse: common.PageResponse{
			Total:   to,
			Size:    req.PageSize,
			Pages:   p,
			Current: req.Current,
		},
	}
	return list, nil
}

// 分页查询图片视图
func (s *PictureService) ListPictureVOByPage(req *reqPicture.PictureQueryRequest) (*resPicture.ListPictureVOResponse, *ecode.ErrorWithCode) {
	//调用PictureList
	list, err := s.ListPictureByPage(req)
	if err != nil {
		return nil, err
	}
	//获取VO对象
	listVO := &resPicture.ListPictureVOResponse{
		PageResponse: list.PageResponse,
		Records:      s.GetPictureVOList(list.Records),
	}
	return listVO, nil
}
