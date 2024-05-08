import store from '@/store';

export async function login(http, body) {
	return request(http, { type: 'post', path: '/login', body });
}

export async function getTweets(http) {
	return request(http, { type: 'get', path: '/api/tweets' });
}

export async function uploadTweet(http, body) {
	return request(http, { type: 'post', path: '/api/tweets', body });
}

export async function deleteTweet(http, body) {
	return request(http, { type: 'delete', path: `/api/tweets/${body.tweetId}`, body });
}

export async function updateTweet(http, body) {
	return request(http, { type: 'patch', path: `/api/tweets/${body.id}`, body });
}

export async function likeTweet(http, tweetId) {
	return request(http, { type: 'post', path: `/api/tweets/${tweetId}/likes` });
}

export async function dislikeTweet(http, tweetId) {
	return request(http, { type: 'delete', path: `/api/tweets/${tweetId}/likes` });
}

export async function getUsersTweets(http, body) {
	return request(http, { type: 'get', path: `/api/tweets/user/${body.id}` });
}

export async function getMe(http, body) {
	return request(http, { type: 'get', path: '/api/users/me', body });
}

export async function updateUser(http, body) {
	return request(http, { type: 'put', path: '/api/users/me', body });
}

export async function uploadMedia(http, body) {
	return request(http, { type: 'post', path: '/api/medias', body });
}

export async function getTrends(http) {
	return request(http, { type: 'get', path: '/api/trends' });
}

async function request(http, settings) {
	store.commit('setLoadingStatus', true);
	try {
		if (settings.body) {
			const response = await http[settings.type](settings.path, settings.body);
			store.commit('setLoadingStatus', false);
			return response;
		}
		const response = await http[settings.type](settings.path);
		store.commit('setLoadingStatus', false);
		return response;
	} catch (err) {
		store.commit('setLoadingStatus', false);
		throw err;
	}
}
