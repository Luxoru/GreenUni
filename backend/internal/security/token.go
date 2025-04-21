package security

import (
	"backend/internal/models"
	"fmt"
	"github.com/google/uuid"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// TODO: maybe move to env?
// Private token for signing
var secretKey = []byte("SUPERSECRETTOKEN")

func GenerateJWT(model *models.UserInfoModel) (string, error) {

	claims := jwt.MapClaims{
		"uuid":     model.UUID,
		"username": model.Username,
		"role":     model.Role.String(),
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // Token expiration time = 3 days
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseJWT(tokenStr string) (*models.UserInfoModel, error) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		uuidStr, ok1 := claims["uuid"].(string)
		username, ok2 := claims["username"].(string)
		roleName, ok3 := claims["role"].(string)

		if !ok1 || !ok2 || !ok3 {
			return nil, fmt.Errorf("invalid token claims")
		}

		userUUID, err := uuid.Parse(uuidStr)
		if err != nil {
			return nil, err
		}

		roleType, err := models.ParseRoleType(roleName)
		if err != nil {
			return nil, err
		}

		return &models.UserInfoModel{
			UUID:     userUUID,
			Username: username,
			Role:     roleType,
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}
