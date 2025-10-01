-- name: CreateUser :one
INSERT INTO
	users (
		"user_email",
		"user_name",
		"password",
		"profile_image",
		"created_at"
	)
VALUES
	($1, $2, $3, $4, $5)
RETURNING
	"user_id",
	"user_email",
	"user_name",
	"profile_image",
	"created_at";

-- name: UpdateUserInfo :one
UPDATE users
SET
	user_name = $1,
	password = $2,
	profile_image = $3
WHERE
	user_id = $4
RETURNING
	"user_id",
	"user_email",
	"user_name",
	"profile_image";

-- name: FindUserById :one
SELECT
	*
FROM
	users
WHERE
	user_id = $1;

-- name: FindUserInfoById :one
SELECT
	user_id,
	user_name,
	profile_image
FROM
	users
WHERE
	user_id = $1;

-- name: GetAllUserInfo :many
SELECT
	user_id,
	user_name,
	profile_image
FROM
	users
ORDER BY
	user_id ASC
LIMIT
	$1
OFFSET
	$2;

-- name: GetUserInfoByName :many
SELECT
	user_id,
	user_name,
	profile_image
FROM
	users
WHERE
	user_name LIKE ('%' || $1 || '%')
ORDER BY
	user_id
LIMIT
	$2
OFFSET
	$3;

-- name: FindUserByEmail :one
SELECT
	*
FROM
	users
WHERE
	user_email = $1;
