<script setup>
// Importa las props del componente usando defineProps.
// - stock: objeto con los datos de la acción a mostrar (obligatorio).
// - financialDetail: booleano que indica si se deben mostrar indicadores financieros (true) o de mercado (false).
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
    <!-- Si financialDetail es true, muestra los indicadores financieros -->
    <div v-if="financialDetail" class="w-full flex flex-col gap-8 self-stretch p-8 rounded-2xl shadow"
        style="background:rgba(255,255,255,0.5);">
        <!-- Título de la sección -->
        <div class="text-[#1B2821] font-roboto text-[24px] font-bold leading-[52px] tracking-[-0.48px]">
            Indicadores financieros
        </div>
        <!-- Indicador: Radio Precio/Beneficio (TTM) -->
        <p class="flex-1 overflow-hidden text-[#1B2821] font-inter text-[15px] font-medium leading-[20px] tracking-[-0.075px] text-ellipsis whitespace-normal"
            style="display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 1; line-clamp: 1;">
            Radio Precio/Benecicio (TTM): {{Number(stock?.pe_ttm).toFixed(2) }}
        </p>
        <!-- Indicador: Precio / Valor contable -->
        <p class="flex-1 overflow-hidden text-[#1B2821] font-inter text-[15px] font-medium leading-[20px] tracking-[-0.075px] text-ellipsis whitespace-normal"
            style="display: -webkit-box; -webkit-box-orient: vertical;  -webkit-line-clamp: 1; line-clamp: 1;">
            Precio /Valor contable: {{ Number(stock?.pb).toFixed(2) }}
        </p>
        <!-- Indicador: Ganancias por acción (TTM) -->
        <p class="flex-1 overflow-hidden text-[#1B2821] font-inter text-[15px] font-medium leading-[20px] tracking-[-0.075px] text-ellipsis whitespace-normal"
            style="display: -webkit-box; -webkit-box-orient: vertical;  -webkit-line-clamp: 1; line-clamp: 1;">
            Ganancias por acción (TTM): {{ Number(stock?.eps_ttm).toFixed(2) }}
        </p>
        <!-- Indicador: Margen neto (TTM) -->
        <p class="flex-1 overflow-hidden text-[#1B2821] font-inter text-[15px] font-medium leading-[20px] tracking-[-0.075px] text-ellipsis whitespace-normal"
            style="display: -webkit-box; -webkit-box-orient: vertical;  -webkit-line-clamp: 1; line-clamp: 1;">
            Margen neto (TTM): {{ Number(stock?.net_profit_margin_ttm).toFixed(2) }} %
        </p>
    </div>

    <!-- Si financialDetail es false, muestra los indicadores del mercado -->
    <div v-else class="w-full flex flex-col gap-8 self-stretch p-8 rounded-2xl shadow"
        style="background:rgba(255,255,255,0.5);">
        <!-- Título de la sección -->
        <div class="text-[#1B2821] font-roboto text-[24px] font-bold leading-[52px] tracking-[-0.48px]">
            Indicadores del Mercado
        </div>
        <!-- Indicador: Máximo en 52 semanas -->
        <p class="flex-1 overflow-hidden text-[#1B2821] font-inter text-[15px] font-medium leading-[20px] tracking-[-0.075px] text-ellipsis whitespace-normal"
            style="display: -webkit-box; -webkit-box-orient: vertical;  -webkit-line-clamp: 1; line-clamp: 1;;">
            Máximo en 52 semanas: {{Number(stock?.week_52_high).toFixed(2) }}
        </p>
        <!-- Indicador: Mínimo en 52 semanas -->
        <p class="flex-1 overflow-hidden text-[#1B2821] font-inter text-[15px] font-medium leading-[20px] tracking-[-0.075px] text-ellipsis whitespace-normal"
            style="display: -webkit-box; -webkit-box-orient: vertical;  -webkit-line-clamp: 1; line-clamp: 1;">
            Mínimo en 52 semanas: {{ Number(stock?.week_52_low).toFixed(2) }}
        </p>
        <!-- Indicador: Volatilidad respecto al mercado (beta) -->
        <p class="flex-1 overflow-hidden text-[#1B2821] font-inter text-[15px] font-medium leading-[20px] tracking-[-0.075px] text-ellipsis whitespace-normal"
            style="display: -webkit-box; -webkit-box-orient: vertical;  -webkit-line-clamp: 1; line-clamp: 1;">
            Volatilidad respecto al mercado: {{ Number(stock?.beta).toFixed(2) }}
        </p>
        <!-- Indicador: Crecimiento de ingresos interanual -->
        <p class="flex-1 overflow-hidden text-[#1B2821] font-inter text-[15px] font-medium leading-[20px] tracking-[-0.075px] text-ellipsis whitespace-normal"
            style="display: -webkit-box; -webkit-box-orient: vertical;  -webkit-line-clamp: 1; line-clamp: 1;">
            Crecimiento de ingresos interanual: {{ Number(stock?.revenue_growth_ttm_yoy).toFixed(2) }} %
        </p>
    </div>
</template>