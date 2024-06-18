import { defineStore } from "pinia";


export const useAppStore = defineStore('app', {
  state: () => ({
    searchFlag: false,
    loginFlag: false,
    registerFlag: false,
    collapsed: false, // 侧边栏折叠（移动端）
  }),
  getters: () => {

  }, 
  actions: () => {
    
  }
})