package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/logeshwarann-dev/bits-bank_auth-service/db"
	"github.com/logeshwarann-dev/bits-bank_auth-service/session"
	"github.com/logeshwarann-dev/bits-bank_auth-service/utils"
	"gorm.io/gorm"
)

var PgDb *gorm.DB

type DwollaCustomerInfo struct {
	CustomerId  string `json:"customer_id"`
	CustomerUrl string `json:"customer_url"`
}

func SignUp(c *gin.Context) {
	var signupForm db.SignUpForm
	if err := c.ShouldBindJSON(&signupForm); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUserID := fmt.Sprintf("BITS%s%s", strings.ToUpper(signupForm.FirstName[:3]), time.Now().Format("20060102150405"))
	signupForm.UserId = newUserID
	newAccount := signupForm.ConvertToUser()

	if err := utils.CreateAccount(PgDb, *newAccount); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	newUser := db.LoggedInUser{
		Username:    newAccount.Email,
		Email:       newAccount.Email,
		FirstName:   newAccount.FirstName,
		LastName:    newAccount.LastName,
		Address1:    newAccount.Address1,
		City:        newAccount.City,
		State:       newAccount.State,
		PostalCode:  newAccount.PostalCode,
		DateOfBirth: newAccount.DateOfBirth,
		AadharNo:    newAccount.AadharNo,
		UserId:      newAccount.UserID,
	}

	var dwollaCustomerResponse DwollaCustomerInfo
	if err := utils.SendPostRequest(utils.DwollaCreateCustomerUrl, newUser, &dwollaCustomerResponse); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("unable to create dwolla account: %v", err.Error())})
		return
	}

	if err := utils.UpdateUserWithDwollaInfo(dwollaCustomerResponse.CustomerId, dwollaCustomerResponse.CustomerUrl, PgDb, newUser.Email); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("unable to update dwolla info in db: %v", err.Error())})
		return
	}

	CompletedUserAccount, err := db.GetRecordUsingEmail(PgDb, newUser.Email)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("unable to get updated user info from db: %v", err.Error())})
		return
	}

	sessionToken, err := utils.CreateSession(CompletedUserAccount.Email)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("unable to create session: %v", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "User signed up successfully",
		"session_token": sessionToken,
		"user":          CompletedUserAccount,
	})
}

func SignIn(c *gin.Context) {
	var newLogin db.SignInForm
	if err := c.ShouldBindJSON(&newLogin); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint("error binding request: ", err.Error())})
		return
	}

	isValid, validUser, err := utils.VerifyCredentials(PgDb, newLogin)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error verifying user: %v", err.Error())})
		return
	}

	if !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid user credentials. Please try again."})
		return
	}

	var sessionToken string
	authHeader := c.GetHeader("Authorization")
	log.Println("Auth header: ", authHeader)
	if authHeader != "" {
		splitToken := strings.Split(authHeader, " ")
		if len(splitToken) != 2 || splitToken[0] != "Bearer" {
			log.Println("error: Invalid token format")
		}
		sessionToken = splitToken[1]
	}

	if len(sessionToken) > 10 {
		log.Println("Sesion token present")
		username, err := session.VerifyToken(sessionToken)
		if err == nil {
			if err = utils.ValidateSession(username, sessionToken); err == nil {
				log.Println("Valid Session")
				c.JSON(http.StatusOK, gin.H{"message": "User authentication successful!"})
				return
			}
			log.Println("Session expired: ", err.Error())
		}
		log.Println("Invalid token", err.Error())
	}

	sessionToken, err = utils.CreateSession(newLogin.Email)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("unable to create session: %v", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User Authentication successful",
		"session_token": sessionToken,
		"user":          validUser,
	})

}

func GetLoggedInUser(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")
	log.Println("Auth header: ", authHeader)
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	splitToken := strings.Split(authHeader, " ")
	if len(splitToken) != 2 || splitToken[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}
	sessionToken := splitToken[1]

	username, err := session.VerifyToken(sessionToken)
	if err != nil {
		log.Println(err.Error())
		log.Println("error verifying token: ", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if err = utils.ValidateSession(username, sessionToken); err != nil {
		log.Println(err.Error())
		log.Println("error validation session: ", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Session expired: %v", err.Error())})
		return
	}

	loggedInUser, err := utils.GetUserDetails(username, PgDb)
	if err != nil {
		log.Println(err.Error())
		log.Println("error getting user: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("unable to get user: %v", err.Error())})
		return
	}

	c.JSON(http.StatusOK, loggedInUser)
}

func SignOut(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	splitToken := strings.Split(authHeader, " ")
	if len(splitToken) != 2 || splitToken[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}
	sessionToken := splitToken[1]

	username, err := session.VerifyToken(sessionToken)
	if err != nil {
		log.Println(err.Error())
		log.Println("error verifying token: ", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if err = utils.DeleteSession(username); err != nil {
		log.Println(err.Error())
		log.Println("Error deleting session from Redis:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete session"})
		return

	}

	c.JSON(http.StatusOK, gin.H{"message": "user session deleted"})

}
