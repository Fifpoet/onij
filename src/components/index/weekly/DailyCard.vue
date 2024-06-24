<script setup lang="ts">
import NCard from "naive-ui";
import NButton from "naive-ui";
import {onMounted, ref} from "vue";

let initDone = ref(true);
onMounted(() => {
})

</script>

<template>
<!--  <n-flex>-->
<!--    <div v-for="n in 7" :key="n">-->
<!--      <n-card title="周一" hoverable>-->
<!--        <template #header-extra>-->
<!--          5月20日-->
<!--        </template>-->
<!--        卡片内容-->
<!--        <template #footer>-->
<!--          #footer-->
<!--        </template>-->
<!--        <template #action>-->
<!--          一些按钮-->
<!--        </template>-->
<!--      </n-card>-->
<!--    </div>-->
<!--  </n-flex>-->
  <div v-if="initDone">
  <h1>good</h1>
  <n-button>Default</n-button>
  </div>
</template>

<style scoped>

</style>
