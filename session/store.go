package session

import (
	"os"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
)

var secretKey = []byte("bitsbank")

func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func InitConnToRedis() *redis.Client {
	dbAddr := os.Getenv("REDIS_DB_ADDR")
	redis := redis.NewClient(&redis.Options{
		Addr:     dbAddr,
		Password: "",
		DB:       0,
	})
	return redis

}
