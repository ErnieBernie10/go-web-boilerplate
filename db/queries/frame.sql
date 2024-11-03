-- name: GetFrame :one
select *
from frame
  left join file f on f.id = frame.file_id
where frame.id = $1
  and user_id = $2;

-- name: GetFrames :many
select *
from frame
    left join file f on f.id = frame.file_id
where user_id = $1;

-- name: SaveFrame :one
insert into frame (
    id,
    title,
    description,
    created_at,
    user_id,
    frame_status,
    file_id,
    content_type,
    content
  )
values ($1, $2, $3, NOW(), $4, $5, $6, $7, $8) on conflict (id) DO
UPDATE
set title = $2,
  description = $3,
  user_id = $4,
  frame_status = $5,
  file_id = $6,
  modified_at = NOW(),
  content_type = $7,
  content = $8
RETURNING id;

-- name: DeleteFrame :exec
delete from frame
where id = $1
  and user_id = $2;
