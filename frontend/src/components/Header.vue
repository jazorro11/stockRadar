<script setup>
import HeaderHome from './HeaderHome.vue'
import Search from './Search.vue';
import { useStocksStore } from '@/stores/stocks'
import { useRouter } from 'vue-router'
import { ref} from 'vue'

const foundStock = ref(null)
const notFoundTicker = ref(null)
const store = useStocksStore()
const router = useRouter()
store.fetchStocks()

function handleNotFound(ticker) {
  foundStock.value = null
  notFoundTicker.value = ticker
  console.log('Stock HEADER no encontrado:', ticker)
}
function handleFound(stock) {
  foundStock.value = stock
  notFoundTicker.value = null
  console.log('Stock HEADER encontrado:', stock)
  router.push({ name: 'StockDetails', params: { ticker: stock.ticker } })
}
</script>

<template>
  <header class="flex items-center justify-between px-12 py-3 border-b-[1.5px] border-[#CCEBDB] bg-white w-full">
    <HeaderHome />
     <Search :is="false" :stocks="store.stocks" @found="handleFound" @notfound="handleNotFound" />
  </header>
</template>
