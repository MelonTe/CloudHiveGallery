package user

import "chg/internal/common"

type ListUserVOResponse struct {
	common.PageResponse
	Records []UserVO `json:"records" `
}
