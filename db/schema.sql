CREATE TABLE IF NOT EXISTS user (
  id integer PRIMARY KEY,
  fullname text NOT NULL,
  email text NOT NULL,
  password_hash text NOT NULL
);

CREATE TABLE IF NOT EXISTS board (
  id integer PRIMARY KEY,
  name text NOT NULL,
  slug text NOT NULL,
  user_id integer NOT NULL
);

CREATE TABLE IF NOT EXISTS column (
  id integer PRIMARY KEY,
  name text NOT NULL,
  element_order integer NOT NULL,
  board_id integer NOT NULL
);

CREATE TABLE IF NOT EXISTS item (
  id integer PRIMARY KEY,
  name text NOT NULL,
  element_order integer NOT NULL,
  column_id integer NOT NULL
);