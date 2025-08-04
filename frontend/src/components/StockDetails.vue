<script setup lang="ts">
import { useRoute } from 'vue-router'
import { useStocksStore } from '@/stores/stocks'
import { ref, watch } from 'vue'

const route = useRoute()
const store = useStocksStore()
// const ticker = route.params.ticker as string
// const stock = (store.stocks as Stock[]).find(s => s.ticker === ticker)
const stock = ref<Stock | undefined>(undefined)

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
}

function updateStock() {
  stock.value = (store.stocks as Stock[]).find(s => s.ticker === (route.params.ticker as string))
}

updateStock()
watch(() => route.params.ticker, updateStock)

function getRecommendation(normalized: number): string {
  if (normalized < 0.33) return 'Fuerte Venta'
  if (normalized < 0.66) return 'Mantener'
  return 'Fuerte Compra'
}

function getColor(normalized: number): string {
  if (normalized < 0.33) return '#ef4444' // rojo
  if (normalized < 0.66) return '#facc15' // amarillo
  return '#22c55e' // verde
}
</script>

<template>

  <div class="max-w-2xl mx-auto mt-12 p-8 bg-white rounded-2xl shadow">
    <h1 class="text-3xl font-bold mb-4">Detalles de {{ stock?.ticker }}</h1>
    <div v-if="stock">
      <p><strong>Brokerage:</strong> {{ stock?.brokerage }}</p>
      <p><strong>Acción:</strong> {{ stock?.action }}</p>
      <p><strong>Rating (de):</strong> {{ stock?.rating_from }}</p>
      <p><strong>Rating (a):</strong> {{ stock?.rating_to }}</p>
      <p><strong>Target (de):</strong> {{ stock?.target_from }}</p>
      <p><strong>Target (a):</strong> {{ stock?.target_to }}</p>
      <p><strong>Fecha:</strong> {{ stock?.time }}</p>
      <p><strong>Market Cap:</strong> {{ stock?.market_cap }}</p>
      <p><strong>EPS (TTM):</strong> {{ stock?.eps_ttm }}</p>
      <p><strong>P/E (TTM):</strong> {{ stock?.pe_ttm }}</p>
      <p><strong>P/B:</strong> {{ stock?.pb }}</p>
      <p><strong>Dividend Yield:</strong> {{ stock?.dividend_yield }}</p>
      <p><strong>52 Week High:</strong> {{ stock?.week_52_high }}</p>
      <p><strong>52 Week Low:</strong> {{ stock?.week_52_low }}</p>
      <p><strong>Crecimiento ingresos YoY (TTM):</strong> {{ stock?.revenue_growth_ttm_yoy }}</p>
      <p><strong>Margen neto (TTM):</strong> {{ stock?.net_profit_margin_ttm }}</p>
      <p><strong>Beta:</strong> {{ stock?.beta }}</p>
      <p><strong>Precio actual:</strong> {{ stock?.current_price }}</p>
      <p><strong>Score:</strong> {{ stock?.score }}</p>
      <p><strong>Normalized:</strong> {{ stock?.normalized }}</p>
    </div>
    <div v-else>
      <p class="text-red-600">No se encontró información para este ticker.</p>
    </div>
  </div>

  <div class="w-full flex flex-col items-center my-8">
    <!-- Barra de colores base -->
    <div class="w-full max-w-md h-4 rounded-full flex overflow-hidden">
      <div class="bg-red-500 h-full" style="width:33.33%"></div>
      <div class="bg-yellow-400 h-full" style="width:33.33%"></div>
      <div class="bg-green-500 h-full" style="width:33.34%"></div>
    </div>
    <!-- Indicador de posición -->
    <div class="relative w-full max-w-md" style="height: 12px;">
      <div class="absolute top-0 left-0 -mt-2" :style="{ left: `calc(${(stock?.normalized ?? 0) * 100}% - 12px)` }">
        <svg width="24" height="12">
          <polygon points="12,0 24,12 0,12" :fill="getColor(stock?.normalized ?? 0)" />
        </svg>
      </div>
    </div>
    <!-- Texto de recomendación -->
    <div class="text-center mt-2 text-2xl font-semibold text-gray-700">
      {{ getRecommendation(stock?.normalized ?? 0) }}
    </div>
  </div>
</template>