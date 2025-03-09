package cmd

import (
	_ "chg/config"
	"chg/pkg/db"
	"chg/router"
	_ "fmt"

	"github.com/gin-gonic/gin"
)

func Main() {
	r := gin.Default()
	router.SetupRoutes(r)
	db.LoadDB()
	r.Run(":8080")
}
