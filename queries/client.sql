-- name: CreateClient :one
INSERT INTO client (pk_client_id, name)
VALUES (?, ?)
RETURNING *;

-- name: ListClients :many
SELECT *
FROM client;

-- name: GetClient :one
SELECT *
FROM client
WHERE pk_client_id = ?;

-- name: UpdateClient :one
UPDATE client
SET name = ?
WHERE pk_client_id = ?
RETURNING *;

-- name: DeleteClient :exec
DELETE FROM client
WHERE pk_client_id = ?;
