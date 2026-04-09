CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR(64) PRIMARY KEY,
    customer_name VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    total_amount BIGINT NOT NULL
);

INSERT INTO orders (id, customer_name, status, total_amount)
VALUES
    ('ord-001', 'Alice', 'pending', 125000),
    ('ord-002', 'Bob', 'paid', 340000)
ON CONFLICT (id) DO NOTHING;
