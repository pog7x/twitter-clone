import { createApp } from 'vue';
import App from '@/App.vue';
import router from '@/router';
import store from '@/store';
import { notification } from '@/utils/notification';
import axios from 'axios';
import { getCookie } from '@/services/cookies';

const app = createApp(App);

app.config.globalProperties.$notification = notification;

app.use(router);
app.use(store);
app.mount('#app');

app.config.globalProperties.axios = axios.create({
	baseURL: 'http://0.0.0.0:8080',
	headers: {
		Authorization: 'Bearer ' + getCookie('access_token'),
	},
});
