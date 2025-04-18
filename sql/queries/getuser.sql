-- name: GetUser :one

SELECT *
FROM users u
WHERE u.name = $1;