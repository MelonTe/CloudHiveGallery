package space

// SpaceAddRequest 添加空间请求
type SpaceAddRequest struct {
	SpaceName  string `json:"spaceName"`  // 空间名称
	SpaceLevel int    `json:"spaceLevel"` // 空间级别：0-普通版 1-专业版 2-旗舰版
	SpaceType  int    `json:"spaceType"`  // 空间类型：0-个人空间 1-团队空间
}
