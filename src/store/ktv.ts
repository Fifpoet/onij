import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useKtvStore = defineStore(
  'ktv',
  () => {
    const active = ref(false)
    /** 显示歌词（词） */
    const lyricsOn = ref(true)
    /** 播放伴奏（伴） */
    const accompanimentOn = ref(true)

    function enterKtv() {
      active.value = true
    }

    function leaveKtv() {
      active.value = false
    }

    function toggleLyrics() {
      lyricsOn.value = !lyricsOn.value
    }

    function toggleAccompaniment() {
      accompanimentOn.value = !accompanimentOn.value
    }

    return {
      active,
      lyricsOn,
      accompanimentOn,
      enterKtv,
      leaveKtv,
      toggleLyrics,
      toggleAccompaniment,
    }
  },
  {
    persist: {
      paths: ['lyricsOn', 'accompanimentOn'],
    },
  },
)
