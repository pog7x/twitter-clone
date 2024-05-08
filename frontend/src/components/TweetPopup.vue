<template>
	<div ref="popupWrapper" class="tweet-popup" @click="handleClickOutside">
		<div ref="popup" class="tweet-popup-wrapper">
			<div class="tweet-popup-header">
				<div class="close-button" @click="$store.commit('toggleTweetButton')">
					<base-icon name="close" />
				</div>
			</div>
			<add-tweet @submit-click="handleSubmit" />
		</div>
	</div>
</template>

<script>
import AddTweet from '@/components/AddTweet.vue';
import BaseIcon from '@/components/Icons/BaseIcon.vue';
import { getTweets } from '@/services/api';

export default {
	name: 'TweetPopup',
	components: {
		AddTweet,
		BaseIcon,
	},
	methods: {
		handleClickOutside(e) {
			const object = {
				target: e.target,
				ref: this.$refs.popupWrapper,
			};
			if (object.target !== object.ref) return;
			this.$store.commit('toggleTweetButton');
		},
		async handleSubmit() {
			this.$store.commit('toggleTweetButton');
			this.$store.commit('setMobileMenuState', false);

			this.emitter.emit('get-tweets', '');
		},
		async getTweets() {
			const response = await getTweets(this.axios);
			this.tweetData = response.data.result;
		},
	},
};
</script>

<style lang="scss">
@import 'src/assets/theme/colors.scss';
@import 'src/assets/variables.scss';

.tweet-popup {
	position: fixed;
	top: 0;
	left: 0;
	height: 100vh;
	width: 100%;
	z-index: 9999;
	display: flex;
	align-items: flex-start;
	justify-content: center;
	background-color: rgba($color: $color-dark-gray, $alpha: 0.3);
	&-header {
		padding: 0.5rem;
		.close-button {
			margin-left: auto;
			width: 2rem;
			width: 2rem;
			padding: 5px;
			cursor: pointer;
			border-radius: 999px;
			svg {
				width: 100%;
				height: 100%;
				fill: $tweet-action-red;
			}
			&:hover {
				background: rgba($color: $tweet-action-red, $alpha: 0.3);
			}
		}
	}
	&-wrapper {
		margin-top: 50px;
		border-radius: 1rem;
		max-width: 450px;
		max-height: 50rem;
		background-color: $color-bg;
	}
}
@media screen and (max-width: $phone) {
	.tweet-popup {
		&-wrapper {
			min-width: unset;
			width: 90%;
		}
	}
}
</style>
