<template>
  <div class="w-full flex justify-center">
    <div class="w-full max-w-4xl px-4 py-8">
      <div class="mb-8">
        <h3 class="text-xl font-semibold mb-4 border-l-4 border-blue-500 pl-2">歌手</h3>
        <div class="artists-container">
          <div v-if="searchResults.artists.length === 0" class="text-gray-500 py-4">
            没有找到相关歌手
          </div>
          <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
            <div
              v-for="artist in searchResults.artists"
              :key="artist.id"
              class="flex items-center p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded cursor-pointer"
              @click="goToArtist(artist.id)"
            >
              <img
                :src="artist.picUrl || artist.img1v1Url || defaultAvatar"
                :alt="artist.name"
                class="w-12 h-12 rounded-full object-cover"
              >
              <div class="ml-3">
                <div class="font-medium">{{ artist.name }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="mb-8">
        <h3 class="text-xl font-semibold mb-4 border-l-4 border-green-500 pl-2">专辑</h3>
        <div class="albums-container">
          <div v-if="searchResults.albums.length === 0" class="text-gray-500 py-4">
            没有找到相关专辑
          </div>
          <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
            <div
              v-for="album in searchResults.albums"
              :key="album.id"
              class="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded cursor-pointer"
              @click="goToAlbum(album.id)"
            >
              <img
                :src="album.picUrl || defaultAlbumCover"
                :alt="album.name"
                class="w-full aspect-square object-cover rounded"
              >
              <div class="mt-2">
                <div class="font-medium truncate">{{ album.name }}</div>
                <div class="text-sm text-gray-500 truncate">{{ album.artist.name }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="mb-8">
        <h3 class="text-xl font-semibold mb-4 border-l-4 border-red-500 pl-2">单曲</h3>
        <div class="songs-container">
          <div v-if="searchResults.songDetails.length === 0" class="text-gray-500 py-4">
            没有找到相关单曲
          </div>
          <div v-else class="space-y-2">
            <div
              v-for="detail in searchResults.songDetails"
              :key="detail.id"
              class="flex items-center p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded cursor-pointer"
            >
              <img
                :src="detail.al.picUrl"
                :alt="detail.name"
                class="w-12 h-12 object-cover rounded"
              >
              <div class="ml-3 flex-1 min-w-0">
                <div class="font-medium truncate">{{ detail.name }}</div>
                <div class="text-sm text-gray-500 truncate">
                  {{ detail.ar.map(a => a.name).join('/') }} - {{ detail.al.name }}
                </div>
              </div>
              <div class="text-sm text-gray-500 ml-2">
                {{ formatDuration(detail.dt) }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import {router} from "@/router.ts";
import { onMounted, reactive, watch } from 'vue'
import { useSearch } from '@/composables/searchMusic'
import type { SongSearchItem, AlbumSearchItem, ArtistSearchItem } from '@/api/netease/search'
import { getSongDetail } from '@/api/netease/search'
import type { SongDetail } from '@/api/netease/result'

const route = useRoute()
const { searchValue, handleSearch } = useSearch()

// 默认图片
const defaultAvatar = 'https://p1.music.126.net/VnZiScyynLG7atLIZ2YPkw==/18686200114669622.jpg'
const defaultAlbumCover = 'https://p1.music.126.net/6y-UleORITEDbvrOLV0Q8A==/5639395138885805.jpg'

// 搜索结果
const searchResults = reactive({
  songs: [] as SongSearchItem[],
  albums: [] as AlbumSearchItem[],
  artists: [] as ArtistSearchItem[],
  songDetails: [] as SongDetail[]
})

// 格式化时长
const formatDuration = (duration: number) => {
  const seconds = Math.floor(duration / 1000)
  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = seconds % 60
  return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`
}

// 跳转到专辑详情页
const goToAlbum = (albumId: number) => {
  router.push(`/album?id=${albumId}`)
}

// 跳转到歌手详情页
const goToArtist = (artistId: number) => {
  router.push(`/artist?ids=${artistId}`)
}

// 执行搜索
const performSearch = async (keywords: string) => {
  if (!keywords.trim()) return

  try {
    const results = await handleSearch([1, 10, 100])
    const allSongIds: number[] = []
    
    results?.forEach(res => {
      if (res.result.songs) {
        searchResults.songs = res.result.songs
        // 提取所有歌曲的 id
        const ids = res.result.songs.map(song => song.id)
        allSongIds.push(...ids)
      }
      if (res.result.albums) {
        searchResults.albums = res.result.albums
      }
      if (res.result.artists) {
        searchResults.artists = res.result.artists
      }
    })

    // 如果有歌曲ID，调用详情接口
    if (allSongIds.length > 0) {
      try {
        const detailResponse = await getSongDetail(allSongIds)
        if (detailResponse.code === 200 && detailResponse.songs) {
          searchResults.songDetails = detailResponse.songs
        }
      } catch (error) {
        console.error('获取歌曲详情出错:', error)
      }
    }
  } catch (error) {
    console.error('搜索出错:', error)
  }
}

watch(
    () => route.query.q,
    (newQuery) => {
      if (newQuery) {
        searchValue.value = newQuery as string
        performSearch(newQuery as string)
      }
    }
)

onMounted(() => {
  if (route.query.q) {
    searchValue.value = route.query.q as string
    performSearch(route.query.q as string)
  }
})

// 监听路由变化
// 如果需要在搜索页面内部再次搜索，可以使用这个方法
</script>

<style scoped>
.artists-container,
.albums-container,
.songs-container {
  min-height: 100px;
}
</style>
