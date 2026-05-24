<template>
  <PageContent>
    <div class="w-full max-w-lg">
      <h1 class="browse-section-title mb-6 border-l-4 border-green-500 pl-2">练习</h1>

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
            class="aspect-square min-h-[2.5rem]"
          >
            <span v-if="cell.day === null" class="block h-full w-full" />
            <button
              v-else
              type="button"
              class="prac-day-btn flex h-full w-full items-center justify-center rounded-lg border border-transparent text-[0.9375rem] font-medium text-gray-700 hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-gray-700"
              :class="{
                'border-blue-500 text-blue-600 dark:text-blue-400': isToday(cell.dateKey),
                '!border-blue-600 !bg-blue-600 !text-white hover:!bg-blue-700': isSelected(cell.dateKey),
              }"
              @click="toggleDate(cell.dateKey)"
            >
              {{ cell.day }}
            </button>
          </div>
        </div>
      </section>
    </div>
  </PageContent>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import PageContent from '@/components/layout/PageContent.vue'

const weekdays = ['一', '二', '三', '四', '五', '六', '日'] as const

const viewYear = ref(new Date().getFullYear())
const viewMonth = ref(new Date().getMonth())
const selectedKeys = ref<Set<string>>(new Set())

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

function toDateKey(year: number, month: number, day: number) {
  return `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`
}

function todayKey() {
  const t = new Date()
  return toDateKey(t.getFullYear(), t.getMonth(), t.getDate())
}

function isSelected(dateKey: string) {
  return selectedKeys.value.has(dateKey)
}

function isToday(dateKey: string) {
  return dateKey === todayKey()
}

function toggleDate(dateKey: string) {
  const next = new Set(selectedKeys.value)
  if (next.has(dateKey)) next.delete(dateKey)
  else next.add(dateKey)
  selectedKeys.value = next
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
</script>
