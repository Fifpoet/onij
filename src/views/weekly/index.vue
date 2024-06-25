<script setup lang="ts">
import {NButton} from 'naive-ui'
import {NCard} from 'naive-ui'
import DailyCard from "@/views/weekly/components/DailyCard.vue";
</script>

<template>
  <div class="flex">
  <div v-for="n in 7">
    <DailyCard></DailyCard>
  </div>
  </div>
</template>