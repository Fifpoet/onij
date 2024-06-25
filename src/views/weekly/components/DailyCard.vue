<script setup lang="ts">

import {NCard} from "naive-ui";
</script>

<template>
  <n-card title="卡片插槽示例" class="w-60 mx-4 my-2">
    <template #header-extra>
      #header-extra
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

</style>