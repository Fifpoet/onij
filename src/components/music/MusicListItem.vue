<template>
  <div class="flex items-center py-3 border-b border-gray-100 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700 rounded px-2">
    <!-- 序号 -->
    <div v-if="index" class="w-8 text-center text-gray-500 mr-2">
      {{ index }}
    </div>

    <!-- 专辑封面 -->
    <img
        v-if="showCover"
        :src="song.cover_file_url || defaultAlbumCover"
        :alt="song.name"
        class="w-12 h-12 object-cover rounded mr-4"
    >

    <div class="flex-1 min-w-0">
      <!-- 歌曲名称 -->
      <div class="font-bold truncate">{{ song.name }}</div>

      <!-- 歌手名（可点击） -->
      <div v-if="song.artists && song.artists.length > 0" class="text-sm text-gray-500 truncate">
        <span
            v-for="(artist, i) in song.artists"
            :key="artist.artist_id"
            class="hover:text-blue-500 hover:underline cursor-pointer"
            @click.stop="onArtistClick(artist.artist_id)"
        >
          {{ artist.artist_name }}{{ i < song.artists.length - 1 ? '/' : '' }}
        </span>
        <template v-if="showAlbumName">
          <span class="mx-1">-</span>
          <span class="truncate">{{ song.album_name }}</span>
        </template>
      </div>
    </div>

    <!-- 时长 -->
    <div class="text-sm text-gray-500 ml-2">
      {{ formatDuration(song.time_long) }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import type { ViewMusicListItem } from '@/api/view/music.ts'

// 默认专辑封面
const defaultAlbumCover = 'https://p1.music.126.net/6y-UleORITEDbvrOLV0Q8A==/5639395138885805.jpg'

const router = useRouter()

const props = withDefaults(defineProps<{
  song: ViewMusicListItem
  index?: number
  showCover?: boolean
  showAlbumName?: boolean
}>(), {
  showCover: true,
  showAlbumName: true,
})

// 格式化时长 (从秒转换为 x:xx 格式)
const formatDuration = (seconds: number) => {
  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = Math.floor(seconds % 60)
  return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`
}

// 歌手点击事件 - 跳转到歌手详情页
const onArtistClick = (id: number) => {
  router.push(`/artist?ids=${id}`)
}
</script>
