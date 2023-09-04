-- name: GetUser :one
SELECT * FROM user
WHERE email = ?
LIMIT 1;

-- name: CreateUser :one
INSERT INTO user (fullname, email, password_hash)
VALUES (?, ?, ?)
RETURNING id;

-- name: GetBoards :many
SELECT id, name, slug 
FROM board
WHERE user_id = ?
ORDER BY id;

-- name: GetBoard :one
SELECT * from board
WHERE slug = ?
LIMIT 1;

-- name: CreateBoard :exec
INSERT INTO board (name, slug, user_id)
VALUES (?, ?, ?);

-- name: UpdateBoard :exec
UPDATE board
SET name = ?, slug = ?
WHERE slug = ? AND user_id = ?;

-- name: DeleteBoard :exec
DELETE FROM board
WHERE slug = ? AND user_id = ?;

-- name: CountColumns :one
SELECT count(*)
FROM column
WHERE board_id = ?;

-- name: GetColumns :many
SELECT id, name, board_id
FROM column
WHERE board_id = ?
ORDER BY element_order;

-- name: GetColumn :one
SELECT id, name, board_id, element_order
FROM column
WHERE id = ?
LIMIT 1;

-- name: CreateColumn :one
INSERT INTO column (name, element_order, board_id)
VALUES (?, ?, ?)
RETURNING id;

-- name: UpdateColumn :exec
UPDATE column
SET name = ?, element_order = ?
WHERE id = ?;

-- name: SetColumnOrder :exec
UPDATE column
SET element_order = ?
WHERE id = ?;

-- name: DeleteColumn :exec
DELETE FROM column
WHERE id = ?;

-- name: CountItems :one
SELECT count(*)
FROM item
WHERE column_id = ?;

-- name: GetItems :many
SELECT id, name, column_id
FROM item
WHERE column_id = ?
ORDER BY element_order;

-- name: CreateItem :one
INSERT INTO item (name, element_order, column_id)
VALUES (?, ?, ?)
RETURNING id;

-- name: UpdateItem :exec
UPDATE item
SET name = ? 
WHERE id = ?;

-- name: SetItemOrder :exec
UPDATE item
SET element_order = ?, column_id = ?
WHERE id = ?;

-- name: DeleteItem :exec
DELETE FROM item
WHERE id = ?;