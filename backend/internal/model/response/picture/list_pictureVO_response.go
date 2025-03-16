package picture

import "chg/internal/common"

type ListPictureVOResponse struct {
	common.PageResponse
	Records []PictureVO `json:"records" `
}
