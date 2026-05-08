<!-- 移动端导航栏：单根包裹，否则父级 lg:hidden 无法 fallthrough -->
<template>
  <div class="lg:hidden">
    <Transition name="slide-fade" appear>
      <header
        class="mobile-top-bar fixed inset-x-0 top-0 z-50 flex items-center justify-between gap-3 border-b border-gray-200/80 bg-white/95 px-3 shadow-sm backdrop-blur-md dark:border-dark-600/80 dark:bg-dark-800/95"
      >
        <div class="flex min-w-0 flex-1 items-center gap-2">
          <RouterLink
            to="/"
            class="shrink-0 text-[17px] font-bold tracking-tight text-gray-900 dark:text-gray-100"
          >
            ONIJ
          </RouterLink>
          <PlayerHeaderControls />
        </div>
        <div class="flex shrink-0 items-center gap-1">
          <button
            type="button"
            class="flex h-10 w-10 items-center justify-center rounded-full text-gray-700 active:bg-gray-100 dark:text-gray-200 dark:active:bg-dark-700"
            aria-label="搜索"
            @click="openSearchOverlay"
          >
            <n-icon size="22">
              <SearchOutline />
            </n-icon>
          </button>
          <button
            type="button"
            class="flex h-10 w-10 items-center justify-center rounded-full text-gray-700 active:bg-gray-100 dark:text-gray-200 dark:active:bg-dark-700"
            aria-label="菜单"
            @click="toggleMenu"
          >
            <n-icon size="22">
              <component :is="isMenuOpen ? CloseOutline : MenuOutline" />
            </n-icon>
          </button>
        </div>
      </header>
    </Transition>

    <!-- 全屏搜索：覆盖下方所有内容 -->
    <Teleport to="body">
      <Transition name="search-overlay">
        <div
          v-if="isSearchOverlayOpen"
          class="fixed inset-0 z-[300] flex flex-col bg-white dark:bg-dark-800"
          role="dialog"
          aria-modal="true"
          aria-label="搜索"
        >
          <div
            class="flex shrink-0 items-center gap-3 border-b border-gray-200 px-3 pb-3 pt-[max(0.75rem,env(safe-area-inset-top))] dark:border-dark-600"
          >
            <n-input
              ref="searchInputRef"
              v-model:value="searchValue"
              class="search-overlay-input min-w-0 flex-1"
              size="large"
              placeholder="搜索歌曲、歌手、专辑"
              clearable
              @keyup.enter="submitSearch"
            >
              <template #prefix>
                <n-icon :component="SearchOutline" class="opacity-60" />
              </template>
            </n-input>
            <button
              type="button"
              class="shrink-0 rounded-lg px-3 py-2 text-[15px] font-medium text-gray-600 active:bg-gray-100 dark:text-gray-300 dark:active:bg-dark-700"
              @click="closeSearchOverlay"
            >
              取消
            </button>
          </div>
          <div class="flex flex-1 flex-col bg-gray-50 dark:bg-dark-900" />
        </div>
      </Transition>
    </Teleport>

    <!-- 导航菜单 -->
    <Transition name="slide-down">
      <div
        v-if="isMenuOpen"
        class="mobile-menu-panel fixed inset-x-0 z-40 overflow-y-auto border-b border-gray-200/90 bg-white/98 shadow-lg backdrop-blur-md dark:border-dark-600 dark:bg-dark-800/98"
        :style="{
          top: menuTop,
          maxHeight: 'min(72vh, calc(100dvh - 56px - env(safe-area-inset-top, 0px)))',
        }"
      >
        <div class="px-3 pb-[max(1rem,env(safe-area-inset-bottom))] pt-3">
          <template v-for="(items, group) in menuItems" :key="group">
            <div class="mb-5 last:mb-0">
              <h3 class="mb-2 px-1 text-xs font-semibold uppercase tracking-wide text-gray-400 dark:text-gray-500">
                {{ group }}
              </h3>
              <div class="grid grid-cols-3 gap-2 sm:grid-cols-4">
                <button
                  v-for="item in items"
                  :key="item.key"
                  type="button"
                  class="flex flex-col items-center justify-center rounded-xl bg-gray-50 py-3 active:scale-[0.98] dark:bg-dark-700/80"
                  @click="handleMenuClick(item)"
                >
                  <n-icon size="26" class="mb-1 text-gray-700 dark:text-gray-200">
                    <component :is="item.icon" />
                  </n-icon>
                  <span class="max-w-full truncate px-1 text-center text-xs text-gray-600 dark:text-gray-400">
                    {{ item.label }}
                  </span>
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
import { ref, watch, nextTick, onUnmounted } from 'vue'
import { useMusicStore } from '@/store'
import { RouterLink, useRouter } from 'vue-router'
import PlayerHeaderControls from '@/components/layout/PlayerHeaderControls.vue'
import { NIcon, NInput } from 'naive-ui'
import type { InputInst } from 'naive-ui'
import { useSearch } from '@/composables/searchMusic'
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
  DiamondOutline,
} from '@vicons/ionicons5'

const musicStore = useMusicStore()
const router = useRouter()
const { searchValue } = useSearch()

const isMenuOpen = ref(false)
const isSearchOverlayOpen = ref(false)
const searchInputRef = ref<InputInst | null>(null)

/** 与顶栏实际高度一致，避免菜单与 bar 重叠（含安全区） */
const menuTop = 'calc(56px + env(safe-area-inset-top, 0px))'

const openSearchOverlay = () => {
  isMenuOpen.value = false
  isSearchOverlayOpen.value = true
}

const closeSearchOverlay = () => {
  isSearchOverlayOpen.value = false
}

const submitSearch = async () => {
  const q = searchValue.value?.trim()
  if (!q) return
  closeSearchOverlay()
  await router.push({ name: 'Search', query: { q } })
}

watch(isSearchOverlayOpen, async (open) => {
  document.body.style.overflow = open ? 'hidden' : ''
  if (open) {
    await nextTick()
    searchInputRef.value?.focus()
  }
})

onUnmounted(() => {
  document.body.style.overflow = ''
})

interface MenuItem {
  key: string
  label: string
  icon: unknown
  onClick?: () => void
}

const menuItems: Record<string, MenuItem[]> = {
  mbox: [
    { key: 'relay', label: 'relay', icon: MailOutline },
    { key: 'trending', label: 'trending', icon: TrendingUpOutline },
    { key: 'footprint', label: 'footprint', icon: FootstepsOutline },
  ],
  goroutine: [
    {
      key: 'file',
      label: 'file',
      icon: FolderOutline,
      onClick: () => musicStore.setMidShowWhat(MidShowWhat.ShowFileList),
    },
    { key: 'record', label: 'record', icon: RecordingOutline },
  ],
  recall: [
    {
      key: 'memo',
      label: 'memo',
      icon: DocumentTextOutline,
      onClick: () => musicStore.setMidShowWhat(MidShowWhat.ShowMemoList),
    },
    { key: 'poem', label: 'poem', icon: BookOutline },
    { key: 'treasure', label: 'treasure', icon: DiamondOutline },
  ],
}

const toggleMenu = () => {
  if (isSearchOverlayOpen.value) return
  isMenuOpen.value = !isMenuOpen.value
}

const handleMenuClick = (item: MenuItem) => {
  item.onClick?.()
  isMenuOpen.value = false
}
</script>

<style scoped>
.mobile-top-bar {
  min-height: calc(56px + env(safe-area-inset-top, 0px));
  padding-top: env(safe-area-inset-top, 0px);
}

.slide-down-enter-active,
.slide-down-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.slide-down-enter-from,
.slide-down-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: opacity 0.25s ease;
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  opacity: 0;
}

.search-overlay-enter-active,
.search-overlay-leave-active {
  transition: opacity 0.2s ease;
}

.search-overlay-enter-from,
.search-overlay-leave-to {
  opacity: 0;
}

.search-overlay-input {
  max-width: 100%;
}

:deep(.search-overlay-input.n-input) {
  min-width: 0;
}
</style>
