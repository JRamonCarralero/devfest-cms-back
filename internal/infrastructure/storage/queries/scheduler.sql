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
    t.title as talk_title,
    p.first_name as speaker_name,
    p.last_name as speaker_last_name,
    tr.name as track_name,
    tr.event_date
FROM scheduler sch
LEFT JOIN talks t ON sch.talk_id = t.id
LEFT JOIN speakers s ON t.speaker_id = s.id
LEFT JOIN persons p ON s.person_id = p.id
JOIN tracks tr ON sch.track_id = tr.id
WHERE sch.id = $1 LIMIT 1;

-- name: ListScheduleByTrack :many
SELECT 
    sch.*, 
    t.title as talk_title,
    p.first_name as speaker_name,
    p.last_name as speaker_last_name
FROM scheduler sch
LEFT JOIN talks t ON sch.talk_id = t.id
LEFT JOIN speakers s ON t.speaker_id = s.id
LEFT JOIN persons p ON s.person_id = p.id
WHERE sch.track_id = $1
ORDER BY sch.start_time ASC;

-- name: ListFullEventSchedule :many
SELECT 
    sch.*, 
    tr.name as track_name,
    tr.event_date,
    t.title as talk_title,
    p.first_name as speaker_name,
    p.last_name as speaker_last_name
FROM scheduler sch
JOIN tracks tr ON sch.track_id = tr.id
LEFT JOIN talks t ON sch.talk_id = t.id
LEFT JOIN speakers s ON t.speaker_id = s.id
LEFT JOIN persons p ON s.person_id = p.id
WHERE tr.event_id = $1
ORDER BY tr.event_date ASC, sch.start_time ASC;

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