<template>
  <PageContent>
    <div class="py-4 md:py-8">
      <div v-if="loading" class="bg-white dark:bg-gray-800 rounded-lg flex items-center justify-center">
        <div class="text-center py-8"></div>
      </div>
      <div v-else-if="album" class="bg-white dark:bg-gray-800 rounded-lg overflow-hidden min-w-0">
        <!-- 专辑信息：手机竖排，避免封面与文字重叠 -->
        <div class="flex flex-col gap-4 p-4 sm:flex-row sm:items-start sm:gap-6 sm:p-6">
          <div class="flex shrink-0 justify-center sm:justify-start">
            <img
                :src="album.picUrl || album.blurPicUrl"
                :alt="album.name"
                class="h-36 w-36 rounded-lg object-cover shadow-md sm:h-48 sm:w-48 md:h-64 md:w-64"
            >
          </div>
          <div class="min-w-0 flex-1 text-left select-none sm:pl-0">
            <h1 class="mb-2 text-2xl font-bold leading-tight text-gray-900 dark:text-gray-100 sm:mb-3 sm:text-3xl md:text-5xl">{{ album.name }}</h1>
            <p class="text-base sm:text-lg text-gray-900 dark:text-gray-100 mb-2">
              <span class="font-normal">Album by </span>
              <RouterLink
                :to="{ path: '/artist', query: { ids: String(album.artist.id) } }"
                class="font-bold !text-gray-900 dark:!text-gray-100 hover:!text-gray-900 dark:hover:!text-gray-100 hover:underline"
              >
                {{ album.artist.name }}
              </RouterLink>
            </p>
            <p class="text-base md:text-lg text-gray-500 dark:text-gray-400 tabular-nums mb-4">
              {{ albumMetaYear }}
              <span class="mx-1 text-gray-300 dark:text-gray-600">·</span>
              {{ albumSongCount }} 首歌, {{ albumTotalMinutes }} 分钟
            </p>

            <!-- 简介 -->
            <div class="mt-4">
              <div
                  class="text-base md:text-lg leading-relaxed text-gray-600 dark:text-gray-400 cursor-pointer"
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
        <div class="border-t border-gray-200 px-4 py-4 dark:border-gray-700 sm:px-6">
          <div v-if="songDetails.length === 0" class="text-base text-gray-500 py-4 md:text-lg">
            暂无歌曲
          </div>
          <div v-else>
            <!-- 按 cd 分组显示 -->
            <template v-for="(discSongs, discNumber) in groupedSongs" :key="discNumber">
              <div v-if="shouldShowDiscTitle(discNumber)" class="mb-4">
                <h3 class="text-xl md:text-2xl font-semibold text-gray-600 dark:text-gray-400 mb-2">
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
                  variant="album"
                />
              </div>
            </template>
          </div>
        </div>
      </div>
    </div>
  </PageContent>
</template>

<script setup lang="ts">
import {ref, onMounted, computed} from 'vue'
import PageContent from '@/components/layout/PageContent.vue'
import { useRoute, RouterLink } from 'vue-router'
import { getSongDetail } from '@/api/netease/search'
import type { SongDetail } from '@/api/netease/result'
import SongItem from '@/components/music/MusicListItem.vue'
import { getAlbumDetail } from '@/api/netease/album'
import type { AlbumDetail, AlbumSong } from '@/api/netease/album'
import {ViewMusicListItem} from "@/api/view/music.ts";

// 默认封面
const defaultAlbumCover = 'https://p1.music.126.net/6y-UleORITEDbvrOLV0Q8A==/5639395138885805.jpg'

const route = useRoute()
const album = ref<AlbumDetail>()
const songs = ref<AlbumSong[]>([])
const songDetails = ref<SongDetail[]>([])
const loading = ref(true)
const error = ref('')
const showFullDescription = ref(false)

const albumMetaYear = computed(() => {
  const t = album.value?.publishTime
  if (!t) return '—'
  const y = new Date(t).getFullYear()
  return Number.isFinite(y) ? String(y) : '—'
})

const albumSongCount = computed(() => {
  if (!album.value) return 0
  return songDetails.value.length > 0 ? songDetails.value.length : album.value.size
})

/** 曲目总时长（分），优先用 songDetail，否则用专辑接口返回的 songs.dt（毫秒） */
const albumTotalMinutes = computed(() => {
  const list = songDetails.value.length > 0 ? songDetails.value : songs.value
  if (!list.length) return 0
  const sec = list.reduce((acc, d) => acc + (Number(d.dt) || 0) / 1000, 0)
  return Math.max(0, Math.round(sec / 60))
})

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
      display_artists: detail.ar.map((a) => ({
        artist_id: a.id,
        artist_name: a.name,
      })),
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

