<template>
  <Transition name="fade-scale" mode="out-in">
    <div 
      v-if="musicStore.midShowWhat == MidShowWhat.ShowMusicDetail" 
      class="fixed inset-0 z-0 flex items-center justify-center pb-20"
    >
      <div class="bg-white dark:bg-dark-800 max-w-4xl w-full mx-4 p-6 flex gap-6">
        <!-- 封面 -->
        <div class="w-200px h-200px flex-shrink-0">
          <img 
            :src="musicStore.current.detail?.cover_url" 
            :alt="musicStore.current.detail?.name" 
            class="w-full h-full object-cover rounded-lg"
          >
        </div>

        <!-- 信息区域 -->
        <div class="flex-grow flex flex-col min-w-0">
          <!-- 标题区域 -->
          <div class="mb-3">
            <h1 class="text-xl font-bold truncate">{{ musicStore.current.detail?.name }}</h1>
            <p class="text-gray-500 dark:text-gray-400 text-sm mt-1">{{ musicStore.current.detail?.album_name }}</p>
          </div>

          <!-- 创作者信息 -->
          <div class="space-y-1 mb-4">
            <div class="flex items-center gap-2 text-sm">
              <span class="text-gray-600 dark:text-gray-400">作词：</span>
              <router-link 
                :to="`/artist/${musicStore.current.detail?.writer_id}`"
                class="text-primary hover:text-primary-600 transition-colors"
              >
                {{ musicStore.current.detail?.writer_name }}
              </router-link>
            </div>
            <div class="flex items-center gap-2 text-sm">
              <span class="text-gray-600 dark:text-gray-400">作曲：</span>
              <router-link 
                :to="`/artist/${musicStore.current.detail?.composer_id}`"
                class="text-primary hover:text-primary-600 transition-colors"
              >
                {{ musicStore.current.detail?.composer_name }}
              </router-link>
            </div>
          </div>

          <!-- 歌词容器 -->
          <div 
            ref="lyricsContainerRef"
            class="flex-grow h-200px overflow-hidden relative"
          >
            <div 
              class="absolute inset-0 overflow-y-auto scrollbar-hide"
              :style="{ transform: `translateY(${scrollOffset}px)` }"
            >
              <div 
                v-for="(line, index) in lyrics" 
                :key="index"
                :ref="el => { if (index === currentLyricIndex) activeLyricRef = el as HTMLElement }"
                :class="[
                  'py-1 transition-all duration-300 text-center text-sm',
                  currentLyricIndex === index 
                    ? 'text-primary text-base font-medium' 
                    : 'text-gray-500 dark:text-gray-400'
                ]"
              >
                {{ line.text }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { MidShowWhat } from "@/api/types";
import { useMusicStore } from "@/store/music";
import { ref, watch, onMounted, onUnmounted, computed } from 'vue';

const musicStore = useMusicStore();
const lyrics = ref<{ time: number; text: string }[]>([]);
const currentLyricIndex = ref(-1);
const lyricsContainerRef = ref<HTMLElement | null>(null);
const activeLyricRef = ref<HTMLElement | null>(null);
const scrollOffset = ref(0);

// 解析歌词
const parseLyrics = (lyricsText: string) => {
  const lines = lyricsText.split('\n');
  const timeRegex = /\[(\d{2}):(\d{2})\.(\d{1,3})\](.*)/;
  
  return lines
    .map(line => {
      const match = line.match(timeRegex);
      if (!match) return null;
      
      const [, minutes, seconds, milliseconds] = match;
      const timeInMs = (parseInt(minutes) * 60 + parseInt(seconds)) * 1000 + parseInt(milliseconds);
      const text = match[4].trim();
      
      return { time: timeInMs, text };
    })
    .filter((item): item is { time: number; text: string } => item !== null)
    .sort((a, b) => a.time - b.time);
};

// 获取歌词
const fetchLyrics = async (url: string) => {
  try {
    const response = await fetch(url);
    if (!response.ok) throw new Error('Failed to fetch lyrics');
    const text = await response.text();
    lyrics.value = parseLyrics(text);
  } catch (error) {
    console.error('Error fetching lyrics:', error);
    lyrics.value = [];
  }
};

// 计算当前歌词
const updateCurrentLyric = () => {
  if (!musicStore.current.detail || lyrics.value.length === 0) return;
  
  const currentTimeMs = (musicStore.current.progress / 100) * getDuration();
  
  const newIndex = lyrics.value.findIndex((lyric, index) => {
    const nextLyric = lyrics.value[index + 1];
    return lyric.time <= currentTimeMs && (!nextLyric || nextLyric.time > currentTimeMs);
  });

  if (currentLyricIndex.value !== newIndex) {
    currentLyricIndex.value = newIndex;
    updateScroll();
  }
};

// 获取歌曲时长
const getDuration = () => {
  return musicStore.current.detail?.duration || 240000; // 默认4分钟
};

// 更新滚动位置
const updateScroll = () => {
  if (!lyricsContainerRef.value || !activeLyricRef.value) return;
  
  const containerHeight = lyricsContainerRef.value.clientHeight;
  const lyricHeight = activeLyricRef.value.clientHeight;
  
  // 计算目标偏移量，使当前歌词位于容器中央
  const targetOffset = -(activeLyricRef.value.offsetTop - (containerHeight / 2) + (lyricHeight / 2));
  
  // 使用动画平滑滚动
  scrollOffset.value = targetOffset;
};

// 设置定时器
let updateInterval: number | null = null;

onMounted(() => {
  updateInterval = window.setInterval(updateCurrentLyric, 100);
});

onUnmounted(() => {
  if (updateInterval !== null) {
    clearInterval(updateInterval);
  }
});

// 监听歌曲变化
watch(
  () => musicStore.current.detail,
  async (newVal) => {
    if (newVal?.lyrics_file_url) {
      await fetchLyrics(newVal.lyrics_file_url);
      currentLyricIndex.value = -1;
      scrollOffset.value = 0;
    } else {
      lyrics.value = [];
    }
  },
  { immediate: true }
);

// 监听进度变化
watch(() => musicStore.current.progress, updateCurrentLyric);
</script>

<style scoped>
.fade-scale-enter-active,
.fade-scale-leave-active {
  transition: all 0.3s ease-out;
}

.fade-scale-enter-from,
.fade-scale-leave-to {
  opacity: 0;
  transform: scale(0.9);
}

.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.scrollbar-hide::-webkit-scrollbar {
  display: none;
}
</style>