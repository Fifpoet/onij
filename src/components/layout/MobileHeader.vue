<!-- 移动端导航栏：单根包裹，否则父级 lg:hidden 无法 fallthrough，大屏会与桌面导航重复一行 -->
<template>
  <div class="lg:hidden">
  <Transition name="slide-fade" appear>
    <div class="fixed inset-x-0 top-0 z-50 h-[60px] flex items-center justify-between px-4 py-2 bg-white/95 dark:bg-dark-800/95 backdrop-blur-sm">
      <div class="flex items-center gap-2 min-w-0">
        <RouterLink to="/" class="text-[18px] font-bold shrink-0">
          ONIJ
        </RouterLink>
        <PlayerHeaderDisc />
      </div>
      <!-- 右上角图标 -->
      <div class="flex items-center gap-3">
        <!-- 搜索按钮 -->
        <button @click=""
          class="w-8 h-8 flex items-center justify-center rounded-full hover:bg-gray-100 dark:hover:bg-dark-700">
          <n-icon size="20">
            <SearchOutline />
          </n-icon>
        </button>
        <!-- 菜单按钮 -->
        <button @click="toggleMenu" 
          class="w-8 h-8 flex items-center justify-center rounded-full hover:bg-gray-100 dark:hover:bg-dark-700">
          <n-icon size="20">
            <component :is="isMenuOpen ? CloseOutline : MenuOutline" />
          </n-icon>
        </button>
      </div>
    </div>
  </Transition>

  <!-- 下拉菜单 -->
  <Transition name="slide-down">
    <div v-if="isMenuOpen" class="fixed inset-x-0 top-[60px] z-40 bg-white/95 dark:bg-dark-800/95 backdrop-blur-sm">
      <div class="px-4 py-3">
        <template v-for="(items, group) in menuItems" :key="group">
          <div class="mb-4">
            <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400 mb-2">{{ group }}</h3>
            <div class="grid grid-cols-3 gap-2">
              <button v-for="item in items" :key="item.key"
                class="flex flex-col items-center justify-center p-3 rounded-lg bg-gray-50 dark:bg-dark-700 hover:bg-gray-100 dark:hover:bg-dark-600"
                @click="handleMenuClick(item)">
                <n-icon size="24" class="mb-1">
                  <component :is="item.icon" />
                </n-icon>
                <span class="text-sm">{{ item.label }}</span>
              </button>
            </div>
          </div>
        </template>
      </div>
    </div>
  </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMusicStore } from '@/store'
import { RouterLink } from 'vue-router'
import PlayerHeaderDisc from '@/components/layout/PlayerHeaderDisc.vue'
import { NIcon } from 'naive-ui'
import { MidShowWhat } from '@/api/types'
import {
  SearchOutline,
  MenuOutline,
  CloseOutline,
  MailOutline,
  TrendingUpOutline,
  FootstepsOutline,
  FolderOutline,
  RecordingOutline,
  DocumentTextOutline,
  BookOutline,
  DiamondOutline
} from '@vicons/ionicons5'

const musicStore = useMusicStore()
const isMenuOpen = ref(false)

interface MenuItem {
  key: string
  label: string
  icon: any
  onClick?: () => void
}

const menuItems: Record<string, MenuItem[]> = {
  "mbox": [
    { key: 'relay', label: 'relay', icon: MailOutline },
    { key: 'trending', label: 'trending', icon: TrendingUpOutline },
    { key: 'footprint', label: 'footprint', icon: FootstepsOutline },
  ],
  "goroutine": [
    { 
      key: 'file', 
      label: 'file', 
      icon: FolderOutline,
      onClick: () => musicStore.setMidShowWhat(MidShowWhat.ShowFileList)
    },
    { key: 'record', label: 'record', icon: RecordingOutline },
  ],
  "recall": [
    { 
      key: 'memo', 
      label: 'memo', 
      icon: DocumentTextOutline,
      onClick: () => musicStore.setMidShowWhat(MidShowWhat.ShowMemoList)
    },
    { key: 'poem', label: 'poem', icon: BookOutline },
    { key: 'treasure', label: 'treasure', icon: DiamondOutline },
  ]
}

const toggleMenu = () => {
  isMenuOpen.value = !isMenuOpen.value
}

const handleMenuClick = (item: MenuItem) => {
  if (item.onClick) {
    item.onClick()
  }
  isMenuOpen.value = false
}
</script>

<style scoped>
.slide-down-enter-active,
.slide-down-leave-active {
  transition: all 0.3s ease-out;
}

.slide-down-enter-from {
  transform: translateY(-20px);
  opacity: 0;
}

.slide-down-leave-to {
  transform: translateY(-20px);
  opacity: 0;
}

.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 0.3s ease-out;
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  opacity: 0;
}
</style> 