package analyze

//通用空间分析请求
type SpaceAnalyzeRequest struct {
	SpaceID     uint64 `json:"spaceId,string" swaggertype:"string"` //空间ID
	QueryPublic bool   `json:"queryPublic"`                         //是否查询公开空间
	QueryAll    bool   `json:"queryAll"`                            //是否查询所有空间
}
