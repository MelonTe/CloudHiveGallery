package analyze

//空间图片分类分析请求
type SpaceUserAnalyzeRequest struct {
	SpaceAnalyzeRequest
	UserID        uint64 `json:"userId,string" swaggertype:"string"` //用户ID
	TimeDimension string `json:"timeDimension"`                      //时间维度：day/week/month
}
