-- +goose Up
CREATE TABLE group_members (
	"group_id" uuid NOT NULL REFERENCES groups (group_id),
	"member_id" uuid NOT NULL REFERENCES users (user_id),
	"joined_at" timestamp NOT NULL,
	UNIQUE (group_id, member_id)
);

CREATE INDEX idx_group_members_member_id ON group_members (member_id);

-- +goose Down
DROP INDEX idx_group_members_member_id;

DROP TABLE group_members;
