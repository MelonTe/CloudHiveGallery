package common

//分页响应的通用结构
type PageResponse struct {
	Total   int `json:"total"`   //总记录数
	Current int `json:"current"` //当前页数
	Pages   int `json:"pages"`   //页面大小
	Size    int `json:"size"`    //页面大小
}
