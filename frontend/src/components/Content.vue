<script setup>
import Search from './Search.vue'
import StockCard from './StockCard.vue'
import { useStocksStore } from '@/stores/stocks'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

const foundStock = ref(null)
const notFoundTicker = ref(null)
const store = useStocksStore()
const router = useRouter()
store.fetchStocks()

// Obtener las 3 acciones con mayor normalized
const topStocks = computed(() => {
  return [...store.stocks]
    .sort((a, b) => b.normalized - a.normalized)
    .slice(0, 3)
})


function handleNotFound(ticker) {
  foundStock.value = null
  notFoundTicker.value = ticker
  console.log('Stock no encontrado:', ticker)
}
function handleFound(stock) {
  foundStock.value = stock
  notFoundTicker.value = null
  console.log('Stock encontrado:', stock)
  router.push({ name: 'StockDetails', params: { ticker: stock.ticker } })
}

</script>

<template>
  <main class="flex flex-col items-center gap-[96px] px-[104px] pt-[96px] pb-[136px] w-full"
    style="background: radial-gradient(193.34% 155.87% at 50% -1.55%, rgba(250,254,252,0.70) 31.17%, rgba(18,156,82,0.11) 100%), #FAFEFC;">
    <p class="text-[#1B2821] text-center font-roboto text-[80px] font-semibold leading-[88px] tracking-[-2px]">
      Análisis técnico profesional para cada acción
    </p>
    <Search :is="true" :stocks="store.stocks" @found="handleFound" @notfound="handleNotFound" />

    <p class="text-[#1B2821] font-roboto text-[56px] font-semibold leading-[52px] tracking-[-1.12px]">
      Recomendaciones de inversión para hoy
    </p>
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 w-full">
      <StockCard v-for="stock in topStocks" :key="stock.ticker" :stock="stock" />
      <div v-if="topStocks.length === 0">Cargando o no encontrado</div>
    </div>
  </main>
</template>