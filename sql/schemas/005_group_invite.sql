-- +goose Up
CREATE TABLE group_invites (
	invite_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	group_id uuid NOT NULL REFERENCES groups (group_id),
	sender_id uuid NOT NULL REFERENCES users (user_id),
	receiver_id uuid NOT NULL REFERENCES users (user_id),
	accepted bool DEFAULT FALSE,
	created_at timestamp NOT NULL,
	UNIQUE (group_id, sender_id, receiver_id)
);

CREATE INDEX idx_group_invites_sender_id ON group_invites (sender_id);

CREATE INDEX idx_group_invites_receiver_id ON group_invites (receiver_id);

-- +goose Down
DROP INDEX idx_group_invites_sender_id;

DROP INDEX idx_group_invites_receiver_id;

DROP TABLE group_invites;
