-- name: CreateVertifyEmail :one
INSERT INTO "vertify_email"(
    user_id,
    email,
    secret_code
) VALUES(
    $1,$2,$3
) RETURNING *;


-- name: GetVertifyEmailByUserIdAndCode :one
SELECT * FROM "vertify_email"
WHERE 
user_id = sqlc.arg(user_id)
AND secret_code = sqlc.arg(secret_code);



-- name: GetVertifyEmails :many
SELECT * FROM "vertify_email"
ORDER BY user_id
LIMIT $1
OFFSET $2;

-- name: UpdateVertifyEmail :one
UPDATE "vertify_email"
SET 
    is_used = COALESCE(sqlc.narg(is_used),is_used),
    used_date = COALESCE(sqlc.narg(used_date),used_date)
WHERE id = sqlc.arg(id)
RETURNING *;