package utils

import (
	"errors"
	"net/http"
	"os"
)

func ImageSave(img []byte, fileName string) (err error) {
	if len(fileName) == 0 || img == nil {
		return errors.New("fileName and img must not empty")
	}
	Mkdir(fileName)
	content_type := http.DetectContentType(img)
	if content_type != "image/png" {
		return errors.New("img content type must image/png")
	}

	f, _ := os.Create(fileName)
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			return
		}
	}(f)

	_, err = f.Write(img)
	if err != nil {
		return err
	}
	return nil
}
