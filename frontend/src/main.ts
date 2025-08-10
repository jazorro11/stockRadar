// Importa los estilos globales de la aplicación.
import './assets/main.css'

// Importa la función para crear la aplicación Vue.
import { createApp } from 'vue'
// Importa la función para crear el store global Pinia.
import { createPinia } from 'pinia'
// Importa el componente raíz de la aplicación.
import App from './App.vue'
// Importa el enrutador de la aplicación.
import router from './router'

// Crea la instancia principal de la aplicación Vue usando el componente App.
const app = createApp(App)

// Registra Pinia como sistema de gestión de estado global.
app.use(createPinia())
// Registra el enrutador para la navegación entre vistas.
app.use(router)

// Monta la aplicación en el elemento con id 'app' en el