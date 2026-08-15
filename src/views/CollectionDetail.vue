<template>
  <PageContent>
    <div class="py-4 md:py-8">
      <div v-if="loading" class="bg-white dark:bg-gray-800 rounded-lg flex items-center justify-center">
        <div class="text-center py-8" />
      </div>
      <div v-else-if="collection" class="bg-white dark:bg-gray-800 rounded-lg overflow-hidden min-w-0">
        <div class="flex flex-col gap-4 p-4 sm:flex-row sm:items-start sm:gap-6 sm:p-6">
          <div class="flex shrink-0 justify-center sm:justify-start relative">
            <img
              :src="displayCover"
              :alt="collection.name"
              class="h-36 w-36 rounded-lg object-cover shadow-md sm:h-48 sm:w-48 md:h-64 md:w-64"
            >
            <label
              v-if="editing"
              class="absolute inset-0 flex cursor-pointer items-center justify-center rounded-lg bg-black/40 text-sm text-white opacity-0 transition-opacity hover:opacity-100"
            >
              更换封面
              <input
                type="file"
                accept="image/*"
                class="sr-only"
                @change="onCoverPicked"
              >
            </label>
          </div>

          <div class="min-w-0 flex-1 text-left select-none">
            <input
              v-if="editing"
              v-model="nameDraft"
              class="mb-3 w-full max-w-xl rounded border border-gray-300 bg-transparent px-2 py-1 text-2xl font-bold text-gray-900 outline-none focus:border-blue-500 dark:border-gray-600 dark:text-gray-100 sm:text-3xl md:text-5xl"
              @blur="commitName"
              @keydown.enter.prevent="commitName"
            >
            <h1
              v-else
              class="mb-2 text-2xl font-bold leading-tight text-gray-900 dark:text-gray-100 sm:mb-3 sm:text-3xl md:text-5xl"
            >
              {{ collection.name }}
            </h1>

            <p class="mb-4 text-base text-gray-500 dark:text-gray-400 md:text-lg">
              {{ songs.length }} 首歌
              <template v-if="totalMinutes > 0">
                <span class="mx-1 text-gray-300 dark:text-gray-600">·</span>
                {{ totalMinutes }} 分钟
              </template>
            </p>

            <div class="mb-4 flex flex-wrap items-center gap-2">
              <n-button
                v-if="songs.length > 0"
                type="primary"
                :focusable="false"
                @click="onPlayAll"
              >
                <template #icon>
                  <n-icon :component="PlayOutline" />
                </template>
                全部播放
              </n-button>

              <n-button :focusable="false" @click="editing = !editing">
                <template #icon>
                  <n-icon :component="editing ? CheckmarkOutline : CreateOutline" />
                </template>
                {{ editing ? '完成' : '编辑' }}
              </n-button>

              <div v-if="editing" class="relative collection-add-wrap">
                <n-button :focusable="false" @click.stop="toggleAddPanel">
                  <template #icon>
                    <n-icon :component="AddOutline" />
                  </template>
                </n-button>
                <div
                  v-if="addPanelOpen"
                  class="collection-add-panel absolute left-0 z-20 mt-2 w-72 max-h-80 overflow-auto rounded-lg border border-gray-200 bg-white p-2 shadow-lg dark:border-gray-600 dark:bg-gray-800"
                  @click.stop
                >
                  <p v-if="!playingSongs.length" class="px-2 py-3 text-sm text-gray-500">
                    当前没有正在播放 / 队列中的歌曲
                  </p>
                  <button
                    v-for="song in playingSongs"
                    :key="song.id"
                    type="button"
                    class="flex w-full items-center gap-2 rounded px-2 py-2 text-left hover:bg-gray-100 dark:hover:bg-gray-700"
                    @click="addSongFromQueue(song)"
                  >
                    <div class="min-w-0 flex-1">
                      <div class="truncate text-sm font-medium text-gray-900 dark:text-gray-100">
                        {{ song.name }}
                      </div>
                      <div class="truncate text-xs text-gray-500">
                        {{ artistLine(song) }}
                      </div>
                    </div>
                    <span class="shrink-0 text-xs text-blue-600">加到本合集</span>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="border-t border-gray-200 px-4 py-4 dark:border-gray-700 sm:px-6">
          <div v-if="songs.length === 0" class="py-4 text-base text-gray-500 md:text-lg">
            暂无歌曲，进入编辑后可从播放列表添加
          </div>
          <div v-else>
            <div
              v-for="(song, idx) in songs"
              :key="song.id"
              class="flex items-center gap-2"
            >
              <SongItem
                class="min-w-0 flex-1"
                :song="song"
                :index="idx + 1"
                :show-cover="false"
                :show-album-name="false"
                variant="album"
                show-mv
              />
              <button
                v-if="editing"
                type="button"
                class="shrink-0 px-2 py-1 text-sm text-red-500"
                @click="removeSong(song.id)"
              >
                移除
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </PageContent>
</template>

<script setup lang="ts">
import { computed, onMounted, onBeforeUnmount, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useMessage, NButton, NIcon } from 'naive-ui'
import { AddOutline, CheckmarkOutline, CreateOutline, PlayOutline } from '@vicons/ionicons5'
import PageContent from '@/components/layout/PageContent.vue'
import SongItem from '@/components/music/MusicListItem.vue'
import {
  AddCollectionSongs,
  DetailCollection,
  RemoveCollectionSong,
  UpdateCollection,
  type CollectionDTO,
} from '@/api/collection'
import { getSongDetail } from '@/api/netease/search'
import type { ViewMusicListItem } from '@/api/view/music'
import { usePlayQueueStore } from '@/store/playQueue'
import { uploadToQiniu } from '@/util/qiniu'
import { prefetchMusicMvForQueue } from '@/player/musicMvCache'

const defaultCover = 'https://p1.music.126.net/6y-UleORITEDbvrOLV0Q8A==/5639395138885805.jpg'

const route = useRoute()
const message = useMessage()
const playQueue = usePlayQueueStore()

const loading = ref(true)
const editing = ref(false)
const addPanelOpen = ref(false)
const collection = ref<CollectionDTO | null>(null)
const songs = ref<ViewMusicListItem[]>([])
const nameDraft = ref('')
let committingName = false

const displayCover = computed(() => collection.value?.cover_url || songs.value[0]?.cover_file_url || defaultCover)

const totalMinutes = computed(() => {
  const sec = songs.value.reduce((acc, s) => acc + (Number(s.time_long) || 0), 0)
  return Math.max(0, Math.round(sec / 60))
})

const playingSongs = computed<ViewMusicListItem[]>(() => {
  const list: ViewMusicListItem[] = []
  if (playQueue.nowPlaying) list.push(playQueue.nowPlaying)
  for (const s of playQueue.queue) list.push(s)
  return list
})

function artistLine(song: ViewMusicListItem) {
  const arts = song.display_artists?.length ? song.display_artists : song.artists
  return (arts ?? []).map((a) => a.artist_name).join('/')
}

async function loadDetail(id: string) {
  loading.value = true
  try {
    const resp = await DetailCollection(id)
    collection.value = resp.item
    nameDraft.value = resp.item?.name ?? ''
    const ids = resp.item?.song_ids ?? []
    if (!ids.length) {
      songs.value = []
      return
    }
    const detail = await getSongDetail(ids)
    const byId = new Map((detail.songs ?? []).map((s) => [s.id, s]))
    songs.value = ids
      .map((sid) => {
        const d = byId.get(sid)
        if (!d) return null
        return {
          id: d.id,
          name: d.name,
          time_long: d.dt / 1000,
          album_id: d.al.id,
          album_name: d.al.name,
          cover_file_url: d.al.picUrl || defaultCover,
          artists: d.ar.map((a) => ({ artist_id: a.id, artist_name: a.name })),
          display_artists: d.ar.map((a) => ({ artist_id: a.id, artist_name: a.name })),
          mv_url: undefined,
        } as ViewMusicListItem
      })
      .filter(Boolean) as ViewMusicListItem[]
    prefetchMusicMvForQueue(ids)
  } catch (e) {
    message.error(e instanceof Error ? e.message : '加载合集失败')
    collection.value = null
    songs.value = []
  } finally {
    loading.value = false
  }
}

function onPlayAll() {
  if (!songs.value.length) return
  const r = playQueue.enqueueMany(songs.value)
  if (r === 'playing') message.success('开始播放')
  else if (r === 'added') message.success(`已加入 ${songs.value.length} 首`)
}

async function commitName() {
  if (!collection.value || committingName) return
  const name = nameDraft.value.trim()
  if (!name || name === collection.value.name) {
    nameDraft.value = collection.value.name
    return
  }
  committingName = true
  try {
    const resp = await UpdateCollection({ id: collection.value.id, name })
    if (resp.item) collection.value = { ...collection.value, ...resp.item }
  } catch (e) {
    message.error(e instanceof Error ? e.message : '改名失败')
    nameDraft.value = collection.value.name
  } finally {
    committingName = false
  }
}

async function onCoverPicked(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file || !collection.value) return
  try {
    const up = await uploadToQiniu(file, ['collection'])
    const resp = await UpdateCollection({ id: collection.value.id, cover_url: up.url })
    if (resp.item) collection.value = { ...collection.value, ...resp.item }
    message.success('封面已更新')
  } catch (err) {
    message.error(err instanceof Error ? err.message : '封面上传失败')
  }
}

function toggleAddPanel() {
  addPanelOpen.value = !addPanelOpen.value
}

async function addSongFromQueue(song: ViewMusicListItem) {
  if (!collection.value) return
  if ((collection.value.song_ids ?? []).includes(song.id)) {
    message.warning('已在合集中')
    return
  }
  try {
    const resp = await AddCollectionSongs({
      id: collection.value.id,
      songs: [
        {
          song_id: song.id,
          song_name: song.name,
          artist_names: artistLine(song),
        },
      ],
    })
    if (resp.item) collection.value = resp.item
    if (!collection.value.cover_url && song.cover_file_url) {
      try {
        const u = await UpdateCollection({
          id: collection.value.id,
          cover_url: song.cover_file_url,
        })
        if (u.item) collection.value = { ...collection.value, ...u.item }
      } catch {
        /* ignore */
      }
    }
    await loadDetail(collection.value.id)
    message.success('已加入合集')
  } catch (e) {
    message.error(e instanceof Error ? e.message : '添加失败')
  }
}

async function removeSong(songId: number) {
  if (!collection.value) return
  try {
    const resp = await RemoveCollectionSong({ id: collection.value.id, song_id: songId })
    if (resp.item) collection.value = resp.item
    songs.value = songs.value.filter((s) => s.id !== songId)
  } catch (e) {
    message.error(e instanceof Error ? e.message : '移除失败')
  }
}

function onDocClick(e: MouseEvent) {
  if (!addPanelOpen.value) return
  const t = e.target as HTMLElement | null
  if (t?.closest('.collection-add-wrap')) return
  addPanelOpen.value = false
}

watch(
  () => route.query.id,
  (id) => {
    if (id) void loadDetail(String(id))
  },
)

watch(editing, (v) => {
  if (!v) addPanelOpen.value = false
  else nameDraft.value = collection.value?.name ?? ''
})

onMounted(() => {
  document.addEventListener('click', onDocClick)
  const id = route.query.id as string
  if (id) void loadDetail(id)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', onDocClick)
})
</script>
