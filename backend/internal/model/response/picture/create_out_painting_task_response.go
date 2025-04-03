package picture

import "chg/internal/api/aliyunai/model"

//扩图处理结果响应
type CreateOutPaintingTaskResponse struct {
	Output    *Output `json:"output,omitempty"`  // 任务输出信息（成功时返回）
	Code      string  `json:"code,omitempty"`    // 错误码（失败时返回）
	Message   string  `json:"message,omitempty"` // 错误信息（失败时返回）
	RequestID string  `json:"requestId"`         // 请求唯一标识符
}

// Output 任务输出信息（仅在成功时返回）
type Output struct {
	TaskStatus string `json:"taskStatus"` // 任务状态：PENDING、RUNNING、SUSPENDED、SUCCEEDED、FAILED、UNKNOWN
	TaskID     string `json:"taskId"`     // 任务的唯一标识符
}

//将阿里云的任务状态转化为要发送给前端的任务状态
//阿里云的响应转化为要发送前端的响应
func AOutPaintResToF(response *model.CreateOutPaintingTaskResponse) *CreateOutPaintingTaskResponse {
	if response == nil {
		return nil
	}
	resp := &CreateOutPaintingTaskResponse{
		Code:      response.Code,
		Message:   response.Message,
		RequestID: response.RequestID,
	}
	if response.Output != nil {
		resp.Output = &Output{
			TaskStatus: response.Output.TaskStatus,
			TaskID:     response.Output.TaskID,
		}
	}
	return resp
}
