import type { ApiResponse, PageQuery, Tweet } from '~/types'
import { PAGE_SIZE } from '~/composables/useApi'

type TweetPageFetcher = (query: PageQuery) => Promise<ApiResponse<Tweet[]>>

/**
 * Drives an offset-paginated tweet list: the first page replaces the list,
 * every following one is appended. A page shorter than PAGE_SIZE means the end
 * of the list has been reached, which is how "load more" knows to disappear.
 */
export const usePaginatedTweets = (fetchPage: TweetPageFetcher) => {
  const toast = useToast()

  const tweets = ref<Tweet[]>([])
  const isLoading = ref(false)
  const hasMore = ref(false)

  const load = async (offset: number) => {
    if (isLoading.value) return

    isLoading.value = true
    try {
      const { result } = await fetchPage({ limit: PAGE_SIZE, offset })
      const batch = result ?? []

      tweets.value = offset === 0 ? batch : [...tweets.value, ...batch]
      hasMore.value = batch.length === PAGE_SIZE
    } catch {
      toast.add({
        title: 'Error',
        description: 'Failed to load posts',
        color: 'error',
        icon: 'i-lucide-alert-circle',
      })
    } finally {
      isLoading.value = false
    }
  }

  const refresh = () => load(0)

  const loadMore = () => load(tweets.value.length)

  return { tweets, isLoading, hasMore, refresh, loadMore }
}
