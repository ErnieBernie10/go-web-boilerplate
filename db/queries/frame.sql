-- name: GetFrame :one
select *
from frame
where id = $1;
-- name: GetFrames :many
select *
from frame;
-- name: SaveFrame :one
insert into frame (id, title, description, created_at)
values ($1, $2, $3, NOW()) on conflict (id) DO
UPDATE
set title = $2,
  description = $3,
  modified_at = NOW()
RETURNING id;