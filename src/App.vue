<script setup lang="ts">
import { computed } from 'vue';
import { RouterView, useRoute } from 'vue-router';
import PageHeader from './components/layout/PageHeader.vue'
// import lbAudio from './components/audio/index.vue';
import AudioContainer from "@/components/music/AudioContainer.vue";
import MusicDetail from "@/components/music/MusicDetail.vue";
import MusicDetailMobile from "@/components/music/MusicDetailMobile.vue";
import MusicLike from "@/components/music/MusicLike.vue";
import {NMessageProvider} from "naive-ui"
import NeteaseAudioHost from '@/components/player/NeteaseAudioHost.vue'
import LyricsFullScreen from '@/components/player/LyricsFullScreen.vue'
import MobileBottomPlayerBar from '@/components/player/MobileBottomPlayerBar.vue'
import MobileLyricsSheet from '@/components/player/MobileLyricsSheet.vue'
import MemoList from './components/memo/MemoList.vue';
import VoiceTranscriptBar from '@/components/voice/VoiceTranscriptBar.vue'

const route = useRoute()
const showMobilePlayerBar = computed(() => route.path !== '/ktv')
</script>

<template>
  <div class="app-shell min-h-screen w-full text-left">
    <VoiceTranscriptBar />
    <!-- 顶部导航栏 -->
  <PageHeader />

  <n-message-provider>
    <NeteaseAudioHost />
    <LyricsFullScreen />
    <MobileBottomPlayerBar v-if="showMobilePlayerBar" />
    <MobileLyricsSheet />
<!--    <MusicDetail />-->
<!--      <MusicDetailMobile />-->
<!--    <MusicLike />-->
<!--    <MemoList />-->
<!--    <ChatWindow />-->
      
      <!-- 路由视图：为固定顶栏留出空间（移动端含安全区） -->
      <div class="app-main">
        <RouterView />
      </div>

      <footer class="app-icp">
        <a
          href="https://beian.miit.gov.cn/"
          target="_blank"
          rel="noopener noreferrer"
          class="app-icp__link"
        >
          鄂ICP备2024057031号
        </a>
      </footer>
      
<!--  <AudioContainer></AudioContainer>-->
  </n-message-provider>
  </div>
</template>

<style scoped>
.app-shell {
  text-align: left;
}

.app-icp {
  text-align: center;
  font-size: 0.75rem;
  line-height: 1.5;
  color: rgb(156 163 175);
  padding: 0.5rem var(--page-gutter-x) calc(var(--app-player-bar-offset) + 0.25rem);
}

@media (min-width: 1024px) {
  .app-icp {
    padding-bottom: 0.75rem;
  }
}

.app-icp__link {
  color: inherit;
  text-decoration: none;
}

.app-icp__link:hover {
  color: rgb(107 114 128);
  text-decoration: underline;
}

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
