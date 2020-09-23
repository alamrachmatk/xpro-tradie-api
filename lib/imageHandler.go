package lib

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func UploadImage(file *multipart.FileHeader, target string) error {
	// Open file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Prepare destination file
	dst, err := os.Create(target)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func Remove(path string) error {
	files, err := filepath.Glob(path+"*")
	if err != nil {
		return err
	}

	for _, delfile := range files {
		err = os.Remove(delfile)
		 if err != nil {
			return err
		}
	}
	return nil
}
