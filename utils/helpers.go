package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/logeshwarann-dev/bits-bank_auth-service/db"
	"github.com/logeshwarann-dev/bits-bank_auth-service/session"
	"gorm.io/gorm"
)

func CreateAccount(authdb *gorm.DB, user db.BankUser) error {
	if err := db.AddUser(authdb, user); err != nil {
		return fmt.Errorf("error adding user: %v", err.Error())
	}
	return nil

}

func VerifyCredentials(authdb *gorm.DB, loginDetails db.SignInForm) (bool, db.LoggedInUser, error) {

	user, err := db.GetRecordUsingEmail(authdb, loginDetails.Email)
	if err != nil {
		return false, db.LoggedInUser{}, err
	}
	if user.Password == loginDetails.Password {
		fmt.Println("Valid User")
		return true, db.LoggedInUser{
			Username:    user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Address1:    user.Address1,
			City:        user.City,
			State:       user.State,
			PostalCode:  user.PostalCode,
			DateOfBirth: user.DateOfBirth,
			AadharNo:    user.AadharNo,
		}, nil
	}
	return false, db.LoggedInUser{}, nil

}

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: ", err.Error())
	}

	db.DB_HOST = os.Getenv("DB_HOST")
	db.DB_PWD = os.Getenv("DB_PWD")
	db.DB_NAME = os.Getenv("DB_NAME")
	db.DB_PORT = os.Getenv("DB_PORT")
	db.DB_USER = os.Getenv("DB_USER")
	db.DB_SSL = os.Getenv("DB_SSL")
	session.SECRET_KEY = []byte(os.Getenv("JWT_SECRET_KEY"))

}

func CreateSession(username string) (string, error) {
	jwtToken, err := session.CreateToken(username)
	if err != nil {
		fmt.Println("error creating token: ", err.Error())
		return "", fmt.Errorf("error creating token: %v", err.Error())
	}

	if err = session.SetSessionInRedis(username, jwtToken); err != nil {
		fmt.Println("error while setting session in redis")
		return "", err
	}

	fmt.Println("User session has been created!")
	return jwtToken, nil
}

func ValidateSession(user string, sessionToken string) error {
	sessionKey := fmt.Sprintf("session:%s", user)
	storedToken, err := session.GetSessionFromRedis(sessionKey)
	if err != nil || storedToken != sessionToken {
		return fmt.Errorf("session validation failed: %v", err)
	}
	return nil
}

func GetUserDetails(username string, authdb *gorm.DB) (db.LoggedInUser, error) {
	user, err := db.GetRecordUsingEmail(authdb, username)
	if err != nil {
		return db.LoggedInUser{}, err
	}

	return db.LoggedInUser{
		Username:    user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Address1:    user.Address1,
		City:        user.City,
		State:       user.State,
		PostalCode:  user.PostalCode,
		DateOfBirth: user.DateOfBirth,
		AadharNo:    user.AadharNo,
	}, nil

}

func DeleteSession(user string) error {
	sessionKey := fmt.Sprintf("session:%s", user)
	if err := session.DeleteSessionInRedis(sessionKey); err != nil {
		return err
	}
	return nil

}
