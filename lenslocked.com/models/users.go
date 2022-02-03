package models

import (
	"errors"
	"github.com/username/project-name/hash"
	"github.com/username/project-name/rand"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"strings"
)

var (
	/*ErrNotFound is returned when a resource cannot be found in the database*/
	ErrNotFound = errors.New("models: resource not found")
	/*ErrInvalidID is returned when invalid ID is provided to a method like Delete*/
	ErrInvalidID = errors.New("models: ID must be other than 0")
	/*ErrInvalidEmailOrPassword is returned when invalid password is typed in*/
	ErrInvalidEmailOrPassword = errors.New("models: Invalid email or password. Please try again")
	ErrInvalidToken           = errors.New("models: The token is invalid or empty")

	ErrEmailRequired     = errors.New("models: Email address is required")
	ErrEmailInvalid      = errors.New("models: Email address is invalid")
	ErrEmailIsRegistered = errors.New("modules: Email is already registered")
)

const hmacSecretKey = "secret-hmac-key"

/*User represents the user's model stored in our database. It's used for users accounts.
Storing both an email and password so users can log in and gain access to their content*/
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null; uniqueIndex"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null; uniqueIndex"`
}

/*UserDB is used to interact with the database*/
type UserDB interface {
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	AutoMigrate() error
	ResetDB() error
}

type UserService interface {
	Authenticate(email, password string) (*User, error)
	UserDB
}

func NewUserService(connectionInfo string) (*userService, error) {

	ug, err := newUserGorm(connectionInfo)

	if err != nil {
		return nil, err
	}

	hmac := hash.NewHMAC(hmacSecretKey)

	uv := newUserValidator(ug, hmac)

	return &userService{
		UserDB: uv,
	}, nil
}

var _ UserService = &userService{}

type userService struct {
	UserDB
}

func (us *userService) Authenticate(email, password string) (*User, error) {
	user, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidEmailOrPassword
		default:
			return nil, err
		}
	}
	return user, nil
}

type userValFunc func(*User) error

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

var _ UserDB = &userValidator{}

func newUserValidator(udb UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB:     udb,
		hmac:       hmac,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

type userValidator struct {
	UserDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
}

func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	if err := runUserValFuncs(&user, uv.normalizeEmail); err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runUserValFuncs(&user, uv.hmacRemember); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

func (uv *userValidator) Create(user *User) error {
	if err := runUserValFuncs(
		user,
		uv.bcryptPassword,
		uv.setRememberIfUnset,
		uv.hmacRemember,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail,
	); err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

func (uv *userValidator) Update(user *User) error {
	if err := runUserValFuncs(
		user,
		uv.bcryptPassword,
		uv.hmacRemember,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail,
	); err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := runUserValFuncs(&user, uv.idGreaterThan(0))
	if err != nil {
		return err
	}
	return uv.UserDB.Delete(id)
}

func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""

	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

type userGorm struct {
	db *gorm.DB
}

func (uv *userValidator) setRememberIfUnset(user *User) error {
	if user.Remember != "" {
		return nil
	}

	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

func (uv *userValidator) idGreaterThan(n uint) userValFunc {
	return userValFunc(func(user *User) error {
		if user.ID <= 0 {
			return ErrInvalidID
		}
		return nil
	})
}

func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	if user.Email == "" {
		return nil
	}
	if !uv.emailRegex.MatchString(user.Email) {
		return ErrEmailInvalid
	}
	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

func (uv *userValidator) emailIsAvail(user *User) error {
	existing, err := uv.ByEmail(user.Email)
	if err == ErrNotFound {
		return nil
	}
	if err != nil {
		return err
	}
	if user.ID != existing.ID {
		return ErrEmailIsRegistered
	}
	return nil
}

var _ UserDB = &userGorm{}

func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open(postgres.Open(connectionInfo), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	return &userGorm{
		db: db,
	}, nil
}

func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

/*ByRemember looks up the user with a token and returns that user.*/
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	if rememberHash == "" {
		return nil, ErrInvalidToken
	}

	var user User
	db := ug.db.Where("remember_hash = ?", rememberHash)
	err := first(db, &user)

	if err != nil {
		return nil, ErrInvalidToken
	}

	return &user, nil
}

func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

/*byID will look up a user with a provided ID. nil is returned if the user is found otherwise
ErrNotFound will be returned. Other errors like connecting to a database would be returned with
500 status code.*/
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

func (ug *userGorm) ResetDB() error {
	if err := ug.db.Migrator().DropTable(&User{}); err != nil {
		return err
	}
	if err := ug.AutoMigrate(); err != nil {
		return err
	}
	return nil
}

func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}); err != nil {
		return err
	}
	return nil
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}
