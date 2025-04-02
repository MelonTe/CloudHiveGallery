package picture

type PictureSearchByColorRequest struct {
	PicColor string `json:"picColor"`                            // 图片颜色
	SpaceID  uint64 `json:"spaceId,string" swaggertype:"string"` //空间ID
}
