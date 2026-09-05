import { defineStore } from 'pinia'
import { ref } from 'vue'
import { neteasePlayerControl } from '@/player/neteasePlayerControl'

export const useVideoPlayerStore = defineStore('videoPlayer', () => {
  const fileId = ref(0)
  const opened = ref(false)
  const currentTime = ref(0)
  const duration = ref(0)

  function open(id: number) {
    fileId.value = id
    opened.value = true
    currentTime.value = 0
    duration.value = 0
    neteasePlayerControl.pause()
  }

  function close() {
    opened.value = false
    fileId.value = 0
    currentTime.value = 0
    duration.value = 0
  }

  function setProgress(t: number, d: number) {
    currentTime.value = t
    duration.value = d
  }

  return { fileId, opened, currentTime, duration, open, close, setProgress }
})
