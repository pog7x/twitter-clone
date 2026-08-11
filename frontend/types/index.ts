/** Shape embedded into other resources: a tweet author, a followers row. */
export interface UserBrief {
  id: number
  username: string
  nickname: string
  name: string
  pic: string
  pic_cover: string
  created_at: string
}

/** Full profile: the brief fields plus aggregates and viewer-relative flags. */
export interface User extends UserBrief {
  website: string
  description: string
  followers_count: number
  followings_count: number
  tweets_count: number
  is_following: boolean
  is_me: boolean
}

export interface Tweet {
  id: number
  content: string
  author: UserBrief
  attachments: string[]
  likes_count: number
  is_liked: boolean
  is_owner: boolean
  created_at: string
}

/** Returned by the like endpoints so a card can be updated in place. */
export interface TweetLikes {
  tweet_id: number
  likes_count: number
  is_liked: boolean
}

export interface Trend {
  id: number
  name: string
  tweets_count: number
  created_at: string
}

export interface Media {
  media_id: number
  link: string
}

export interface ApiResponse<T> {
  success: boolean
  result: T
}

export interface ApiErrorResponse {
  success: boolean
  error_type: string
  error_message: string
}

export interface LoginCredentials {
  username: string
  password: string
}

/** Query params accepted by every list endpoint. */
export interface PageQuery {
  limit?: number
  offset?: number
}

export interface CreateTweetPayload {
  tweet_data: string
  tweet_media_ids: number[]
}

export interface UpdateTweetPayload {
  id: number
  tweet_data: string
}

/**
 * Every field is optional: omitting one keeps the current value, sending an
 * empty string clears it. Images are referenced by the media id returned from
 * the upload endpoint.
 */
export interface UpdateUserPayload {
  name?: string
  description?: string
  website?: string
  pic_media_id?: number
  pic_cover_media_id?: number
}
