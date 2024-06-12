package main

import (
	"example.com/mod/web_project/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.LoadHTMLGlob("view/*.html")
	routes.InitializeRoutes(r)
	r.Run(":3000")
}
