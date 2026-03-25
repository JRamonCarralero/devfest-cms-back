-- name: CountOrganizersByEvent :one
SELECT COUNT(*) 
FROM organizers c
JOIN persons p ON c.person_id = p.id
WHERE c.event_id = $1 
AND (
    p.first_name ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%' OR 
    p.last_name ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%' OR 
    p.email ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%'
);

-- name: CreateOrganizer :one
INSERT INTO organizers (
  person_id, 
  event_id, 
  company,
  role_description,
  created_by, 
  updated_by
)
VALUES (
  $1, $2, $3, $4, $5, $5
)
RETURNING *;

-- name: GetOrganizerByID :one
SELECT 
  c.*, 
  p.first_name, p.last_name, p.email, p.avatar_url, p.github_user, p.twitter_url, p.linkedin_url, p.website_url
FROM organizers c
JOIN persons p ON c.person_id = p.id
WHERE c.id = $1 LIMIT 1;

-- name: ListOrganizersByEvent :many
SELECT 
  c.*, 
  p.first_name, p.last_name, p.email, p.avatar_url, p.github_user, p.twitter_url, p.linkedin_url, p.website_url
FROM organizers c
JOIN persons p ON c.person_id = p.id
WHERE c.event_id = $1;

-- name: ListOrganizersByEventPaged :many
SELECT 
  c.*, 
  p.first_name, p.last_name, p.email, p.avatar_url, p.github_user, p.twitter_url, p.linkedin_url, p.website_url
FROM organizers c
JOIN persons p ON c.person_id = p.id
WHERE c.event_id = $1 
AND (
    p.first_name ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%' OR 
    p.last_name ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%' OR
    p.email ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%'
)
ORDER BY p.first_name ASC
LIMIT $2 OFFSET $3;

-- name: UpdateOrganizer :one
UPDATE organizers SET 
  company = COALESCE(sqlc.narg('company'), company),
  role_description = COALESCE(sqlc.narg('role_description'), role_description),
  updated_at = NOW(), 
  updated_by = $2
WHERE id = $1 
RETURNING *;

-- name: DeleteOrganizer :exec
DELETE FROM organizers 
WHERE id = $1;

-- name: GetOrganizerByPersonAndEvent :one
SELECT id FROM organizers 
WHERE person_id = $1 AND event_id = $2;