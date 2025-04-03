package analyze

//空间图片分类分析请求
type SpaceUserAnalyzeResponse struct {
	Period string `json:"period"` //时间周期
	Count  int64  `json:"count"`  //周期内上传的图片数量
}
