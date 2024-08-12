package database

import (
	"context"
	"errors"
	"peanutserver/types"
	"time"
)

func GetUser(id int) (*types.User, error) {
	db := Establish()
	defer db.Close()

	var uid int
	var name string
	var rank int
	var createdAt time.Time

	err := db.QueryRow(context.Background(),
		"SELECT user_id, name, rank, created_at FROM users WHERE user_id = $1", id).
		Scan(&uid, &name, &rank, &createdAt)
	if err != nil {
		return nil, err
	}

	user := types.User{
		ID:        uid,
		Username:  name,
		Rank:      rank,
		CreatedAt: createdAt.String(),
	}

	return &user, nil
}

// CreateUser - create a new user with the given string, sha256 hashed password, and rank.
// Return the ID of the new user or an error.
func CreateUser(name string, hashedPassword string, rank int) (id int, err error) {
	db := Establish()
	defer db.Close()

	// Create new, return id
	rows, err := db.Query(context.Background(),
		"INSERT INTO users (name, password, rank, created_at) VALUES ($1, $2, $3, $4) RETURNING user_id",
		name, hashedPassword, rank, time.Now())
	defer rows.Close()
	if err != nil {
		return 0, err
	}

	condition := rows.Next()
	if !condition {
		return 0, errors.New("user already exists")
	}

	err = rows.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// CheckAuthentication - verify if a username and hashed password combo aligns with a db entry.
func CheckAuthentication(name string, hashedPassword string) (id int, err error) {
	db := Establish()
	defer db.Close()

	var userID int
	err = db.QueryRow(context.Background(), "SELECT user_id FROM users WHERE name = $1 AND password = $2",
		name, hashedPassword).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// CheckUserRank - check if user rank is geq requiredRank, error if it is not.
func CheckUserRank(id int, requiredRank int) (err error) {
	db := Establish()
	defer db.Close()

	var userRank int

	err = db.QueryRow(context.Background(), "SELECT rank FROM users WHERE user_id = $1",
		id).Scan(&userRank)
	if err != nil {
		return err
	}
	if userRank < requiredRank {
		return errors.New("user is not the correct rank")
	}

	return nil
}
