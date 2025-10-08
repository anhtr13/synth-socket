-- +goose Up
CREATE TABLE users (
	user_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	user_email varchar(255) NOT NULL UNIQUE,
	user_name varchar(255) NOT NULL,
	password varchar(64) NOT NULL,
	profile_image varchar(255),
	created_at timestamp without time zone DEFAULT now()
);

CREATE INDEX idx_users_user_name ON users (user_name);

-- +goose Down
DROP INDEX idx_users_user_name;

DROP TABLE users;
