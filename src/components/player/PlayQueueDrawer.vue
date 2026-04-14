<template>
  <n-drawer
    :show="playQueue.drawerVisible"
    :width="420"
    placement="left"
    @update:show="playQueue.setDrawerVisible"
  >
    <n-drawer-content title="播放队列">
      <div class="play-queue-drawer flex flex-col gap-4 pb-4">
        <!-- 1 正在播放：与下方列表相同的 SongItem -->
        <div
          v-if="playQueue.nowPlaying"
          class="flex items-stretch border-b border-gray-100 dark:border-gray-700 rounded min-w-0 -mx-1"
        >
          <div class="flex-1 min-w-0">
            <SongItem
              :song="playQueue.nowPlaying"
              :show-queue-add="false"
              :show-row-border="false"
            />
          </div>
        </div>
        <div
          v-else
          class="flex items-center gap-3 py-3 px-2 text-sm text-gray-500 dark:text-gray-400 border-b border-gray-100 dark:border-gray-700 rounded"
        >
          暂无播放
        </div>

        <PlayTransportStrip />

        <!-- 4 待播放(n) + 清空 -->
        <div class="flex items-center justify-between gap-2 pt-1 border-t border-gray-100 dark:border-gray-700">
          <span class="text-sm font-medium text-gray-700 dark:text-gray-200">
            待播放（{{ playQueue.queue.length }}）
          </span>
          <n-tooltip trigger="hover">
            <template #trigger>
              <n-button
                text
                circle
                size="small"
                :focusable="false"
                :disabled="!playQueue.queue.length"
                aria-label="清空待播"
                @click="playQueue.clearQueue()"
              >
                <n-icon :component="TrashOutline" :size="18" />
              </n-button>
            </template>
            清空待播
          </n-tooltip>
        </div>

        <!-- 5 待播列表 -->
        <div v-if="playQueue.queue.length === 0" class="text-gray-500 text-center text-sm py-6">
        </div>
        <div v-else class="flex flex-col min-w-0 -mx-1">
          <div
            v-for="song in playQueue.queue"
            :key="song.id"
            class="flex items-stretch border-b border-gray-100 dark:border-gray-700 last:border-b-0 group/qrow hover:bg-gray-50 dark:hover:bg-gray-800/40 rounded"
          >
            <div class="flex-1 min-w-0">
              <SongItem
                :song="song"
                :show-queue-add="false"
                :show-row-border="false"
              />
            </div>
            <div class="flex items-center pr-1 shrink-0">
              <n-button
                text
                circle
                size="tiny"
                :focusable="false"
                class="opacity-0 pointer-events-none group-hover/qrow:opacity-100 group-hover/qrow:pointer-events-auto"
                aria-label="从队列移除"
                @click="playQueue.removeFromQueue(song.id)"
              >
                <template #icon>
                  <n-icon :component="CloseOutline" />
                </template>
              </n-button>
            </div>
          </div>
        </div>
      </div>
    </n-drawer-content>
  </n-drawer>
</template>

<script setup lang="ts">
import { NDrawer, NDrawerContent, NButton, NIcon, NTooltip } from 'naive-ui'
import { CloseOutline, TrashOutline } from '@vicons/ionicons5'
import SongItem from '@/components/music/MusicListItem.vue'
import PlayTransportStrip from '@/components/player/PlayTransportStrip.vue'
import { usePlayQueueStore } from '@/store/playQueue'

const playQueue = usePlayQueueStore()
</script>

<style scoped>
/* 列表内按钮与 PlayTransportStrip 一致：去描边 */
.play-queue-drawer :deep(.n-button),
.play-queue-drawer :deep(.n-button:focus),
.play-queue-drawer :deep(.n-button:focus-visible),
.play-queue-drawer :deep(.n-button:active) {
  outline: none !important;
  box-shadow: none !important;
}
.play-queue-drawer :deep(.n-button .n-button__border),
.play-queue-drawer :deep(.n-button .n-button__state-border) {
  display: none !important;
}
</style>
