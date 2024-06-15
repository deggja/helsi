import { createApp } from 'vue'
import App from './App.vue'
import axios from 'axios'

axios.defaults.baseURL = 'https://welcomed-usefully-porpoise.ngrok-free.app/api';

createApp(App).mount('#app')