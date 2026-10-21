<template>
  <div class="presence-peer" :class="{ 'presence-peer--idle': !song }" data-testid="presence-peer">
    <button
      type="button"
      class="presence-peer__cover"
      :disabled="!song"
      :aria-label="song ? '同步对方正在听的歌' : '未在播放'"
      @click="onCover"
    >
      <img :src="coverUrl" alt="" class="presence-peer__img">
    </button>
    <div class="presence-peer__meta min-w-0">
      <p class="presence-peer__nick">
        <PresenceAvatar :seed="peer.clientId" :title="peer.nickname" size="sm" />
        <span class="presence-peer__nick-text" data-testid="presence-peer-nick">{{ peer.nickname }}</span>
      </p>
      <template v-if="song">
        <p class="presence-peer__artists truncate">
          <template v-for="(artist, i) in artists" :key="artist.artist_id">
            <RouterLink
              :to="{ path: '/artist', query: { ids: String(artist.artist_id) } }"
              class="presence-peer__link"
              @click="emit('close')"
            >
              {{ artist.artist_name }}
            </RouterLink>
            <span v-if="i < artists.length - 1" class="presence-peer__sep">/</span>
          </template>
          <span v-if="!artists.length">未知歌手</span>
        </p>
        <p class="presence-peer__title truncate">
          <RouterLink
            v-if="song.album_id"
            :to="{ path: '/album', query: { id: String(song.album_id) } }"
            class="presence-peer__link presence-peer__link--title"
            @click="emit('close')"
          >
            {{ song.name }}
          </RouterLink>
          <span v-else>{{ song.name }}</span>
        </p>
      </template>
      <p v-else class="presence-peer__idle">未在播放</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { usePresenceStore, type PresencePeer } from '@/store/presence'
import PresenceAvatar from '@/components/presence/PresenceAvatar.vue'

const props = defineProps<{
  peer: PresencePeer
}>()

const emit = defineEmits<{
  close: []
}>()

const presence = usePresenceStore()

const defaultAlbumCover =
  'https://p1.music.126.net/6y-UleORITEDbvrOLV0Q8A==/5639395138885805.jpg'

const song = computed(() => props.peer.nowPlaying)
const artists = computed(() => {
  const s = song.value
  if (!s) return []
  if (s.display_artists?.length) return s.display_artists
  return s.artists ?? []
})
const coverUrl = computed(() => song.value?.cover_file_url || defaultAlbumCover)

function onCover() {
  if (!song.value) return
  presence.adoptPeer(props.peer)
  emit('close')
}
</script>

<style scoped>
.presence-peer {
  display: flex;
  align-items: center;
  gap: 0.65rem;
  padding: 0.5rem 0.25rem;
}
.presence-peer--idle {
  opacity: 0.55;
}
.presence-peer__cover {
  flex-shrink: 0;
  width: 2.5rem;
  height: 2.5rem;
  padding: 0;
  border: 0;
  border-radius: 0.375rem;
  overflow: hidden;
  background: rgb(229 231 235);
  cursor: pointer;
}
.presence-peer__cover:disabled {
  cursor: default;
}
.presence-peer__img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.presence-peer__meta {
  flex: 1;
  min-width: 0;
  text-align: left;
}
.presence-peer__nick {
  display: flex;
  align-items: center;
  gap: 0.3rem;
  margin: 0 0 0.1rem;
  font-size: 0.6875rem;
  line-height: 1.2;
  color: rgb(156 163 175);
}
.presence-peer__nick-text {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.presence-peer__artists,
.presence-peer__title,
.presence-peer__idle {
  margin: 0;
  font-size: 0.8125rem;
  line-height: 1.35;
}
.presence-peer__artists {
  color: rgb(75 85 99);
}
.presence-peer__title {
  font-weight: 600;
  color: rgb(17 24 39);
}
.presence-peer__idle {
  color: rgb(107 114 128);
}
.presence-peer__link {
  color: inherit;
  text-decoration: none;
}
.presence-peer__link:hover {
  text-decoration: underline;
  color: rgb(37 99 235);
}
.presence-peer__sep {
  color: rgb(156 163 175);
}
.dark .presence-peer__artists,
.dark .presence-peer__idle {
  color: rgb(156 163 175);
}
.dark .presence-peer__title {
  color: rgb(243 244 246);
}
.dark .presence-peer__cover {
  background: rgb(55 65 81);
}
</style>
