CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  code VARCHAR(50) NOT NULL UNIQUE,
  description TEXT NOT NULL,
  balance INTEGER NOT NULL CHECK (balance >= 0)
);