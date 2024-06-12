package routes

import (
	"example.com/mod/web_project/controllers"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes(r *gin.Engine) {

	videoRoutes := r.Group("/main")
	{
		videoRoutes.GET("/", controllers.GetVideoList)
		videoRoutes.GET("/:id", controllers.GetVideo)
		videoRoutes.POST("/form", controllers.UploadVideo)
		videoRoutes.GET("/form", controllers.GetForm)
		videoRoutes.GET("/search", controllers.SearchVideos)
		// Add more routes as needed
	}
}
