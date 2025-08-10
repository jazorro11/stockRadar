<script setup>
// Importa el componente Footer, que se mostrará en la parte inferior de la aplicación.
import Footer from './components/Footer.vue'
// Importa el componente Header, que se mostrará en la parte superior de la aplicación.
import Header from './components/Header.vue'
// Las siguientes importaciones están comentadas porque no se usan directamente aquí:
// import { RouterLink, RouterView } from 'vue-router'
// import Content from './components/Content.vue'
// import HelloWorld from './components/HelloWorld.vue'
</script>

<template>
  <!-- Contenedor principal de la aplicación -->
  <div>
    <!-- Header de la aplicación, siempre visible en la parte superior -->
    <Header />
    
    <main>
      <!-- Aquí se renderizan las vistas según la ruta activa usando Vue Router -->
      <!-- <Content /> está comentado porque el contenido principal se gestiona por rutas -->
      <router-view />
    </main>

    <!-- Footer de la aplicación, siempre visible en la parte inferior -->
    <Footer />
  </div>
</template>