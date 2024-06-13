<template>
	<div class="profile-body">
		<div class="sections">
			<div class="sections-item active">Tweets</div>
			<div class="sections-item">Tweets & replies</div>
			<div class="sections-item">Media</div>
			<div class="sections-item">Likes</div>
		</div>
		<div v-if="userTweets" class="tweets-wrapper">
			<tweet v-for="tweet in userTweets" :key="tweet.id" :tweet-data="tweet" @delete-tweet="handleTweetDelete" />
		</div>
	</div>
</template>

<script>
import Tweet from '@/components/Tweet/Tweet.vue';
import { getUsersTweets } from '@/services/api';

export default {
	name: 'ProfileBody',
	components: {
		Tweet,
	},
	props: {
		profileId: {
			type: String,
			required: true,
		},
	},
	watch: {
		profileId(newValue) {
			this.getUsersTweets();
		},
	},
	data() {
		return {
			userTweets: [],
		};
	},
	mounted() {
		this.getUsersTweets();
	},
	methods: {
		handleTweetDelete() {
			this.getUsersTweets();
		},
		async getUsersTweets() {
			try {
				const response = await getUsersTweets(this.axios, {
					id: this.profileId,
				});
				this.userTweets = response.data.result;
				this.$store.commit('setProfileTweetCount', response.data.result?.length || 0);
			} catch (err) {
				this.$notification({
					type: 'error',
					message: 'Error when fetching tweets',
				});
			}
		},
	},
};
</script>

<style lang="scss">
@import 'src/assets/theme/colors.scss';

.profile-body {
	.sections {
		border-bottom: $border-dark;
		display: flex;
		&-item {
			width: calc(100% / 4);
			text-align: center;
			padding: 1.5rem 0;
			color: $color-dark-gray;
			font-weight: bold;
			&.active {
				border-bottom: 2px solid $color-blue;
				color: $color-blue;
			}
		}
	}
}
</style>
