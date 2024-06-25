<script setup lang="ts">

import {NCard} from "naive-ui";
import {defineProps} from "vue";

const props = defineProps({
  weekName: String,
  mouthDay: String,
  isToday: Boolean,
})

</script>

<template>
  <n-card :title="weekName" class="w-72 mx-4 my-2" :class="{ 'today-style': isToday }">
    <template #header-extra>
      {{mouthDay}}
    </template>
    卡片内容
    <template #footer>
      #footer
    </template>
    <template #action>
      #action
    </template>
  </n-card>
</template>

<style scoped>
.today-style {
  /* 今天的样式 */
  border: 2px solid red;
}

</style>