DROP TABLE IF EXISTS payments;

CREATE TABLE payments (
                          id UUID PRIMARY KEY,
                          order_id UUID NOT NULL,
                          user_id UUID NOT NULL,
                          provider VARCHAR(32) NOT NULL,
                          amount DOUBLE PRECISION NOT NULL,
                          currency VARCHAR(8) NOT NULL,
                          status VARCHAR(16) NOT NULL,
                          created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                          updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                          CHECK (amount > 0),
                          CHECK (status IN ('pending','success','failed'))
);

CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments (order_id);
CREATE INDEX IF NOT EXISTS idx_payments_user_id ON payments (user_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments (status);