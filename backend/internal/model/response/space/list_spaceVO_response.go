package space

import (
	"chg/internal/common"
)

type ListSpaceVOResponse struct {
	common.PageResponse
	Records []SpaceVO `json:"records"`
}
