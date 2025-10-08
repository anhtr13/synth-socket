-- name: CreateRoomMember :one
INSERT INTO
	room_members ("room_id", "member_id")
VALUES
	($1, $2)
RETURNING
	*;

-- name: GetRoomMembersTable :many
SELECT
	*
FROM
	room_members;

-- name: FindRoomMember :one
SELECT
	*
FROM
	room_members
WHERE
	room_id = $1
	AND member_id = $2;

-- name: CountMemmbersInRoom :one
SELECT
	count(member_id)
FROM
	room_members
WHERE
	room_id = $1;

-- name: DeleteRoomMemmber :exec
DELETE FROM room_members
WHERE
	room_id = $1
	AND member_id = $2;

-- name: GetAllMemberInfoByRoomId :many
SELECT
	u.user_id,
	u.user_name,
	u.profile_image,
	r.joined_at
FROM
	(
		SELECT
			*
		FROM
			room_members
		WHERE
			room_id = $1
	) r
	LEFT JOIN users u ON r.member_id = u.user_id
ORDER BY
	r.joined_at
LIMIT
	$2
OFFSET
	$3;

-- name: GetRoomMemberInfoByRoomId :many
SELECT
	u.user_id,
	u.user_name,
	u.profile_image,
	r.joined_at
FROM
	(
		SELECT
			*
		FROM
			room_members
		WHERE
			room_id = $1
	) r
	LEFT JOIN users u ON r.member_id = u.user_id
ORDER BY
	r.joined_at
LIMIT
	$2
OFFSET
	$3;

-- name: GetRoomMemberInfoByRoomAndUserName :many
SELECT
	u.user_id,
	u.user_name,
	u.profile_image,
	r.joined_at
FROM
	(
		SELECT
			*
		FROM
			room_members
		WHERE
			room_id = $1
		ORDER BY
			joined_at
		LIMIT
			$2
		OFFSET
			$3
	) r
	LEFT JOIN users u ON r.member_id = u.user_id
WHERE
	u.user_name LIKE ('%' || sqlc.arg (name) || '%');
