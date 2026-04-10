-- name: CountTalksByEvent :one
SELECT COUNT(*) 
FROM talks t
WHERE t.event_id = $1 
AND (
    t.title ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%' OR 
    t.description ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%'
);

-- name: CreateTalk :one
INSERT INTO talks (
    event_id, 
    title, 
    description,
    tags,
    created_by, 
    updated_by
)
VALUES (
    $1, $2, $3, $4, $5, $5
)
RETURNING *;

-- name: GetTalkByID :one
SELECT 
    t.*,
    (
        SELECT json_agg(json_build_object(
            'id', s.id,
            'first_name', p.first_name,
            'last_name', p.last_name,
            'email', p.email,
            'avatar_url', p.avatar_url,
            'company', s.company,
            'bio', s.bio
        ))
        FROM talk_speakers ts
        JOIN speakers s ON ts.speaker_id = s.id
        JOIN persons p ON s.person_id = p.id
        WHERE ts.talk_id = t.id
    ) AS speakers
FROM talks t
WHERE t.id = $1 LIMIT 1;

-- name: ListTalksByEvent :many
SELECT 
    t.*,
    (
        SELECT json_agg(json_build_object(
            'id', s.id,
            'first_name', p.first_name,
            'last_name', p.last_name,
            'email', p.email,
            'avatar_url', p.avatar_url,
            'company', s.company,
            'bio', s.bio
        ))
        FROM talk_speakers ts
        JOIN speakers s ON ts.speaker_id = s.id
        JOIN persons p ON s.person_id = p.id
        WHERE ts.talk_id = t.id
    ) AS speakers
FROM talks t
WHERE t.event_id = $1
ORDER BY t.created_at DESC;

-- name: ListTalksByEventPaged :many
SELECT 
    t.*,
    (
        SELECT json_agg(json_build_object(
            'id', s.id,
            'first_name', p.first_name,
            'last_name', p.last_name,
            'email', p.email,
            'avatar_url', p.avatar_url,
            'company', s.company,
            'bio', s.bio
        ))
        FROM talk_speakers ts
        JOIN speakers s ON ts.speaker_id = s.id
        JOIN persons p ON s.person_id = p.id
        WHERE ts.talk_id = t.id
    ) AS speakers
FROM talks t
WHERE t.event_id = $1 
AND (
    t.title ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%' OR 
    t.description ILIKE '%' || COALESCE(sqlc.narg('search'), '') || '%'
)
ORDER BY t.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateTalk :one
UPDATE talks SET 
    title = COALESCE(sqlc.narg('title'), title),
    description = COALESCE(sqlc.narg('description'), description),
    tags = COALESCE(sqlc.narg('tags'), tags),
    updated_at = NOW(), 
    updated_by = $2
WHERE id = $1 
RETURNING *;

-- name: DeleteTalk :exec
DELETE FROM talks 
WHERE id = $1;

-- name: AddSpeakerToTalk :exec
INSERT INTO talk_speakers (talk_id, speaker_id, created_by)
VALUES ($1, $2, $3);

-- name: RemoveSpeakerFromTalk :exec
DELETE FROM talk_speakers 
WHERE talk_id = $1 AND speaker_id = $2;