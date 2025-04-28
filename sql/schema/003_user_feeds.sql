-- +goose Up
CREATE TABLE user_feeds (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_feed UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE user_feeds;
