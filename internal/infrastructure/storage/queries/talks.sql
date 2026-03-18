-- name: CreateTalk :one
INSERT INTO talks (
    event_id,
    speaker_id,
    title,
    description,
    tags,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $6
)
RETURNING *;

-- name: GetTalkByID :one
SELECT 
    t.*, 
    s.bio, s.company,
    p.first_name, p.last_name, p.avatar_url
FROM talks t
JOIN speakers s ON t.speaker_id = s.id
JOIN persons p ON s.person_id = p.id
WHERE t.id = $1 LIMIT 1;

-- name: ListTalksByEvent :many
SELECT 
    t.*, 
    s.bio, s.company,
    p.first_name, p.last_name, p.avatar_url
FROM talks t
JOIN speakers s ON t.speaker_id = s.id
JOIN persons p ON s.person_id = p.id
WHERE t.event_id = $1
ORDER BY t.created_at DESC;

-- name: ListTalksByEventPaged :many
SELECT 
    t.*, 
    s.bio, s.company,
    p.first_name, p.last_name, p.avatar_url
FROM talks t
JOIN speakers s ON t.speaker_id = s.id
JOIN persons p ON s.person_id = p.id
WHERE t.event_id = $1 
AND t.title ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%' OR 
    p.first_name ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%' OR 
    p.last_name ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%'
)
ORDER BY t.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountTalksByEvent :one
SELECT COUNT(*) 
FROM talks t
JOIN speakers s ON t.speaker_id = s.id
JOIN persons p ON s.person_id = p.id
WHERE t.event_id = $1 
AND (
    t.title ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%' OR 
    p.first_name ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%' OR 
    p.last_name ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%'
);

-- name: UpdateTalk :one
UPDATE talks
SET 
    title = COALESCE(sqlc.narg('title'), title),
    description = COALESCE(sqlc.narg('description'), description),
    tags = COALESCE(sqlc.narg('tags'), tags),
    speaker_id = COALESCE(sqlc.narg('speaker_id'), speaker_id),
    updated_at = NOW(),
    updated_by = $2
WHERE id = $1
RETURNING *;

-- name: DeleteTalk :exec
DELETE FROM talks
WHERE id = $1;