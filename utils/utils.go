package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"mime/multipart"
	"os"
)

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func GetSHA256Hash(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func GetFileHeader(file *os.File) (*multipart.FileHeader, error) {
	// get file size
	fileStat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// create *multipart.FileHeader
	return &multipart.FileHeader{
		Filename: fileStat.Name(),
		Size:     fileStat.Size(),
	}, nil
}
