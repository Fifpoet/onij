<template>
  <div
    class="flex items-center hover:bg-gray-50 dark:hover:bg-gray-700 rounded px-2"
    :class="[
      rowPaddingClass,
      showRowBorder ? 'border-b border-gray-100 dark:border-gray-700' : '',
      showQueueAdd ? 'group' : '',
    ]"
    @dblclick="onRowDblClick"
  >
    <!-- 序号 -->
    <div
      v-if="index"
      class="w-10 text-center text-gray-500 mr-2 shrink-0"
      :class="isLargeRow ? 'text-base' : 'text-sm'"
    >
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
      <!-- 专辑页：歌名 + 非专辑歌手（超链接）同一行 -->
      <div
        v-if="variant === 'album'"
        class="flex flex-wrap items-baseline gap-x-1.5 gap-y-1 min-w-0"
      >
        <span class="font-bold text-lg text-gray-900 dark:text-gray-100 shrink min-w-0">{{ song.name }}</span>
        <template v-if="song.artists && song.artists.length > 0">
          <span class="text-gray-400 dark:text-gray-500 shrink-0 text-base">·</span>
          <span class="inline-flex flex-wrap items-baseline gap-x-0 min-w-0 text-base">
            <template v-for="(artist, i) in song.artists" :key="artist.artist_id">
              <RouterLink
                :to="{ path: '/artist', query: { ids: String(artist.artist_id) } }"
                class="text-blue-600 dark:text-blue-400 hover:underline shrink-0"
                @click.stop
              >
                {{ artist.artist_name }}
              </RouterLink>
              <span
                v-if="i < song.artists.length - 1"
                class="text-gray-500 dark:text-gray-400 mx-0.5"
              >/</span>
            </template>
          </span>
        </template>
      </div>

      <!-- 播放队列抽屉：大号字 + 双行（参考专辑列表体量） -->
      <template v-else-if="variant === 'queue'">
        <div class="font-bold text-lg text-gray-900 dark:text-gray-100 truncate">{{ song.name }}</div>
        <div v-if="showSubtitleRow" class="text-base text-gray-500 dark:text-gray-400 truncate mt-0.5">
          <span
              v-for="(artist, i) in lineArtists"
              :key="artist.artist_id"
              class="music-list-skip-dblplay hover:text-blue-500 hover:underline cursor-pointer"
              @click.stop="onArtistClick(artist.artist_id)"
          >
            {{ artist.artist_name }}{{ i < lineArtists.length - 1 ? '/' : '' }}
          </span>
          <template v-if="showAlbumName && song.album_name">
            <span v-if="lineArtists.length > 0" class="mx-1">-</span>
            <span
              class="music-list-skip-dblplay truncate hover:text-blue-500 hover:underline cursor-pointer"
              @click.stop="onAlbumClick(song.album_id)"
            >
              {{ song.album_name }}
            </span>
          </template>
        </div>
      </template>

      <template v-else>
        <!-- 歌曲名称 -->
        <div class="font-bold truncate">{{ song.name }}</div>

        <!-- 歌手 / 专辑（无歌手时仍可只显示专辑名） -->
        <div v-if="showSubtitleRow" class="text-sm text-gray-500 truncate">
          <span
              v-for="(artist, i) in lineArtists"
              :key="artist.artist_id"
              class="music-list-skip-dblplay hover:text-blue-500 hover:underline cursor-pointer"
              @click.stop="onArtistClick(artist.artist_id)"
          >
            {{ artist.artist_name }}{{ i < lineArtists.length - 1 ? '/' : '' }}
          </span>
          <template v-if="showAlbumName && song.album_name">
            <span v-if="lineArtists.length > 0" class="mx-1">-</span>
            <span 
              class="music-list-skip-dblplay truncate hover:text-blue-500 hover:underline cursor-pointer"
              @click.stop="onAlbumClick(song.album_id)"
            >
              {{ song.album_name }}
            </span>
          </template>
        </div>
      </template>
    </div>

    <div class="flex items-center shrink-0 gap-1 ml-2">
      <div
        v-if="showQueueAdd"
        class="w-8 h-8 flex items-center justify-center"
      >
        <n-button
          quaternary
          circle
          size="small"
          class="!opacity-0 group-hover:!opacity-100 transition-opacity duration-150"
          :focusable="false"
          @click.stop="onAddToQueue"
        >
          <template #icon>
            <n-icon :component="AddOutline" :size="18" />
          </template>
        </n-button>
      </div>
      <!-- 时长 -->
      <div
        class="text-gray-500 tabular-nums min-w-[2.5rem] text-right"
        :class="isLargeRow ? 'text-base' : 'text-sm'"
      >
        {{ formatDuration(song.time_long) }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { useMessage, NButton, NIcon } from 'naive-ui'
import { AddOutline } from '@vicons/ionicons5'
import type { ViewMusicListItem } from '@/api/view/music.ts'
import { usePlayQueueStore } from '@/store/playQueue'

// 默认专辑封面
const defaultAlbumCover = 'https://p1.music.126.net/6y-UleORITEDbvrOLV0Q8A==/5639395138885805.jpg'

const router = useRouter()
const message = useMessage()
const playQueue = usePlayQueueStore()

const props = withDefaults(defineProps<{
  song: ViewMusicListItem
  index?: number
  showCover?: boolean
  showAlbumName?: boolean
  /** 专辑曲目列表：歌名后展示非专辑歌手链接 */
  variant?: 'default' | 'album' | 'queue'
  /** 悬停时在时长前显示加入队列 */
  showQueueAdd?: boolean
  /** 底部分割线 */
  showRowBorder?: boolean
  /** 双击空白处（非链接/按钮）插队立即播放 */
  dblClickPlayNow?: boolean
}>(), {
  showCover: true,
  showAlbumName: true,
  variant: 'default',
  showQueueAdd: true,
  showRowBorder: true,
  dblClickPlayNow: true,
})

const isLargeRow = computed(() => props.variant === 'album' || props.variant === 'queue')

const rowPaddingClass = computed(() =>
  isLargeRow.value ? 'py-3.5' : 'py-3',
)

/** 副标题行歌手：专辑页入队后队列用 display_artists（全员），专辑行 variant=album 仍只用 artists（合作歌手） */
const lineArtists = computed(() => {
  const d = props.song.display_artists
  if (d && d.length > 0) return d
  return props.song.artists ?? []
})

const showSubtitleRow = computed(
  () =>
    lineArtists.value.length > 0 || (props.showAlbumName && !!props.song.album_name),
)

function onAddToQueue() {
  const r = playQueue.enqueue(props.song)
  if (r === 'duplicate') {
    message.warning('该歌曲已在队列或正在播放')
  } else {
    message.success('已加入播放队列')
  }
}

function onRowDblClick(e: MouseEvent) {
  if (!props.dblClickPlayNow) return
  const el = e.target as HTMLElement | null
  if (!el) return
  if (el.closest('a, button, [role="button"], .n-button, .music-list-skip-dblplay, img')) return
  playQueue.playNow(props.song)
}

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

// 专辑点击事件 - 跳转到专辑详情页
const onAlbumClick = (id: number) => {
  router.push(`/album?id=${id}`)
}
</script>
