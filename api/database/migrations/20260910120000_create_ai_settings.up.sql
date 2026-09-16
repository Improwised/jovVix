-- +migrate Up
CREATE TABLE ai_settings (
  user_id CHAR(20) PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  provider VARCHAR(80) NOT NULL, base_url VARCHAR(300) NOT NULL, model VARCHAR(200) NOT NULL,
  encrypted_api_key BYTEA NOT NULL, salt BYTEA NOT NULL, nonce BYTEA NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);
