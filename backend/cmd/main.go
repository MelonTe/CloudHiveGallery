package cmd

import (
	_ "chg/config"
	"chg/router"
	_ "fmt"

	"github.com/gin-gonic/gin"
)

func Main() {
	r := gin.Default()
	router.SetupRoutes(r)
	r.Run(":8080")
}
