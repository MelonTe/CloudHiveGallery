package tcos

import (
	"chg/config"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// 图片详细数据结构体
type PicInfo struct {
	Format        string `json:"format"`
	Width         int    `json:"width,string"`
	Height        int    `json:"height,string"`
	Size          int64  `json:"size,string"`
	MD5           string `json:"md5"`
	FrameCount    int    `json:"frame_count,string"`
	BitDepth      int    `json:"bit_depth,string"`
	VerticalDPI   int    `json:"vertical_dpi,string"`
	HorizontalDPI int    `json:"horizontal_dpi,string"`
}

var tcos *cos.Client

func init() {
	// 将 examplebucket-1250000000 和 COS_REGION 修改为用户真实的信息
	// 存储桶名称，由 bucketname-appid 组成，appid 必须填入，可以在 COS 控制台查看存储桶名称。https://console.cloud.tencent.com/cos5/bucket
	// COS_REGION 可以在控制台查看，https://console.cloud.tencent.com/cos5/bucket, 关于地域的详情见 https://cloud.tencent.com/document/product/436/6224
	c := config.LoadConfig()
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", c.Tcos.BucketName, c.Tcos.Region))
	// 用于 Get Service 查询，默认全地域 service.cos.myqcloud.com
	su, _ := url.Parse(fmt.Sprintf("https://cos.%s.myqcloud.com", c.Tcos.Region))
	b := &cos.BaseURL{BucketURL: u, ServiceURL: su}
	log.Println("密钥显示：", os.Getenv("SECRETID"))
	// 1.永久密钥
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("SECRETID"),  // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
			SecretKey: os.Getenv("SECRETKEY"), // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
		},
	})
	tcos = client
	//测试连接，失败则panic
	s, _, err := client.Service.Get(context.Background())
	if err != nil {
		panic(err)
	}

	for _, b := range s.Buckets {
		fmt.Printf("%#v\n", b)
	}
}
func LoadDB() *cos.Client {
	return tcos
}

// 上传本地对象到COS服务器中，key是对象在存储桶中的唯一标识，例如"doc/test.txt"，path是本地文件路径。
func PutObjectFromLocal(key, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	opt := &cos.ObjectPutOptions{}

	_, err = tcos.Object.Put(context.Background(), key, f, opt)
	return err
}

// 上传实现了io.Reader接口的数据，key是对象在存储桶中的唯一标识，例如"doc/test.txt"。
func PutObject(f io.Reader, key string) error {
	opt := &cos.ObjectPutOptions{}
	_, err := tcos.Object.Put(context.Background(), key, f, opt)
	if err != nil {
		return err
	}
	return nil
}

// 从 COS 获取文件流（流式传输）
func GetObject(key string) (io.ReadCloser, error) {
	resp, err := tcos.Object.Get(context.Background(), key, nil)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// 上传图片对象，返回原始响应体
// key是对象在存储桶中的唯一标识，例如"doc/test.jpg"。
func PutPicture(f io.Reader, key string) (*cos.Response, error) {
	pic := &cos.PicOperations{
		IsPicInfo: 1,
	}
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			XOptionHeader: &http.Header{},
		},
	}
	opt.XOptionHeader.Add("Pic-Operations", cos.EncodePicOperations(pic))
	opt.XOptionHeader.Add("x-cos-return-response", "true")
	res, err := tcos.Object.Put(context.Background(), key, f, opt)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 上传图片对象并且进行上传时压缩，保存成webp格式，返回原始响应体。
// 原图不会被覆盖。
// key是对象在存储桶中的唯一标识，例如"doc/test.png"。
// 同时，会添加一个缩略图，会缩略原图至宽高至多为256，参考网址为doc/test_thumbnail.png。
func PutPictureWithCompress(f io.Reader, key string) (*cos.Response, error) {
	//取出key的后缀，修改为webp
	lastIdx := strings.LastIndex(key, ".")
	var newKey string
	var thumbnailKey string
	//确保安全性
	if lastIdx != -1 {
		keyNoType := key[:lastIdx]
		keyType := key[lastIdx:]
		newKey = keyNoType + ".webp"
		thumbnailKey = keyNoType + "_thumbnail" + keyType
	}
	pic := &cos.PicOperations{
		IsPicInfo: 1,
		Rules: []cos.PicOperationsRules{
			{
				Rule:   "imageMogr2/format/webp",
				FileId: "/" + newKey,
			},
			{
				Rule:   fmt.Sprintf("imageMogr2/thumbnail/%dx%d>", 256, 256),
				FileId: "/" + thumbnailKey,
			},
		},
	}
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			XOptionHeader: &http.Header{},
		},
	}
	opt.XOptionHeader.Add("Pic-Operations", cos.EncodePicOperations(pic))
	res, err := tcos.Object.Put(context.Background(), key, f, opt)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 获取图片详细信息，返回详细信息的结构体
// key是对象在存储桶中的唯一标识，例如"doc/test.jpg"。
func GetPictureInfo(key string) (*PicInfo, error) {
	operation := "imageInfo"
	resp, err := tcos.CI.Get(context.Background(), key, operation, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	info, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var picInfo PicInfo
	err = json.Unmarshal(info, &picInfo)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &picInfo, nil
}

// 获取图片的主色调，返回十六进制的图片主色调，例如：0x736246
func GetPictureColor(key string) (string, error) {
	operation := "imageAve"
	resp, err := tcos.CI.Get(context.Background(), key, operation, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	//获取响应体
	var result map[string]string
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	//解析JSON
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	//获取RGB值
	rgb := result["RGB"]
	if rgb == "" {
		return "", fmt.Errorf("获取图片主色调失败")
	}
	return rgb, nil
}

// 删除对象，key为唯一标识
func DeleteObject(key string) error {
	_, err := tcos.Object.Delete(context.Background(), key)
	if err != nil {
		return err
	}
	return nil
}
