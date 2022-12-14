-- Add UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Set timezone
-- For more information, please visit:
-- https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
SET TIMEZONE="Asia/Jakarta";

-- Create users table
CREATE TABLE users (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    updated_at TIMESTAMP NULL,
    email VARCHAR (255) NOT NULL UNIQUE,
    password_hash VARCHAR (255) NOT NULL,
    user_status INT NOT NULL,
    user_role VARCHAR (25) NOT NULL
);

-- Create books table
CREATE TABLE books (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    updated_at TIMESTAMP NULL,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title VARCHAR (255) NOT NULL,
    author VARCHAR (255) NOT NULL,
    book_status INT NOT NULL,
    book_attrs JSONB NOT NULL
);

--create cart table
CREATE TABLE carts (
    id INT DEFAULT AUTO_INCREMENT () PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    updated_at TIMESTAMP NULL,
    product_id UUID NOT NULL REFERENCES produks (id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    cart_status INT NOT NULL,
    cart_attrs JSONB NOT NULL
);

--Create invoice table
CREATE TABLE invoices (
    id INT DEFAULT AUTO_INCREMENT  () PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    updated_at TIMESTAMP NULL,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    invoice_status INT NOT NULL,
    invoice_attrs JSONB NOT NULL
);
--Create invoice_item table
CREATE TABLE invoice_items (
    id INT DEFAULT AUTO_INCREMENT () PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    updated_at TIMESTAMP NULL,
    invoice_id UUID NOT NULL REFERENCES invoices (id) ON DELETE CASCADE,
    produk_id UUID NOT NULL REFERENCES produks (id) ON DELETE CASCADE,
    price INT NOT NULL,
    qty INT NOT NULL,
    item_status INT NOT NULL,
    item_attrs JSONB NOT NULL
);
--Create produk table
CREATE TABLE produks (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    updated_at TIMESTAMP NULL,
    total INT NOT NULL,
    produk_status INT NOT NULL,
    produk_attrs JSONB NOT NULL
);

-- Add indexes
CREATE INDEX active_users ON users (id) WHERE user_status = 1;
CREATE INDEX active_books ON books (title) WHERE book_status = 1;
CREATE INDEX active_invoice ON invoices (id) WHERE invoice_status = 1;
CREATE INDEX active_invoice_item ON invoice_items (id) WHERE item_status = 1;
CREATE INDEX active_produk ON produks (id) WHERE produk_status = 1;