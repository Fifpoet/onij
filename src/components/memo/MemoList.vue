<template>
  <Transition name="fade">
    <div v-if="musicStore.midShowWhat == MidShowWhat.ShowMemoList" class="flex justify-center items-center min-h-screen w-full p-5 box-border">
    <div class="grid grid-cols-[repeat(auto-fill,minmax(250px,1fr))] gap-5 max-w-1200px w-full justify-center">
      <!-- 示例回忆卡片 -->
      <div class="bg-gray-50 border border-gray-200 rounded-b-lg shadow-md text-center overflow-hidden">
        <img src="https://assets.msn.cn/weathermapdata/1/static/weather/Icons/taskbar_v10/Condition_Card/PartlyCloudyNightV2.svg" alt="Memo Image" class="w-full h-32 object-cover rounded-t-lg">
        <div class="p-5">
          <h3 class="mt-0 text-gray-800">回忆标题 1</h3>
          <p class="text-gray-600">这是第一张回忆卡片的内容。</p>
        </div>
      </div>
      <div class="bg-gray-50 border border-gray-200 rounded-b-lg shadow-md text-center overflow-hidden">
        <img src="https://assets.msn.cn/weathermapdata/1/static/weather/Icons/taskbar_v10/Condition_Card/PartlyCloudyNightV2.svg" alt="Memo Image" class="w-full h-32 object-cover rounded-t-lg">
        <div class="p-5">
          <h3 class="mt-0 text-gray-800">回忆标题 2</h3>
          <p class="text-gray-600">这是第二张回忆卡片的内容。</p>
        </div>
      </div>
      <div class="bg-gray-50 border border-gray-200 rounded-b-lg shadow-md text-center overflow-hidden">
        <img src="https://assets.msn.cn/weathermapdata/1/static/weather/Icons/taskbar_v10/Condition_Card/PartlyCloudyNightV2.svg" alt="Memo Image" class="w-full h-32 object-cover rounded-t-lg">
        <div class="p-5">
          <h3 class="mt-0 text-gray-800">回忆标题 3</h3>
          <p class="text-gray-600">这是第三张回忆卡片的内容。</p>
        </div>
      </div>
    </div>
  </div>
  </Transition>
</template>

<script setup lang="ts">
import { MidShowWhat } from "@/api/types";
import { useMusicStore } from '@/store/music';

const musicStore = useMusicStore();
// 你可以在这里添加组件的逻辑，例如从API获取回忆数据
</script>

<style scoped>
/* UnoCSS 通常不需要大量的 scoped CSS，样式通过类名直接应用 */
</style>