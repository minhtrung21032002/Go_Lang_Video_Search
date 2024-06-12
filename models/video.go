package models

import (
	"context"
	"fmt"
	"net/http"

	"encoding/json"

	"example.com/mod/web_project/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Video struct {
	VideoName        string ` gorm:"column:video_name"`
	VideoS3          string ` gorm:"column:video_s3"`
	VideoDescription string ` gorm:"column:video_description"`
}

var DB *gorm.DB
var AWS *s3.Client
var rdb *redis.Client
var ctx = context.Background()

func init() {

	// Initialize the database connection
	var err error
	DB, err = config.Protgres()
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	AWS = config.Aws()
	rdb = config.RedisClient()
}

func GetAllVideo() ([]Video, error) {
	var videos []Video
	result := DB.Find(&videos)

	err := result.Error // Assign the error returned by DB.Find() to err

	// Check if the error is not nil
	if err != nil {
		// Handle the error
		// For example, log the error or return an empty slice
		fmt.Println("Error fetching videos:", err)
		return []Video{}, err
	} else {
		return videos, nil
	}
}

// Function to get a single user by ID
func GetVideoByID(videoID string) (Video, error) {
	var video Video
	if err := DB.First(&video, videoID).Error; err != nil {
		return Video{}, err
	}
	return video, nil
}

func (v Video) MarshalBinary() (data []byte, err error) {
	bytes, err := json.Marshal(v)
	return bytes, err
}

// Function to get a single user by ID
func UploadVideo(c *gin.Context) (Video, error) {

	// file, _ := c.FormFile("VideoS3")

	// // Open the uploaded file
	// uploadedFile, _ := file.Open()

	// defer uploadedFile.Close()

	// // Read the file content
	// fileContent, _ := io.ReadAll(uploadedFile)

	// // Print the file content
	// println("File Content:", string(fileContent))

	// // Respond to the request
	// c.String(200, "File content printed")

	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		fmt.Println("ERROR")
	}

	videoName := c.PostForm("VideoName")
	videoDescription := c.PostForm("VideoDescription")

	file, err := c.FormFile("VideoS3")
	if err != nil {
		fmt.Println(err)
	}

	uploader := manager.NewUploader(AWS)
	f, openErr := file.Open()

	if openErr != nil {
		fmt.Println("Open error")
	}

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("test-igolang"),
		Key:    aws.String(file.Filename),
		Body:   f,
		ACL:    "public-read",
	})
	if err != nil {
		fmt.Println(err)
	}

	video := Video{videoName, result.Location, videoDescription}
	videoBytes, err := video.MarshalBinary()
	if err := DB.Create(&video).Error; err != nil {
		return Video{}, err
	}
	videoID := fmt.Sprintf("video:%s", videoName)
	// videoMeta := map[string]interface{}{
	// 	"VideoName":        videoName,
	// 	"VideoDescription": videoDescription,
	// 	"VideoS3":          video.VideoS3,
	// }

	err = rdb.HSet(ctx, videoID, videoBytes, 0).Err()
	if err != nil {
		return Video{}, fmt.Errorf("error setting cache: %v", err)
	}
	return video, nil
}

func SearchVideo(c *gin.Context) error {
	query := c.Query("q")
	if query == "" {
		var videos []Video
		DB.Find(&videos)
		c.JSON(http.StatusOK, videos)
		return nil
	}

	keys, err := rdb.Keys(ctx, "video:*"+query+"*").Result()
	if err != nil {
		return err
	}
	fmt.Println(keys)

	keys2, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(keys2)

	videos := make([]Video, 0)
	for _, key := range keys {
		videoData, err := rdb.HGetAll(ctx, key).Result()
		fmt.Println(videoData)
		if err != nil {
			continue
		}
		fmt.Println(videoData["VideoName"])

		if len(videoData) > 0 {
			for jsonString := range videoData {
				var video Video
				// Unmarshal the JSON string into the Video struct
				err = json.Unmarshal([]byte(jsonString), &video)
				if err != nil {
					fmt.Printf("Error unmarshaling data for key %s: %v\n", key, err)
					continue
				}

				// Print the constructed Video struct
				fmt.Println("Constructed Video:", video)

				// Append the Video struct to the slice
				videos = append(videos, video)
			}
		}
		// video := Video{
		// 	VideoName:        videoData["VideoName"],
		// 	VideoS3:          videoData["VideoS3"],
		// 	VideoDescription: videoData["VideoDescription"],
		// }
		// fmt.Println(video)
		// videos = append(videos, video)
	}

	c.JSON(http.StatusOK, videos)
	return nil
}
