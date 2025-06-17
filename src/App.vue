<script setup lang="ts">
import { ref } from 'vue';
import { RouterView } from 'vue-router';
import PageHeader from './components/layout/PageHeader.vue'
// import lbAudio from './components/audio/index.vue';
import AudioContainer from "@/components/music/AudioContainer.vue";
import MusicDetail from "@/components/music/MusicDetail.vue";
import MusicLike from "@/components/music/MusicLike.vue";
import {NMessageProvider} from "naive-ui"
import MemoList from './components/memo/MemoList.vue';
import FileList from './components/file/FileList.vue';
import { FileDetail } from '@/api/types/file';
import { NButton, NIcon } from 'naive-ui';

// 文件夹导航状态
const currentFolderId = ref(0);
const currentFolderName = ref('根目录');
const folderHistory = ref<number[]>([0]);
const folderNameHistory = ref<string[]>(['根目录']);

// 处理文件点击事件
const handleFileClick = (file: FileDetail) => {
  if (file.format === 99) { // FileType.FT_Folder
    // 进入文件夹
    currentFolderId.value = file.id;
    currentFolderName.value = file.name;
    folderHistory.value.push(file.id);
    folderNameHistory.value.push(file.name);
    console.log('进入文件夹:', file.name, 'ID:', file.id);
  }
};

// 返回上级文件夹
const goBack = () => {
  if (folderHistory.value.length > 1) {
    folderHistory.value.pop();
    folderNameHistory.value.pop();
    currentFolderId.value = folderHistory.value[folderHistory.value.length - 1];
    currentFolderName.value = folderNameHistory.value[folderNameHistory.value.length - 1];
  }
};

</script>

<template>
  <div>
    <!-- 顶部导航栏 -->
    <PageHeader />
    
    <n-message-provider>
      <MusicDetail />
      <MusicLike />
      <MemoList />
      
      <!-- 文件管理区域 -->
      <div class="min-h-screen pt-[60px]">
        <!-- 面包屑导航 -->
        <div class="px-5 pt-5 pb-3">
          <div v-if="currentFolderId !== 0" class="flex items-center space-x-3">
            <n-button 
              size="small" 
              type="primary" 
              ghost 
              @click="goBack"
              class="flex items-center"
            >
              <template #icon>
                <n-icon>
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="m15 18-6-6 6-6"/>
                  </svg>
                </n-icon>
              </template>
              返回上级
            </n-button>
            <span class="text-gray-400">|</span>
            <span class="text-sm text-gray-600">{{ currentFolderName }}</span>
          </div>
          <div v-else class="h-8 flex items-center">
            <span class="text-sm text-gray-600">根目录</span>
          </div>
        </div>
        
        <!-- 文件列表 -->
        <FileList 
          :parent-id="currentFolderId"
          @file-click="handleFileClick"
        />
      </div>
      
      <AudioContainer></AudioContainer>
    </n-message-provider>
  </div>
</template>

<style scoped>
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}
.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}
.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}
</style>
