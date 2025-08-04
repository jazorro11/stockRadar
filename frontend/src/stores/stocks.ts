import { defineStore } from 'pinia'
import axios from 'axios'

export const useStocksStore = defineStore('stocks', {
  state: () => ({
    stocks: [],
    loading: false,
    error: null as null | unknown,
  }),
  actions: {
    async fetchStocks() {
      this.loading = true
      this.error = null
      try {
        const res = await axios.get('http://localhost:8080/stocks')
        this.stocks = res.data
      } catch (e) {
        this.error = e
      } finally {
        this.loading = false
      }
    },
  },
})