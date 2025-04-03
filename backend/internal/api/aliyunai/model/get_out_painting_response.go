package model

//查询扩图任务状态响应
type GetOutPaintingResponse struct {
	RequestID string           `json:"request_id"`      // 请求唯一标识
	Output    TaskDetailOutput `json:"output"`          // 任务输出信息（一定包含）
	Usage     *Usage           `json:"usage,omitempty"` // 图像统计信息（仅在成功时返回）
}

// Output 任务输出信息
type TaskDetailOutput struct {
	TaskID         string       `json:"task_id"`                    // 任务的唯一标识符
	TaskStatus     string       `json:"task_status"`                // 任务状态：PENDING、RUNNING、SUSPENDED、SUCCEEDED、FAILED、UNKNOWN
	TaskMetrics    *TaskMetrics `json:"task_metrics,omitempty"`     // 任务结果统计（执行中时返回）
	SubmitTime     string       `json:"submit_time,omitempty"`      // 任务提交时间（成功或失败时返回）
	ScheduledTime  string       `json:"scheduled_time,omitempty"`   // 任务调度时间（成功或失败时返回）
	EndTime        string       `json:"end_time,omitempty"`         // 任务完成时间（成功或失败时返回）
	OutputImageURL string       `json:"output_image_url,omitempty"` // 输出图像的 URL（成功时返回）
	Code           string       `json:"code,omitempty"`             // 错误码（失败时返回）
	Message        string       `json:"message,omitempty"`          // 错误信息（失败时返回）
}

// TaskMetrics 任务统计信息（仅在任务执行中时返回）
type TaskMetrics struct {
	Total     int `json:"TOTAL"`     // 总任务数
	Succeeded int `json:"SUCCEEDED"` // 成功任务数
	Failed    int `json:"FAILED"`    // 失败任务数
}

// Usage 图像统计信息（仅在成功时返回）
type Usage struct {
	ImageCount int `json:"image_count"` // 生成的图片数量
}
