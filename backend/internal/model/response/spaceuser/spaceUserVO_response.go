package spaceuser

import (
	"chg/internal/model/entity"
	resSpace "chg/internal/model/response/space"
	resUser "chg/internal/model/response/user"
	"time"
)

type SpaceUserVO struct {
	ID         uint64           `json:"id,string" swaggertype:"string"`
	SpaceID    uint64           `json:"spaceId,string" swaggertype:"string"`
	UserID     uint64           `json:"userId,string" swaggertype:"string"`
	SpaceRole  string           `json:"spaceRole"`
	CreateTime time.Time        `json:"createTime"`
	UpdateTime time.Time        `json:"updateTime"`
	User       resUser.UserVO   `json:"user"`  // 用户信息
	SpaceVO    resSpace.SpaceVO `json:"space"` // 空间信息
}

func VOToEntity(vo SpaceUserVO) entity.SpaceUser {
	return entity.SpaceUser{
		ID:        vo.ID,
		SpaceID:   vo.SpaceID,
		UserID:    vo.UserID,
		SpaceRole: vo.SpaceRole,
	}
}

func EntityToVO(entity entity.SpaceUser, userVO resUser.UserVO, spaceVO resSpace.SpaceVO) SpaceUserVO {
	return SpaceUserVO{
		ID:         entity.ID,
		SpaceID:    entity.SpaceID,
		UserID:     entity.UserID,
		SpaceRole:  entity.SpaceRole,
		CreateTime: entity.CreateTime,
		UpdateTime: entity.UpdateTime,
		User:       userVO,
		SpaceVO:    spaceVO,
	}
}
