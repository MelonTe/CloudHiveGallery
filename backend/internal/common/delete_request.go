package common

type DeleteRequest struct {
	Id uint64 `json:"id,string" binding:"required" comment:"删除的ID" swaggertype:"string"`
}
