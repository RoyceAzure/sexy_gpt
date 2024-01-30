-- name: CreateSession :one
INSERT INTO "session"(
    user_id,
    refresh_token,
    user_agent,
    client_ip,
    expired_at
) VALUES(
    $1,$2,$3,$4,$5
) RETURNING *;


-- name: GetSessionByUserId :one
SELECT * FROM "session"
WHERE 
user_id = sqlc.arg(user_id) LIMIT 1;



-- name: GetSessions :many
SELECT * FROM "session"
ORDER BY user_id
LIMIT $1
OFFSET $2;

-- name: UpdateSession :one
UPDATE "session"
SET 
    is_blocked = COALESCE(sqlc.narg(is_blocked),is_blocked)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteSession :one
Delete FROM "session"
WHERE id = sqlc.arg(id)
RETURNING *;

