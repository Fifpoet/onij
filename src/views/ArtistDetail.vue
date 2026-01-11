<template>
  <div class="w-full flex justify-center min-w-0">
    <div class="w-[80vw] px-4 py-8">
      <div v-if="loading" class="bg-white dark:bg-gray-800 rounded-lg  flex items-center justify-center">
        <div class="text-center py-8"></div>
      </div>
      <div v-else-if="artist" class="bg-white dark:bg-gray-800 rounded-lg  overflow-hidden">
        <!-- 歌手信息 -->
        <div class="flex p-6">
          <div class="w-1/3 flex justify-center">
            <img
                :src="artist.picUrl || artist.img1v1Url"
                :alt="artist.name"
                class="w-64 h-64 rounded-full object-cover shadow-md"
            >
          </div>
          <div class="w-2/3 pl-6 flex flex-col justify-center text-left">
            <h1 class="text-5xl font-bold mb-4 ">{{ artist.name }}</h1>
            <div class="grid grid-cols-3 gap-3 text-gray-700 dark:text-gray-300">
              <div>专辑: {{ artist.albumSize }}</div>
              <div>歌曲: {{ artist.musicSize }}</div>
              <div>MV: {{ artist.mvSize }}</div>
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
          <h2 class="text-2xl font-bold mb-4">热门歌曲</h2>
          <div v-if="songDetails.length === 0" class="text-gray-500 py-4">
            暂无歌曲
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
import {ref, onMounted, computed, watch} from 'vue'
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
