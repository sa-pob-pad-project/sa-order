CREATE TYPE payment_status AS ENUM ('pending','success','failed');
CREATE TYPE payment_method AS ENUM ('credit_card','promptpay');

CREATE TABLE payment_informations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL,                      -- cross-service to user_service
  type payment_method NOT NULL,
  details jsonb NOT NULL,
  version int NOT NULL DEFAULT 1 CHECK (version > 0),
  created_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT unique_payment_profile UNIQUE (user_id, type, version)
);

CREATE TABLE payment_attempts (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id uuid NOT NULL,                    
  payment_information_id uuid,               
  method payment_method NOT NULL,
  status payment_status NOT NULL DEFAULT 'pending',
  created_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT fk_attempt_profile FOREIGN KEY (payment_information_id) REFERENCES payment_informations(id) ON DELETE SET NULL
);

CREATE TABLE payments (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  attempt_id uuid NOT NULL REFERENCES payment_attempts(id) ON DELETE CASCADE,
  amount numeric(12,2) NOT NULL CHECK (amount >= 0),
  order_id uuid NOT NULL,                   
  paid_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_attempts_order ON payment_attempts(order_id);
CREATE INDEX idx_payments_order ON payments(order_id);
