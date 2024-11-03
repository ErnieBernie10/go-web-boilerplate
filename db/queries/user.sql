-- name: GetUserByEmail :one
select *
from app_user
where email like $1;
-- name: Register :one
insert into app_user (email, password_hash)
values ($1, $2)
returning id;
