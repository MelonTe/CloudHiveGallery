package fetcher

import (
	"chg/config"
	"chg/internal/api/aliyunai/model"
	"chg/internal/ecode"
	"encoding/json"
	"fmt"
	"log"
	"resty.dev/v3"
)

//该包实现了对阿里云API的请求和响应的封装

const (
	CreateOutPaintingTaskURL = "https://dashscope.aliyuncs.com/api/v1/services/aigc/image2image/out-painting" // 扩图请求地址
	GetOutPaintingTaskURL    = "https://dashscope.aliyuncs.com/api/v1/tasks/%s"                               // 获取扩图任务结果地址，需要拼接上task_id
)

// 创建任务
func CreateOutPaintingTask(req *model.CreateOutPaintingTaskRequest) (*model.CreateOutPaintingTaskResponse, *ecode.ErrorWithCode) {
	// 参数校验
	if req == nil {
		return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "请求参数不能为空")
	}

	// 创建Resty客户端
	client := resty.New()

	// 将请求参数序列化为JSON
	body, err := json.Marshal(req)
	if err != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "请求参数序列化失败")
	}
	cfg := config.LoadConfig()
	// 发送POST请求
	resp, err := client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", cfg.AliYunAi.ApiKey)).
		SetHeader("X-DashScope-Async", "enable").
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(CreateOutPaintingTaskURL)
	if err != nil {
		log.Println("请求失败:", err)
		return nil, ecode.GetErrWithDetail(ecode.OPERATION_ERROR, "AI扩图失败")
	}
	// 检查响应状态码
	if resp.StatusCode() != 200 {
		return nil, ecode.GetErrWithDetail(ecode.OPERATION_ERROR, "AI扩图失败")
	}

	// 解析响应
	var result model.CreateOutPaintingTaskResponse
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "响应解析失败")
	}
	if result.Code != "" {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "请求异常，AI扩图失败："+result.Message)
	}
	return &result, nil
}

// 获取执行扩图的任务的状态信息
func GetOutPaintingTaskResponse(taskId string) (*model.GetOutPaintingResponse, *ecode.ErrorWithCode) {
	if taskId == "" {
		return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "任务ID不能为空")
	}
	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", config.LoadConfig().AliYunAi.ApiKey)).
		Get(fmt.Sprintf(GetOutPaintingTaskURL, taskId))
	if err != nil {
		log.Println("请求失败:", err)
		return nil, ecode.GetErrWithDetail(ecode.OPERATION_ERROR, "获取任务状态失败")
	}
	// 检查响应状态码
	if resp.StatusCode() != 200 {
		return nil, ecode.GetErrWithDetail(ecode.OPERATION_ERROR, "获取任务状态失败")
	}
	// 解析响应
	var result model.GetOutPaintingResponse
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "响应解析失败")
	}
	return &result, nil
}
