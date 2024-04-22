<template>
	<component :is="iconComponent" />
</template>

<script>
const icons = {
	back: require('@/components/Icons/icons/back.vue'),
	bookmarks: require('@/components/Icons/icons/bookmarks.vue'),
	calendar: require('@/components/Icons/icons/calendar.vue'),
	close: require('@/components/Icons/icons/close.vue'),
	comment: require('@/components/Icons/icons/comment.vue'),
	editTweet: require('@/components/Icons/icons/editTweet.vue'),
	explore: require('@/components/Icons/icons/explore.vue'),
	gif: require('@/components/Icons/icons/gif.vue'),
	graph: require('@/components/Icons/icons/graph.vue'),
	hamburger: require('@/components/Icons/icons/hamburger.vue'),
	help: require('@/components/Icons/icons/help.vue'),
	home: require('@/components/Icons/icons/home.vue'),
	image: require('@/components/Icons/icons/image.vue'),
	left: require('@/components/Icons/icons/left.vue'),
	like: require('@/components/Icons/icons/like.vue'),
	link: require('@/components/Icons/icons/link.vue'),
	lists: require('@/components/Icons/icons/lists.vue'),
	messages: require('@/components/Icons/icons/messages.vue'),
	moments: require('@/components/Icons/icons/moments.vue'),
	more: require('@/components/Icons/icons/more.vue'),
	notifications: require('@/components/Icons/icons/notifications.vue'),
	pen: require('@/components/Icons/icons/pen.vue'),
	profile: require('@/components/Icons/icons/profile.vue'),
	retweet: require('@/components/Icons/icons/retweet.vue'),
	right: require('@/components/Icons/icons/right.vue'),
	schedule: require('@/components/Icons/icons/schedule.vue'),
	search: require('@/components/Icons/icons/search.vue'),
	settings: require('@/components/Icons/icons/settings.vue'),
	share: require('@/components/Icons/icons/share.vue'),
	tick: require('@/components/Icons/icons/tick.vue'),
	topics: require('@/components/Icons/icons/topics.vue'),
	twitter: require('@/components/Icons/icons/twitter.vue'),
};

export default {
	props: {
		name: {
			type: String,
			required: true,
			validator(value) {
				return Object.prototype.hasOwnProperty.call(icons, value);
			},
		},
	},

	computed: {
		iconComponent() {
			return icons[this.name];
		},
	},
};
</script>
