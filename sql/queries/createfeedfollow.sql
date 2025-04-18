-- name: CreateFeedFollow :one
WITH inserted AS (
  INSERT INTO feed_follows(id, created_at, updated_at, user_id, feed_id)
  VALUES ($1, $2, $3, $4, $5)
  RETURNING *
)
SELECT
  inserted.*,
  u.name AS user_name,
  f.name AS feed_name
FROM inserted
INNER JOIN users u ON inserted.user_id = u.id
INNER JOIN feeds f ON inserted.feed_id = f.id;







