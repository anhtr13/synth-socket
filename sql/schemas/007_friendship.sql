-- +goose Up
CREATE TABLE friendships (
	friendship_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	user1_id uuid NOT NULL REFERENCES users (user_id),
	user2_id uuid NOT NULL REFERENCES users (user_id),
	created_at timestamp,
	UNIQUE (user1_id, user2_id),
	CONSTRAINT chk_different_users CHECK (user1_id != user2_id)
);

-- +goose Down
DROP TABLE friendships;
