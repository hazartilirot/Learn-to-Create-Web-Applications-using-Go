package main

import (
	"fmt"
	"github.com/username/project-name/models"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "lenslocked_dev"
)

func main() {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	us, err := models.NewUserService(dsn)
	if err != nil {
		panic(err)
	}
	us.ResetDB()

	user := models.User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	if err := us.Create(&user); err != nil {
		panic(err)
	}

	user.Email = "johndoe@example.com"

	if err = us.Update(&user); err != nil {
		panic(err)
	}

	if err := us.Delete(user.ID); err != nil {
		panic(err)
	}

	userByID, err := us.ByID(user.ID)
	if err != nil {
		panic(err)
	}

	fmt.Println(userByID)
}
