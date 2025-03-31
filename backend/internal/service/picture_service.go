package service

import (
	"chg/internal/common"
	"chg/internal/consts"
	"chg/internal/ecode"
	"chg/internal/manager"
	"chg/internal/model/dto/file"
	"chg/internal/model/entity"
	reqPicture "chg/internal/model/request/picture"
	resPicture "chg/internal/model/response/picture"
	resUser "chg/internal/model/response/user"
	"chg/internal/repository"
	"chg/pkg/cache"
	"chg/pkg/db"
	"chg/pkg/rds"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
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
// 修改为接收接口类型，可以是URL地址或者文件（multipartFile）
func (s *PictureService) UploadPicture(picFile interface{}, PictureUploadRequest *reqPicture.PictureUploadRequest, loginUser *entity.User) (*resPicture.PictureVO, *ecode.ErrorWithCode) {
	//判断图片是需要新增还是需要更新
	picId := uint64(0)
	if PictureUploadRequest.ID != 0 {
		picId = PictureUploadRequest.ID
	}
	var space *entity.Space
	//校验空间ID是否存在
	//若存在，则需要校验空间是否存在以及是否有权限上传
	if PictureUploadRequest.SpaceID != 0 {
		var err error
		space, err = repository.NewSpaceRepository().GetSpaceById(nil, PictureUploadRequest.SpaceID)
		if err != nil {
			return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库异常")
		}
		if space == nil {
			return nil, ecode.GetErrWithDetail(ecode.NOT_FOUND_ERROR, "空间不存在")
		}
		//仅允许空间管理员上传图片
		if space.UserID != loginUser.ID {
			return nil, ecode.GetErrWithDetail(ecode.NO_AUTH_ERROR, "没有空间权限")
		}
		//校验额度
		if space.TotalCount >= space.MaxCount {
			return nil, ecode.GetErrWithDetail(ecode.NO_AUTH_ERROR, "空间图片数量已满")
		}
		if space.TotalSize >= space.MaxSize {
			return nil, ecode.GetErrWithDetail(ecode.NO_AUTH_ERROR, "空间图片大小已满")
		}
	}
	//若更新图片，则需要校验图片是否存在，以及空间id是否跟原本的一致
	if picId != 0 {
		oldpic, err := s.PictureRepo.FindById(nil, picId)
		if err != nil {
			return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库异常")
		}
		if oldpic == nil {
			return nil, ecode.GetErrWithDetail(ecode.NOT_FOUND_ERROR, "图片不存在")
		}
		//权限校验，仅本人或管理员可以编辑
		if loginUser.UserRole != consts.ADMIN_ROLE && loginUser.ID != oldpic.UserID {
			return nil, ecode.GetErrWithDetail(ecode.NO_AUTH_ERROR, "权限不足")
		}
		//校验空间是否一致
		if space != nil && oldpic.SpaceID != PictureUploadRequest.SpaceID {
			return nil, ecode.GetErrWithDetail(ecode.NO_AUTH_ERROR, "空间不一致")
		}
		//没传spaceID，则复用原有图片的spaceID（兼容了公共图库）
		if space == nil {
			PictureUploadRequest.SpaceID = oldpic.SpaceID
		}
	}
	//上传图片，得到信息
	//去要区分上传到公共图库还是私人图库
	var uploadPathPrefix string
	if PictureUploadRequest.SpaceID == 0 {
		uploadPathPrefix = fmt.Sprintf("public/%d", loginUser.ID)
	} else {
		//存在space，则上传到私人图库
		uploadPathPrefix = fmt.Sprintf("space/%d", PictureUploadRequest.SpaceID)
	}

	var info *file.UploadPictureResult
	var err *ecode.ErrorWithCode
	//根据参数的不同类型，调用不同的方法。请保证传入的正确性。
	switch v := picFile.(type) {
	case *multipart.FileHeader:
		info, err = manager.UploadPicture(v, uploadPathPrefix)
	case string:
		info, err = manager.UploadPictureByURL(v, uploadPathPrefix, PictureUploadRequest.PicName)
	default:
		return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "参数错误")
	}
	if err != nil {
		return nil, err
	}
	//构造插入数据库的实体
	pic := &entity.Picture{
		URL:          info.URL,
		ThumbnailURL: info.ThumbnailURL,
		Name:         info.PicName,
		PicSize:      info.PicSize,
		PicWidth:     info.PicWidth,
		PicHeight:    info.PicHeight,
		PicScale:     info.PicScale,
		PicFormat:    info.PicFormat,
		UserID:       loginUser.ID,
		EditTime:     time.Now(),
		SpaceID:      PictureUploadRequest.SpaceID, //指定空间id
	}
	//补充审核校验参数
	s.FillReviewParamsInPic(pic, loginUser)
	//若是更新，则需要更新ID
	if picId != 0 {
		pic.ID = picId
	}
	//开启事务
	tx := s.PictureRepo.BeginTransaction()
	//进行插入或者更新操作，即save
	errr := s.PictureRepo.SavePicture(tx, pic)
	if errr != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库错误")
	}
	//修改空间的额度
	if space != nil {
		//设置更新字段
		updateMap := make(map[string]interface{}, 2)
		updateMap["total_count"] = gorm.Expr("total_count + 1")
		updateMap["total_size"] = gorm.Expr("total_size + ?", pic.PicSize)
		err := NewSpaceService().SpaceRepo.UpdateSpaceById(tx, space.ID, updateMap)
		if err != nil {
			return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库错误")
		}
	}
	//提交事务
	errr = tx.Commit().Error
	if err != nil {
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
	//补充审核字段条件
	if req.ReviewStatus != nil {
		query = query.Where("review_status = ?", *req.ReviewStatus)
	}
	if req.ReviewMessage != "" {
		query = query.Where("review_message LIKE ?", "%"+req.ReviewMessage+"%")
	}
	if req.ReviewerID != 0 {
		query = query.Where("reviewer_id = ?", req.ReviewerID)
	}
	if req.SpaceID != 0 {
		query = query.Where("space_id = ?", req.SpaceID)
	}
	if req.IsNullSpaceID {
		query = query.Where("space_id IS NULL")
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
	user, err := repository.NewUserRepository().FindById(nil, Picture.UserID)
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
			user, err := repository.NewUserRepository().FindById(nil, Picture.UserID)
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
	Picture, err := s.PictureRepo.FindById(nil, id)
	if err != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库错误")
	}
	if Picture == nil {
		return nil, ecode.GetErrWithDetail(ecode.NOT_FOUND_ERROR, "图片不存在")
	}
	return Picture, nil
}

// 删除图片，会校验权限
func (s *PictureService) DeletePicture(loginUser *entity.User, deleReq *common.DeleteRequest) *ecode.ErrorWithCode {
	//判断id图片是否存在
	oldPic, err := s.GetPictureById(deleReq.Id)
	if err != nil {
		return err
	}
	//权限校验
	if err := s.CheckPictureAuth(loginUser, oldPic); err != nil {
		return err
	}
	//开启事务
	tx := s.PictureRepo.BeginTransaction()
	//进行删除图片操作
	errr := s.PictureRepo.DeleteById(tx, deleReq.Id)
	if errr != nil {
		return ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库错误")
	}
	var space *entity.Space
	if oldPic.SpaceID != 0 {
		space, errr = repository.NewSpaceRepository().GetSpaceById(nil, oldPic.SpaceID)
		if errr != nil {
			return ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库错误")
		}
	}
	//修改空间的额度
	if space != nil {
		//设置更新字段
		updateMap := make(map[string]interface{}, 2)
		updateMap["total_count"] = gorm.Expr("total_count - 1")
		updateMap["total_size"] = gorm.Expr("total_size - ?", oldPic.PicSize)
		err := NewSpaceService().SpaceRepo.UpdateSpaceById(tx, space.ID, updateMap)
		if err != nil {
			return ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库错误")
		}
	}
	//提交事务
	errr = tx.Commit().Error
	if err != nil {
		return ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库错误")
	}
	return nil
}

// 更新图片，会进行权限校验
func (s *PictureService) UpdatePicture(updateReq *reqPicture.PictureUpdateRequest, loginUser *entity.User) *ecode.ErrorWithCode {
	//查询图片是否存在
	oldPic, err := s.GetPictureById(updateReq.ID)
	if err != nil {
		return err
	}
	//权限校验
	if err := s.CheckPictureAuth(loginUser, oldPic); err != nil {
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
	//填充审核参数
	s.FillReviewParamsInMap(oldPic, loginUser, updateMap)
	//更新
	if err := s.PictureRepo.UpdateById(nil, updateReq.ID, updateMap); err != nil {
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

// 分页查询图片视图（带缓存、多级缓存模式ristretto + redis）
func (s *PictureService) ListPictureVOByPageWithCache(req *reqPicture.PictureQueryRequest) (*resPicture.ListPictureVOResponse, *ecode.ErrorWithCode) {
	//获取redis客户端，和本地缓存
	redisClient := rds.GetRedisClient()
	localCache := cache.GetCache()
	// 将req转为json字符串，并用md5进行压缩
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "参数序列化失败")
	}
	//进一步将请求转化为缓存key
	hash := md5.Sum(reqBytes)
	m5Str := hex.EncodeToString(hash[:])
	cacheKey := fmt.Sprintf("chg:ListPictureVOByPage:%s", m5Str)
	//尝试获取缓存
	//1.本地缓存获取
	dataInterface, found := localCache.Get(cacheKey)
	if found && dataInterface != nil {
		//断言，保证数据为Byte数组
		dataBytes, _ := dataInterface.([]byte)
		//本地缓存命中，直接返回
		var cachedList resPicture.ListPictureVOResponse
		if err := json.Unmarshal(dataBytes, &cachedList); err == nil {
			log.Println("本地缓存命中，数据成功返回")
			return &cachedList, nil
		}
	}
	//2.分布式缓存获取
	cachedData, err := redisClient.Get(context.Background(), cacheKey).Result()
	if rds.IsNilErr(err) {
		log.Println("缓存未命中，查询数据库...")
	} else if err != nil {
		log.Printf("Redis 读取失败: %v\n", err) // 仅做警告，不中断流程
	} else if cachedData != "" {
		var cachedList resPicture.ListPictureVOResponse
		if err := json.Unmarshal([]byte(cachedData), &cachedList); err == nil {
			log.Println("缓存命中，数据成功返回")
			return &cachedList, nil
		} else {
			log.Println("缓存解析失败，重新查询数据库")
		}
	}

	//缓存未击中，正常流程，并将结果放入缓存

	data, Eerr := s.ListPictureVOByPage(req)
	if Eerr != nil {
		return nil, Eerr
	}
	//数据序列化，加入缓存中，允许存储空值
	dataBytes, err := json.Marshal(data)
	if err != nil {
		//序列化失败，不影响正常流程
		log.Println("数据序列化失败，错误为", err)
		return data, nil
	}
	//设置过期时间，为5分钟~10分钟
	expireTime := time.Duration(rand.IntN(300)+300) * time.Second
	expireTime2 := time.Duration(rand.IntN(200)+300) * time.Second
	_, err = redisClient.Set(context.Background(), cacheKey, dataBytes, expireTime).Result()
	if err != nil {
		log.Println("设置缓存失败，错误为", err)
	}
	localCache.SetWithTTL(cacheKey, dataBytes, 1, expireTime2)
	//返回数据
	return data, nil
}

// 图片审核功能，需要记录审核用户
func (s *PictureService) DoPictureReview(req *reqPicture.PictureReviewRequest, user *entity.User) *ecode.ErrorWithCode {
	//校验参数
	if req.ID == 0 || req.ReviewStatus == nil || !consts.ReviewValueExist(*req.ReviewStatus) {
		return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "参数错误")
	}
	//判断图片是否存在
	oldPic, err := s.PictureRepo.FindById(nil, req.ID)
	if err != nil {
		return ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库错误")
	}
	if oldPic == nil {
		return ecode.GetErrWithDetail(ecode.NOT_FOUND_ERROR, "图片不存在")
	}
	//校验审核状态是否重复
	//若当前请求的状态，和图片原有状态一致，返回重复审核异常
	if oldPic.ReviewStatus == *req.ReviewStatus {
		return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "请勿重复审核")
	}
	//数据库操作

	//记录要更新的值，防止全部更新效率过低
	updateMap := make(map[string]interface{}, 8)
	updateMap["review_status"] = *req.ReviewStatus
	updateMap["reviewer_id"] = user.ID
	updateMap["review_time"] = time.Now()
	updateMap["review_message"] = req.ReviewMessage
	//执行更新
	if err := s.PictureRepo.UpdateById(nil, req.ID, updateMap); err != nil {
		return ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库错误")
	}
	return nil
}

// 填充审核参数到指定的map中
func (s *PictureService) FillReviewParamsInMap(Pic *entity.Picture, LoginUser *entity.User, UpdateMap map[string]interface{}) {
	if LoginUser.UserRole == consts.ADMIN_ROLE {
		UpdateMap["review_status"] = consts.PASS
		UpdateMap["reviewer_id"] = LoginUser.ID
		UpdateMap["review_time"] = time.Now()
		UpdateMap["review_message"] = "管理员自动过审"
	} else {
		UpdateMap["review_status"] = consts.REVIEWING
	}
}

// 填充审核参数到指定的Pic中
func (s *PictureService) FillReviewParamsInPic(Pic *entity.Picture, LoginUser *entity.User) {
	if LoginUser.UserRole == consts.ADMIN_ROLE {
		Pic.ReviewStatus = consts.PASS
		Pic.ReviewerID = LoginUser.ID
		Pic.ReviewTime = time.Now()
		Pic.ReviewMessage = "管理员自动过审"
	} else {
		Pic.ReviewStatus = consts.REVIEWING
	}
}

// 批量爬取图片，返回成功数量。
func (s *PictureService) UploadPictureByBatch(req *reqPicture.PictureUploadByBatchRequest, loginUser *entity.User) (int, *ecode.ErrorWithCode) {
	//1.校验参数
	if req.Count > 30 {
		return 0, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "一次最多上传30张图片")
	}
	if req.NamePrefix == "" {
		req.NamePrefix = req.SearchText
	}
	//2.抓取内容
	//searchText需要编码，否则无法抓取中文
	encodedSearchText := url.QueryEscape(req.SearchText)
	//设置一定的页面偏移量
	randInt := rand.IntN(100)
	fetchUrl := fmt.Sprintf("https://cn.bing.com/images/async?q=%s&mmasync=1&first=%d", encodedSearchText, randInt)
	//创建链接
	res, err := http.Get(fetchUrl)
	if err != nil {
		return 0, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "网络请求失败")
	}
	defer res.Body.Close()
	//3.解析内容
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println("解析失败，错误为", err)
		return 0, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "解析失败")
	}
	//提取图片的div
	div := doc.Find(".dgControl").First()
	if div.Length() == 0 {
		log.Println("解析失败，错误为", err)
		return 0, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "解析失败")
	}
	//遍历图片标签
	uploadCount := 0
	div.Find("img.mimg").EachWithBreak(func(i int, img *goquery.Selection) bool {
		//获取src属性，即图片URL
		fileUrl, exists := img.Attr("src")
		if !exists || strings.TrimSpace(fileUrl) == "" {
			log.Println("当前链接为空，已跳过")
			return true // 继续循环
		}
		//去掉url里面的参数，获取原始的图片地址
		if idx := strings.Index(fileUrl, "?"); idx != -1 {
			fileUrl = fileUrl[:idx]
		}
		//4.上传图片
		//编写一个请求，模拟前端调用上传
		uploadReq := &reqPicture.PictureUploadRequest{
			FileUrl: fileUrl,
			PicName: req.NamePrefix,
		}
		if _, err := s.UploadPicture(fileUrl, uploadReq, loginUser); err != nil {
			log.Println("上传失败，错误为", err)
		} else {
			log.Println("上传成功")
			uploadCount++
		}
		return uploadCount < req.Count
	})
	return uploadCount, nil
}

//增加的空间逻辑

// 校验操作图片权限，公共图库仅本人或管理员可以操作，私人图库仅空间管理员可以操作
func (s *PictureService) CheckPictureAuth(loginUser *entity.User, picture *entity.Picture) *ecode.ErrorWithCode {
	//公共图库，仅本人或管理员可以操作
	if picture.SpaceID == 0 {
		if loginUser.ID != picture.UserID && loginUser.UserRole != consts.ADMIN_ROLE {
			return ecode.GetErrWithDetail(ecode.NO_AUTH_ERROR, "没有权限")
		}
	} else {
		//私人图库，仅空间管理员可以操作
		if picture.UserID != loginUser.ID {
			return ecode.GetErrWithDetail(ecode.NO_AUTH_ERROR, "没有权限")
		}
	}
	return nil
}
