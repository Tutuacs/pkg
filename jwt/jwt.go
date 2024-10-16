package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/Tutuacs/pkg/config"
	"github.com/Tutuacs/pkg/enums"
)

func CreateJWT(email string, userID int64, role enums.Role) (string, error) {
	expiration := time.Second * time.Duration(config.GetJWT().JWT_EXP)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(int(userID)),
		"email":     email,
		"role":      role,
		"expiresAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.GetJWT().JWT_SECRET))
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(string(config.GetJWT().JWT_SECRET)), nil
	})
}
