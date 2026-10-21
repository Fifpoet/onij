<template>
  <n-popover
    v-model:show="open"
    trigger="click"
    placement="bottom-end"
    :show-arrow="false"
    raw
  >
    <template #trigger>
      <PresenceFriendsBtn />
    </template>
    <div class="presence-desk-panel">
      <p class="presence-desk-panel__title">在线</p>
      <p v-if="!peers.length" class="presence-desk-panel__empty" data-testid="presence-empty">暂无其他在线用户</p>
      <div v-else class="presence-desk-panel__list">
        <PresencePeerRow
          v-for="peer in peers"
          :key="peer.clientId"
          :peer="peer"
          @close="open = false"
        />
      </div>
    </div>
  </n-popover>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { NPopover } from 'naive-ui'
import { usePresenceStore } from '@/store/presence'
import PresencePeerRow from '@/components/presence/PresencePeerRow.vue'
import PresenceFriendsBtn from '@/components/presence/PresenceFriendsBtn.vue'

const open = ref(false)
const presence = usePresenceStore()
const { peers } = storeToRefs(presence)
watch(open, (shown) => {
  if (shown) presence.refreshPeers()
})
</script>

<style scoped>
.presence-desk-panel {
  width: 17.5rem;
  max-height: min(24rem, 70vh);
  overflow: auto;
  padding: 0.65rem 0.75rem 0.75rem;
  border-radius: 0.75rem;
  background: #fff;
  box-shadow: 0 10px 28px rgb(0 0 0 / 0.12);
}
.dark .presence-desk-panel {
  background: rgb(31 41 55);
  box-shadow: 0 10px 28px rgb(0 0 0 / 0.4);
}
.presence-desk-panel__title {
  margin: 0 0 0.35rem;
  font-size: 0.75rem;
  font-weight: 600;
  letter-spacing: 0.04em;
  color: rgb(156 163 175);
}
.presence-desk-panel__empty {
  margin: 0.5rem 0 0.25rem;
  font-size: 0.8125rem;
  color: rgb(107 114 128);
}
.presence-desk-panel__list {
  display: flex;
  flex-direction: column;
}
</style>
