<template>
  <div class="lg:hidden">
    <Transition name="slide-fade" appear>
      <header
        class="mobile-top-bar fixed inset-x-0 top-0 z-50 flex items-center justify-between gap-2 border-b border-gray-200/80 bg-white/95 px-3 shadow-sm backdrop-blur-md dark:border-dark-600/80 dark:bg-dark-800/95"
      >
        <div class="flex min-w-0 flex-1 items-center gap-2">
          <RouterLink
            to="/search"
            class="shrink-0 text-[17px] font-bold tracking-tight text-gray-900 dark:text-gray-100"
          >
            ONIJ
          </RouterLink>
          <PlayerHeaderControls />
        </div>
        <div class="flex shrink-0 items-center gap-1">
          <button type="button" class="mobile-header-icon-btn" aria-label="搜索" @click="openSearchOverlay">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="7" /><path d="M20 20L17 17" />
            </svg>
          </button>
          <button type="button" class="mobile-header-icon-btn" aria-label="菜单" @click="toggleMenu">
            <svg v-if="isMenuOpen" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M6 6l12 12M18 6L6 18" />
            </svg>
            <svg v-else width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M4 7h16M4 12h16M4 17h16" />
            </svg>
          </button>
        </div>
      </header>
    </Transition>

    <Teleport to="body">
      <Transition name="search-overlay">
        <div
          v-if="isSearchOverlayOpen"
          class="fixed inset-0 z-[300] flex flex-col bg-white dark:bg-dark-800"
          role="dialog"
          aria-modal="true"
        >
          <div
            class="flex shrink-0 items-center gap-3 border-b border-gray-200 px-3 pb-3 pt-[max(0.75rem,env(safe-area-inset-top))] dark:border-dark-600"
          >
            <input
              ref="searchInputRef"
              v-model="searchValue"
              type="search"
              class="min-w-0 flex-1 rounded-lg border border-gray-200 bg-gray-50 px-3 py-2.5 text-base outline-none focus:border-blue-500 dark:border-dark-600 dark:bg-dark-900 dark:text-gray-100"
              placeholder="搜索歌曲、歌手、专辑"
              @keydown.enter="submitSearch"
            >
            <button type="button" class="shrink-0 px-3 py-2 text-[15px] font-medium text-gray-600" @click="closeSearchOverlay">
              取消
            </button>
          </div>
        </div>
      </Transition>
    </Teleport>

    <Transition name="slide-down">
      <nav
        v-if="isMenuOpen"
        class="mobile-menu-panel fixed inset-x-0 z-40 border-b border-gray-200/90 bg-white/98 px-3 pb-[max(1rem,env(safe-area-inset-bottom))] pt-3 shadow-lg backdrop-blur-md dark:border-dark-600 dark:bg-dark-800/98"
        :style="{ top: menuTop }"
      >
        <p class="mb-2 px-1 text-xs font-semibold uppercase tracking-wide text-gray-400">工具</p>
        <RouterLink to="/tran" class="mobile-menu-link" @click="isMenuOpen = false">中转站</RouterLink>
        <RouterLink to="/prac" class="mobile-menu-link" @click="isMenuOpen = false">练习</RouterLink>
        <RouterLink to="/cloud" class="mobile-menu-link" @click="isMenuOpen = false">云盘</RouterLink>
        <RouterLink to="/monitor" class="mobile-menu-link" @click="isMenuOpen = false">监控</RouterLink>
      </nav>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onUnmounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import PlayerHeaderControls from '@/components/layout/PlayerHeaderControls.vue'
import { useSearch } from '@/composables/searchMusic'

const router = useRouter()
const { searchValue } = useSearch()
const isMenuOpen = ref(false)
const isSearchOverlayOpen = ref(false)
const searchInputRef = ref<HTMLInputElement | null>(null)
const menuTop = 'calc(56px + env(safe-area-inset-top, 0px))'

function openSearchOverlay() {
  isMenuOpen.value = false
  isSearchOverlayOpen.value = true
}
function closeSearchOverlay() {
  isSearchOverlayOpen.value = false
}
async function submitSearch() {
  const q = searchValue.value?.trim()
  if (!q) return
  closeSearchOverlay()
  await router.push({ name: 'Search', query: { q } })
}
function toggleMenu() {
  if (!isSearchOverlayOpen.value) isMenuOpen.value = !isMenuOpen.value
}

watch(isSearchOverlayOpen, (open) => {
  document.body.style.overflow = open ? 'hidden' : ''
  if (open) setTimeout(() => searchInputRef.value?.focus(), 50)
})
onUnmounted(() => {
  document.body.style.overflow = ''
})
</script>

<style scoped>
.mobile-top-bar {
  min-height: calc(56px + env(safe-area-inset-top, 0px));
  padding-top: env(safe-area-inset-top, 0px);
}
.mobile-header-icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.5rem;
  height: 2.5rem;
  border: 0;
  border-radius: 9999px;
  background: transparent;
  color: rgb(55 65 81);
}
.mobile-menu-link {
  display: block;
  padding: 0.85rem 1rem;
  border-radius: 0.625rem;
  font-size: 1rem;
  font-weight: 500;
  color: rgb(17 24 39);
  text-decoration: none;
  background: rgb(249 250 251);
}
.mobile-menu-link.router-link-active {
  background: rgb(219 234 254);
  color: rgb(29 78 216);
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
</style>
