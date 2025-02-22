package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/logeshwarann-dev/bits-bank_auth-service/db"
	"gorm.io/gorm"
)

var PostDb *gorm.DB

func SignUp(c *gin.Context) {
	var signupForm db.SignUpForm
	if err := c.ShouldBindJSON(&signupForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// userID := fmt.Sprintf("%s%s", newUser.FirstName[:3], time.Now().Format("20060102150405000"))
	newAccount := signupForm.ConvertToUser()

	if err := CreateAccount(PostDb, *newAccount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, newAccount)
}

func SignIn(c *gin.Context) {
	var newLogin db.SignInForm
	if err := c.ShouldBindJSON(&newLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint("error binding request: ", err.Error())})
		return
	}

	validUser, err := VerifyCredentials(PostDb, newLogin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error verifying user: %v", err.Error())})
		return
	}

	if !validUser {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user credentials. Please try again."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User Authentication successful"})

}
