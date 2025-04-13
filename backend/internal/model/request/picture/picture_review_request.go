package picture

//图片审核请求
type PictureReviewRequest struct {
	ID            uint64 `json:"id,string" swaggertype:"string"` //图片ID
	ReviewStatus  *int   `json:"reviewStatus"`                   //审核状态
	ReviewMessage string `json:"reviewMessage"`                  //审核信息
}
