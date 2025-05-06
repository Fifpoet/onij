<template>
  <Transition name="fade">
    <div v-if="show" class="music-detail-container">
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
          <!-- 音乐名称 -->
          <h1 class="text-3xl font-bold mb-4">{{ musicStore.current.detail?.name }}</h1>

          <!-- 专辑信息 -->
          <div class="album-info mb-4">
            <router-link 
              :to="`/album/${musicStore.current.detail?.album_id}`"
              class="text-blue-500 hover:text-blue-700 text-lg"
            >
              专辑：{{ musicStore.current.detail?.album_name }}
            </router-link>
          </div>

          <!-- 艺术家信息 -->
          <div class="artist-info mb-4">
            <span class="text-gray-600">演唱：</span>
            <span v-for="(name, index) in musicStore.current.detail?.artist_names" :key="index">
              <router-link 
                :to="`/artist/${musicStore.current.detail?.artist_ids[index]}`"
                class="text-blue-500 hover:text-blue-700"
              >
                {{ name }}
              </router-link>
              <span v-if="index < musicStore.current.detail?.artist_names.length - 1">, </span>
            </span>
          </div>

          <!-- 创作者信息 -->
          <div class="creator-info mb-4">
            <div class="mb-2">
              <span class="text-gray-600">作曲：</span>
              <router-link 
                :to="`/artist/${musicStore.current.detail?.composer_id}`"
                class="text-blue-500 hover:text-blue-700"
              >
                {{ musicStore.current.detail?.composer_name }}
              </router-link>
            </div>
            <div>
              <span class="text-gray-600">作词：</span>
              <router-link 
                :to="`/artist/${musicStore.current.detail?.writer_id}`"
                class="text-blue-500 hover:text-blue-700"
              >
                {{ musicStore.current.detail?.writer_name }}
              </router-link>
            </div>
          </div>

          <!-- 歌词占位 -->
          <div class="lyrics-container mt-6 h-[200px] overflow-y-auto">
            <p class="text-gray-500 italic">
              歌词加载中...
            </p>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { useMusicStore } from "@/store/music.ts";
import { watch, ref } from 'vue';

const musicStore = useMusicStore();
const show = ref(false);

// 监听当前音乐的变化
watch(
  () => musicStore.current.detail,
  (newVal) => {
    if (newVal) {
      show.value = true;
    }
  }
);

const handleClose = () => {
  show.value = false;
}
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
</style>