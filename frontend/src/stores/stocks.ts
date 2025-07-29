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

import { defineStore } from 'pinia' // Importamos la función para definir un store de estado global con Pinia

// Definimos la estructura que tendrá cada acción bursátil que llega desde la API o la base de datos
export interface StockInfo {
  ticker: string        // Símbolo bursátil de la empresa (ej. AAPL, TSLA)
  company: string       // Nombre completo de la empresa
  brokerage: string     // Nombre de la firma de corretaje que hizo la recomendación
  action: string        // Tipo de acción sugerida (Buy, Hold, Sell, etc.)
  rating_from: string   // Calificación anterior del activo
  rating_to: string     // Nueva calificación otorgada al activo
  target_from: string   // Precio objetivo anterior (puede venir como texto)
  target_to: string     // Precio objetivo actualizado (puede venir como texto)
  time: string          // Fecha/hora de la recomendación (en formato ISO 8601 generalmente)
}
// Definimos el store principal de acciones bursátiles usando Pinia
export const useStocksStore = defineStore('stocks', {
  // --- STATE: almacena el estado reactivo del store ---
  state: () => ({
    stocks: [] as StockInfo[],        // Lista de acciones recuperadas desde el backend
    loading: false,                   // Indica si está en proceso de carga
    error: null as string | null,     // Guarda mensajes de error en caso de fallos
    search: '',                       // Texto de búsqueda para filtrar acciones
    sortKey: 'time',                  // Campo por el cual ordenar ('time' por defecto)
    sortAsc: false,                   // Indica si el orden es ascendente (true) o descendente (false)
  }),

  // --- ACTIONS: métodos para mutar el estado y ejecutar lógica asíncrona ---
  actions: {
    // fetchStocks(): obtiene la lista de acciones desde el backend
    async fetchStocks(): Promise<void> {
      this.loading = true         // Indicamos que comenzó la carga
      this.error = null           // Limpiamos errores anteriores

      try {
        // Llamada al backend local en el endpoint /stocks (ajustar si es necesario)
        const res = await fetch('http://localhost:8080/stocks')

        // Validamos que la respuesta sea OK (status 200–299)
        if (!res.ok) throw new Error('Failed to fetch stocks')

        // Si la respuesta es válida, parseamos el JSON y lo almacenamos en el estado
        this.stocks = await res.json()
      } catch (e: any) {
        // Si ocurre un error (network, parsing, etc), lo guardamos en error
        this.error = e.message
      } finally {
        // Finalizamos el estado de carga, sin importar éxito o fallo
        this.loading = false
      }
    },

    // setSearch(query): actualiza el término de búsqueda en el estado
    setSearch(query: string): void {
      this.search = query
    },

    // setSort(key): actualiza el campo por el que se ordena la lista
    setSort(key: string): void {
      if (this.sortKey === key) {
        // Si el usuario vuelve a hacer clic sobre el mismo campo, invierte el orden
        this.sortAsc = !this.sortAsc
      } else {
        // Si cambia de campo, lo asignamos y reiniciamos a descendente
        this.sortKey = key
        this.sortAsc = false
      }
    }
  },

  // --- GETTERS: funciones computadas basadas en el estado ---
  getters: {
    // filteredStocks(): devuelve la lista filtrada por búsqueda y ordenada
    filteredStocks(state): StockInfo[] {
      let result = state.stocks

      // Si hay un término de búsqueda, filtramos por ticker, empresa o bróker
      if (state.search) {
        const q = state.search.toLowerCase()
        result = result.filter(s =>
          s.ticker.toLowerCase().includes(q) ||
          s.company.toLowerCase().includes(q) ||
          s.brokerage.toLowerCase().includes(q)
        )
      }

      // Ordenamos la lista por el campo `sortKey` y la dirección `sortAsc`
      result = [...result].sort((a, b) => {
        const key = state.sortKey as keyof StockInfo

        if (a[key] < b[key]) return state.sortAsc ? -1 : 1
        if (a[key] > b[key]) return state.sortAsc ? 1 : -1
        return 0
      })

      return result
    }
  }
})