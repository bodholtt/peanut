package types

// Post - details of a particular post.
// ID: post id
// Tags: whitespace separated list of tags
// CreatedAt: Timestamp of post creation
type Post struct {
	ID        string `json:"id"`
	Tags      string `json:"tags"`
	CreatedAt string `json:"created_at"`
	ImagePath string `json:"image_path"`
	AuthorID  string `json:"author_id"`
}

// PostThumb - thumbnail of post viewed on the browsing page
// ID: post id
type PostThumb struct {
	ID        string `json:"id"`
	ImagePath string `json:"image_path"`
}

// PostThumbs - list of post thumbs
type PostThumbs struct {
	Thumbs []PostThumb
}

// User - user information
type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Rank      string `json:"rank"`
	CreatedAt string `json:"created_at"`
}
