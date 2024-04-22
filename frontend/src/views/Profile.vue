<template>
	<div class="profile">
		<profile-header />
		<profile-body />
		<EditProfilePopup v-if="getEditProfileStatus" />
	</div>
</template>

<script>
import ProfileBody from '@/components/Profile/ProfileBody.vue';
import ProfileHeader from '@/components/Profile/ProfileHeader.vue';
import EditProfilePopup from '@/components/EditProfilePopup.vue';
import { mapGetters } from 'vuex';

export default {
	name: 'ProfileView',
	components: {
		ProfileBody,
		ProfileHeader,
		EditProfilePopup,
	},
	data() {
		return {};
	},
	computed: {
		...mapGetters(['getMyProfileId', 'getEditProfileStatus']),
	},
	async mounted() {},
};
</script>

<style lang="scss">
@import 'src/assets/theme/colors.scss';

.profile {
	.profile-cover-pic {
		img {
			width: 100%;
		}
	}
	&-header {
		padding: 1rem;
	}
}
</style>
