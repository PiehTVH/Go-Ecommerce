package helper

import (
	"errors"
	"os"
	"time"

	"github.com/PiehTVH/go-ecommerce/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func CheckUserValidation(u types.UserClient) error {
	if u.Email == "" {
		return errors.New("email can't be empty")
	}
	if u.Name == "" {
		return errors.New("name can't be empty")
	}
	if u.Phone == "" {
		return errors.New("phone can't be empty")
	}
	if u.Password == "" {
		return errors.New("password can't be empty")
	}
	return nil
}

func EncryptPassword(s string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return ""
	}
	return string(bytes) // Example: $2a$04$ky/fH/EcjZwcI0ZIJbYA8eUMjStsW.0D3ETbIxAX7HvR7TL1d.7x2
}

func ComparePassword(hashedPwd string, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	return err == nil
}

func IsUserAdmin(c *gin.Context, tokenString string) (bool, error) {

	return true, nil

}

func GenerateToken(userId string, email string, userType string) (string, error) {
	secretKey := os.Getenv("secretKey")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"type":   userType,
		"exp":    time.Now().Add(time.Hour * 2).Unix(), // Ex: 1729368780
	})

	return token.SignedString([]byte(secretKey))
}

func VerifyToken(tokenString string) (string, string, error) {
	secretKey := os.Getenv("secretKey")
	token, err := jwt.Parse((tokenString), func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("could not parse claims")
	}

	email := claims["type"].(string)
	userType := claims["type"].(string)

	return email, userType, nil
}
