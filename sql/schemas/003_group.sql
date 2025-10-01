-- +goose Up
CREATE TABLE groups (
	"group_id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	"group_name" varchar(255) NOT NULL,
	"group_picture" varchar(255),
	"created_by" uuid NOT NULL REFERENCES users (user_id),
	"created_at" timestamp NOT NULL
);

CREATE INDEX idx_groups_group_name ON groups (group_name);

-- +goose Down
DROP INDEX idx_groups_group_name;

DROP TABLE groups;
