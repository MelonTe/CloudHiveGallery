package analyze

//空间标签分析响应
type SpaceTagAnalyzeResponse struct {
	Tag   string `json:"tag"`   //标签名称
	Count int64  `json:"count"` //标签数量
}
