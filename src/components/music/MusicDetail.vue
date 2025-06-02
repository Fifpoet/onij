<template>
  <Transition name="fade">
    <div v-if="musicStore.midShowWhat == MidShowWhat.ShowMusicDetail" class="music-detail-container">
      <div class="music-detail-content relative p-8 flex max-w-4xl w-full">
        <!-- 左侧封面 -->
        <div class="cover-container w-[300px] h-[300px] mr-8">
          <img 
            src="https://assets.msn.cn/weathermapdata/1/static/weather/Icons/taskbar_v10/Condition_Card/PartlyCloudyNightV2.svg" 
            alt="album cover" 
            class="w-full h-full object-cover rounded-lg shadow-md"
          >
        </div>

        <!-- 右侧信息 -->
        <div class="info-container flex-1">
          <!-- 音乐名称和专辑信息 -->
          <div class="flex items-center mb-4">
            <h1 class="text-3xl font-bold">{{ musicStore.current.detail?.name }}</h1>
            <span class="mx-2 text-gray-400">-</span>
            <span class="text-lg text-gray-600">{{ musicStore.current.detail?.album_name }}</span>
          </div>

          <!-- 创作者信息 -->
          <div class="creator-info mb-4 flex items-center">
            <span class="text-gray-600 mr-2">作词：</span>
            <router-link 
              :to="`/artist/${musicStore.current.detail?.writer_id}`"
              class="text-blue-500 hover:text-blue-700"
            >
              {{ musicStore.current.detail?.writer_name }}
            </router-link>
            <span class="mx-2 text-gray-400">  </span>
            <span class="text-gray-600 mr-2">作曲：</span>
            <router-link 
              :to="`/artist/${musicStore.current.detail?.composer_id}`"
              class="text-blue-500 hover:text-blue-700"
            >
              {{ musicStore.current.detail?.composer_name }}
            </router-link>
          </div>

          <!-- 歌词容器 -->
          <div class="lyrics-container mt-6 h-[200px] overflow-y-auto" ref="lyricsContainerRef">
            <div v-if="lyrics.length > 0">
              <p 
                v-for="(line, index) in lyrics" 
                :key="index"
                :class="{ 'active-lyric': currentLyricIndex === index }"
                class="lyric-line py-1"
              >
                {{ line.text }}
              </p>
            </div>
            <p v-else class="text-gray-500 italic">
              歌词加载中...
            </p>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { MidShowWhat } from "@/api/types";
import { useMusicStore } from "@/store/music.ts";
import { watch, ref, onMounted, onUnmounted } from 'vue';

const musicStore = useMusicStore();
const lyrics = ref<{ time: number; text: string }[]>([]);
const currentLyricIndex = ref(-1);
const lyricsContainerRef = ref<HTMLElement | null>(null);

// 解析歌词文件
const parseLyrics = (lyricsText: string) => {
  const lines = lyricsText.split('\n');
  const result: { time: number; text: string }[] = [];
  
  // 正则表达式匹配时间戳 [mm:ss.ms]
  const timeRegex = /\[(\d{2}):(\d{2})\.(\d{1,3})\](.*)/;
  
  lines.forEach(line => {
    const match = line.match(timeRegex);
    if (match) {
      const minutes = parseInt(match[1]);
      const seconds = parseInt(match[2]);
      const milliseconds = parseInt(match[3]);
      
      // 转换为总毫秒数
      const timeInMs = minutes * 60 * 1000 + seconds * 1000 + milliseconds;
      
      // 提取歌词文本并去除前后空格
      const text = match[4].trim();
      
      result.push({
        time: timeInMs,
        text: text
      });
    }
  });
  
  // 按时间排序
  return result.sort((a, b) => a.time - b.time);
};

// 获取歌词文件
const fetchLyrics = async (url: string) => {
  try {
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error('获取歌词失败');
    }
    const text = await response.text();
    lyrics.value = parseLyrics(text);
  } catch (error) {
    console.error('获取歌词出错:', error);
    lyrics.value = [];
  }
};

// 更新当前歌词索引
const updateCurrentLyric = () => {
  if (!musicStore.current.detail || lyrics.value.length === 0) return;
  
  // 计算当前播放时间（毫秒）
  const currentTimeMs = musicStore.current.progress / 100 * getDuration();
  
  // 查找当前应该显示的歌词
  let index = -1;
  for (let i = 0; i < lyrics.value.length; i++) {
    if (lyrics.value[i].time <= currentTimeMs) {
      index = i;
    } else {
      break;
    }
  }
  
  if (currentLyricIndex.value !== index) {
    currentLyricIndex.value = index;
    scrollToActiveLyric();
  }
};

// 获取歌曲总时长（毫秒）
const getDuration = () => {
  // 这里需要根据实际情况获取歌曲时长
  // 如果 musicStore 中有 duration 属性，可以直接使用
  // 暂时使用一个固定值（4分钟）
  return 4 * 60 * 1000;
};

// 滚动到当前歌词
const scrollToActiveLyric = () => {
  if (currentLyricIndex.value >= 0 && lyricsContainerRef.value) {
    const container = lyricsContainerRef.value;
    const activeElement = container.querySelector('.active-lyric');
    
    if (activeElement) {
      const containerHeight = container.clientHeight;
      const elementTop = (activeElement as HTMLElement).offsetTop;
      const elementHeight = (activeElement as HTMLElement).clientHeight;
      
      // 将当前歌词滚动到容器中间位置
      container.scrollTop = elementTop - containerHeight / 2 + elementHeight / 2;
    }
  }
};

// 定时更新当前歌词
let intervalId: number | null = null;

onMounted(() => {
  // 每100毫秒更新一次当前歌词
  intervalId = window.setInterval(updateCurrentLyric, 100);
});

onUnmounted(() => {
  if (intervalId !== null) {
    clearInterval(intervalId);
  }
});

// 监听当前音乐的变化
watch(
  () => musicStore.current.detail,
  (newVal) => {
    if (newVal) {
      // 当歌曲变化时，获取新的歌词
      if (newVal.lyrics_file_url) {
        fetchLyrics(newVal.lyrics_file_url);
      } else {
        lyrics.value = [];
      }
    }
  },
  { immediate: true }
);

// 监听播放进度变化
watch(
  () => musicStore.current.progress,
  () => {
    updateCurrentLyric();
  }
);


</script>

<style scoped>
.music-detail-container {
  width: 100%;
  padding: 2rem;
}

.music-detail-content {
  background-color: transparent;
  margin: 0 auto;
}

.lyrics-container {
  /* 自定义滚动条样式 */
  scrollbar-width: thin;
  scrollbar-color: #cbd5e0 #f7fafc;
}

.lyrics-container::-webkit-scrollbar {
  width: 6px;
}

.lyrics-container::-webkit-scrollbar-track {
  background: #f7fafc;
}

.lyrics-container::-webkit-scrollbar-thumb {
  background-color: #cbd5e0;
  border-radius: 3px;
}

/* 添加过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: all 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

/* 歌词样式 */
.lyric-line {
  transition: all 0.3s ease;
  color: rgba(75, 85, 99, 0.8);
  font-size: 14px;
}

.active-lyric {
  color: #3b82f6;
  font-size: 16px;
  font-weight: 600;
}
</style>