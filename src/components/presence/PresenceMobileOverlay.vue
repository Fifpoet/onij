<template>
  <Teleport to="body">
    <Transition name="presence-overlay">
      <div
        v-if="open"
        class="fixed inset-0 z-[300] flex flex-col bg-white dark:bg-dark-800"
        role="dialog"
        aria-modal="true"
        aria-label="在线好友"
      >
        <div
          class="flex shrink-0 items-center gap-3 border-b border-gray-200 px-3 pb-3 pt-[max(0.75rem,env(safe-area-inset-top))] dark:border-dark-600"
        >
          <p class="min-w-0 flex-1 text-[17px] font-semibold text-gray-900 dark:text-gray-100">
            在线
          </p>
          <button
            type="button"
            class="shrink-0 px-3 py-2 text-[15px] font-medium text-gray-600 dark:text-gray-300"
            @click="emit('close')"
          >
            关闭
          </button>
        </div>
        <div class="min-h-0 flex-1 overflow-y-auto px-3 py-2">
          <p v-if="!peers.length" class="px-1 py-8 text-center text-[15px] text-gray-400">
            暂无其他在线用户
          </p>
          <PresencePeerRow
            v-else
            v-for="peer in peers"
            :key="peer.clientId"
            :peer="peer"
            @close="emit('close')"
          />
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { watch } from 'vue'
import { usePresenceStore } from '@/store/presence'
import PresencePeerRow from '@/components/presence/PresencePeerRow.vue'

const props = defineProps<{
  open: boolean
}>()

const emit = defineEmits<{
  close: []
}>()

const presence = usePresenceStore()
const { peers } = storeToRefs(presence)
watch(
  () => props.open,
  (shown) => {
    if (shown) presence.refreshPeers()
  },
)
</script>

<style scoped>
.presence-overlay-enter-active,
.presence-overlay-leave-active {
  transition: opacity 0.2s ease;
}
.presence-overlay-enter-from,
.presence-overlay-leave-to {
  opacity: 0;
}
</style>
