import { defineStore } from "pinia";
import {ref} from "vue";


export const useAppStore = defineStore('app', {
  state: () => ({
    weekly_ShowWeeklyModal: ref(false)
  }),
  getters: {

  }, 
  actions: {
    changeWeeklyModal(index: number) {
      this.weekly_ShowWeeklyModal = !this.weekly_ShowWeeklyModal
    }
  }
})