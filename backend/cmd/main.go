package cmd

import (
	_ "fmt"
	_ "shg/config"
	"shg/router"

	"github.com/gin-gonic/gin"
)

func Main() {
	r := gin.Default()
	router.SetupRoutes(r)
	r.Run(":8080")
}
