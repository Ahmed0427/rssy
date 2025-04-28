-- name: CreateUserFeed :one
WITH inserted_user_feed AS (
    INSERT INTO user_feeds (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)

SELECT 
    inserted_user_feed.*,
    users.name AS user_name,
    feeds.name AS feed_name
FROM inserted_user_feed 
INNER JOIN users ON users.id = inserted_user_feed.user_id
INNER JOIN feeds ON feeds.id = inserted_user_feed.feed_id;


-- name: GetUserFeedsForUser :many
SELECT 
    *,
    users.name AS user_name,
    feeds.name AS feed_name
FROM user_feeds
INNER JOIN users ON users.id = user_feeds.user_id
INNER JOIN feeds ON feeds.id = user_feeds.feed_id
WHERE users.name = $1;

-- name: DeleteUserFeedByUserAndURL :exec
DELETE FROM user_feeds
WHERE user_id = (
    SELECT id FROM users
    WHERE users.name = $1
    LIMIT 1
)
AND feed_id = (
    SELECT id FROM feeds
    WHERE feeds.url = $2
    LIMIT 1
);
