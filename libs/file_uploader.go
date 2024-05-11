package libs

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func UploadFile(url string, filename string, filePath string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	ext := filepath.Ext(filePath)

	localFilePath := "files/" + filename + ext

	localFile, err := os.Create(localFilePath)
	if err != nil {
		return "", err
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, response.Body)
	if err != nil {
		return "", err
	}

	return localFilePath, nil
}

func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}
