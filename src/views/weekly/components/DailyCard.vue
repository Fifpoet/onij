<script setup lang="ts">

import {h} from 'vue'
import {Icon} from '@iconify/vue';

import {NButton, NCard, NIcon, NList, NListItem, NSpace, NTag, NThing} from "naive-ui";
import {calculateTodayIndex} from "@/util/day.ts";
import {useAppStore} from "@/store/app.ts";


const todayIndex = calculateTodayIndex();
const appStatus = useAppStore()

const props = defineProps({
  weekName: String,
  mouthDay: String,
  index: Number,
})

const dayTimePassed = (() => {
  const now = new Date();

  const totalMinutes = 24 * 60; // 一天的总分钟数
  const currentMinutes = now.getHours() * 60 + now.getMinutes(); // 从午夜开始到现在的分钟数
  return Math.floor(((currentMinutes / totalMinutes) * 10000)) / 100;
})();


const renderIcon = (icon: string) => {
  return () => {
    return h(Icon, {
      icon: icon,
      class: 'text-xl'
    }, {
      default: () => h(icon)
    })
  }
}

const handleClickTodo = (todo: any) => {

}


// weekly todo列表的业务属性： 
// 1. 可以查看弹框详情，具有基本的check功能， 同步状态。 
// 2. 食谱功能跳转 下厨房链接； 绑定食材和属性，添加时筛选； 周结果可以发b站
// 3. 健身功能 ； 分类， tt绑定视频url
// 4. 音乐功能； 笛子 琴， 绑定谱子； 歌唱
// 5. 语言：routine之一
// 6. 读书：
const todoList = [
  {
    title: '今天要写代码',
    content: ['今天要写代码', 'kwhald'],
    link: 'www.xcf.com',
    record: 'feiuu.com',
    tag: [
      {
        id: 1,
        name: '冷冻鸡腿',
        type: 2,
      },
      {
        id: 2,
        name: '洋葱',
        type: 2,
      }
    ],
    checked: true
  },
  {
    title: '鸡腿',
    content: ['今天要写代码', 'kwhald'],
    link: 'www.xcf.com',
    record: 'feiuu.com',
    tag: [
      {
        id: 1,
        name: '冷冻鸡腿',
        type: 2,
      },
      {
        id: 2,
        name: '洋葱',
        type: 2,
      }
    ],
    checked: false
  },

]


</script>

<template>

  <n-card :title="weekName" class="w-72 mx-4 my-2"
          :class="{ 'today-style': todayIndex === index, 'passed': todayIndex > index }">
    <template #header-extra>
      {{ mouthDay }}
    </template>
    <!-- 卡片内容 -->
    <n-list hoverable clickable>
      <!-- todo项目 -->
      <div v-for="todo in todoList">
        <n-list-item :class="{ 'todo-checked': todo.checked}">

          <n-thing :title="todo.title">

            <template #description>
              <n-space size="small">

                <!-- 可见的标签 -->
                <div v-for="tag in todo.tag">
                  <n-tag :bordered="false" type="info" size="small">
                    {{ tag.name }}
                  </n-tag>
                </div>

                最新的打印机<br>
                复制着彩色傀儡<br>
              </n-space>
            </template>
          </n-thing>
        </n-list-item>
      </div>
    </n-list>
    <!-- 底部按钮栏 -->
    <template #action>
      <n-button strong secondary circle type="info" @click="appStatus.changeWeeklyModal(index)">
        <template #icon>
          <n-icon :component="renderIcon('mdi:archive')"/>
        </template>
      </n-button>
    </template>
  </n-card>
</template>

<style scoped>
.today-style {
  /* 今天的样式 */
  border: 2px solid red;
}

.passed {
  /* 今天的样式 */
  border: 2px solid blue;
}

.todo-checked {
  text-decoration: line-through;
  opacity: 0.5;
}
</style>