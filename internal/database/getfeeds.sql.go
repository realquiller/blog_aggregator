// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: getfeeds.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const getFeeds = `-- name: GetFeeds :many

SELECT
    f.id,
    f.name,
    f.url,
    f.user_id,
    u.name as user_name
FROM feeds f
INNER JOIN users u ON f.user_id = u.id
`

type GetFeedsRow struct {
	ID       uuid.UUID
	Name     string
	Url      string
	UserID   uuid.UUID
	UserName string
}

func (q *Queries) GetFeeds(ctx context.Context) ([]GetFeedsRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedsRow
	for rows.Next() {
		var i GetFeedsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Url,
			&i.UserID,
			&i.UserName,
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
