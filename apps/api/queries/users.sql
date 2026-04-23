-- name: GetUser :one
SELECT
    user_id,
    nickname,
    introduce
FROM
    users
WHERE
    users.user_id = @UserID :: UUID;

-- name: UpdateUser :exec
UPDATE
    users
SET
    nickname = @Nickname,
    introduce = @Introduce
WHERE
    users.user_id = @UserID :: UUID;
