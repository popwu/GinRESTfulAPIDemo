package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID uint, nickname string) (string, error) {
	// 从环境变量中获取 JWT_SIGNING_KEY
	signingKey := os.Getenv("JWT_SIGNING_KEY")
	if signingKey == "" {
		return "", errors.New("JWT_SIGNING_KEY is not set in the environment")
	}

	claims := jwt.MapClaims{
		"user_id":  userID,
		"nickname": nickname,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", errors.New("failed to sign the token")
	}

	return signedToken, nil
}

func ParseToken(tokenString string) (userID uint, nickname string, err error) {
	signingKey := os.Getenv("JWT_SIGNING_KEY")
	if signingKey == "" {
		return 0, "", errors.New("JWT_SIGNING_KEY is not set in the environment")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return 0, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return 0, "", errors.New("invalid user_id in token")
		}
		nickname, ok := claims["nickname"].(string)
		if !ok {
			return 0, "", errors.New("invalid nickname in token")
		}
		return uint(userID), nickname, nil
	}

	return 0, "", errors.New("invalid token")
}
