// Importa la función defineStore de Pinia para crear un store reactivo.
import { defineStore } from 'pinia'
// Importa axios para realizar peticiones HTTP.
import axios from 'axios'

// Define y exporta el store 'stocks' usando Pinia.
export const useStocksStore = defineStore('stocks', {
  // Estado reactivo del store.
  state: () => ({
    // stocks: almacena la lista de acciones obtenidas del backend.
    stocks: [],
    // loading: indica si se está realizando una petición de datos.
    loading: false,
    // error: almacena cualquier error ocurrido durante la petición.
    error: null as null | unknown,
  }),
  // Acciones del store (funciones que pueden modificar el estado).
  actions: {
    // Acción asíncrona para obtener la lista de acciones desde el backend.
    async fetchStocks() {
      // Marca el estado como cargando.
      this.loading = true
      // Reinicia el error previo.
      this.error = null
      try {
        // Realiza la petición GET al endpoint de acciones.
        const res = await axios.get('http://localhost:8080/stocks')
        // Almacena los datos recibidos en el estado.
        this.stocks = res.data
      } catch (e) {
        // Si ocurre un error, lo guarda en el estado.
        this.error = e
      } finally {
        // Marca que la carga ha terminado, haya éxito o error.
        this.loading = false
      }
    },
  },
})