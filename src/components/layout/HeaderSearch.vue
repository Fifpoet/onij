<template>
  <form class="header-search" @submit.prevent="onSubmit">
    <span class="header-search__icon" aria-hidden="true">
      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="11" cy="11" r="7" />
        <path d="M20 20L17 17" />
      </svg>
    </span>
    <input v-model="searchValue" type="search" class="header-search__input" placeholder="搜索" />
  </form>
</template>

<script setup lang="ts">
import { useSearch } from '@/composables/searchMusic'
import { router } from '@/router'

const { searchValue } = useSearch()

async function onSubmit() {
  const q = searchValue.value?.trim()
  if (!q) return
  await router.push({ name: 'Search', query: { q } })
}
</script>

<style scoped>
.header-search {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  width: 100%;
  max-width: 18.75rem;
  height: 2.25rem;
  padding: 0 0.75rem;
  border: 1px solid rgb(229 231 235);
  border-radius: 0.5rem;
  background: rgb(249 250 251);
}
.dark .header-search {
  border-color: rgb(75 85 99);
  background: rgb(17 24 39);
}
.header-search__icon {
  color: rgb(156 163 175);
}
.header-search__input {
  flex: 1;
  min-width: 0;
  border: 0;
  background: transparent;
  font-size: 0.9375rem;
  outline: none;
}
.dark .header-search__input {
  color: rgb(243 244 246);
}
</style>
