<template>
  <div @scroll="handleScroll"
       class="music-list absolute bg-white shadow-lg rounded-lg p-4 w-[400px] bottom-[70px] left-0 max-h-[300px] overflow-y-auto">
    <n-list hoverable clickable>
      <div v-for="music in musicStore.musicList" :key="music?.id" class="mb-0.5 cursor-pointer"
           @dblclick="playMusic(music?.id)">
        <n-list-item>
          <n-thing :title="music?.name" content-style="margin-top: 10px;">
            <template #description>
              <n-space size="small" style="margin-top: 2px">
                <n-tag :bordered="false" type="info" size="small">
                  暑夜
                </n-tag>
                <n-tag :bordered="false" type="info" size="small">
                  晚春
                </n-tag>
              </n-space>
            </template>
            <template #header-extra>
              {{music?.artist }}
            </template>
          </n-thing>
        </n-list-item>
      </div>
    </n-list>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { NList, NListItem, NThing, NSpace, NTag } from 'naive-ui';
import { useMusicStore } from "@/store/music.ts";
import { GetMusicDetail, GetMusicList } from '@/api';
import { MusicSortType } from '@/api/types';

// 直接使用 musicStore 而不通过 props 传递
const musicStore = useMusicStore();

// 定义事件
const emit = defineEmits(['play-music']);

// 防止重复请求
const fetching = ref(false);

// 播放音乐
const playMusic = async (id: number) => {
  emit('play-music', id);
};

// 处理滚动事件
const handleScroll = (event: Event) => {
  const target = event.target as HTMLElement;
  // 检查是否滚动到底部
  if (target.scrollHeight - target.scrollTop <= target.clientHeight + 10) {
    if (!fetching.value && musicStore.musicListContinue) {
      fetching.value = true; // 设置为正在加载状态
      fetchMusicList().finally(() => {
        fetching.value = false; // 重置加载状态
      });
    }
  }
};

// 获取音乐列表
const fetchMusicList = async () => {
  try {
    const response = await GetMusicList({
      sort_type: MusicSortType.MST_Created_At_Desc,
      page: musicStore.musicListPage,
      limit: 20
    });
    if (response?.musics) {
      musicStore.append(response.musics);
      musicStore.incrPage()
    } else {
      musicStore.stopPage()
    }
  } catch (error) {
    console.error("加载音乐列表失败:", error);
  }
};
</script>