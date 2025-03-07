package controller

import (
	"shg/internal/common"
	"shg/internal/ecode"

	"github.com/gin-gonic/gin"
)

// HelloWorld godoc
// @Summary      Say Hello World
// @Description  This is a simple API endpoint that returns a "Hello, World!" message.
// @Tags         example
// @Accept       json
// @Produce      json
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Router       /hello [get]
func HelloWorld(ctx *gin.Context) {
	common.Error(ctx, ecode.FORBIDDEN_ERROR)
}
