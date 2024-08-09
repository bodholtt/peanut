package files

import (
	"io"
	"mime/multipart"
	"os"
	"peanutserver/pcfg"
)

// UploadImage - handle uploading an image to local storage
func UploadImage(fileData multipart.File, filename string) error {

	localFile, err := os.Create(pcfg.Cfg.Images.RootLocation + "images/" + filename)
	if err != nil {
		return err
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, fileData)
	if err != nil {
		return err
	}
	return nil
}

// DeleteImage - handle deletion of an image from local storage
func DeleteImage(filepath string) error {

	err := os.Remove(pcfg.Cfg.Images.RootLocation + "images/" + filepath)
	if err != nil {
		return err
	}
	return nil
}
