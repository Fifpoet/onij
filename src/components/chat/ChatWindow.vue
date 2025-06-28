<template>
  <Transition name="fade-scale" mode="out-in">
    <div 
      v-if="musicStore.midShowWhat === MidShowWhat.ShowChatWindow" 
      class="fixed inset-x-0 top-1/2 -translate-y-1/2 bg-white dark:bg-dark-800 z-0 hidden lg:block"
    >
      <div class="max-w-4xl mx-auto p-6">
        <McLayout class="flex flex-col h-[500px] rounded-lg border dark:border-gray-700">
          <!-- 聊天内容区域 -->
          <McLayoutContent class="flex-1 overflow-auto p-4 space-y-4">
            <template v-for="(msg, idx) in messages" :key="idx">
              <McBubble
                v-if="msg.from === 'user'"
                :content="msg.content"
                :align="'right'"
                :variant="'filled'"
                :avatarConfig="{
                  imgSrc: '/assets/avatar/user.svg',
                  width: 40,
                  height: 40,
                  isRound: true
                }"
              />
              <McBubble 
                v-else 
                :content="msg.content"
                :align="'left'"
                :variant="'filled'"
                :avatarConfig="{
                  imgSrc: '/assets/avatar/ai.svg',
                  width: 40,
                  height: 40,
                  isRound: true
                }"
              />
            </template>
          </McLayoutContent>

          <!-- 输入区域 -->
          <McLayoutSender class="border-t dark:border-gray-700 p-4">
            <McInput 
              v-model:value="inputValue"
              :maxLength="20000"
              :loading="loading"
              @submit="onSubmit"
              placeholder="输入您的问题..."
              class="w-full"
            >
              <template #head>
                <div v-if="attachments.length" class="flex gap-2 p-2">
                  <div 
                    v-for="(file, index) in attachments" 
                    :key="file.id"
                    class="px-2 py-1 bg-gray-100 dark:bg-dark-700 rounded flex items-center gap-2 text-sm"
                  >
                    <span class="truncate max-w-[200px]">{{ file.name }}</span>
                    <button 
                      class="text-gray-500 hover:text-red-500 transition-colors"
                      @click="removeAttachment(index)"
                    >
                      <i class="icon-code-editor-close"></i>
                    </button>
                  </div>
                </div>
              </template>
              <template #extra>
                <div class="flex justify-between items-center w-full px-2">
                  <div class="flex items-center gap-4 text-sm text-gray-600">
                    <div class="flex items-center gap-1 cursor-pointer hover:text-primary transition-colors">
                      <input 
                        type="file" 
                        ref="fileInput"
                        class="hidden"
                        @change="handleFileUpload"
                      >
                      <i class="icon-appendix"></i>
                      <span @click="fileInput?.click()">附件</span>
                    </div>
                    <span class="text-gray-500">{{ inputValue.length }}/20000</span>
                  </div>
                  <div class="flex gap-2">
                    <button 
                      class="px-3 py-1 rounded text-sm hover:bg-gray-100 transition-colors"
                      :class="[inputValue ? 'text-primary' : 'text-gray-400']"
                      @click="inputValue = ''"
                      :disabled="!inputValue"
                    >
                      清空
                    </button>
                  </div>
                </div>
              </template>
            </McInput>
          </McLayoutSender>
        </McLayout>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref, nextTick } from 'vue';
import { useMusicStore } from "@/store/music";
import { MidShowWhat } from "@/api/types";

interface Attachment {
  id: number;
  name: string;
  url: string;
}

const musicStore = useMusicStore();
const inputValue = ref('');
const loading = ref(false);
const fileInput = ref<HTMLInputElement | null>(null);
const attachments = ref<Attachment[]>([]);
const messages = ref<Array<{from: 'user' | 'model', content: string}>>([]);

// 模拟文件上传
const uploadFile = async (file: File): Promise<string> => {
  return new Promise((resolve) => {
    setTimeout(() => {
      // 模拟返回文件URL
      resolve(`https://example.com/files/${file.name}`);
    }, 1000);
  });
};

const handleFileUpload = async (event: Event) => {
  const input = event.target as HTMLInputElement;
  if (!input.files?.length) return;

  const file = input.files[0];
  loading.value = true;

  try {
    const url = await uploadFile(file);
    attachments.value.push({
      id: Date.now(),
      name: file.name,
      url: url
    });
  } catch (error) {
    console.error('文件上传失败:', error);
  } finally {
    loading.value = false;
    // 清空 input 值，这样相同文件可以重复上传
    input.value = '';
  }
};

const removeAttachment = (index: number) => {
  attachments.value.splice(index, 1);
};

const onSubmit = (content: string) => {
  if (!content.trim() && !attachments.value.length) return;
  
  // 构建消息内容
  let messageContent = content.trim();
  if (attachments.value.length) {
    messageContent += '\n附件：\n' + attachments.value.map(file => `- ${file.name}: ${file.url}`).join('\n');
  }

  // 添加用户消息
  messages.value.push({
    from: 'user',
    content: messageContent
  });

  // 清空输入和附件
  inputValue.value = '';
  attachments.value = [];

  // 模拟AI回复
  loading.value = true;
  setTimeout(() => {
    messages.value.push({
      from: 'model',
      content: `收到您的消息：${messageContent}`
    });

    loading.value = false;

    // 滚动到底部
    nextTick(() => {
      const container = document.querySelector('.mc-layout-content');
      if (container) {
        container.scrollTop = container.scrollHeight;
      }
    });
  }, 500);
};
</script>

<style scoped>
.fade-scale-enter-active,
.fade-scale-leave-active {
  transition: all 0.3s ease-out;
}

.fade-scale-enter-from,
.fade-scale-leave-to {
  opacity: 0;
  transform: scale(0.9);
}

:deep(.mc-layout-content) {
  scrollbar-width: thin;
  scrollbar-color: rgba(156, 163, 175, 0.5) transparent;
}

:deep(.mc-layout-content::-webkit-scrollbar) {
  width: 6px;
}

:deep(.mc-layout-content::-webkit-scrollbar-track) {
  background: transparent;
}

:deep(.mc-layout-content::-webkit-scrollbar-thumb) {
  background-color: rgba(156, 163, 175, 0.5);
  border-radius: 3px;
}

:deep(.mc-bubble-content) {
  max-width: 80%;
  word-break: break-word;
  white-space: pre-wrap;
}
</style>
