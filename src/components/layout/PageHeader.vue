<script setup lang="ts">
import { ref } from 'vue'
import { Icon } from '@iconify/vue';
import { useAppStore, useUserStore } from '@/store';
import { RouterLink } from 'vue-router';
import {NButton} from 'naive-ui'
import {NDropdown} from 'naive-ui'
import { h } from 'vue'
// import MobileSideBar from './MobileSideBar.vue';

const appStore = useAppStore()
const userStore = useUserStore()

const navClass = ref('nav')
const barShow = ref(true)



const renderIcon = (icon: string) => {
  return () => {
    return h(Icon, {
      icon: icon,
      class: 'text-xl'
    }, {
      default: () => h(icon)
    })
  }
}


const options = [
  {
    label: '用户资料',
    key: 'profile',
    icon: renderIcon('mdi:archive'),
  },
  {
    label: '编辑用户资料',
    key: 'editProfile',
    icon: renderIcon('mdi:archive'),
  },
  {
    label: '退出登录',
    key: 'logout',
    icon: renderIcon('mdi:archive'),
  }
]



</script>


<template>
  <!-- 移动端顶部导航栏 -->
  <Transition name="slide-fade" appear>
    <div v-if="barShow" :class="navClass"
      class="fixed inset-x-0 top-0 z-11 h-[60px] flex items-center justify-between px-4 py-2 lg:hidden">
      <!-- 左上角标题 -->
      <RouterLink to="/" class="text-[18px] font-bold">
        {{ appStore.blog_config.website_author }}
      </RouterLink>
      <!-- 右上角图标 -->
      <div class="flex items-center gap-2 text-2xl">
        <button @click="appStore.setSearchFlag(true)">
          <Icon icon="ic:round-search" />
        </button>
        <button @click="appStore.setCollapsed(true)">
          <Icon icon="ic:sharp-menu" />
        </button>
      </div>
    </div>
  </Transition>
  <!-- 侧边栏 -->
  <!-- <MobileSideBar /> -->
  <!-- 电脑端顶部导航栏 -->
  <Transition name="slide-fade" appear>
    <div v-if="barShow" :class="navClass" class="fixed inset-x-0 top-0 z-11 hidden h-[60px] lg:block">
      <div class="h-full flex items-center justify-between px-9">
        <!-- 左上角标题 -->
        <RouterLink to="/" class="text-xl font-bold">
          哭哭
        </RouterLink>
        <!-- 右上角菜单 -->
        <div class="flex items-center space-x-4">
          <n-dropdown :options="options">
            <n-button>用户资料</n-button>
          </n-dropdown>
          
          
        </div>
      </div>
    </div>
  </Transition>
</template>



<style scoped></style>