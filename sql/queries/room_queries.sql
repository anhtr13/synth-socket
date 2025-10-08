-- name: CreateRoom :one
INSERT INTO
	rooms (room_name, room_picture, created_by)
VALUES
	($1, $2, $3)
RETURNING
	*;

-- name: FindRoomById :one
SELECT
	*
FROM
	rooms
WHERE
	room_id = $1;

-- name: FindRoomByCreatorAndName :one
SELECT
	*
FROM
	rooms
WHERE
	created_by = $1
	AND room_name = $2;

-- name: GetRoomsByCreator :many
SELECT
	*
FROM
	rooms
WHERE
	created_by = $1
ORDER BY
	room_name
LIMIT
	$2
OFFSET
	$3;

-- name: GetRoomsByCreatorAndName :many
SELECT
	*
FROM
	rooms
WHERE
	created_by = $1
	AND room_name LIKE ('%' || sqlc.arg (room_name) || '%')
ORDER BY
	room_name
LIMIT
	$2
OFFSET
	$3;

-- name: GetRoomsDataByMemberIdAndRoomName :many
SELECT
	r.*,
	rm.joined_at,
	u.user_name AS owner_name,
	u.profile_image AS owner_image
FROM
	(
		SELECT
			*
		FROM
			room_members
		WHERE
			member_id = $1
		ORDER BY
			joined_at DESC
	) rm
	LEFT JOIN rooms r ON rm.room_id = r.room_id
	LEFT JOIN users u ON r.created_by = u.user_id
WHERE
	r.room_name LIKE ('%' || sqlc.arg (name) || '%')
LIMIT
	$2
OFFSET
	$3;

-- name: GetRoomsDataByMemberId :many
SELECT
	r.*,
	rm.joined_at,
	u.user_name AS owner_name,
	u.profile_image AS owner_image
FROM
	(
		SELECT
			*
		FROM
			room_members
		WHERE
			member_id = $1
		ORDER BY
			joined_at DESC
		LIMIT
			$2
		OFFSET
			$3
	) rm
	LEFT JOIN rooms r ON rm.room_id = r.room_id
	LEFT JOIN users u ON r.created_by = u.user_id;
