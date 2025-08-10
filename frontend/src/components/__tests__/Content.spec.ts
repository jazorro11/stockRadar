
// Este archivo contiene pruebas unitarias para el componente Content.vue utilizando Vitest y Vue Test Utils.
// Se mockean (simulan) dependencias externas: router (vue-router), store (Pinia) y componentes hijos (Search, StockCard)
// para aislar el comportamiento interno de Content.vue y validar:
//  1. Que al montar el componente se llama fetchStocks() y sólo renderiza las 3 acciones con mayor 'normalized'.
//  2. Que al emitirse el evento 'found' desde Search se navega correctamente usando router.push().
//  3. Que al emitirse el evento 'notfound' no se realiza navegación.

import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { reactive, nextTick } from 'vue'

// Mock router
const pushMock = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({ push: pushMock })
}))

// Mock datos (4 acciones para validar tope de 3)
const mockStocks = reactive([
  { ticker: 'AAA', normalized: 0.95 },
  { ticker: 'BBB', normalized: 0.80 },
  { ticker: 'CCC', normalized: 0.60 },
  { ticker: 'DDD', normalized: 0.99 }
])

// Mock store
const fetchStocksMock = vi.fn()
vi.mock('@/stores/stocks', () => ({
  useStocksStore: () => ({
    stocks: mockStocks,
    fetchStocks: fetchStocksMock
  })
}))

// Stubs de componentes hijos (mockeamos los módulos que Content importa)
vi.mock('../Search.vue', () => ({
  default: {
    name: 'Search',
    props: ['is', 'stocks'],
    template: `<div class="search-stub">SEARCH</div>`
  }
}))
vi.mock('../StockCard.vue', () => ({
  default: {
    name: 'StockCard',
    props: ['stock'],
    template: `<div class="stock-card">{{ stock.ticker }}</div>`
  }
}))

// Importar después de mocks
import Content from '../Content.vue'

describe('Content.vue', () => {
  beforeEach(() => {
    pushMock.mockClear()
    fetchStocksMock.mockClear()
  })

  it('llama fetchStocks y muestra solo las 3 acciones top ordenadas', async () => {
    const wrapper = mount(Content)
    // Esperar ciclo de actualización
    await nextTick()

    expect(fetchStocksMock).toHaveBeenCalled()

    const cards = wrapper.findAll('.stock-card')
    expect(cards.length).toBe(3)

    const tickers = cards.map(c => c.text())
    // Orden esperado: DDD (0.99), AAA (0.95), BBB (0.80)
    expect(tickers).toEqual(['DDD', 'AAA', 'BBB'])
  })

  it('navega al emitir found', async () => {
    const wrapper = mount(Content)
    await nextTick()

    // Buscar stub de Search
    const search = wrapper.find('.search-stub')
    expect(search.exists()).toBe(true)

    // Emitir manualmente desde el stub:
    // Como es un stub simple, necesitamos acceder al componente mockeado
    wrapper.findComponent({ name: 'Search' }).vm.$emit('found', { ticker: 'AAA' })

    expect(pushMock).toHaveBeenCalledWith({
      name: 'StockDetails',
      params: { ticker: 'AAA' }
    })
  })

  it('no navega al emitir notfound', async () => {
    const wrapper = mount(Content)
    await nextTick()

    wrapper.findComponent({ name: 'Search' }).vm.$emit('notfound', 'ZZZZ')
    expect(pushMock).not.toHaveBeenCalled()
  })
})