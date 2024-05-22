CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA pg_catalog;

--users
CREATE TABLE IF NOT EXISTS users (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v1mc(),
  email TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL
);

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

--registrations
CREATE TABLE IF NOT EXISTS registrations (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v1mc(),
  event_id uuid,
  user_id uuid,
  FOREIGN KEY(event_id) REFERENCES events(id),
  FOREIGN KEY(user_id) REFERENCES users(id)
);

