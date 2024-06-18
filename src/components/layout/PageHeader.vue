<script setup lang="ts">
// import { ref } from 'vue'
import { Icon } from '@iconify/vue';
import { useAppStore, useUserStore } from '@/store';

const appStore = useAppStore()
const userStore = useUserStore()

const barShow = true

const menuOptions = [
  { text: '首页', icon: 'mdi:home', path: '/' },
  {
    text: '发现',
    icon: 'mdi:apple-safari',
    subMenu: [
      { text: '归档', icon: 'mdi:archive', path: '/archives' },
      { text: '分类', icon: 'mdi:menu', path: '/categories' },
      { text: '标签', icon: 'mdi:tag', path: '/tags' },
    ],
  }
]

</script>


<template>
  <Transition name="slide-fade" appear>
    <div v-if="barShow" :class="navClass" class="fixed inset-x-0 top-0 z-11 hidden h-[60px] lg:block">
      <div class="h-full flex items-center justify-between px-9">
        <!-- 左上角标题 -->
        <RouterLink to="/" class="text-xl font-bold">
          {{ appStore.blogConfig.website_author }}
        </RouterLink>
        <!-- 右上角菜单 -->
        <div class="flex items-center space-x-4">
          <!-- 搜索 -->
          <div class="menus-item">
            <a class="menu-btn flex items-center" @click="appStore.setSearchFlag(true)">
              <Icon icon="mdi:magnify" class="text-xl" />
              <span class="ml-1"> 搜索 </span>
            </a>
          </div>
          <div v-for="item of menuOptions" :key="item.text" class="menus-item">
            <!-- 不包含子菜单 -->
            <RouterLink v-if="!item.subMenu" :to="item.path" class="menu-btn flex items-center">
              <Icon :icon="item.icon" class="text-xl" />
              <span class="ml-1"> {{ item.text }} </span>
            </RouterLink>
            <!-- 包含子菜单 -->
            <div v-else class="menu-btn">
              <div class="flex items-center">
                <Icon :icon="item.icon" class="text-xl" />
                <span class="mx-1"> {{ item.text }} </span>
                <Icon icon="ep:arrow-down-bold" class="text-xl" />
              </div>
              <ul class="menus-submenu">
                <RouterLink v-for="sub of item.subMenu" :key="sub.text" :to="sub.path">
                  <div class="flex items-center">
                    <Icon :icon="sub.icon" class="text-xl" />
                    <span class="ml-1"> {{ sub.text }} </span>
                  </div>
                </RouterLink>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>


<style scoped></style>