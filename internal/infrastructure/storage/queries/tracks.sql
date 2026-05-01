-- name: CreateTrack :one
INSERT INTO tracks (
    event_id,
    name,
    event_date,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $4
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

-- name: GetFullEventSchedule :many
SELECT 
    tr.id AS track_id,
    tr.name AS track_name,
    tr.event_date,
    (
        SELECT json_agg(json_build_object(
            'id', sch.id,
            'start_time', sch.start_time,
            'end_time', sch.end_time,
            'room', sch.room,
            'talk', CASE WHEN t.id IS NOT NULL THEN json_build_object(
                'id', t.id,
                'title', t.title,
                'description', t.description,
                'speakers', (
                    SELECT json_agg(json_build_object(
                        'id', s.id,
                        'first_name', p.first_name,
                        'last_name', p.last_name,
                        'avatar_url', p.avatar_url,
                        'company', s.company
                    ))
                    FROM talk_speakers ts
                    JOIN speakers s ON ts.speaker_id = s.id
                    JOIN persons p ON s.person_id = p.id
                    WHERE ts.talk_id = t.id
                )
            ) ELSE NULL END
        ) ORDER BY sch.start_time ASC)
        FROM scheduler sch
        LEFT JOIN talks t ON sch.talk_id = t.id
        WHERE sch.track_id = tr.id
    ) AS entries
FROM tracks tr
WHERE tr.event_id = $1
ORDER BY tr.event_date ASC, tr.name ASC;

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
    event_date = COALESCE(sqlc.narg('event_date'), event_date),
    updated_at = NOW(),
    updated_by = $2
WHERE id = $1
RETURNING *;

-- name: DeleteTrack :exec
DELETE FROM tracks
WHERE id = $1;