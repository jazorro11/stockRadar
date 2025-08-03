
<script setup>
// import { RouterLink, RouterView } from 'vue-router'
import Footer from './components/Footer.vue'
import Header from './components/Header.vue'
import Content from './components/Content.vue'
// import HelloWorld from './components/HelloWorld.vue'
</script>

<template>
  <div class="app-container">
    <Header />
    
    <main>
      <!-- Aquí iría el contenido principal o <RouterView /> si usas rutas -->
      <Content />
    </main>

    <Footer />
  </div>
</template>
