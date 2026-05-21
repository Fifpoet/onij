import { ref, computed, watch, type Ref } from 'vue'

export interface CoverTheme {
  bg: string
  fg: string
  fgMuted: string
  fgSubtle: string
  controlFg: string
  isDark: boolean
}

const FALLBACK: CoverTheme = {
  bg: '#1a1a1e',
  fg: '#f4f4f5',
  fgMuted: 'rgba(244, 244, 245, 0.72)',
  fgSubtle: 'rgba(244, 244, 245, 0.48)',
  controlFg: 'rgba(244, 244, 245, 0.92)',
  isDark: true,
}

function mix(a: number, b: number, t: number) {
  return Math.round(a + (b - a) * t)
}

function buildTheme(r: number, g: number, b: number): CoverTheme {
  const luminance = (0.299 * r + 0.587 * g + 0.114 * b) / 255
  const isDark = luminance < 0.52

  const bgR = mix(r, isDark ? 0 : 255, isDark ? 0.72 : 0.55)
  const bgG = mix(g, isDark ? 0 : 255, isDark ? 0.72 : 0.55)
  const bgB = mix(b, isDark ? 0 : 255, isDark ? 0.72 : 0.55)
  const bg = `rgb(${bgR}, ${bgG}, ${bgB})`

  const fg = isDark ? '#f8fafc' : '#0f172a'
  const fgMuted = isDark ? 'rgba(248, 250, 252, 0.72)' : 'rgba(15, 23, 42, 0.72)'
  const fgSubtle = isDark ? 'rgba(248, 250, 252, 0.48)' : 'rgba(15, 23, 42, 0.48)'
  const controlFg = isDark ? 'rgba(248, 250, 252, 0.92)' : 'rgba(15, 23, 42, 0.88)'

  return { bg, fg, fgMuted, fgSubtle, controlFg, isDark }
}

let sampleToken = 0

async function sampleCover(url: string): Promise<CoverTheme> {
  if (!url) return { ...FALLBACK }

  return new Promise((resolve) => {
    const img = new Image()
    img.crossOrigin = 'anonymous'
    img.decoding = 'async'
    img.onload = () => {
      try {
        const size = 28
        const canvas = document.createElement('canvas')
        canvas.width = size
        canvas.height = size
        const ctx = canvas.getContext('2d', { willReadFrequently: true })
        if (!ctx) {
          resolve({ ...FALLBACK })
          return
        }
        ctx.drawImage(img, 0, 0, size, size)
        const data = ctx.getImageData(0, 0, size, size).data
        let r = 0
        let g = 0
        let b = 0
        let n = 0
        for (let i = 0; i < data.length; i += 4) {
          const a = data[i + 3]
          if (a < 16) continue
          r += data[i]
          g += data[i + 1]
          b += data[i + 2]
          n++
        }
        if (!n) {
          resolve({ ...FALLBACK })
          return
        }
        resolve(buildTheme(Math.round(r / n), Math.round(g / n), Math.round(b / n)))
      } catch {
        resolve({ ...FALLBACK })
      }
    }
    img.onerror = () => resolve({ ...FALLBACK })
    img.src = url
  })
}

export function useCoverTheme(coverSrc: Ref<string>) {
  const theme = ref<CoverTheme>({ ...FALLBACK })

  const themeStyle = computed(() => ({
    '--cover-bg': theme.value.bg,
    '--cover-fg': theme.value.fg,
    '--cover-fg-muted': theme.value.fgMuted,
    '--cover-fg-subtle': theme.value.fgSubtle,
    '--cover-control-fg': theme.value.controlFg,
    color: theme.value.fg,
  }))

  watch(
    coverSrc,
    (url) => {
      const token = ++sampleToken
      void sampleCover(url).then((t) => {
        if (token === sampleToken) theme.value = t
      })
    },
    { immediate: true },
  )

  return { theme, themeStyle }
}
