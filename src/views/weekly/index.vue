<script setup lang="ts">
import { ref } from 'vue';
import DailyCard from "@/views/weekly/components/DailyCard.vue";
import ExtendCard from './components/ExtendCard.vue';
import { getWeekName, getMouthDay, calculateTodayIndex } from '@/util/day';


const today = ref(new Date());
const weekNames = getWeekName(today.value);
const mouthDays = getMouthDay(today.value);


</script>

<template>
<div class="flex flex-wrap">
  <div v-for="(weekName, index) in weekNames" :key="index">
    <DailyCard 
      :weekName="weekName" 
      :mouthDay="mouthDays[index]" 
      :isToday=" index === calculateTodayIndex()"
    ></DailyCard>
  </div>
  
  <ExtendCard></ExtendCard>
</div>
</template>