import { onMounted, onUnmounted, ref } from 'vue'
import { useRoute } from 'vue-router'

const SEQ = ['o', 'n', 'i', 'j'] as const
const RESET_MS = 1500

/** 初始隐藏；首页非输入框连续输入 onij 切换 */
export const navToolsVisible = ref(false)

function isHomePath(path: string): boolean {
  return path === '/search' || path === '/'
}

function isTypingTarget(target: EventTarget | null): boolean {
  if (!(target instanceof HTMLElement)) return false
  const tag = target.tagName
  if (tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT') return true
  return target.isContentEditable
}

export function useNavToolsVisible() {
  return { navToolsVisible }
}

/** 仅在 App 壳挂一次：首页敲 onij 切换导航 AI / 工具 */
export function useOnijNavToggle() {
  const route = useRoute()
  let idx = 0
  let timer: ReturnType<typeof setTimeout> | null = null

  function resetSeq() {
    idx = 0
    if (timer) {
      clearTimeout(timer)
      timer = null
    }
  }

  function onKeydown(e: KeyboardEvent) {
    if (!isHomePath(route.path)) {
      resetSeq()
      return
    }
    if (e.isComposing || e.metaKey || e.ctrlKey || e.altKey) return
    if (isTypingTarget(e.target)) {
      resetSeq()
      return
    }
    if (e.key.length !== 1) return
    const key = e.key.toLowerCase()
    if (key === SEQ[idx]) {
      idx += 1
      if (timer) clearTimeout(timer)
      if (idx === SEQ.length) {
        resetSeq()
        navToolsVisible.value = !navToolsVisible.value
        return
      }
      timer = setTimeout(resetSeq, RESET_MS)
      return
    }
    if (key === SEQ[0]) {
      idx = 1
      if (timer) clearTimeout(timer)
      timer = setTimeout(resetSeq, RESET_MS)
      return
    }
    resetSeq()
  }

  onMounted(() => {
    window.addEventListener('keydown', onKeydown)
  })
  onUnmounted(() => {
    resetSeq()
    window.removeEventListener('keydown', onKeydown)
  })

  return { navToolsVisible }
}
