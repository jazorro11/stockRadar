# StockRadar Backend

Este backend está desarrollado en Go y expone una API REST para consultar recomendaciones bursátiles enriquecidas con datos cuantitativos y fundamentales. Utiliza la API de Finnhub para obtener métricas financieras y almacena los resultados en una base de datos CockroachDB.

---

## Requisitos previos

- Go 1.20 o superior
- CockroachDB (puedes usar un cluster local o en la nube)
- Acceso a la API de Finnhub (requiere API Key)
- Archivo `.env` con las siguientes variables:
  ```env
  API_URL=https://<tu_api_de_recomendaciones>
  API_KEY=<tu_api_key_de_recomendaciones>
  FINNHUB_API_KEY=<tu_api_key_de_finnhub>
  COCKROACHDB_URL=postgresql://<usuario>:<contraseña>@<host>:<puerto>/<db>?sslmode=disable
---

## Instalación y ejecución local

1. **Clona el repositorio y entra a la carpeta del backend:**
 ```sh
 git clone <url-del-repo>
 cd StockRadar/backend
 ```
2. **Instala las dependencias:**
 ```sh
 go mod tidy
 ```
3. **Configura el archivo `.env`:**
    - Completa las variables necesarias, como `FINNHUB_API_KEY` y los detalles de conexión a CockroachDB.
4. **Ejecuta las migraciones de la base de datos:**
 ```sh
 migrate -path db/migrations -database "cockroachdb://<usuario>:<password>@<host>:<puerto>/<nombre-db>?sslmode=disable" up
 ```
5. **Inicia el servidor:**
 ```sh
 go run main.go
 ```
6. **Accede a la API:**
   - La API debería estar corriendo en `http://localhost:8080`
   - Puedes usar herramientas como Postman o curl para interactuar con ella.

---
## ¿Por qué el servidor está alojado localmente?

**El backend se ejecuta localmente** para facilitar el desarrollo, pruebas y depuración.  
Puedes conectarlo fácilmente con el frontend (por ejemplo, en Vue) y con tu base de datos CockroachDB, ya sea local o en la nube.  
En producción, puedes desplegarlo en cualquier servidor compatible con Go (VPS, Docker, servicios cloud, etc).

---
## Comandos útiles
Instalar dependencias:
 ```sh
 go mod tidy
 ```
Ejecutar el servidor:
 ```sh
go run main.go
 ```
Actualizar dependencias:
 ```sh
go get -u
 ```
---
## Diccionario

### Funciones principales

| Función | Parámetros | Output | Descripción |
|--------|------------|--------|-------------|
| `main()` | Ninguno | Ninguno *(efectos colaterales: ejecuta todo el flujo del backend)* | Orquesta la carga de variables de entorno, consulta la API externa, enriquece los datos con Finnhub, calcula el score, almacena en CockroachDB y expone la API REST. |
| `enrichWithFinnhub(stock *StockInfo) error` | `stock`: puntero a `StockInfo` | `error`: si falla la consulta a Finnhub | Enriquece los datos con información de Finnhub: precio actual, market cap, EPS, P/E, P/B, dividend yield, 52W high/low, revenue growth, net profit margin y beta. |
| `scoreStock(stock StockInfo) float64` | `stock`: estructura con campos numéricos | `float64`: puntaje calculado | Calcula un puntaje cuantitativo para cada acción según sus métricas. |
| `getStocksHandler(w http.ResponseWriter, r *http.Request)` | `w`: ResponseWriter, `r`: Request | Ninguno *(responde con JSON)* | Expone los datos de acciones vía API REST en `/stocks`. |
| `parseDollarString(val string) (float64, error)` | `val`: string con `$` y/o comas | `float64`: valor numérico <br> `error`: si la conversión falla | Convierte strings con símbolos de dólar y comas a `float64`. |
| `UnmarshalJSON` *(método de `StockInfo`)* | `data`: bytes JSON | `error`: si la conversión falla | Convierte `target_from` y `target_to` de string con `$` a `float64` al deserializar JSON. |

### Estructuras principales

---

### `StockInfo`

Representa una acción bursátil con todos los campos relevantes.

| Campo                     | Tipo     | Descripción                                      |
|--------------------------|----------|--------------------------------------------------|
| `ticker`                 | string   | Símbolo bursátil de la acción                    |
| `company`                | string   | Nombre de la empresa                             |
| `brokerage`              | string   | Nombre del bróker que reporta la recomendación   |
| `action`                 | string   | Tipo de recomendación (ej: "Buy", "Hold", etc.)  |
| `rating_from`            | string   | Calificación anterior del bróker                 |
| `rating_to`              | string   | Nueva calificación otorgada                      |
| `target_from`            | float64  | Precio objetivo anterior                         |
| `target_to`              | float64  | Precio objetivo nuevo                            |
| `time` / `stock_time`    | string   | Momento de la recomendación (formato ISO 8601)   |
| `market_cap`             | float64  | Capitalización de mercado                        |
| `eps_ttm`                | float64  | Ganancias por acción (últimos 12 meses)          |
| `pe_ttm`                 | float64  | Relación precio/ganancias (últimos 12 meses)     |
| `pb`                     | float64  | Relación precio/valor contable                   |
| `dividend_yield`         | float64  | Rendimiento por dividendo                        |
| `week_52_high`           | float64  | Máximo de 52 semanas                             |
| `week_52_low`            | float64  | Mínimo de 52 semanas                             |
| `revenue_growth_ttm_yoy`| float64  | Crecimiento de ingresos interanual               |
| `net_profit_margin_ttm` | float64  | Margen de beneficio neto                         |
| `beta`                   | float64  | Beta (volatilidad relativa)                      |
| `current_price`          | float64  | Precio actual de la acción                       |
| `score`                  | float64  | Puntaje calculado según criterios cuantitativos  |

---

### `ApiResponse`

| Campo       | Tipo         | Descripción                                 |
|-------------|--------------|---------------------------------------------|
| `items`     | []StockInfo  | Lista de acciones recibidas del API         |
| `next_page` | string       | Token de paginación (opcional)              |

### Módulos y paquetes usados

| Paquete | Descripción |
|--------|-------------|
| `encoding/json` | Serialización y deserialización de datos JSON. |
| `net/http` | Cliente y servidor HTTP. |
| `os` | Manejo de variables de entorno y utilidades del sistema operativo. |
| `strconv`, `strings` | Conversión y manipulación de cadenas y números. |
| [`github.com/joho/godotenv`](https://github.com/joho/godotenv) | Carga de variables de entorno desde archivos `.env`. |
| [`github.com/jackc/pgx/v5`](https://github.com/jackc/pgx) | Driver para conectar y operar con bases de datos PostgreSQL y CockroachDB. |
| [`finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"`](https://github.com/Finnhub-Stock-API/finnhub-go) | SDK para consumir la API de Finnhub y obtener datos financieros. |

---
## Endpoints

### `GET /stocks`

- **Descripción:** Devuelve un arreglo JSON con todas las acciones almacenadas y sus métricas cuantitativas.
- **Formato de respuesta:**
  ```json
  [
    {
      "ticker": "AKBA",
      "company": "Akebia Therapeutics",
      "brokerage": "HC Wainwright",
      "action": "initiated by",
      "rating_from": "Buy",
      "rating_to": "Buy",
      "target_from": 8,
      "target_to": 8,
      "time": "2025-06-04T19:30:05-05:00",
      "market_cap": 934.98,
      "eps_ttm": -0.2151,
      "pe_ttm": 0,
      "pb": 20.03,
      "dividend_yield": 0,
      "week_52_high": 4.08,
      "week_52_low": 1.07,
      "revenue_growth_ttm_yoy": -1.23,
      "net_profit_margin_ttm": -24.51,
      "beta": 1.06,
      "current_price": 3.53,
      "score": 119.97
    }
  ]

---
## Despliegue en producción

Para desplegar este backend en un entorno de producción, sigue estos pasos:

1. **Configura un servidor con Go, CockroachDB y acceso a la API de Finnhub.**
2. **Clona el repositorio y repite los pasos de instalación, ajustando las configuraciones necesarias para producción.**
3. **Asegúrate de que las variables de entorno en el archivo `.env` estén correctamente configuradas para el entorno de producción.**
4. **Configura un servicio como systemd o supervisord para mantener el backend corriendo.**
5. **Configura un proxy inverso como Nginx o Apache para manejar las solicitudes HTTP y HTTPS.**
6. **Asegúrate de que el firewall del servidor permita el tráfico en los puertos necesarios (por defecto, 80 y 443 para HTTP/HTTPS).**