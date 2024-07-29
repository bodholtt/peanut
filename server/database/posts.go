package database

import (
	"context"
	"log"
	"peanutserver/types"
	"strconv"
	"time"
)

func GetPostCount() (int, error) {
	db := Establish()
	defer db.Close()

	var count int
	err := db.QueryRow(context.Background(), "SELECT count(*) FROM posts").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetPostThumbs(limit int, offset int) (*types.PostThumbs, error) {
	db := Establish()
	defer db.Close()

	var id int
	var imagePath string

	rows, err := db.Query(context.Background(), "SELECT post_id, image_path FROM posts ORDER BY post_id LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	thumbs := types.PostThumbs{}

	for rows.Next() {
		err := rows.Scan(&id, &imagePath)
		if err != nil {
			return nil, err
		}
		thumbs.Thumbs = append(thumbs.Thumbs, types.PostThumb{
			ID:        strconv.Itoa(id),
			ImagePath: imagePath,
		})
	}

	return &thumbs, nil
}

func CreatePost() (int, error) {
	db := Establish()
	defer db.Close()

	var id int
	imagePath := "/static/images/test.jpg"

	err := db.QueryRow(
		context.Background(),
		"INSERT INTO posts (created_at, image_path, author_id) VALUES ($1, $2, $3) RETURNING post_id",
		time.Now(), imagePath, 1).
		Scan(&id)

	if err != nil {
		return 0, err
	}

	log.Println("Inserted new post at time :", time.Now())

	return id, nil
}

func GetPost(id int) (*types.Post, error) {
	db := Establish()
	defer db.Close()

	var tags string
	var createdAt time.Time
	var imagePath string
	var authorID int

	err := db.QueryRow(context.Background(),
		"SELECT tags, created_at, image_path, author_id FROM posts WHERE post_id = $1", id).
		Scan(&tags, &createdAt, &imagePath, &authorID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	post := &types.Post{
		ID:        strconv.Itoa(id),
		Tags:      tags,
		CreatedAt: createdAt.String(),
		ImagePath: imagePath,
		AuthorID:  strconv.Itoa(authorID),
	}

	return post, nil
}
