package space

import (
	"chg/internal/common"
	"chg/internal/model/entity"
)

type ListSpaceResponse struct {
	common.PageResponse
	Records []entity.Space `json:"records"`
}
