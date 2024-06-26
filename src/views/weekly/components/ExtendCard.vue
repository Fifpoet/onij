<script setup lang="ts">

import {NCard} from "naive-ui";
import {NSwitch} from "naive-ui";
</script>

<template>
  <n-card title="卡片插槽示例" class="w-72 mx-4 my-2">
    <template #header-extra>
      #header-extra
    </template>
    <n-switch :value="active" size="large">
      <template #icon>
        🥵
      </template>
    </n-switch>
    <n-switch :value="active" size="large">
      <template #icon>
        🤤
      </template>
    </n-switch>
        <n-switch :value="active" size="large">
      <template #icon>
        🥰
      </template>
    </n-switch>
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