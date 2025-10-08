-- name: CreateFriendship :one
INSERT INTO
	friendships (user1_id, user2_id)
VALUES
	($1, $2)
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
		user1_id = sqlc.arg (user_id)
		AND user2_id = sqlc.arg (friend_id)
	)
	OR (
		user2_id = sqlc.arg (user_id)
		AND user1_id = sqlc.arg (friend_id)
	);

-- name: DeleteFriendshipById :one
DELETE FROM friendships
WHERE
	friendship_id = $1
RETURNING
	*;

-- name: DeleteFriendshipByUserId :one
DELETE FROM friendships
WHERE
	(
		user1_id = sqlc.arg (user_id)
		AND user2_id = sqlc.arg (friend_id)
	)
	OR (
		user2_id = sqlc.arg (user_id)
		AND user1_id = sqlc.arg (friend_id)
	)
RETURNING
	*;

-- name: GetFriendshipTable :many
SELECT
	*
FROM
	friendships;

-- name: GetFriendshipsByUserId :many
SELECT
	*
FROM
	friendships
WHERE
	user1_id = sqlc.arg (user_id)
	OR user2_id = sqlc.arg (user_id);

-- name: GetAllFriendInfoByUserId :many
SELECT
	u.user_id,
	u.user_name,
	u.profile_image
FROM
	(
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
		LIMIT
			$1
		OFFSET
			$2
	) fr
	LEFT JOIN users u ON fr.fr_id = u.user_id;

-- name: GetFriendInfoByUserAndFriendName :many
SELECT
	u.user_id,
	u.user_name,
	u.profile_image
FROM
	(
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
		LIMIT
			$1
		OFFSET
			$2
	) fr
	LEFT JOIN users u ON fr.fr_id = u.user_id
WHERE
	u.user_name LIKE ('%' || sqlc.arg (name) || '%');
