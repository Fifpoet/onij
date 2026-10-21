<template>
  <button
    type="button"
    class="presence-friends-btn"
    data-testid="presence-friends-btn"
    :aria-label="ariaLabel"
    title="在线好友"
    @click="emit('click')"
  >
    <n-icon :component="PeopleOutline" :size="22" />
    <span
      v-if="count > 0"
      class="presence-friends-badge"
      data-testid="presence-badge"
    >{{ count > 99 ? '99+' : count }}</span>
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { NIcon } from 'naive-ui'
import { PeopleOutline } from '@vicons/ionicons5'
import { usePresenceStore } from '@/store/presence'

const emit = defineEmits<{
  click: []
}>()

const { peers } = storeToRefs(usePresenceStore())
const count = computed(() => peers.value.length)
const ariaLabel = computed(() =>
  count.value > 0 ? `在线好友，${count.value} 人在线` : '在线好友',
)
</script>

<style scoped>
.presence-friends-btn {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.25rem;
  height: 2.25rem;
  padding: 0;
  border: 0;
  border-radius: 0.5rem;
  background: transparent;
  color: rgb(55 65 81);
  cursor: pointer;
}
.dark .presence-friends-btn {
  color: rgb(229 231 235);
}
.presence-friends-btn:hover {
  background: rgb(243 244 246);
}
.dark .presence-friends-btn:hover {
  background: rgb(55 65 81);
}
.presence-friends-badge {
  position: absolute;
  top: 0.05rem;
  right: 0.05rem;
  min-width: 0.95rem;
  height: 0.95rem;
  padding: 0 0.2rem;
  border-radius: 999px;
  background: #ef4444;
  color: #fff;
  font-size: 0.625rem;
  font-weight: 700;
  line-height: 0.95rem;
  text-align: center;
  pointer-events: none;
}
</style>
