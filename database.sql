CREATE TABLE categories (
    category_id SERIAL PRIMARY KEY,
    category_uuid UUID NOT NULL,
    category_name VARCHAR(255) NOT NULL,
    category_slug VARCHAR(255) NOT NULL,
    category_status VARCHAR(10) CHECK (category_status IN ('active', 'deleted')) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
    product_id SERIAL PRIMARY KEY,
    product_uuid UUID NOT NULL,
    category_id INTEGER NOT NULL REFERENCES categories(category_id),
    product_name VARCHAR(255) NOT NULL,
    product_slug VARCHAR(255) NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    stock_quantity INTEGER NOT NULL,
    description TEXT,
    product_status VARCHAR(10) CHECK (product_status IN ('active', 'deleted')) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE customers (
    customer_id SERIAL PRIMARY KEY,
    customer_uuid UUID NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    address TEXT,
    customer_status VARCHAR(10) CHECK (customer_status IN ('active', 'deleted')) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE carts (
    cart_id SERIAL PRIMARY KEY,
    cart_uuid UUID NOT NULL,
    customer_id INTEGER NOT NULL REFERENCES customers(customer_id),
    cart_status VARCHAR(20) CHECK (cart_status IN ('active', 'converted', 'abandoned')) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cart_details (
    cart_detail_id SERIAL PRIMARY KEY,
    cart_detail_uuid UUID NOT NULL,
    cart_id INTEGER NOT NULL REFERENCES carts(cart_id),
    product_id INTEGER NOT NULL REFERENCES products(product_id),
    quantity INTEGER NOT NULL,
    cart_detail_status VARCHAR(30) CHECK (cart_detail_status IN ('active', 'deleted_by_customer')) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "orders" (
    order_id SERIAL PRIMARY KEY,
    order_uuid UUID NOT NULL,
    invoice_number VARCHAR(50) NOT NULL,
    customer_id INTEGER NOT NULL REFERENCES customers(customer_id),
    cart_id INTEGER NOT NULL REFERENCES carts(cart_id),
    order_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    order_payment_type VARCHAR(20) CHECK (order_payment_type IN ('alfamart', 'indomaret', 'virtual_account', 'jenius')) NOT NULL,
    order_payment_status VARCHAR(10) CHECK (order_payment_status IN ('unpaid', 'completed')) NOT NULL,
    order_status VARCHAR(10) CHECK (order_status IN ('pending', 'scheduled')) NOT NULL,
    grand_total NUMERIC(10, 2) NOT NULL,
    expired_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_details (
    order_detail_id SERIAL PRIMARY KEY,
    order_detail_uuid UUID NOT NULL,
    order_id INTEGER NOT NULL REFERENCES "orders"(order_id),
    product_id INTEGER NOT NULL REFERENCES products(product_id),
    quantity INTEGER NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    final_price NUMERIC(10, 2) NOT NULL,
    order_detail_status VARCHAR(10) CHECK (order_detail_status IN ('active', 'deleted')) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO categories (category_uuid, category_name, category_slug, category_status, created_at, updated_at)
VALUES
    (gen_random_uuid(), 'Electronics', 'electronics', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), 'Books', 'books', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), 'Clothing', 'clothing', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), 'Home & Kitchen', 'home-kitchen', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), 'Sports', 'sports', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
   
   
INSERT INTO products (product_uuid, category_id, product_name, product_slug, price, stock_quantity, description, product_status, created_at, updated_at)
VALUES
    (gen_random_uuid(), 1, 'Smartphone', 'smartphone', 6999900, 50, 'Latest model with advanced features', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), 1, 'Laptop', 'laptop', 12999900, 30, 'High performance laptop for work and play', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), 2, 'Science Fiction Novel', 'science-fiction-novel', 159900, 100, 'A thrilling science fiction story', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), 3, 'Mens T-Shirt', 'mens-tshirt', 49990, 200, 'Comfortable cotton t-shirt', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), 4, 'Blender', 'blender', 499900, 75, 'High speed blender for smoothies and more', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), 5, 'Yoga Mat', 'yoga-mat', 259900, 150, 'Non-slip yoga mat for home workouts', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
