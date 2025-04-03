package analyze

//空间资源使用分析响应
// SpaceUsageAnalyzeResponse 表示空间资源使用分析的响应结构
type SpaceUsageAnalyzeResponse struct {
	UsedSize        int64   `json:"usedSize"`                  // 已使用的空间大小
	MaxSize         int64   `json:"maxSize,omitempty"`         // 最大空间大小
	SizeUsageRatio  float64 `json:"sizeUsageRatio,omitempty"`  // 空间使用比例
	UsedCount       int64   `json:"usedCount"`                 // 已使用的资源数量
	MaxCount        int64   `json:"maxCount,omitempty"`        // 最大资源数量
	CountUsageRatio float64 `json:"countUsageRatio,omitempty"` // 资源数量使用比例
}
