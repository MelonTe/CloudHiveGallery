package service

import (
	"chg/internal/ecode"
	"chg/internal/manager"
	"chg/internal/model/entity"
	reqPicture "chg/internal/model/request/picture"
	resPicture "chg/internal/model/response/picture"
	resUser "chg/internal/model/response/user"
	"chg/internal/repository"
	"fmt"
	"mime/multipart"
	"time"
)

type PictureService struct {
	PictureRepo *repository.PictureRepository
}

func NewPictureService() *PictureService {
	return &PictureService{
		PictureRepo: repository.NewPictureRepository(),
	}
}

// 该服务用于修改或插入图片数据到服务器中
func (s *PictureService) UploadPicture(multipartFile *multipart.FileHeader, PictureUploadRequest *reqPicture.PictureUploadRequest, loginUser *entity.User) (*resPicture.PictureVO, *ecode.ErrorWithCode) {
	//判断图片是需要新增还是需要更新
	picId := uint64(0)
	if PictureUploadRequest.ID != 0 {
		picId = PictureUploadRequest.ID
	}
	//若更新图片，则需要校验图片是否存在
	if picId != 0 {
		_, err := s.PictureRepo.FindById(picId)
		if err != nil {
			return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "图片不存在")
		}
	}
	//上传图片，得到信息
	uploadPathPrefix := fmt.Sprintf("public/%d", loginUser.ID)
	info, err := manager.UploadPicture(multipartFile, uploadPathPrefix)
	if err != nil {
		return nil, err
	}
	//构造插入数据库的实体
	pic := &entity.Picture{
		URL:       info.URL,
		Name:      info.PicName,
		PicSize:   info.PicSize,
		PicWidth:  info.PicWidth,
		PicHeight: info.PicHeight,
		PicScale:  info.PicScale,
		PicFormat: info.PicFormat,
		UserID:    loginUser.ID,
		EditTime:  time.Now(),
	}
	//若是更新，则需要更新ID
	if picId != 0 {
		pic.ID = picId
	}
	//进行插入或者更新操作，即save
	errr := s.PictureRepo.SavePicture(pic)
	if errr != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库错误")
	}
	userVO := resUser.GetUserVO(*loginUser)
	picVO := resPicture.EntityToVO(*pic, userVO)
	return &picVO, nil
}
