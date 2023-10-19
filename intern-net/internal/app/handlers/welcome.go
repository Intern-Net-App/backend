package handlers

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	// Obtain the session token from the requests cookies, which come with every request
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the JWT string from cookie
	tokenString := cookie.Value

	//Initialize a new instance of Claims
	claims := &Claims{}

	// Parse the JWT string and store the result in claims.
	// Note: Passing the key in method as well, casuing the method to return an error
	// if the token is invalid (if it expired according to expiring time set on login),
	// or if the signature does not match
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Finally, return welcome message to the user, along with their username given in token
	w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
}
