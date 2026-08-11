import type {
  ApiResponse,
  LoginCredentials,
  CreateTweetPayload,
  UpdateTweetPayload,
  UpdateUserPayload,
  PageQuery,
  Tweet,
  TweetLikes,
  User,
  UserBrief,
  Trend,
  Media,
} from '~/types'

/** Page size used by every list view. Matches the backend default. */
export const PAGE_SIZE = 20

export const useApi = () => {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBaseUrl as string
  const authStore = useAuthStore()

  const apiFetch = <T>(url: string, options: Parameters<typeof $fetch>[1] = {}) => {
    return $fetch<ApiResponse<T>>(url, {
      baseURL,
      headers: {
        ...(authStore.token ? { Authorization: `Bearer ${authStore.token}` } : {}),
      },
      ...options,
    })
  }

  const login = (credentials: LoginCredentials) =>
    apiFetch<string>('/login', { method: 'POST', body: credentials })

  const logout = () => apiFetch<null>('/api/logout', { method: 'POST' })

  const getMe = () => apiFetch<User>('/api/users/me')

  const getUser = (id: number | string) => apiFetch<User>(`/api/users/${id}`)

  const updateUser = (data: UpdateUserPayload) =>
    apiFetch<User>('/api/users/me', { method: 'PUT', body: data })

  const followUser = (id: number | string) =>
    apiFetch<User>(`/api/users/${id}/follow`, { method: 'POST' })

  const unfollowUser = (id: number | string) =>
    apiFetch<User>(`/api/users/${id}/follow`, { method: 'DELETE' })

  const getFollowers = (id: number | string, query: PageQuery = {}) =>
    apiFetch<UserBrief[]>(`/api/users/${id}/followers`, { query })

  const getFollowings = (id: number | string, query: PageQuery = {}) =>
    apiFetch<UserBrief[]>(`/api/users/${id}/followings`, { query })

  const getTweets = (query: PageQuery = {}) => apiFetch<Tweet[]>('/api/tweets', { query })

  const getUserTweets = (userId: number | string, query: PageQuery = {}) =>
    apiFetch<Tweet[]>(`/api/tweets/user/${userId}`, { query })

  const getTweet = (id: number) => apiFetch<Tweet>(`/api/tweets/${id}`)

  const createTweet = (data: CreateTweetPayload) =>
    apiFetch<Tweet>('/api/tweets', { method: 'POST', body: data })

  const deleteTweet = (id: number) =>
    apiFetch<null>(`/api/tweets/${id}`, { method: 'DELETE' })

  const updateTweet = (data: UpdateTweetPayload) =>
    apiFetch<Tweet>(`/api/tweets/${data.id}`, { method: 'PATCH', body: data })

  const likeTweet = (tweetId: number) =>
    apiFetch<TweetLikes>(`/api/tweets/${tweetId}/likes`, { method: 'POST' })

  const unlikeTweet = (tweetId: number) =>
    apiFetch<TweetLikes>(`/api/tweets/${tweetId}/likes`, { method: 'DELETE' })

  const uploadMedia = (formData: FormData) =>
    apiFetch<Media>('/api/medias', { method: 'POST', body: formData })

  const getTrends = (query: PageQuery = {}) => apiFetch<Trend[]>('/api/trends', { query })

  return {
    login,
    logout,
    getMe,
    getUser,
    updateUser,
    followUser,
    unfollowUser,
    getFollowers,
    getFollowings,
    getTweets,
    getUserTweets,
    getTweet,
    createTweet,
    deleteTweet,
    updateTweet,
    likeTweet,
    unlikeTweet,
    uploadMedia,
    getTrends,
  }
}
