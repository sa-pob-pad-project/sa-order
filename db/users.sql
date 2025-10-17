CREATE TYPE roles AS ENUM ('patient','doctor','admin');
CREATE TYPE gender_enum AS ENUM ('male','female');

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  password text NOT NULL,
  first_name text NOT NULL,
  last_name text NOT NULL,
  gender gender_enum NOT NULL,
  phone_number text,
  role roles NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);

CREATE TABLE patients (
  user_id uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  hospital_id text UNIQUE,
  birth_date timestamptz,
  id_card_number varchar(13),
  address text,
  allergies text,
  emergency_contact text,
  blood_type varchar(5),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz,
  CONSTRAINT patients_id_card_len CHECK (id_card_number IS NULL OR length(id_card_number)=13)
);

CREATE TABLE healthcare_entitlements (
  healthcare_entitlement text PRIMARY KEY
);

CREATE TABLE user_healthcare_entitlement (
  patient_id uuid NOT NULL REFERENCES patients(user_id) ON DELETE CASCADE,
  healthcare_entitlement text NOT NULL REFERENCES healthcare_entitlements(healthcare_entitlement) ON DELETE RESTRICT,
  PRIMARY KEY (patient_id, healthcare_entitlement)
);

CREATE TABLE doctors (
  user_id uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  username text UNIQUE NOT NULL,
  specialty text,
  bio text,
  years_experience int CHECK (years_experience IS NULL OR years_experience >= 0),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);

CREATE TABLE admins (
  user_id uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  username text UNIQUE NOT NULL
);
