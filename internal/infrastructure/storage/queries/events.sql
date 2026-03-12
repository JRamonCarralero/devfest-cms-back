-- queries/events.sql

-- name: GetEventByID :one
SELECT * FROM events
WHERE id = $1 LIMIT 1;

-- name: GetEventBySlug :one
SELECT * FROM events
WHERE slug = $1 LIMIT 1;

-- name: ListEvents :many
SELECT * FROM events
ORDER BY created_at DESC;

-- name: ListEventsPaged :many
SELECT * FROM events
WHERE 
    (name ILIKE '%' || $1::text || '%' OR $1::text = '')
ORDER BY 
    CASE WHEN $4::text = 'name_asc' THEN name END ASC,
    CASE WHEN $4::text = 'name_desc' THEN name END DESC,
    CASE WHEN $4::text = 'created_at_asc' THEN created_at END ASC,
    CASE WHEN $4::text = 'created_at_desc' THEN created_at END DESC
LIMIT $2 OFFSET $3;

-- name: CountEvents :one
SELECT COUNT(*) FROM events
WHERE 
    (name ILIKE '%' || $1::text || '%' OR $1::text = '');

-- name: ListActiveEvents :many
SELECT * FROM events
WHERE is_active = true
ORDER BY created_at DESC;

-- name: CreateEvent :one
INSERT INTO events (
    name, 
    slug, 
    is_active, 
    created_by, 
    updated_by
) VALUES (
    $1, $2, $3, $4, $4
)
RETURNING *;

-- name: UpdateEvent :one
UPDATE events
SET name = $2,
    slug = $3,
    is_active = $4,
    updated_by = $5
WHERE id = $1
RETURNING *;

-- name: DeleteEvent :exec
DELETE FROM events
WHERE id = $1;