<template>
  <Transition name="fade">
    <div v-if="musicStore.midShowWhat == MidShowWhat.ShowMusicLike" class="music-like-container p-8">
      <div class="music-like-content w-full max-w-4xl mx-auto">
        <h2 class="text-2xl font-bold mb-6 text-center">加星收藏</h2>

        <!-- 歌手列表 -->
        <div class="mb-6">
          <h3 class="text-lg font-semibold mb-3">歌手</h3>
          <div class="flex flex-wrap gap-2">
            <div
              v-for="artist in musicStore.current.detail?.artist_names"
              :key="artist"
              class="flex items-center bg-gray-100 px-3 py-1 rounded-full text-sm"
            >
              <span>{{ artist }}</span>
              <NIcon :component="HeartOutline" class="ml-2 cursor-pointer hover:text-red-500" @click="toggleArtistLike(artist)" />
            </div>
          </div>
        </div>

        <!-- 标签组 -->
        <div v-for="group in tagGroups" :key="group.key" class="mb-6">
          <h3 class="text-lg font-semibold mb-3">{{ group.label }}</h3>
          <div class="flex flex-wrap gap-2">
            <div
              v-for="tag in group.children"
              :key="tag.value"
              class="flex items-center bg-gray-100 px-3 py-1 rounded-full text-sm"
            >
              <span>{{ tag.label }}</span>
              <NIcon :component="HeartOutline" class="ml-2 cursor-pointer hover:text-red-500" @click="toggleTagLike(tag)" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { MidShowWhat } from "@/api/types";
import { ref, computed } from 'vue';
import { useMusicStore } from '@/store/music';
import { tagOpts } from '@/api/types';
import { TagBiz, TagGroup, TagType } from '@/api/types/tag';
import { NButton, NIcon } from 'naive-ui';
import { HeartOutline } from '@vicons/ionicons5';

const musicStore = useMusicStore();

// 过滤出 tagOpts 中 type 为 'group' 的项，并确保 children 存在
const tagGroups = computed(() => {
  return tagOpts.filter(opt => opt.type === 'group' && opt.children);
});

const toggleArtistLike = async (artistName: string) => {
  console.log(`Toggling like for artist: ${artistName}`);
  // 这里需要根据实际情况获取 artist_id，目前 MusicDetail 中只有 artist_ids 和 artist_names
  // 如果需要精确到某个 artist_id，可能需要修改 MusicDetail 接口或通过其他方式获取
  // 暂时使用一个示例 TagDetail 结构
  const tagDetail = {
    resource_id: musicStore.current.detail?.id || 0, // 音乐ID
    resource_type: 1, // 假设1代表音乐
    tag_biz: TagBiz.TB_Music,
    tag_group: TagGroup.TG_ArtistStar, // 假设这是歌手收藏的标签组
    tag_type: TagType.TT_Star, // 假设这是收藏的标签类型
    extra: JSON.stringify({ artist_name: artistName }), // 存储歌手名称
    list_show: true,
  };

  try {
    const resp = await UploadTag({ tag_detail: tagDetail });
    if (resp.code === 0) {
      window.$message.success(`成功收藏歌手: ${artistName}`);
    } else {
      window.$message.error(`收藏歌手失败: ${resp.message}`);
    }
  } catch (error) {
    console.error('收藏歌手出错:', error);
    window.$message.error('收藏歌手出错');
  }
};

const toggleTagLike = async (tag: any) => {
  console.log(`Toggling like for tag: ${tag.label}`);
  let parsedValue: any = {};
  try {
    parsedValue = JSON.parse(tag.value);
  } catch (e) {
    console.error('Failed to parse tag value:', tag.value, e);
    window.$message.error('标签数据解析失败');
    return;
  }

  const tagDetail = {
    resource_id: musicStore.current.detail?.id || 0, // 音乐ID
    resource_type: 1, // 假设1代表音乐
    tag_biz: TagBiz.TB_Music,
    tag_group: getTagGroupFromTagType(parsedValue.tag_type), // 根据tag_type获取tag_group
    tag_type: parsedValue.tag_type, // 使用标签中定义的tag_type
    extra: JSON.stringify({ label: tag.label }), // 存储标签名称
    list_show: true,
  };

  try {
    const resp = await UploadTag({ tag_detail: tagDetail });
    if (resp.code === 0) {
      window.$message.success(`成功收藏标签: ${tag.label}`);
    } else {
      window.$message.error(`收藏标签失败: ${resp.message}`);
    }
  } catch (error) {
    console.error('收藏标签出错:', error);
    window.$message.error('收藏标签出错');
  }
};

// 辅助函数：根据 TagType 获取 TagGroup
// 这是一个简化的映射，您可能需要根据实际的 TagType 范围进行更精确的映射
const getTagGroupFromTagType = (tagType: TagType): TagGroup => {
  if (tagType >= TagType.TT_NcBridge && tagType < TagType.TT_Star) {
    return TagGroup.TG_MusicHighlight; // 示例：假设 TT_NcBridge 属于高亮
  } else if (tagType === TagType.TT_Star) {
    return TagGroup.TG_ArtistStar; // 示例：假设 TT_Star 属于歌手收藏
  }
  // 默认返回未知或根据您的业务逻辑返回其他组
  return TagGroup.TG_Unknown;
};

// 移除 handleClose 和 emit('close')
// const emit = defineEmits(['close']);

// const handleClose = () => {
//   emit('close');
// };
</script>

<style scoped>
.music-like-container {
  /* UnoCSS 已经处理了大部分样式，这里可以添加一些额外的 */
}

.music-like-content {
  /* UnoCSS 已经处理了大部分样式，这里可以添加一些额外的 */
}

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