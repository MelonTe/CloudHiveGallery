package controller

import (
	"chg/internal/common"
	"chg/internal/ecode"
	"chg/internal/manager"
	"chg/pkg/tcos"
	"fmt"
	"io"
	"log"
	"path"

	"github.com/gin-gonic/gin"
)

// TestUploadFile godoc
// @Summary      测试文件上传接口「管理员」
// @Tags         file
// @Accept       mpfd
// @Produce      json
// @Param        file formData file true "文件"
// @Success      200  {object}  common.Response{data=string} "响应文件存储在COS的KEY"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /v1/file/test/upload [POST]
func TestUploadFile(c *gin.Context) {
	file, _ := c.FormFile("file")
	manager.UploadPicture(file, "test")
	// 处理文件逻辑，采用不落地存储的形式
	//打开文件流
	src, err := file.Open()
	if err != nil {
		common.BaseResponse(c, nil, "文件打开失败", ecode.PARAMS_ERROR)
		return
	}
	defer src.Close()
	//上传到COS
	key := fmt.Sprintf("test/%s", file.Filename)
	err = tcos.PutObject(src, key)
	if err != nil {
		log.Print(err)
		common.BaseResponse(c, nil, "上传失败", ecode.SYSTEM_ERROR)
		return
	}
	common.Success(c, "上传成功")
}

// TestDownloadFile godoc
// @Summary      测试文件下载接口「管理员」
// @Tags         file
// @Produce      octet-stream
// @Param        key query string true "文件存储在 COS 的 KEY"
// @Success      200 {file} file "返回文件流"
// @Failure      400 {object} common.Response "下载失败，详情见响应中的 code"
// @Router       /v1/file/test/download [GET]
func TestDownloadFile(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		common.BaseResponse(c, nil, "缺少 key 参数", ecode.PARAMS_ERROR)
		return
	}

	// 从 COS 获取文件流
	reader, err := tcos.GetObject(key)
	if err != nil {
		log.Printf("文件下载失败: %v", err)
		common.BaseResponse(c, nil, "文件下载失败", ecode.SYSTEM_ERROR)
		return
	}
	defer reader.Close() //关闭文件流，发送TCP FIN包

	// 设置 HTTP 头，返回流式数据
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", path.Base(key)))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Transfer-Encoding", "chunked")

	// 通过 io.Copy() 流式传输数据
	if _, err := io.Copy(c.Writer, reader); err != nil {
		log.Printf("流式传输失败: %v", err)
	}
}
