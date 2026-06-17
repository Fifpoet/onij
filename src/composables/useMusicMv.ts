import { computed, type Ref } from 'vue'
import { useMessage } from 'naive-ui'
import { UpdateMusicMv } from '@/api/music'
import { usePlayQueueStore } from '@/store/playQueue'
import { validateMvUrl } from '@/util/musicMv'

export function useMusicMv(songId: Ref<number | null | undefined>) {
  const message = useMessage()
  const playQueue = usePlayQueueStore()

  const mvUrl = computed(() => playQueue.nowPlaying?.mv_url?.trim() ?? '')
  const hasMv = computed(() => validateMvUrl(mvUrl.value).ok && !!mvUrl.value)

  async function saveMvUrl(raw: string) {
    const id = songId.value
    if (!id) return

    const result = validateMvUrl(raw)
    if (!result.ok) {
      message.warning(result.error ?? 'URL 无效')
      return
    }

    try {
      await UpdateMusicMv({
        third_id: id,
        mv_url: result.url,
      })
    } catch {
      /* 后端未就绪时仍更新本地队列 */
    }
    playQueue.patchMvUrl(id, result.url)
  }

  return {
    mvUrl,
    hasMv,
    saveMvUrl,
  }
}
