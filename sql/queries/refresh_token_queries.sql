-- name: CreateRefreshToken :one
INSERT INTO
	refresh_tokens (token, user_id, user_email, user_name, expired_at)
VALUES
	($1, $2, $3, $4, $5)
RETURNING
	*;

-- name: FindRefreshTokenByToken :one
SELECT
	*
FROM
	refresh_tokens
WHERE
	token = $1;

-- name: FindRefreshTokenByUserId :one
SELECT
	*
FROM
	refresh_tokens
WHERE
	user_id = $1;

-- name: UpdateTokenExpiratedTimeByUserId :one
UPDATE refresh_tokens
SET
	expired_at = $1
WHERE
	user_id = $2
RETURNING
	*;

-- name: UpdateTokenExpiratedTimeByToken :one
UPDATE refresh_tokens
SET
	expired_at = $1
WHERE
	token = $2
RETURNING
	*;

-- name: DeleteRefreshTokenByUserId :one
DELETE FROM refresh_tokens
WHERE
	user_id = $1
RETURNING
	*;

-- name: DeleteRefreshToken :one
DELETE FROM refresh_tokens
WHERE
	token = $1
RETURNING
	*;
