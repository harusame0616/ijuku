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
