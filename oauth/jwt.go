package oauth

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var HmacSecret []byte = []byte("jwt-secrie")

func CreateJWT(payload string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": payload,
		"nbf":    time.Now().Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(HmacSecret)
}

// func IsValidToken(token string) (*gql.User, bool) {
// 	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
// 		// from github.com/dhrijalva/jwt-go/hmac.go we should return a []byte
// 		// as we only use one single key, we just return it
// 		return HmacSecret, nil
// 	})
// 	if err != nil {
// 		return nil, false
// 	}
// 	claim, ok := parsedToken.Claims.(jwt.MapClaims)
// 	if !ok {
// 		return nil, false
// 	}
// 	userID, ok := claim["userId"]
// 	if !ok {
// 		return nil, false
// 	}
// 	user, err := database.FindOneUserById(userID)
// 	if err != nil {
// 		return nil, false
// 	}
// 	return user, true
// }

// func IsValidToken(token string) (string, error) {
// 	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
// 		// from github.com/dhrijalva/jwt-go/hmac.go we should return a []byte
// 		// as we only use one single key, we just return it
// 		return hmacSecret, nil
// 	})
// 	if err != nil {

// 	}
// 	claim := parsedToken.Claims.(jwt.MapClaims)

// 	userID := claim["userId"]
// 	return userID.(string), nil
// }
