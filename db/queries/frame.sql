-- name: GetFrame :one
select *
from frame
where id = $1 and user_id = $2;
-- name: GetFrames :many
select *
from frame where user_id = $1;
-- name: SaveFrame :one
insert into frame (
    id,
    title,
    description,
    created_at,
    user_id,
    frame_status,
    file_id
  )
values ($1, $2, $3, NOW(), $4, $5, $6) on conflict (id) DO
UPDATE
set title = $2,
  description = $3,
  user_id = $4,
  frame_status = $5,
  file_id = $6,
  modified_at = NOW()
RETURNING id;
-- name: DeleteFrame :exec
delete from frame where id = $1 and user_id = $2;
