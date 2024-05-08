import { getCookie } from '@/services/cookies';
import axios from 'axios';

export const instance = axios.create({
	baseURL: 'http://0.0.0.0:8080',
	headers: {
		Authorization: `Bearer ${getCookie('access_token')}`,
	},
});
