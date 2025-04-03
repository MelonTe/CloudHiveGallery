package picture

import "chg/internal/api/aliyunai/model"

// 查询扩图任务状态响应
type GetOutPaintingResponse struct {
	RequestID string           `json:"requestId"`       // 请求唯一标识
	Output    TaskDetailOutput `json:"output"`          // 任务输出信息（一定包含）
	Usage     *Usage           `json:"usage,omitempty"` // 图像统计信息（仅在成功时返回）
}

// Output 任务输出信息
type TaskDetailOutput struct {
	TaskID         string       `json:"taskId"`                   // 任务的唯一标识符
	TaskStatus     string       `json:"taskStatus"`               // 任务状态：PENDING、RUNNING、SUSPENDED、SUCCEEDED、FAILED、UNKNOWN
	TaskMetrics    *TaskMetrics `json:"taskMetrics,omitempty"`    // 任务结果统计（执行中时返回）
	SubmitTime     string       `json:"submitTime,omitempty"`     // 任务提交时间（成功或失败时返回）
	ScheduledTime  string       `json:"scheduledTime,omitempty"`  // 任务调度时间（成功或失败时返回）
	EndTime        string       `json:"endTime,omitempty"`        // 任务完成时间（成功或失败时返回）
	OutputImageURL string       `json:"outputImageUrl,omitempty"` // 输出图像的 URL（成功时返回）
	Code           string       `json:"code,omitempty"`           // 错误码（失败时返回）
	Message        string       `json:"message,omitempty"`        // 错误信息（失败时返回）
}

// TaskMetrics 任务统计信息（仅在任务执行中时返回）
type TaskMetrics struct {
	Total     int `json:"total"`     // 总任务数
	Succeeded int `json:"succeeded"` // 成功任务数
	Failed    int `json:"failed"`    // 失败任务数
}

// Usage 图像统计信息（仅在成功时返回）
type Usage struct {
	ImageCount int `json:"imageCount"` // 生成的图片数量
}

//将阿里云的任务状态转化为要发送给前端的任务状态
func AGetOutPaintResToF(Aresponse *model.GetOutPaintingResponse) *GetOutPaintingResponse {
	if Aresponse == nil {
		return nil
	}
	resp := &GetOutPaintingResponse{
		RequestID: Aresponse.RequestID,
	}
	resp.Output = TaskDetailOutput{
		TaskID:         Aresponse.Output.TaskID,
		TaskStatus:     Aresponse.Output.TaskStatus,
		SubmitTime:     Aresponse.Output.SubmitTime,
		ScheduledTime:  Aresponse.Output.ScheduledTime,
		EndTime:        Aresponse.Output.EndTime,
		OutputImageURL: Aresponse.Output.OutputImageURL,
		Code:           Aresponse.Output.Code,
		Message:        Aresponse.Output.Message,
	}
	if Aresponse.Output.TaskMetrics != nil {
		resp.Output.TaskMetrics = &TaskMetrics{
			Total:     Aresponse.Output.TaskMetrics.Total,
			Succeeded: Aresponse.Output.TaskMetrics.Succeeded,
			Failed:    Aresponse.Output.TaskMetrics.Failed,
		}
	}
	if Aresponse.Usage != nil {
		resp.Usage = &Usage{
			ImageCount: Aresponse.Usage.ImageCount,
		}
	}
	return resp
}
