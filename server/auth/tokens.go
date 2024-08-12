package auth

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"peanutserver/database"
	"peanutserver/pcfg"
	"peanutserver/types"
	"strings"
	"time"
)

var secretKey = []byte(pcfg.Cfg.Server.SecretKey)

// userClaims - JWT claims allowing for registered claims + a user's username and id.
type userClaims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
	UserID   int    `json:"user_id"`
}

// CreateToken - Create a JWT for API auth for a user
// username - User's username
// expiration - Time until expiration in hours > 0
func CreateToken(username string, id int, expiryHours int) (string, error) {

	if expiryHours < 1 {
		return "", errors.New("expiryHours must be greater than zero")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		userClaims{
			Username: username,
			UserID:   id,
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryHours) * time.Hour)),
			},
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// verifyToken - verify that a token is valid and return the username in the token
func verifyToken(tokenString string) (username string, userID int, err error) {
	claims := &userClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		log.Println(err)
		return "", 0, err
	}

	if !token.Valid {
		return "", 0, errors.New("invalid token")
	}

	username = claims.Username
	userID = claims.UserID

	return username, userID, nil
}

// getTokenString - accept an http request and retrieve the jwt from the authorization header
func getTokenString(r *http.Request) (tokenString string, err error) {
	authheader := r.Header.Get("Authorization")
	// Authorization: Bearer <token>
	tokenString, found := strings.CutPrefix(authheader, "Bearer ")
	if !found {
		return "", errors.New("invalid authorization header")
	}
	return tokenString, nil
}

// GetUserIDFromAuthHeader - Public method to get the username from the auth header.
func GetUserIDFromAuthHeader(r *http.Request) (userID int, err error) {
	tokenString, err := getTokenString(r)
	if err != nil {
		return 0, err
	}
	_, userID, err = verifyToken(tokenString)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// RankMiddleware - Middleware to check if a user is allowed to access a resource according to their rank.
func RankMiddleware(next http.Handler, rank int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if rank > 0 {
			tokenString, err := getTokenString(r)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(types.APIResponse{
					Body:  nil,
					Error: err.Error(),
				})
				return
			}

			_, id, err := verifyToken(tokenString)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(types.APIResponse{
					Body:  nil,
					Error: err.Error(),
				})
				return
			}

			err = database.CheckUserRank(id, rank)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(types.APIResponse{
					Body:  nil,
					Error: err.Error(),
				})
				return
			}

		}

		next.ServeHTTP(w, r)
	})
}
