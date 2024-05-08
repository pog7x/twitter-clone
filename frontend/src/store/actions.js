import { defaultUser } from '@/services/functions';

export default {
	setLoginInfo({ commit }, payload) {
		commit('setMe', payload);
		commit('setLoginStatus', true);
	},
	setLogOut({ commit }) {
		commit('setMe', defaultUser());
		commit('setLoginStatus', false);
	},
	setMyInfo({ commit }, payload) {
		commit('editProfileInfo', payload);
		commit('setEditProfileStatus', false);
	},
	setLightbox({ commit }, payload) {
		commit('setLightboxState', true);
		commit('setLightboxImages', payload.tweetImages);
		commit('setLightboxIndex', payload.index);
	},
	closeLightbox({ commit }) {
		commit('setLightboxState', false);
		commit('setLightboxImages', []);
		commit('setLightboxIndex', 0);
	},
};
