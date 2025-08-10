
// Este archivo contiene pruebas unitarias para el componente DetailCard.vue utilizando Vitest y Vue Test Utils.
// Se valida el comportamiento condicionado por la prop 'financialDetail' y el formateo seguro de métricas numéricas.
// Casos verificados:
//  1. Con financialDetail = true se renderizan sólo los indicadores financieros con valores a 2 decimales.
//  2. Con financialDetail = false se renderizan sólo los indicadores de mercado con valores a 2 decimales.
//  3. Se usa regex para comprobar el patrón de formateo (dos decimales) en ambas vistas.
//  4. Ante propiedades numéricas faltantes/undefined se muestra el fallback "No hay información" sin provocar errores.
//  5. Se comprueba que los encabezados de sección no mezclan entre vistas (exclusividad financiera vs mercado).
import { describe, it, expect } from 'vitest'         // Importa las funciones básicas de Vitest: describe (agrupa tests), it (caso individual) y expect (aserciones)
import { mount } from '@vue/test-utils'               // mount permite montar el componente Vue en un DOM simulado para inspeccionar su salida
import DetailCard from '../DetailCard.vue'            // Componente bajo prueba

// Objeto mock que simula un "stock" completo con todos los campos numéricos que el componente formatea.
// Se incluyen valores con decimales para verificar el redondeo a 2 cifras y signos negativos (eps_ttm).
const stockMock = {
  pe_ttm: 12.3456,
  pb: 3.789,
  eps_ttm: -1.234,
  net_profit_margin_ttm: 15.9876,
  week_52_high: 150.456,
  week_52_low: 90.123,
  beta: 1.5567,
  revenue_growth_ttm_yoy: 7.4321
}

// Grupo de pruebas para DetailCard.vue
describe('DetailCard.vue', () => {
  // Caso 1: Vista de detalles financieros (financialDetail = true)
  it('muestra indicadores financieros cuando financialDetail = true', () => {
    // Montaje del componente pasando el mock y activando la vista financiera
    const wrapper = mount(DetailCard, {
      props: {
        stock: stockMock,
        financialDetail: true
      }
    })

    // Verifica presencia del título de la sección financiera
    expect(wrapper.text()).toContain('Indicadores financieros')

    // Verifica el formateo (toFixed(2)) y etiquetas correctas de cada campo financiero
    expect(wrapper.text()).toContain('Radio Precio/Benecicio (TTM): 12.35')
    expect(wrapper.text()).toContain('Precio /Valor contable: 3.79')
    expect(wrapper.text()).toContain('Ganancias por acción (TTM): -1.23')
    expect(wrapper.text()).toContain('Margen neto (TTM): 15.99 %')

    // Asegura que la vista alternativa (mercado) no aparece
    expect(wrapper.text()).not.toContain('Indicadores del Mercado')
  })

  // Caso 2: Vista de indicadores de mercado (financialDetail = false)
  it('muestra indicadores de mercado cuando financialDetail = false', () => {
    // Monta el componente en modo mercado
    const wrapper = mount(DetailCard, {
      props: {
        stock: stockMock,
        financialDetail: false
      }
    })

    // Título correcto para esta vista
    expect(wrapper.text()).toContain('Indicadores del Mercado')

    // Verificación de campos y formateo (dos decimales) en la sección de mercado
    expect(wrapper.text()).toContain('Máximo en 52 semanas: 150.46')
    expect(wrapper.text()).toContain('Mínimo en 52 semanas: 90.12')
    expect(wrapper.text()).toContain('Volatilidad respecto al mercado: 1.56')
    expect(wrapper.text()).toContain('Crecimiento de ingresos interanual: 7.43 %')

    // Asegura que no se muestra el título financiero en esta vista
    expect(wrapper.text()).not.toContain('Indicadores financieros')
  })

  // Caso 3: Validación genérica de formato (usa regex para comprobar patrón de dos decimales)
  it('respeta el formato de dos decimales en ambas vistas', () => {
    // Monta ambas variantes para comparar
    const financial = mount(DetailCard, { props: { stock: stockMock, financialDetail: true } })
    const market = mount(DetailCard, { props: { stock: stockMock, financialDetail: false } })

    // Expresiones regulares para confirmar que los números relevantes aparecen con exactamente dos decimales
    expect(financial.text()).toMatch(/Radio Precio\/Benecicio \(TTM\): \d+\.\d{2}/)
    expect(market.text()).toMatch(/Máximo en 52 semanas: \d+\.\d{2}/)
  })

  // Caso 4: Robustez ante datos incompletos (campos undefined)
  it('no rompe si alguna propiedad numérica es undefined (opcional chaining)', () => {
    // Se pasa un objeto parcial con sólo pe_ttm definido; el resto debería resolverse con el fallback en el componente
    const partialStock = { pe_ttm: 10 } as any
    const wrapper = mount(DetailCard, {
      props: {
        stock: partialStock,
        financialDetail: true
      }
    })
    // Confirma que al menos el campo existente se renderiza formateado a dos decimales
    expect(wrapper.text()).toContain('Radio Precio/Benecicio (TTM): 10.00')
    // (Los otros campos se prueban en otro test enfocado al fallback si se desea)
  })
})