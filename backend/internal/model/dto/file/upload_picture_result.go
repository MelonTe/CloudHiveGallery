package file

//定义上传图片返回结果

//用于接收图片解析信息
type UploadPictureResult struct {
	URL          string  `json:"url"`
	ThumbnailURL string  `json:"thumbnailUrl"`
	PicName      string  `json:"picName"`
	PicSize      int64   `json:"picSize"`
	PicWidth     int     `json:"picWidth"`
	PicHeight    int     `json:"picHeight"`
	PicScale     float64 `json:"picScale"`
	PicFormat    string  `json:"picFormat"`
}
