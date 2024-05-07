CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA pg_catalog;

--events
CREATE TABLE events (
  id uuid NOT NULL DEFAULT uuid_generate_v1mc(),
  event_name text NOT NULL,
  event_description text NOT NULL,
  event_location text NOT NULL,
  event_date date NOT NULL,
  CONSTRAINT events_pk PRIMARY KEY (id),
);