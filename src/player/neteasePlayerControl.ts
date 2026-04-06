/**
 * 全局音频由 NeteaseAudioHost 绑定；抽屉/顶栏通过此处触发播放暂停与进度跳转。
 */
type ToggleFn = () => void
type SeekPercentFn = (percent: number) => void

let togglePlayImpl: ToggleFn | null = null
let seekPercentImpl: SeekPercentFn | null = null

export const neteasePlayerControl = {
  bindTogglePlay(fn: ToggleFn) {
    togglePlayImpl = fn
  },
  togglePlay() {
    togglePlayImpl?.()
  },
  bindSeekPercent(fn: SeekPercentFn) {
    seekPercentImpl = fn
  },
  /** 0–100，对应整首进度 */
  seekToPercent(percent: number) {
    seekPercentImpl?.(percent)
  },
}
