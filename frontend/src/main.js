import { createApp } from 'vue';
import App from '@/App.vue';
import router from '@/router';
import store from '@/store';
import { notification } from '@/services/notification';
import { instance } from '@/services/http';
import mitt from 'mitt';

const emitter = mitt();

const app = createApp(App);

app.config.globalProperties.$notification = notification;
app.config.globalProperties.emitter = emitter;
app.config.globalProperties.axios = instance;
app.config.globalProperties.baseUrl = import.meta.env.VITE_BASE_API_URL;

app.use(router);
app.use(store);
app.mount('#app');
