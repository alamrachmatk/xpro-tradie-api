package lib

import (
	"net/http"
	jwt "github.com/dgrijalva/jwt-go"
)

func ExtractClaims(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, nil)
	if err == nil {
		return nil, CustomError(http.StatusBadGateway, "error parse token")
	} 
	claims, _ := token.Claims.(jwt.MapClaims)

	return claims, nil
}
