<template>
  <div
    class="player-header-wave flex h-5 w-5 shrink-0 items-end justify-center gap-[2px]"
    :class="active ? 'player-header-wave--active' : 'player-header-wave--idle'"
    aria-hidden="true"
  >
    <span v-for="i in 4" :key="i" class="player-header-wave__bar" :style="{ '--i': i }" />
  </div>
</template>

<script setup lang="ts">
defineProps<{
  active?: boolean
}>()
</script>

<style scoped>
.player-header-wave__bar {
  display: block;
  width: 2px;
  border-radius: 1px;
  background: currentColor;
  height: 4px;
  transform-origin: bottom;
}

.player-header-wave--idle {
  color: rgba(107, 114, 128, 0.45);
}
.dark .player-header-wave--idle {
  color: rgba(156, 163, 175, 0.4);
}

.player-header-wave--idle .player-header-wave__bar {
  height: 5px;
}

.player-header-wave--active {
  color: #3b82f6;
}
.dark .player-header-wave--active {
  color: #60a5fa;
}

.player-header-wave--active .player-header-wave__bar {
  animation: wave-bounce 0.85s ease-in-out infinite;
  animation-delay: calc((var(--i) - 1) * 0.12s);
}

@keyframes wave-bounce {
  0%,
  100% {
    height: 4px;
    opacity: 0.55;
  }
  50% {
    height: 14px;
    opacity: 1;
  }
}
</style>
