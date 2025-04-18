-- name: GetFeeds :many

SELECT
    f.id,
    f.name,
    f.url,
    f.user_id,
    u.name as user_name
FROM feeds f
INNER JOIN users u ON f.user_id = u.id;