import { createApp } from 'vue';
import App from '@/App.vue';
import router from '@/router';
import store from '@/store';
import { notification } from '@/utils/notification';
import { instance } from '@/services/http';
import mitt from 'mitt';

const emitter = mitt();

const app = createApp(App);

app.config.globalProperties.$notification = notification;
app.config.globalProperties.emitter = emitter;
app.config.globalProperties.axios = instance;

app.use(router);
app.use(store);
app.mount('#app');
