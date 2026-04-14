import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useLyricsPanelStore = defineStore('lyricsPanel', () => {
  const visible = ref(false)

  function open() {
    visible.value = true
  }

  function close() {
    visible.value = false
  }

  function toggle() {
    visible.value = !visible.value
  }

  return { visible, open, close, toggle }
})
