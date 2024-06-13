<template>
	<div class="profile">
		<profile-header :profileId="this.profileId" />
		<profile-body :profileId="this.profileId" />
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
	props: {
		profileId: {
			type: String,
			required: true,
		},
	},
	computed: {
		...mapGetters(['getMyProfileId', 'getEditProfileStatus']),
	},
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
