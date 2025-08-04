
<script setup>
// import { RouterLink, RouterView } from 'vue-router'
import Footer from './components/Footer.vue'
import Header from './components/Header.vue'
// import Content from './components/Content.vue'
// import HelloWorld from './components/HelloWorld.vue'
</script>

<template>
  <div>
    <Header />
    
    <main>
      <!-- <Content /> -->
       <router-view />
    </main>

    <Footer />
  </div>
</template>
