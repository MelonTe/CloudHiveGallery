package picture

//定义前端->后端->阿里云的中间请求体
// 该请求体是为了适配阿里云的API而设计的，主要是将前端传来的参数转换为阿里云需要的格式
import (
	"chg/internal/api/aliyunai/model"
)

// 前端扩图请求体
type CreateOutPaintingTaskRequest struct {
	// 图片ID
	PictureID uint64 `json:"pictureId,string" swaggertype:"string"`

	// 图像处理任务的参数
	Parameters ImageParameters `json:"parameters"`
}

// ImageParameters 表示配置图像处理任务的参数
type ImageParameters struct {
	// 图像旋转角度（单位：度）
	Angle int `json:"angle,omitempty"`

	// 输出图像的宽高比（例如："16:9"）
	OutputRatio string `json:"outputRatio,omitempty"`

	// 图像的水平缩放比例，范围在1.0 ~ 3.0
	XScale float32 `json:"xScale"`

	// 图像的垂直缩放比例，范围在1.0 ~ 3.0
	YScale float32 `json:"yScale"`

	// 图像顶部的偏移量
	TopOffset int `json:"topOffset,omitempty"`

	// 图像底部的偏移量
	BottomOffset int `json:"bottomOffset,omitempty"`

	// 图像左侧的偏移量
	LeftOffset int `json:"leftOffset,omitempty"`

	// 图像右侧的偏移量
	RightOffset int `json:"rightOffset,omitempty"`

	// 是否启用最佳质量
	BestQuality bool `json:"bestQuality,omitempty"`

	// 是否限制图像大小
	LimitImageSize bool `json:"limitImageSize,omitempty"`

	// 是否添加水印
	AddWatermark bool `json:"addWatermark,omitempty"`
}

// 转化为请求阿里云的实体类，主要是兼容驼峰命名和下划线命名的转换
func (r *CreateOutPaintingTaskRequest) ToAliAiRequest(imageURL string) *model.CreateOutPaintingTaskRequest {
	ret := &model.CreateOutPaintingTaskRequest{
		Model: "image-out-painting",
		Input: model.ImageInput{
			ImageURL: imageURL,
		},
		Parameters: model.ImageParameters{
			Angle:          r.Parameters.Angle,
			OutputRatio:    r.Parameters.OutputRatio,
			TopOffset:      r.Parameters.TopOffset,
			BottomOffset:   r.Parameters.BottomOffset,
			LeftOffset:     r.Parameters.LeftOffset,
			RightOffset:    r.Parameters.RightOffset,
			BestQuality:    r.Parameters.BestQuality,
			LimitImageSize: r.Parameters.LimitImageSize,
			AddWatermark:   true, // 阿里云的扩图任务必须添加水印
		},
	}
	//调整缩放比例
	ret.Parameters.XScale = 1.0
	ret.Parameters.YScale = 1.0
	if r.Parameters.XScale <= 3.0 && r.Parameters.XScale >= 1.0 {
		ret.Parameters.XScale = r.Parameters.XScale
	}
	if r.Parameters.YScale <= 3.0 && r.Parameters.YScale >= 1.0 {
		ret.Parameters.YScale = r.Parameters.YScale
	}
	return ret
}
