-- name: CreateMetadata :one
INSERT INTO metadata (fk_record_id, fk_created_by_id)
VALUES (?, ?)
RETURNING *;

-- name: ListMetadata :many
SELECT *
FROM metadata;

-- name: GetMetadata :one
SELECT *
FROM metadata
WHERE fk_record_id = ?;

-- name: UpdateMetadata :one
UPDATE metadata
SET updated_at = ?, fk_updated_by_id = ?
WHERE fk_record_id = ?
RETURNING *;

-- name: DeleteMetadata :exec
DELETE FROM metadata
WHERE fk_record_id = ?;