-- name: CreateUserRole :one
INSERT INTO "user_role"(
    user_id,
    role_id,
    cr_user
) VALUES(
    $1,$2,$3
) RETURNING *;


-- name: GetUserRoleByUserId :many
SELECT * FROM "user_role"
WHERE user_id = sqlc.arg(user_id);


-- name: GetUserRoles :many
SELECT * FROM "user_role"
ORDER BY user_id
LIMIT $1
OFFSET $2;

-- name: UpdateUserRole :one
UPDATE "user_role"
SET 
    role_id = COALESCE(sqlc.narg(role_id),role_id),
    up_date = COALESCE(sqlc.narg(up_date),up_date),
    up_user = COALESCE(sqlc.narg(up_user),up_user)
WHERE user_id = sqlc.arg(user_id)
RETURNING *;