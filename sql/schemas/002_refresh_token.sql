-- +goose Up
CREATE TABLE refresh_tokens (
	token varchar(64) PRIMARY KEY,
	user_id uuid NOT NULL REFERENCES users (user_id),
	user_email varchar(255) NOT NULL,
	user_name varchar(255) NOT NULL,
	expired_at timestamp NOT NULL,
	created_at timestamp without time zone DEFAULT now(),
	UNIQUE (token, user_id)
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id);

-- +goose Down
DROP INDEX idx_refresh_tokens_user_id;

DROP TABLE refresh_tokens;
