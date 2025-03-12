package common

type DeleteRequest struct {
	Id uint64 `json:"id" binding:"required" comment:"删除的ID"`
}
