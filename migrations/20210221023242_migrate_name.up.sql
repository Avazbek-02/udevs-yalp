CREATE TYPE auth_user_type AS ENUM (
  'client',
  'admin',
  'superadmin'
);

CREATE TYPE auth_status AS ENUM (
  'active',
  'blocked',
  'in_verify'
);

CREATE TYPE project_status AS ENUM (
  'new',
  'in_verify',
  'verified',
  'canceled',
  'won',
  'lost'
);

CREATE TABLE auth_user (
  id uuid PRIMARY KEY,
  full_name varchar NOT NULL,
  phone_number varchar(10) NOT NULL DEFAULT '',
  username varchar(20) NOT NULL DEFAULT '',
  password varchar,
  user_type auth_user_type NOT NULL DEFAULT 'client',
  status auth_status NOT NULL DEFAULT 'in_verify',
  created_at timestamp NOT NULL DEFAULT 'now()',
  updated_at timestamp NOT NULL DEFAULT 'now()'
);

CREATE UNIQUE INDEX ON auth_user (phone_number, username);

CREATE TABLE auth_session (
  id uuid PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES auth_user (id),
  ip_address varchar NOT NULL,
  user_agent text NOT NULL DEFAULT '',
  is_active bool NOT NULL DEFAULT 'true',
  expires_at timestamp,
  last_active_at timestamp NOT NULL DEFAULT 'now()',
  created_at timestamp NOT NULL DEFAULT 'now()',
  updated_at timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE cycle (
  id uuid PRIMARY KEY,
  title varchar NOT NULL,
  description text NOT NULL,
  start_date timestamp NOT NULL,
  project_accept_end_date timestamp NOT NULL CHECK (project_accept_end_date > start_date),
  moderation_end_date timestamp NOT NULL CHECK (moderation_end_date > project_accept_end_date),
  voting_end_date timestamp NOT NULL CHECK (voting_end_date > moderation_end_date),
  created_at timestamp NOT NULL DEFAULT 'now()',
  updated_at timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE region (
  id uuid PRIMARY KEY,
  title varchar NOT NULL,
  created_at timestamp NOT NULL DEFAULT 'now()',
  updated_at timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE district (
  id uuid PRIMARY KEY,
  region_id uuid NOT NULL REFERENCES region (id),
  title varchar NOT NULL,
  created_at timestamp NOT NULL DEFAULT 'now()',
  updated_at timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE project_type (
  id uuid PRIMARY KEY,
  title varchar NOT NULL,
  slug varchar UNIQUE NOT NULL,
  active bool NOT NULL DEFAULT 'true',
  created_at timestamp NOT NULL DEFAULT 'now()',
  updated_at timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE project (
  id uuid PRIMARY KEY,
  title varchar NOT NULL,
  description varchar NOT NULL,
  images varchar[] NOT NULL DEFAULT '{}',
  district_id uuid NOT NULL REFERENCES district (id),
  type_id uuid NOT NULL REFERENCES project_type (id),
  creater_id uuid NOT NULL REFERENCES auth_user (id),
  cycle_id uuid NOT NULL REFERENCES cycle (id),
  status project_status NOT NULL DEFAULT 'new',
  created_at timestamp NOT NULL DEFAULT 'now()',
  updated_at timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE vote (
    id uuid PRIMARY KEY,
    creater_id uuid NOT NULL REFERENCES auth_user (id),
    project_id uuid NOT NULL REFERENCES project (id),
    cycle_id uuid NOT NULL REFERENCES cycle (id),
    phone_number varchar NOT NULL,
    passport_series varchar NOT NULL,
    created_at timestamp NOT NULL DEFAULT 'now()',
    updated_at timestamp NOT NULL DEFAULT 'now()'
);

CREATE UNIQUE INDEX ON vote (cycle_id, phone_number);
CREATE UNIQUE INDEX ON vote (cycle_id, passport_series);