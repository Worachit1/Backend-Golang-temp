package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

type user struct {
	ID int64 `json:"id"`
}

func GetUserByToken(ctx *gin.Context) (int64, error) {
	claims, exist := ctx.Get("claims")
	if !exist {
		return 0, nil
	}
	var user user
	err := json.Unmarshal(claims.([]byte), &user)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

// คืน Cloudinary instance พร้อมใช้งาน
func NewCloudinary() (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return nil, fmt.Errorf("cloudinary config error: %v", err)
	}
	return cld, nil
}

// อัปโหลดรูปภาพ พร้อมคืน URL
func UploadImageToCloudinary(file interface{}) (string, error) {
	cld, err := NewCloudinary()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:   "user",
		PublicID: fmt.Sprintf("user_%d", time.Now().UnixNano()),
	})
	if err != nil {
		return "", fmt.Errorf("upload failed: %v", err)
	}

	return uploadResult.SecureURL, nil
}

func UploadToCloudinary(fileHeader *multipart.FileHeader, folder string) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return "", err
	}

	uploadParams := uploader.UploadParams{
		Folder: folder,
	}
	uploadResult, err := cld.Upload.Upload(nil, file, uploadParams)
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}
