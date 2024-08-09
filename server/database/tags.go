package database

import (
	"context"
	"errors"
	"log"
	"peanutserver/types"
	"strconv"
)

func GetTag(id int) (tag *types.Tag, err error) {
	db := Establish()
	defer db.Close()

	var tagName string
	var description string

	err = db.QueryRow(context.Background(),
		"SELECT name, description FROM tags WHERE tag_id = $1", id).
		Scan(&tagName, &description)

	if err != nil {
		return nil, err
	}

	tag = &types.Tag{
		ID:          strconv.Itoa(id),
		Name:        tagName,
		Description: description,
	}

	return tag, nil
}

func GetTagsByPostID(id int) (tags *[]types.Tag, err error) {
	db := Establish()
	defer db.Close()

	rows, err := db.Query(context.Background(), "SELECT tag_id FROM post_tags WHERE post_id = $1", id)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var postTags []types.Tag
	var tagID int
	var tagName string

	for rows.Next() {
		err := rows.Scan(&tagID)
		if err != nil {
			return nil, err
		}

		err = db.QueryRow(context.Background(),
			"SELECT name FROM tags WHERE tag_id = $1", tagID).Scan(&tagName)
		if err != nil {
			return nil, err
		}
		postTags = append(postTags, types.Tag{
			ID:   strconv.Itoa(tagID),
			Name: tagName,
		})
	}

	return &postTags, nil
}

func DeleteTag(id int) (err error) {
	db := Establish()
	defer db.Close()

	log.Println("Deleting tag", id)

	rows, err := db.Query(context.Background(), "DELETE FROM post_tags WHERE tag_id = $1", id)
	defer rows.Close()
	if err != nil {
		return err
	}

	rows, err = db.Query(context.Background(), "DELETE FROM tags WHERE tag_id = $1", id)
	defer rows.Close()
	if err != nil {
		return err
	}

	return nil
}

func CreateTag(name string) (id int, err error) {
	db := Establish()
	defer db.Close()

	rows, err := db.Query(context.Background(),
		"INSERT INTO tags (name) VALUES ($1) RETURNING tag_id",
		name)
	defer rows.Close()

	condition := rows.Next()
	if !condition {
		log.Println("Tag", name, "already exists")
		return 0, errors.New("tag already exists")
	}
	err = rows.Scan(&id)

	log.Println("Created new tag", name, "with id", id)

	return id, nil
}
