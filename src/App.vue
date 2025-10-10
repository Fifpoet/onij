<script setup lang="ts">
import { ref } from 'vue';
import { RouterView } from 'vue-router';
import PageHeader from './components/layout/PageHeader.vue'
// import lbAudio from './components/audio/index.vue';
import AudioContainer from "@/components/music/AudioContainer.vue";
import MusicDetail from "@/components/music/MusicDetail.vue";
import MusicDetailMobile from "@/components/music/MusicDetailMobile.vue";
import MusicLike from "@/components/music/MusicLike.vue";
import {NMessageProvider} from "naive-ui"
import MemoList from './components/memo/MemoList.vue';
import FileList from './components/file/FileList.vue';
import { FileDetail } from '@/api/types/file';

// 当前文件夹ID
const currentFolderId = ref(0);

// 处理文件点击事件
const handleFileClick = (file: FileDetail) => {
  if (file.format === 99) { // FileType.FT_Folder
    // 更新当前文件夹ID
    currentFolderId.value = file.id;
    console.log('切换到文件夹:', file.name, 'ID:', file.id);
  }
};
</script>

<template>
  <div>
    <!-- 顶部导航栏 -->
  <PageHeader />

  <n-message-provider>
<!--    <MusicDetail />-->
<!--      <MusicDetailMobile />-->
<!--    <MusicLike />-->
<!--    <MemoList />-->
<!--    <ChatWindow />-->
      
      <!-- 路由视图 -->
      <div class="pt-[60px]">
        <RouterView />
      </div>

<!--      &lt;!&ndash; 文件管理区域 &ndash;&gt;-->
<!--      <div class="h-full">-->
<!--        &lt;!&ndash; 文件列表 &ndash;&gt;-->
<!--        <FileList -->
<!--          :parent-id="currentFolderId"-->
<!--          @file-click="handleFileClick"-->
<!--        />-->
<!--      </div>-->
      
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
