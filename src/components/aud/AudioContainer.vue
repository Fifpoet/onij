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

    <!-- 音乐播放器主体（后续实现） -->
    <div class="audio-content flex-grow pl-5">
      <!-- 播放器内容区域 -->
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import apiClient from '@/util/http.ts'; // 引入 axios 实例
import { useMusicStore } from "@/store/music.ts";

const audioContainer = ref<HTMLDivElement | null>(null);
let isDragging = false;
let offset = { x: 0, y: 0 };

const listMusicReq= {
  "name": "",
}

const fetchMusicList = async () => {
  console.log(listMusicReq)
  const musicList = await apiClient.post('/music/list', listMusicReq);
  useMusicStore().setMusicList(musicList);
};

onMounted(()=> {
  fetchMusicList();
})

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
