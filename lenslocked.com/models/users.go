package models

import (
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	/*ErrNotFound is returned when a resource cannot be found in the database*/
	ErrNotFound = errors.New("models: resource not found")
	/*ErrInvalidID is returned when invalid ID is provided to a method like Delete*/
	ErrInvalidID = errors.New("models: ID must be other than 0")
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

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}

	return us.db.Delete(&user).Error
}

func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

/*byID will look up a user with a provided ID. nil is returned if the user is found otherwise
ErrNotFound will be returned. Other errors like connecting to a database would be returned with
500 status code.*/
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

func (us *UserService) Create(user *User) error {
	return us.db.Create(&user).Error
}

func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
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
