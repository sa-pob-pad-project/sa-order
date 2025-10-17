CREATE TYPE delivery_status AS ENUM ('pending','in_transit','delivered','failed');
CREATE TYPE delivery_method_enum AS ENUM ('flash','pick_up');

CREATE TABLE delivery_informations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL,                        
  address text NOT NULL,
  phone_number text NOT NULL,
  version int NOT NULL DEFAULT 1 CHECK (version > 0),
  delivery_method delivery_method_enum NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT unique_user_version UNIQUE (user_id, delivery_method, version)
);

CREATE TABLE deliveries (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id uuid NOT NULL,                         
  delivery_information uuid NOT NULL REFERENCES delivery_informations(id) ON DELETE RESTRICT,
  tracking_number text,
  status delivery_status NOT NULL DEFAULT 'pending',
  delivered_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_deliveries_order ON deliveries(order_id);
