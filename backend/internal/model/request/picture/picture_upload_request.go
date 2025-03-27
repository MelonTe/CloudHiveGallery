package picture

type PictureUploadRequest struct {
	ID      uint64 `json:"id,string" swaggertype:"string"`      //图片ID
	FileUrl string `json:"fileUrl"`                             //图片地址
	PicName string `json:"picName"`                             //图片名称
	SpaceID uint64 `json:"spaceId,string" swaggertype:"string"` //空间ID
}
