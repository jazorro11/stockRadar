<script setup>
// Importa funciones reactivas y utilidades de Vue.
import { ref, computed } from 'vue'
// Importa helpers para definir props y eventos personalizados.
import { defineProps, defineEmits } from 'vue'
// Importa el botón personalizado para búsqueda.
import ButtonClick from './ButtonClick.vue'
// Importa el store de acciones para acceder a los datos de las acciones.
import { useStocksStore } from '@/stores/stocks'

// Instancia el store y solicita la carga de acciones al montar el componente.
const store = useStocksStore()
store.fetchStocks()

// Define las props del componente:
// - is: determina el modo de visualización (completo o compacto).
// - stocks: lista de acciones (no se usa directamente, se usa store.stocks).
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

// Define los eventos personalizados que puede emitir el componente:
// - found: cuando se encuentra una acción.
// - notfound: cuando no se encuentra la acción buscada.
const emit = defineEmits(['found', 'notfound'])

// searchValue almacena el texto ingresado por el usuario en el input de búsqueda.
const searchValue = ref('')

// filteredStocks es una lista computada de acciones que coinciden con el texto ingresado.
// Si no hay coincidencias, muestra un mensaje de "No hay resultados".
const filteredStocks = computed(() => {
  const query = searchValue.value.trim().toUpperCase()
  if (!query) return []
  const matches = store.stocks.filter(stock => stock.ticker.toUpperCase().startsWith(query))
  if (matches.length === 0) {
    return [{ ticker: '', company: 'No hay resultados', noResult: true }]
  }
  return matches
})

// handleKeyup se ejecuta al presionar una tecla en el input.
// Si la tecla es Enter, realiza la búsqueda y limpia el input.
function handleKeyup(e) {
  if (e.key === 'Enter') {
    searchTicker()
    console.log('Enter pressed, searching for:', searchValue.value)
    searchValue.value = '' // Borra el input
  }
}

// handleButtonClick se ejecuta al hacer clic en el botón de búsqueda.
function handleButtonClick() {
  searchTicker()
  console.log('Click Handle Button')
}

// searchTicker busca el ticker exacto ingresado.
// Si lo encuentra, emite el evento 'found' con el stock encontrado.
// Si no lo encuentra, emite el evento 'notfound' con el ticker buscado.
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

// selectStock se ejecuta al seleccionar una acción de la lista de sugerencias.
// Actualiza el input, realiza la búsqueda y limpia el input.
function selectStock(stock) {
  searchValue.value = stock.ticker
  searchTicker()
  searchValue.value = '' // Borra el input
}
</script>

<template>
  <!-- Modo completo: barra de búsqueda centrada y botón grande -->
  <div v-if="is" class="w-full flex flex-col items-center gap-[72px]">
    <div class="flex items-center p-6 gap-2.5 rounded-[40px] border border-[#91B6A2] bg-[#FAFEFC]  w-full">
      <!-- Icono de búsqueda -->
      <img src="@/assets/search.svg" alt="Buscar" class="w-8 h-8 aspect-square" />

      <!-- Input de búsqueda y lista de sugerencias -->
      <div class="relative w-full">
        <input v-model="searchValue" @keyup="handleKeyup" type="text" placeholder="Ingresa el Ticker de la acción"
          class="bg-transparent outline-none border-none text-[#91B6A2] text-[20px] font-normal leading-none w-full" />
        <!-- Lista de sugerencias filtradas -->
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
    <!-- Botón de búsqueda -->
    <div class="flex justify-center w-full">
      <ButtonClick @click="handleButtonClick" />
    </div>
  </div>

  <!-- Modo compacto: barra de búsqueda horizontal y botón pequeño -->
  <div v-else class="flex items-center gap-2 relative " style="width:340px;">
    <!-- Icono de búsqueda clickeable -->
    <img src="@/assets/search.svg" alt="Buscar" class="w-8 h-8 aspect-square cursor-pointer"
      @click="handleButtonClick" />
    <!-- Input de búsqueda -->
    <input v-model="searchValue" @keyup="handleKeyup" type="text" placeholder="Ingresa el Ticker de la acción"
      class="bg-transparent outline-none border-none text-[#91B6A2] text-[20px] font-normal leading-none px-2 flex-1"
      style="width:280px; transition: width 0.2s;" />
    <!-- Lista de sugerencias filtradas -->
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