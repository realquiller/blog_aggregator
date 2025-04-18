-- name: GetNextFeedToFetch :one

SELECT feeds.*
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;