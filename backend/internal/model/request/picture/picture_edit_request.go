package picture

//供用户使用的图片更新请求
type PictureEditRequest struct {
	ID           uint64   `json:"id,string" swaggertype:"string"`
	Name         string   `json:"name"`
	Introduction string   `json:"introduction"`
	Category     string   `json:"category"`
	Tags         []string `json:"tags"`
	SpaceId      uint64   `json:"spaceId,string" swaggertype:"string"` //空间ID
}
