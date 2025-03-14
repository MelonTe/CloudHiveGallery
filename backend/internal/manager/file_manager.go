package manager

import (
	"chg/config"
	"chg/internal/ecode"
	"chg/internal/model/dto/file"
	"chg/pkg/tcos"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"mime/multipart"
	"strings"
	"time"

	"github.com/google/uuid"
)

//该文件用于处理文件上传服务，不与具体业务逻辑耦合

// 上传图片，获取图片的解析信息
// multipartFile: 上传的文件
// uploadPrefix: 上传文件的目录的前缀
func UploadPicture(multipartFile *multipart.FileHeader, uploadPrefix string) (*file.UploadPictureResult, *ecode.ErrorWithCode) {
	//1.校验文件
	if err := ValidPicture(multipartFile); err != nil {
		return nil, err
	}
	//2.图片上传地址
	//生成url的随机种子
	u := uuid.New()
	hash := md5.Sum(u[:])
	id := hex.EncodeToString(hash[:])[:16]
	//文件后缀
	fileType := multipartFile.Filename[strings.LastIndex(multipartFile.Filename, ".")+1:]
	//文件名
	uploadFileName := fmt.Sprintf("%s_%s.%s", time.Now().Format("2006-01-02"), id, fileType)
	fileNameNoType := uploadFileName[:strings.LastIndex(uploadFileName, ".")]
	//最终文件名
	uploadPath := fmt.Sprintf("%s/%s", uploadPrefix, uploadFileName)
	//3.解析结果并且返回
	//打开文件流
	src, _ := multipartFile.Open()
	defer src.Close()
	_, err := tcos.PutPicture(src, uploadPath)
	if err != nil {
		log.Print(err)
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "上传失败")
	}
	//获取图片信息结构体
	picInfo, err := tcos.GetPictureInfo(uploadPath)
	if err != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "获取图片信息失败")
	}
	return &file.UploadPictureResult{
		URL:       config.LoadConfig().Tcos.Host + "/" + uploadPath,
		PicName:   fileNameNoType,
		PicSize:   picInfo.Size,
		PicWidth:  picInfo.Width,
		PicHeight: picInfo.Height,
		PicScale:  math.Round(float64(picInfo.Width)/float64(picInfo.Height)*100) / 100,
		PicFormat: picInfo.Format,
	}, nil
}

func ValidPicture(multipartFile *multipart.FileHeader) *ecode.ErrorWithCode {
	if multipartFile == nil {
		return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "文件为空")
	}
	//校验文件大小
	fileSize := multipartFile.Size
	ONE_MB := int64(1024 * 1024)
	if fileSize > 2*ONE_MB {
		return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "文件过大，不能超过2MB")
	}
	//校验文件类型
	lastDotIndex := strings.LastIndex(multipartFile.Filename, ".")
	if lastDotIndex == -1 {
		return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "文件不是图片")
	}
	fileType := multipartFile.Filename[lastDotIndex:]
	//允许的文件类型
	allowType := []string{".jpg", ".jpeg", ".png", ".webp"}
	isAllow := false
	for _, v := range allowType {
		if fileType == v {
			isAllow = true
			break
		}
	}
	if !isAllow {
		return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "文件类型不支持")
	}
	return nil
}
