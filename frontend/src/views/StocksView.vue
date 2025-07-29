<!--
  ---------------------------------------------------------------------------
  StocksView.vue
  ---------------------------------------------------------------------------
  Esta vista representa la interfaz principal del sistema de recomendaciones
  bursátiles. Su objetivo es mostrar una tabla con información detallada de 
  las acciones, incluyendo la compañía, bróker, calificación y objetivos de precio.

  Funcionalidades clave:
  - Recupera los datos desde el backend con Pinia (fetchStocks).
  - Permite buscar por ticker, empresa o bróker.
  - Permite ordenar por cualquier columna haciendo clic en el encabezado.
  - Muestra mensajes de carga o error si corresponde.

  Este componente depende del store de Pinia definido en `useStocksStore`,
  y utiliza Tailwind CSS para estilos rápidos y adaptables.
-->

<template>
  <div class="p-4">
    <!-- Título de la página -->
    <h1 class="text-2xl font-bold mb-4">Stock Recommendations</h1>

    <!-- Campo de búsqueda reactivo vinculado al store -->
    <input
      v-model="store.search"
      type="text"
      placeholder="Search by ticker, company, or brokerage"
      class="input input-bordered w-full mb-4"
    />

    <!-- Muestra mensaje de carga mientras se recuperan los datos -->
    <div v-if="store.loading" class="text-center">Loading...</div>

    <!-- Una vez cargados los datos, se muestra la tabla -->
    <div v-else>
      <table class="min-w-full table-auto border">
        <thead>
          <tr>
            <!-- Cada encabezado es clickeable para ordenar por esa columna -->
            <th @click="store.setSort('ticker')" class="cursor-pointer">Ticker</th>
            <th @click="store.setSort('company')" class="cursor-pointer">Company</th>
            <th @click="store.setSort('brokerage')" class="cursor-pointer">Brokerage</th>
            <th @click="store.setSort('action')" class="cursor-pointer">Action</th>
            <th @click="store.setSort('rating_from')" class="cursor-pointer">From</th>
            <th @click="store.setSort('rating_to')" class="cursor-pointer">To</th>
            <th @click="store.setSort('target_from')" class="cursor-pointer">Target From</th>
            <th @click="store.setSort('target_to')" class="cursor-pointer">Target To</th>
            <th @click="store.setSort('time')" class="cursor-pointer">Time</th>
          </tr>
        </thead>
        <tbody>
          <!-- Iteramos sobre la lista filtrada y ordenada del store -->
          <tr v-for="stock in store.filteredStocks" :key="stock.ticker + stock.time">
            <td>{{ stock.ticker }}</td>
            <td>{{ stock.company }}</td>
            <td>{{ stock.brokerage }}</td>
            <td>{{ stock.action }}</td>
            <td>{{ stock.rating_from }}</td>
            <td>{{ stock.rating_to }}</td>
            <td>{{ stock.target_from }}</td>
            <td>{{ stock.target_to }}</td>
            <!-- Mostramos la fecha en formato legible local -->
            <td>{{ new Date(stock.time).toLocaleString() }}</td>
          </tr>
        </tbody>
      </table>

      <!-- Si ocurre un error en la carga de datos, se muestra aquí -->
      <div v-if="store.error" class="text-red-500 mt-2">{{ store.error }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
// Importamos onMounted para ejecutar lógica al cargar la vista
import { onMounted } from 'vue'

// Importamos nuestro store de acciones bursátiles
import { useStocksStore } from '@/stores/stocks'

// Inicializamos el store
const store = useStocksStore()

// Al montar el componente, se hace fetch de los datos desde el backend
onMounted(() => store.fetchStocks())
</script>
