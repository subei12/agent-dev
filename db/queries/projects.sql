-- name: CreateProject :one
insert into projects (
  id,
  name,
  description,
  status
) values (
  $1, $2, $3, $4
)
returning id, name, description, status, created_at, updated_at;

-- name: GetProject :one
select id, name, description, status, created_at, updated_at
from projects
where id = $1;
