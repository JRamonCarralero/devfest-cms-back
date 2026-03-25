-- name: CountSpeakersByEvent :one
SELECT COUNT(*) 
FROM speakers s
JOIN persons p ON s.person_id = p.id
WHERE s.event_id = $1 
AND (
    p.first_name ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%' OR 
    p.last_name ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%' OR 
    p.email ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%'
);

-- name: CreateSpeaker :one
INSERT INTO speakers (
    person_id,
    event_id,
    bio,
    company,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $5
)
RETURNING *;

-- name: GetSpeakerByID :one
SELECT 
    s.*, 
    p.first_name, p.last_name, p.email, p.avatar_url, p.github_user, p.twitter_url, p.linkedin_url, p.website_url
FROM speakers s
JOIN persons p ON s.person_id = p.id
WHERE s.id = $1 LIMIT 1;

-- name: ListSpeakersByEvent :many
SELECT 
    s.*, 
    p.first_name, p.last_name, p.email, p.avatar_url, p.github_user, p.twitter_url, p.linkedin_url, p.website_url
FROM speakers s
JOIN persons p ON s.person_id = p.id
WHERE s.event_id = $1;

-- name: ListSpeakersByEventPaged :many
SELECT 
    s.*, 
    p.first_name, p.last_name, p.email, p.avatar_url, p.github_user, p.twitter_url, p.linkedin_url, p.website_url
FROM speakers s
JOIN persons p ON s.person_id = p.id
WHERE s.event_id = $1 
AND (
    p.first_name ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%' OR 
    p.last_name ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%' OR 
    p.email ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%'
)
ORDER BY p.first_name ASC
LIMIT $2 OFFSET $3;

-- name: UpdateSpeaker :one
UPDATE speakers
SET 
    bio = COALESCE(sqlc.narg('bio'), bio),
    company = COALESCE(sqlc.narg('company'), company),
    updated_at = NOW(),
    updated_by = $2
WHERE id = $1
RETURNING *;

-- name: DeleteSpeaker :exec
DELETE FROM speakers
WHERE id = $1;

-- name: GetSpeakerByPersonAndEvent :one
SELECT id FROM speakers 
WHERE person_id = $1 AND event_id = $2;