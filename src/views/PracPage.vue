<template>
  <PageContent>
    <div class="prac-layout">
      <h1 class="browse-section-title mb-6 border-l-4 border-green-500 pl-2">练习</h1>

      <div class="prac-split">
        <!-- 左：日历 -->
        <section
          class="rounded-xl border border-gray-200 bg-white p-4 shadow-sm dark:border-gray-600 dark:bg-gray-800"
          aria-label="练习日历"
        >
          <div class="mb-4 flex items-center justify-between gap-2">
            <button
              type="button"
              class="prac-nav-btn flex h-9 w-9 shrink-0 items-center justify-center rounded-lg text-gray-500 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700"
              aria-label="上一月"
              @click="prevMonth"
            >
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M15 18l-6-6 6-6" />
              </svg>
            </button>
            <h2 class="m-0 flex-1 text-center text-lg font-bold text-gray-900 dark:text-gray-100">
              {{ monthTitle }}
            </h2>
            <button
              type="button"
              class="prac-nav-btn flex h-9 w-9 shrink-0 items-center justify-center rounded-lg text-gray-500 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700"
              aria-label="下一月"
              @click="nextMonth"
            >
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 18l6-6-6-6" />
              </svg>
            </button>
          </div>

          <div class="mb-2 grid grid-cols-7 gap-1 text-center text-sm font-semibold text-gray-500 dark:text-gray-400">
            <span v-for="w in weekdays" :key="w">{{ w }}</span>
          </div>

          <div class="grid grid-cols-7 gap-1">
            <div
              v-for="cell in monthCells"
              :key="cell.key"
              class="aspect-square min-h-[2.25rem]"
            >
              <span v-if="cell.day === null" class="block h-full w-full" />
              <button
                v-else
                type="button"
                class="prac-day-btn relative flex h-full w-full items-center justify-center rounded-lg border border-transparent text-[0.875rem] font-medium text-gray-700 hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-gray-700"
                :class="{
                  'border-blue-500 text-blue-600 dark:text-blue-400': isToday(cell.dateKey),
                  '!border-blue-600 !bg-blue-600 !text-white hover:!bg-blue-700': selectedDateKey === cell.dateKey,
                }"
                @click="selectDate(cell.dateKey)"
              >
                {{ cell.day }}
                <span
                  v-if="daysWithPractice.has(cell.day!)"
                  class="absolute bottom-1 left-1/2 h-1 w-1 -translate-x-1/2 rounded-full bg-blue-500"
                  :class="{ 'bg-white': selectedDateKey === cell.dateKey }"
                />
              </button>
            </div>
          </div>
        </section>

        <!-- 右：练习列表 + 内容 -->
        <div class="prac-right flex min-h-0 flex-col gap-4">
          <section class="rounded-xl border border-gray-200 bg-white p-4 shadow-sm dark:border-gray-600 dark:bg-gray-800">
            <div class="mb-3 flex items-center justify-between gap-2">
              <h3 class="m-0 text-base font-bold text-gray-900 dark:text-gray-100">
                {{ selectedDateLabel }} 的练习
              </h3>
              <button
                type="button"
                class="rounded-lg bg-blue-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
                :disabled="saving"
                @click="startNewPractice"
              >
                新建
              </button>
            </div>

            <p v-if="loading" class="browse-muted py-4 text-center text-sm">加载中…</p>
            <p v-else-if="!dayPractices.length" class="browse-muted py-4 text-center text-sm">暂无练习</p>
            <ul v-else class="m-0 flex list-none flex-col gap-1 p-0">
              <li v-for="p in dayPractices" :key="p.id">
                <button
                  type="button"
                  class="w-full rounded-lg border px-3 py-2 text-left text-sm transition-colors"
                  :class="
                    selectedPractice?.id === p.id
                      ? 'border-blue-500 bg-blue-50 text-blue-800 dark:bg-blue-950/40 dark:text-blue-200'
                      : 'border-transparent hover:bg-gray-50 dark:hover:bg-gray-700/60'
                  "
                  @click="selectPractice(p)"
                >
                  <span class="font-semibold">{{ practiceLabel(p) }}</span>
                  <span v-if="p.duration" class="ml-2 text-gray-500 dark:text-gray-400">{{ p.duration }} 分钟</span>
                </button>
              </li>
            </ul>
          </section>

          <section
            v-if="editing"
            class="flex min-h-[12rem] flex-1 flex-col rounded-xl border border-gray-200 bg-white p-4 shadow-sm dark:border-gray-600 dark:bg-gray-800"
          >
            <div class="mb-3 flex flex-wrap items-center gap-2">
              <select
                v-model.number="editForm.practice_type"
                class="rounded-lg border border-gray-200 bg-white px-2 py-1.5 text-sm dark:border-gray-600 dark:bg-gray-900 dark:text-gray-100"
              >
                <option v-for="(label, val) in practiceTypeOptions" :key="val" :value="Number(val)">
                  {{ label }}
                </option>
              </select>
              <input
                v-model.number="editForm.duration"
                type="number"
                min="0"
                placeholder="时长(分)"
                class="w-24 rounded-lg border border-gray-200 bg-white px-2 py-1.5 text-sm dark:border-gray-600 dark:bg-gray-900 dark:text-gray-100"
              />
              <label class="ml-auto cursor-pointer rounded-lg border border-gray-200 px-3 py-1.5 text-sm hover:bg-gray-50 dark:border-gray-600 dark:hover:bg-gray-700">
                插入图片
                <input type="file" accept="image/*" class="hidden" @change="onImagePick" />
              </label>
              <button
                type="button"
                class="rounded-lg bg-blue-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
                :disabled="saving || uploadingImage"
                @click="savePractice"
              >
                {{ saving ? '保存中…' : '保存' }}
              </button>
              <button
                v-if="editForm.id"
                type="button"
                class="rounded-lg px-3 py-1.5 text-sm text-red-600 hover:bg-red-50 dark:hover:bg-red-950/30"
                :disabled="saving"
                @click="removePractice"
              >
                删除
              </button>
            </div>
            <textarea
              ref="contentRef"
              v-model="editForm.content"
              class="min-h-[10rem] flex-1 resize-y rounded-lg border border-gray-200 bg-gray-50 p-3 text-sm leading-relaxed dark:border-gray-600 dark:bg-gray-900 dark:text-gray-100"
              placeholder="练习内容，可用 ![文件id] 插入图片"
            />
            <p v-if="uploadingImage" class="mt-2 text-xs text-gray-500">图片上传中…</p>
          </section>

          <section
            v-else-if="selectedPractice"
            class="flex min-h-[12rem] flex-1 flex-col rounded-xl border border-gray-200 bg-white p-4 shadow-sm dark:border-gray-600 dark:bg-gray-800"
          >
            <div class="mb-3 flex items-center justify-between">
              <h3 class="m-0 text-base font-bold text-gray-900 dark:text-gray-100">
                {{ practiceLabel(selectedPractice) }}
              </h3>
              <button
                type="button"
                class="text-sm text-blue-600 hover:underline dark:text-blue-400"
                @click="editSelected"
              >
                编辑
              </button>
            </div>
            <div
              class="prac-content-body flex-1 overflow-auto text-sm leading-relaxed text-gray-800 dark:text-gray-200"
              v-html="contentHtml"
            />
          </section>

          <section
            v-else
            class="flex min-h-[12rem] flex-1 items-center justify-center rounded-xl border border-dashed border-gray-300 dark:border-gray-600"
          >
            <p class="browse-muted text-sm">选择一条练习查看内容</p>
          </section>
        </div>
      </div>
    </div>
  </PageContent>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import PageContent from '@/components/layout/PageContent.vue'
import { GetMonthlyPractice, UploadPractice, DeletePractice } from '@/api/practice'
import { DownloadFiles, UploadFile } from '@/api/file'
import type { Practice } from '@/api/types/practice'
import { PracticeType, PRACTICE_TYPE_LABELS } from '@/api/types/enums'
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

const selectedPractice = ref<Practice | null>(null)
const editing = ref(false)
const contentRef = ref<HTMLTextAreaElement | null>(null)
const imageUrlMap = ref<Record<number, string>>({})

const editForm = ref({
  id: undefined as number | undefined,
  practice_type: PracticeType.PCT_Other,
  content: '',
  duration: 0,
})

const practiceTypeOptions = PRACTICE_TYPE_LABELS

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

function isToday(dateKey: string) {
  return dateKey === todayDateKey()
}

function practiceLabel(p: Practice) {
  return PRACTICE_TYPE_LABELS[p.practice_type as PracticeType] ?? '练习'
}

function selectDate(dateKey: string) {
  selectedDateKey.value = dateKey
  selectedPractice.value = null
  editing.value = false
  const list = monthlyDays.value.get(Number(dateKey.split('-')[2])) ?? []
  if (list.length) selectPractice(list[0])
}

function selectPractice(p: Practice) {
  selectedPractice.value = p
  editing.value = false
  void loadImageUrls(p.content)
}

function startNewPractice() {
  editing.value = true
  selectedPractice.value = null
  editForm.value = {
    id: undefined,
    practice_type: PracticeType.PCT_Other,
    content: '',
    duration: 0,
  }
}

function editSelected() {
  if (!selectedPractice.value) return
  editing.value = true
  editForm.value = {
    id: selectedPractice.value.id,
    practice_type: selectedPractice.value.practice_type,
    content: selectedPractice.value.content,
    duration: selectedPractice.value.duration,
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
    if (list.length && !selectedPractice.value && !editing.value) {
      selectPractice(list[0])
    } else if (!list.length) {
      selectedPractice.value = null
    }
  } finally {
    loading.value = false
  }
}

async function loadImageUrls(content: string) {
  const ids = extractImageIds(content)
  if (!ids.length) {
    imageUrlMap.value = {}
    return
  }
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

async function onImagePick(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return

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
    await fetchMonthly()
    const id = editForm.value.id ?? resp.practice_id
    const list = monthlyDays.value.get(selectedDay.value) ?? []
    const saved = list.find((p) => p.id === id) ?? list[list.length - 1]
    if (saved) selectPractice(saved)
    editing.value = false
  } finally {
    saving.value = false
  }
}

async function removePractice() {
  if (!editForm.value.id) return
  saving.value = true
  try {
    await DeletePractice({ id: editForm.value.id })
    editing.value = false
    selectedPractice.value = null
    await fetchMonthly()
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
  grid-template-columns: minmax(0, 22rem) minmax(0, 1fr);
  gap: 1.25rem;
  align-items: start;
}

@media (max-width: 768px) {
  .prac-split {
    grid-template-columns: 1fr;
  }
}

.prac-right {
  min-height: 20rem;
}

:deep(.prac-content-img) {
  display: block;
  max-width: 100%;
  margin: 0.5rem 0;
  border-radius: 0.5rem;
}
</style>
