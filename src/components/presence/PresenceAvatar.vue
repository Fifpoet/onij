<template>
  <span
    class="presence-avatar"
    :class="{ 'presence-avatar--sm': size === 'sm' }"
    :style="{ background: visual.background }"
    :title="title"
    aria-hidden="true"
  >{{ visual.emoji }}</span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { presenceAvatar } from '@/util/presenceAvatar'

const props = withDefaults(
  defineProps<{
    seed: string
    title?: string
    size?: 'sm' | 'md'
  }>(),
  { size: 'md' },
)

const visual = computed(() => presenceAvatar(props.seed))
</script>

<style scoped>
.presence-avatar {
  display: inline-flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 1.75rem;
  height: 1.75rem;
  border-radius: 999px;
  font-size: 0.85rem;
  line-height: 1;
  user-select: none;
}
.presence-avatar--sm {
  width: 1.35rem;
  height: 1.35rem;
  font-size: 0.7rem;
}
</style>
