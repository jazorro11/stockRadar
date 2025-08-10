/**
 * Search.spec.ts
 * Pruebas unitarias para el componente Search.vue usando Vitest + Vue Test Utils.
 *
 * Objetivos de las pruebas:
 *  1. Verificar que al montar el componente se invoca fetchStocks() del store (inicialización de datos).
 *  2. Confirmar que no se muestra la lista desplegable cuando el campo de búsqueda está vacío.
 *  3. Validar el filtrado case-insensitive por prefijo de ticker (ej: 'a' filtra AAPL y AMZN).
 *  4. Mostrar el mensaje "No hay resultados" cuando no existen coincidencias.
 *  5. Emitir el evento 'found' al presionar Enter con un ticker exacto (y limpiar el input).
 *  6. Emitir el evento 'notfound' al presionar Enter con un ticker inexistente.
 *  7. Emitir 'found' al hacer clic sobre una sugerencia válida de la lista.
 *  8. Emitir 'found' al pulsar el botón de búsqueda con un ticker válido.
 *  9. Emitir 'notfound' al pulsar el botón con un ticker inexistente.
 *
 * Estrategia:
 *  - Se mockea el store (useStocksStore) para: evitar dependencias reales y controlar los datos (stocksMock).
 *  - Se stubbea ButtonClick.vue para simplificar su comportamiento y emitir un evento 'click'.
 *  - Se usa mount() para renderizar Search.vue y manipular el input simulando interacciones del usuario.
 *  - Se emplea nextTick() tras cambiar valores reactivos para asegurar que el DOM se haya actualizado.
 */

import { describe, it, expect, vi, beforeEach } from 'vitest'          // Herramientas de Vitest: definir suites/casos, expectativas y mocks
import { mount } from '@vue/test-utils'                                // Utilidad para montar componentes Vue en entorno de prueba
import { nextTick } from 'vue'                                         // nextTick asegura esperar el ciclo de render tras cambios reactivos

// Mock del store Pinia usado dentro de Search.vue
const fetchStocksMock = vi.fn()                                        // Espía para verificar que fetchStocks() se llame al montar
const stocksMock = [                                                   // Datos simulados de acciones (varias para probar filtrado)
  { ticker: 'AAPL', company: 'Apple Inc.' },
  { ticker: 'AMZN', company: 'Amazon.com Inc.' },
  { ticker: 'MSFT', company: 'Microsoft Corp.' }
]
vi.mock('@/stores/stocks', () => ({                                    // Mock del módulo del store
  useStocksStore: () => ({                                             // Sustituye la implementación de useStocksStore
    stocks: stocksMock,                                                // Retorna el arreglo mock
    fetchStocks: fetchStocksMock                                       // Método mockeado para espiar su invocación
  })
}))

// Stub de ButtonClick (el real emite 'button-click', aquí emitimos 'click' para que el handler funcione)
// Simplificamos el componente hijo para aislar el test a la lógica de Search.vue
vi.mock('../ButtonClick.vue', () => ({
  default: {
    name: 'ButtonClick',                                               // Nombre (importante para resolución de componentes)
    template: `<button class="btn-stub" @click="$emit('click')">Buscar</button>` // Botón mínimo que emite 'click'
  }
}))

// Importar después de mocks (asegura que las sustituciones anteriores apliquen)
import Search from '../Search.vue'

// Helper para montar el componente con props base y permitir extenderlas por test
function mountSearch(extraProps: Record<string, any> = {}) {
  return mount(Search, {
    props: { is: true, stocks: stocksMock, ...extraProps }             // 'is' y 'stocks' replican los props esperados
  })
}

describe('Search.vue', () => {
  beforeEach(() => {                                                   // Se ejecuta antes de cada test
    fetchStocksMock.mockClear()                                        // Limpia contador y llamadas previas del mock
  })

  it('llama fetchStocks al montar', () => {                            // Test 1: inicialización correcta
    mountSearch()                                                      // Monta el componente
    expect(fetchStocksMock).toHaveBeenCalled()                         // Verifica que se invocó fetchStocks()
  })

  it('no muestra lista si el input está vacío', async () => {          // Test 2: ausencia de dropdown sin criterio
    const wrapper = mountSearch()                                      // Monta componente
    expect(wrapper.find('ul').exists()).toBe(false)                    // No debe existir la lista (ul)
  })

  it('filtra acciones por prefijo del ticker (case-insensitive)', async () => { // Test 3: filtrado por prefijo
    const wrapper = mountSearch()                                      // Monta componente
    const input = wrapper.find('input')                                // Localiza input de búsqueda
    await input.setValue('a')                                          // Escribe 'a' en minúscula
    await nextTick()                                                   // Espera actualización del DOM
    const items = wrapper.findAll('ul li')                             // Obtiene las filas de la lista
    const texts = items.map(li => li.text())                           // Extrae el texto de cada item
    expect(texts).toContain('AAPL - Apple Inc.')                       // Debe contener AAPL
    expect(texts).toContain('AMZN - Amazon.com Inc.')                  // Debe contener AMZN
    expect(texts).not.toContain('MSFT - Microsoft Corp.')              // No debe incluir MSFT (no coincide el prefijo)
  })

  it('muestra "No hay resultados" cuando no hay coincidencias', async () => { // Test 4: mensaje vacío
    const wrapper = mountSearch()
    await wrapper.find('input').setValue('ZZ')                         // Valor sin coincidencias
    await nextTick()
    const items = wrapper.findAll('ul li')                             // Lista renderizada
    expect(items.length).toBe(1)                                       // Solo un elemento (el mensaje)
    expect(items[0].text()).toBe('No hay resultados')                  // Texto esperado
  })

  it('emite "found" al presionar Enter con un ticker exacto', async () => { // Test 5: Enter match exacto
    const wrapper = mountSearch()
    const input = wrapper.find('input')
    await input.setValue('MSFT')                                       // Escribe ticker exacto
    await input.trigger('keyup', { key: 'Enter' })                     // Simula Enter
    const emitted = wrapper.emitted('found')                           // Revisa eventos emitidos
    expect(emitted).toBeTruthy()                                       // Debe haberse emitido 'found'
    expect(emitted![0][0]).toEqual(stocksMock[2])                      // Payload: objeto MSFT
    expect((input.element as HTMLInputElement).value).toBe('')         // Input se limpia tras búsqueda exitosa
  })

  it('emite "notfound" al presionar Enter con un ticker inexistente', async () => { // Test 6: Enter sin match
    const wrapper = mountSearch()
    const input = wrapper.find('input')
    await input.setValue('XXXX')                                       // Ticker inexistente
    await input.trigger('keyup', { key: 'Enter' })                     // Simula Enter
    const emitted = wrapper.emitted('notfound')                        // Verifica evento alternativo
    expect(emitted).toBeTruthy()                                       // Debe existir
    expect(emitted![0][0]).toBe('XXXX')                                // Payload: cadena buscada
  })

  it('al hacer clic en una sugerencia válida emite "found"', async () => { // Test 7: click en sugerencia
    const wrapper = mountSearch()
    await wrapper.find('input').setValue('A')                          // Prefijo que genera varias coincidencias
    await nextTick()
    const first = wrapper.find('ul li')                                // Primer resultado
    await first.trigger('click')                                       // Simula clic
    const emitted = wrapper.emitted('found')                           // Eventos 'found'
    expect(emitted).toBeTruthy()
    expect(emitted![0][0]).toEqual(stocksMock[0])                      // Primer match: AAPL
  })

  it('al clicar el botón ejecuta búsqueda y emite "found"', async () => { // Test 8: botón con match
    const wrapper = mountSearch()
    const input = wrapper.find('input')
    await input.setValue('AAPL')                                       // Ticker válido
    await wrapper.find('.btn-stub').trigger('click')                   // Clic en botón stub
    const emitted = wrapper.emitted('found')
    expect(emitted).toBeTruthy()
    expect(emitted![0][0]).toEqual(stocksMock[0])                      // Payload correcto
  })

  it('al clicar el botón con ticker inexistente emite "notfound"', async () => { // Test 9: botón sin match
    const wrapper = mountSearch()
    await wrapper.find('input').setValue('QQQQ')                       // Ticker inválido
    await wrapper.find('.btn-stub').trigger('click')                   // Clic en botón
    const emitted = wrapper.emitted('notfound')                        // Evento esperado
    expect(emitted).toBeTruthy()
    expect(emitted![0][0]).toBe('QQQQ')                                // Verifica payload
  })
})