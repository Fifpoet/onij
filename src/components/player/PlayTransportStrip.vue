<template>
  <div
    class="play-queue-drawer flex w-full flex-col gap-3"
    :class="surface === 'lyrics' ? 'play-transport--lyrics' : ''"
  >
    <div class="flex w-full items-center gap-2">
      <span
        class="min-w-[3rem] shrink-0 text-right text-xs tabular-nums"
        :class="surface === 'lyrics' ? 'text-white/55' : 'text-gray-600 dark:text-gray-400'"
      >
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
      <span
        class="min-w-[3rem] shrink-0 text-xs tabular-nums"
        :class="surface === 'lyrics' ? 'text-white/55' : 'text-gray-600 dark:text-gray-400'"
      >
        {{ progressDurationLabel }}
      </span>
    </div>

    <div
      class="mx-auto flex w-full max-w-full flex-wrap items-center justify-center gap-x-2 gap-y-2"
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

      <n-tooltip trigger="hover">
        <template #trigger>
          <n-button
            text
            circle
            size="medium"
            :focusable="false"
            :disabled="nextDisabled"
            aria-label="下一曲"
            @click="playQueue.skipToNext()"
          >
            <n-icon :component="PlaySkipForwardOutline" :size="22" />
          </n-button>
        </template>
        下一曲
      </n-tooltip>

      <div class="flex items-center gap-2 pl-2">
        <n-icon
          :component="volumeIcon"
          :size="18"
          class="shrink-0"
          :class="surface === 'lyrics' ? 'text-white/45' : 'text-gray-500 dark:text-gray-400'"
        />
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
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { NButton, NIcon, NTooltip, NSlider } from 'naive-ui'
import {
  PlayCircleOutline,
  PauseCircleOutline,
  ListOutline,
  ShuffleOutline,
  PlaySkipForwardOutline,
  VolumeHighOutline,
  VolumeMediumOutline,
  VolumeLowOutline,
  VolumeMuteOutline,
} from '@vicons/ionicons5'
import { usePlayQueueStore } from '@/store/playQueue'
import { neteasePlayerControl } from '@/player/neteasePlayerControl'
import { formatTime } from '@/util/time'

withDefaults(
  defineProps<{
    /** 歌词全屏等深色底上使用浅色时间/滑轨 */
    surface?: 'default' | 'lyrics'
  }>(),
  { surface: 'default' },
)

const playQueue = usePlayQueueStore()

function formatPlaybackClock(sec: number): string {
  const n = Number(sec)
  if (!Number.isFinite(n) || n < 0) return '00:00'
  return formatTime(Math.floor(n + 1e-6))
}

const progressCurrentLabel = computed(() => formatPlaybackClock(playQueue.playbackCurrentSec))
const progressDurationLabel = computed(() => formatPlaybackClock(playQueue.playbackDurationSec))

const progressPercentRounded = computed(() => {
  const d = playQueue.playbackDurationSec
  if (d <= 0) return 0
  const pct = (playQueue.playbackCurrentSec / d) * 100
  return Math.min(100, Math.max(0, Math.round(pct)))
})

const playButtonDisabled = computed(
  () => !playQueue.nowPlaying && playQueue.queue.length === 0,
)

const nextDisabled = computed(
  () => !playQueue.nowPlaying || playQueue.queue.length === 0,
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
  playQueue.playMode === 'shuffle' ? '随机播放' : '顺序播放',
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
.volume-slider-drawer :deep(.n-slider),
.drawer-progress-slider :deep(.n-slider) {
  --n-handle-size: 10px;
  --n-rail-height: 3px;
}

.volume-slider-drawer :deep(.n-slider-rail__fill),
.drawer-progress-slider :deep(.n-slider-rail__fill) {
  max-width: calc(100% - var(--n-handle-size) / 2);
}

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

.play-transport--lyrics :deep(.n-slider-rail) {
  background-color: rgba(255, 255, 255, 0.18) !important;
}
</style>
