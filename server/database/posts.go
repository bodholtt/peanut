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

	rows, err := db.Query(context.Background(),
		"SELECT post_id, image_path FROM posts WHERE post_id < (SELECT count(*) FROM posts) - $1 ORDER BY post_id DESC LIMIT $2", offset, limit)
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

func CreatePost(extension string) (postID int, imagePath string, err error) {
	db := Establish()
	defer db.Close()

	var id int

	rows, err := db.Query(context.Background(),
		"INSERT INTO posts (created_at, author_id) VALUES ($1, $2) RETURNING post_id",
		time.Now(), 1)
	defer rows.Close()
	condition := rows.Next()
	if !condition {
		return 0, "", err
	}
	err = rows.Scan(&id)

	imagePath = "/static/images/" + strconv.Itoa(id) + extension

	rows2, err := db.Query(context.Background(), "UPDATE posts SET image_path = $1 WHERE post_id = $2", imagePath, id)
	defer rows2.Close()

	if err != nil {
		return 0, "", err
	}

	log.Println("Inserted new post", id, "at time", time.Now())

	return id, imagePath, nil
}

func GetPost(id int) (*types.Post, error) {
	db := Establish()
	defer db.Close()

	var tags string
	var createdAt time.Time
	var imagePath string
	var authorID int
	var source string

	err := db.QueryRow(context.Background(),
		"SELECT tags, created_at, image_path, author_id, source FROM posts WHERE post_id = $1", id).
		Scan(&tags, &createdAt, &imagePath, &authorID, &source)

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
		Source:    source,
	}

	return post, nil
}

func GetNextAndPreviousPostIDs(id int, query string) (prev int, next int) {
	db := Establish()
	defer db.Close()

	// TODO: Implement doing the search query here

	var prevID int
	var nextID int

	err := db.QueryRow(context.Background(),
		"SELECT post_id FROM posts WHERE post_id = (select max(post_id) from posts where post_id < $1)", id).
		Scan(&prevID)
	if err != nil {
		prevID = 0
	}

	err = db.QueryRow(context.Background(),
		"SELECT post_id FROM posts WHERE post_id = (SELECT min(post_id) FROM posts WHERE post_id > $1)", id).
		Scan(&nextID)
	if err != nil {
		nextID = 0
	}

	return prevID, nextID
}

func DeletePost(id int) error {
	db := Establish()
	defer db.Close()

	_, err := db.Query(context.Background(), "DELETE FROM posts WHERE post_id = $1", id)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
