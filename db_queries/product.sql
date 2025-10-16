CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    barcode VARCHAR(20) NOT NULL UNIQUE, 
    name VARCHAR(255),
    brand_name VARCHAR(255),
    category VARCHAR(100),
    sub_category VARCHAR(100),
    image_url TEXT,
    price NUMERIC(10,2),
    packaging_material VARCHAR(100),
    manufacturing_location VARCHAR(255),
    disposal_method VARCHAR(100)
);

CREATE TABLE product_requests (
    id SERIAL PRIMARY KEY,
    barcode VARCHAR(100) NOT NULL,
    name VARCHAR(255),
    brand_name VARCHAR(255),
    image_url TEXT,
    user_id INT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

INSERT INTO products 
(barcode, name, brand_name, category, sub_category, image_url, price, packaging_material, manufacturing_location, disposal_method)
VALUES
('1234567890123', 'Mineral Water 1L', 'FreshCo', 'Beverages', 'Water', 'http://example.com/image.jpg', 25.50, 'Plastic Bottle', 'Dhaka, Bangladesh', 'Recycle');

