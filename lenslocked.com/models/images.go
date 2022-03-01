package models

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Image is not stored in database
type Image struct {
	GalleryID uint
	Filename  string
}

func (i *Image) Path() string {
	temp := url.URL{
		Path: "/" + i.RelativePath(),
	}
	return temp.String()
}

func (i *Image) RelativePath() string {
	return fmt.Sprintf("images/galleries/%v/%v", i.GalleryID, i.Filename)
}

type ImageService interface {
	Create(galleryID uint, r io.ReadCloser, filename string) error
	ByGalleryID(galleryID uint) ([]Image, error)
	Delete(i *Image) error
}

func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct{}

func (is *imageService) Create(galleryID uint, r io.ReadCloser, filename string) error {
	defer r.Close()
	path, err := is.mkImagePath(galleryID)
	if err != nil {
		return err
	}
	dst, err := os.Create(path + filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, r)
	if err != nil {
		return err
	}
	return nil
}

func (is *imageService) imagePath(galleryID uint) string {
	return fmt.Sprintf("images/galleries/%v/", galleryID)
}

func (is *imageService) ByGalleryID(galleryID uint) ([]Image, error) {
	path := is.imagePath(galleryID)
	imgStrings, err := filepath.Glob(path + "*")
	if err != nil {
		return nil, err
	}
	ret := make([]Image, len(imgStrings))

	for i := range imgStrings {
		imgStrings[i] = strings.Replace(imgStrings[i], path, "", 1)
		ret[i] = Image{
			Filename:  imgStrings[i],
			GalleryID: galleryID,
		}
	}
	return ret, nil
}

func (is *imageService) Delete(i *Image) error {
	return os.Remove(i.RelativePath())
}

func (is *imageService) mkImagePath(galleryID uint) (string, error) {
	galleryPath := is.imagePath(galleryID)
	err := os.MkdirAll(galleryPath, 0755)
	if err != nil {
		return "", err
	}
	return galleryPath, nil
}
