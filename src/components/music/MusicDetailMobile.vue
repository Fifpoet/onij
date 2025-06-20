<!-- 移动端音乐详情 -->
<template>
  <Transition name="fade-scale" mode="out-in">
    <div 
      v-if="musicStore.midShowWhat == MidShowWhat.ShowMusicDetail" 
      class="fixed inset-0 bg-white dark:bg-dark-800 z-0 lg:hidden pt-[60px] flex flex-col"
    >
      <!-- 顶部信息区 -->
      <div class="px-4 py-3 flex-shrink-0">
        <!-- 歌曲名 -->
        <h1 class="text-lg font-bold truncate m-0">{{ musicStore.current.detail?.name }}</h1>
        
        <!-- 歌手信息和作词作曲 -->
        <div class="mt-2 flex justify-end items-center gap-3 text-sm">
          <!-- 歌手信息 -->
          <div class="text-gray-600 dark:text-gray-300">
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

          <span class="text-gray-300">|</span>

          <!-- 作词作曲 -->
          <div class="flex items-center gap-2 text-gray-500">
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
      </div>

      <!-- 主要内容区 -->
      <div class="flex-1 relative" @click="toggleView">
        <!-- 封面 -->
        <Transition name="fade" mode="out-in">
          <div 
            v-show="!showLyrics"
            class="absolute inset-0 flex items-center justify-center p-8"
          >
            <div class="w-full aspect-square relative">
              <img 
                :src="musicStore.current.detail?.cover_url" 
                :alt="musicStore.current.detail?.name" 
                class="w-full h-full object-cover rounded-lg"
              >
              <div class="absolute inset-0 flex items-center justify-center bg-black/20 rounded-lg opacity-0 hover:opacity-100 transition-opacity">
                <span class="text-white text-sm">点击查看歌词</span>
              </div>
            </div>
          </div>
        </Transition>

        <!-- 歌词 -->
        <Transition name="fade" mode="out-in">
          <div 
            v-show="showLyrics"
            class="absolute inset-0 flex flex-col items-center justify-center p-6"
          >
            <div class="w-full flex flex-col items-center space-y-2">
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
            <div class="absolute inset-0 flex items-center justify-center bg-black/5 opacity-0 hover:opacity-100 transition-opacity">
              <span class="text-gray-600 dark:text-gray-300 text-sm">点击查看封面</span>
            </div>
          </div>
        </Transition>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { MidShowWhat } from "@/api/types";
import { useMusicStore } from "@/store/music";
import { ref, watch, computed } from 'vue';

const musicStore = useMusicStore();
const showLyrics = ref(false);
const lyrics = ref<{ time: number; text: string }[]>([]);
const currentLyricIndex = ref(-1);

const toggleView = () => {
  showLyrics.value = !showLyrics.value;
};

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

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease-out;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style> 