package helpers

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"strings"
	"time"
)

const ExpireJWT = 1200000
const ACCESS_TOKEN_EXPIRE = 1800    // 30 minute
const REFRESH_TOKEN_EXPIRE = 604800 // week

var secretKey = []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImJhcnNpbCIsImV4cCI6MTUxMjE5OTk3NzExNSwiaWF0IjoxNTEyMTk4Nzc3MTE1LCJuYmYiOjE1MTIxOTg3NzcxMTV9.VqaOlTi8Z2kH_SrAY-qmzWnKV5E_ADk0iN4TFnQDFLc")

type MyCustomClaims struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Last     []byte `json:"last"`
	jwt.StandardClaims
}

type AccessTokenClaims struct {
	UUID string `json:"uuid"`
	jwt.StandardClaims
}

type RefreshTokenClaims struct {
	UUID string `json:"uuid"`
	jwt.StandardClaims
}

func GetStClaims(notBefore, issueAt, expireAt int64) jwt.StandardClaims {
	return jwt.StandardClaims{IssuedAt: issueAt, ExpiresAt: expireAt, NotBefore: notBefore}
}

func GetExpireInMSEC() int64 {
	return (time.Now().Unix()) + ExpireJWT
}

func GetAccessTokenExpireInMSEC() int64 {
	return (time.Now().Unix()) + ACCESS_TOKEN_EXPIRE
}

func GetRefreshTokenExpireInMSEC() int64 {
	return (time.Now().Unix()) + REFRESH_TOKEN_EXPIRE
}

func GetNowInMSEC() int64 {
	return (time.Now().Unix())
}

func GetNewClaims(id *uuid.UUID, username string) MyCustomClaims {
	nbf := GetNowInMSEC() - 100
	iat := GetNowInMSEC()
	exp := GetExpireInMSEC()
	last, _ := time.Now().MarshalText()
	str_id := id.String()

	stClaims := GetStClaims(nbf, iat, exp)
	return MyCustomClaims{str_id, username, last, stClaims}
}

func GetNewAccessClaims(id *uuid.UUID) (AccessTokenClaims, int64) {
	nbf := GetNowInMSEC() - 100
	iat := GetNowInMSEC()
	exp := GetAccessTokenExpireInMSEC()
	str_id := id.String()

	stClaims := GetStClaims(nbf, iat, exp)
	return AccessTokenClaims{str_id, stClaims}, exp
}

func GetNewRefreshClaims(id *uuid.UUID) RefreshTokenClaims {
	nbf := GetNowInMSEC() - 200
	iat := GetNowInMSEC()
	exp := GetRefreshTokenExpireInMSEC()
	str_id := id.String()

	stClaims := GetStClaims(nbf, iat, exp)
	return RefreshTokenClaims{str_id, stClaims}
}

func GenerateToken(c MyCustomClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenString, _ := token.SignedString(secretKey)

	return tokenString
}

func GenerateAccessToken(c AccessTokenClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenString, _ := token.SignedString(secretKey)

	return tokenString
}

func GenerateRefreshToken(c RefreshTokenClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenString, _ := token.SignedString(secretKey)

	return tokenString
}

func ValidateToken(tokenStr string) (*MyCustomClaims, error) {

	tokenStr = strings.TrimSpace(tokenStr)
	token, err := jwt.ParseWithClaims(tokenStr, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected siging method")
		}
		return secretKey, nil
	})

	if err != nil {
		return &MyCustomClaims{}, err
	}

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return claims, errors.New("Token is not valid")
	}
}

func ValidateAccessToken(tokenStr string) (*AccessTokenClaims, error) {

	tokenStr = strings.TrimSpace(tokenStr)
	token, err := jwt.ParseWithClaims(tokenStr, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected siging method")
		}
		return secretKey, nil
	})

	if err != nil {
		return &AccessTokenClaims{}, err
	}

	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims, nil
	} else {
		return claims, errors.New("Access Token is not valid")
	}
}

func ValidateRefreshToken(tokenStr string) (*RefreshTokenClaims, error) {

	tokenStr = strings.TrimSpace(tokenStr)
	token, err := jwt.ParseWithClaims(tokenStr, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected siging method")
		}
		return secretKey, nil
	})

	if err != nil {
		return &RefreshTokenClaims{}, err
	}

	if claims, ok := token.Claims.(*RefreshTokenClaims); ok && token.Valid {
		return claims, nil
	} else {
		return claims, errors.New("Refresh Token is not valid")
	}
}
