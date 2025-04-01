package model

type APIResponse struct {
	Status int       `json:"status"`
	Data   ImageData `json:"data"`
}

// 定义结构体，用于从FirstURL获取信息
type ImageData struct {
	List []ImageSearchResult `json:"list"`
}

// 定义图片搜索结果结构体
type ImageSearchResult struct {
	ThumbURL string `json:"thumbURL"` // 缩略图地址
	FromURL  string `json:"fromURL"`  // 来源地址
}
