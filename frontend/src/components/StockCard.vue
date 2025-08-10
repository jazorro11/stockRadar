<script setup>
// Importa el hook useRouter de Vue Router para navegación programática.
import { useRouter } from 'vue-router'
// Define las props del componente, esperando un objeto stock obligatorio.
const props = defineProps({
  stock: {
    type: Object,
    required: true
  }
})
// Obtiene la instancia del router para poder navegar entre rutas.
const router = useRouter()
// Función que navega a la vista de detalles del stock seleccionado.
function goToDetails() {
  router.push({ name: 'StockDetails', params: { ticker: props.stock.ticker } })
}
</script>

<template>
  <!-- Tarjeta principal del stock, con fondo blanco semitransparente y sombra -->
  <div class="flex flex-col gap-2 self-stretch p-8 rounded-2xl shadow w-full "
    style="background:rgba(255,255,255,0.5);">
    <!-- Encabezado: ticker y precio actual alineados horizontalmente -->
    <div class="flex justify-between items-center">
      <!-- Ticker de la acción, texto grande y en negrita, con elipsis si es muy largo -->
      <span
        class="align-self-stretch text-[#1B2821] font-inter text-[31px] 
        font-semibold leading-[36px] whitespace-nowrap overflow-hidden text-ellipsis"
      >
        {{ stock.ticker }}
      </span>
      <!-- Precio actual, texto grande en verde, con elipsis si es muy largo -->
      <div
        class="text-[#0E7B41] font-inter text-[31px] font-semibold leading-[36px]
        [letter-spacing:-0.085px] whitespace-nowrap overflow-hidden text-ellipsis"
      >
        ${{ Number(stock.current_price).toFixed(2) }}
    </div>
    </div>
    <!-- Nombre de la compañía, texto mediano, con elipsis si es muy largo -->
    <div
      class="align-self-stretch text-[#4B4A5C] font-inter text-[15px] font-medium leading-[24px]
      [display:-webkit-box] [overflow:hidden] [text-overflow:ellipsis] [white-space:normal]
      [-webkit-line-clamp:1] [-webkit-box-orient:vertical]"
    >
      {{ stock.company }}
    </div>
    <!-- Nombre del bróker, texto pequeño, con elipsis si es muy largo -->
    <div
      class="align-self-stretch text-[#1B2821] font-inter text-[15px] font-normal leading-[20px]
        [letter-spacing:-0.075px] [display:-webkit-box] [overflow:hidden] [text-overflow:ellipsis] [white-space:normal]
        [-webkit-line-clamp:1] [-webkit-box-orient:vertical]"
    >
      {{ stock.brokerage }}
    </div>

    <!-- Botón para ver detalles, alineado a la derecha -->
    <div class="flex justify-end mt-4">
      <button
        class="flex justify-center items-center px-6 py-2 rounded-full bg-[#14AE5C] text-[#010502] font-inter text-[15px] font-medium leading-[20px] hover:bg-[#189e53] transition-colors cursor-pointer"
        type="button"
        @click="goToDetails"
      >
        Ver detalles
      </button>
    </div>

  </div>
</template>