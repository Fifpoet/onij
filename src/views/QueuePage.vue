<template>
  <PageContent>
  <div class="queue-page">
    <section class="queue-page__section">
      <div class="queue-page__section-head">
        <h1 class="queue-page__title border-blue-500">正在播放</h1>
        <span v-if="remainingLabel" class="queue-page__remaining">{{ remainingLabel }}</span>
      </div>
      <div class="queue-page__list">
        <SongItem
          v-if="playQueue.nowPlaying"
          :song="playQueue.nowPlaying"
          variant="queuePage"
          :show-queue-add="false"
        />
      </div>
    </section>

    <section class="queue-page__section">
      <h1 class="queue-page__title border-green-500">队列</h1>
      <div class="queue-page__list">
        <SongItem
          v-for="song in playQueue.queue"
          :key="song.id"
          :song="song"
          variant="queuePage"
          :show-queue-add="false"
          :show-queue-pin="true"
          :requeue-previous-on-play-now="false"
        />
      </div>
    </section>
  </div>
  </PageContent>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import PageContent from '@/components/layout/PageContent.vue'
import SongItem from '@/components/music/MusicListItem.vue'
import { usePlayQueueStore } from '@/store/playQueue'

const playQueue = usePlayQueueStore()

function songDurationSec(song: { time_long?: number }): number {
  const t = song.time_long ?? 0
  return Number.isFinite(t) && t > 0 ? t : 0
}

const remainingSec = computed(() => {
  let total = 0
  const current = playQueue.nowPlaying
  if (current) {
    const dur =
      playQueue.playbackDurationSec > 0
        ? playQueue.playbackDurationSec
        : songDurationSec(current)
    const remain = Math.max(0, dur - playQueue.playbackCurrentSec)
    total += remain
  }
  for (const song of playQueue.queue) {
    total += songDurationSec(song)
  }
  return total
})

const remainingLabel = computed(() => {
  if (!playQueue.nowPlaying) return ''
  const sec = Math.floor(remainingSec.value)
  const h = Math.floor(sec / 3600)
  const m = Math.floor((sec % 3600) / 60)
  if (h > 0) return `剩余 ${h} h ${m} m`
  return `剩余 ${m} m`
})
</script>

<style scoped>
.queue-page {
  box-sizing: border-box;
  display: flex;
  width: 100%;
  flex-direction: column;
  align-items: flex-start;
  justify-content: flex-start;
  margin: 0;
  padding: 0;
  text-align: left;
}

.queue-page__section {
  width: 100%;
  align-self: flex-start;
  margin-bottom: 2.5rem;
}

.queue-page__section:last-child {
  margin-bottom: 0;
}

.queue-page__section-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 0.75rem;
  margin-bottom: 1rem;
}

@media (min-width: 768px) {
  .queue-page__section-head {
    margin-bottom: 1.25rem;
  }
}

.queue-page__remaining {
  flex-shrink: 0;
  font-size: 1rem;
  line-height: 1.5rem;
  color: rgb(107 114 128);
}

.dark .queue-page__remaining {
  color: rgb(156 163 175);
}

.queue-page__title {
  margin: 0;
  padding-left: 0.75rem;
  border-left-width: 4px;
  border-left-style: solid;
  font-size: 1.75rem;
  font-weight: 700;
  line-height: 1.25;
  letter-spacing: -0.02em;
  color: rgb(17 24 39);
}

.dark .queue-page__title {
  color: rgb(243 244 246);
}

@media (min-width: 768px) {
  .queue-page__title {
    font-size: 2.125rem;
  }
}

.queue-page__list {
  display: flex;
  width: 100%;
  flex-direction: column;
  align-items: stretch;
  align-self: flex-start;
}
</style>
