CREATE TABLE orders
(
    id          SERIAL PRIMARY KEY,
    customer_id INT NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE order_products
(
    order_id   INT NOT NULL,
    product_id INT NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders (id),
    PRIMARY KEY (order_id, product_id)
);

CREATE TABLE inventory
(
    product_id INT PRIMARY KEY,
    quantity   INT NOT NULL CHECK (quantity >= 0),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO inventory (product_id, quantity)
VALUES (1, 10),
       (2, 10),
       (3, 10);
