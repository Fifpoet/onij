<script setup lang="ts">

import { NCard } from "naive-ui";
import { NList } from "naive-ui";
import { NListItem } from "naive-ui";
import { NSpace } from "naive-ui";
import { NTag } from "naive-ui";
import { NThing } from "naive-ui";
import { NProgress } from "naive-ui";

const props = defineProps({
  weekName: String,
  mouthDay: String,
  isToday: Boolean,
})

const dayTimePassed = (() => {
  const now = new Date();

  const totalMinutes = 24 * 60; // 一天的总分钟数
  const currentMinutes = now.getHours() * 60 + now.getMinutes(); // 从午夜开始到现在的分钟数
  return Math.floor(((currentMinutes / totalMinutes)* 10000)) / 100;
})();

</script>

<template>
  <n-card :title="weekName" class="w-72 mx-4 my-2" :class="{ 'today-style': isToday }">
    <template #header-extra>
      {{ mouthDay }}
    </template>
    <!-- 卡片内容 -->
    <n-list hoverable clickable>

      <div v-if="isToday">
        <n-progress style="margin: 0 2px 2px 0; width: 80px; height: 80px;" type="circle" :percentage="dayTimePassed"  :stroke-width="6"/>
        <n-list-item>
          <n-thing title="相见恨晚" content-style="margin-top: 10px;">
            <template #description>
              <n-space size="small" style="margin-top: 4px">
                <n-tag :bordered="false" type="info" size="small">
                  暑夜
                </n-tag>
                <n-tag :bordered="false" type="info" size="small">
                  晚春
                </n-tag>
              </n-space>
            </template>
            奋勇呀然后休息呀<br>
            完成你伟大的人生
          </n-thing>
        </n-list-item>
      </div>
      <n-list-item>
        <n-thing title="他在时间门外" content-style="margin-top: 10px;">
          <template #description>
            <n-space size="small" style="margin-top: 4px">
              <n-tag :bordered="false" type="info" size="small">
                环形公路
              </n-tag>
              <n-tag :bordered="false" type="info" size="small">
                潜水艇司机
              </n-tag>
            </n-space>
          </template>
          最新的打印机<br>
          复制着彩色傀儡<br>
          早上好我的罐头先生<br>
          让他带你去被工厂敲击
        </n-thing>
      </n-list-item>
    </n-list>
    <template #footer>
      #footer
    </template>
    <template #action>
      #action
    </template>
  </n-card>
</template>

<style scoped>
.today-style {
  /* 今天的样式 */
  border: 2px solid red;
}
</style>