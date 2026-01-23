<template>
  <div class="w-full flex justify-center min-w-0">
    <div class="w-[80vw] px-4 py-8">
      <div v-if="loading" class="bg-white dark:bg-gray-800 rounded-lg flex items-center justify-center">
        <div class="text-center py-8"></div>
      </div>
      <div v-else-if="album" class="bg-white dark:bg-gray-800 rounded-lg overflow-hidden">
        <!-- 专辑信息 -->
        <div class="flex p-6">
          <div class="w-1/3 flex justify-center">
            <img
                :src="album.picUrl || album.blurPicUrl"
                :alt="album.name"
                class="w-64 h-64 rounded-lg object-cover shadow-md"
            >
          </div>
          <div class="w-2/3 pl-6 flex flex-col justify-center text-left">
            <h1 class="text-5xl font-bold mb-4">{{ album.name }}</h1>
            <div class="grid grid-cols-3 gap-3 text-gray-700 dark:text-gray-300 mb-4">
              <div>艺术家: {{ album.artist.name }}</div>
              <div>发行时间: {{ formatDate(album.publishTime) }}</div>
              <div>歌曲数: {{ album.size }}</div>
            </div>

            <!-- 简介 -->
            <div class="mt-4">
              <div
                  class="text-gray-600 dark:text-gray-400 cursor-pointer"
                  @click="showFullDescription = !showFullDescription"
              >
                {{ displayedDescription }}
                <span v-if="shouldShowMore" class="text-blue-500">
                  {{ showFullDescription ? '收起' : '...更多' }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- 歌曲列表 -->
        <div class="px-6 py-4 border-t border-gray-200 dark:border-gray-700">
          <div v-if="songDetails.length === 0" class="text-gray-500 py-4">
            暂无歌曲
          </div>
          <div v-else>
            <!-- 按 cd 分组显示 -->
            <template v-for="(discSongs, discNumber) in groupedSongs" :key="discNumber">
              <div v-if="shouldShowDiscTitle(discNumber)" class="mb-4">
                <h3 class="text-lg font-semibold text-gray-600 dark:text-gray-400 mb-2">
                  Disc {{ discNumber }}
                </h3>
              </div>
              <div class="mb-6">
                <SongItem
                  v-for="(song, idx) in discSongs"
                  :key="song.id"
                  :song="song"
                  :index="idx + 1"
                  :show-cover="false"
                  :show-album-name="false"
                />
              </div>
            </template>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {ref, onMounted, computed} from 'vue'
import { useRoute } from 'vue-router'
import { getSongDetail } from '@/api/netease/search'
import type { SongDetail } from '@/api/netease/result'
import SongItem from '@/components/music/MusicListItem.vue'
import { getAlbumDetail } from '@/api/netease/album'
import type { AlbumDetail, AlbumSong } from '@/api/netease/album'
import {ViewMusicListItem} from "@/api/view/music.ts";
import { formatDate } from '@/util/time'

// 默认封面
const defaultAlbumCover = 'https://p1.music.126.net/6y-UleORITEDbvrOLV0Q8A==/5639395138885805.jpg'

const route = useRoute()
const album = ref<AlbumDetail>()
const songs = ref<AlbumSong[]>([])
const songDetails = ref<SongDetail[]>([])
const loading = ref(true)
const error = ref('')
const showFullDescription = ref(false)

// 获取专辑创作者ID列表
const albumArtistIds = computed(() => {
  if (!album.value) return new Set<number>()
  const ids = new Set<number>()
  // 添加主艺术家
  if (album.value.artist?.id) {
    ids.add(album.value.artist.id)
  }
  // 添加所有艺术家
  if (album.value.artists) {
    album.value.artists.forEach(artist => {
      if (artist.id) {
        ids.add(artist.id)
      }
    })
  }
  return ids
})

// 按 cd 字段分组歌曲
const groupedSongs = computed(() => {
  const groups: Record<string, ViewMusicListItem[]> = {}
  const artistIds = albumArtistIds.value
  
  songDetails.value.forEach((detail) => {
    const discNumber = detail.cd || '0'
    if (!groups[discNumber]) {
      groups[discNumber] = []
    }
    
    // 过滤掉专辑创作者的歌手
    const filteredArtists = detail.ar
      .filter(artist => !artistIds.has(artist.id))
      .map(artist => ({
        artist_id: artist.id,
        artist_name: artist.name,
      }))
    
    groups[discNumber].push({
      id: detail.id,
      name: detail.name,
      time_long: detail.dt / 1000,
      album_id: detail.al.id,
      album_name: detail.al.name,
      cover_file_url: detail.al.picUrl || defaultAlbumCover,
      artists: filteredArtists,
    })
  })
  
  // 按 disc 序号排序
  const sortedGroups: Record<string, ViewMusicListItem[]> = {}
  Object.keys(groups).sort((a, b) => {
    const numA = parseInt(a) || 0
    const numB = parseInt(b) || 0
    return numA - numB
  }).forEach(key => {
    sortedGroups[key] = groups[key]
  })
  
  return sortedGroups
})

// 判断是否应该显示 Disc 标题
const shouldShowDiscTitle = (discNumber: string) => {
  // 如果有多个 disc，显示所有 disc 的标题
  // 如果只有一个 disc 且不是 '0' 或空字符串，也显示标题
  const discKeys = Object.keys(groupedSongs.value)
  if (discKeys.length > 1) {
    return true // 多个 disc 时，都显示标题
  }
  // 只有一个 disc 时，只有当它不是 '0' 或空字符串时才显示
  return discNumber !== '0' && discNumber !== '' && discNumber !== '1'
}

// 显示的简介（根据是否展开显示全部或部分）
const displayedDescription = computed(() => {
  const desc = album.value?.description || ''
  if (showFullDescription.value) {
    return desc
  }
  return desc.length > 100 ? desc.substring(0, 100) : desc
})

// 是否应该显示"更多"按钮
const shouldShowMore = computed(() => {
  return (album.value?.description || '').length > 100
})

// 获取专辑详情
const fetchAlbumDetail = async (id: string) => {
  try {
    loading.value = true
    error.value = ''

    // 获取专辑信息
    const response = await getAlbumDetail(id)

    if (response.code === 200) {
      album.value = response.album
      songs.value = response.songs || []

      // 提取所有歌曲ID，调用详情接口获取完整信息（包括专辑封面）
      if (songs.value.length > 0) {
        const songIds = songs.value.map(song => song.id)
        try {
          const detailResponse = await getSongDetail(songIds)
          if (detailResponse.code === 200 && detailResponse.songs) {
            songDetails.value = detailResponse.songs
          }
        } catch (err) {
          console.error('获取歌曲详情出错:', err)
        }
      }
    } else {
      throw new Error('获取专辑详情失败')
    }
  } catch (err) {
    error.value = '获取专辑详情失败'
    console.error('获取专辑详情失败:', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  const albumId = route.query.id as string
  if (albumId) {
    fetchAlbumDetail(albumId)
  }
})
</script>

