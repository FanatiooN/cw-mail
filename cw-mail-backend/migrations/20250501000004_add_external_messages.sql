-- +goose Up
CREATE TABLE external_messages (
  id SERIAL PRIMARY KEY,
  external_account_id INTEGER NOT NULL REFERENCES external_mail_accounts(id) ON DELETE CASCADE,
  message_id VARCHAR(255) NOT NULL,
  subject VARCHAR(255),
  "from" VARCHAR(255) NOT NULL,
  "to" TEXT NOT NULL,
  body TEXT,
  body_html TEXT,
  is_read BOOLEAN NOT NULL DEFAULT FALSE,
  has_attachments BOOLEAN NOT NULL DEFAULT FALSE,
  received_at TIMESTAMP WITH TIME ZONE NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX idx_external_messages_account_id ON external_messages(external_account_id);
CREATE INDEX idx_external_messages_message_id ON external_messages(message_id);
CREATE INDEX idx_external_messages_received_at ON external_messages(received_at);

CREATE TABLE external_attachments (
  id SERIAL PRIMARY KEY,
  external_message_id INTEGER NOT NULL REFERENCES external_messages(id) ON DELETE CASCADE,
  filename VARCHAR(255) NOT NULL,
  content_type VARCHAR(100) NOT NULL,
  size BIGINT NOT NULL,
  storage_path VARCHAR(255) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX idx_external_attachments_message_id ON external_attachments(external_message_id);

-- +goose Down
DROP TABLE external_attachments;
DROP TABLE external_messages; 