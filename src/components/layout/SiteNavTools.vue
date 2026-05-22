<template>
  <div class="site-nav-tools" @mouseenter="open = true" @mouseleave="open = false">
    <button
      type="button"
      class="site-nav-tools__trigger"
      :class="{ 'site-nav-tools__trigger--active': isTranRoute }"
      @click="onTriggerClick"
    >
      工具
      <span class="site-nav-tools__caret" aria-hidden="true" />
    </button>
    <Transition name="site-nav-drop">
      <div v-show="open" class="site-nav-tools__panel">
        <RouterLink to="/tran" class="site-nav-tools__item" @click="open = false">
          中转站
        </RouterLink>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const open = ref(false)
const isTranRoute = computed(() => route.path === '/tran')

function onTriggerClick() {
  if (typeof window !== 'undefined' && window.matchMedia('(hover: none)').matches) {
    void router.push('/tran')
  } else {
    open.value = !open.value
  }
}
</script>

<style scoped>
.site-nav-tools {
  position: relative;
}
.site-nav-tools__trigger {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  height: 2.25rem;
  padding: 0 0.85rem;
  border: 0;
  border-radius: 0.5rem;
  background: transparent;
  font-size: 0.9375rem;
  font-weight: 500;
  color: rgb(55 65 81);
  cursor: pointer;
}
.dark .site-nav-tools__trigger {
  color: rgb(229 231 235);
}
.site-nav-tools__trigger:hover,
.site-nav-tools__trigger--active {
  background: rgb(243 244 246);
}
.dark .site-nav-tools__trigger:hover,
.dark .site-nav-tools__trigger--active {
  background: rgb(55 65 81);
}
.site-nav-tools__caret {
  width: 0;
  height: 0;
  border-left: 4px solid transparent;
  border-right: 4px solid transparent;
  border-top: 5px solid currentColor;
  opacity: 0.65;
}
.site-nav-tools__panel {
  position: absolute;
  top: calc(100% + 6px);
  left: 0;
  z-index: 60;
  min-width: 9rem;
  padding: 0.35rem;
  border: 1px solid rgb(229 231 235);
  border-radius: 0.5rem;
  background: #fff;
  box-shadow: 0 10px 24px rgb(0 0 0 / 0.08);
}
.dark .site-nav-tools__panel {
  border-color: rgb(75 85 99);
  background: rgb(31 41 55);
}
.site-nav-tools__item {
  display: block;
  padding: 0.55rem 0.75rem;
  border-radius: 0.375rem;
  font-size: 0.9375rem;
  color: rgb(55 65 81);
  text-decoration: none;
}
.dark .site-nav-tools__item {
  color: rgb(229 231 235);
}
.site-nav-tools__item:hover,
.site-nav-tools__item.router-link-active {
  background: rgb(243 244 246);
}
.site-nav-drop-enter-active,
.site-nav-drop-leave-active {
  transition: opacity 0.14s ease, transform 0.14s ease;
}
.site-nav-drop-enter-from,
.site-nav-drop-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
