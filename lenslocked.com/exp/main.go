package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "lenslocked_dev"
)

type User struct {
	gorm.Model
	Age       int
	FirstName string
	LastName  string
	Email     string `gorm:"not null; uniqueIndex"`
	Orders    []Order
}

type Order struct {
	gorm.Model
	UserId      uint
	Amount      int
	Description string
}

func main() {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	/*if err := db.Migrator().DropTable(&User{}); err != nil {
		panic(err)
	}*/

	if err := db.AutoMigrate(&User{}, &Order{}); err != nil {
		panic(err)
	}

	var u User
	if err := db.Preload("Orders").First(&u).Error; err != nil {
		panic(err)
	}

	/*createOrder(db, u, 1001, "Fake description #1")
	createOrder(db, u, 9999, "Fake description #2")
	createOrder(db, u, 2000, "Fake description #3")
	createOrder(db, u, 5555, "Fake description #4")*/
	fmt.Println(u)
}

func createOrder(db *gorm.DB, user User, amount int, desc string) {
	err := db.Create(&Order{
		UserId:      user.ID,
		Amount:      amount,
		Description: desc,
	}).Error

	if err != nil {
		panic(err)
	}
}
