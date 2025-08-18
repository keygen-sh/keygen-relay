-- name: CreatePool :one
INSERT INTO pools (name)
VALUES (?)
RETURNING *;

-- name: GetPoolByID :one
SELECT *
FROM pools
WHERE id = ?;

-- name: GetPoolByName :one
SELECT *
FROM pools
WHERE name = ?;

-- name: DeletePoolByID :one
DELETE FROM pools
WHERE id = ?
RETURNING *;
