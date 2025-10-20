package utils

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateAccessToken(tokenString string) (jwt.MapClaims, error){

	// parsing token to check if the alg is HS256/hmac

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error)  {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok{
			return nil, fmt.Errorf(":Unexpected signing method %v", token.Header["alg"])
		}
		// returning the token in byte for further verify
		return []byte(config.AppConfig.AccessSecretKey), nil
	})

	if err != nil{
		return nil, err
	}

	//  returning the payload as 'claims'

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid{
		return claims, nil
	}

	// if token is not valid
	return nil, fmt.Errorf("invalid token")
}