<template>
  <div class="music-detail-container flex justify-center items-center min-h-screen">
    <div class="music-detail-content bg-white rounded-lg shadow-lg p-8 flex max-w-4xl w-full">
      <!-- 左侧封面 -->
      <div class="cover-container w-[300px] h-[300px] mr-8">
        <img 
          src="/default-cover.jpg" 
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
</template>

<script setup lang="ts">
import { useMusicStore } from "@/store/music.ts";

const musicStore = useMusicStore();

// 移除了 Props 接口和 defineProps，因为我们现在直接使用 store
</script>

<style scoped>
.music-detail-container {
  background-color: #f5f5f5;
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
</style>