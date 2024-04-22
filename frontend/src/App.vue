<template>
	<component :is="getLoginStatus ? 'layout' : 'div'">
		<router-view />
		<loading v-if="getLoadingStatus" />
		<notification v-for="(notification, i) in getActiveNotifications" :key="i" :index="i" :data="notification" />
	</component>
</template>

<script>
import 'normalize.css';
import Layout from '@/views/Layout.vue';
import Loading from '@/components/Loading.vue';
import Notification from '@/components/Notification.vue';
import { mapGetters } from 'vuex';

export default {
	components: {
		Layout,
		Loading,
		Notification,
	},
	data() {
		return {
			globalNotification: false,
		};
	},
	computed: {
		...mapGetters(['getLoginStatus', 'getLoadingStatus', 'getActiveNotifications']),
	},
};
</script>

<style lang="scss">
* {
	box-sizing: border-box;
	&::before,
	&::after {
		box-sizing: inherit;
	}
}

html,
body,
:root {
	font-size: 14px;
}
@import 'src/assets/theme/colors.scss';

body {
	background-color: $color-bg;
	font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
}
</style>
