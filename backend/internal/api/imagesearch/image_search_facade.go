package imagesearch

import (
	"chg/internal/api/imagesearch/fetcher"
	"chg/internal/api/imagesearch/model"
	"chg/internal/ecode"
)

//门面模式，通过统一接口简化多个接口的调用
//该文件整合封装了百度搜图API的功能

func SearchImage(imageURL string) ([]model.ImageSearchResult, *ecode.ErrorWithCode) {
	//1.获取图片页面URL
	imagePageURL, err := fetcher.GetImagePageURL(imageURL)
	if err != nil {
		return nil, err
	}
	firstUrl, err := fetcher.GetImageFirstURL(imagePageURL)
	if err != nil {
		return nil, err
	}
	//2.获取图片搜索结果
	imageSearchResult, err := fetcher.GetImageList(firstUrl)
	if err != nil {
		return nil, err
	}
	return imageSearchResult, nil
}
