package db

import (
	"errors"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB_USER string
	DB_HOST string
	DB_NAME string
	DB_PWD  string
	DB_PORT string
	DB_SSL  string
)

func ConnectToDB() *gorm.DB {

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s", DB_USER, DB_PWD, DB_NAME, DB_HOST, DB_PORT, DB_SSL)
	gormDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Error connection to DB: ", err.Error())
	}

	fmt.Println("DB Connection Successful!")
	return gormDb
}

func AddUser(authdb *gorm.DB, bankUser BankUser) error {
	if err := authdb.Create(&bankUser).Error; err != nil {
		fmt.Println(err.Error())
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

func UpdateRecord(authdb *gorm.DB, record BankUser, field string, newValue string) error {
	result := authdb.Model(&record).Update(field, newValue)
	if result.Error != nil {
		return fmt.Errorf("unable to update record: %v", result.Error.Error())
	}
	return nil
}
