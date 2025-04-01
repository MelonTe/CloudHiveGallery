package fetcher

//该文件用于实现调用百度搜图API的功能
import (
	"chg/internal/api/imagesearch/model"
	"chg/internal/ecode"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"regexp"
	"resty.dev/v3"
	"strings"
	"time"
)

// 定义返回数据的结构体
type BaiduResponse struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   struct {
		URL  string `json:"url"`
		Sign string `json:"sign"`
	} `json:"data"`
}

// 调用百度以图搜图接口获取原始搜索结果URL
// 上传原始图片地址
// 仅支持20M以下jpg，jpeg，png，bmp，gif等格式的图片
// 请传入webp图片，会添加解析成为PNG格式
func GetImagePageURL(imageURL string) (string, *ecode.ErrorWithCode) {
	//1.准备请求参数
	//处理图片地址，添加格式转换
	imageURL = imageURL + "?imageMogr2/format/png"
	//需要进行转义
	imageURL = url.QueryEscape(imageURL)
	formData := map[string]string{}
	formData["image"] = imageURL //图片地址
	formData["tn"] = "pc"
	formData["from"] = "pc" //图片类型
	formData["image_source"] = "PC_UPLOAD_URL"
	uptime := fmt.Sprintf("%d", time.Now().UnixMilli())
	//请求URL
	reqUrl := "https://graph.baidu.com/upload?uptime=" + uptime
	//2.发送POST请求至API
	client := resty.New()
	resp, err := client.R().SetHeader("Acs-Token", fmt.Sprintf("%d", rand.IntN(1000))).SetFormData(formData).SetTimeout(5 * time.Second).Post(reqUrl)
	if err != nil || resp.StatusCode() != http.StatusOK {
		return "", ecode.GetErrWithDetail(ecode.OPERATION_ERROR, "请求搜图接口失败")
	}
	//3.处理响应结果
	var baiduResp BaiduResponse
	if err := json.Unmarshal(resp.Bytes(), &baiduResp); err != nil {
		return "", ecode.GetErrWithDetail(ecode.OPERATION_ERROR, "解析响应结果失败")
	}
	if baiduResp.Status != 0 || baiduResp.Data.URL == "" {
		return "", ecode.GetErrWithDetail(ecode.OPERATION_ERROR, "获取图片页面失败，可能是图片格式不支持")
	}
	return baiduResp.Data.URL, nil
}

// 通过GetImagePageURL获取的URL，来获取简略图片信息的请求接口，即FirstURL
func GetImageFirstURL(searchResultURL string) (string, *ecode.ErrorWithCode) {
	// 发送 HTTP GET 请求
	resp, err := http.Get(searchResultURL)
	if err != nil {
		return "", ecode.GetErrWithDetail(ecode.OPERATION_ERROR, "请求图片页面失败")
	}
	defer resp.Body.Close()

	// 解析 HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", ecode.GetErrWithDetail(ecode.OPERATION_ERROR, "解析图片页面失败")
	}

	// 查找所有 <script> 标签
	scriptElements := doc.Find("script")
	reg := regexp.MustCompile(`"firstUrl"\s*:\s*"(.*?)"`)

	var firstURL string
	scriptElements.EachWithBreak(func(i int, s *goquery.Selection) bool {
		scriptContent := s.Text()
		matches := reg.FindStringSubmatch(scriptContent)
		if len(matches) > 1 {
			firstURL = strings.ReplaceAll(matches[1], "\\/", "/")
			return false // 找到后退出循环
		}
		return true
	})

	if firstURL == "" {
		return "", ecode.GetErrWithDetail(ecode.OPERATION_ERROR, "搜索失败")
	}
	return firstURL, nil
}

// 第三部，根据FIRSTURL，获取图片列表
func GetImageList(firstURL string) ([]model.ImageSearchResult, *ecode.ErrorWithCode) {
	resp, err := http.Get(firstURL)
	if err != nil {
		return nil, ecode.GetErrWithDetail(ecode.OPERATION_ERROR, "请求图片列表失败")
	}
	defer resp.Body.Close()
	//读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ecode.GetErrWithDetail(ecode.OPERATION_ERROR, "读取响应数据失败")
	}
	// 解析 JSON
	var apiResp model.APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, ecode.GetErrWithDetail(ecode.OPERATION_ERROR, "解析JSON数据失败")
	}

	return apiResp.Data.List, nil
}
