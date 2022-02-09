ALTER TABLE receipts ADD COLUMN user_id uuid REFERENCES users(id);
ALTER TABLE expenses ADD COLUMN user_id uuid REFERENCES users(id);