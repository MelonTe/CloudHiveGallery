package space

// SpaceUpdateRequest 更新空间请求
type SpaceUpdateRequest struct {
	ID         uint64 `json:"id,string" swaggertype:"string"` // Space ID
	SpaceName  string `json:"spaceName"`                      // Space name
	SpaceLevel int    `json:"spaceLevel"`                     // Space level: 0-普通版 1-专业版 2-旗舰版
	MaxSize    int64  `json:"maxSize"`                        // Maximum total size of space images
	MaxCount   int64  `json:"maxCount"`                       // Maximum number of space images
}
