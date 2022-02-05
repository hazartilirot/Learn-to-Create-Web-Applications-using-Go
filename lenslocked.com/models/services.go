package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open(postgres.Open(connectionInfo), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Services{
		User: NewUserService(db),
	}, nil
}

type Services struct {
	Gallery GalleryService
	User    UserService
}
