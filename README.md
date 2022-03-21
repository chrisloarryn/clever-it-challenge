# Clever IT Challenge

## Especificaciones
Para ejecutar este aplicativo se necesita:

### Desarrollo
Solo es necesario el api token:
````shell
API_TOKEN=<YOUR_API_TOKEN>
````

### Producción
Se necesita una DB postgres configurada, por default toma los siguientes valores:
````shell
DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=postgres
DB_USER=postgres
DB_PASSWORD=postgres
````

## Problema

Bender es fanático de las cervezas y quiere tener un registro de todas las cervezas que prueba y como calcular el precio que necesita para comprar una caja de algún tipo especifico de cervezas. Para esto necesita una API REST con esta información que posteriormente compartirá con sus amigos.

### Descripción

Se solicita crear un API REST basándonos en la definición que se encuentra en el archivo **openapi.yaml**.

#### Funcionalidad

- GET /Beers: Lista todas las cervezas que se encuentran en el sistema.
- POST /Beers: Permite ingresar una nueva cerveza.
- GET /beers/{beerID}: Lista un detalle de una cerveza especifica.
- GET /beets/{beerID}/boxprice: Entrega el valor que cuesta una caja específica de cerveza dependiendo de los parámetros ingresados, esto quiere decir que multiplique el precio por la cantidad una vez se homologara la moneda a lo que se ingreso por parámetro.
    - Quantity: Cantidad de cervezas a comprar (valor por defecto 6).
    - Currency: Tipo de moneda con la que desea pagar, para este caso se recomienda que utilice esta API <https://currencylayer.com/>

### Requisitos

- Puede usar alguno de los siguientes lenguajes Java, NodeJS, Go o Python. Aunque valoramos el uso de GO.
- Usar Docker y Docker Compose para los diferentes servicios.
- Se puede usar librarías externas y frameworks
- Requisito un 70% de cobertura de código
- Completa libertad para agregar nuevas funcionalidades.

### Entrega

- Enviar el link del repositorio donde se realiza este ejercicio.