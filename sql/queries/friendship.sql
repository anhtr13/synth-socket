-- name: CreateFriendship :one
INSERT INTO
	friendships (user1_id, user2_id, created_at)
VALUES
	($1, $2, $3)
RETURNING
	*;

-- name: FindFriendshipById :one
SELECT
	*
FROM
	friendships
WHERE
	friendship_id = $1;

-- name: FindFriendshipByUserIds :one
SELECT
	*
FROM
	friendships
WHERE
	(
		user1_id = $1
		AND user2_id = $2
	)
	OR (
		user2_id = $1
		AND user1_id = $2
	);

-- name: DeleteFriendshipById :one
DELETE FROM friendships
WHERE
	friendship_id = $1
RETURNING
	*;

-- name: DeleteFriendship :one
DELETE FROM friendships
WHERE
	(
		user1_id = $1
		AND user2_id = $2
	)
	OR (
		user2_id = $1
		AND user1_id = $2
	)
RETURNING
	*;

-- name: GetFriendshipTable :many
SELECT
	*
FROM
	friendships;

-- name: GetAllUserFriendIds :many
SELECT
	*
FROM
	friendships
WHERE
	user1_id = $1
	OR user2_id = $1;

-- name: GetAllUserFriendInfo :many
SELECT
	u.user_id,
	u.user_name,
	u.profile_image
FROM
	(
		SELECT
			CASE
				WHEN user1_id = $1 THEN user2_id
				WHEN user2_id = $1 THEN user1_id
			END AS fr_id
		FROM
			friendships
		WHERE
			user1_id = $1
			OR user2_id = $1
	) fr
	LEFT JOIN users u ON fr.fr_id = u.user_id
LIMIT
	$2
OFFSET
	$3;

-- name: GetUserFriendByName :many
SELECT
	u.user_id,
	u.user_name,
	u.profile_image
FROM
	(
		SELECT
			CASE
				WHEN user1_id = $1 THEN user2_id
				WHEN user2_id = $1 THEN user1_id
			END AS fr_id
		FROM
			friendships
		WHERE
			user1_id = $1
			OR user2_id = $1
	) fr
	LEFT JOIN users u ON fr.fr_id = u.user_id
WHERE
	u.user_name LIKE ('%' || $2 || '%')
LIMIT
	$3
OFFSET
	$4;
