package service

import (
	"chg/internal/consts"
	"chg/internal/ecode"
	"chg/internal/model/entity"
	reqSpaceAnalyze "chg/internal/model/request/space/analyze"
	resSpaceAnalyze "chg/internal/model/response/space/analyze"
	"chg/internal/repository"
	"chg/pkg/db"
	"encoding/json"
	"math"
	"sort"

	"gorm.io/gorm"
)

type SpaceAnalyzeService struct {
	SpaceAnalyzeRepo *repository.SpaceAnalyzeRepository
}

func NewSpaceAnalyzeService() *SpaceAnalyzeService {
	return &SpaceAnalyzeService{
		SpaceAnalyzeRepo: repository.NewSpaceAnalyzeRepository(),
	}
}

// 校验空间分析权限
func (s *SpaceAnalyzeService) CheckSpaceAnalyzeAuth(SpaceAnalyzeReq *reqSpaceAnalyze.SpaceAnalyzeRequest, loginUser *entity.User) *ecode.ErrorWithCode {
	//校验查询列表
	if SpaceAnalyzeReq.QueryAll || SpaceAnalyzeReq.QueryPublic {
		//需要校验是否是管理员
		if loginUser.UserRole != consts.ADMIN_ROLE {
			return ecode.GetErrWithDetail(ecode.NO_AUTH_ERROR, "没有权限")
		}
	} else {
		//私有空间权限校验
		if SpaceAnalyzeReq.SpaceID <= 0 {
			return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "空间ID不能为空")
		}
		space, err := NewSpaceService().GetSpaceById(SpaceAnalyzeReq.SpaceID)
		if err != nil {
			return err
		}
		//仅管理员或空间管理者可以查询空间分析
		if space.UserID != loginUser.ID && loginUser.UserRole != consts.ADMIN_ROLE {
			return ecode.GetErrWithDetail(ecode.NO_AUTH_ERROR, "没有权限")
		}
	}
	return nil
}

// 填充空间分析链式查询条件
func (s *SpaceAnalyzeService) FillAnalyzeQueryWrapper(PictureQuery *gorm.DB, req *reqSpaceAnalyze.SpaceAnalyzeRequest) (*gorm.DB, *ecode.ErrorWithCode) {
	query := PictureQuery.Session(&gorm.Session{})
	//全空间分析
	if req.QueryAll {
		//全空间分析不需要任何条件
		return query, nil
	}
	//公共图库分析
	if req.QueryPublic {
		//需要查询spaceId为null的图片
		query = query.Where("space_id IS NULL")
		return query, nil
	}
	//特定空间分析
	if req.SpaceID > 0 {
		query = query.Where("space_id = ?", req.SpaceID)
		return query, nil
	}
	//未指定查询范围
	return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "查询范围不明确")
}

// 空间使用情况分析
func (s *SpaceAnalyzeService) GetSpaceUsageAnalyze(req *reqSpaceAnalyze.SpaceUsageAnalyzeRequest, loginUser *entity.User) (*resSpaceAnalyze.SpaceUsageAnalyzeResponse, *ecode.ErrorWithCode) {
	//校验参数
	//全空间或公共图库需要从picture查询
	res := &resSpaceAnalyze.SpaceUsageAnalyzeResponse{}
	if req.QueryAll || req.QueryPublic {
		//权限校验
		if err := s.CheckSpaceAnalyzeAuth(&req.SpaceAnalyzeRequest, loginUser); err != nil {
			return nil, err
		}
		query := db.LoadDB()
		//只获取pic_size字段
		var picSize []int64
		//补充空间字段
		query, err := s.FillAnalyzeQueryWrapper(query, &req.SpaceAnalyzeRequest)
		if err != nil {
			return nil, err
		}
		//查询
		query.Model(&entity.Picture{}).Pluck("pic_size", &picSize)
		//字段填充
		sumSize := int64(0)
		for _, size := range picSize {
			sumSize += size
		}
		res.UsedCount = int64(len(picSize))
		res.UsedSize = sumSize
		return res, nil
	} else {
		//私有空间可以从Space查询
		//参数校验和权限校验
		if err := s.CheckSpaceAnalyzeAuth(&req.SpaceAnalyzeRequest, loginUser); err != nil {
			return nil, err
		}
		space, err := NewSpaceService().GetSpaceById(req.SpaceID)
		if err != nil {
			return nil, err
		}
		res.UsedCount = space.TotalCount
		res.UsedSize = space.TotalSize
		res.MaxCount = space.MaxCount
		res.MaxSize = space.MaxSize
		res.SizeUsageRatio = math.Round(float64(space.TotalSize)/float64(space.MaxSize)*100*100) / 100
		res.CountUsageRatio = math.Round(float64(space.TotalCount)/float64(space.MaxCount)*100*100) / 100
		return res, nil
	}
}

// 获取空间分类统计分析
func (s *SpaceAnalyzeService) GetSpaceCategoryAnalyze(req *reqSpaceAnalyze.SpaceCategoryAnalyzeRequest, loginUser *entity.User) ([]resSpaceAnalyze.SpaceCategoryAnalyzeResponse, *ecode.ErrorWithCode) {
	//权限校验
	if err := s.CheckSpaceAnalyzeAuth(&req.SpaceAnalyzeRequest, loginUser); err != nil {
		return nil, err
	}
	//获取查询对象
	query := db.LoadDB()
	query = query.Model(&entity.Picture{})
	//补充空间字段
	query, err := s.FillAnalyzeQueryWrapper(query, &req.SpaceAnalyzeRequest)
	if err != nil {
		return nil, err
	}
	//查询分类统计
	var result []resSpaceAnalyze.SpaceCategoryAnalyzeResponse
	//SQL语句，匹配结构体字段昵称的snake_case形式
	if originErr := query.Select("COALESCE(NULLIF(category,''),'未分类') AS category, COUNT(*) AS count, SUM(pic_size) as total_size").
		Group("category").
		Scan(&result).Error; originErr != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库查询失败")
	}
	return result, nil
}

// 获取空间标签统计分析
func (s *SpaceAnalyzeService) GetSpaceTagAnalyze(req *reqSpaceAnalyze.SpaceTagAnalyzeRequest, loginUser *entity.User) ([]resSpaceAnalyze.SpaceTagAnalyzeResponse, *ecode.ErrorWithCode) {
	//权限校验
	if err := s.CheckSpaceAnalyzeAuth(&req.SpaceAnalyzeRequest, loginUser); err != nil {
		return nil, err
	}
	//获取查询对象
	query := db.LoadDB()
	query = query.Model(&entity.Picture{})
	//补充空间字段
	query, err := s.FillAnalyzeQueryWrapper(query, &req.SpaceAnalyzeRequest)
	if err != nil {
		return nil, err
	}
	//查询原始标签
	var OriginTags []string
	if originErr := query.
		Where("tags IS NOT NULL").
		Where("tags != ''").
		Pluck("tags", &OriginTags).Error; originErr != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库查询失败")
	}
	TagCount := make(map[string]int64)
	//解析标签
	for _, tags := range OriginTags {
		var tagList []string
		if err := json.Unmarshal([]byte(tags), &tagList); err != nil {
			return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "标签解析失败")
		}
		for _, tag := range tagList {
			TagCount[tag]++
		}
	}
	var result []resSpaceAnalyze.SpaceTagAnalyzeResponse
	for tag, count := range TagCount {
		res := resSpaceAnalyze.SpaceTagAnalyzeResponse{
			Tag:   tag,
			Count: count,
		}
		result = append(result, res)
	}
	return result, nil
}

// 获取空间大小统计分析
func (s *SpaceAnalyzeService) GetSpaceSizeAnalyze(req *reqSpaceAnalyze.SpaceSizeAnalyzeRequest, loginUser *entity.User) ([]resSpaceAnalyze.SpaceSizeAnalyzeResponse, *ecode.ErrorWithCode) {
	//权限校验
	if err := s.CheckSpaceAnalyzeAuth(&req.SpaceAnalyzeRequest, loginUser); err != nil {
		return nil, err
	}
	//获取查询对象
	query := db.LoadDB()
	query = query.Model(&entity.Picture{})
	//补充空间字段
	query, err := s.FillAnalyzeQueryWrapper(query, &req.SpaceAnalyzeRequest)
	if err != nil {
		return nil, err
	}
	//查询大小统计
	var picsSize []int64
	if originErr := query.Pluck("pic_size", &picsSize).Error; originErr != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库查询失败")
	}
	sort.Slice(picsSize, func(i, j int) bool {
		return picsSize[i] < picsSize[j]
	})
	var result []resSpaceAnalyze.SpaceSizeAnalyzeResponse
	//区间划分为：[0,100KB), [100KB,500KB), [500KB,1MB),[1MB,*]
	target := []int{100 * 1024, 500 * 1024, 1024 * 1024} //边界定义
	targetToRange := []string{"<100KB", "100KB-500KB", "500KB-1MB", ">1MB"}
	needSub := 0
	for i := 0; i < len(target); i++ {
		//找到第一个大于等于target[i]的元素
		index := sort.Search(len(picsSize), func(j int) bool {
			if picsSize[j] < int64(target[i]) {
				return false
			}
			return true
		})
		if index < len(picsSize) {
			//找到了，index左侧都是<target[i]的元素
			result = append(result, resSpaceAnalyze.SpaceSizeAnalyzeResponse{
				SizeRange: targetToRange[i],
				Count:     int64(index - needSub),
			})
			//若是最后一个，要处理剩余的元素
			if i == len(target)-1 {
				result = append(result, resSpaceAnalyze.SpaceSizeAnalyzeResponse{
					SizeRange: targetToRange[i+1],
					Count:     int64(len(picsSize) - index),
				})
			}
			needSub = index
		} else {
			//没有找到，说明该区间没有元素，直接添加
			result = append(result, resSpaceAnalyze.SpaceSizeAnalyzeResponse{
				SizeRange: targetToRange[i],
				Count:     0,
			})
			if i == len(target)-1 {
				//最后一个区间，添加剩余元素
				result = append(result, resSpaceAnalyze.SpaceSizeAnalyzeResponse{
					SizeRange: targetToRange[i+1],
					Count:     0,
				})
			}
		}
	}
	return result, nil
}

// 获取规定时间周期内，用户上传图片的情况
func (s *SpaceAnalyzeService) GetSpaceUserAnalyze(req *reqSpaceAnalyze.SpaceUserAnalyzeRequest, loginUser *entity.User) ([]resSpaceAnalyze.SpaceUserAnalyzeResponse, *ecode.ErrorWithCode) {
	//权限校验
	if err := s.CheckSpaceAnalyzeAuth(&req.SpaceAnalyzeRequest, loginUser); err != nil {
		return nil, err
	}
	//获取查询对象
	query := db.LoadDB()
	query = query.Model(&entity.Picture{})
	//补充空间字段
	query, err := s.FillAnalyzeQueryWrapper(query, &req.SpaceAnalyzeRequest)
	if err != nil {
		return nil, err
	}
	if req.UserID != 0 {
		query = query.Where("user_id = ?", req.UserID)
	}
	//根据需要分析的时间维度，进行分组
	switch req.TimeDimension {
	case "day":
		//DATE_FORMAT将时间格式化为YYYY-MM-DD
		query = query.Select("DATE_FORMAT(create_time, '%Y-%m-%d') AS period, COUNT(*) AS count")
	case "week":
		//YEARWEEK将时间格式化为第几年的第几周，如202511
		query = query.Select("YEARWEEK(create_time) AS period, COUNT(*) AS count")
	case "month":
		query = query.Select("DATE_FORMAT(create_time, '%Y-%m') AS period, COUNT(*) AS count")
	default:
		return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "时间维度不合法")
	}
	var result []resSpaceAnalyze.SpaceUserAnalyzeResponse
	if originErr := query.Group("period").Order("period").Scan(&result).Error; originErr != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库查询失败")
	}
	return result, nil
}

// 查询空间使用排名前topN的用户
func (s *SpaceAnalyzeService) GetSpaceRankAnalyze(req *reqSpaceAnalyze.SpaceRankAnalyzeRequest, loginUser *entity.User) ([]entity.Space, *ecode.ErrorWithCode) {
	if loginUser.UserRole != consts.ADMIN_ROLE {
		return nil, ecode.GetErrWithDetail(ecode.NO_AUTH_ERROR, "没有权限")
	}
	if req.TopN == 0 {
		req.TopN = 10
	}
	//获取查询对象
	query := db.LoadDB()
	var result []entity.Space
	if originErr := query.Model(&entity.Space{}).Select("id", "space_name", "user_id", "total_size").Order("total_size DESC").Limit(req.TopN).Scan(&result).Error; originErr != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库查询失败")
	}
	return result, nil
}
