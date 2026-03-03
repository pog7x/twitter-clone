export interface User {
  id: number
  username: string
  name: string
  nickname: string
  description: string
  website: string
  pic: string
  pic_cover: string
  created_at: string
  followers: User[]
  followings: User[]
}

export interface Tweet {
  id: number
  content: string
  author: User
  likes: Like[]
  attachments: string[]
  created_at: string
}

export interface Like {
  id: number
  user_id: number
  tweet_id: number
}

export interface Trend {
  id: number
  name: string
  tweets_count: number
}

export interface Media {
  media_id: number
}

export interface ApiResponse<T> {
  success: boolean
  result: T
}

export interface LoginCredentials {
  username: string
  password: string
}

export interface CreateTweetPayload {
  tweet_data: string
  tweet_media_ids: number[]
}

export interface UpdateTweetPayload {
  id: number
  tweet_data: string
}

export interface UpdateUserPayload {
  name: string
  description: string
  website: string
}
