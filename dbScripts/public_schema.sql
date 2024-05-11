CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA pg_catalog;

--events
CREATE TABLE IF NOT EXISTS events (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v1mc(),
  event_name text NOT NULL,
  event_description text NOT NULL,
  event_location text NOT NULL,
  event_date date NOT NULL,
  user_id uuid NOT NULL,
  FOREIGN KEY(user_id) REFERENCES users(id)
);

--users
CREATE TABLE IF NOT EXISTS users (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v1mc(),
  email TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL
)