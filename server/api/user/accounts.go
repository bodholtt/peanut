package user

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"net/http"
	"peanutserver/auth"
	"peanutserver/database"
	"peanutserver/pcfg"
	"peanutserver/types"
)

// hashPassword - hash a password with SHA256 encryption
func hashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// HandleAccountsOPTIONS - OPTIONS for /login, /signup, /createUser
func HandleAccountsOPTIONS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", pcfg.Cfg.Client.Host)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
}

// HandleLogin - Accept a POST request with a username and password, validate, create a JWT and return it in the response.
func HandleLogin(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", pcfg.Cfg.Client.Host)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	user := &types.User{}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  nil,
			Error: err.Error(),
		})
		return
	}

	hashedPassword := hashPassword(user.Password)

	id, err := database.CheckAuthentication(user.Username, hashedPassword)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			//	User and password match not found - login failed
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(types.APIResponse{
				Body:  nil,
				Error: "invalid user credentials",
			})
			return
		}
		// any other error
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  nil,
			Error: err.Error(),
		})
		return
	}

	token, err := auth.CreateToken(user.Username, id, 168) // 7 days
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  nil,
			Error: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)

	b := struct {
		Token   string `json:"token"`
		Message string `json:"message"`
	}{
		Token:   token,
		Message: "Login successful",
	}

	json.NewEncoder(w).Encode(types.APIResponse{
		Body:  b,
		Error: "",
	})
}

// HandleSignup - Accept a POST request with a username and password and create it in the database.
// Respond with the user's ID if successful.
func HandleSignup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", pcfg.Cfg.Client.Host)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	user := &types.User{}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  "",
			Error: err.Error(),
		})
		return
	}

	user.Rank = pcfg.Perms.DefaultRank
	hashedPassword := hashPassword(user.Password)

	userID, err := database.CreateUser(user.Username, hashedPassword, user.Rank)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  "",
			Error: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(types.APIResponse{
		Body:  userID,
		Error: "",
	})
}

// HandleCreateUser - Accept a POST request with a username, password, and rank. For manual account creation by privileged users.
// Respond with the user's ID if successful.
func HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", pcfg.Cfg.Client.Host)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	user := &types.User{}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  nil,
			Error: err.Error(),
		})
		return
	}

	// a user cannot be rank 0 (anonymous)
	if user.Rank == 0 {
		user.Rank = pcfg.Perms.DefaultRank
	}

	hashedPassword := hashPassword(user.Password)
	userID, err := database.CreateUser(user.Username, hashedPassword, user.Rank)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  nil,
			Error: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(types.APIResponse{
		Body:  userID,
		Error: "",
	})
}
