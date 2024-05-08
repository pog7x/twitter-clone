import { getCookie } from '@/services/cookies';
import axios from 'axios';

export const instance = axios.create({
	baseURL: import.meta.env.VITE_BASE_API_URL,
	headers: {
		Authorization: `Bearer ${getCookie('access_token')}`,
	},
});
