<script setup>
import HeaderHome from './HeaderHome.vue'
</script>

<template>
  <header class="flex items-center justify-between px-12 py-3 border-b-[1.5px] border-[#CCEBDB] bg-white w-full">
    <HeaderHome />
    <div class="flex items-center gap-12">
      <img src="@/assets/search.svg" alt="Buscar" class="w-6 h-6" />
    </div>
  </header>
</template>
