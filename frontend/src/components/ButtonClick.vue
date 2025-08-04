<template>
  <div>
    <button
      @click="$emit('click')"
      class="flex justify-center items-center px-[40px] py-[24px] rounded-[40px] bg-[#14AE5C] transition-colors duration-300 cursor-pointer hover:bg-[#189e53]"
    >
      <span class="pointer-events-none text-[#010502] font-roboto text-[24px] font-normal leading-[24px]">
        Buscar
      </span>
    </button>
  </div>
</template>

<script>
export default {
  name: 'ButtonClick'
};
</script>
