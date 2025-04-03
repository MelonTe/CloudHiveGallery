package analyze

//空间分类分析响应
type SpaceCategoryAnalyzeResponse struct {
	Category  string `json:"category"`  //分类名称
	Count     int64  `json:"count"`     //分类数量
	TotalSize int64  `json:"totalSize"` //分类总大小
}
