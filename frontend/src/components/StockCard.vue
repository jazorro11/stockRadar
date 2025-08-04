<template>
  <div class="flex flex-col gap-2 self-stretch p-8 rounded-2xl shadow w-full "
    style="background:rgba(255,255,255,0.5);">
    <div class="flex justify-between items-center">
      <span
        class="align-self-stretch text-[#1B2821] font-inter text-[31px] font-semibold leading-[24px]
        [display:-webkit-box] [overflow:hidden] [text-overflow:ellipsis] [white-space:normal]
        [-webkit-line-clamp:1] [-webkit-box-orient:vertical]"
      >
        {{ stock.ticker }}
      </span>
      <span
        class="text-[#0E7B41] font-inter text-[31px] font-semibold leading-[24px]
        [letter-spacing:-0.085px] [display:-webkit-box] [overflow:hidden] [text-overflow:ellipsis] [white-space:normal]
        [-webkit-line-clamp:1] [-webkit-box-orient:vertical]"
      >
        ${{ Number(stock.current_price).toFixed(2) }}
      </span>
    </div>
    <div
      class="align-self-stretch text-[#4B4A5C] font-inter text-[15px] font-medium leading-[24px]
      [display:-webkit-box] [overflow:hidden] [text-overflow:ellipsis] [white-space:normal]
      [-webkit-line-clamp:1] [-webkit-box-orient:vertical]"
    >
      {{ stock.company }}
    </div>
        <div
      class="align-self-stretch text-[#1B2821] font-inter text-[15px] font-normal leading-[20px]
        [letter-spacing:-0.075px] [display:-webkit-box] [overflow:hidden] [text-overflow:ellipsis] [white-space:normal]
        [-webkit-line-clamp:1] [-webkit-box-orient:vertical]"
    >
      {{ stock.brokerage }}
    </div>
    <div
      class="flex justify-between items-center text-[#1B2821] font-inter text-[15px] font-normal leading-[20px]"
    >
      <button
        class="flex justify-center items-center px-6 py-2 rounded-full bg-[#14AE5C] mt-2"
        type="button"
      >
        <span class="text-[#010502] font-inter text-[15px] font-medium leading-[20px]">
        Ver detalles
      </span>
    </button>
    </div>
  </div>
</template>

<script setup>
defineProps({
  stock: {
    type: Object,
    required: true
  }
})
</script>