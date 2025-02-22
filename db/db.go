package db

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDB() *gorm.DB {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: ", err.Error())
	}
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s", os.Getenv("DB_USER"), os.Getenv("DB_PWD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_SSL"))
	gormDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connection to DB: ", err.Error())
	}

	fmt.Println("DB Connection Successful!")
	return gormDb
}

func AddUser(authdb *gorm.DB, bankUser BankUser) error {
	if err := authdb.Create(&bankUser).Error; err != nil {
		return err
	}
	return nil
}

func GetRecordUsingEmail(authdb *gorm.DB, email string) (BankUser, error) {
	var user BankUser
	result := authdb.Where("email = ?", email).First(&user)
	if result.Error != nil {
		fmt.Println("Error: ", result.Error)
		return BankUser{}, errors.New("no records found")
	}
	return user, nil
}
