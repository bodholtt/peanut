CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(256) NOT NULL, -- sha256 hashed
    rank INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

INSERT INTO users (name, password, rank)
VALUES ('administrator', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8', 3)
ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS tags (
    tag_id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    category TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS posts (
    post_id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    image_path VARCHAR(255) NOT NULL DEFAULT '',
    author_id INT REFERENCES users(user_id),
    source VARCHAR(255) NOT NULL DEFAULT '',
    md5 VARCHAR(32) NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS post_tags (
    post_id INT REFERENCES posts(post_id),
    tag_id INTEGER REFERENCES tags(tag_id),
    PRIMARY KEY (post_id, tag_id)
)