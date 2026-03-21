-- name: CreateSponsor :one
INSERT INTO sponsors (
    event_id,
    name,
    logo_url,
    website_url,
    tier,
    order_priority,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $7
)
RETURNING *;

-- name: GetSponsorByID :one
SELECT * FROM sponsors
WHERE id = $1 LIMIT 1;

-- name: ListSponsorsByEvent :many
SELECT * FROM sponsors
WHERE event_id = $1
ORDER BY order_priority DESC, name ASC;

-- name: ListSponsorsByEventPaged :many
SELECT * FROM sponsors
WHERE event_id = $1
AND (
    name ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%'
)
ORDER BY order_priority DESC, name ASC
LIMIT $2 OFFSET $3;

-- name: CountSponsorsByEvent :one
SELECT COUNT(*) FROM sponsors
WHERE event_id = $1
AND (
    name ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%'
);

-- name: UpdateSponsor :one
UPDATE sponsors
SET 
    name = COALESCE(sqlc.narg('name'), name),
    logo_url = COALESCE(sqlc.narg('logo_url'), logo_url),
    website_url = COALESCE(sqlc.narg('website_url'), website_url),
    tier = COALESCE(sqlc.narg('tier'), tier),
    order_priority = COALESCE(sqlc.narg('order_priority'), order_priority),
    updated_at = NOW(),
    updated_by = $2
WHERE id = $1
RETURNING *;

-- name: DeleteSponsor :exec
DELETE FROM sponsors
WHERE id = $1;