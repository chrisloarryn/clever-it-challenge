# REST Client Testing Guide

Este directorio contiene archivos para probar la API usando la extensión REST Client de VS Code.

## Archivos Disponibles

### 📁 `test.REST`
Archivo principal con **todos los casos de prueba** de la API:
- ✅ Casos funcionales (CRUD completo)
- ❌ Casos de error (validaciones)
- 🔄 Rutas legacy (compatibilidad)
- 🧪 Casos edge (límites y performance)

### 📁 `test-variables.REST`
Archivo simplificado con **variables configurables** para diferentes entornos:
- 🏠 Desarrollo (localhost)
- 🚀 Staging (comentado)
- 🌍 Producción (comentado)

## Prerequisitos

### 1. Instalar REST Client Extension
```
Ext Name: REST Client
Publisher: Huachao Mao
VS Code Extension ID: humao.rest-client
```

### 2. Iniciar el Servidor
```bash
# Terminal 1: Iniciar servidor
make dev
# o
export DB_TYPE=inmemory && go run cmd/main.go
```

## Cómo Usar

### Opción 1: Ejecutar Todo el Test Suite
1. Abrir `test.REST` en VS Code
2. Click en **"Send Request"** en cada sección
3. Ver respuestas en panel derecho

### Opción 2: Pruebas Rápidas con Variables
1. Abrir `test-variables.REST` 
2. Modificar variables si es necesario:
   ```http
   @baseURL = http://localhost:8080
   @testBeerID = 999
   @testPrice = 1500
   ```
3. Ejecutar requests individuales

### Opción 3: Cambiar Entorno
Para probar en staging/producción:
1. Abrir `test-variables.REST`
2. Comentar variables de desarrollo
3. Descomentar variables del entorno deseado
4. Ejecutar tests

## Variables Disponibles

| Variable | Descripción | Ejemplo |
|----------|-------------|---------|
| `@baseURL` | URL base del servidor | `http://localhost:8080` |
| `@apiURL` | URL del API v1 | `{{baseURL}}/api/v1` |
| `@contentType` | Content-Type header | `application/json` |
| `@testBeerID` | ID para pruebas | `999` |
| `@testBrewery` | Cervecería de prueba | `Test Brewery` |
| `@testCountry` | País de prueba | `Chile` |
| `@testPrice` | Precio de prueba | `1500` |
| `@testCurrency` | Moneda de prueba | `CLP` |

## Estructura del Request

### Formato Correcto ✅
```http
### Descripción del Test
# @name nombreDelTest
GET {{baseURL}}/ping HTTP/1.1
Accept: {{contentType}}
```

### Formato con Body ✅
```http
### Crear Cerveza
# @name createBeer
POST {{apiURL}}/beers HTTP/1.1
Content-Type: {{contentType}}

{
  "id": {{testBeerID}},
  "name": "Test Beer",
  "brewery": "{{testBrewery}}",
  "country": "{{testCountry}}",
  "price": {{testPrice}},
  "currency": "{{testCurrency}}"
}
```

## Tests Incluidos

### 🟢 Funcionales (Happy Path)
- Health check
- Crear cervezas (múltiples monedas)
- Obtener todas las cervezas
- Obtener cerveza por ID
- Calcular precio de caja
- Conversión de monedas

### 🔴 Casos de Error
- ID duplicado (409 Conflict)
- Datos faltantes (400 Bad Request)
- Datos inválidos (400 Bad Request)
- Recurso no encontrado (404 Not Found)
- Moneda inválida (400 Bad Request)

### 🔄 Compatibilidad
- Rutas API v1: `/api/v1/beers`
- Rutas legacy: `/beers`
- Headers opcionales
- CORS testing

### 🧪 Edge Cases
- Valores mínimos/máximos
- Strings largos
- Cantidades grandes
- Diferentes content-types

## Respuestas Esperadas

### Success (2xx)
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

### Error (4xx)
```json
{
  "error": "VALIDATION_ERROR",
  "message": "Price must be greater than 0",
  "code": "INVALID_PRICE"
}
```

## Tips de Uso

### 🔄 Re-ejecutar Tests
- **Ctrl/Cmd + Shift + P** → "Rest Client: Send Request"
- Click en **"Send Request"** sobre cada endpoint

### 📋 Copiar Response
- Click derecho en response → "Copy Response"
- Usar response en otros tests

### 🔍 Histórico de Requests
- **Ctrl/Cmd + Shift + P** → "Rest Client: Request History"
- Ver requests anteriores

### ⚡ Shortcuts
- **Ctrl/Cmd + Alt + R**: Send request
- **Ctrl/Cmd + Alt + K**: Clear cache
- **Ctrl/Cmd + Alt + H**: Toggle history

## Troubleshooting

### ❌ "Failed to connect"
```
✅ Verificar que el servidor esté corriendo: make dev
✅ Verificar URL en @baseURL
✅ Verificar puerto (8080 por defecto)
```

### ❌ "404 Not Found"
```
✅ Verificar endpoint en la URL
✅ Verificar método HTTP (GET/POST)
✅ Verificar rutas configuradas en el servidor
```

### ❌ "400 Bad Request"
```
✅ Verificar JSON syntax
✅ Verificar headers requeridos
✅ Verificar datos de request
```

## Integración con Desarrollo

### Durante Desarrollo
1. Mantener `test.REST` abierto
2. Ejecutar tests después de cambios
3. Verificar responses esperadas

### Para CI/CD
```bash
# Usar script automatizado para CI
./scripts/test-api.sh
```

### Para Load Testing
```bash
# Usar Makefile para carga
make api-load-test
```
