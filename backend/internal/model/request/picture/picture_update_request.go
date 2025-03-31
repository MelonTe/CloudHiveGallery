package picture

//供管理员使用的图片更新请求
type PictureUpdateRequest struct {
	ID           uint64   `json:"id,string" swaggertype:"string"`
	Name         string   `json:"name"`
	Introduction string   `json:"introduction"`
	Category     string   `json:"category"`
	Tags         []string `json:"tags"`
	SpaceId      uint64   `json:"spaceId,string" swaggertype:"string"` //空间ID
}
