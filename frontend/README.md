# StockRadar Frontend

Este proyecto es el frontend de **StockRadar**, una aplicación web para visualizar información y análisis técnico de acciones. Está construido con [Vue 3](https://vuejs.org/), [Pinia](https://pinia.vuejs.org/) para gestión de estado, [Vue Router](https://router.vuejs.org/) para enrutamiento y utiliza [Vite](https://vitejs.dev/) como herramienta de desarrollo.

---

## Estructura del Proyecto

```
frontend/
│
├── public/                  # Archivos públicos (favicon, etc.)
├── src/
│   ├── assets/              # Imágenes y estilos globales
│   ├── components/          # Componentes Vue reutilizables
│   │   ├── __tests__/       # Pruebas unitarias de componentes
│   │   └── icons/           # Iconos SVG como componentes Vue
│   ├── router/
│   │   └── index.ts         # Definición de rutas de la aplicación
│   ├── stores/
│   │   └── stocks.ts        # Store Pinia para gestión de acciones
│   ├── App.vue              # Componente raíz
│   └── main.ts              # Punto de entrada de la aplicación
├── index.html               # HTML principal
├── package.json             # Dependencias y scripts
├── tsconfig*.json           # Configuración de TypeScript
├── vite.config.ts           # Configuración de Vite
└── README.md                # Documentación del proyecto
```

---

## Principales Componentes

- **App.vue**  
  Componente raíz. Incluye el header, el footer y el área principal donde se renderizan las vistas según la ruta.

- **Header.vue**  
  Barra superior con el logo y la barra de búsqueda.

- **HeaderHome.vue**  
  Logo y nombre de la aplicación, clickeable para volver al inicio.

- **Footer.vue**  
  Pie de página simple con información de copyright.

- **Content.vue**  
  Vista principal o de inicio, puede mostrar información general o destacados.

- **StockCard.vue**  
  Tarjeta individual para mostrar información resumida de una acción.

- **StockDetails.vue**  
  Vista de detalle de una acción, muestra información principal, análisis técnico, gráfica de recomendación y detalles financieros/mercado.

- **DetailCard.vue**  
  Tarjeta para mostrar indicadores financieros o de mercado de una acción.

- **Search.vue**  
  Barra de búsqueda de acciones con autocompletado y sugerencias.

- **ButtonClick.vue**  
  Botón estilizado reutilizable.

---

## Gestión de Estado

- **Pinia** se utiliza para el manejo global del estado.
- El store principal es `stores/stocks.ts`, que gestiona la lista de acciones (`stocks`), el estado de carga (`loading`) y los errores (`error`).
- La acción `fetchStocks` obtiene los datos desde el backend usando Axios.

---

## Enrutamiento

- **Vue Router** define dos rutas principales en `router/index.ts`:
  - `/` → Componente `Content.vue` (vista principal)
  - `/:ticker` → Componente `StockDetails.vue` (detalle de una acción según su ticker)

---

## Estilos

- Los estilos globales se encuentran en `src/assets/main.css` y `base.css`.
- Se utilizan clases utilitarias (por ejemplo, de Tailwind CSS) para el diseño de los componentes.

---

## Scripts Principales

- **main.ts**  
  Punto de entrada de la app. Crea la instancia de Vue, registra Pinia y Vue Router, y monta la aplicación.

---

## Cómo ejecutar el proyecto

1. **Instalar dependencias**
   ```bash
   npm install
   ```

2. **Iniciar el servidor de desarrollo**
   ```bash
   npm run dev
   ```

3. **Abrir en el navegador**
   - Accede a [http://localhost:5173](http://localhost:5173) (o el puerto que indique la terminal).

---

## Pruebas

- Las pruebas unitarias de componentes se encuentran en `src/components/__tests__/`.
- Se utiliza [Vitest](https://vitest.dev/) para ejecutar las pruebas.

---

## Notas adicionales

- El frontend espera que el backend esté corriendo en `http://localhost:8080` para obtener la lista de acciones.
- Puedes personalizar los estilos y la estructura de los componentes según tus necesidades.
- El proyecto está preparado para trabajar con TypeScript.

---

