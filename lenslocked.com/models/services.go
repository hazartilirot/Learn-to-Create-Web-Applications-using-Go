package models

import (
	"gorm.io/gorm"
)

type ServicesConfig func(*Services) error

func WithGorm(dialect gorm.Dialector) ServicesConfig {
	return func(s *Services) error {
		db, err := gorm.Open(dialect, &gorm.Config{})
		if err != nil {
			return err
		}
		s.db = db
		return nil
	}
}

func WithUser(hmacKey string) ServicesConfig {
	return func(s *Services) error {
		s.User = NewUserService(s.db, hmacKey)
		return nil
	}
}

func WithGallery() ServicesConfig {
	return func(s *Services) error {
		s.Gallery = NewGalleryService(s.db)
		return nil
	}
}

func WithImage() ServicesConfig {
	return func(s *Services) error {
		s.Image = NewImageService()
		return nil
	}
}

func NewServices(cfgs ...func(*Services) error) (*Services, error) {
	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

type Services struct {
	Gallery GalleryService
	User    UserService
	Image   ImageService
	db      *gorm.DB
}

/*ResetDB drops all tables and then recreates them*/
func (s *Services) ResetDB() error {
	s.db.Migrator().DropTable(&User{}, &Gallery{}, &pwReset{})
	if err := s.db.AutoMigrate(); err != nil {
		return err
	}
	return s.AutoMigrate()
}

/*AutoMigrate will create tables, indexes, etc */
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}, &pwReset{})
}
