-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title,
    url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.*
FROM posts
INNER JOIN user_feeds
ON posts.feed_id = user_feeds.feed_id
WHERE user_feeds.user_id = $1
ORDER BY posts.published_at
LIMIT $2;
