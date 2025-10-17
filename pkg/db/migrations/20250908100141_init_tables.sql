-- +goose Up
-- +goose StatementBegin

DO $$ BEGIN
  CREATE TYPE delivery_status AS ENUM ('pending','in_transit','delivered','failed');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE delivery_method_enum AS ENUM ('flash','pick_up');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE order_status AS ENUM ('pending','approved','rejected','paid','processing','shipped','delivered','cancelled');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

CREATE TABLE IF NOT EXISTS medicines (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name text NOT NULL,
  price numeric(12,2) NOT NULL CHECK (price >= 0),
  stock numeric(12,2) NOT NULL CHECK (stock >= 0),
  unit text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);

CREATE TABLE IF NOT EXISTS orders (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  patient_id uuid NOT NULL,
  doctor_id uuid,
  total_amount numeric(12,2) NOT NULL CHECK (total_amount >= 0),
  note text,
  submitted_at timestamptz,
  reviewed_at timestamptz,
  status order_status NOT NULL DEFAULT 'pending',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_orders_patient_created
  ON orders (patient_id, created_at);

CREATE TABLE IF NOT EXISTS delivery_informations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL,
  address text NOT NULL,
  phone_number text NOT NULL,
  version int NOT NULL DEFAULT 1 CHECK (version > 0),
  delivery_method delivery_method_enum NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT unique_user_version UNIQUE (user_id, delivery_method, version)
);

CREATE TABLE IF NOT EXISTS deliveries (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id uuid NOT NULL UNIQUE,
  delivery_information uuid NOT NULL,
  tracking_number text,
  status delivery_status NOT NULL DEFAULT 'pending',
  delivered_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT fk_deliveries_order
    FOREIGN KEY (order_id)
    REFERENCES orders(id)
    ON DELETE CASCADE,
  CONSTRAINT fk_deliveries_delivery_info
    FOREIGN KEY (delivery_information)
    REFERENCES delivery_informations(id)
    ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS order_items (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id uuid NOT NULL,
  medicine_id uuid NOT NULL,
  quantity numeric(12,2) NOT NULL CHECK (quantity > 0),
  CONSTRAINT fk_order_items_order
    FOREIGN KEY (order_id)
    REFERENCES orders(id)
    ON DELETE CASCADE,
  CONSTRAINT fk_order_items_medicine
    FOREIGN KEY (medicine_id)
    REFERENCES medicines(id)
    ON DELETE RESTRICT,
  CONSTRAINT unique_item_per_medicine UNIQUE (order_id, medicine_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS order_items CASCADE;
DROP TABLE IF EXISTS deliveries CASCADE;
DROP TABLE IF EXISTS delivery_informations CASCADE;
DROP TABLE IF EXISTS orders CASCADE;
DROP TABLE IF EXISTS medicines CASCADE;

DROP TYPE IF EXISTS order_status CASCADE;
DROP TYPE IF EXISTS delivery_method_enum CASCADE;
DROP TYPE IF EXISTS delivery_status CASCADE;

-- +goose StatementEnd
