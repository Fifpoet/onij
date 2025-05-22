<script setup lang="ts">
import { RouterView } from 'vue-router';
import PageHeader from './components/layout/PageHeader.vue'
// import lbAudio from './components/audio/index.vue';
import AudioContainer from "@/components/music/AudioContainer.vue";
import MusicDetail from "@/components/music/MusicDetail.vue";
import {NMessageProvider} from "naive-ui"

</script>

<template>
  <div class="h-full w-full flex flex-col">

    <!-- 顶部导航栏 -->
  <PageHeader />
  </div>

<!--  <article class="flex flex-1 flex-col">-->
<!--      <RouterView v-slot="{ Component, route }">-->
<!--        <component :is="Component" :key="route.path" />-->
<!--      </RouterView>-->
<!--    </article>-->

  <n-message-provider>
  <MusicDetail />
  <AudioContainer></AudioContainer>
  <!-- <img src="https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=gQGY8DwAAAAAAAAAAS5odHRwOi8vd2VpeGluLnFxLmNvbS9xLzAyRVMydjVzLXBmeEgxYlNxSnhFMTIAAgTmjC1oAwQQDgAA" alt="QR Code" class="mx-auto my-2"> -->
<!--  <lbAudio class="music-player" :musicList="musicList" :index="8" :lyrics="true">-->
<!--  </lbAudio>-->
<!--  <WeeklyModal></WeeklyModal>-->

  </n-message-provider>
  <br>
  <br>
  <br>
  <br>
  <br>
  <br>

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
