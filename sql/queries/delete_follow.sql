-- name: DeleteFollow :exec
DELETE FROM feed_follows
WHERE user_id = (
    SELECT id from users where users.name = $1
) AND feed_id = (
    SELECT id from feeds where feeds.url = $2
);