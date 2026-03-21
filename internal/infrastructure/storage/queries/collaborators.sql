-- name: CountCollaboratorsByEvent :one
SELECT COUNT(*) 
FROM collaborators c
JOIN persons p ON c.person_id = p.id
WHERE c.event_id = $1 
AND (
    p.first_name ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%' OR 
    p.last_name ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%' OR 
    p.email ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%'
);

-- name: CreateCollaborator :one
INSERT INTO collaborators (
  person_id, 
  event_id, 
  area, 
  created_by, 
  updated_by
)
VALUES (
  $1, $2, $3, $4, $4
)
RETURNING *;

-- name: GetCollaboratorByID :one
SELECT 
  c.*, 
  p.first_name, p.last_name, p.email, p.avatar_url, p.github_user, p.twitter_url, p.linkedin_url, p.website_url
FROM collaborators c
JOIN persons p ON c.person_id = p.id
WHERE c.id = $1 LIMIT 1;

-- name: ListCollaboratorsByEvent :many
SELECT 
  c.*, 
  p.first_name, p.last_name, p.email, p.avatar_url, p.github_user, p.twitter_url, p.linkedin_url, p.website_url
FROM collaborators c
JOIN persons p ON c.person_id = p.id
WHERE c.event_id = $1;

-- name: ListCollaboratorsByEventPaged :many
SELECT 
  c.*, 
  p.first_name, p.last_name, p.email, p.avatar_url, p.github_user, p.twitter_url, p.linkedin_url, p.website_url
FROM collaborators c
JOIN persons p ON c.person_id = p.id
WHERE c.event_id = $1 
AND (
    p.first_name ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%' OR 
    p.last_name ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%' OR
    p.email ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%'
)
ORDER BY p.first_name ASC
LIMIT $2 OFFSET $3;

-- name: UpdateCollaborator :one
UPDATE collaborators SET 
  area = COALESCE(sqlc.narg('area'), area), 
  updated_at = NOW(), 
  updated_by = $2
WHERE id = $1 
RETURNING *;

-- name: DeleteCollaborator :exec
DELETE FROM collaborators 
WHERE id = $1;

-- name: GetCollaboratorByPersonAndEvent :one
SELECT id FROM collaborators 
WHERE person_id = $1 AND event_id = $2;