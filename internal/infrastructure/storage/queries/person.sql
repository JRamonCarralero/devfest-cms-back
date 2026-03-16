-- name: CreatePerson :one
INSERT INTO persons (
    first_name, last_name, email, avatar_url, 
    github_user, linkedin_url, twitter_url, website_url,
    created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $9
)
RETURNING *;

-- name: GetPersonByID :one
SELECT * FROM persons
WHERE id = $1 LIMIT 1;

-- name: GetPersonByEmail :one
SELECT * FROM persons
WHERE email = $1 LIMIT 1;

-- name: ListPersons :many
SELECT * FROM persons
ORDER BY last_name ASC, first_name ASC;

-- name: ListPersonsPaged :many
-- Usamos un listado con búsqueda básica por nombre/email y paginación
SELECT * FROM persons
WHERE 
    (first_name ILIKE '%' || sqlc.arg('search')::text || '%' OR 
     last_name ILIKE '%' || sqlc.arg('search')::text || '%' OR 
     email ILIKE '%' || sqlc.arg('search')::text || '%')
ORDER BY last_name ASC, first_name ASC
LIMIT $1 OFFSET $2;

-- name: CountPersons :one
SELECT count(*) FROM persons
WHERE 
    (first_name ILIKE '%' || sqlc.arg('search')::text || '%' OR 
     last_name ILIKE '%' || sqlc.arg('search')::text || '%' OR 
     email ILIKE '%' || sqlc.arg('search')::text || '%');

-- name: UpdatePerson :one
UPDATE persons
SET 
    first_name = $2,
    last_name = $3,
    email = $4,
    avatar_url = $5,
    github_user = $6,
    linkedin_url = $7,
    twitter_url = $8,
    website_url = $9,
    updated_by = $10,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeletePerson :exec
DELETE FROM persons
WHERE id = $1;