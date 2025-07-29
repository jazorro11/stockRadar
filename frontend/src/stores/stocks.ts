/**
 * Archivo: stocks.ts
 * ------------------
 * Este archivo define la interfaz TypeScript `StockInfo`, que representa la estructura de los datos
 * de acciones bursátiles (stocks) que se manejan en la aplicación StockRadar UI.
 *
 * Esta interfaz permite definir claramente el formato que debe tener cada objeto de acción
 * obtenido desde el backend (escrito en Go y conectado a CockroachDB), o directamente desde una API.
 *
 * Gracias a TypeScript, este tipo de estructura ayuda a:
 * - Validar y restringir el tipo de datos (evita errores)
 * - Mejorar la autocompletación y documentación en el editor (VS Code, WebStorm, etc.)
 * - Mantener coherencia en los componentes de Vue, el store de Pinia y los servicios de API
 *
 * Esta interfaz será utilizada dentro del Pinia Store para almacenar y manipular la lista de acciones disponibles,
 * y también puede ser reutilizada en componentes Vue como `StockCard.vue`, `Home.vue`, etc.
 */

import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

// Definimos la estructura que tendrá cada acción bursátil que llega desde la API o la base de datos
export interface StockInfo {
  ticker: string
  company: string
  brokerage: string
  action: string
  rating_from: string
  rating_to: string
  target_from: number
  target_to: number
  time: string
  market_cap: number
  eps_ttm: number
  pe_ttm: number
  pb: number
  dividend_yield: number
  week_52_high: number
  week_52_low: number
  revenue_growth_ttm_yoy: number
  net_profit_margin_ttm: number
  beta: number
  current_price: number
  score: number
}

// Definimos el store principal de acciones bursátiles usando Pinia
export const useStocksStore = defineStore('stocks', () => {
  // --- STATE: almacena el estado reactivo del store ---
  const stocks = ref<StockInfo[]>([])        // Lista de acciones recuperadas desde el backend
  const loading = ref(false)                   // Indica si está en proceso de carga
  const error = ref<string | null>(null)     // Guarda mensajes de error en caso de fallos
  const search = ref('')                       // Texto de búsqueda para filtrar acciones
  const sortKey = ref<keyof StockInfo>('time')                  // Campo por el cual ordenar ('time' por defecto)
  const sortAsc = ref(false)                   // Indica si el orden es ascendente (true) o descendente (false)

  // Filtros adicionales
  const filterAction = ref('')
  const filterDividend = ref(0)
  const filterBeta = ref(0)
  const filterScore = ref(0)

  // --- ACTIONS: métodos para mutar el estado y ejecutar lógica asíncrona ---
  async function fetchStocks(): Promise<void> {
    loading.value = true         // Indicamos que comenzó la carga
    error.value = null           // Limpiamos errores anteriores

    try {
      // Llamada al backend local en el endpoint /stocks (ajustar si es necesario)
      const res = await fetch('http://localhost:8080/stocks')

      // Validamos que la respuesta sea OK (status 200–299)
      if (!res.ok) throw new Error('Failed to fetch stocks')

      // Si la respuesta es válida, parseamos el JSON y lo almacenamos en el estado
      stocks.value = await res.json()
    } catch (e: any) {
      // Si ocurre un error (network, parsing, etc), lo guardamos en error
      error.value = e.message
    } finally {
      // Finalizamos el estado de carga, sin importar éxito o fallo
      loading.value = false
    }
  }

  function setSearch(query: string): void {
    search.value = query
  }

  function setSort(key: keyof StockInfo): void {
    if (sortKey.value === key) {
      // Si el usuario vuelve a hacer clic sobre el mismo campo, invierte el orden
      sortAsc.value = !sortAsc.value
    } else {
      // Si cambia de campo, lo asignamos y reiniciamos a descendente
      sortKey.value = key
      sortAsc.value = false
    }
  }

  // --- GETTERS: funciones computadas basadas en el estado ---
  const filteredStocks = computed(() =>
    stocks.value
      .filter(stock =>
        (!filterAction.value || stock.action === filterAction.value) &&
        (filterDividend.value === 0 || stock.dividend_yield >= filterDividend.value) &&
        (filterBeta.value === 0 || stock.beta <= filterBeta.value) &&
        (filterScore.value === 0 || stock.score >= filterScore.value) &&
        (
          stock.ticker.toLowerCase().includes(search.value.toLowerCase()) ||
          stock.company.toLowerCase().includes(search.value.toLowerCase()) ||
          stock.brokerage.toLowerCase().includes(search.value.toLowerCase())
        )
      )
      .sort((a, b) => {
        const key = sortKey.value
        if (a[key] < b[key]) return sortAsc.value ? -1 : 1
        if (a[key] > b[key]) return sortAsc.value ? 1 : -1
        return 0
      })
  )

  // --- RETURN ---
  return {
    stocks,
    loading,
    error,
    search,
    sortKey,
    sortAsc,
    filterAction,
    filterDividend,
    filterBeta,
    filterScore,
    fetchStocks,
    setSearch,
    setSort,
    filteredStocks,
  }
})