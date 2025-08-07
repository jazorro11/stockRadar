<script setup lang="ts">
import DetailCard from './DetailCard.vue'
import { useRoute } from 'vue-router'
import { useStocksStore } from '@/stores/stocks'
import { ref, watch } from 'vue'
// import AllDetails from '@/components/AllDetails.vue'

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
  type: string
}

function updateStock() {
  stock.value = (store.stocks as Stock[]).find(s => s.ticker === (route.params.ticker as string))
}

updateStock()
watch(() => route.params.ticker, updateStock)

function getRecommendation(normalized: number): string {
  if (normalized < 0.33) return 'Venta Fuerte'
  if (normalized < 0.66) return 'Mantener'
  return 'Compra Fuerte'
}

function getColor(normalized: number): string {
  if (normalized < 0.33) return '#ef4444' // rojo
  if (normalized < 0.66) return '#facc15' // amarillo
  return '#22c55e' // verde
}

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
  <main class="flex flex-col items-center gap-[96px] px-[104px] pt-[96px] pb-[136px] w-full"
    style="background: radial-gradient(193.34% 155.87% at 50% -1.55%, rgba(250,254,252,0.70) 31.17%, rgba(18,156,82,0.11) 100%), #FAFEFC;">
    <div v-if="stock" class="w-full flex flex-col ">
      <div class="flex items-center gap-4 self-stretch">
        <div class="flex flex-col items-start gap-4 flex-1">
          <div class="flex items-center gap-2">
            <p
              class="flex flex-col justify-center w-[128px] h-[42px] flex-shrink-0 overflow-hidden text-[#1B2821] text-ellipsis whitespace-nowrap font-inter text-[48px] font-semibold leading-[24px]">
              {{ stock?.ticker }}
            </p>
            <p class="text-[#4B4A5C] font-inter text-[15px] font-medium leading-[24px] overflow-hidden text-ellipsis whitespace-normal"
              style="display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 1;">
              {{ toTitleCase(stock?.company) }}
            </p>
          </div>
          <p class="self-stretch text-[#1B2821] font-inter text-[15px] font-normal leading-[20px] tracking-[-0.075px] overflow-hidden text-ellipsis whitespace-normal"
            style="display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 1;">
            {{ stock?.type }}
          </p>
        </div>
        <div
          class="flex items-center justify-center min-w-[180px] px-4 h-[56px] bg-white rounded-xl text-[#0E7B41] font-inter text-[36px] font-semibold leading-[48px] tracking-[-0.23px] overflow-hidden text-ellipsis whitespace-nowrap">
          $ {{ Number(stock?.current_price).toFixed(2) }}
        </div>
      </div>
    </div>

    <div v-else>
      <p class="text-red-600">No se encontró información para este ticker.</p>
    </div>

    <div class="w-full flex flex-row gap-8 self-stretch p-8 rounded-2xl shadow"
      style="background:rgba(255,255,255,0.5);">

      <!-- Texto principal -->
      <div class="flex flex-col items-start gap-6 flex-1 justify-center">
        <p class="text-[#1B2821] font-inter text-[48px] font-semibold leading-[48px]">
          Análisis técnico
        </p>
        <p class="self-stretch text-[#1B2821] font-inter text-[15px] font-normal leading-[20px] tracking-[-0.075px] overflow-hidden text-ellipsis whitespace-normal"
          style="display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 1;">
          Broker: {{ stock?.brokerage }}
        </p>
      </div>
      <!-- Gráfica -->
      <div class="flex flex-col justify-end items-end gap-6">
        <div class="w-full flex flex-col items-center my-8">
          <!-- Barra de colores base -->
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
          <!-- Texto de recomendación -->
          <div class="text-center mt-2 text-2xl font-semibold text-gray-700">
            {{ getRecommendation(stock?.normalized ?? 0) }}
          </div>
        </div>
      </div>
    </div>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6 w-full">

      <!-- <AllDetails />  -->
      <DetailCard :financialDetail="true" :stock="stock" />
      <DetailCard :financialDetail="false" :stock="stock" />

    </div>




  </main>

</template>