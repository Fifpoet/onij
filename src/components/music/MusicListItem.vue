<template>
  <div class="flex items-center py-3 border-b border-gray-100 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700 rounded px-2">
    <img
        :src="song.album?.name || defaultAlbumCover"
        :alt="song.name"
        class="w-12 h-12 object-cover rounded mr-4"
        @error="handleImageError"
    >
    <div class="flex-1 min-w-0">
      <div class="font-medium truncate">{{ song.name }}</div>
      <div class="text-sm text-gray-500 truncate">
        {{ getArtists(song) }} - {{ song.album?.name }}
      </div>
    </div>
    <div class="text-sm text-gray-500 ml-2">
      {{ formatDuration(song.duration || 0) }}
    </div>
  </div>
</template>

<script setup lang="ts">
import type { SongSearchItem } from '@/api/netease/search'

// 默认专辑封面
const defaultAlbumCover = 'https://p1.music.126.net/6y-UleORITEDbvrOLV0Q8A==/5639395138885805.jpg'

const props = defineProps<{
  song: SongSearchItem
}>()

// 获取艺术家名称
const getArtists = (song: SongSearchItem): string => {
  if (!song.artists || song.artists.length === 0) {
    return '未知艺术家'
  }
  return song.artists.map(a => a.name).join('/')
}

// 格式化时长
const formatDuration = (duration: number) => {
  const seconds = Math.floor(duration / 1000)
  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = seconds % 60
  return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`
}

// 处理图片加载错误
const handleImageError = (event: Event) => {
  const imgElement = event.target as HTMLImageElement
  if (imgElement.src !== defaultAlbumCover) {
    imgElement.src = defaultAlbumCover
  }
}
</script>
