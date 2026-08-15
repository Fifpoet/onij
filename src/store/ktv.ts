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
    /** 音高线降一个八度，方便跟唱；伴奏不移调 */
    const octaveDown = ref(false)

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

    function toggleOctaveDown() {
      octaveDown.value = !octaveDown.value
    }

    return {
      active,
      lyricsOn,
      accompanimentOn,
      octaveDown,
      enterKtv,
      leaveKtv,
      toggleLyrics,
      toggleAccompaniment,
      toggleOctaveDown,
    }
  },
  {
    persist: {
      paths: ['lyricsOn', 'accompanimentOn', 'octaveDown'],
    },
  },
)
