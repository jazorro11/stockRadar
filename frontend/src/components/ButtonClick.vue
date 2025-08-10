<template>
  <div>
    <!-- Botón personalizado que emite un evento 'click' al hacer clic -->
    <button
      @click="$emit('button-click')"
      
      class="flex justify-center items-center px-[40px] py-[24px] rounded-[40px] bg-[#14AE5C] transition-colors duration-300 cursor-pointer hover:bg-[#189e53]"
    >
      <!-- Texto del botón con estilos personalizados y sin eventos de puntero -->
      <span class="pointer-events-none text-[#010502] font-roboto text-[24px] font-normal leading-[24px]">
        Buscar
      </span>
    </button>
  </div>
</template>

<script>
// Componente ButtonClick
// - Renderiza un botón estilizado
// - Emite un evento 'click' al hacer clic para que el componente padre lo maneje
export default {
  name: 'ButtonClick'
};
</script>