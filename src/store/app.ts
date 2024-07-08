import { defineStore } from "pinia";
import {ref} from "vue";


export const useAppStore = defineStore('app', {
  state: () => ({
    showWeeklyModal: ref(false)
  }),
  getters: {

  }, 
  actions: {
    changeWeeklyModal() {
      this.showWeeklyModal = !this.showWeeklyModal
      console.log('I have changeWeeklyModal', this.showWeeklyModal)
    }
  }
})