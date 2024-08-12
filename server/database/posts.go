package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"peanutserver/types"
	"strconv"
	"strings"
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

func GetPostThumbs(limit int, offset int, tags []string) (thumbs *types.PostThumbs, err error) {
	db := Establish()
	defer db.Close()

	var id int
	var rownum int // trash var so that scan works
	var imagePath string

	var rows pgx.Rows

	if len(tags) != 0 {

		// TODO: actually get the tag ids
		tagIDs := []string{"3", "4"}

		queryString :=
			`WITH ordered_posts AS (
				SELECT post_id, image_path, ROW_NUMBER() OVER (ORDER BY post_id) AS row_num
				FROM posts
			) 
			SELECT op.post_id, image_path, row_num FROM ordered_posts op
			INNER JOIN post_tags pt ON op.post_id = pt.post_id
			WHERE pt.tag_id IN (` + strings.Join(tagIDs, ",") + `)
			AND row_num <= (SELECT count(*) FROM ordered_posts) - ` + strconv.Itoa(offset) + `
			GROUP BY op.post_id, op.image_path, op.row_num HAVING COUNT(DISTINCT pt.tag_id) = ` + strconv.Itoa(len(tagIDs)) + `
			ORDER BY op.row_num DESC LIMIT ` + strconv.Itoa(limit)
		rows, err = db.Query(context.Background(), queryString)

	} else {
		rows, err = db.Query(context.Background(),
			`WITH ordered_posts AS (
				SELECT post_id, image_path, ROW_NUMBER() OVER (ORDER BY post_id) AS row_num
				FROM posts
			)
			SELECT post_id, image_path, row_num FROM ordered_posts
			WHERE row_num <= (SELECT count(*) FROM ordered_posts) - $1 
			ORDER BY row_num DESC LIMIT $2`,
			offset, limit)
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	thumbs = &types.PostThumbs{}

	for rows.Next() {

		err := rows.Scan(&id, &imagePath, &rownum)
		if err != nil {
			return nil, err
		}
		thumbs.Thumbs = append(thumbs.Thumbs, types.PostThumb{
			ID:        id,
			ImagePath: imagePath,
		})
	}

	return thumbs, nil
}

func CreatePost(extension string, uploaderID int) (postID int, imagePath string, err error) {
	db := Establish()
	defer db.Close()

	var id int

	rows, err := db.Query(context.Background(),
		"INSERT INTO posts (created_at, author_id) VALUES ($1, $2) RETURNING post_id",
		time.Now(), uploaderID)
	defer rows.Close()
	if err != nil {
		return 0, "", err
	}

	condition := rows.Next()
	if !condition {
		return 0, "", err
	}

	err = rows.Scan(&id)
	if err != nil {
		return 0, "", err
	}

	imagePath = strconv.Itoa(id) + extension

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

	var createdAt time.Time
	var imagePath string
	var authorID int
	var source string

	err := db.QueryRow(context.Background(),
		"SELECT created_at, image_path, author_id, source FROM posts WHERE post_id = $1", id).
		Scan(&createdAt, &imagePath, &authorID, &source)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	rows, err := db.Query(context.Background(),
		"SELECT tag_id FROM post_tags WHERE post_id = $1", id)
	defer rows.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var tags []types.Tag
	var tagID int
	var tagName string

	for rows.Next() {
		err = rows.Scan(&tagID)
		if err != nil {
			log.Println(err)
			break
		}

		err = db.QueryRow(context.Background(), "SELECT name FROM tags WHERE tag_id = $1", tagID).Scan(&tagName)
		if err != nil {
			log.Println(err)
			break
		}

		tags = append(tags, types.Tag{ID: tagID, Name: tagName})
	}

	post := &types.Post{
		ID:        id,
		Tags:      tags,
		CreatedAt: createdAt.String(),
		ImagePath: imagePath,
		AuthorID:  authorID,
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

func DeletePost(id int) (filepath string, err error) {
	db := Establish()
	defer db.Close()

	rows, err := db.Query(context.Background(), "DELETE FROM posts WHERE post_id = $1 RETURNING image_path", id)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return "", err
	}
	rows.Next()
	err = rows.Scan(&filepath)
	if err != nil {
		return "", err
	}

	log.Println("Deleted post", id, "at time", time.Now())
	return filepath, nil
}
