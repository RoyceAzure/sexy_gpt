-- name: CreateRole :one
INSERT INTO "role"(
    role_name,
    cr_user
) VALUES(
    $1,$2
) RETURNING *;


-- name: GetRoleByRoleName :one
SELECT * FROM "role"
WHERE role_name = sqlc.arg(role_name) LIMIT 1;

-- name: GetRole :one
SELECT * FROM "role"
WHERE role_id = sqlc.arg(role_id) LIMIT 1;

-- name: GetRoles :many
SELECT * FROM "role"
ORDER BY role_id
LIMIT $1
OFFSET $2;

-- name: UpdateRole :one
UPDATE "role"
SET 
    role_name = COALESCE(sqlc.narg(role_name),role_name),
    is_enable = COALESCE(sqlc.narg(is_enable),is_enable),
    up_date = COALESCE(sqlc.narg(up_date),up_date),
    up_user = COALESCE(sqlc.narg(up_user),up_user)
WHERE role_id = sqlc.arg(role_id)
RETURNING *;