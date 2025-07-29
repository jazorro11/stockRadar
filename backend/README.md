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
## Diccionario de funciones y módulos

### Funciones principales

---

#### `main()`

- **Parámetros:** ninguno  
- **Output:** ninguno *(efectos colaterales: ejecuta todo el flujo del backend)*  
- **Descripción:**  
  Orquesta la carga de variables de entorno, consulta la API externa, enriquece los datos con Finnhub, calcula el score, almacena en CockroachDB y expone la API REST.

---

#### `enrichWithFinnhub(stock *StockInfo) error`

- **Parámetros:**
  - `stock`: puntero a una estructura `StockInfo` que será enriquecida con datos de Finnhub.
- **Output:**
  - `error`: retorna un error si falla la consulta a Finnhub, en caso contrario retorna `nil`.
- **Descripción:**  
  Enriquecimiento de cada acción con datos de Finnhub: precio actual, market cap, EPS, P/E, P/B, dividend yield, 52W high/low, revenue growth, net profit margin y beta.

---

#### `scoreStock(stock StockInfo) float64`

- **Parámetros:**
  - `stock`: estructura `StockInfo` con todos los campos numéricos relevantes.
- **Output:**
  - `float64`: el puntaje calculado para la acción.
- **Descripción:**  
  Calcula un puntaje cuantitativo para cada acción usando todos los campos numéricos relevantes.

---

#### `getStocksHandler(w http.ResponseWriter, r *http.Request)`

- **Parámetros:**
  - `w`: `ResponseWriter` de HTTP para enviar la respuesta.
  - `r`: `Request` de HTTP recibido.
- **Output:** ninguno *(efectos colaterales: responde con JSON)*  
- **Descripción:**  
  Handler HTTP que expone los datos almacenados en la base de datos vía API REST en `/stocks`.

---

#### `parseDollarString(val string) (float64, error)`

- **Parámetros:**
  - `val`: string que puede contener un valor numérico con símbolo `$` y/o comas.
- **Output:**
  - `float64`: valor numérico convertido.
  - `error`: error si la conversión falla.
- **Descripción:**  
  Convierte strings con símbolo `$` a `float64`.

---

#### `UnmarshalJSON` (método de `StockInfo`)

- **Parámetros:**
  - `data`: bytes JSON a deserializar.
- **Output:**
  - `error`: error si la conversión falla.
- **Descripción:**  
  Método personalizado para convertir los campos `target_from` y `target_to` de string (con `$`) a `float64` al deserializar el JSON.

---
## Módulos y paquetes usados

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