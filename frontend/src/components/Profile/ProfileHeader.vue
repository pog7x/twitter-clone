<template>
	<header v-if="profile.id">
		<div class="profile-cover-pic">
			<img :src="profile.pic_cover ? baseUrl + profile.pic_cover : ''" />
		</div>
		<div class="profile-header">
			<div class="profile-actions">
				<div class="profile-actions-image">
					<img :src="profile.pic ? baseUrl + profile.pic : ''" />
				</div>
				<div v-if="isMe" class="profile-actions-edit">
					<div class="edit-button" @click="$store.commit('setEditProfileStatus', true)">Edit profile</div>
				</div>
			</div>
			<div class="profile-info">
				<p class="profile-info-name">
					{{ profile.name }}
				</p>
				<span class="profile-info-username">
					{{ profile.nickname }}
				</span>
			</div>
			<div class="profile-description">
				{{ profile.description }}
			</div>
			<div class="profile-created-at">
				<span>
					<base-icon name="link" />
					<a :href="profileWebsite.full_website">{{ profileWebsite.website }}</a>
				</span>
				<span>
					<base-icon name="calendar" />
					Joined {{ joinedAtDate }}
				</span>
			</div>
			<div class="profile-follower-counts">
				<p>
					{{ profile.followings?.length }}
					<span>Following</span>
				</p>
				<p>
					{{ profile.followers?.length }}
					<span>Followers</span>
				</p>
			</div>
		</div>
	</header>
</template>

<script>
import BaseIcon from '@/components/Icons/BaseIcon.vue';
import { getUser } from '@/services/api';
import moment from 'moment';
import { mapGetters } from 'vuex';

export default {
	name: 'ProfileHeader',
	components: {
		BaseIcon,
	},
	props: {
		profileId: {
			type: String,
			required: true,
		},
	},
	data() {
		return { profile: {} };
	},
	watch: {
		profileId(newValue) {
			this.profile = this.isMe ? this.me : this.profile;
		},
	},
	computed: {
		...mapGetters({
			me: 'getMe',
		}),
		profileWebsite() {
			return {
				website: new URL(new URL(this.profile.website)).host,
				full_website: this.profile.website,
			};
		},
		joinedAtDate() {
			return `${moment(this.profile.created_at).format('MMM YYYY')}`;
		},
		isMe() {
			return this.me.id === Number(this.profileId);
		},
	},
	async mounted() {
		if (this.isMe) {
			this.profile = this.me;
			return;
		}
		try {
			const response = await getUser(this.axios, { id: this.profileId });
			this.profile = response.data.result;
		} catch (err) {
			this.$notification({
				type: 'error',
				message: 'Error when fetching user data',
			});
		}
	},
};
</script>

<style lang="scss">
@import 'src/assets/theme/colors.scss';

.profile {
	&-cover-pic {
		border-bottom: $border-dark;
		img {
			vertical-align: middle;
		}
	}
	&-actions {
		display: flex;
		align-items: center;
		justify-content: space-between;
		&-image {
			width: 130px;
			height: 130px;
			margin-top: -80px;
			img {
				border-radius: 999px;
				width: 100%;
			}
		}
		&-edit {
			.edit-button {
				border-radius: 999px;
				border: 1px solid $color-blue;
				color: $color-blue;
				font-weight: bold;
				font-size: 1rem;
				padding: 1rem;
				cursor: pointer;
				transition: background-color 80ms ease;
				&:hover {
					background-color: rgba($color: $color-blue, $alpha: 0.1);
				}
			}
		}
	}
	&-info {
		margin-top: 1rem;
		&-name {
			color: #fff;
			margin: 0;
			font-weight: bold;
			font-size: 1.5rem;
		}
		&-username {
			font-size: 1.2rem;
			color: $color-dark-gray;
		}
	}
	&-description {
		margin-top: 1rem;
		color: #fff;
	}
	&-created-at {
		margin-top: 1rem;
		display: flex;
		align-items: center;
		color: $color-dark-gray;
		a {
			color: $color-blue;
			&:hover {
				text-decoration: underline;
			}
		}
		span {
			display: flex;
			align-items: center;
			& + span {
				margin-left: 2rem;
			}
			svg {
				fill: $color-dark-gray;
				margin-right: 0.5rem;
				width: 1.2rem;
				height: 1.2rem;
			}
		}
	}
	&-follower-counts {
		display: flex;
		color: #fff;
		cursor: pointer;
		margin-top: 1rem;
		p {
			margin: 0;
			& + p {
				margin-left: 1rem;
			}
			span {
				color: $color-dark-gray;
			}
			&:hover {
				text-decoration: underline;
			}
		}
	}
}
</style>
