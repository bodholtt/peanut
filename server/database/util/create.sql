CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(256) NOT NULL, -- sha256 hashed
    rank INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL
);

INSERT INTO users VALUES (1, 'administrator', '6b3a55e0261b0304143f805a24924d0c1c44524821305f31d9277843b8a10f4e', 3) ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS tags (
    tag_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL DEFAULT '',
    category VARCHAR(255) NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS posts (
    post_id SERIAL PRIMARY KEY,
    tags VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL,
    image_path VARCHAR(255) NOT NULL,
    author_id INT REFERENCES users(user_id),
    source VARCHAR(255) NOT NULL DEFAULT ''
);

