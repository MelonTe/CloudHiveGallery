package model

//扩图处理结果响应
type CreateOutPaintingTaskResponse struct {
	Output    *Output `json:"output,omitempty"`  // 任务输出信息（成功时返回）
	Code      string  `json:"code,omitempty"`    // 错误码（失败时返回）
	Message   string  `json:"message,omitempty"` // 错误信息（失败时返回）
	RequestID string  `json:"request_id"`        // 请求唯一标识符
}

// Output 任务输出信息（仅在成功时返回）
type Output struct {
	TaskStatus string `json:"task_status"` // 任务状态：PENDING、RUNNING、SUSPENDED、SUCCEEDED、FAILED、UNKNOWN
	TaskID     string `json:"task_id"`     // 任务的唯一标识符
}

// TaskStatusType 定义任务状态的枚举值
// 任务状态包括：排队中、处理中、挂起、成功、失败和未知
const (
	TaskStatusPending   = "PENDING"   // 任务排队中
	TaskStatusRunning   = "RUNNING"   // 任务处理中
	TaskStatusSuspended = "SUSPENDED" // 任务挂起
	TaskStatusSucceeded = "SUCCEEDED" // 任务执行成功
	TaskStatusFailed    = "FAILED"    // 任务执行失败
	TaskStatusUnknown   = "UNKNOWN"   // 任务状态未知或任务不存在
)
