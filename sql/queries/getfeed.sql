-- name: GetFeed :one

SELECT *
FROM feeds f
WHERE f.url = $1;