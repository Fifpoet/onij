<template>
  <PageContent>
    <div class="cloud-page">
      <h1 class="browse-section-title mb-4 border-l-4 border-blue-500 pl-2">云盘</h1>
      <FileList />
    </div>
  </PageContent>
</template>

<script setup lang="ts">
import PageContent from '@/components/layout/PageContent.vue'
import FileList from '@/components/file/FileList.vue'
</script>

<style scoped>
.cloud-page {
  width: 100%;
}
</style>
