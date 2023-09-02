CREATE TABLE IF NOT EXISTS user (
  id integer PRIMARY KEY,
  fullname text NOT NULL,
  email text NOT NULL,
  password_hash text NOT NULL
);