CREATE TYPE order_status AS ENUM ('pending','approved','rejected','paid','processing','shipped','delivered','cancelled');

CREATE TABLE medicines (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name text NOT NULL,
  price numeric(12,2) NOT NULL CHECK (price >= 0),
  stock numeric(12,2) NOT NULL CHECK (stock >= 0),  -- หรือ int ถ้าต้องการจำนวนเต็ม
  unit text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);

CREATE TABLE orders (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  patient_id uuid NOT NULL,
  doctor_id uuid,          
  reference_order_id uuid, 
  total_amount numeric(12,2) NOT NULL CHECK (total_amount >= 0),
  note text,
  submitted_at timestamptz,
  reviewed_at timestamptz,
  status order_status NOT NULL DEFAULT 'pending',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT fk_orders_ref FOREIGN KEY (reference_order_id) REFERENCES orders(id) ON DELETE SET NULL
);

CREATE TABLE order_items (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id uuid NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  medicine_id uuid NOT NULL REFERENCES medicines(id) ON DELETE RESTRICT,
  quantity numeric(12,2) NOT NULL CHECK (quantity > 0),
  CONSTRAINT unique_item_per_medicine UNIQUE (order_id, medicine_id)
);

CREATE INDEX idx_orders_patient_created ON orders (patient_id, created_at);
