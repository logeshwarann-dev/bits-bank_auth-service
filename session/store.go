package session

import (
	"context"
	"fmt"
	"os"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
)

var (
	RedisClient *redis.Client
	SECRET_KEY  []byte
	Ctx         context.Context
)

func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})
	tokenString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) (string, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return SECRET_KEY, nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token: %v", err.Error())
	}

	username := claims["username"].(string)
	return username, nil
}

func InitConnToRedis() error {
	connString := fmt.Sprintf("%s:%s", os.Getenv("REDIS_DB_ADDR"), os.Getenv("REDIS_DB_PORT"))
	redis := redis.NewClient(&redis.Options{
		Addr:     connString,
		Password: "",
		DB:       0,
	})
	RedisClient = redis
	Ctx = context.Background()
	fmt.Println("Redis connection successful!")
	return nil
}

func SetSessionInRedis(username string, jwtToken string) error {
	sessionKey := fmt.Sprintf("session:%s", username)
	if err := RedisClient.Set(Ctx, sessionKey, jwtToken, 24*time.Hour).Err(); err != nil {
		return fmt.Errorf("error while setting session in redis: %v", err.Error())
	}
	fmt.Println("Session was set successfully!")
	return nil
}

func GetSessionFromRedis(sessionKey string) (string, error) {
	storedToken, err := RedisClient.Get(Ctx, sessionKey).Result()
	if err != nil {
		return "", fmt.Errorf("error Session expired: %v", err.Error())
	}
	return storedToken, nil
}

func DeleteSessionInRedis(sessionKey string) error {
	deletedCount, err := RedisClient.Del(Ctx, sessionKey).Result()
	if err != nil {
		return fmt.Errorf("error deleting session: %v", err.Error())
	}
	fmt.Println(deletedCount, " Session deleted in Redis")
	return nil
}
