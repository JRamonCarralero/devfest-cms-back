-- name: CreateTrack :one
INSERT INTO tracks (
    event_id,
    name,
    event_date
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetTrackByID :one
SELECT * FROM tracks
WHERE id = $1 LIMIT 1;

-- name: ListTracksByEvent :many
SELECT * FROM tracks
WHERE event_id = $1
ORDER BY event_date ASC, name ASC;

-- name: ListTracksByEventPaged :many
SELECT * FROM tracks
WHERE event_id = $1
AND (
    name ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%'
)
ORDER BY event_date ASC, name ASC
LIMIT $2 OFFSET $3;

-- name: CountTracksByEvent :one
SELECT COUNT(*) FROM tracks
WHERE event_id = $1
AND (
    name ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%'
);

-- name: UpdateTrack :one
UPDATE tracks
SET 
    name = COALESCE(sqlc.narg('name'), name),
    event_date = COALESCE(sqlc.narg('event_date'), event_date)
WHERE id = $1
RETURNING *;

-- name: DeleteTrack :exec
DELETE FROM tracks
WHERE id = $1;