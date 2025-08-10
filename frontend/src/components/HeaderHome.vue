<script setup>
// Importa el hook useRouter de Vue Router para navegación programática.
import { useRouter } from 'vue-router'
// Obtiene la instancia del router para poder navegar entre rutas.
const router = useRouter()
// Función que navega a la página principal ("/") cuando se llama.
function goHome() {
  router.push({ path: '/' })
}
</script>

<template>
  <!-- Contenedor principal del header, muestra el logo y el nombre de la app.
       - Usa flex para alinear horizontalmente el logo y el texto.
       - gap-1 para separar el logo y el texto.
       - h-8 para altura fija.
       - cursor-pointer para indicar que es clickeable.
       - Al hacer click, ejecuta goHome para navegar al home.
  -->
  <div
    class="flex items-center h-8 gap-1 cursor-pointer"
    @click="goHome"
  >
    <!-- Logo de la aplicación, imagen SVG de un cerdito -->
    <img
      src="@/assets/logomark.svg"
      alt="Cerdito"
      class="w-8 h-8"
    />
    <!-- Nombre de la aplicación con estilos de color, fuente y tamaño -->
    <span class="text-[#0F8446] font-[400] text-[24px] leading-none tracking-[-0.48px] font-roboto">
      Stocky
    </span>
  </div>
</template>