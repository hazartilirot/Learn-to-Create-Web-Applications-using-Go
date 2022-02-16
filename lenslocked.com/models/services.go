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
		Gallery: NewGalleryService(db),
		User:    NewUserService(db),
		Image:   NewImageService(),
		db:      db,
	}, nil
}

type Services struct {
	Gallery GalleryService
	User    UserService
	Image   ImageService
	db      *gorm.DB
}

/*ResetDB drops all tables and then recreates them*/
func (s *Services) ResetDB() error {
	s.db.Migrator().DropTable(&User{}, &Gallery{})
	if err := s.db.AutoMigrate(); err != nil {
		return err
	}
	return s.AutoMigrate()
}

/*AutoMigrate will create tables, indexes, etc */
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{})
}
