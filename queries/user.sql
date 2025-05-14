-- name: CreateUser :one
INSERT INTO user (pk_user_id, first_name, last_name)
VALUES (?, ?, ?)
RETURNING *;

-- name: ListUsers :many
SELECT *
FROM user;

-- name: GetUser :one
SELECT *
FROM user
WHERE pk_user_id = ?;

-- name: UpdateUser :one
UPDATE user
SET first_name = ?, last_name = ?
WHERE pk_user_id = ?
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM user
WHERE pk_user_id = ?;