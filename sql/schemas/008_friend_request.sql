-- +goose Up
CREATE TABLE friend_requests (
	request_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	sender_id uuid NOT NULL REFERENCES users (user_id),
	receiver_id uuid NOT NULL REFERENCES users (user_id),
	accepted bool DEFAULT FALSE,
	created_at timestamp,
	UNIQUE (sender_id, receiver_id),
	CONSTRAINT chk_different_users CHECK (sender_id != receiver_id)
);

CREATE INDEX idx_friend_requests_sender_id ON friend_requests (sender_id);

CREATE INDEX idx_friend_requests_receiver_id ON friend_requests (receiver_id);

-- +goose Down
DROP INDEX idx_friend_requests_sender_id;

DROP INDEX idx_friend_requests_receiver_id;

DROP TABLE friend_requests;
