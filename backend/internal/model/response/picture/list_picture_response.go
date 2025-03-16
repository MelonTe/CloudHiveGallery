package picture

import (
	"chg/internal/common"
	"chg/internal/model/entity"
)

type ListPictureResponse struct {
	common.PageResponse
	Records []entity.Picture `json:"records" `
}
