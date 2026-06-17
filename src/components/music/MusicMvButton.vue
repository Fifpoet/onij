<template>
  <n-popover
    v-model:show="popoverOpen"
    trigger="manual"
    placement="bottom"
    :show-arrow="true"
    raw
    class="music-mv-popover"
    @clickoutside="onPopoverClickOutside"
  >
    <template #trigger>
      <button
        type="button"
        class="music-mv-btn"
        :class="{ 'music-mv-btn--active': hasMv }"
        :title="hasMv ? '打开 MV（双击编辑）' : '设置 MV 链接'"
        :aria-label="hasMv ? '打开 MV' : '设置 MV 链接'"
        @click.stop="onIconClick"
        @dblclick.stop="onIconDblClick"
      >
        <svg
          class="music-mv-btn__icon"
          viewBox="0 0 24 24"
          fill="currentColor"
          aria-hidden="true"
        >
          <path
            d="M21 7.5a1.5 1.5 0 0 0-1.5-1.5h-15A1.5 1.5 0 0 0 3 7.5v9A1.5 1.5 0 0 0 4.5 18h15a1.5 1.5 0 0 0 1.5-1.5v-9ZM8 16.5v-9l8 4.5-8 4.5Z"
          />
        </svg>
      </button>
    </template>

    <div class="music-mv-bubble">
      <input
        ref="inputRef"
        v-model="draftUrl"
        type="text"
        class="music-mv-bubble__input"
        :class="{ 'music-mv-bubble__input--error': !!errorMsg }"
        placeholder="https://..."
        spellcheck="false"
        autocomplete="off"
        @keydown.enter.prevent="commitSave"
        @blur="onInputBlur"
      />
      <p v-if="errorMsg" class="music-mv-bubble__error">{{ errorMsg }}</p>
    </div>
  </n-popover>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { NPopover } from 'naive-ui'
import { openMvUrl, validateMvUrl } from '@/util/musicMv'

const props = defineProps<{
  hasMv: boolean
  url?: string
}>()

const emit = defineEmits<{
  save: [url: string]
}>()

const popoverOpen = ref(false)
const inputRef = ref<HTMLInputElement | null>(null)
const draftUrl = ref('')
const lastCommitted = ref('')
const errorMsg = ref('')
let singleClickTimer: ReturnType<typeof setTimeout> | null = null
let skipBlurCommit = false

watch(
  () => props.url ?? '',
  (v) => {
    draftUrl.value = v
    lastCommitted.value = v
  },
  { immediate: true },
)

watch(popoverOpen, async (open) => {
  if (open) {
    errorMsg.value = ''
    await nextTick()
    inputRef.value?.focus()
    inputRef.value?.select()
  }
})

function clearSingleClickTimer() {
  if (singleClickTimer) {
    clearTimeout(singleClickTimer)
    singleClickTimer = null
  }
}

function onIconClick() {
  if (!props.hasMv) {
    popoverOpen.value = true
    return
  }

  clearSingleClickTimer()
  singleClickTimer = setTimeout(() => {
    singleClickTimer = null
    openMvUrl(props.url ?? '')
  }, 220)
}

function onIconDblClick() {
  clearSingleClickTimer()
  popoverOpen.value = true
}

function commitSave(): boolean {
  const result = validateMvUrl(draftUrl.value)
  if (!result.ok) {
    errorMsg.value = result.error ?? 'URL 无效'
    return false
  }

  errorMsg.value = ''
  if (result.url === lastCommitted.value) {
    popoverOpen.value = false
    return true
  }

  lastCommitted.value = result.url
  draftUrl.value = result.url
  emit('save', result.url)
  popoverOpen.value = false
  return true
}

function onInputBlur() {
  if (skipBlurCommit) {
    skipBlurCommit = false
    return
  }
  if (!commitSave()) popoverOpen.value = true
}

function onPopoverClickOutside() {
  skipBlurCommit = true
  if (!commitSave()) popoverOpen.value = true
}
</script>

<style scoped>
.music-mv-btn {
  display: inline-flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 1.75rem;
  height: 1.75rem;
  padding: 0;
  border: 0;
  border-radius: 0.375rem;
  background: transparent;
  color: rgb(156 163 175);
  cursor: pointer;
  transition: color 0.15s ease, transform 0.15s ease;
}

.music-mv-btn:hover {
  color: rgb(107 114 128);
  transform: scale(1.06);
}

.music-mv-btn--active {
  color: rgb(251 114 153);
}

.music-mv-btn--active:hover {
  color: rgb(244 63 94);
}

.music-mv-btn__icon {
  width: 1.125rem;
  height: 1.125rem;
}

.music-mv-bubble {
  padding: 0.125rem;
}

.music-mv-bubble__input {
  box-sizing: border-box;
  display: block;
  width: min(72vw, 16rem);
  height: 2rem;
  padding: 0 0.875rem;
  border: 1px solid rgb(229 231 235);
  border-radius: 9999px;
  background: #fff;
  box-shadow:
    0 4px 14px rgb(0 0 0 / 0.08),
    0 1px 3px rgb(0 0 0 / 0.06);
  font-size: 0.8125rem;
  line-height: 1.25rem;
  color: rgb(31 41 55);
  outline: none;
}

.music-mv-bubble__input::placeholder {
  color: rgb(156 163 175);
}

.music-mv-bubble__input:focus {
  border-color: rgb(251 114 153);
  box-shadow:
    0 4px 14px rgb(251 114 153 / 0.15),
    0 0 0 2px rgb(251 114 153 / 0.2);
}

.music-mv-bubble__input--error {
  border-color: rgb(248 113 113);
}

.music-mv-bubble__input--error:focus {
  box-shadow:
    0 4px 14px rgb(248 113 113 / 0.12),
    0 0 0 2px rgb(248 113 113 / 0.2);
}

.music-mv-bubble__error {
  margin: 0.375rem 0.875rem 0.125rem;
  font-size: 0.6875rem;
  line-height: 1.3;
  color: rgb(239 68 68);
}
</style>
