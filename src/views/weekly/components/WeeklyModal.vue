
<template>
  <n-modal v-model:show="appStatus.showWeeklyModal">
    <n-card
        style="width: 600px"
        title="模态框"
        :bordered="true"
        size="huge"
        role="dialog"
        aria-modal="true"
    >
      <template #header-extra>
        噢！
      </template>
      内容
      <template #footer>
        <n-button @click="appStatus.changeWeeklyModal()" type="info">
          来吧
        </n-button>
      </template>
    </n-card>
  </n-modal>
</template>



<script setup lang="ts">
import { useAppStore } from '@/store';
import {NModal, NCard} from 'naive-ui'
import { ref } from 'vue'

  const appStatus = useAppStore();
</script>