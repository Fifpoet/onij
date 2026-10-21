<template>
  <div class="presence-nick" data-testid="presence-nick">
    <PresenceAvatar :seed="clientId" :title="nickname" size="sm" />
    <input
      v-if="editing"
      ref="inputRef"
      v-model="draft"
      class="presence-nick__input"
      maxlength="32"
      aria-label="修改昵称"
      data-testid="presence-nick-input"
      @keydown.enter.prevent="commit"
      @keydown.esc.prevent="cancel"
      @blur="commit"
    >
    <button
      v-else
      type="button"
      class="presence-nick__btn"
      :title="nickname || '点击改名'"
      aria-label="修改昵称"
      data-testid="presence-nick-btn"
      @click="startEdit"
    >
      {{ nickname || '未命名' }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { nextTick, ref } from 'vue'
import { storeToRefs } from 'pinia'
import { usePresenceStore } from '@/store/presence'
import PresenceAvatar from '@/components/presence/PresenceAvatar.vue'

defineProps<{
  truncated?: boolean
}>()

const presence = usePresenceStore()
const { nickname, clientId } = storeToRefs(presence)
presence.ensureIdentity()

const editing = ref(false)
const draft = ref('')
const inputRef = ref<HTMLInputElement | null>(null)

function startEdit() {
  draft.value = nickname.value
  editing.value = true
  void nextTick(() => {
    inputRef.value?.focus()
    inputRef.value?.select()
  })
}

function commit() {
  if (!editing.value) return
  presence.setNickname(draft.value)
  editing.value = false
}

function cancel() {
  editing.value = false
}
</script>

<style scoped>
.presence-nick {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  flex-shrink: 0;
  min-width: 0;
}
.presence-nick__btn,
.presence-nick__input {
  height: 2.25rem;
  padding: 0 0.35rem 0 0;
  border: 0;
  border-radius: 0.5rem;
  background: transparent;
  font-size: 0.8125rem;
  font-weight: 500;
  line-height: 1;
  color: rgb(55 65 81);
  white-space: nowrap;
}
.presence-nick__btn {
  text-align: left;
  cursor: pointer;
}
.presence-nick__input {
  width: 11rem;
  outline: none;
  box-shadow: inset 0 0 0 1px rgb(209 213 219);
  background: rgb(249 250 251);
  padding: 0 0.5rem;
}
.dark .presence-nick__btn,
.dark .presence-nick__input {
  color: rgb(229 231 235);
}
.dark .presence-nick__input {
  background: rgb(31 41 55);
  box-shadow: inset 0 0 0 1px rgb(75 85 99);
}
.presence-nick__btn:hover {
  background: rgb(243 244 246);
}
.dark .presence-nick__btn:hover {
  background: rgb(55 65 81);
}
</style>
