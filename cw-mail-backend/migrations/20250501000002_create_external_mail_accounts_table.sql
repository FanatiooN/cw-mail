-- +goose Up
CREATE TABLE external_mail_accounts (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id) ON DELETE CASCADE,
  email VARCHAR(255) NOT NULL,
  account_type VARCHAR(10) NOT NULL,
  server VARCHAR(255) NOT NULL,
  port INT NOT NULL,
  username VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL,
  last_sync_time TIMESTAMP WITH TIME ZONE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX idx_external_mail_accounts_user_id ON external_mail_accounts(user_id);

-- +goose Down
DROP TABLE external_mail_accounts; 