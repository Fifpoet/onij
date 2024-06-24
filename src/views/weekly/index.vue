<script setup lang="ts">


import DailyCard from "@/views/weekly/components/DailyCard.vue";
</script>

<template>
  <div class="flex flex-wrap justify-center w-full"> <!--flex-wrap空间不足时换行-->
    <!--    w-full：确保在小屏幕上，每个DailyCard占据全部宽度。md:w-1/4：在中等屏幕尺寸及以上，每个DailyCard占据容器宽度的四分之一。-->
    <div v-for="n in 7" :key="n" class="w-full md:w-1/5 px-2 py-4">
      <DailyCard></DailyCard>
    </div>
  </div>
</template>

<style scoped>

</style>
