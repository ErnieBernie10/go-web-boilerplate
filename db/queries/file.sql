-- name: GetFileByID :one
select * from file where id = $1;

-- name: CreateFile :exec
insert into file (id, file_name) values ($1, $2);

-- name: DeleteFile :exec
delete from file where id = $1;
