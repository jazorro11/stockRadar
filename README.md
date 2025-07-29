# StockRadar

**StockRadar** es una plataforma para el análisis y visualización de recomendaciones bursátiles, combinando datos de brokers y métricas cuantitativas/fundamentales de cada acción. El proyecto está dividido en dos grandes módulos: **backend** y **frontend**, cada uno ubicado en su propia carpeta y con su propio README para detalles técnicos específicos.

---

## Estructura del proyecto

StockRadar/
├── backend/ # API REST en Go, integración con Finnhub y CockroachDB
│ └── README.md
├── frontend/ # Aplicación web en Vue 3 + Pinia + Tailwind CSS
│ └── README.md
└── README.md # (este archivo)

---

## ¿Qué contiene cada módulo?

### Backend (`/backend`)

- API REST desarrollada en Go.
- Integración con la API de Finnhub para enriquecer los datos de acciones.
- Almacenamiento de datos en CockroachDB.
- Lógica de cálculo de **score cuantitativo** para cada acción.
- Endpoint principal: `/stocks` para consultar todas las acciones y sus métricas.
- Documentación técnica y despliegue en [`/backend/README.md`](./backend/README.md).

### Frontend (`/frontend`)

- Aplicación web desarrollada con **Vue 3**.
- Gestión de estado con **Pinia**.
- Estilos responsivos con **Tailwind CSS**.
- Visualización de la tabla de acciones con filtros y ordenamiento.
- Consumo del endpoint `/stocks` del backend.
- Detalles técnicos y de desarrollo en [`/frontend/README.md`](./frontend/README.md).

---

## ¿Cómo correr el proyecto?

### 1. Clona el repositorio

```bash
git clone <url-del-repo>
cd StockRadar
```
### 2. Configura y ejecuta el backend

Sigue las instrucciones en [`/backend/README.md`](./backend/README.md) para:

- Instalar dependencias
- Configurar variables de entorno
- Ejecutar el servidor Go

---

### 3. Configura y ejecuta el frontend

Sigue las instrucciones en [`/frontend/README.md`](./frontend/README.md) para:

- Instalar dependencias
- Levantar el servidor de desarrollo de Vue

---

### 4. Accede a la aplicación

- **Frontend:** [http://localhost:5173](http://localhost:5173)
- **Backend:** [http://localhost:8080/stocks](http://localhost:8080/stocks)


