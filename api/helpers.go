package api

import (
	"fmt"

	"github.com/logeshwarann-dev/bits-bank_auth-service/db"
	"gorm.io/gorm"
)

func CreateAccount(authdb *gorm.DB, user db.BankUser) error {
	if err := db.AddUser(authdb, user); err != nil {
		return fmt.Errorf("error adding user: %v", err.Error())
	}
	return nil

}

func VerifyCredentials(authdb *gorm.DB, loginDetails db.SignInForm) (bool, error) {

	user, err := db.GetRecordUsingEmail(authdb, loginDetails.Email)
	if err != nil {
		return false, err
	}
	if user.Password == loginDetails.Password {
		fmt.Println("Valid User")
		return true, nil
	}
	return false, nil

}
