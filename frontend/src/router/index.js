import { getMe } from '@/services/api';
import { instance } from '@/services/http';
import store from '@/store';
import { createRouter, createWebHistory } from 'vue-router';

const routes = [
	{
		path: '/',
		name: 'Home',
		component: () => import('@/views/Home.vue'),
	},
	{
		path: '/login',
		name: 'Login',
		beforeEnter: loginGuardian,
		component: () => import('@/views/Login.vue'),
	},
	{
		path: '/profile/:profileId',
		name: 'Profile',
		component: () => import('@/views/Profile.vue'),
		props: true,
	},
];

const router = createRouter({
	history: createWebHistory(),
	routes,
});

router.beforeEach(async (to, from, next) => {
	if (to.path !== '/login') {
		try {
			const response = await getMe(instance, {});
			store.commit('setLoginStatus', true);
			store.commit('setMe', response.data.result);
			next();
		} catch (err) {
			next({ path: '/login' });
		}
		return
	}
	next();
});

function loginGuardian(to, from, next) {
	const isLoggedIn = store.getters.getLoginStatus;
	if (isLoggedIn) {
		next(from);
	} else {
		next();
	}
}

export default router;
