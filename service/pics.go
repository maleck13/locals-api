package service

import (
	"github.com/maleck13/locals-api/Godeps/_workspace/src/github.com/disintegration/imaging"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"log"
	"mime/multipart"
	"os"
	"time"
)

/**
takes a multipart file. Does not close it. creates a thumbnail and returns the file path
*/
func ThumbnailMultipart(file multipart.File, fileName string) (string, error) {
	var (
		err      error
		img      image.Image
		thumbImg *image.NRGBA
	)
	if _, err := file.Seek(0, 0); err != nil {
		log.Printf("failed to seek to beginning of img " + err.Error())
		return "", err
	}
	var thumbPath string = "/tmp/" + time.Now().String() + fileName
	img, _, err = image.Decode(file)
	if nil != err {
		log.Printf("failed to decode img " + err.Error())
		return "", err
	}

	thumbImg = imaging.Thumbnail(img, 300, 300, imaging.Lanczos)

	out, err := os.Create(thumbPath)

	if nil != err {
		log.Printf("failed to create thumb path " + err.Error())
		return "", err
	}

	defer out.Close()

	// write new image to file
	err = jpeg.Encode(out, thumbImg, nil)

	return thumbPath, err
}

func Thumbnail(filePath string, fileName string) (string, error) {
	var (
		file     *os.File
		err      error
		img      image.Image
		thumbImg *image.NRGBA
	)

	file, err = os.Open(filePath)
	if nil != err {
		log.Printf("failed to decode img " + err.Error())
		return "", err
	}

	var thumbPath string = "/tmp/" + time.Now().String() + fileName

	img, _, err = image.Decode(file)
	if nil != err {
		log.Printf("failed to decode img " + err.Error())
		return "", err
	}

	thumbImg = imaging.Thumbnail(img, 300, 300, imaging.Lanczos)

	out, err := os.Create(thumbPath)

	if nil != err {
		log.Printf("failed to decode img " + err.Error())
		return "", err
	}

	defer out.Close()

	// write new image to file
	err = jpeg.Encode(out, thumbImg, nil)

	file.Close()
	os.Remove(filePath)

	return thumbPath, err
}

func SaveUploadedFile(file multipart.File, fileName string) (string, error) {
	tmpFile := "/tmp/" + fileName
	f, err := os.OpenFile(tmpFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println("error uploading pic " + err.Error())
		return "", err
	}
	_, err = io.Copy(f, file)
	if err != nil {
		log.Println("failed to copy file  " + err.Error())
		return "", err
	}
	file.Close()
	f.Close()
	return tmpFile, err
}
