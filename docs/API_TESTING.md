# Beer API - cURL Test Commands

## Prerequisites
Start the server first:
```bash
make dev
# or
export DB_TYPE=inmemory && go run cmd/main.go
```

The API will be available at `http://localhost:8080`

## Health Check
```bash
curl -X GET http://localhost:8080/ping
```

## Create Beers

### Create Corona Extra
```bash
curl -X POST http://localhost:8080/api/v1/beers \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "name": "Corona Extra",
    "brewery": "Cervecería Modelo",
    "country": "Mexico",
    "price": 1200,
    "currency": "CLP"
  }'
```

### Create Heineken
```bash
curl -X POST http://localhost:8080/api/v1/beers \
  -H "Content-Type: application/json" \
  -d '{
    "id": 2,
    "name": "Heineken",
    "brewery": "Heineken N.V.",
    "country": "Netherlands",
    "price": 1500,
    "currency": "CLP"
  }'
```

### Create Budweiser (USD)
```bash
curl -X POST http://localhost:8080/api/v1/beers \
  -H "Content-Type: application/json" \
  -d '{
    "id": 3,
    "name": "Budweiser",
    "brewery": "Anheuser-Busch",
    "country": "USA",
    "price": 25.99,
    "currency": "USD"
  }'
```

### Create Stella Artois (EUR)
```bash
curl -X POST http://localhost:8080/api/v1/beers \
  -H "Content-Type: application/json" \
  -d '{
    "id": 4,
    "name": "Stella Artois",
    "brewery": "Anheuser-Busch InBev",
    "country": "Belgium",
    "price": 3.50,
    "currency": "EUR"
  }'
```

### Create Local Beer
```bash
curl -X POST http://localhost:8080/api/v1/beers \
  -H "Content-Type: application/json" \
  -d '{
    "id": 5,
    "name": "Escudo",
    "brewery": "CCU",
    "country": "Chile",
    "price": 900,
    "currency": "CLP"
  }'
```

## Get Beers

### Get All Beers
```bash
curl -X GET http://localhost:8080/api/v1/beers \
  -H "Accept: application/json"
```

### Get Beer by ID
```bash
curl -X GET http://localhost:8080/api/v1/beers/1 \
  -H "Accept: application/json"
```

### Get Non-existent Beer (404 Test)
```bash
curl -X GET http://localhost:8080/api/v1/beers/999 \
  -H "Accept: application/json"
```

## Calculate Box Prices

### Corona 6-pack in CLP (same currency)
```bash
curl -X GET "http://localhost:8080/api/v1/beers/1/boxprice?quantity=6&currency=CLP" \
  -H "Accept: application/json"
```

### Corona 12-pack in USD (currency conversion)
```bash
curl -X GET "http://localhost:8080/api/v1/beers/1/boxprice?quantity=12&currency=USD" \
  -H "Accept: application/json"
```

### Heineken 24-pack in EUR
```bash
curl -X GET "http://localhost:8080/api/v1/beers/2/boxprice?quantity=24&currency=EUR" \
  -H "Accept: application/json"
```

### Budweiser 6-pack in USD (same currency)
```bash
curl -X GET "http://localhost:8080/api/v1/beers/3/boxprice?quantity=6&currency=USD" \
  -H "Accept: application/json"
```

## Error Cases

### Duplicate Beer ID
```bash
curl -X POST http://localhost:8080/api/v1/beers \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "name": "Duplicate Beer",
    "brewery": "Test Brewery",
    "country": "Chile",
    "price": 1000,
    "currency": "CLP"
  }'
```

### Invalid Beer Data (Missing Name)
```bash
curl -X POST http://localhost:8080/api/v1/beers \
  -H "Content-Type: application/json" \
  -d '{
    "id": 10,
    "brewery": "Test Brewery",
    "country": "Chile",
    "price": 1000,
    "currency": "CLP"
  }'
```

### Invalid Beer Data (Negative Price)
```bash
curl -X POST http://localhost:8080/api/v1/beers \
  -H "Content-Type: application/json" \
  -d '{
    "id": 11,
    "name": "Invalid Beer",
    "brewery": "Test Brewery",
    "country": "Chile",
    "price": -100,
    "currency": "CLP"
  }'
```

### Invalid Currency
```bash
curl -X POST http://localhost:8080/api/v1/beers \
  -H "Content-Type: application/json" \
  -d '{
    "id": 12,
    "name": "Invalid Currency Beer",
    "brewery": "Test Brewery",
    "country": "Chile",
    "price": 1000,
    "currency": "INVALID"
  }'
```

### Invalid Box Price Calculation (Zero Quantity)
```bash
curl -X GET "http://localhost:8080/api/v1/beers/1/boxprice?quantity=0&currency=CLP" \
  -H "Accept: application/json"
```

### Invalid Currency for Box Price
```bash
curl -X GET "http://localhost:8080/api/v1/beers/1/boxprice?quantity=6&currency=INVALID" \
  -H "Accept: application/json"
```

## Legacy Routes (Backward Compatibility)

### Create Beer (Legacy Route)
```bash
curl -X POST http://localhost:8080/beers \
  -H "Content-Type: application/json" \
  -d '{
    "id": 100,
    "name": "Legacy Beer",
    "brewery": "Legacy Brewery",
    "country": "Chile",
    "price": 800,
    "currency": "CLP"
  }'
```

### Get All Beers (Legacy Route)
```bash
curl -X GET http://localhost:8080/beers \
  -H "Accept: application/json"
```

### Get Beer by ID (Legacy Route)
```bash
curl -X GET http://localhost:8080/beers/100 \
  -H "Accept: application/json"
```

### Calculate Box Price (Legacy Route)
```bash
curl -X GET "http://localhost:8080/beers/100/boxprice?quantity=6&currency=CLP" \
  -H "Accept: application/json"
```

## Testing Scripts

### Quick Test All Endpoints
```bash
#!/bin/bash

echo "=== Beer API Test Script ==="
echo

echo "1. Health Check"
curl -s http://localhost:8080/ping | jq .
echo

echo "2. Create Test Beer"
curl -s -X POST http://localhost:8080/api/v1/beers \
  -H "Content-Type: application/json" \
  -d '{
    "id": 999,
    "name": "Test Beer",
    "brewery": "Test Brewery",
    "country": "Chile",
    "price": 1000,
    "currency": "CLP"
  }' | jq .
echo

echo "3. Get All Beers"
curl -s http://localhost:8080/api/v1/beers | jq .
echo

echo "4. Get Beer by ID"
curl -s http://localhost:8080/api/v1/beers/999 | jq .
echo

echo "5. Calculate Box Price"
curl -s "http://localhost:8080/api/v1/beers/999/boxprice?quantity=6&currency=USD" | jq .
echo

echo "=== Test Complete ==="
```

### Load Test (Create Multiple Beers)
```bash
#!/bin/bash

for i in {1000..1010}; do
  echo "Creating beer $i"
  curl -s -X POST http://localhost:8080/api/v1/beers \
    -H "Content-Type: application/json" \
    -d "{
      \"id\": $i,
      \"name\": \"Beer $i\",
      \"brewery\": \"Brewery $i\",
      \"country\": \"Chile\",
      \"price\": $((RANDOM % 2000 + 500)),
      \"currency\": \"CLP\"
    }" | jq .
done
```

## Expected Response Examples

### Successful Beer Creation (201)
```json
{
  "id": 1,
  "name": "Corona Extra",
  "brewery": "Cervecería Modelo",
  "country": "Mexico",
  "price": 1200,
  "currency": "CLP",
  "created_at": "2025-08-24T10:30:00Z",
  "updated_at": "2025-08-24T10:30:00Z"
}
```

### Box Price Response (200)
```json
{
  "beer_id": 1,
  "beer_name": "Corona Extra",
  "quantity": 6,
  "unit_price": 1200,
  "unit_currency": "CLP",
  "total_price": 7200,
  "target_currency": "CLP",
  "exchange_rate": 1,
  "converted_total": 7200
}
```

### Error Response (400)
```json
{
  "error": "VALIDATION_ERROR",
  "message": "Price must be greater than 0",
  "code": "INVALID_PRICE"
}
```

### Not Found Response (404)
```json
{
  "error": "NOT_FOUND",
  "message": "Beer with ID 999 not found"
}
```
