-- name: GetFrame :one
select * from frame where id = $1;
