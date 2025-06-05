<template>
  <Transition name="fade">
    <div v-if="musicStore.midShowWhat == MidShowWhat.ShowMusicLike" class="music-like-container p-8">
      <div class="music-like-content w-full max-w-4xl mx-auto">
        <!-- 歌手列表 -->
        <div class="mb-6">
          <h3 class="text-lg font-semibold mb-3">歌手</h3>
          <div class="flex flex-wrap gap-2">
            <div v-for="(artistName, index) in musicStore.current.detail?.artist_names" :key="artistName"
              class="flex items-center bg-gray-100 px-3 py-1 rounded-full text-sm">
              <span>{{ artistName }}</span>
              <NIcon :component="musicStore.likeMap[musicStore.current.detail?.artist_ids[index]] ? Heart : HeartOutline" 
              class="ml-2 cursor-pointer hover:text-red-500" @click="toggleArtistLike(musicStore.current.detail?.artist_ids[index], artistName)" />
            </div>
          </div>
        </div>

        <!-- 标签组 -->
        <div v-for="group in tagOpts" :key="String(group.key)" class="mb-6">
          <h3 class="text-lg font-semibold mb-3">{{ group.label }}</h3>
          <div class="flex flex-wrap gap-2">
            <div v-for="tag in group.children" :key="tag.value"
              class="flex items-center bg-gray-100 px-3 py-1 rounded-full text-sm">
              <span>{{ tag.label }}</span>
              <NIcon :component="musicStore.likeMap[tag.label] ? Heart: HeartOutline" class="ml-2 cursor-pointer hover:text-red-500"
                @click="toggleTagLike(group, tag)" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { MidShowWhat } from "@/api/types";
import { useMusicStore } from '@/store/music';
import { tagOpts } from '@/api/types';
import { ResourceType, SelectTagExtra, TagBiz, TagDetail, TagGroup, TagType } from '@/api/types/tag';
import { NButton, NIcon, SelectGroupOption } from 'naive-ui';
import { HeartOutline, Heart } from '@vicons/ionicons5';
import { UploadTag } from "@/api";
import { SelectBaseOption } from "naive-ui/es/select/src/interface";
import { onMounted, ref } from "vue";

const musicStore = useMusicStore();

const toggleTagLike = async ( group: SelectGroupOption, tag: SelectBaseOption) => {
  const extra = JSON.parse(String(tag.value)) as SelectTagExtra;
  const tagDetail: TagDetail = {
    resource_id: musicStore.current.detail?.id,
    resource_type: ResourceType.RT_Music,
    tag_biz: TagBiz.TB_Music,
    tag_group: group.key as TagGroup,
    tag_type: extra.tag_type!,
    target_id: 0,
    target_type: 0,
    extra: "",
    list_show: true,
  }
  try {
    const response = await UploadTag({ tag_detail: tagDetail });
  } catch (error) {
    console.error('收藏失败', error);
  }
};

const toggleArtistLike = async (artistId: number | undefined, artistName: string) => {
  if (artistId === undefined) {
    console.error('Artist ID is undefined, cannot like artist.');
    return;
  }

  const tagDetail: TagDetail = {
    resource_id: artistId,
    resource_type: ResourceType.RT_Artist,
    tag_biz: TagBiz.TB_Music, 
    tag_group: TagGroup.TG_ArtistStar,
    tag_type: TagType.TT_Star,
    target_id: 0,
    target_type: 0,
    extra: "",
    list_show: true,
  };
  try {
    const response = await UploadTag({ tag_detail: tagDetail });
  } catch (error) {
    console.error(`收藏艺术家 ${artistName} 失败`, error);
  }
};

onMounted(() => {
  setInterval (() => {
    console.log('111'); // 3000ms 后输出
    console.log(musicStore.likeMap); // 3000ms 后输出
  }, 3000);
});

</script>

<style scoped>
/* 添加过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: all 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}
</style>