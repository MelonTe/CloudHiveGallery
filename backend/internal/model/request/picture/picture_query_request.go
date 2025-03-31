package picture

import "chg/internal/common"

type PictureQueryRequest struct {
	ID           uint64   `json:"id,string" swaggertype:"string"` //图片ID
	Name         string   `json:"name"`
	Introduction string   `json:"introduction"`
	Category     string   `json:"category"`
	Tags         []string `json:"tags"`
	PicSize      int64    `json:"picSize"`
	PicWidth     int      `json:"picWidth"`
	PicHeight    int      `json:"picHeight"`
	PicScale     float64  `json:"picScale"`
	PicFormat    string   `json:"picFormat"`
	UserID       uint64   `json:"userId,string" swaggertype:"string"` //图片上传人信息
	SearchText   string   `json:"searchText"`                         //搜索词
	common.PageRequest
	//新增审核字段
	ReviewStatus  *int   `json:"reviewStatus,string" swaggertype:"string"` //审核状态
	ReviewerID    uint64 `json:"reviewerId,string" swaggertype:"string"`   //审核人ID
	ReviewMessage string `json:"reviewMessage"`
	//新增空间筛选字段
	SpaceID       uint64 `json:"spaceId,string" swaggertype:"string"` //空间ID
	IsNullSpaceID bool   `json:"isNullSpaceId"`                       //是否查询空间ID为空的图片
}
