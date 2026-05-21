<template>
  <PageContent>
    <div class="py-4 md:py-8">
      <div v-if="loading" class="bg-white dark:bg-gray-800 rounded-lg  flex items-center justify-center">
        <div class="text-center py-8"></div>
      </div>
      <div v-else-if="artist" class="bg-white dark:bg-gray-800 rounded-lg overflow-hidden min-w-0">
        <!-- 歌手信息：手机竖排，避免与背景圆/文字重叠 -->
        <div class="flex flex-col gap-4 p-4 sm:flex-row sm:items-start sm:gap-6 sm:p-6">
          <div class="flex shrink-0 justify-center sm:justify-start">
            <img
                :src="artist.picUrl || artist.img1v1Url"
                :alt="artist.name"
                class="h-36 w-36 rounded-full object-cover shadow-md sm:h-48 sm:w-48 md:h-64 md:w-64"
            >
          </div>
          <div class="min-w-0 flex-1 text-left">
            <h1 class="mb-3 text-2xl font-bold leading-tight text-gray-900 dark:text-gray-100 sm:text-3xl md:mb-4 md:text-5xl">{{ artist.name }}</h1>
            <p class="text-base md:text-lg text-gray-600 dark:text-gray-400 tabular-nums select-none">
              {{ artist.musicSize }} 首歌
              <span class="mx-1 text-gray-300 dark:text-gray-600">·</span>
              {{ artist.albumSize }} 张专辑
              <span class="mx-1 text-gray-300 dark:text-gray-600">·</span>
              {{ artist.mvSize }} 个 MV
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

        <!-- 歌曲列表：手机单列，宽屏再网格 -->
        <div class="border-t border-gray-200 px-4 py-4 dark:border-gray-700 sm:px-6">
          <h2 class="browse-section-title mb-4">热门歌曲</h2>
          <div v-if="songDetails.length === 0" class="text-base text-gray-500 py-4 md:text-lg">
            暂无歌曲
          </div>
          <div v-else class="flex min-w-0 flex-col md:grid md:grid-cols-2 md:gap-3 lg:grid-cols-3 lg:gap-4">
            <div
              v-for="song in viewSongs"
              :key="song.id"
              class="min-w-0 w-full"
            >
              <SongItem :song="song" variant="browse" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </PageContent>
</template>

<script setup lang="ts">
import {ref, onMounted, computed, watch} from 'vue'
import PageContent from '@/components/layout/PageContent.vue'
import { useRoute } from 'vue-router'
import { getNetease } from '@/util'
import { getSongDetail } from '@/api/netease/search'
import type { SongDetail } from '@/api/netease/result'
import SongItem from '@/components/music/MusicListItem.vue'
import {ArtistDetailResponse} from "@/api/netease/result.ts";
import {ArtistDetail, ArtistHotSong} from "@/api/netease/artists.ts";
import {ViewMusicListItem} from "@/api/view/music.ts";

// 默认封面
const defaultAlbumCover = 'https://p1.music.126.net/6y-UleORITEDbvrOLV0Q8A==/5639395138885805.jpg'

const route = useRoute()
const artist = ref<ArtistDetail>()
const hotSongs = ref<ArtistHotSong[]>([])
const songDetails = ref<SongDetail[]>([])
const albumCoverMap = ref<Map<number, string>>(new Map())
const loading = ref(true)
const error = ref('')
const showFullDescription = ref(false)
const fullDescription = ref('')

// 将歌曲详情转换为通用列表项结构
const viewSongs = computed<ViewMusicListItem[]>(() => {
  return songDetails.value.map((detail) => {
    return {
      id: detail.id,
      name: detail.name,
      time_long: detail.dt / 1000,
      album_id: detail.al.id,
      album_name: detail.al.name,
      cover_file_url: detail.al.picUrl || defaultAlbumCover,
      artists: detail.ar.map(artist => ({
        artist_id: artist.id,
        artist_name: artist.name,
      })),
    }
  })
})

// 显示的简介（根据是否展开显示全部或部分）
const displayedDescription = computed(() => {
  if (showFullDescription.value) {
    return fullDescription.value
  }
  return fullDescription.value.length > 100 ? fullDescription.value.substring(0, 100) : fullDescription.value
})

// 是否应该显示"更多"按钮
const shouldShowMore = computed(() => {
  return fullDescription.value.length > 100
})

// 获取歌手详情
const fetchArtistDetail = async (id: string) => {
  try {
    loading.value = true
    error.value = ''

    // 获取歌手信息
    const timestamp = Date.now()
    const response = await getNetease<ArtistDetailResponse>('/artists', {
      id: id,
      timestamp: timestamp
    })

    if (response.code === 200) {
      artist.value = response.artist
      hotSongs.value = response.hotSongs || []
      fullDescription.value = response.artist.briefDesc || ''

      // 提取所有歌曲ID，调用详情接口获取完整信息（包括专辑封面）
      if (hotSongs.value.length > 0) {
        const songIds = hotSongs.value.map(song => song.id)
        try {
          const detailResponse = await getSongDetail(songIds)
          if (detailResponse.code === 200 && detailResponse.songs) {
            songDetails.value = detailResponse.songs
            // 构建专辑ID到封面URL的映射
            detailResponse.songs.forEach(song => {
              if (song.al.picUrl) {
                albumCoverMap.value.set(song.al.id, song.al.picUrl)
              }
            })
          }
        } catch (err) {
          console.error('获取歌曲详情出错:', err)
        }
      }
    } else {
      throw new Error('获取歌手详情失败')
    }
  } catch (err) {
    error.value = '获取歌手详情失败'
    console.error('获取歌手详情失败:', err)
  } finally {
    loading.value = false
  }
}

watch(
    () => route.query.q,
    (newQuery) => {
      if (newQuery) {
        // fetchArtistDetail(newQuery)
      }
    }
)

onMounted(() => {
  const artistId = route.query.ids as string
  if (artistId) {
    fetchArtistDetail(artistId)
  }
})
</script>
