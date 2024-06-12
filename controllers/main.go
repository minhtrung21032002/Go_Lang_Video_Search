package controllers

import (
	"fmt"
	"net/http"

	"example.com/mod/web_project/models"
	"github.com/gin-gonic/gin"
)

type Video struct {
	VideoID          int32  ` gorm:"column:id"`
	VideoName        string ` gorm:"column:video_name"`
	VideoS3          string ` gorm:"column:video_s3"`
	VideoDescription string ` gorm:"column:video_description"`
}

// Handler for getting user list
func GetVideoList(c *gin.Context) {
	videos, error := models.GetAllVideo()

	if error != nil {
		fmt.Println("GET ALL VIDEO WRONG")

	}

	fmt.Println(videos)
	c.HTML(http.StatusOK, "main.html", gin.H{
		"videos": videos,
	})
}

func GetForm(c *gin.Context) {

	c.HTML(http.StatusOK, "submit.html", gin.H{})
}

// Handler for getting a single user
func GetVideo(c *gin.Context) {
	userID := c.Param("id")
	user, error := models.GetVideoByID(userID)
	if error != nil {
		fmt.Println("GET VIDEO BY ID WRONG")
	}
	c.JSON(http.StatusOK, user)
}

func UploadVideo(c *gin.Context) {

	message, error := models.UploadVideo(c)
	if error != nil {
		fmt.Println(error)
	}
	c.JSON(http.StatusOK, message)
}

func SearchVideos(c *gin.Context) {
	err := models.SearchVideo(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error searching videos"})
		return
	}
}
