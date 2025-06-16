<template>
  <div class="min-h-screen bg-gray-50">
    <n-page-header title="文件管理" subtitle="查看和管理您的文件">
      <template #extra>
        <n-space>
          <n-button type="primary" @click="refreshList">
            <template #icon>
              <n-icon>
                <RefreshOutline />
              </n-icon>
            </template>
            刷新
          </n-button>
          <n-button @click="goBack" v-if="currentParentId !== 0">
            <template #icon>
              <n-icon>
                <ArrowBackOutline />
              </n-icon>
            </template>
            返回上级
          </n-button>
        </n-space>
      </template>
    </n-page-header>
    
    <FileList 
      ref="fileListRef"
      :parent-id="currentParentId" 
      @file-click="handleFileClick" 
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { NPageHeader, NButton, NSpace, NIcon } from 'naive-ui';
import { RefreshOutline, ArrowBackOutline } from '@vicons/ionicons5';
import FileList from './FileList.vue';
import { FileType, FileDetail } from '@/api/types/file';

// 当前父级文件夹ID
const currentParentId = ref(0);
const fileListRef = ref<InstanceType<typeof FileList>>();

// 刷新列表
const refreshList = () => {
  fileListRef.value?.refresh();
};

// 返回上级
const goBack = () => {
  // TODO: 实现返回上级逻辑
  currentParentId.value = 0;
};

// 处理文件点击
const handleFileClick = (file: FileDetail) => {
  if (file.format === FileType.FT_Folder) {
    // 如果是文件夹，进入该文件夹
    currentParentId.value = file.id;
  }
};
</script> 