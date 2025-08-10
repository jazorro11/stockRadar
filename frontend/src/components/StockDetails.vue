<script setup lang="ts">
// Importa el componente DetailCard para mostrar detalles financieros y de mercado.
import DetailCard from './DetailCard.vue'
// Importa useRoute para acceder a los parámetros de la ruta (por ejemplo, el ticker).
import { useRoute } from 'vue-router'
// Importa el store de acciones para acceder a la lista de stocks.
import { useStocksStore } from '@/stores/stocks'
// Importa ref y watch de Vue para reactividad y observación de cambios.
import { ref, watch, onMounted } from 'vue'
// import AllDetails from '@/components/AllDetails.vue' // (comentado, posible uso futuro)

// Obtiene la instancia de la ruta actual.
const route = useRoute()
// Obtiene la instancia del store de acciones.
const store = useStocksStore()
// const ticker = route.params.ticker as string // (comentado, ejemplo de obtención de ticker)
// const stock = (store.stocks as Stock[]).find(s => s.ticker === ticker) // (comentado, ejemplo de búsqueda)
// Define una referencia reactiva para el stock seleccionado.
const stock = ref<Stock | undefined>(undefined)

// Define el tipo Stock con todas las propiedades relevantes.
type Stock = {
  ticker: string
  company: string
  brokerage: string
  action: string
  rating_from: string
  rating_to: string
  target_from: number
  target_to: number
  time: string
  market_cap: number
  eps_ttm: number
  pe_ttm: number
  pb: number
  dividend_yield: number
  week_52_high: number
  week_52_low: number
  revenue_growth_ttm_yoy: number
  net_profit_margin_ttm: number
  beta: number
  current_price: number
  score: number
  normalized: number
  type: string
}

// Función que actualiza el stock seleccionado según el ticker de la ruta.
function updateStock() {
  stock.value = (store.stocks as Stock[]).find(s => s.ticker === (route.params.ticker as string))
}
// Nueva función: asegura que haya datos (recarga) y luego selecciona
async function ensureDataAndSelect() {
  if (store.stocks.length === 0) {
    await store.fetchStocks()
  }
  updateStock()
}

// // Llama a updateStock al cargar el componente.
// updateStock()
// // Observa cambios en el parámetro ticker de la ruta y actualiza el stock cuando cambie.
// watch(() => route.params.ticker, updateStock)
// Reemplaza la llamada directa anterior:
// updateStock()
onMounted(() => {
  ensureDataAndSelect()
})

// Observa cambios en el parámetro ticker y re-evalúa (por si navegación directa)
watch(() => route.params.ticker, () => ensureDataAndSelect())

// Devuelve una recomendación textual según el valor normalizado.
function getRecommendation(normalized: number): string {
  if (normalized < 0.33) return 'Venta Fuerte'
  if (normalized < 0.66) return 'Mantener'
  return 'Compra Fuerte'
}

// Devuelve un color según el valor normalizado (rojo, amarillo o verde).
function getColor(normalized: number): string {
  if (normalized < 0.33) return '#ef4444' // rojo
  if (normalized < 0.66) return '#facc15' // amarillo
  return '#22c55e' // verde
}

// Convierte un texto a Title Case (primera letra de cada palabra en mayúscula).
function toTitleCase(text?: string): string {
  if (!text) return ''
  return text
    .toLowerCase()
    .split(' ')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ')
}
</script>

<template>
  <!-- Contenedor principal -->
  <main class="flex flex-col items-center gap-[96px] px-[104px] pt-[96px] pb-[136px] w-full"
    style="background: radial-gradient(193.34% 155.87% at 50% -1.55%, rgba(250,254,252,0.70) 31.17%, rgba(18,156,82,0.11) 100%), #FAFEFC;">

    <!-- Estado de carga -->
    <div v-if="store.loading" class="text-gray-500">Cargando datos...</div>

    <!-- Cabecera si existe el stock -->
    <div v-else-if="stock" class="w-full flex flex-col ">
      <!-- ...existing code (cabecera del stock)... -->
      <div class="flex items-center gap-4 self-stretch">
        <div class="flex flex-col items-start gap-4 flex-1">
          <div class="flex items-center gap-2">
            <p
              class="flex flex-col justify-center w-[128px] h-[42px] flex-shrink-0 overflow-hidden text-[#1B2821] text-ellipsis whitespace-nowrap font-inter text-[48px] font-semibold leading-[24px]">
              {{ stock?.ticker }}
            </p>
            <p class="text-[#4B4A5C] font-inter text-[15px] font-medium leading-[24px] overflow-hidden text-ellipsis whitespace-normal"
              style="display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 1; line-clamp: 1;">
              {{ toTitleCase(stock?.company) }}
            </p>
          </div>
          <p class="self-stretch text-[#1B2821] font-inter text-[15px] font-normal leading-[20px] tracking-[-0.075px] overflow-hidden text-ellipsis whitespace-normal"
            style="display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 1; line-clamp: 1;">
            {{ stock?.type }}
          </p>
        </div>
        <div
          class="flex items-center justify-center min-w-[180px] px-4 h-[56px] bg-white rounded-xl text-[#0E7B41] font-inter text-[36px] font-semibold leading-[48px] tracking-[-0.23px] overflow-hidden text-ellipsis whitespace-nowrap">
          $ {{ Number(stock?.current_price).toFixed(2) }}
        </div>
      </div>
    </div>

    <!-- Mensaje si no se encontró el stock tras cargar -->
    <div v-else>
      <p class="text-red-600">No se encontró información para este ticker.</p>
    </div>

    <!-- Sección análisis técnico (solo si hay stock) -->
    <div v-if="stock" class="w-full flex flex-row gap-8 self-stretch p-8 rounded-2xl shadow"
      style="background:rgba(255,255,255,0.5);">
      <!-- ...existing code (análisis técnico + barra)... -->
      <div class="flex flex-col items-start gap-6 flex-1 justify-center">
        <p class="text-[#1B2821] font-inter text-[48px] font-semibold leading-[48px]">
          Análisis técnico
        </p>
        <p class="self-stretch text-[#1B2821] font-inter text-[15px] font-normal leading-[20px] tracking-[-0.075px] overflow-hidden text-ellipsis whitespace-normal"
          style="display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 1; line-clamp: 1;">
          Broker: {{ stock?.brokerage }}
        </p>
      </div>
      <div class="flex flex-col justify-end items-end gap-6">
        <div class="w-full flex flex-col items-center my-8">
          <div class=" max-w-xl h-4 rounded-full flex overflow-hidden w-[800px]">
            <div class="bg-red-500 h-full" style="width:33.33%"></div>
            <div class="bg-yellow-400 h-full" style="width:33.33%"></div>
            <div class="bg-green-500 h-full" style="width:33.34%"></div>
          </div>
          <div class="relative w-full max-w-xl" style="height: 12px;">
            <div class="absolute top-0 left-0 -mt-2"
              :style="{ left: `calc(${(stock?.normalized ?? 0) * 100}% - 12px)` }">
              <svg width="24" height="12">
                <polygon points="12,0 24,12 0,12" :fill="getColor(stock?.normalized ?? 0)" />
              </svg>
            </div>
          </div>
          <div class="text-center mt-2 text-2xl font-semibold text-gray-700">
            {{ getRecommendation(stock?.normalized ?? 0) }}
          </div>
        </div>
      </div>
    </div>

    <!-- Tarjetas de detalles (solo si hay stock) -->
    <div v-if="stock" class="grid grid-cols-1 md:grid-cols-2 gap-6 w-full">
      <DetailCard :financialDetail="true" :stock="stock" />
      <DetailCard :financialDetail="false" :stock="stock" />
    </div>

  </main>
</template>