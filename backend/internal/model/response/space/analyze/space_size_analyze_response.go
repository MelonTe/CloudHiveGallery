package analyze

//空间分类分析响应
type SpaceSizeAnalyzeResponse struct {
	SizeRange string `json:"sizeRange"` //大小范围，格式为"<100KB","100KB-500KB","500KB-1MB",">1MB"
	Count     int64  `json:"count"`     //分类数量
}
