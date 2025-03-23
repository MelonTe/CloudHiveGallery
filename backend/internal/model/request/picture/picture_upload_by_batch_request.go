package picture

//图片批量上传请求
type PictureUploadByBatchRequest struct {
	SearchText string `json:"searchText"` //搜索词
	Count      int    `json:"count"`      //图片数量
	NamePrefix string `json:"namePrefix"` //图片名称前缀，默认为SearchText
}
