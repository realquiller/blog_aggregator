-- name: GetPostsForUser :many

select *
from posts p
INNER JOIN feeds f ON p.feed_id = f.id
INNER JOIN feed_follows ff ON f.id = ff.feed_id
WHERE ff.user_id = $1
ORDER BY COALESCE(p.published_at, p.created_at) DESC
LIMIT $2;