package services

import (
	"context"
	"gin-mongo-api/configs"
	"gin-mongo-api/utils"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
)

func UploadFile(bucketName string, extention string, file *multipart.FileHeader, customFileName ...string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Make a new bucket called {{bucketName}}.
	location := "us-east-1"
	err := configs.MinioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		_, _ = configs.MinioClient.BucketExists(ctx, bucketName)
		// exists, errBucketExists := configs.MinioClient.BucketExists(ctx, bucketName)
		// if errBucketExists == nil && exists {
		// 	log.Printf("We already own %s\n", bucketName)
		// } else {
		// 	log.Println(err)
		// }
	}
	// else {
	// log.Printf("Successfully created %s\n", bucketName)
	// }

	// Copy file in local and Upload the file
	var fileExtention string = "png"
	if file.Header.Get("Content-Type") != "" {
		fileExtention = strings.Split(file.Header["Content-Type"][0], "/")[1]
	}
	fileName := utils.GetMD5Hash(file.Filename) + "." + fileExtention
	if len(customFileName) > 0 {
		fileName = customFileName[0] + "." + fileExtention
	}
	filePath := "public/" + fileName
	_ = SaveUploadedFile(file, filePath)
	fileNameInStorage := utils.GetMD5Hash(file.Filename) + extention
	// log.Println(err)

	var contentType string = "image/png"
	if file.Header.Get("Content-Type") != "" {
		contentType = file.Header["Content-Type"][0]
	}

	// Upload the zip file with FPutObject
	_, _ = configs.MinioClient.FPutObject(ctx, bucketName, fileNameInStorage, filePath, minio.PutObjectOptions{ContentType: contentType})
	// if err != nil {
	// 	log.Println(err)
	// }

	// log.Printf("Successfully uploaded %s of size %d\n", fileName, info.Size)

	// Delete the local file
	_ = DeleteFile(fileName)
	// if deleteErr != nil {
	// 	log.Println(deleteErr)
	// }

	link := configs.MinioClient.EndpointURL().String() + "/" + bucketName + "/" + fileNameInStorage
	return link
}

func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func DeleteFile(fileName string) error {
	filePath := "public/" + fileName
	err := os.Remove(filePath)
	return err
}
