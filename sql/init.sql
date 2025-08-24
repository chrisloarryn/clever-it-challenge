-- Drop table if exists (for development purposes)
DROP TABLE IF EXISTS beer;

-- Create beer table with improved schema
CREATE TABLE beer
(
    id         INTEGER PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    brewery    VARCHAR(100) NOT NULL,
    country    VARCHAR(100) NOT NULL,
    currency   CHAR(3)      NOT NULL,
    price      DECIMAL(10, 6) NOT NULL CHECK (price >= 0),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX idx_beer_name ON beer(name);
CREATE INDEX idx_beer_brewery ON beer(brewery);
CREATE INDEX idx_beer_country ON beer(country);
CREATE INDEX idx_beer_currency ON beer(currency);
CREATE INDEX idx_beer_created_at ON beer(created_at);

-- Add some sample data for testing
INSERT INTO beer (id, name, brewery, country, currency, price, created_at, updated_at) VALUES
(1, 'Cerveza Cristal', 'CCU', 'Chile', 'CLP', 1200.00, NOW(), NOW()),
(2, 'Escudo', 'CCU', 'Chile', 'CLP', 1100.00, NOW(), NOW()),
(3, 'Heineken', 'Heineken N.V.', 'Netherlands', 'EUR', 2.50, NOW(), NOW()),
(4, 'Corona Extra', 'Grupo Modelo', 'Mexico', 'MXN', 35.00, NOW(), NOW()),
(5, 'Budweiser', 'Anheuser-Busch', 'United States', 'USD', 4.50, NOW(), NOW()),
(6, 'Stella Artois', 'Anheuser-Busch InBev', 'Belgium', 'EUR', 3.20, NOW(), NOW()),
(7, 'Guinness', 'Guinness Brewery', 'Ireland', 'EUR', 4.80, NOW(), NOW()),
(8, 'Asahi Super Dry', 'Asahi Breweries', 'Japan', 'JPY', 250.00, NOW(), NOW());

-- Create a function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_beer_updated_at 
    BEFORE UPDATE ON beer 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();