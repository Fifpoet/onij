import { getNetease } from '@/util/http'

/** 网易云 lyric 接口常见结构 */
export interface LyricLrcBlock {
  version?: number
  lyric?: string
}

export interface LyricApiResponse {
  code?: number
  lrc?: LyricLrcBlock
  tlyric?: LyricLrcBlock
  yrc?: LyricLrcBlock
}

export async function fetchLyricRaw(id: number): Promise<string> {
  const data = await getNetease<LyricApiResponse>('/lyric', { id })
  const raw = data?.lrc?.lyric
  return typeof raw === 'string' ? raw : ''
}
