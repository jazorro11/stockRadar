<script setup>



const props = defineProps({
    stock: {
        type: Object,
        required: true
    },
    financialDetail: {
        type: Boolean,
        required: true
    }
})

</script>

<template>
    <div v-if="financialDetail" class="w-full flex flex-col gap-8 self-stretch p-8 rounded-2xl shadow"
        style="background:rgba(255,255,255,0.5);">
        <div class="text-[#1B2821] font-roboto text-[24px] font-bold leading-[52px] tracking-[-0.48px]">
            Indicadores financieros
        </div>
        <p class="flex-1 overflow-hidden text-[#1B2821] font-inter text-[15px] font-medium leading-[20px] tracking-[-0.075px] text-ellipsis whitespace-normal"
            style="display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 1;">
            Radio Precio/Benecicio (TTM): {{Number(stock?.pe_ttm).toFixed(2) }}
        </p>
        <p class="flex-1 overflow-hidden text-[#1B2821] font-inter text-[15px] font-medium leading-[20px] tracking-[-0.075px] text-ellipsis whitespace-normal"
            style="display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 1;">
            Precio /Valor contable: {{ Number(stock?.pb).toFixed(2) }}
        </p>
        <p class="flex-1 overflow-hidden text-[#1B2821] font-inter text-[15px] font-medium leading-[20px] tracking-[-0.075px] text-ellipsis whitespace-normal"
            style="display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 1;">
            Ganancias por acción (TTM): {{ Number(stock?.eps_ttm).toFixed(2) }}
        </p>
        <p class="flex-1 overflow-hidden text-[#1B2821] font-inter text-[15px] font-medium leading-[20px] tracking-[-0.075px] text-ellipsis whitespace-normal"
            style="display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 1;">
            Margen neto (TTM): {{ Number(stock?.net_profit_margin_ttm).toFixed(2) }} %
        </p>
    </div>

    <div v-else class="w-full flex flex-col gap-8 self-stretch p-8 rounded-2xl shadow"
        style="background:rgba(255,255,255,0.5);">
        <div class="text-[#1B2821] font-roboto text-[24px] font-bold leading-[52px] tracking-[-0.48px]">
            Indicadores del Mercado
        </div>
        <p class="flex-1 overflow-hidden text-[#1B2821] font-inter text-[15px] font-medium leading-[20px] tracking-[-0.075px] text-ellipsis whitespace-normal"
            style="display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 1;">
            Máximo en 52 semanas: {{Number(stock?.week_52_high).toFixed(2) }}
        </p>
        <p class="flex-1 overflow-hidden text-[#1B2821] font-inter text-[15px] font-medium leading-[20px] tracking-[-0.075px] text-ellipsis whitespace-normal"
            style="display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 1;">
            Mínimo en 52 semanas: {{ Number(stock?.week_52_low).toFixed(2) }}
        </p>
        <p class="flex-1 overflow-hidden text-[#1B2821] font-inter text-[15px] font-medium leading-[20px] tracking-[-0.075px] text-ellipsis whitespace-normal"
            style="display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 1;">
            Volatilidad respecto al mercado: {{ Number(stock?.beta).toFixed(2) }}
        </p>
        <p class="flex-1 overflow-hidden text-[#1B2821] font-inter text-[15px] font-medium leading-[20px] tracking-[-0.075px] text-ellipsis whitespace-normal"
            style="display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 1;">
            Crecimiento de ingresos interanual: {{ Number(stock?.revenue_growth_ttm_yoy).toFixed(2) }} %
        </p>
    </div>
</template>