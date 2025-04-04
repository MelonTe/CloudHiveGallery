package space

import (
	"chg/internal/model/entity"
	resUser "chg/internal/model/response/user"
	"time"
)

type SpaceVO struct {
	ID         uint64         `json:"id,string" swaggertype:"string"` // Space ID
	SpaceName  string         `json:"spaceName"`
	SpaceLevel int            `json:"spaceLevel"`
	MaxSize    int64          `json:"maxSize"`
	MaxCount   int64          `json:"maxCount"`
	TotalSize  int64          `json:"totalSize"`
	TotalCount int64          `json:"totalCount"`
	UserID     uint64         `json:"userId,string" swaggertype:"string"` // User ID
	CreateTime time.Time      `json:"createTime"`
	EditTime   time.Time      `json:"editTime"`
	UpdateTime time.Time      `json:"updateTime"`
	User       resUser.UserVO `json:"user"`
	SpaceType  int            `json:"spaceType"` // Space type: 0 - 私人空间, 1 - 团队空间
}

// Convert SpaceVO to entity.Space
func VOToEntity(vo SpaceVO) entity.Space {
	return entity.Space{
		ID:         vo.ID,
		SpaceName:  vo.SpaceName,
		SpaceLevel: vo.SpaceLevel,
		MaxSize:    vo.MaxSize,
		MaxCount:   vo.MaxCount,
		TotalSize:  vo.TotalSize,
		TotalCount: vo.TotalCount,
		UserID:     vo.UserID,
		CreateTime: vo.CreateTime,
		EditTime:   vo.EditTime,
		UpdateTime: vo.UpdateTime,
		SpaceType:  vo.SpaceType,
	}
}

// Convert entity.Space to SpaceVO
func EntityToVO(entity entity.Space, userVO resUser.UserVO) SpaceVO {
	return SpaceVO{
		ID:         entity.ID,
		SpaceName:  entity.SpaceName,
		SpaceLevel: entity.SpaceLevel,
		MaxSize:    entity.MaxSize,
		MaxCount:   entity.MaxCount,
		TotalSize:  entity.TotalSize,
		TotalCount: entity.TotalCount,
		UserID:     entity.UserID,
		CreateTime: entity.CreateTime,
		EditTime:   entity.EditTime,
		UpdateTime: entity.UpdateTime,
		User:       userVO,
		SpaceType:  entity.SpaceType,
	}
}
