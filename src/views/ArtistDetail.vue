<template>
  <div class="w-full flex justify-center">
    <div class="w-1/2 px-4 py-8">
      <div v-if="loading" class="text-center py-8">
        加载中...
      </div>
      <div v-else-if="error" class="text-center py-8 text-red-500">
        {{ error }}
      </div>
      <div v-else-if="artist" class="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden">
        <!-- 歌手信息 -->
        <div class="flex p-6">
          <div class="w-1/3 flex justify-center">
            <img
                :src="artist.picUrl || artist.img1v1Url"
                :alt="artist.name"
                class="w-48 h-48 rounded-full object-cover shadow-md"
                @error="handleImageError"
            >
          </div>
          <div class="w-2/3 pl-6 flex flex-col justify-center">
            <h1 class="text-3xl font-bold mb-4">{{ artist.name }}</h1>
            <div class="grid grid-cols-2 gap-3 text-gray-700 dark:text-gray-300">
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
          <SongItem
              v-for="song in hotSongs"
              :key="song.id"
              :song="song"
              class="cursor-pointer"
              @click="playSong(song)"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {ref, onMounted, computed, watch} from 'vue'
import { useRoute } from 'vue-router'
import { getNetease } from '@/util'
import SongItem from '@/components/music/MusicListItem.vue'
import {ArtistDetailResponse} from "@/api/netease/result.ts";
import {ArtistDetail, ArtistHotSong} from "@/api/netease/artists.ts";

// 默认头像
const defaultAvatar = 'https://p1.music.126.net/VnZiScyynLG7atLIZ2YPkw==/18686200114669622.jpg'
const defaultAlbumCover = 'https://p1.music.126.net/6y-UleORITEDbvrOLV0Q8A==/5639395138885805.jpg'

const route = useRoute()
const artist = ref<ArtistDetail>()
const hotSongs = ref<ArtistHotSong[]>([])
const loading = ref(true)
const error = ref('')
const showFullDescription = ref(false)
const fullDescription = ref('')

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

// 处理图片加载错误
const handleImageError = (event: Event) => {
  const imgElement = event.target as HTMLImageElement
  // 根据元素的类名判断是头像还是专辑封面
  if (imgElement.classList.contains('rounded-full')) {
    // 歌手头像
    if (imgElement.src !== defaultAvatar) {
      imgElement.src = defaultAvatar
    }
  } else {
    // 专辑封面
    if (imgElement.src !== defaultAlbumCover) {
      imgElement.src = defaultAlbumCover
    }
  }
}

// 播放歌曲（后续可以连接到播放功能）
const playSong = (song: ArtistHotSong) => {
  console.log('播放歌曲:', song)
  // 这里可以连接到音乐播放功能
}

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
