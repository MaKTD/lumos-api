ALTER TABLE lumos.users ADD COLUMN IF NOT EXISTS last_transaction_id varchar NOT NULL DEFAULT '';
ALTER TABLE lumos.users ADD COLUMN IF NOT EXISTS updated_at timestamptz NOT NULL DEFAULT NOW();

CREATE OR REPLACE TRIGGER set_timestamp_users
  BEFORE UPDATE ON lumos.users
  FOR EACH ROW
  EXECUTE FUNCTION trigger_set_timestamp();
