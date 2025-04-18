// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: getpostsforuser.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const getPostsForUser = `-- name: GetPostsForUser :many

select p.id, p.created_at, p.updated_at, title, p.url, description, published_at, p.feed_id, f.id, f.created_at, f.updated_at, name, f.url, f.user_id, last_fetched_at, ff.id, ff.created_at, ff.updated_at, ff.user_id, ff.feed_id
from posts p
INNER JOIN feeds f ON p.feed_id = f.id
INNER JOIN feed_follows ff ON f.id = ff.feed_id
WHERE ff.user_id = $1
ORDER BY COALESCE(p.published_at, p.created_at) DESC
LIMIT $2
`

type GetPostsForUserParams struct {
	UserID uuid.UUID
	Limit  int32
}

type GetPostsForUserRow struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Title         string
	Url           string
	Description   string
	PublishedAt   sql.NullTime
	FeedID        uuid.UUID
	ID_2          uuid.UUID
	CreatedAt_2   time.Time
	UpdatedAt_2   time.Time
	Name          string
	Url_2         string
	UserID        uuid.UUID
	LastFetchedAt sql.NullTime
	ID_3          uuid.UUID
	CreatedAt_3   time.Time
	UpdatedAt_3   time.Time
	UserID_2      uuid.UUID
	FeedID_2      uuid.UUID
}

func (q *Queries) GetPostsForUser(ctx context.Context, arg GetPostsForUserParams) ([]GetPostsForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getPostsForUser, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostsForUserRow
	for rows.Next() {
		var i GetPostsForUserRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Url,
			&i.Description,
			&i.PublishedAt,
			&i.FeedID,
			&i.ID_2,
			&i.CreatedAt_2,
			&i.UpdatedAt_2,
			&i.Name,
			&i.Url_2,
			&i.UserID,
			&i.LastFetchedAt,
			&i.ID_3,
			&i.CreatedAt_3,
			&i.UpdatedAt_3,
			&i.UserID_2,
			&i.FeedID_2,
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
