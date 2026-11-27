import { defineStore } from 'pinia'
import { ref } from 'vue'

/**
 * 全局 UI 状态。
 * 周报弹层目前只有监控页在用，别处需要再加状态时也放这里。
 */
export const useAppStore = defineStore(
  'app',
  () => {
    const weekly_ShowWeeklyModal = ref(false)

    function changeWeeklyModal(_index?: number) {
      weekly_ShowWeeklyModal.value = !weekly_ShowWeeklyModal.value
    }

    return {
      weekly_ShowWeeklyModal,
      changeWeeklyModal,
    }
  },
  {
    persist: true,
  },
)
