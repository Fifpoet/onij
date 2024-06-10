<script setup lang="ts">
// import { ref } from 'vue'
import { Icon } from '@iconify/vue';

const menuOptions = [

]

</script>


<template>
  <!-- 移动端顶部导航栏 -->
  <Transition name="slide-fade" appear>
    <div class="fixed inset-x-0 top-0 z-11 h-[60px] flex items-center justify-between px-4 py-2 lg:hidden">
      <!-- 这里的lg:hidden是只在移动端展示 -->
      <!-- 左上角标题 -->
      <RouterLink to="/" class="text-[18px] font-bold">
        TITLE
      </RouterLink>
      <!-- 右上角图标 -->
      <div class="flex items-center gap-2 text-2xl">
        <button>
          <Icon icon="ic:round-search" />
        </button>
        <button>
          <Icon icon="ic:sharp-menu" />
        </button>
      </div>
    </div>
  </Transition>

  <!-- 电脑端顶部导航栏 -->
  <Transition name="slide-fade" appear>
    <div class="fixed inset-x-0 top-0 z-11 hidden h-[60px] lg:block">
      <!-- 顶部左侧图标和菜单 -->
      <div class="flex items-center gap-2 text-2xl">
        <RouterLink to="/" class="text-[18px] font-bold">
          onij
        </RouterLink>
        <button>
          <Icon icon="ic:round-search" />
        </button>
        <button>
          <Icon icon="ic:sharp-menu" />
        </button>
      </div>
      <!-- 顶部右侧菜单 -->

    </div>
  </Transition>
</template>


<style scoped></style>