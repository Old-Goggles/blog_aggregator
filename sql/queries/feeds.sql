-- name: CreateFeed :one
INSERT INTO feeds (name, url, user_id)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetFeed :one
SELECT * FROM feeds
WHERE url = $1;

-- name: DeleteFeed :exec
DELETE FROM feeds
WHERE id = $1;

-- name: GetFeeds :many
SELECT * FROM feeds
WHERE user_id = $1;

-- name: GetAllFeeds :many
SELECT * FROM feeds;