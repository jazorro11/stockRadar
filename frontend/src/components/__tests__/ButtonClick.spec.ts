/**
 * ButtonClick.spec.ts
 * Test para ButtonClick.vue — comprobamos:
 *  1) que el botón emite el evento personalizado 'button-click' exactamente una vez al hacer clic
 *  2) que muestra el texto "Buscar"
 */

import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ButtonClick from '../ButtonClick.vue'

describe('ButtonClick.vue', () => {

  it('emite el evento "button-click" una sola vez al hacer clic', async () => {
    // Montamos el componente en memoria
    const wrapper = mount(ButtonClick)

    // Encontramos el <button> y simulamos un clic.
    // Usamos trigger('click') porque el evento emitido ahora tiene nombre único
    // y no colisiona con el evento DOM nativo.
    await wrapper.find('button').trigger('click')

    // Recuperamos las emisiones del evento 'button-click'
    const emitted = wrapper.emitted('button-click')

    // Comprobamos que se haya emitido algo
    expect(emitted).toBeTruthy()

    // Comprobamos que se haya emitido exactamente una vez
    expect(emitted!.length).toBe(1)
  })

  it('muestra el texto "Buscar" en el botón', () => {
    // Montamos de nuevo (test aislado)
    const wrapper = mount(ButtonClick)

    // Buscamos el <span> y verificamos su texto
    expect(wrapper.find('span').text()).toBe('Buscar')
  })
})
