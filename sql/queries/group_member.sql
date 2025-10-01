-- name: CreateGroupMember :one
INSERT INTO
	group_members ("group_id", "member_id", "joined_at")
VALUES
	($1, $2, $3)
RETURNING
	*;

-- name: GetGroupMembersTable :many
SELECT
	*
FROM
	group_members;

-- name: FindGroupMember :one
SELECT
	*
FROM
	group_members
WHERE
	group_id = $1
	AND member_id = $2;

-- name: CountMemmbersInGroup :one
SELECT
	count(member_id)
FROM
	group_members
WHERE
	group_id = $1;

-- name: DeleteGroupMemmber :exec
DELETE FROM group_members
WHERE
	group_id = $1
	AND member_id = $2;

-- name: GetAllGroupMemberInfo :many
SELECT
	u.user_id,
	u.user_name,
	u.profile_image,
	g.joined_at
FROM
	(
		SELECT
			*
		FROM
			group_members
		WHERE
			group_id = $1
	) g
	LEFT JOIN users u ON g.member_id = u.user_id
ORDER BY
	g.joined_at
LIMIT
	$2
OFFSET
	$3;

-- name: GetGroupMemberInfoByName :many
SELECT
	u.user_id,
	u.user_name,
	u.profile_image,
	g.joined_at
FROM
	(
		SELECT
			*
		FROM
			group_members
		WHERE
			group_id = $1
	) g
	LEFT JOIN users u ON g.member_id = u.user_id
WHERE
	u.user_name LIKE ('%' || $2 || '%')
ORDER BY
	g.joined_at
LIMIT
	$3
OFFSET
	$4;

-- name: GetAllUserGroups :many
SELECT
	g.*,
	cm.joined_at
FROM
	(
		SELECT
			*
		FROM
			group_members
		WHERE
			member_id = $1
		ORDER BY
			group_id
	) cm
	LEFT JOIN groups g ON cm.group_id = g.group_id
LIMIT
	$2
OFFSET
	$3;

-- name: GetAllUserGroupIds :many
SELECT
	group_id
FROM
	group_members
WHERE
	member_id = $1;
