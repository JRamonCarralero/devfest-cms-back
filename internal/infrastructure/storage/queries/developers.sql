-- name: CountDevelopersByEvent :one
SELECT COUNT(*) 
FROM developers d
JOIN persons p ON d.person_id = p.id
WHERE d.event_id = $1 
AND (
    p.first_name ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%' OR 
    p.last_name ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%' OR 
    p.email ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%'
);

-- name: CreateDeveloper :one
INSERT INTO developers (
  person_id, 
  event_id, 
  role_description, 
  created_by, 
  updated_by)
VALUES (
  $1, $2, $3, $4, $4
)
RETURNING *;

-- name: GetDeveloperByID :one
SELECT 
  d.*, 
  p.first_name, p.last_name, p.email, p.avatar_url, p.github_user, p.twitter_url, p.linkedin_url, p.website_url
FROM developers d
JOIN persons p ON d.person_id = p.id
WHERE d.id = $1 LIMIT 1;

-- name: ListDevelopersByEvent :many
SELECT 
  d.*, 
  p.first_name, p.last_name, p.email, p.avatar_url, p.github_user, p.twitter_url, p.linkedin_url, p.website_url
FROM developers d
JOIN persons p ON d.person_id = p.id
WHERE d.event_id = $1;

-- name: ListDevelopersByEventPaged :many
SELECT 
  d.*, 
  p.first_name, p.last_name, p.email, p.avatar_url, p.github_user, p.twitter_url, p.linkedin_url, p.website_url
FROM developers d
JOIN persons p ON d.person_id = p.id
WHERE d.event_id = $1 
AND (
    p.first_name ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%' OR 
    p.last_name ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%' OR 
    p.email ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%'
)
ORDER BY p.first_name ASC
LIMIT $2 OFFSET $3;

-- name: UpdateDeveloper :one
UPDATE developers 
SET 
  role_description = COALESCE(sqlc.narg('role_description'), role_description), 
  updated_at = NOW(), 
  updated_by = $2
WHERE id = $1 RETURNING *;

-- name: DeleteDeveloper :exec
DELETE FROM developers 
WHERE id = $1;

-- name: GetDeveloperByPersonAndEvent :one
SELECT id FROM developers 
WHERE person_id = $1 AND event_id = $2;