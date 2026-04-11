-- name: CreateScheduleEntry :one
INSERT INTO scheduler (
    track_id,
    talk_id,
    start_time,
    end_time,
    room
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetScheduleEntryByID :one
SELECT 
    sch.*, 
    tr.name as track_name,
    tr.event_date,
    t.title as talk_title,
    t.description as talk_description,
    (
        SELECT json_agg(json_build_object(
            'first_name', p.first_name,
            'last_name', p.last_name,
            'avatar_url', p.avatar_url
        ))
        FROM talk_speakers ts
        JOIN speakers s ON ts.speaker_id = s.id
        JOIN persons p ON s.person_id = p.id
        WHERE ts.talk_id = sch.talk_id
    ) as speakers
FROM scheduler sch
JOIN tracks tr ON sch.track_id = tr.id
LEFT JOIN talks t ON sch.talk_id = t.id
WHERE sch.id = $1 LIMIT 1;

-- name: ListScheduleByTrack :many
SELECT 
    sch.*, 
    t.title as talk_title,
    (
        SELECT json_agg(json_build_object(
            'first_name', p.first_name,
            'last_name', p.last_name
        ))
        FROM talk_speakers ts
        JOIN speakers s ON ts.speaker_id = s.id
        JOIN persons p ON s.person_id = p.id
        WHERE ts.talk_id = sch.talk_id
    ) as speakers
FROM scheduler sch
LEFT JOIN talks t ON sch.talk_id = t.id
WHERE sch.track_id = $1
ORDER BY sch.start_time ASC;

-- name: UpdateScheduleEntry :one
UPDATE scheduler
SET 
    track_id = COALESCE(sqlc.narg('track_id'), track_id),
    talk_id = COALESCE(sqlc.narg('talk_id'), talk_id),
    start_time = COALESCE(sqlc.narg('start_time'), start_time),
    end_time = COALESCE(sqlc.narg('end_time'), end_time),
    room = COALESCE(sqlc.narg('room'), room),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteScheduleEntry :exec
DELETE FROM scheduler
WHERE id = $1;