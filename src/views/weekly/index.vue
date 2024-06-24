<script setup lang="ts">


import DailyCard from "@/views/weekly/components/DailyCard.vue";
</script>

<template>
  <div class="flex flex-wrap justify-start">
    <div v-for="n in 7">
    <DailyCard></DailyCard>
    </div>
  </div>
</template>

<style scoped>

</style>
