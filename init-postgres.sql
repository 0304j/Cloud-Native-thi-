CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE IF EXISTS payments;
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    provider VARCHAR(32) NOT NULL,
    amount NUMERIC(10,2) NOT NULL,
    currency VARCHAR(8) NOT NULL,
    status VARCHAR(16) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Example data

INSERT INTO payments (provider, amount, currency, status)
VALUES
  ('paypal', 19.99, 'EUR', 'success'),
  ('stripe', 49.50, 'USD', 'pending'),
  ('ueberweisung', 10.00, 'EUR', 'failed');
