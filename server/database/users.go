package database

import (
	"context"
	"peanutserver/types"
	"strconv"
	"time"
)

func GetUser(id int) (*types.User, error) {
	db := Establish()
	defer db.Close()

	var uid int
	var name string
	var rank int
	var createdAt time.Time

	err := db.QueryRow(context.Background(), "SELECT user_id, name, rank, created_at FROM users WHERE user_id = $1", id).Scan(&uid, &name, &rank, &createdAt)
	if err != nil {
		return nil, err
	}

	user := types.User{
		ID:        strconv.Itoa(uid),
		Username:  name,
		Rank:      strconv.Itoa(rank),
		CreatedAt: createdAt.String(),
	}

	return &user, nil
}
