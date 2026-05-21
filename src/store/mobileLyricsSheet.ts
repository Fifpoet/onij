import { defineStore } from 'pinia'
import { ref } from 'vue'

/** 手机端底部上滑歌词详情 */
export const useMobileLyricsSheetStore = defineStore('mobileLyricsSheet', () => {
  const visible = ref(false)
  /** false：封面 + 操作；true：点击封面后展示滚动歌词 */
  const showLyricsScroll = ref(false)

  function open() {
    visible.value = true
    showLyricsScroll.value = false
  }

  function close() {
    visible.value = false
    showLyricsScroll.value = false
  }

  function enterLyricsScroll() {
    showLyricsScroll.value = true
  }

  function backToCover() {
    showLyricsScroll.value = false
  }

  return { visible, showLyricsScroll, open, close, enterLyricsScroll, backToCover }
})
