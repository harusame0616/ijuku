-- name: CountApiKeyByUserID :one
SELECT
    COUNT(*)
FROM
    apikeys
WHERE
    apikeys.user_id = @UserID :: uuid;

-- name: InsertApiKey :exec
INSERT INTO
    apikeys (
        apikey_id,
        key_hash,
        user_id,
        plain_suffix,
        expired_at
    )
VALUES
    (
        @apikey_id,
        @key_hash,
        @user_id,
        @key_plain_suffix,
        @expired_at
    );

-- name: ListApiKeysByUserID :many
SELECT
    apikey_id,
    plain_suffix,
    _created_at,
    expired_at
FROM
    apikeys
WHERE
    apikeys.user_id = @UserID :: uuid;
