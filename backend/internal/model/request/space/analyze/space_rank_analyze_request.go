package analyze

//通用空间分析请求
type SpaceRankAnalyzeRequest struct {
	TopN int `json:"top_n"` //排名前N的空间
}
