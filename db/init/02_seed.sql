INSERT INTO payments (id, order_id, user_id, provider, amount, currency, status)
VALUES
    (gen_random_uuid(), gen_random_uuid(), gen_random_uuid(), 'stripe', 19.99, 'EUR', 'pending'),
    (gen_random_uuid(), gen_random_uuid(), gen_random_uuid(), 'paypal',  9.49, 'EUR', 'success');