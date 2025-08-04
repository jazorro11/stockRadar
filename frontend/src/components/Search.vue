<script setup>
import { ref, computed } from 'vue'
import { defineProps, defineEmits } from 'vue'
import ButtonClick from './ButtonClick.vue'
import { useStocksStore } from '@/stores/stocks'

const store = useStocksStore()
store.fetchStocks()

const props = defineProps({
  is: {
    type: Boolean,
    default: true
  },
  stocks: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['found', 'notfound'])

const searchValue = ref('')

const filteredStocks = computed(() => {
  const query = searchValue.value.trim().toUpperCase()
  if (!query) return []
  const matches = store.stocks.filter(stock => stock.ticker.toUpperCase().startsWith(query))
  if (matches.length === 0) {
    return [{ ticker: '', company: 'No hay resultados', noResult: true }]
  }
  return matches
})

function handleKeyup(e) {
  if (e.key === 'Enter') {
    searchTicker()
    console.log('Enter pressed, searching for:', searchValue.value)
    searchValue.value = '' // Borra el input
  }
}

function handleButtonClick() {
  searchTicker()
  console.log('Click Handle Button')
}

function searchTicker() {
  const ticker = searchValue.value.trim().toUpperCase()
  const found = store.stocks.find(stock => stock.ticker.toUpperCase() === ticker)
  if (found) {
    console.log('Encontrado:', found)
    emit('found', found)
  } else {
    console.log('No encontrado:', ticker)
    emit('notfound', ticker)
  }
}

function selectStock(stock) {
  searchValue.value = stock.ticker
  searchTicker()
  searchValue.value = '' // Borra el input
}

</script>

<template>
  <div v-if="is" class="w-full flex flex-col items-center gap-[72px]">
    <div class="flex items-center p-6 gap-2.5 rounded-[40px] border border-[#91B6A2] bg-[#FAFEFC]  w-full">
      <img src="@/assets/search.svg" alt="Buscar" class="w-8 h-8 aspect-square" />

      <div class="relative w-full">
        <input v-model="searchValue" @keyup="handleKeyup" type="text" placeholder="Ingresa el Ticker de la acción"
          class="bg-transparent outline-none border-none text-[#91B6A2] text-[20px] font-normal leading-none w-full" />
        <ul v-if="searchValue"
          class="absolute left-0 top-full z-10 w-full bg-white border border-[#91B6A2] rounded-lg shadow max-h-60 overflow-auto">
          <li v-for="stock in filteredStocks" :key="stock.ticker + stock.company"
            class="px-4 py-2 hover:bg-[#f0fdf4] cursor-pointer" @click="!stock.noResult && selectStock(stock)">
            <span v-if="stock.noResult" class="text-gray-400">{{ stock.company }}</span>
            <span v-else>{{ stock.ticker }} - {{ stock.company }}</span>
          </li>
        </ul>
      </div>
    </div>
    <div class="flex justify-center w-full">
      <ButtonClick @click="handleButtonClick" />
    </div>
  </div>

  <div v-else class="flex items-center gap-2 relative " style="width:340px;">
    <img src="@/assets/search.svg" alt="Buscar" class="w-8 h-8 aspect-square cursor-pointer"
      @click="handleButtonClick" />
    <input v-model="searchValue" @keyup="handleKeyup" type="text" placeholder="Ingresa el Ticker de la acción"
      class="bg-transparent outline-none border-none text-[#91B6A2] text-[20px] font-normal leading-none px-2 flex-1"
      style="width:280px; transition: width 0.2s;" />
    <ul v-if="searchValue"
      class="absolute left-0 top-full z-10 w-full bg-white border border-[#91B6A2] rounded-lg shadow max-h-60 overflow-auto">
      <li v-for="stock in filteredStocks" :key="stock.ticker + stock.company"
        class="px-4 py-2 hover:bg-[#f0fdf4] cursor-pointer" @click="!stock.noResult && selectStock(stock)">
        <span v-if="stock.noResult" class="text-gray-400">{{ stock.company }}</span>
        <span v-else>{{ stock.ticker }} - {{ stock.company }}</span>
      </li>
    </ul>
  </div>

</template>
