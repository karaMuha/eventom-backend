package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var PrivateKey *rsa.PrivateKey
var ProtectedRoutes map[string]bool

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func ReadPrivateKeyFromFile(filename string) (*rsa.PrivateKey, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}

	data, _ := pem.Decode(buffer)
	privateKey, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func GenerateJwt(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"userId": userId,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(time.Hour).Unix(),
	})

	return token.SignedString(PrivateKey)
}

func VerifyJwt(jwtToken string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return PrivateKey.Public(), nil
	})

	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	return parsedToken, nil
}
