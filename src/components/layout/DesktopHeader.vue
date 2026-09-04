<template>
  <Transition name="slide-fade" appear>
    <header class="desktop-header fixed inset-x-0 top-0 z-50 h-[60px] border-b border-gray-200/80 bg-white/95 backdrop-blur-sm dark:border-dark-600/80 dark:bg-dark-800/95">
      <div class="desktop-header__inner">
        <div class="desktop-header__left min-w-0">
          <RouterLink to="/search" class="desktop-header__logo shrink-0">ONIJ</RouterLink>
          <PlayerHeaderControls />
        </div>
        <div class="desktop-header__right">
          <RouterLink
            to="/ai"
            class="desktop-header__ai"
            :class="{ 'desktop-header__ai--active': isAi }"
          >
            AI
          </RouterLink>
          <SiteNavTools />
          <HeaderSearch />
        </div>
      </div>
    </header>
  </Transition>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import PlayerHeaderControls from '@/components/layout/PlayerHeaderControls.vue'
import SiteNavTools from '@/components/layout/SiteNavTools.vue'
import HeaderSearch from '@/components/layout/HeaderSearch.vue'

const route = useRoute()
const isAi = computed(() => route.path === '/ai')
</script>

<style scoped>
.desktop-header__inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1.5rem;
  height: 100%;
  padding: 0 2.25rem;
}
.desktop-header__left {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex: 1;
  min-width: 0;
}
.desktop-header__logo {
  font-size: 1.25rem;
  font-weight: 700;
  color: rgb(17 24 39);
  text-decoration: none;
}
.dark .desktop-header__logo {
  color: rgb(243 244 246);
}
.desktop-header__right {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex-shrink: 0;
}
.desktop-header__ai {
  display: inline-flex;
  flex-shrink: 0;
  align-items: center;
  height: 2.25rem;
  padding: 0 0.85rem;
  white-space: nowrap;
  border-radius: 0.5rem;
  font-size: 0.9375rem;
  font-weight: 500;
  color: rgb(55 65 81);
  text-decoration: none;
}
.dark .desktop-header__ai {
  color: rgb(229 231 235);
}
.desktop-header__ai:hover,
.desktop-header__ai--active {
  background: rgb(243 244 246);
}
.dark .desktop-header__ai:hover,
.dark .desktop-header__ai--active {
  background: rgb(55 65 81);
}
.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: opacity 0.3s ease;
}
.slide-fade-enter-from,
.slide-fade-leave-to {
  opacity: 0;
}
</style>
