package space

import "chg/internal/common"

// SpaceQueryRequest 表示空间的查询请求
type SpaceQueryRequest struct {
	common.PageRequest        // 嵌入 PageRequest 以支持分页字段
	ID                 uint64 `json:"id,string" swaggertype:"string"`     // 空间 ID
	UserID             uint64 `json:"userId,string" swaggertype:"string"` // 用户 ID
	SpaceName          string `json:"spaceName"`                          // 空间名称
	SpaceLevel         int    `json:"spaceLevel"`                         // 空间级别：0-普通版 1-专业版 2-旗舰版
}
