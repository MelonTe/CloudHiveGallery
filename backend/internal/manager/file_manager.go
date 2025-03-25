package manager

import (
	"chg/config"
	"chg/internal/ecode"
	"chg/internal/model/dto/file"
	"chg/pkg/tcos"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
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
	//存储到COS的Key
	uploadPath := fmt.Sprintf("%s/%s", uploadPrefix, uploadFileName)
	//3.解析结果并且返回
	//打开文件流
	src, _ := multipartFile.Open()
	defer src.Close()
	//调用压缩图片上传请求
	_, err := tcos.PutPictureWithCompress(src, uploadPath)
	if err != nil {
		log.Print(err)
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "上传失败")
	}
	//最终要获取的图片信息，为压缩后的图片信息。因此需要修改文件的后缀
	uploadPath = strings.Replace(uploadPath, fileType, "webp", 1)
	//缩略图路径
	thumbnailUrl := strings.Replace(uploadPath, ".webp", "_thumbnail."+fileType, 1)
	//获取图片信息结构体
	picInfo, err := tcos.GetPictureInfo(uploadPath)
	if err != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "获取图片信息失败")
	}
	return &file.UploadPictureResult{
		URL:          config.LoadConfig().Tcos.Host + "/" + uploadPath,
		ThumbnailURL: config.LoadConfig().Tcos.Host + "/" + thumbnailUrl,
		PicName:      fileNameNoType,
		PicSize:      picInfo.Size,
		PicWidth:     picInfo.Width,
		PicHeight:    picInfo.Height,
		PicScale:     math.Round(float64(picInfo.Width)/float64(picInfo.Height)*100) / 100,
		PicFormat:    picInfo.Format,
	}, nil
}

// 校验图片文件是否合法，不合法则返回原因
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

// 通过URL上传图片，获取图片的解析信息
// multipartFile: 上传的图片地址
// uploadPrefix: 上传文件的目录的前缀
// picName: 上传的图片的名称，若为空，则默认为“临时图片” + 一个随机种子。不为空，则只添加一个随机种子。
func UploadPictureByURL(fileURL string, uploadPrefix string, picName string) (*file.UploadPictureResult, *ecode.ErrorWithCode) {
	//1.校验文件

	//图片的暂时昵称用“临时图片”代替
	if picName == "" {
		picName = "临时图片"
	}
	//不论如何，都需要为picName添加一个随机种子，防止并发覆盖
	picName = picName + strconv.Itoa(rand.Intn(999999)+1)
	//传入picName，补充图片的后缀
	if err := ValidPictureByURL(fileURL, &picName); err != nil {
		return nil, err
	}
	//2.校验成功，将图片获取到本地tempfile中
	//localFilePath的格式为:tempfile/picName（包含了后缀）
	localFilePath, err := downLoadPictureByURL(fileURL, picName)
	if err != nil {
		return nil, err
	}
	//最终执行删除函数
	defer deleteTempFile(localFilePath)
	//3.图片上传地址
	//生成COSurl的随机种子
	u := uuid.New()
	hash := md5.Sum(u[:])
	id := hex.EncodeToString(hash[:])[:16]
	//文件后缀
	fileType := localFilePath[strings.LastIndex(localFilePath, ".")+1:]
	//上传到COS的文件名（id）
	uploadFileName := fmt.Sprintf("%s_%s.%s", time.Now().Format("2006-01-02"), id, fileType)
	//最终上传COS的key，包含了前缀（id）
	uploadPath := fmt.Sprintf("%s/%s", uploadPrefix, uploadFileName)
	//4.解析结果并且返回
	//打开文件流
	src, _ := os.Open(localFilePath)
	defer src.Close()
	_, errr := tcos.PutPictureWithCompress(src, uploadPath)
	if errr != nil {
		log.Print(errr)
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "上传失败")
	}
	//获取图片信息结构体，经过了压缩需要获取webp格式的图片信息
	uploadPath = strings.Replace(uploadPath, fileType, "webp", 1)
	//缩略图路径
	thumbnailUrl := strings.Replace(uploadPath, ".webp", "_thumbnail."+fileType, 1)
	picInfo, errr := tcos.GetPictureInfo(uploadPath)
	if errr != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "获取图片信息失败")
	}
	//picName去除后缀，作为图片昵称
	picNameNoType := picName[:strings.LastIndex(picName, ".")]
	return &file.UploadPictureResult{
		URL:          config.LoadConfig().Tcos.Host + "/" + uploadPath,
		ThumbnailURL: config.LoadConfig().Tcos.Host + "/" + thumbnailUrl,
		PicName:      picNameNoType,
		PicSize:      picInfo.Size,
		PicWidth:     picInfo.Width,
		PicHeight:    picInfo.Height,
		PicScale:     math.Round(float64(picInfo.Width)/float64(picInfo.Height)*100) / 100,
		PicFormat:    picInfo.Format,
	}, nil
}

// 下载图片到本地存储tempfile中，请确保地址已经校验过。
// 返回存储的本地路径，例如：tempfile/picName.jpg
func downLoadPictureByURL(fileURL string, picName string) (string, *ecode.ErrorWithCode) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return "", ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "下载图片失败")
	}
	defer resp.Body.Close()
	//创建临时文件
	tempDir := "tempfile"
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		return "", ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "创建临时文件夹失败")
	}
	//获取文件原始昵称，拼接存储路径
	tempFilePath := tempDir + "/" + picName
	//创建文件
	file, err := os.Create(tempFilePath)
	if err != nil {
		log.Println(err)
		return "", ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "链接格式不支持")
	}
	defer file.Close()

	//写入文件
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "写入文件失败")
	}
	//写入成功
	return tempFilePath, nil
}

// 执行本地临时文件删除
func deleteTempFile(tempFilePath string) {
	if tempFilePath != "" {
		os.Remove(tempFilePath)
	}
}

// 通过URL的Head信息校验图片是否合法，若合法则补充图片的后缀，不合法返回原因。
func ValidPictureByURL(fileURL string, picName *string) *ecode.ErrorWithCode {
	//1.校验链接是否为空
	if fileURL == "" {
		return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "URL为空")
	}
	//2.校验URL格式
	_, err := url.ParseRequestURI(fileURL)
	if err != nil {
		return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "URL格式错误")
	}
	//3.校验URL的协议
	if !strings.HasPrefix(fileURL, "http") || !strings.HasPrefix(fileURL, "https") {
		return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "仅支持 HTTP 或 HTTPS 协议的文件地址")
	}
	//发送HEAD请求验证文件是否存在
	resp, err := http.Head(fileURL)
	if err != nil {
		//未正常返回，无需进行其他判断。
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil
	}
	//4.文件存在，文件类型校验
	contentType := resp.Header.Get("Content-Type")
	//不为空，才校验是否合法
	if contentType != "" {
		//校验文件类型
		allowType := []string{"image/jpeg", "image/jpg", "image/png", "image/webp"}
		isAllow := false
		for _, v := range allowType {
			if contentType == v {
				*picName = *picName + "." + strings.Split(contentType, "/")[1]
				isAllow = true
				break
			}
		}
		if !isAllow {
			return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "文件类型不支持")
		}
	}
	//5.文件大小校验
	contentLength := resp.Header.Get("Content-Length")
	//不为空校验
	if contentLength != "" {
		size, err := strconv.ParseUint(contentLength, 10, 64)
		if err != nil {
			return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "文件大小格式异常")
		}
		ONE_M := uint64(1024 * 1024)
		if size > 2*ONE_M {
			return ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "文件过大，不能超过2MB")
		}
	}
	return nil
}
