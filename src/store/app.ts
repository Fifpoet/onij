import { defineStore } from "pinia";


export const useAppStore = defineStore('app', {
  state: () => ({
    searchFlag: false,
    loginFlag: false,
    registerFlag: false,
    collapsed: false, // 侧边栏折叠（移动端）

    page_list: [], // 页面数据
    blogInfo: {
      article_count: 0,
      category_count: 0,
      tag_count: 0,
      view_count: 0,
      user_count: 0,
    },
    blog_config: {
      website_name: '阵、雨的个人博客',
      website_author: '阵、雨',
      website_intro: '往事随风而去',
      website_avatar: '',
    },
  }),
  getters: () => {

  }, 
  actions: () => {
    
  }
})