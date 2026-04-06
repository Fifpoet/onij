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

        <!-- 2 进度条：左当前时间 / 右总时长；tooltip=false 关掉拖拽时的「79」百分比浮层（show-tooltip 关不掉） -->
        <div class="flex items-center gap-2 w-full">
          <span class="text-xs tabular-nums text-gray-600 dark:text-gray-400 min-w-[3rem] shrink-0 text-right">
            {{ progressCurrentLabel }}
          </span>
          <div class="drawer-progress-slider flex-1 min-w-0">
            <n-slider
              class="w-full"
              :value="progressPercentRounded"
              :min="0"
              :max="100"
              :step="1"
              :disabled="progressDisabled"
              :tooltip="false"
              :show-tooltip="false"
              aria-label="播放进度"
              :on-dragstart="onProgressDragStart"
              :on-dragend="onProgressDragEnd"
              @update:value="onProgressChange"
            />
          </div>
          <span class="text-xs tabular-nums text-gray-600 dark:text-gray-400 min-w-[3rem] shrink-0">
            {{ progressDurationLabel }}
          </span>
        </div>

        <!-- 3 随机 | 播放 | 音量：居中紧凑一组 -->
        <div
          class="mx-auto flex items-center gap-3 rounded-full bg-gray-100/90 dark:bg-gray-800/80 px-3 py-2 border border-gray-200/80 dark:border-gray-700/80"
        >
          <n-tooltip trigger="hover">
            <template #trigger>
              <n-button
                text
                circle
                size="medium"
                :focusable="false"
                :type="playQueue.playMode === 'shuffle' ? 'primary' : 'default'"
                :aria-label="modeButtonLabel"
                @click="playQueue.togglePlayMode()"
              >
                <n-icon
                  :component="playQueue.playMode === 'shuffle' ? ShuffleOutline : ListOutline"
                  :size="22"
                />
              </n-button>
            </template>
            {{ modeButtonLabel }}
          </n-tooltip>

          <n-tooltip trigger="hover">
            <template #trigger>
              <n-button
                text
                type="primary"
                circle
                size="large"
                :focusable="false"
                :disabled="playButtonDisabled"
                :aria-label="playButtonLabel"
                @click="onPlayButton"
              >
                <n-icon :component="playButtonIcon" :size="30" />
              </n-button>
            </template>
            {{ playButtonLabel }}
          </n-tooltip>

          <div class="flex items-center gap-1.5 pl-1 border-l border-gray-200 dark:border-gray-600">
            <n-icon :component="volumeIcon" :size="18" class="text-gray-500 dark:text-gray-400 shrink-0" />
            <div class="volume-slider-drawer w-[4.5rem] shrink-0">
              <n-slider
                :value="playQueue.volume"
                :step="1"
                :min="0"
                :max="100"
                class="w-full"
                :tooltip="false"
                :show-tooltip="false"
                aria-label="音量"
                @update:value="playQueue.setVolume"
              />
            </div>
          </div>
        </div>

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
import { computed } from 'vue'
import {
  NDrawer,
  NDrawerContent,
  NButton,
  NIcon,
  NTooltip,
  NSlider,
} from 'naive-ui'
import {
  CloseOutline,
  PlayCircleOutline,
  PauseCircleOutline,
  ListOutline,
  ShuffleOutline,
  TrashOutline,
  VolumeHighOutline,
  VolumeMediumOutline,
  VolumeLowOutline,
  VolumeMuteOutline,
} from '@vicons/ionicons5'
import SongItem from '@/components/music/MusicListItem.vue'
import { usePlayQueueStore } from '@/store/playQueue'
import { neteasePlayerControl } from '@/player/neteasePlayerControl'
import { formatTime } from '@/util/time'

const playQueue = usePlayQueueStore()

function formatPlaybackClock(sec: number): string {
  const n = Number(sec)
  if (!Number.isFinite(n) || n < 0) return '00:00'
  return formatTime(Math.floor(n + 1e-6))
}

const progressCurrentLabel = computed(() => formatPlaybackClock(playQueue.playbackCurrentSec))
const progressDurationLabel = computed(() => formatPlaybackClock(playQueue.playbackDurationSec))

/** 进度条绑定用整数百分比，避免 Naive 内部浮点展示异常 */
const progressPercentRounded = computed(() => {
  const d = playQueue.playbackDurationSec
  if (d <= 0) return 0
  const pct = (playQueue.playbackCurrentSec / d) * 100
  return Math.min(100, Math.max(0, Math.round(pct)))
})

const playButtonDisabled = computed(
  () => !playQueue.nowPlaying && playQueue.queue.length === 0,
)

const playButtonLabel = computed(() => {
  if (!playQueue.nowPlaying && playQueue.queue.length > 0) return '播放队列'
  if (playQueue.nowPlaying && playQueue.isPlaying) return '暂停'
  if (playQueue.nowPlaying) return '播放'
  return '播放'
})

const playButtonIcon = computed(() => {
  if (playQueue.nowPlaying && playQueue.isPlaying) return PauseCircleOutline
  return PlayCircleOutline
})

const modeButtonLabel = computed(() =>
  playQueue.playMode === 'shuffle' ? '随机播放（点击切换为顺序）' : '顺序播放（点击切换为随机）',
)

const volumeIcon = computed(() => {
  const v = playQueue.volume
  if (v <= 0) return VolumeMuteOutline
  if (v < 34) return VolumeLowOutline
  if (v < 67) return VolumeMediumOutline
  return VolumeHighOutline
})

const progressDisabled = computed(
  () => !playQueue.nowPlaying || playQueue.playbackDurationSec <= 0,
)

function onPlayButton() {
  if (!playQueue.nowPlaying && playQueue.queue.length > 0) {
    playQueue.startIfIdle()
    return
  }
  neteasePlayerControl.togglePlay()
}

function onProgressDragStart() {
  playQueue.setPlaybackSeeking(true)
}

function onProgressDragEnd() {
  playQueue.setPlaybackSeeking(false)
}

function onProgressChange(v: number | number[]) {
  const pct = Math.round(Array.isArray(v) ? v[0] : v)
  neteasePlayerControl.seekToPercent(pct)
  const d = playQueue.playbackDurationSec
  if (d > 0) {
    playQueue.setPlaybackProgress((pct / 100) * d, d)
  }
}
</script>

<style scoped>
/*
 * 只改 CSS 变量，勿强制 wrapper 尺寸，否则与 Naive 内联百分比错位，最右端「坨」与绿条对不齐。
 * 绿条在 100% 时会铺满整轨，手柄中心实际在轨端内收 handle/2，用 max-width 让填充右缘与手柄中心对齐。
 */
.volume-slider-drawer :deep(.n-slider),
.drawer-progress-slider :deep(.n-slider) {
  --n-handle-size: 10px;
  --n-rail-height: 3px;
}

.volume-slider-drawer :deep(.n-slider-rail__fill),
.drawer-progress-slider :deep(.n-slider-rail__fill) {
  max-width: calc(100% - var(--n-handle-size) / 2);
}

/* 去掉点击/聚焦时 Naive 按钮描边与浏览器 outline */
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
