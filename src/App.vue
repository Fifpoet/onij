<script setup lang="ts">
import { computed } from 'vue';
import { RouterView, useRoute } from 'vue-router';
import PageHeader from './components/layout/PageHeader.vue'
import { useOnijNavToggle } from '@/composables/useNavToolsUnlock'
import {NMessageProvider} from "naive-ui"
import NeteaseAudioHost from '@/components/player/NeteaseAudioHost.vue'
import VideoPlayerHost from '@/components/player/VideoPlayerHost.vue'
import LyricsFullScreen from '@/components/player/LyricsFullScreen.vue'
import MobileBottomPlayerBar from '@/components/player/MobileBottomPlayerBar.vue'
import MobileLyricsSheet from '@/components/player/MobileLyricsSheet.vue'
import VoiceTranscriptBar from '@/components/voice/VoiceTranscriptBar.vue'
import { useVideoPlayerStore } from '@/store/videoPlayer'
import { usePresenceStore } from '@/store/presence'
import { usePlaySourcesStore } from '@/store/playSources'

const route = useRoute()
const showMobilePlayerBar = computed(() => route.path !== '/ktv')
const showIcp = computed(() => route.path !== '/ktv')
const videoPlayer = useVideoPlayerStore()
const { navToolsVisible } = useOnijNavToggle()
usePresenceStore().startHeartbeat()
usePlaySourcesStore().loadFromServer()
</script>

<template>
  <div class="app-shell min-h-screen w-full text-left">
    <VoiceTranscriptBar v-if="navToolsVisible && !videoPlayer.opened" />
  <PageHeader />

  <n-message-provider class="app-body">
    <NeteaseAudioHost />
    <VideoPlayerHost />
    <LyricsFullScreen />
    <MobileBottomPlayerBar v-if="showMobilePlayerBar" />
    <MobileLyricsSheet />
      <div class="app-main">
        <RouterView />
      </div>

      <footer v-if="showIcp" class="app-icp">
        <a
          href="https://beian.miit.gov.cn/"
          target="_blank"
          rel="noopener noreferrer"
          class="app-icp__link"
        >
          鄂ICP备2024057031号
        </a>
      </footer>
  </n-message-provider>
  </div>
</template>

<style scoped>
.app-shell {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  min-height: 100dvh;
  text-align: left;
}

.app-main {
  flex: 1 0 auto;
}

:deep(.app-body) {
  display: flex;
  flex: 1 1 auto;
  flex-direction: column;
  min-height: 0;
  width: 100%;
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
</style>
