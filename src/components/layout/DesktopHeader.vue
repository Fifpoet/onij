<!-- 桌面端导航栏 -->
<template>
  <Transition name="slide-fade" appear>
    <div class="fixed inset-x-0 top-0 z-50 h-[60px] bg-white/95 dark:bg-dark-800/95 backdrop-blur-sm">
      <div class="h-full flex items-center justify-between px-9">
        <div class="flex items-center gap-3 min-w-0">
          <RouterLink to="/" class="text-xl font-bold shrink-0">
            ONIJ
          </RouterLink>
          <PlayerHeaderControls />
        </div>
        <!-- 右上角菜单 -->
        <div class="flex items-center justify-end space-x-4">
          <div v-for="(items, group) in menuItems" :key="group" class="flex items-center">
            <n-dropdown :options="createDropdownOptions(items)" trigger="hover">
              <n-button>{{ group }}</n-button>
            </n-dropdown>
          </div>
        </div>
        <n-input
            placeholder="搜索"
            :style="{ width: '300px' }"
            @keyup.enter="onSearch"
            v-model:value="searchValue"
        >
          <template #prefix>
            <n-icon :component="SearchOutline" />
          </template>
        </n-input>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { h } from 'vue'
import { useMusicStore } from '@/store'
import { RouterLink } from 'vue-router'
import { NButton, NDropdown, NIcon, NInput } from 'naive-ui'
import { SearchOutline } from '@vicons/ionicons5'
import type { DropdownOption } from 'naive-ui'
import { MidShowWhat } from '@/api/types'
import {
  MailOutline,
  TrendingUpOutline,
  FootstepsOutline,
  FolderOutline,
  RecordingOutline,
  DocumentTextOutline,
  BookOutline,
  DiamondOutline,
  ChatboxOutline
} from '@vicons/ionicons5'
import {useSearch} from "@/composables/searchMusic.ts";
import {router} from "@/router.ts";
import PlayerHeaderControls from '@/components/layout/PlayerHeaderControls.vue'

const musicStore = useMusicStore()

interface MenuItem {
  key: string
  label: string
  icon: any
  onClick?: () => void
}

const { searchValue } = useSearch()

const onSearch = async () => {
  if (!searchValue.value) {
    return
  }
  await router.push({
    name: 'Search',
    query: {
      q: searchValue.value
    }
  })
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
    { 
      key: 'chat', 
      label: 'chat', 
      icon: ChatboxOutline,
      onClick: () => musicStore.setMidShowWhat(MidShowWhat.ShowChatWindow)
    },
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

const createDropdownOptions = (items: MenuItem[]): DropdownOption[] => {
  return items.map(item => ({
    key: item.key,
    label: () => h('div', { class: 'flex items-center' }, [
      h(NIcon, { class: 'mr-2' }, { default: () => h(item.icon) }),
      h('span', { class: 'text-sm' }, item.label)
    ]),
    ...(item.onClick ? { props: { onClick: item.onClick } } : {})
  }))
}
</script>

<style scoped>
.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 0.3s ease-out;
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  opacity: 0;
}
</style> 