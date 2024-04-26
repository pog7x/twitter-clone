<template>
	<div class="home">
		<add-tweet @submit-click="handleTweetSubmit" />
		<hr class="gap" />
		<div v-if="tweetData" class="tweet-wrapper">
			<tweet v-for="tweet in tweetData" :key="tweet.id" :tweet-data="tweet" @delete-tweet="handleTweetDelete" @get-tweets="getTweets" />
		</div>
	</div>
</template>

<script>
import AddTweet from '@/components/AddTweet.vue';
import Tweet from '@/components/Tweet/Tweet.vue';
import { getTweets } from '@/services/api';
import { mapGetters } from 'vuex';

export default {
	components: {
		AddTweet,
		Tweet,
	},
	data: function () {
		return {
			tweetData: [],
		};
	},
	computed: {
		...mapGetters(['getMe']),
	},
	mounted: async function () {
		this.getTweets();
	},
	methods: {
		async handleTweetSubmit() {
			try {
				await this.getTweets();
			} catch (err) {
				this.$notification({
					type: 'error',
					message: 'Error in get tweet',
				});
			}
		},
		async getTweets() {
			const response = await getTweets(this.axios);
			this.tweetData = response.data.result;
		},
		async handleTweetDelete() {
			this.handleTweetSubmit();
		},
	},
};
</script>

<style lang="scss">
@import 'src/assets/theme/colors.scss';

hr.gap {
	background-color: $hr-color;
	padding: 4px 0;
	margin: 0;
	border: none;
}
</style>
