<script setup>
import { ref } from 'vue'
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

function handleKeyup(e) {
  if (e.key === 'Enter') {
    searchTicker()
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


</script>

<template>
  <div v-if="is" class="w-full flex flex-col items-center gap-[72px]">
    <div class="flex items-center p-6 gap-2.5 rounded-[40px] border border-[#91B6A2] bg-[#FAFEFC] w-full">
      <img 
      src="@/assets/search.svg" 
      alt="Buscar" 
      class="w-8 h-8 aspect-square" 
      />
      <input 
      v-model="searchValue" 
      @keyup="handleKeyup" 
      type="text" 
      placeholder="Ingresa el Ticker de la acción"
      class="bg-transparent outline-none border-none text-[#91B6A2] text-[20px] font-normal leading-none w-full" 
      />
    </div>
    <div class="flex justify-center w-full">
      <ButtonClick @click="handleButtonClick" />
    </div>
  </div>

<div v-else class="flex items-center gap-2" style="width:340px;">
  <img
    src="@/assets/search.svg"
    alt="Buscar"
    class="w-8 h-8 aspect-square cursor-pointer"
    @click="handleButtonClick"
  />
  <input
    v-model="searchValue"
    @keyup="handleKeyup"
    type="text"
    placeholder="Ingresa el Ticker de la acción"
    class="bg-transparent outline-none border-none text-[#91B6A2] text-[20px] font-normal leading-none px-2 flex-1"
    style="width:280px; transition: width 0.2s;"
  />
</div>

</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>