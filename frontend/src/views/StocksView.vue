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
      <div class="flex">
        <!-- Sidebar de filtros -->
        <aside class="w-64 mr-6 p-4 bg-gray-50 border rounded">
          <h2 class="font-bold mb-2">Filtros</h2>
          <!-- Filtro por tipo de acción -->
          <label class="block mb-2">
            <span class="text-sm">Tipo de acción:</span>
            <select v-model="store.filterAction" class="input input-bordered w-full">
              <option value="">Todas</option>
              <option value="Buy">Buy</option>
              <option value="Hold">Hold</option>
              <option value="Sell">Sell</option>
              <!-- Agrega más según tus datos -->
            </select>
          </label>
          <!-- Filtro por dividendos -->
          <label class="block mb-2">
            <span class="text-sm">Dividend Yield mínimo:</span>
            <input type="number" v-model.number="store.filterDividend" class="input input-bordered w-full" min="0" step="0.01" />
          </label>
          <!-- Filtro por beta -->
          <label class="block mb-2">
            <span class="text-sm">Beta máximo:</span>
            <input type="number" v-model.number="store.filterBeta" class="input input-bordered w-full" min="0" step="0.01" />
          </label>
          <!-- Filtro por score -->
          <label class="block mb-2">
            <span class="text-sm">Score mínimo:</span>
            <input type="number" v-model.number="store.filterScore" class="input input-bordered w-full" min="0" step="0.01" />
          </label>
        </aside>

        <!-- Contenido principal: tabla -->
        <div class="flex-1">
          <table class="min-w-full table-auto border">
        <thead>
          <tr>
            <th @click="store.setSort('ticker')" class="cursor-pointer">Ticker</th>
            <th @click="store.setSort('company')" class="cursor-pointer">Company</th>
            <th @click="store.setSort('brokerage')" class="cursor-pointer">Brokerage</th>
            <th @click="store.setSort('action')" class="cursor-pointer">Action</th>
            <th @click="store.setSort('rating_from')" class="cursor-pointer">From</th>
            <th @click="store.setSort('rating_to')" class="cursor-pointer">To</th>
            <th @click="store.setSort('target_from')" class="cursor-pointer">Target From</th>
            <th @click="store.setSort('target_to')" class="cursor-pointer">Target To</th>
            <th @click="store.setSort('time')" class="cursor-pointer">Time</th>
            <th @click="store.setSort('market_cap')" class="cursor-pointer">Market Cap</th>
            <th @click="store.setSort('eps_ttm')" class="cursor-pointer">EPS (TTM)</th>
            <th @click="store.setSort('pe_ttm')" class="cursor-pointer">P/E (TTM)</th>
            <th @click="store.setSort('pb')" class="cursor-pointer">P/B</th>
            <th @click="store.setSort('dividend_yield')" class="cursor-pointer">Dividend Yield</th>
            <th @click="store.setSort('week_52_high')" class="cursor-pointer">52W High</th>
            <th @click="store.setSort('week_52_low')" class="cursor-pointer">52W Low</th>
            <th @click="store.setSort('revenue_growth_ttm_yoy')" class="cursor-pointer">Revenue Growth YoY</th>
            <th @click="store.setSort('net_profit_margin_ttm')" class="cursor-pointer">Net Profit Margin</th>
            <th @click="store.setSort('beta')" class="cursor-pointer">Beta</th>
            <th @click="store.setSort('current_price')" class="cursor-pointer">Current Price</th>
            <th @click="store.setSort('score')" class="cursor-pointer">Score</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="stock in store.filteredStocks" :key="stock.ticker + stock.time">
            <td>{{ stock.ticker }}</td>
            <td>{{ stock.company }}</td>
            <td>{{ stock.brokerage }}</td>
            <td>{{ stock.action }}</td>
            <td>{{ stock.rating_from }}</td>
            <td>{{ stock.rating_to }}</td>
            <td>{{ stock.target_from }}</td>
            <td>{{ stock.target_to }}</td>
            <td>{{ new Date(stock.time).toLocaleString() }}</td>
            <td>{{ stock.market_cap }}</td>
            <td>{{ stock.eps_ttm }}</td>
            <td>{{ stock.pe_ttm }}</td>
            <td>{{ stock.pb }}</td>
            <td>{{ stock.dividend_yield }}</td>
            <td>{{ stock.week_52_high }}</td>
            <td>{{ stock.week_52_low }}</td>
            <td>{{ stock.revenue_growth_ttm_yoy }}</td>
            <td>{{ stock.net_profit_margin_ttm }}</td>
            <td>{{ stock.beta }}</td>
            <td>{{ stock.current_price }}</td>
            <td class="font-bold">{{ stock.score }}</td>
          </tr>
        </tbody>
      </table>
          <!-- ...tu tabla aquí... -->
        </div>
      </div>
      

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
