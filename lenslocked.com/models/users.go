package models

import (
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("models: resource not found")
)

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open(postgres.Open(connectionInfo), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	return &UserService{db: db}, nil
}

type UserService struct {
	db *gorm.DB
}

/*byID will look up a user with a provided ID. nil is returned if the user is found otherwise
ErrNotFound will be returned. Other errors like connecting to a database would be returned with
500 status code.*/
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	err := us.db.Where("id = ?", id).First(&user).Error

	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (us *UserService) ResetDB() {
	us.db.Migrator().DropTable(&User{})
	us.db.AutoMigrate(&User{})
}

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null; uniqueIndex"`
}
