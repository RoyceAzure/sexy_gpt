-- name: CreateUser :one
INSERT INTO "user"(
    user_name,
    email,
    hashed_password,
    sso_identifer,
    is_internal,
    cr_user
)  VALUES(
    $1,$2,$3,$4,$5,$6
)  RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM "user"
WHERE email = sqlc.arg(email) LIMIT 1;

-- name: GetUser :one
SELECT * FROM "user"
WHERE user_id = sqlc.arg(user_id) LIMIT 1;

-- name: GetUsers :many
SELECT * FROM "user"
ORDER BY user_id
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE "user"
SET 
    user_name = COALESCE(sqlc.narg(user_name),user_name),
    email = COALESCE(sqlc.narg(email),email),
    is_email_verified = COALESCE(sqlc.narg(is_email_verified),is_email_verified),
    hashed_password = COALESCE(sqlc.narg(hashed_password),hashed_password),
    password_changed_at = COALESCE(sqlc.narg(password_changed_at),password_changed_at),
    sso_identifer = COALESCE(sqlc.narg(sso_identifer),sso_identifer),
    is_internal = COALESCE(sqlc.narg(is_internal),is_internal),
    up_date = COALESCE(sqlc.narg(up_date),up_date),
    up_user = COALESCE(sqlc.narg(up_user),up_user)
WHERE user_id = sqlc.arg(user_id)
RETURNING *;
