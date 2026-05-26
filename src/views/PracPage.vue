<template>
  <PageContent>
    <div class="prac-layout">
      <h1 class="browse-section-title mb-6 border-l-4 border-green-500 pl-2">练习</h1>

      <div class="prac-split">
        <section class="prac-card prac-calendar" aria-label="练习日历">
          <div class="prac-calendar__head">
            <button type="button" class="prac-cal-nav" aria-label="上一月" @click="prevMonth">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M15 18l-6-6 6-6" />
              </svg>
            </button>
            <h2 class="prac-calendar__title">{{ monthTitle }}</h2>
            <button type="button" class="prac-cal-nav" aria-label="下一月" @click="nextMonth">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 18l6-6-6-6" />
              </svg>
            </button>
          </div>

          <div class="prac-calendar__weekdays">
            <span v-for="w in weekdays" :key="w">{{ w }}</span>
          </div>
          <div class="prac-calendar__grid">
            <div
              v-for="cell in monthCells"
              :key="cell.key"
              class="prac-calendar__cell-wrap"
            >
              <span v-if="cell.day === null" class="prac-calendar__cell prac-calendar__cell--empty" />
              <button
                v-else
                type="button"
                class="prac-day-btn prac-calendar__cell"
                :class="{ 'prac-calendar__cell--selected': selectedDateKey === cell.dateKey }"
                :aria-label="isToday(cell.dateKey) ? '今天' : `${cell.day}日`"
                @click="selectDate(cell.dateKey)"
              >
                {{ dayCellLabel(cell) }}
                <span
                  v-if="daysWithPractice.has(cell.day!)"
                  class="prac-day-dot"
                  :class="{ 'prac-day-dot--selected': selectedDateKey === cell.dateKey }"
                />
              </button>
            </div>
          </div>
        </section>

        <section class="prac-card prac-panel">
          <p class="prac-panel__date">{{ selectedDateLabel }}</p>

          <p v-if="loading" class="browse-muted py-6 text-center text-sm">加载中…</p>
          <template v-else>
            <div class="prac-tab-row">
              <div class="prac-tab-row__scroll">
                <button
                  v-for="p in dayPractices"
                  :key="p.id"
                  type="button"
                  class="prac-tab"
                  :class="{
                    'prac-tab--active': activeTab === tabKey(p.id) && !isEditing,
                    'prac-tab--editing': isEditing && editingId === p.id,
                  }"
                  title="双击编辑"
                  @click="openPractice(p.id)"
                  @dblclick.stop="onTabDblClick(p)"
                >
                  <component :is="practiceIcon(p.practice_type)" class="prac-tab__icon" />
                  <span>{{ practiceLabel(p) }}</span>
                </button>
                <button
                  v-if="isEditing && editingId === null"
                  type="button"
                  class="prac-tab prac-tab--active"
                >
                  <component :is="practiceIcon(editForm.practice_type)" class="prac-tab__icon" />
                  <span>新练习</span>
                </button>
              </div>
              <button
                type="button"
                class="prac-tab-add"
                title="新建练习"
                aria-label="新建练习"
                :disabled="saving"
                @click="startNewPractice"
              >
                <span class="prac-tab-add__icon" aria-hidden="true">+</span>
              </button>
            </div>

            <div class="prac-panel__divider" />

            <div
              v-if="isEditing"
              ref="editPanelRef"
              class="prac-edit"
              @paste="onEditPaste"
            >
              <div class="prac-edit__toolbar">
                <select
                  v-model.number="editForm.practice_type"
                  class="prac-field"
                >
                  <option
                    v-for="(label, val) in practiceTypeOptions"
                    :key="val"
                    :value="Number(val)"
                  >
                    {{ label }}
                  </option>
                </select>
                <label class="prac-duration">
                  <input
                    v-model.number="editForm.duration"
                    type="number"
                    min="0"
                    step="10"
                    class="prac-field prac-field--num"
                  >
                  <span>分钟</span>
                </label>
                <div class="prac-edit__actions">
                  <input
                    ref="imageInputRef"
                    type="file"
                    accept="image/*"
                    class="sr-only"
                    @change="onImagePick"
                  >
                  <button
                    type="button"
                    class="prac-icon-btn"
                    aria-label="插入图片"
                    title="插入图片"
                    :disabled="uploadingImage"
                    @click="imageInputRef?.click()"
                  >
                    <ImageOutline class="prac-icon-btn__svg" />
                  </button>
                  <button
                    type="button"
                    class="prac-icon-btn prac-icon-btn--primary"
                    aria-label="保存"
                    title="保存"
                    :disabled="saving || uploadingImage"
                    @click="savePractice"
                  >
                    <SaveOutline class="prac-icon-btn__svg" />
                  </button>
                  <button
                    v-if="editForm.id"
                    type="button"
                    class="prac-icon-btn prac-icon-btn--danger"
                    aria-label="删除"
                    title="删除"
                    :disabled="saving"
                    @click="removePractice"
                  >
                    <TrashOutline class="prac-icon-btn__svg" />
                  </button>
                </div>
                <button
                  v-if="editForm.id"
                  type="button"
                  class="prac-btn prac-btn--ghost"
                  @click="cancelEdit"
                >
                  取消
                </button>
              </div>
              <textarea
                ref="contentRef"
                v-model="editForm.content"
                class="prac-textarea"
                placeholder="练习内容，可用 ![文件id] 插入图片；Ctrl+V 可粘贴图片"
              />
              <p v-if="uploadingImage" class="prac-hint">图片上传中…</p>
            </div>

            <template v-else-if="selectedPractice">
              <div class="prac-content-body" v-html="contentHtml" />
            </template>

            <p v-else class="browse-muted py-12 text-center text-sm">暂无练习，点击右侧 + 添加</p>
          </template>
        </section>
      </div>
    </div>
  </PageContent>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { ImageOutline, SaveOutline, TrashOutline } from '@vicons/ionicons5'
import PageContent from '@/components/layout/PageContent.vue'
import { GetMonthlyPractice, UploadPractice, DeletePractice } from '@/api/practice'
import { DownloadFiles, UploadFile } from '@/api/file'
import type { Practice } from '@/api/types/practice'
import { PracticeType, PRACTICE_TYPE_LABELS } from '@/api/types/enums'
import { practiceIcon } from '@/util/practiceIcons'
import { uploadToQiniu, getFileTypeFromFile } from '@/util/qiniu'
import {
  extractImageIds,
  renderPracticeContentHtml,
  dateKeyToPracticeAt,
  todayDateKey,
  toDateKey,
} from '@/util/practiceContent'

const weekdays = ['一', '二', '三', '四', '五', '六', '日'] as const

const viewYear = ref(new Date().getFullYear())
const viewMonth = ref(new Date().getMonth())
const selectedDateKey = ref(todayDateKey())
const monthlyDays = ref<Map<number, Practice[]>>(new Map())
const loading = ref(false)
const saving = ref(false)
const uploadingImage = ref(false)

const activeTab = ref<string | null>(null)
const isEditing = ref(false)
const editingId = ref<number | null>(null)
const selectedPractice = ref<Practice | null>(null)
const imageUrlMap = ref<Record<number, string>>({})

const contentRef = ref<HTMLTextAreaElement | null>(null)
const imageInputRef = ref<HTMLInputElement | null>(null)
const editPanelRef = ref<HTMLElement | null>(null)

const editForm = ref({
  id: undefined as number | undefined,
  practice_type: PracticeType.PCT_Other,
  content: '',
  duration: 10,
})

const practiceTypeOptions = Object.fromEntries(
  Object.entries(PRACTICE_TYPE_LABELS).filter(([k]) => Number(k) !== PracticeType.PCT_Unknown),
) as Record<number, string>

const monthTitle = computed(() => `${viewYear.value}年${viewMonth.value + 1}月`)

interface MonthCell {
  key: string
  day: number | null
  dateKey: string
}

const monthCells = computed((): MonthCell[] => {
  const y = viewYear.value
  const m = viewMonth.value
  const first = new Date(y, m, 1)
  const startOffset = (first.getDay() + 6) % 7
  const daysInMonth = new Date(y, m + 1, 0).getDate()

  const cells: MonthCell[] = []
  for (let i = 0; i < startOffset; i++) {
    cells.push({ key: `pad-${y}-${m}-b-${i}`, day: null, dateKey: '' })
  }
  for (let d = 1; d <= daysInMonth; d++) {
    const dateKey = toDateKey(y, m, d)
    cells.push({ key: dateKey, day: d, dateKey })
  }
  while (cells.length % 7 !== 0) {
    cells.push({ key: `pad-${y}-${m}-a-${cells.length}`, day: null, dateKey: '' })
  }
  return cells
})

const daysWithPractice = computed(() => new Set(monthlyDays.value.keys()))
const selectedDay = computed(() => Number(selectedDateKey.value.split('-')[2]))
const dayPractices = computed(() => monthlyDays.value.get(selectedDay.value) ?? [])

const selectedDateLabel = computed(() => {
  const [y, m, d] = selectedDateKey.value.split('-')
  return `${y}年${Number(m)}月${Number(d)}日`
})

const contentHtml = computed(() => {
  if (!selectedPractice.value) return ''
  return renderPracticeContentHtml(selectedPractice.value.content, imageUrlMap.value)
})

function tabKey(id: number) {
  return `p-${id}`
}

function isToday(dateKey: string) {
  return dateKey === todayDateKey()
}

function dayCellLabel(cell: MonthCell) {
  if (isToday(cell.dateKey)) return '今'
  return cell.day!
}

function practiceLabel(p: Practice) {
  return PRACTICE_TYPE_LABELS[p.practice_type as PracticeType] ?? '练习'
}

function roundDuration(minutes: number) {
  if (minutes <= 0) return 0
  return Math.round(minutes / 10) * 10
}

function openPractice(id: number) {
  const p = dayPractices.value.find((x) => x.id === id)
  if (!p) return
  isEditing.value = false
  editingId.value = null
  activeTab.value = tabKey(id)
  selectPractice(p)
}

function selectDate(dateKey: string) {
  selectedDateKey.value = dateKey
  isEditing.value = false
  editingId.value = null
  const list = monthlyDays.value.get(Number(dateKey.split('-')[2])) ?? []
  if (list.length) {
    openPractice(list[0].id)
  } else {
    activeTab.value = null
    selectedPractice.value = null
  }
}

function selectPractice(p: Practice) {
  selectedPractice.value = p
  void loadImageUrls(p.content)
}

function startNewPractice() {
  isEditing.value = true
  editingId.value = null
  activeTab.value = '__new__'
  selectedPractice.value = null
  editForm.value = {
    id: undefined,
    practice_type: PracticeType.PCT_Other,
    content: '',
    duration: 10,
  }
}

function editSelected() {
  if (!selectedPractice.value) return
  isEditing.value = true
  editingId.value = selectedPractice.value.id
  activeTab.value = tabKey(selectedPractice.value.id)
  editForm.value = {
    id: selectedPractice.value.id,
    practice_type: selectedPractice.value.practice_type,
    content: selectedPractice.value.content,
    duration: roundDuration(selectedPractice.value.duration),
  }
}

function onTabDblClick(p: Practice) {
  if (isEditing.value) return
  if (activeTab.value !== tabKey(p.id)) return
  if (selectedPractice.value?.id !== p.id) {
    selectPractice(p)
  }
  editSelected()
}

function cancelEdit() {
  isEditing.value = false
  editingId.value = null
  if (selectedPractice.value) {
    activeTab.value = tabKey(selectedPractice.value.id)
  } else if (dayPractices.value.length) {
    openPractice(dayPractices.value[0].id)
  } else {
    activeTab.value = null
  }
}

async function insertImageToken(file: File) {
  uploadingImage.value = true
  try {
    const qiniu = await uploadToQiniu(file, ['practice'])
    const resp = await UploadFile({
      parent_id: 0,
      files: [
        {
          filename: file.name,
          store_key: qiniu.key,
          hash: qiniu.hash,
          size: file.size,
          format: getFileTypeFromFile(file),
          exif: '',
        },
      ],
    })
    const fileId = resp.file_ids?.[0]
    if (!fileId) throw new Error(resp.message || '登记文件失败')

    const token = `![${fileId}]`
    const ta = contentRef.value
    if (ta) {
      const start = ta.selectionStart ?? editForm.value.content.length
      const end = ta.selectionEnd ?? start
      const c = editForm.value.content
      editForm.value.content = c.slice(0, start) + token + c.slice(end)
    } else {
      editForm.value.content += token
    }
    imageUrlMap.value = { ...imageUrlMap.value, [fileId]: qiniu.url }
  } finally {
    uploadingImage.value = false
  }
}

async function onImagePick(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file?.type.startsWith('image/')) return
  await insertImageToken(file)
}

async function onEditPaste(e: ClipboardEvent) {
  if (!isEditing.value) return
  const items = e.clipboardData?.items
  if (!items) return
  for (const item of items) {
    if (item.kind !== 'file') continue
    const file = item.getAsFile()
    if (file?.type.startsWith('image/')) {
      e.preventDefault()
      await insertImageToken(file)
      return
    }
  }
}

async function fetchMonthly() {
  loading.value = true
  try {
    const resp = await GetMonthlyPractice({
      year: viewYear.value,
      month: viewMonth.value + 1,
    })
    const map = new Map<number, Practice[]>()
    for (const d of resp.days ?? []) {
      map.set(d.day, d.practices ?? [])
    }
    monthlyDays.value = map

    const list = map.get(selectedDay.value) ?? []
    if (list.length && !isEditing.value) {
      const keep = selectedPractice.value?.id
      const target = keep ? list.find((p) => p.id === keep) ?? list[0] : list[0]
      openPractice(target.id)
    } else if (!list.length && !isEditing.value) {
      selectedPractice.value = null
      activeTab.value = null
    }
  } finally {
    loading.value = false
  }
}

async function loadImageUrls(content: string) {
  const ids = extractImageIds(content)
  if (!ids.length) return
  const missing = ids.filter((id) => !imageUrlMap.value[id])
  if (!missing.length) return
  try {
    const resp = await DownloadFiles({ file_ids: missing })
    const urls = resp.urls ?? []
    const next = { ...imageUrlMap.value }
    missing.forEach((id, i) => {
      if (urls[i]) next[id] = urls[i]
    })
    imageUrlMap.value = next
  } catch {
    /* ignore */
  }
}

async function savePractice() {
  saving.value = true
  try {
    const resp = await UploadPractice({
      id: editForm.value.id,
      practice_type: editForm.value.practice_type,
      content: editForm.value.content,
      practice_at: dateKeyToPracticeAt(selectedDateKey.value),
      duration: editForm.value.duration,
    })
    isEditing.value = false
    editingId.value = null
    await fetchMonthly()
    const id = editForm.value.id ?? resp.practice_id
    openPractice(id)
  } catch (e) {
    alert(e instanceof Error ? e.message : '保存失败')
  } finally {
    saving.value = false
  }
}

async function removePractice() {
  if (!editForm.value.id) return
  if (!window.confirm('确定删除这条练习？')) return
  saving.value = true
  try {
    await DeletePractice({ id: editForm.value.id })
    isEditing.value = false
    editingId.value = null
    selectedPractice.value = null
    activeTab.value = null
    await fetchMonthly()
  } catch (e) {
    alert(e instanceof Error ? e.message : '删除失败')
  } finally {
    saving.value = false
  }
}

function prevMonth() {
  if (viewMonth.value === 0) {
    viewMonth.value = 11
    viewYear.value -= 1
  } else {
    viewMonth.value -= 1
  }
}

function nextMonth() {
  if (viewMonth.value === 11) {
    viewMonth.value = 0
    viewYear.value += 1
  } else {
    viewMonth.value += 1
  }
}

watch([viewYear, viewMonth], () => {
  void fetchMonthly()
})

onMounted(() => {
  void fetchMonthly()
})
</script>

<style scoped>
.prac-layout {
  width: 100%;
  padding-bottom: 2rem;
}

.prac-split {
  display: grid;
  grid-template-columns: minmax(0, 20rem) minmax(0, 1fr);
  gap: 1rem;
  align-items: start;
}

@media (max-width: 768px) {
  .prac-split {
    grid-template-columns: 1fr;
  }
}

.prac-card {
  border: 1px solid rgb(229 231 235);
  border-radius: 0.75rem;
  background: rgb(255 255 255);
  padding: 1rem 1.1rem 1.25rem;
  box-shadow: 0 1px 3px rgb(0 0 0 / 0.06);
}

.dark .prac-card {
  border-color: rgb(55 65 81);
  background: rgb(31 41 55);
}

.prac-calendar__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.prac-calendar__title {
  margin: 0;
  flex: 1;
  text-align: center;
  font-size: 1.0625rem;
  font-weight: 700;
  color: rgb(17 24 39);
}

.dark .prac-calendar__title {
  color: rgb(243 244 246);
}

.prac-cal-nav {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.25rem;
  height: 2.25rem;
  border: 0;
  border-radius: 0.5rem;
  background: transparent;
  color: rgb(107 114 128);
  cursor: pointer;
}

.prac-cal-nav:hover {
  background: rgb(243 244 246);
}

.dark .prac-cal-nav:hover {
  background: rgb(55 65 81);
  color: rgb(229 231 235);
}

.prac-calendar__weekdays {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 0.25rem;
  margin-bottom: 0.35rem;
  text-align: center;
  font-size: 0.75rem;
  font-weight: 600;
  color: rgb(107 114 128);
}

.prac-calendar__grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 0.35rem;
}

.prac-calendar__cell-wrap {
  aspect-ratio: 1;
  min-height: 2.25rem;
}

.prac-calendar__cell--empty {
  display: block;
}

.prac-day-btn {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  border: 1px solid transparent;
  border-radius: 0.5rem;
  background: transparent;
  font-size: 0.875rem;
  font-weight: 500;
  color: rgb(55 65 81);
  cursor: pointer;
  outline: none;
  box-shadow: none;
}

.dark .prac-day-btn {
  color: rgb(229 231 235);
}

.prac-day-btn:hover:not(.prac-calendar__cell--selected) {
  background: rgb(243 244 246);
}

.dark .prac-day-btn:hover:not(.prac-calendar__cell--selected) {
  background: rgb(55 65 81);
}

.prac-day-btn:focus,
.prac-day-btn:focus-visible {
  outline: none;
  box-shadow: none;
}

.prac-calendar__cell--selected {
  border-color: rgb(37 99 235);
  color: rgb(37 99 235);
  font-weight: 600;
}

.prac-calendar__cell--selected:hover {
  background: rgb(239 246 255);
}

.dark .prac-calendar__cell--selected:hover {
  background: rgb(30 58 138 / 0.25);
}

.prac-day-dot {
  position: absolute;
  bottom: 4px;
  left: 50%;
  width: 4px;
  height: 4px;
  transform: translateX(-50%);
  border-radius: 50%;
  background: rgb(37 99 235);
}

.prac-day-dot--selected {
  background: rgb(37 99 235);
}

.prac-panel {
  min-height: 22rem;
}

.prac-panel__date {
  margin: 0 0 0.75rem;
  font-size: 0.8125rem;
  color: rgb(107 114 128);
}

.prac-tab-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.prac-tab-row__scroll {
  display: flex;
  flex: 1;
  gap: 0.35rem;
  overflow-x: auto;
  padding-bottom: 2px;
}

.prac-tab {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  flex-shrink: 0;
  padding: 0.4rem 0.75rem;
  border: 1px solid rgb(229 231 235);
  border-radius: 9999px;
  background: rgb(249 250 251);
  font-size: 0.875rem;
  font-weight: 500;
  color: rgb(55 65 81);
  cursor: pointer;
  transition: background 0.12s, border-color 0.12s, color 0.12s;
}

.dark .prac-tab {
  border-color: rgb(75 85 99);
  background: rgb(55 65 81);
  color: rgb(229 231 235);
}

.prac-tab:hover {
  border-color: rgb(147 197 253);
}

.prac-tab--active {
  background: rgb(37 99 235);
  border-color: rgb(37 99 235);
  color: #fff;
}

.prac-tab--editing {
  border-color: rgb(59 130 246);
  box-shadow: 0 0 0 2px rgb(59 130 246 / 0.25);
}

.prac-tab__icon {
  width: 1rem;
  height: 1rem;
  flex-shrink: 0;
}

.prac-tab-add {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  width: 2.25rem;
  height: 2.25rem;
  padding: 0;
  border: 1px solid rgb(37 99 235);
  border-radius: 0.5rem;
  background: rgb(37 99 235);
  font-size: 1.375rem;
  font-weight: 400;
  line-height: 1;
  color: #fff;
  cursor: pointer;
}

.prac-tab-add__icon {
  display: block;
  line-height: 1;
  transform: translateY(-0.05em);
}

.prac-tab-add:hover {
  background: rgb(29 78 216);
}

.prac-tab-add:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.prac-panel__divider {
  height: 1px;
  margin: 0.75rem 0;
  background: rgb(229 231 235);
}

.dark .prac-panel__divider {
  background: rgb(55 65 81);
}

.prac-edit__toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
}

.prac-edit__actions {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  margin-left: auto;
}

.prac-icon-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 2.25rem;
  height: 2.25rem;
  padding: 0;
  border: 0;
  border-radius: 0.5rem;
  background: transparent;
  color: rgb(107 114 128);
  cursor: pointer;
  transition: background-color 0.12s ease, color 0.12s ease;
}

.prac-icon-btn__svg {
  width: 1.25rem;
  height: 1.25rem;
}

.prac-icon-btn:hover:not(:disabled) {
  background: rgb(243 244 246);
  color: rgb(55 65 81);
}

.dark .prac-icon-btn {
  color: rgb(156 163 175);
}

.dark .prac-icon-btn:hover:not(:disabled) {
  background: rgb(55 65 81);
  color: rgb(243 244 246);
}

.prac-icon-btn--primary {
  color: rgb(37 99 235);
}

.prac-icon-btn--primary:hover:not(:disabled) {
  background: rgb(219 234 254);
  color: rgb(29 78 216);
}

.dark .prac-icon-btn--primary {
  color: rgb(96 165 250);
}

.dark .prac-icon-btn--primary:hover:not(:disabled) {
  background: rgb(30 64 175);
  color: rgb(191 219 254);
}

.prac-icon-btn--danger {
  color: rgb(220 38 38);
}

.prac-icon-btn--danger:hover:not(:disabled) {
  background: rgb(254 242 242);
  color: rgb(185 28 28);
}

.dark .prac-icon-btn--danger {
  color: rgb(248 113 113);
}

.dark .prac-icon-btn--danger:hover:not(:disabled) {
  background: rgb(127 29 29);
  color: rgb(254 202 202);
}

.prac-icon-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.prac-field {
  border: 1px solid rgb(229 231 235);
  border-radius: 0.5rem;
  padding: 0.4rem 0.6rem;
  font-size: 0.875rem;
  background: #fff;
  color: rgb(17 24 39);
}

.dark .prac-field {
  border-color: rgb(75 85 99);
  background: rgb(17 24 39);
  color: rgb(243 244 246);
}

.prac-field--num {
  width: 4rem;
}

.prac-duration {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  font-size: 0.875rem;
  color: rgb(107 114 128);
}

.prac-btn {
  border: 0;
  border-radius: 0.5rem;
  padding: 0.4rem 0.75rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
}

.prac-btn--primary {
  background: rgb(37 99 235);
  color: #fff;
}

.prac-btn--primary:hover {
  background: rgb(29 78 216);
}

.prac-btn--ghost {
  background: rgb(243 244 246);
  color: rgb(55 65 81);
}

.dark .prac-btn--ghost {
  background: rgb(55 65 81);
  color: rgb(229 231 235);
}

.prac-btn--danger {
  background: transparent;
  color: rgb(220 38 38);
}

.prac-textarea {
  width: 100%;
  min-height: 14rem;
  box-sizing: border-box;
  border: 1px solid rgb(229 231 235);
  border-radius: 0.5rem;
  padding: 0.75rem;
  font-size: 0.9375rem;
  line-height: 1.65;
  resize: vertical;
  background: rgb(249 250 251);
  color: rgb(17 24 39);
}

.dark .prac-textarea {
  border-color: rgb(75 85 99);
  background: rgb(17 24 39);
  color: rgb(243 244 246);
}

.prac-hint {
  margin: 0.5rem 0 0;
  font-size: 0.75rem;
  color: rgb(107 114 128);
}

.prac-content-body {
  min-height: 12rem;
  font-size: 0.9375rem;
  line-height: 1.7;
  color: rgb(55 65 81);
  word-break: break-word;
}

.dark .prac-content-body {
  color: rgb(229 231 235);
}

:deep(.prac-content-img) {
  display: block;
  max-width: 100%;
  margin: 0.5rem 0;
  border-radius: 0.5rem;
}

.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}
</style>
