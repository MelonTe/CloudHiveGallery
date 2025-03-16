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
	UserID       uint64   `json:"userId,string" swaggertype:"string"`
	SearchText   string   `json:"searchText"` //搜索词
	common.PageRequest
}
