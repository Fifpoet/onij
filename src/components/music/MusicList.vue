<template>
  <div @scroll="handleScroll"
    class="music-list absolute bg-white shadow-lg rounded-lg p-4 w-[400px] bottom-[70px] left-0 max-h-[300px] overflow-y-auto text-sm">
    <!-- 新增选择器 -->
    <div class="mt-4">
      <n-select multiple v-model:value="selectedTag" :options="tagOptions" :render-label="renderLabel" tag filterable
        placeholder="选择标签" />
    </div>
    <n-list hoverable clickable>
      <div v-for="music in musicStore.musicList" :key="music?.id" class="mb-0.5 cursor-pointer"
        @dblclick="playMusic(music?.id)">
        <n-list-item>
          <n-thing :title="music?.name" content-style="margin-top: 1px;">
            <template #description>
              <n-space size="small" style="margin-top: 1px">
                <n-tag :bordered="false" type="info" size="small">
                  暑夜
                </n-tag>
                <n-tag :bordered="false" type="info" size="small">
                  晚春
                </n-tag>
              </n-space>
            </template>
            <template #header-extra>
              {{ music?.artist }}
            </template>
          </n-thing>
        </n-list-item>
      </div>
    </n-list>

  </div>

</template>

<script lang="ts" setup>
import { ref, h, onMounted } from 'vue';
import { NList, NListItem, NThing, NSpace, NTag, NSelect, NIcon, SelectGroupOption } from 'naive-ui';
import { useMusicStore } from "@/store/music.ts";
import { GetMusicDetail, GetMusicList, SearchArtist } from '@/api';
import { MusicSortType, tagOpts, TagType } from '@/api/types';
import { MusicalNote as MusicIcon } from '@vicons/ionicons5';
import type { SelectOption } from 'naive-ui';

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

// 新增选择器相关代码
const selectedTag = ref(null);

const tagOptions = ref<SelectGroupOption[]>([])

const renderLabel = (option: SelectOption) => {
  if (option.type === 'group') {
    return option.label;
  }
  return [
    h(
      NIcon,
      {
        style: {
          verticalAlign: '-0.15em',
          marginRight: '4px'
        }
      },
      {
        default: () => h(MusicIcon)
      }
    ),
    option.label as string
  ];
};

onMounted(async () => {
  // 获取标签
  try {
    const response = await SearchArtist({
      keyword: '',
      tag_type: TagType.TT_Star,
      page: 1,
      limit: 50
    });
    if (response?.artists) {
      // 艺术家tag
      tagOptions.value.push({
        type: 'group',
        label: '艺术家',
        key: 'artist',
        children: response.artists.map(artist => ({
          label: artist.name,
          value: String(artist.id),
          type: 'success'
        }))
      });
    console.log(tagOptions.value)
    }
  } catch (error) {
    console.error("加载艺术家列表失败:", error);
  }

  // 固定tag
  tagOptions.value.push(...tagOpts)
});
</script>