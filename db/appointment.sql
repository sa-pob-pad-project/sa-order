CREATE TYPE slot_status AS ENUM ('open','booked','cancelled');
CREATE TYPE appointment_status AS ENUM ('scheduled','completed','cancelled');
CREATE TYPE day_of_week AS ENUM ('mon','tue','wed','thu','fri','sat','sun');

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE doctor_shifts (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  doctor_id uuid NOT NULL,                
  weekday day_of_week NOT NULL,
  start_time timestamptz NOT NULL,
  end_time timestamptz NOT NULL,
  duration_min int NOT NULL CHECK (duration_min > 0),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz,
  CONSTRAINT shift_time_window CHECK (end_time > start_time),
  CONSTRAINT shift_unique_per_day UNIQUE (doctor_id, weekday, start_time, end_time)
);

CREATE TABLE doctor_slots (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  doctor_id uuid NOT NULL,
  start_time timestamptz NOT NULL,
  end_time timestamptz NOT NULL,
  status slot_status NOT NULL DEFAULT 'open',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz,
  CONSTRAINT slot_time_window CHECK (end_time > start_time)
);


CREATE TABLE appointments (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  patient_id uuid NOT NULL,                
  doctor_id uuid NOT NULL,                 
  slot_id uuid UNIQUE,                     
  status appointment_status NOT NULL DEFAULT 'scheduled',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz,
  CONSTRAINT unique_patient_doctor_slot UNIQUE (patient_id, doctor_id, slot_id)
);

ALTER TABLE appointments
  ADD CONSTRAINT fk_appointments_slot
  FOREIGN KEY (slot_id) REFERENCES doctor_slots(id) ON DELETE SET NULL;

CREATE INDEX idx_slots_doctor_time ON doctor_slots (doctor_id, start_time) WHERE deleted_at IS NULL;
CREATE INDEX idx_appt_doctor_time ON appointments (doctor_id, created_at) WHERE deleted_at IS NULL;
