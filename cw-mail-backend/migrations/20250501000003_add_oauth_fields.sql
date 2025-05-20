-- +goose Up
ALTER TABLE external_mail_accounts
  ADD COLUMN auth_type VARCHAR(10) NOT NULL DEFAULT 'basic',
  ADD COLUMN access_token TEXT DEFAULT '',
  ADD COLUMN refresh_token TEXT DEFAULT '',
  ADD COLUMN token_expiry TIMESTAMP WITH TIME ZONE,
  ADD COLUMN provider_name VARCHAR(50) DEFAULT '',
  ALTER COLUMN password DROP NOT NULL,
  ALTER COLUMN password SET DEFAULT '';

-- +goose Down
ALTER TABLE external_mail_accounts
  DROP COLUMN auth_type,
  DROP COLUMN access_token,
  DROP COLUMN refresh_token,
  DROP COLUMN token_expiry,
  DROP COLUMN provider_name,
  ALTER COLUMN password SET NOT NULL,
  ALTER COLUMN password DROP DEFAULT; 