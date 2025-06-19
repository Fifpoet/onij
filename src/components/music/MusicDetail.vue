<template>
  <Transition name="fade-scale" mode="out-in">
    <div 
      v-if="musicStore.midShowWhat == MidShowWhat.ShowMusicDetail" 
      class="fixed inset-x-0 top-1/2 -translate-y-1/2 bg-white dark:bg-dark-800 z-0"
    >
      <div class="max-w-4xl mx-auto p-6">
        <div class="flex gap-6 h-[300px]">
          <!-- 封面 -->
          <div class="w-[300px] h-[300px] flex-shrink-0">
          <img 
              :src="musicStore.current.detail?.cover_url" 
              :alt="musicStore.current.detail?.name" 
              class="w-full h-full object-cover rounded-lg"
          >
        </div>

          <!-- 信息区域 -->
          <div class="flex-1 flex flex-col min-w-0 h-full">
            <!-- 标题区域 -->
            <div class="mb-3 flex-shrink-0">
              <div class="flex items-end justify-between">
                <h1 class="text-xl font-bold truncate flex-1 leading-none m-0">{{ musicStore.current.detail?.name }}</h1>
                <div class="flex items-center gap-1 text-sm ml-2">
                  <template v-for="(name, index) in musicStore.current.detail?.artist_names || []" :key="index">
                    <router-link 
                      :to="`/artist/${musicStore.current.detail?.artist_ids[index]}`"
                      class="text-primary hover:text-primary-600 transition-colors"
                    >
                      {{ name }}
                    </router-link>
                    <span v-if="index < (musicStore.current.detail?.artist_names.length || 0) - 1" class="text-gray-400">/</span>
                  </template>
                </div>
          </div>
              <div class="flex justify-end items-center gap-4 mt-1 text-sm text-gray-500">
                <div class="flex items-center gap-1">
                  <span>词:</span>
            <router-link 
              :to="`/artist/${musicStore.current.detail?.writer_id}`"
                    class="text-primary hover:text-primary-600 transition-colors"
            >
              {{ musicStore.current.detail?.writer_name }}
            </router-link>
                </div>
                <div class="flex items-center gap-1">
                  <span>曲:</span>
            <router-link 
              :to="`/artist/${musicStore.current.detail?.composer_id}`"
                    class="text-primary hover:text-primary-600 transition-colors"
            >
              {{ musicStore.current.detail?.composer_name }}
            </router-link>
                </div>
              </div>
          </div>

          <!-- 歌词容器 -->
            <div class="flex-1 flex flex-col justify-center items-center space-y-2">
              <div 
                v-for="lyric in displayLyrics" 
                :key="lyric.index"
                :class="[
                  'transition-all duration-300 text-center',
                  lyric.index === currentLyricIndex 
                    ? 'text-primary text-base font-medium'
                    : 'text-gray-500 dark:text-gray-400 text-sm'
                ]"
              >
                {{ lyric.text }}
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
import { ref, watch, computed } from 'vue';

const musicStore = useMusicStore();
const lyrics = ref<{ time: number; text: string }[]>([]);
const currentLyricIndex = ref(-1);

// 计算要显示的歌词
const displayLyrics = computed(() => {
  if (lyrics.value.length === 0) return [];
  
  const index = currentLyricIndex.value;
  const total = lyrics.value.length;
  const count = 5; // 显示的歌词数量
  
  let indices: number[];
  if (index < 0) {
    indices = [0, 1, 2, 3, 4];
  } else if (index < 2) {
    indices = [0, 1, 2, 3, 4];
  } else if (index > total - 3) {
    indices = [total - 5, total - 4, total - 3, total - 2, total - 1];
  } else {
    indices = [index - 2, index - 1, index, index + 1, index + 2];
  }
  
  return indices
    .filter(i => i >= 0 && i < total)
    .slice(0, count)
    .map(i => ({
      index: i,
      text: lyrics.value[i].text
    }));
  });
  
// 更新当前歌词
watch(() => musicStore.current.progress, () => {
  if (!lyrics.value.length) return;
  
  const currentTime = musicStore.current.progress / 100 * (document.querySelector('audio')?.duration || 240) * 1000;
  
  const index = lyrics.value.findIndex((lyric, i) => {
    const nextTime = lyrics.value[i + 1]?.time ?? Infinity;
    return lyric.time <= currentTime && currentTime < nextTime;
  });
  
  if (index !== currentLyricIndex.value) {
    currentLyricIndex.value = index;
  }
});

// 监听歌曲变化
watch(() => musicStore.current.detail?.lyrics_file_url, async (url) => {
  if (!url) {
    lyrics.value = [];
    return;
  }
  
  try {
    const text = await fetch(url).then(r => r.text());
    lyrics.value = text
      .split('\n')
      .map(line => {
        const match = line.match(/\[(\d{2}):(\d{2})\.(\d{1,3})\](.*)/);
        if (!match) return null;
        
        const [, min, sec, ms] = match;
        return {
          time: (parseInt(min) * 60 + parseInt(sec)) * 1000 + parseInt(ms),
          text: match[4].trim()
        };
      })
      .filter((item): item is { time: number; text: string } => item !== null)
      .sort((a, b) => a.time - b.time);
    
    currentLyricIndex.value = -1;
  } catch (error) {
    console.error('加载歌词失败:', error);
        lyrics.value = [];
      }
}, { immediate: true });
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