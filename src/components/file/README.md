# 文件列表组件

## 功能特性

- 📁 支持多种文件类型图标显示（文件夹、MP3、歌词、图片等）
- 📊 显示文件大小（目前写死，后续接口返回后可替换）
- 🕒 显示文件创建时间
- ⬇️ 文件下载功能
- 🗑️ 文件删除功能（待实现）
- 📱 响应式设计，适配不同屏幕尺寸
- 🔄 支持刷新和文件夹导航
- 🎨 使用 UnoCSS 工具类，无原生 CSS

## 使用方法

### 基础用法

```vue
<template>
  <FileList :parent-id="0" @file-click="handleFileClick" />
</template>

<script setup lang="ts">
import FileList from '@/components/file/FileList.vue';
import { FileType, FileDetail } from '@/api/types/file';

const handleFileClick = (file: FileDetail) => {
  if (file.format === FileType.FT_Folder) {
    // 进入文件夹
    console.log('进入文件夹:', file.id);
  }
};
</script>
```

### 带导航的完整示例

```vue
<template>
  <div class="min-h-screen bg-gray-50">
    <n-page-header title="文件管理">
      <template #extra>
        <n-button @click="refreshList">刷新</n-button>
        <n-button @click="goBack" v-if="currentParentId !== 0">返回上级</n-button>
      </template>
    </n-page-header>
    
    <FileList 
      ref="fileListRef"
      :parent-id="currentParentId" 
      @file-click="handleFileClick" 
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import FileList from '@/components/file/FileList.vue';
import { FileType, FileDetail } from '@/api/types/file';

const currentParentId = ref(0);
const fileListRef = ref<InstanceType<typeof FileList>>();

const refreshList = () => {
  fileListRef.value?.refresh();
};

const goBack = () => {
  currentParentId.value = 0;
};

const handleFileClick = (file: FileDetail) => {
  if (file.format === FileType.FT_Folder) {
    currentParentId.value = file.id;
  }
};
</script>
```

## Props

| 属性名 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| parentId | number | 0 | 父级文件夹ID，0表示根目录 |

## Events

| 事件名 | 参数 | 说明 |
|--------|------|------|
| fileClick | file: FileDetail | 文件被点击时触发 |

## 文件类型支持

- `FT_Folder` - 文件夹
- `FT_Mp3` - MP3音频文件
- `FT_Lyrics` - 歌词文件
- `FT_Png` - PNG图片文件
- `FT_Unknown` - 未知文件类型

## 样式说明

组件完全使用 UnoCSS 工具类，无需编写原生 CSS：

- `p-5` - 内边距
- `w-full max-w-4xl mx-auto` - 宽度和居中
- `border-b border-gray-100` - 底部边框
- `cursor-pointer transition-colors duration-200` - 鼠标指针和过渡效果
- `hover:bg-gray-50` - 悬停背景色
- `flex gap-4 text-gray-500 text-xs` - 弹性布局和文字样式
- `min-w-15 min-w-30` - 最小宽度
- `flex justify-center py-5` - 居中对齐和内边距

## 注意事项

1. 文件大小目前是写死的，需要等后端接口返回实际大小后更新
2. 删除文件功能需要后端提供相应的API接口
3. 组件使用了 Naive UI 的组件库，确保项目中已正确安装
4. 图标使用了 @vicons/ionicons5，确保已安装该依赖
5. 样式完全基于 UnoCSS 工具类，无需额外 CSS 文件

## 技术栈

- Vue 3 + TypeScript
- Naive UI 组件库
- UnoCSS 工具类框架
- @vicons/ionicons5 图标库 