package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/logeshwarann-dev/bits-bank_auth-service/db"
	"github.com/logeshwarann-dev/bits-bank_auth-service/session"
	"gorm.io/gorm"
)

var DwollaCreateCustomerUrl string

func CreateAccount(authdb *gorm.DB, user db.BankUser) error {
	if err := db.AddUser(authdb, user); err != nil {
		log.Println(err.Error())
		return fmt.Errorf("error adding user: %v", err.Error())
	}
	return nil

}

func VerifyCredentials(authdb *gorm.DB, loginDetails db.SignInForm) (bool, db.LoggedInUser, error) {

	user, err := db.GetRecordUsingEmail(authdb, loginDetails.Email)
	if err != nil {
		log.Println(err.Error())
		return false, db.LoggedInUser{}, err
	}
	if user.Password == loginDetails.Password {
		log.Println("Valid User")
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
			UserId:      user.UserID,
			Email:       user.Email,
		}, nil
	}
	return false, db.LoggedInUser{}, nil

}

func LoadEnv() {
	// if err := godotenv.Load(); err != nil {
	// 	log.Println(err.Error())
	// 	log.Fatal("Error loading .env file: ", err.Error())
	// }

	db.DB_HOST = os.Getenv("DB_HOST")
	db.DB_PWD = os.Getenv("DB_PWD")
	db.DB_NAME = os.Getenv("DB_NAME")
	db.DB_PORT = os.Getenv("DB_PORT")
	db.DB_USER = os.Getenv("DB_USER")
	db.DB_SSL = os.Getenv("DB_SSL")
	session.SECRET_KEY = []byte(os.Getenv("JWT_SECRET_KEY"))
	DwollaCreateCustomerUrl = os.Getenv("PLAID_SERVICE_CREATE_CUSTOMER_URL")

}

func CreateSession(username string) (string, error) {
	jwtToken, err := session.CreateToken(username)
	if err != nil {
		log.Println(err.Error())
		log.Println("error creating token: ", err.Error())
		return "", fmt.Errorf("error creating token: %v", err.Error())
	}

	if err = session.SetSessionInRedis(username, jwtToken); err != nil {
		log.Println(err.Error())
		log.Println("error while setting session in redis")
		return "", err
	}

	log.Println("User session has been created!")
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
		log.Println(err.Error())
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
		UserId:      user.UserID,
		Email:       user.Email,
	}, nil

}

func DeleteSession(user string) error {
	sessionKey := fmt.Sprintf("session:%s", user)
	if err := session.DeleteSessionInRedis(sessionKey); err != nil {
		log.Println(err.Error())
		return err
	}
	return nil

}

func SendPostRequest(targetUrl string, payload any, responseContainer any) error {
	requestPayload, _ := json.Marshal(payload)

	response, err := http.Post(targetUrl, "application/json", bytes.NewBuffer(requestPayload))
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("error while sending post request: %v", err.Error())
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("error while reading response: %v", err.Error())
	}

	result := string(responseBody)
	log.Println("Response from Plaid service: ", result)
	if err = json.Unmarshal(responseBody, responseContainer); err != nil {
		log.Println(err.Error())
		return fmt.Errorf("error while unmarshalling response: %v", err.Error())
	}
	return nil
}

func UpdateUserWithDwollaInfo(dwollaCustomerId string, dwollaCustomerUrl string, pgDb *gorm.DB, userEmail string) error {
	var existingRecord db.BankUser
	dwollaCustomerIdColumnName := "dwolla_customer_id"
	dwollaCustomerUrlColumnName := "dwolla_customer_url"
	existingRecord, err := db.GetRecordUsingEmail(pgDb, userEmail)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("unable to get record: %v", err.Error())
	}

	if err = db.UpdateRecord(pgDb, existingRecord, dwollaCustomerIdColumnName, dwollaCustomerId); err != nil {
		log.Println(err.Error())
		return fmt.Errorf("unable to update dwolla customer Id: %v", err.Error())
	}

	existingRecord, err = db.GetRecordUsingEmail(pgDb, userEmail)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("unable to get record: %v", err.Error())
	}

	if err = db.UpdateRecord(pgDb, existingRecord, dwollaCustomerUrlColumnName, dwollaCustomerUrl); err != nil {
		log.Println(err.Error())
		return fmt.Errorf("unable to update dwolla customer url:: %v", err.Error())
	}

	return nil
}
