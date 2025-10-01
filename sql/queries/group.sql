-- name: CreateGroup :one
INSERT INTO
	groups (
		"group_name",
		"group_picture",
		"created_by",
		"created_at"
	)
VALUES
	($1, $2, $3, $4)
RETURNING
	*;

-- name: FindGroupById :one
SELECT
	*
FROM
	groups
WHERE
	group_id = $1;

-- name: FindGroupByCreatorAndName :one
SELECT
	*
FROM
	groups
WHERE
	created_by = $1
	AND group_name = $2;

-- name: GetAllGroup :many
SELECT
	*
FROM
	groups
ORDER BY
	group_name
LIMIT
	$1
OFFSET
	$2;

-- name: GetGroupByGroupName :many
SELECT
	*
FROM
	groups
WHERE
	group_name LIKE ('%' || $1 || '%')
ORDER BY
	group_name
LIMIT
	$2
OFFSET
	$3;

-- name: GetAllGroupIds :many
SELECT
	group_id
FROM
	groups;
