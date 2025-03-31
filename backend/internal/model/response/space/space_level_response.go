package space

//返回给前端空间的等级

type SpaceLevelResponse struct {
	Value    int    `json:"value"`    //空间的等级
	Text     string `json:"text"`     //空间的等级名称
	MaxCount int64  `json:"maxCount"` //空间图片的最大数量
	MaxSize  int64  `json:"maxSize"`  //空间图片的最大总大小
}
