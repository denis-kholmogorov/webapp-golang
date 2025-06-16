package service

import (
	"context"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"web/application/dto"
)

type StorageService struct {
	cloud *cloudinary.Cloudinary
}

func NewStorageService() *StorageService {
	cld, err := cloudinary.NewFromURL(getUrlCloudinary())
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &StorageService{
		cld,
	}
}

func (s StorageService) Upload(c *gin.Context) {
	log.Println("StorageService:upload()")
	form, err := c.MultipartForm()
	if err != nil {
		fmt.Println(err)
	}

	files := form.File["file"]

	var upload *uploader.UploadResult
	for _, file := range files {
		log.Println(file.Filename)
		open, err := file.Open()
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("StorageService:Upload() Cannot open file with"))
			return
		}
		upload, err = s.cloud.Upload.Upload(context.Background(), open, uploader.UploadParams{})
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("StorageService:Upload() Cannot upload to cloud file"))
			return
		}

	}
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("GetCountry() Row with not found"))
	} else {
		c.JSON(http.StatusOK, dto.PhotoDto{FileName: upload.URL})
	}

}

func getUrlCloudinary() string {
	url, exists := os.LookupEnv("CLOUDINARY_URL")
	if !exists {
		url = "localhost:9080"
	}
	return url
}
