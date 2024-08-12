package database

import (
	"context"
	"errors"
	"log"
	"peanutserver/types"
)

// GetTagByID - Get a tag by its ID.
func GetTagByID(id int) (tag *types.Tag, err error) {
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
		ID:          id,
		Name:        tagName,
		Description: description,
	}

	return tag, nil
}

// GetTagByName - Get a tag by its name.
func GetTagByName(name string) (tag *types.Tag, err error) {
	db := Establish()
	defer db.Close()

	var tagID int
	var description string

	err = db.QueryRow(context.Background(),
		"SELECT tag_id, description FROM tags WHERE name = $1", name).
		Scan(&tagID, &description)

	if err != nil {
		return nil, err
	}

	tag = &types.Tag{
		ID:          tagID,
		Name:        name,
		Description: description,
	}

	return tag, nil
}

// GetTagsByPostID - Get a list of tags from a post by the post's ID.
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
			ID:   tagID,
			Name: tagName,
		})
	}

	return &postTags, nil
}

// DeleteTag - Delete tag by ID.
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

// CreateTag - Create a tag by its name and return its ID.
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

// SetPostTags - Accept a post ID and a list of tag names, and link all the tags to the post.
// If a tag doesn't exist, this function will call CreateTag to create it.
func SetPostTags(postID int, tags []string) (err error) {
	db := Establish()
	defer db.Close()

	// this kind of sucks and i would like to do it in a better way.

	var tagsToBeSet []int

	// clear post's current tags so that we don't have to handle tag removals
	_, err = db.Query(context.Background(), "DELETE FROM post_tags WHERE post_id = $1", postID)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		retrievedTag, err := GetTagByName(tag)
		if err != nil {
			createdTagID, err := CreateTag(tag)
			if err != nil {
				continue
			}
			tagsToBeSet = append(tagsToBeSet, createdTagID)
			continue
		}
		tagsToBeSet = append(tagsToBeSet, retrievedTag.ID)
	}

	for _, id := range tagsToBeSet {
		db.QueryRow(context.Background(), "INSERT INTO post_tags (post_id, tag_id) VALUES ($1, $2)", postID, id)
	}

	log.Println("Tags updated for post", postID)
	return nil
}
