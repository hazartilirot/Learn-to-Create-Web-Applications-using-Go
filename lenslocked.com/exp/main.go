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
	//us.ResetDB()
	user, err := us.ByID(1)
	if err != nil {
		panic(err)
	}

	fmt.Println(user)
}
