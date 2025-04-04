package spaceuser

//查询空间成员请求，不需要分页
type SpaceUserQueryRequest struct {
	ID        uint64 `json:"Id,string" swaggertype:"string"`      //表的元组ID
	SpaceID   uint64 `json:"spaceId,string" swaggertype:"string"` //空间ID
	UserID    uint64 `json:"userId,string" swaggertype:"string"`  //用户ID
	SpaceRole string `json:"spaceRole"`                           //空间角色：viewer-查看者 editor-编辑者 admin-管理员
}
