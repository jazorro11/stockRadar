<script setup>
// Importa el componente HeaderHome para mostrar el logo o título principal en el header.
import HeaderHome from './HeaderHome.vue'
// Importa el componente Search para la barra de búsqueda de acciones.
import Search from './Search.vue';
// Importa el store de acciones para acceder a la lista de stocks.
import { useStocksStore } from '@/stores/stocks'
// Importa el router de Vue para navegación programática.
import { useRouter } from 'vue-router'
// Importa ref para variables reactivas locales.
import { ref } from 'vue'

// Variable reactiva para almacenar el stock encontrado por la búsqueda.
const foundStock = ref(null)
// Variable reactiva para almacenar el ticker no encontrado.
const notFoundTicker = ref(null)
// Instancia el store de acciones.
const store = useStocksStore()
// Instancia el router para navegación.
const router = useRouter()
// Solicita la carga de los stocks al store al montar el componente.
store.fetchStocks()

// Maneja el caso cuando no se encuentra un ticker en la búsqueda.
// - Limpia foundStock y guarda el ticker no encontrado.
// - Limpia el input de búsqueda.
// - Muestra un mensaje en consola.
function handleNotFound(ticker) {
  foundStock.value = null
  notFoundTicker.value = ticker
  console.log('Stock HEADER no encontrado:', ticker)
  searchValue.value = '' // Borra el input
}

// Maneja el caso cuando se encuentra un stock en la búsqueda.
// - Guarda el stock encontrado y limpia notFoundTicker.
// - Muestra un mensaje en consola.
// - Navega a la vista de detalles del stock encontrado.
function handleFound(stock) {
  foundStock.value = stock
  notFoundTicker.value = null
  console.log('Stock HEADER encontrado:', stock)
  router.push({ name: 'StockDetails', params: { ticker: stock.ticker } })
}
</script>

<template>
  <!-- Header principal de la aplicación -->
  <header class="flex items-center justify-between px-12 py-3 border-b-[1.5px] border-[#CCEBDB] bg-white w-full">
    <!-- Logo o título principal -->
    <HeaderHome />
    <!-- Barra de búsqueda de acciones
         - Recibe la lista de stocks del store
         - Emite eventos 'found' y 'notfound' según el resultado de la búsqueda
    -->
    <Search :is="false" :stocks="store.stocks" @found="handleFound" @notfound="handleNotFound" />
  </header>
</template>