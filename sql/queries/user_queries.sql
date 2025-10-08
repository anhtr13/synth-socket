-- name: CreateUser :one
INSERT INTO
	users (user_email, user_name, password, profile_image)
VALUES
	($1, $2, $3, $4)
RETURNING
	user_id,
	user_email,
	user_name,
	profile_image,
	created_at;

-- name: UpdateUserInfo :one
UPDATE users
SET
	user_name = $1,
	password = $2,
	profile_image = $3
WHERE
	user_id = $4
RETURNING
	user_id,
	user_email,
	user_name,
	profile_image;

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
	u.*,
	CASE
		WHEN fr_id IS NOT NULL THEN TRUE
		ELSE FALSE
	END AS is_friend
FROM
	(
		SELECT
			user_id,
			user_name,
			profile_image
		FROM
			users
		ORDER BY
			user_name
		LIMIT
			$1
		OFFSET
			$2
	) u
	LEFT JOIN (
		SELECT
			CASE
				WHEN user1_id = sqlc.arg (user_id) THEN user2_id
				WHEN user2_id = sqlc.arg (user_id) THEN user1_id
			END AS fr_id
		FROM
			friendships
		WHERE
			user1_id = sqlc.arg (user_id)
			OR user2_id = sqlc.arg (user_id)
	) fr ON u.user_id = fr.fr_id;

-- name: GetAllUserInfoByName :many
SELECT
	u.*,
	CASE
		WHEN fr_id IS NOT NULL THEN TRUE
		ELSE FALSE
	END AS is_friend
FROM
	(
		SELECT
			user_id,
			user_name,
			profile_image
		FROM
			users
		WHERE
			user_name LIKE ('%' || sqlc.arg (name) || '%')
		ORDER BY
			user_name
		LIMIT
			$1
		OFFSET
			$2
	) u
	LEFT JOIN (
		SELECT
			CASE
				WHEN user1_id = sqlc.arg (user_id) THEN user2_id
				WHEN user2_id = sqlc.arg (user_id) THEN user1_id
			END AS fr_id
		FROM
			friendships
		WHERE
			user1_id = sqlc.arg (user_id)
			OR user2_id = sqlc.arg (user_id)
	) fr ON u.user_id = fr.fr_id;

-- name: FindUserByEmail :one
SELECT
	*
FROM
	users
WHERE
	user_email = $1;
