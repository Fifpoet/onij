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

          <button
            type="button"
            class="prac-stats-btn"
            :class="{ 'prac-stats-btn--active': panelMode === 'stats' }"
            @click="openMonthlyStats"
          >
            月度统计
          </button>

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
                :class="{
                  'prac-calendar__cell--selected': selectedDateKey === cell.dateKey,
                  [`prac-day-btn--heat-${cell.heatLevel}`]: cell.heatLevel > 0,
                }"
                :aria-label="cell.ariaLabel"
                :title="cell.tooltip || undefined"
                @click="selectDate(cell.dateKey)"
              >
                <span class="prac-day-btn__label">{{ dayCellLabel(cell) }}</span>
                <span
                  v-if="cell.typeMark === 'triangle'"
                  class="prac-day-mark prac-day-mark--triangle"
                  aria-hidden="true"
                />
                <span
                  v-else-if="cell.typeMark === 'circle'"
                  class="prac-day-mark prac-day-mark--circle"
                  aria-hidden="true"
                />
              </button>
            </div>
          </div>
        </section>

        <section class="prac-card prac-panel">
          <p class="prac-panel__date">
            {{ panelMode === 'stats' ? `${monthTitle} 统计` : selectedDateLabel }}
          </p>

          <p v-if="loading" class="browse-muted py-6 text-center text-sm">加载中…</p>
          <template v-else-if="panelMode === 'stats'">
            <div class="prac-stats-type-row">
              <button
                type="button"
                class="prac-stats-type"
                :class="{ 'prac-stats-type--active': statsTypeFilter === 'all' }"
                @click="statsTypeFilter = 'all'"
              >
                全部
              </button>
              <button
                v-for="type in monthPracticeTypes"
                :key="type"
                type="button"
                class="prac-stats-type"
                :class="{ 'prac-stats-type--active': statsTypeFilter === type }"
                @click="statsTypeFilter = type"
              >
                <component :is="practiceIcon(type)" class="prac-stats-type__icon" />
                <span>{{ PRACTICE_TYPE_LABELS[type] }}</span>
              </button>
            </div>

            <div class="prac-stats-summary">
              <div class="prac-stats-summary__item">
                <span class="prac-stats-summary__value">{{ filteredMonthStats.sessionCount }}</span>
                <span class="prac-stats-summary__label">练习次数</span>
              </div>
              <div class="prac-stats-summary__item">
                <span class="prac-stats-summary__value">{{ formatDurationShort(filteredMonthStats.durationMin) }}</span>
                <span class="prac-stats-summary__label">总时长</span>
              </div>
              <div class="prac-stats-summary__item">
                <span class="prac-stats-summary__value">{{ filteredMonthStats.activeDays }}</span>
                <span class="prac-stats-summary__label">活跃天数</span>
              </div>
            </div>

            <div class="prac-panel__divider" />

            <div v-if="filteredMonthStats.sessionCount === 0" class="browse-muted py-12 text-center text-sm">
              本月暂无该类型练习记录
            </div>
            <div v-else class="prac-chart-wrap">
              <p class="prac-chart__title">每日练习时长（分钟）</p>
              <svg
                class="prac-chart__svg"
                viewBox="0 0 400 180"
                role="img"
                :aria-label="`${statsTypeLabel} 本月每日练习时长折线图`"
              >
                <line
                  :x1="chartPad.l"
                  :y1="chartPad.t"
                  :x2="chartPad.l"
                  :y2="180 - chartPad.b"
                  class="prac-chart__axis"
                />
                <line
                  :x1="chartPad.l"
                  :y1="180 - chartPad.b"
                  :x2="400 - chartPad.r"
                  :y2="180 - chartPad.b"
                  class="prac-chart__axis"
                />
                <text :x="chartPad.l - 4" :y="180 - chartPad.b + 4" class="prac-chart__tick" text-anchor="end">0</text>
                <text :x="chartPad.l - 4" :y="chartPad.t + 4" class="prac-chart__tick" text-anchor="end">
                  {{ chartGeometry.maxY }}
                </text>
                <polyline
                  v-if="chartGeometry.polyline"
                  :points="chartGeometry.polyline"
                  class="prac-chart__line"
                />
                <circle
                  v-for="pt in chartGeometry.points"
                  :key="pt.day"
                  :cx="pt.x"
                  :cy="pt.y"
                  r="3"
                  class="prac-chart__dot"
                />
              </svg>
              <div class="prac-chart__xlabels">
                <span>1日</span>
                <span>{{ monthStats.daysInMonth }}日</span>
              </div>
            </div>
          </template>
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
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue'
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

type PanelMode = 'detail' | 'stats'
type StatsTypeFilter = 'all' | PracticeType
type DayTypeMark = 'none' | 'triangle' | 'circle'

interface DayStats {
  count: number
  durationMin: number
  typeCount: number
}

const chartPad = { l: 36, r: 12, t: 14, b: 28 }

const panelMode = ref<PanelMode>('detail')
const statsTypeFilter = ref<StatsTypeFilter>('all')

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

let pageActive = true
let fetchSeq = 0

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
  heatLevel: 0 | 1 | 2 | 3 | 4
  typeMark: DayTypeMark
  tooltip: string
  ariaLabel: string
}

function countPracticeTypes(practices: Practice[]): number {
  const types = new Set<PracticeType>()
  for (const p of practices) {
    types.add(p.practice_type as PracticeType)
  }
  return types.size
}

function typeMarkFromCount(typeCount: number): DayTypeMark {
  if (typeCount === 1) return 'triangle'
  if (typeCount >= 2) return 'circle'
  return 'none'
}

function filterPracticesByType(practices: Practice[], filter: StatsTypeFilter): Practice[] {
  if (filter === 'all') return practices
  return practices.filter((p) => p.practice_type === filter)
}

function calcHeatLevel(durationMin: number, maxDuration: number): 0 | 1 | 2 | 3 | 4 {
  if (durationMin <= 0) return 0
  if (maxDuration <= 0) return 1
  const ratio = durationMin / maxDuration
  if (ratio <= 0.25) return 1
  if (ratio <= 0.5) return 2
  if (ratio <= 0.75) return 3
  return 4
}

const monthCells = computed((): MonthCell[] => {
  const y = viewYear.value
  const m = viewMonth.value
  const first = new Date(y, m, 1)
  const startOffset = (first.getDay() + 6) % 7
  const daysInMonth = new Date(y, m + 1, 0).getDate()
  const { dayStatsMap } = monthStats.value

  let maxDuration = 0
  for (const stats of dayStatsMap.values()) {
    if (stats.durationMin > maxDuration) maxDuration = stats.durationMin
  }

  const emptyCell = (key: string): MonthCell => ({
    key,
    day: null,
    dateKey: '',
    heatLevel: 0,
    typeMark: 'none',
    tooltip: '',
    ariaLabel: '',
  })

  const cells: MonthCell[] = []
  for (let i = 0; i < startOffset; i++) {
    cells.push(emptyCell(`pad-${y}-${m}-b-${i}`))
  }
  for (let d = 1; d <= daysInMonth; d++) {
    const dateKey = toDateKey(y, m, d)
    const stats = dayStatsMap.get(d)
    const heatLevel = stats ? calcHeatLevel(stats.durationMin, maxDuration) : 0
    const typeMark = stats ? typeMarkFromCount(stats.typeCount) : 'none'
    let ariaLabel = isToday(dateKey) ? '今天' : `${d}日`
    let tooltip = ''
    if (stats) {
      tooltip = `${d}日 · ${stats.count}次 · ${formatDurationShort(stats.durationMin)}`
      ariaLabel = `${d}日，${stats.count}次练习，${formatDurationShort(stats.durationMin)}`
    }
    cells.push({
      key: dateKey,
      day: d,
      dateKey,
      heatLevel,
      typeMark,
      tooltip,
      ariaLabel,
    })
  }
  while (cells.length % 7 !== 0) {
    cells.push(emptyCell(`pad-${y}-${m}-a-${cells.length}`))
  }
  return cells
})


function statsForPractices(practices: Practice[]): DayStats {
  return {
    count: practices.length,
    durationMin: practices.reduce((sum, p) => sum + (Number(p.duration) || 0), 0),
    typeCount: countPracticeTypes(practices),
  }
}

function buildMonthStats(filter: StatsTypeFilter) {
  let sessionCount = 0
  let durationMin = 0
  let activeDays = 0
  const dayStatsMap = new Map<number, DayStats>()

  for (const [day, practices] of monthlyDays.value) {
    const filtered = filterPracticesByType(practices, filter)
    const stats = statsForPractices(filtered)
    if (stats.count <= 0) continue
    activeDays += 1
    sessionCount += stats.count
    durationMin += stats.durationMin
    dayStatsMap.set(day, stats)
  }

  const daysInMonth = new Date(viewYear.value, viewMonth.value + 1, 0).getDate()
  return { sessionCount, durationMin, activeDays, daysInMonth, dayStatsMap }
}

const monthStats = computed(() => buildMonthStats('all'))

const filteredMonthStats = computed(() => buildMonthStats(statsTypeFilter.value))

const monthPracticeTypes = computed(() => {
  const durationByType = new Map<PracticeType, number>()
  for (const practices of monthlyDays.value.values()) {
    for (const p of practices) {
      const type = p.practice_type as PracticeType
      if (type === PracticeType.PCT_Unknown) continue
      durationByType.set(type, (durationByType.get(type) ?? 0) + (Number(p.duration) || 0))
    }
  }
  return Array.from(durationByType.entries())
    .sort((a, b) => b[1] - a[1])
    .map(([type]) => type)
})

const statsTypeLabel = computed(() => {
  if (statsTypeFilter.value === 'all') return '全部类型'
  return PRACTICE_TYPE_LABELS[statsTypeFilter.value]
})

const chartGeometry = computed(() => {
  const daysInMonth = filteredMonthStats.value.daysInMonth
  const series: { day: number; minutes: number }[] = []
  for (let d = 1; d <= daysInMonth; d++) {
    const practices = monthlyDays.value.get(d) ?? []
    const filtered = filterPracticesByType(practices, statsTypeFilter.value)
    const minutes = filtered.reduce((sum, p) => sum + (Number(p.duration) || 0), 0)
    series.push({ day: d, minutes })
  }

  const maxY = Math.max(...series.map((p) => p.minutes), 1)
  const plotW = 400 - chartPad.l - chartPad.r
  const plotH = 180 - chartPad.t - chartPad.b
  const step = daysInMonth > 1 ? plotW / (daysInMonth - 1) : 0

  const points = series.map((p, i) => ({
    day: p.day,
    minutes: p.minutes,
    x: chartPad.l + i * step,
    y: chartPad.t + plotH - (p.minutes / maxY) * plotH,
  }))

  const polyline = points.map((p) => `${p.x},${p.y}`).join(' ')
  return { points, polyline, maxY: Math.round(maxY) }
})

function formatDurationShort(minutes: number): string {
  const m = Math.max(0, Math.round(minutes))
  if (m < 60) return `${m}分钟`
  const h = Math.floor(m / 60)
  const rest = m % 60
  if (rest === 0) return `${h}小时`
  return `${h}h ${rest}m`
}

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

function openMonthlyStats() {
  panelMode.value = 'stats'
  isEditing.value = false
  editingId.value = null
}

function selectDate(dateKey: string) {
  panelMode.value = 'detail'
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
  panelMode.value = 'detail'
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
  if (!pageActive) return
  uploadingImage.value = true
  try {
    const qiniu = await uploadToQiniu(file, ['practice'])
    if (!pageActive) return
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
    if (pageActive) uploadingImage.value = false
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
  const seq = ++fetchSeq
  if (pageActive) loading.value = true
  try {
    const resp = await GetMonthlyPractice({
      year: viewYear.value,
      month: viewMonth.value + 1,
    })
    if (!pageActive || seq !== fetchSeq) return

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
    if (pageActive && seq === fetchSeq) loading.value = false
  }
}

async function loadImageUrls(content: string) {
  const ids = extractImageIds(content)
  if (!ids.length) return
  const missing = ids.filter((id) => !imageUrlMap.value[id])
  if (!missing.length) return
  try {
    const resp = await DownloadFiles({ file_ids: missing })
    if (!pageActive) return
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

function syncSelectedDateToViewMonth() {
  const day = Number(selectedDateKey.value.split('-')[2])
  const daysInMonth = new Date(viewYear.value, viewMonth.value + 1, 0).getDate()
  selectedDateKey.value = toDateKey(viewYear.value, viewMonth.value, Math.min(day, daysInMonth))
}

function prevMonth() {
  if (viewMonth.value === 0) {
    viewMonth.value = 11
    viewYear.value -= 1
  } else {
    viewMonth.value -= 1
  }
  syncSelectedDateToViewMonth()
}

function nextMonth() {
  if (viewMonth.value === 11) {
    viewMonth.value = 0
    viewYear.value += 1
  } else {
    viewMonth.value += 1
  }
  syncSelectedDateToViewMonth()
}

watch([viewYear, viewMonth], () => {
  if (
    statsTypeFilter.value !== 'all'
    && !monthPracticeTypes.value.includes(statsTypeFilter.value)
  ) {
    statsTypeFilter.value = 'all'
  }
  void fetchMonthly()
})

onMounted(() => {
  void fetchMonthly()
})

onBeforeUnmount(() => {
  pageActive = false
  fetchSeq += 1
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
  margin-bottom: 0.75rem;
}

.prac-stats-btn {
  width: 100%;
  margin-bottom: 0.85rem;
  border: 1px solid rgb(187 247 208);
  border-radius: 0.5rem;
  padding: 0.45rem 0.75rem;
  font-size: 0.8125rem;
  font-weight: 600;
  color: rgb(22 101 52);
  background: rgb(240 253 244);
  cursor: pointer;
  transition: background 0.12s, border-color 0.12s, color 0.12s;
}

.prac-stats-btn:hover {
  background: rgb(220 252 231);
}

.prac-stats-btn--active {
  border-color: rgb(34 197 94);
  background: rgb(34 197 94);
  color: #fff;
}

.dark .prac-stats-btn {
  border-color: rgb(22 101 52 / 0.55);
  color: rgb(187 247 208);
  background: rgb(20 83 45 / 0.25);
}

.dark .prac-stats-btn--active {
  border-color: rgb(34 197 94);
  background: rgb(22 163 74);
  color: rgb(240 253 244);
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
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.15rem;
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

.prac-day-btn__label {
  line-height: 1;
}

.prac-day-mark {
  flex-shrink: 0;
}

.prac-day-mark--triangle {
  width: 0;
  height: 0;
  border-left: 3.5px solid transparent;
  border-right: 3.5px solid transparent;
  border-bottom: 5px solid rgb(22 163 74);
}

.prac-day-mark--circle {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: rgb(22 163 74);
}

.prac-calendar__cell--selected .prac-day-mark--triangle {
  border-bottom-color: rgb(37 99 235);
}

.prac-calendar__cell--selected .prac-day-mark--circle {
  background: rgb(37 99 235);
}

.dark .prac-day-mark--triangle {
  border-bottom-color: rgb(74 222 128);
}

.dark .prac-day-mark--circle {
  background: rgb(74 222 128);
}

.dark .prac-calendar__cell--selected .prac-day-mark--triangle {
  border-bottom-color: rgb(96 165 250);
}

.dark .prac-calendar__cell--selected .prac-day-mark--circle {
  background: rgb(96 165 250);
}

.dark .prac-day-btn {
  color: rgb(229 231 235);
}

.prac-day-btn:hover:not(.prac-calendar__cell--selected):not([class*='prac-day-btn--heat-']) {
  background: rgb(243 244 246);
}

.dark .prac-day-btn:hover:not(.prac-calendar__cell--selected):not([class*='prac-day-btn--heat-']) {
  background: rgb(55 65 81);
}

.prac-day-btn--heat-1 {
  background: rgb(220 252 231);
}

.prac-day-btn--heat-2 {
  background: rgb(187 247 208);
}

.prac-day-btn--heat-3 {
  background: rgb(134 239 172);
}

.prac-day-btn--heat-4 {
  background: rgb(74 222 128);
  color: rgb(20 83 45);
  font-weight: 600;
}

.dark .prac-day-btn--heat-1 {
  background: rgb(20 83 45 / 0.35);
}

.dark .prac-day-btn--heat-2 {
  background: rgb(22 101 52 / 0.5);
}

.dark .prac-day-btn--heat-3 {
  background: rgb(22 163 74 / 0.55);
}

.dark .prac-day-btn--heat-4 {
  background: rgb(34 197 94 / 0.65);
  color: rgb(240 253 244);
}

.prac-day-btn--heat-1:hover:not(.prac-calendar__cell--selected),
.prac-day-btn--heat-2:hover:not(.prac-calendar__cell--selected),
.prac-day-btn--heat-3:hover:not(.prac-calendar__cell--selected),
.prac-day-btn--heat-4:hover:not(.prac-calendar__cell--selected) {
  filter: brightness(0.96);
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

.prac-panel {
  min-height: 22rem;
}

.prac-panel__date {
  margin: 0 0 0.75rem;
  font-size: 0.8125rem;
  color: rgb(107 114 128);
}

.prac-stats-type-row {
  display: flex;
  flex-wrap: wrap;
  gap: 0.35rem;
  margin-bottom: 0.85rem;
}

.prac-stats-type {
  display: inline-flex;
  align-items: center;
  gap: 0.3rem;
  border: 1px solid rgb(229 231 235);
  border-radius: 9999px;
  padding: 0.35rem 0.7rem;
  font-size: 0.8125rem;
  font-weight: 500;
  color: rgb(55 65 81);
  background: rgb(249 250 251);
  cursor: pointer;
}

.dark .prac-stats-type {
  border-color: rgb(75 85 99);
  background: rgb(55 65 81);
  color: rgb(229 231 235);
}

.prac-stats-type--active {
  border-color: rgb(34 197 94);
  background: rgb(220 252 231);
  color: rgb(22 101 52);
}

.dark .prac-stats-type--active {
  border-color: rgb(34 197 94);
  background: rgb(20 83 45 / 0.45);
  color: rgb(187 247 208);
}

.prac-stats-type__icon {
  width: 0.9rem;
  height: 0.9rem;
  flex-shrink: 0;
}

.prac-stats-summary {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0.5rem;
}

.prac-stats-summary__item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.2rem;
  padding: 0.55rem 0.35rem;
  border-radius: 0.5rem;
  background: rgb(240 253 244);
  border: 1px solid rgb(187 247 208);
}

.dark .prac-stats-summary__item {
  background: rgb(20 83 45 / 0.2);
  border-color: rgb(22 101 52 / 0.45);
}

.prac-stats-summary__value {
  font-size: 1rem;
  font-weight: 700;
  color: rgb(22 101 52);
  text-align: center;
}

.dark .prac-stats-summary__value {
  color: rgb(134 239 172);
}

.prac-stats-summary__label {
  font-size: 0.6875rem;
  color: rgb(107 114 128);
}

.prac-chart-wrap {
  margin-top: 0.25rem;
}

.prac-chart__title {
  margin: 0 0 0.5rem;
  font-size: 0.8125rem;
  font-weight: 600;
  color: rgb(107 114 128);
}

.prac-chart__svg {
  display: block;
  width: 100%;
  height: auto;
}

.prac-chart__axis {
  stroke: rgb(209 213 219);
  stroke-width: 1;
}

.dark .prac-chart__axis {
  stroke: rgb(75 85 99);
}

.prac-chart__tick {
  fill: rgb(107 114 128);
  font-size: 10px;
}

.prac-chart__line {
  fill: none;
  stroke: rgb(34 197 94);
  stroke-width: 2;
  stroke-linejoin: round;
  stroke-linecap: round;
}

.dark .prac-chart__line {
  stroke: rgb(74 222 128);
}

.prac-chart__dot {
  fill: rgb(22 163 74);
}

.dark .prac-chart__dot {
  fill: rgb(74 222 128);
}

.prac-chart__xlabels {
  display: flex;
  justify-content: space-between;
  margin-top: 0.25rem;
  padding: 0 0.25rem 0 2rem;
  font-size: 0.6875rem;
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
