package picture

type PictureEditByBatchRequest struct {
	PictureIdList []uint64 `json:"pictureIdList" swaggertype:"array,string"` // 图片ID列表
	SpaceID       uint64   `json:"spaceId,string" swaggertype:"string"`      //空间ID
	Category      string   `json:"category"`                                 //分类
	Tags          []string `json:"tags"`                                     //标签
	NameRule      string   `json:"nameRule"`                                 //名称规则，暂时只支持“名称{序号}的形式，序号将会自动递增”
}
