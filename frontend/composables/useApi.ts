import type {
  ApiResponse,
  LoginCredentials,
  CreateTweetPayload,
  UpdateTweetPayload,
  UpdateUserPayload,
  Tweet,
  User,
  Trend,
  Media,
} from '~/types'

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

  const getMe = () => apiFetch<User>('/api/users/me')

  const getUser = (id: number | string) => apiFetch<User>(`/api/users/${id}`)

  const updateUser = (data: UpdateUserPayload) =>
    apiFetch<User>('/api/users/me', { method: 'PUT', body: data })

  const getTweets = () => apiFetch<Tweet[]>('/api/tweets')

  const getUserTweets = (userId: number | string) =>
    apiFetch<Tweet[]>(`/api/tweets/user/${userId}`)

  const createTweet = (data: CreateTweetPayload) =>
    apiFetch<Tweet>('/api/tweets', { method: 'POST', body: data })

  const deleteTweet = (id: number) =>
    apiFetch<void>(`/api/tweets/${id}`, { method: 'DELETE' })

  const updateTweet = (data: UpdateTweetPayload) =>
    apiFetch<Tweet>(`/api/tweets/${data.id}`, { method: 'PATCH', body: data })

  const likeTweet = (tweetId: number) =>
    apiFetch<void>(`/api/tweets/${tweetId}/likes`, { method: 'POST' })

  const unlikeTweet = (tweetId: number) =>
    apiFetch<void>(`/api/tweets/${tweetId}/likes`, { method: 'DELETE' })

  const uploadMedia = (formData: FormData) =>
    apiFetch<Media>('/api/medias', { method: 'POST', body: formData })

  const getTrends = () => apiFetch<Trend[]>('/api/trends')

  return {
    login,
    getMe,
    getUser,
    updateUser,
    getTweets,
    getUserTweets,
    createTweet,
    deleteTweet,
    updateTweet,
    likeTweet,
    unlikeTweet,
    uploadMedia,
    getTrends,
  }
}
