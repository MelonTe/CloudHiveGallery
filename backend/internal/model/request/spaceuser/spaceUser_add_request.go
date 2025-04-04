package spaceuser

type SpaceUserAddRequest struct {
	SpaceID   uint64 `json:"spaceId,string" swaggertype:"string"` //空间ID
	UserID    uint64 `json:"userId,string" swaggertype:"string"`  //用户ID
	SpaceRole string `json:"spaceRole"`                           //空间角色：viewer-查看者 editor-编辑者 admin-管理员
}
