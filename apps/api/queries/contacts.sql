-- name: InsertContact :exec
INSERT INTO
    contacts (
        contact_id,
        name,
        email,
        phone,
        content,
        ip_address,
        user_agent
    )
VALUES
    (
        @contact_id,
        @name,
        @email,
        @phone,
        @content,
        @ip_address,
        @user_agent
    );
