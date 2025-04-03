package model

//创建扩图任务请求，发送给阿里云的请求体
type CreateOutPaintingTaskRequest struct {
	Model      string          `json:"model" binding:"required"`
	Input      ImageInput      `json:"input" binding:"required"`
	Parameters ImageParameters `json:"parameters" binding:"required"`
}

type ImageInput struct {
	ImageURL string `json:"image_url" binding:"required"`
}

type ImageParameters struct {
	Angle          int     `json:"angle,omitempty"`
	OutputRatio    string  `json:"output_ratio,omitempty"`
	XScale         float32 `json:"x_scale,omitempty"`
	YScale         float32 `json:"y_scale,omitempty"`
	TopOffset      int     `json:"top_offset,omitempty"`
	BottomOffset   int     `json:"bottom_offset,omitempty"`
	LeftOffset     int     `json:"left_offset,omitempty"`
	RightOffset    int     `json:"right_offset,omitempty"`
	BestQuality    bool    `json:"best_quality,omitempty"`
	LimitImageSize bool    `json:"limit_image_size,omitempty"`
	AddWatermark   bool    `json:"add_watermark,omitempty"`
}
