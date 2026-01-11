<template>
  <div class="w-full flex justify-center">
    <div class="w-4/5 px-4 py-8">
      <!-- 歌手 & 专辑 并排 -->
      <div class="mb-8">
        <div class="flex flex-col md:flex-row gap-8">
          <!-- 歌手 -->
          <div class="md:w-1/2">
            <h3 v-if="!isLoading" class="text-xl font-semibold mb-4 border-l-4 border-blue-500 pl-2">歌手</h3>
            <div class="artists-container">
              <div v-if="isLoading"></div>
              <div v-else-if="searchResults.artists.length === 0" class="text-gray-500 py-4">
                没有找到相关歌手
              </div>
              <div v-else class="grid grid-cols-3 gap-4">
                <div
                  v-for="artist in searchResults.artists.slice(0, 6)"
                  :key="artist.id"
                  class="flex flex-col items-center p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded cursor-pointer"
                  @click="goToArtist(artist.id)"
                >
                  <img
                    :src="artist.picUrl || artist.img1v1Url || defaultAvatar"
                    :alt="artist.name"
                    class="w-full aspect-square rounded-full object-cover"
                  >
                  <div class="mt-2 text-center">
                    <div class="font-bold text-sm truncate w-full">{{ artist.name }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 专辑 -->
          <div class="md:w-1/2">
            <h3 v-if="!isLoading" class="text-xl font-semibold mb-4 border-l-4 border-green-500 pl-2">专辑</h3>
            <div class="albums-container">
              <div v-if="isLoading"></div>
              <div v-else-if="searchResults.albums.length === 0" class="text-gray-500 py-4">
                没有找到相关专辑
              </div>
              <div v-else class="grid grid-cols-2 sm:grid-cols-3 gap-4">
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
                  <div class="mt-2 text-center">
                    <div class="font-bold truncate">{{ album.name }}</div>
                    <div 
                      class="text-sm text-gray-500 truncate hover:text-blue-500 hover:underline cursor-pointer inline-block"
                      @click.stop="goToArtist(album.artist.id)"
                    >
                      {{ album.artist.name }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 单曲 -->
      <div class="mb-8">
        <h3 v-if="!isLoading" class="text-xl font-semibold mb-4 border-l-4 border-red-500 pl-2">单曲</h3>
        <div class="songs-container">
          <div v-if="isLoading"></div>
          <div v-else-if="searchResults.songDetails.length === 0" class="text-gray-500 py-4">
            没有找到相关单曲
          </div>
          <div v-else class="grid grid-cols-3 gap-4">
            <div
              v-for="song in viewSongs"
              :key="song.id"
              class="w-full"
            >
              <SongItem :song="song" />
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
import { onMounted, reactive, watch, computed, ref } from 'vue'
import { useSearch } from '@/composables/searchMusic'
import type { SongSearchItem, AlbumSearchItem, ArtistSearchItem } from '@/api/netease/search'
import { getSongDetail } from '@/api/netease/search'
import type { SongDetail } from '@/api/netease/result'
import SongItem from '@/components/music/MusicListItem.vue'
import type { ViewMusicListItem } from '@/api/view/music'

const route = useRoute()
const { searchValue, handleSearch } = useSearch()

// 默认图片
const defaultAvatar = 'https://p1.music.126.net/VnZiScyynLG7atLIZ2YPkw==/18686200114669622.jpg'
const defaultAlbumCover = 'https://p1.music.126.net/6y-UleORITEDbvrOLV0Q8A==/5639395138885805.jpg'

// 加载状态
const isLoading = ref(false)

// 搜索结果
const searchResults = reactive({
  songs: [] as SongSearchItem[],
  albums: [] as AlbumSearchItem[],
  artists: [] as ArtistSearchItem[],
  songDetails: [] as SongDetail[]
})

// 将歌曲详情转换为通用列表项结构，供 MusicListItem 使用
const viewSongs = computed<ViewMusicListItem[]>(() => {
  return searchResults.songDetails.map((detail) => {
    return {
      id: detail.id,
      name: detail.name,
      time_long: detail.dt / 1000,
      album_id: detail.al.id,
      album_name: detail.al.name,
      cover_file_url: detail.al.picUrl,
      artists: detail.ar.map(artist => ({
        artist_id: artist.id,
        artist_name: artist.name,
      })),
    }
  })
})

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

  isLoading.value = true
  
  // 清空之前的结果
  searchResults.songs = []
  searchResults.albums = []
  searchResults.artists = []
  searchResults.songDetails = []

  try {
    const results = await handleSearch([1, 10, 100], 0, 20)
    const allSongIds: number[] = []
    
    results?.forEach(res => {
      if (res.result.songs) {
        searchResults.songs = res.result.songs
        // 提取所有歌曲的 id
        const ids = res.result.songs.map(song => song.id)
        allSongIds.push(...ids)
      }
      if (res.result.albums) {
        searchResults.albums = res.result.albums.filter(album => album.artist.albumSize > 10 || album.artist.musicSize > 10)
      }
      if (res.result.artists) {
        searchResults.artists = res.result.artists.filter(artist => artist.albumSize > 10 || artist.musicSize > 10)
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
  } finally {
    isLoading.value = false
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
</script>

<style scoped>
.artists-container,
.albums-container,
.songs-container {
  min-height: 200px;
}
</style>
