// Importa las funciones necesarias de vue-router para crear el enrutador y el historial de navegación.
import { createRouter, createWebHistory } from 'vue-router'
// Importa el tipo RouteRecordRaw para tipar correctamente las rutas.
import type { RouteRecordRaw } from 'vue-router'
// Importa el componente Content, que será la vista principal o de inicio.
import Content from '@/components/Content.vue'
// Importa el componente StockDetails, que muestra los detalles de una acción.
import StockDetails from '@/components/StockDetails.vue'

// Define el arreglo de rutas de la aplicación, tipado como Array<RouteRecordRaw>.
const routes: Array<RouteRecordRaw> = [
  // Ruta raíz ("/"), renderiza el componente Content.
  { path: '/', component: Content },
  // Ruta dinámica para mostrar detalles de un stock según el ticker en la URL.
  { path: '/:ticker', name: 'StockDetails', component: StockDetails }
  // Ruta alternativa comentada, podría usarse para prefijar con /stock/ en vez de solo /:ticker.
  // { path: '/stock/:ticker', name: 'StockDetails', component: StockDetails }
]

// Crea la instancia del enrutador usando el historial HTML5 y las rutas definidas.
const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Exporta el enrutador para que pueda ser usado en la aplicación principal.
export default router