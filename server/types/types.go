package types

type APIResponse struct {
	Body  any    `json:"body"`
	Error string `json:"error"`
}

// Post - details of a particular post.
// ID: post id
// Tags: whitespace separated list of tags
// CreatedAt: Timestamp of post creation
// ImagePath: Path to static image file
// AuthorID: ID of uploader
// Source: Link to post source
// MD5Hash: Hash of the image for similarity comparison
// Previous: Previous post for pagination
// Next: Next post for pagination
// Query: The query to use for pagination
type Post struct {
	ID        string `json:"id"`
	Tags      []Tag  `json:"tags"`
	CreatedAt string `json:"created_at"`
	ImagePath string `json:"image_path"`
	AuthorID  string `json:"author_id"`
	Source    string `json:"source"`
	MD5Hash   string `json:"md5_hash"`
	Previous  string `json:"previous"` // Not stored in database - for pagination
	Next      string `json:"next"`     // see above
	Query     string `json:"query"`    // see above
}

// PostThumb - thumbnail of post viewed on the browsing page
// ID: post id
// ImagePath: Path to static image
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

type Tag struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
