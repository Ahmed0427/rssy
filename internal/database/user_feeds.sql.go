// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: user_feeds.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createUserFeed = `-- name: CreateUserFeed :one
WITH inserted_user_feed AS (
    INSERT INTO user_feeds (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id, created_at, updated_at, user_id, feed_id
)

SELECT 
    inserted_user_feed.id, inserted_user_feed.created_at, inserted_user_feed.updated_at, inserted_user_feed.user_id, inserted_user_feed.feed_id,
    users.name AS user_name,
    feeds.name AS feed_name
FROM inserted_user_feed 
INNER JOIN users ON users.id = inserted_user_feed.user_id
INNER JOIN feeds ON feeds.id = inserted_user_feed.feed_id
`

type CreateUserFeedParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
}

type CreateUserFeedRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
	UserName  string
	FeedName  string
}

func (q *Queries) CreateUserFeed(ctx context.Context, arg CreateUserFeedParams) (CreateUserFeedRow, error) {
	row := q.db.QueryRowContext(ctx, createUserFeed,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.FeedID,
	)
	var i CreateUserFeedRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
		&i.UserName,
		&i.FeedName,
	)
	return i, err
}

const deleteUserFeedByUserAndURL = `-- name: DeleteUserFeedByUserAndURL :exec
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
)
`

type DeleteUserFeedByUserAndURLParams struct {
	Name string
	Url  string
}

func (q *Queries) DeleteUserFeedByUserAndURL(ctx context.Context, arg DeleteUserFeedByUserAndURLParams) error {
	_, err := q.db.ExecContext(ctx, deleteUserFeedByUserAndURL, arg.Name, arg.Url)
	return err
}

const getUserFeedsForUser = `-- name: GetUserFeedsForUser :many
SELECT 
    user_feeds.id, user_feeds.created_at, user_feeds.updated_at, user_id, feed_id, users.id, users.created_at, users.updated_at, users.name, feeds.id, feeds.created_at, feeds.updated_at, last_fetched_at, feeds.name, url,
    users.name AS user_name,
    feeds.name AS feed_name
FROM user_feeds
INNER JOIN users ON users.id = user_feeds.user_id
INNER JOIN feeds ON feeds.id = user_feeds.feed_id
WHERE users.name = $1
`

type GetUserFeedsForUserRow struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	UserID        uuid.UUID
	FeedID        uuid.UUID
	ID_2          uuid.UUID
	CreatedAt_2   time.Time
	UpdatedAt_2   time.Time
	Name          string
	ID_3          uuid.UUID
	CreatedAt_3   time.Time
	UpdatedAt_3   time.Time
	LastFetchedAt sql.NullTime
	Name_2        string
	Url           string
	UserName      string
	FeedName      string
}

func (q *Queries) GetUserFeedsForUser(ctx context.Context, name string) ([]GetUserFeedsForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserFeedsForUser, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserFeedsForUserRow
	for rows.Next() {
		var i GetUserFeedsForUserRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.FeedID,
			&i.ID_2,
			&i.CreatedAt_2,
			&i.UpdatedAt_2,
			&i.Name,
			&i.ID_3,
			&i.CreatedAt_3,
			&i.UpdatedAt_3,
			&i.LastFetchedAt,
			&i.Name_2,
			&i.Url,
			&i.UserName,
			&i.FeedName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
