<template>
  <div
      ref="audioContainer"
      class="audio-container fixed bg-[rgb(245,245,245)] rounded-lg shadow-lg flex items-center bottom-5 left-5 w-[400px] h-[60px]"
  >
    <!-- 左侧两竖排小点 -->
    <div
        class="drag-handle flex justify-between p-2 cursor-grab"
        @mousedown="startDragging"
    >
      <!-- 第一排小点 -->
      <div class="flex flex-col">
        <div class="dot w-[4px] h-[4px] bg-gray-500 rounded-full mb-1" v-for="n in 3" :key="'left'+n"></div>
      </div>
      <!-- 第二排小点 -->
      <div class="flex flex-col ml-1">
        <div class="dot w-[4px] h-[4px] bg-gray-500 rounded-full mb-1" v-for="n in 3" :key="'right'+n"></div>
      </div>
    </div>

    <!-- 音乐播放器主体，显示当前播放的音乐 -->
    <div class="audio-content flex-grow pl-5">
      <div v-if="currentMusicDetail">
        <strong>{{ currentMusicDetail.title }}</strong> - {{ currentMusicDetail.artist }}
      </div>
      <div v-else>
        <p>请选择一首音乐播放</p>
      </div>
    </div>

    <!-- 右侧展示音乐列表的图标 -->
    <div class="music-list-toggle p-2 cursor-pointer" @click="toggleMusicList">
      🎵 <!-- 可以替换为你想要的图标 -->
    </div>

    <!-- 音乐列表展示 -->
    <div v-if="showMusicList" class="music-list absolute bg-white shadow-lg rounded-lg p-4 w-[400px] bottom-[70px] left-0">
      <ul>
        <li
            v-for="music in musicStore.MusicList"
            :key="music.id"
            class="mb-2 cursor-pointer"
            @dblclick="playMusic(music.id)"
        >
          <strong>{{ music.title }}</strong> - {{ music.artist }}
        </li>
      </ul>
    </div>
  </div>
</template>


<script lang="ts" setup>

import {onMounted, ref} from 'vue';
import apiClient from '@/util/http.ts'; // 引入 axios 实例
import {useMusicStore} from "@/store/music.ts";
// *************** 音乐列表展示逻辑 *************** //
const showMusicList = ref(false);
const currentMusicDetail = ref(null); // 用于保存当前音乐详情
const musicStore = useMusicStore(); // 获取 Pinia store


// *************** API操作 *************** //
const listMusicReq = {
  "title": "",
  "artist": 1,
  "perform_type": 1,
  "page": 1,
  "size": 5
};

const fetchMusicList = async () => {
  const response = await apiClient.post('/music/list', listMusicReq);
  const musicList = response.data.data;
  musicStore.setMusicList(musicList);
};

// 播放音乐，获取音乐详情
const playMusic = async (id: number) => {
  try {
    const response = await apiClient.get(`/music/detail/${id}`); // 获取音乐详情的 API
    currentMusicDetail.value = response.data; // 更新当前音乐详情
    musicStore.setCurrentMusic(id); // 更新 Pinia store 中的当前播放的音乐 id
  } catch (error) {
    console.error('获取音乐详情失败', error);
  }
};

// *************** 拖动操作 *************** //
const audioContainer = ref<HTMLDivElement | null>(null);
let isDragging = false;
let offset = { x: 0, y: 0 };

// 切换音乐列表展示
const toggleMusicList = () => {
  showMusicList.value = !showMusicList.value;
};

// 获取音乐列表
onMounted(() => {
  fetchMusicList();
});



const startDragging = (e: MouseEvent) => {
  if (!audioContainer.value) return;

  isDragging = true;
  offset.x = e.clientX - audioContainer.value.getBoundingClientRect().left;
  offset.y = e.clientY - audioContainer.value.getBoundingClientRect().top;

  // 禁用文本选择
  document.body.style.userSelect = 'none';

  document.addEventListener('mousemove', drag);
  document.addEventListener('mouseup', stopDragging);
};

const stopDragging = () => {
  isDragging = false;

  // 恢复文本选择
  document.body.style.userSelect = '';

  document.removeEventListener('mousemove', drag);
  document.removeEventListener('mouseup', stopDragging);
};

const drag = (e: MouseEvent) => {
  if (!isDragging || !audioContainer.value) return;

  // 计算新的位置
  let newLeft = e.clientX - offset.x;
  let newTop = e.clientY - offset.y;

  // 获取窗口宽高
  const windowWidth = window.innerWidth;
  const windowHeight = window.innerHeight;

  // 限制新位置不能超出窗口范围
  const containerWidth = audioContainer.value.offsetWidth;
  const containerHeight = audioContainer.value.offsetHeight;

  // 限制左、右、上、下的最大最小值
  if (newLeft < 0) newLeft = 0;
  if (newLeft + containerWidth > windowWidth) newLeft = windowWidth - containerWidth;
  if (newTop < 0) newTop = 0;
  if (newTop + containerHeight > windowHeight) newTop = windowHeight - containerHeight;

  // 设置新位置
  audioContainer.value.style.left = `${newLeft}px`;
  audioContainer.value.style.top = `${newTop}px`;
};


</script>

<!-- Unocss classes used -->
