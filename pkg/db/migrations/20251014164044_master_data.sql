-- +goose Up
-- +goose StatementBegin

-- Insert seed medicines
INSERT INTO medicines (name, price, stock, unit) VALUES
  ('Paracetamol', 5.99, 500, 'tablet'),
  ('Ibuprofen', 7.99, 300, 'tablet'),
  ('Amoxicillin', 15.99, 200, 'capsule'),
  ('Metformin', 8.99, 250, 'tablet'),
  ('Lisinopril', 12.99, 150, 'tablet'),
  ('Atorvastatin', 14.99, 180, 'tablet'),
  ('Omeprazole', 9.99, 220, 'capsule'),
  ('Aspirin', 4.99, 400, 'tablet'),
  ('Cetirizine', 6.99, 350, 'tablet'),
  ('Vitamin C', 3.99, 600, 'tablet')
ON CONFLICT DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Delete seed medicines (optional, can be removed if you prefer to keep data)
DELETE FROM medicines WHERE name IN (
  'Paracetamol', 'Ibuprofen', 'Amoxicillin', 'Metformin', 'Lisinopril',
  'Atorvastatin', 'Omeprazole', 'Aspirin', 'Cetirizine', 'Vitamin C'
);

-- +goose StatementEnd
