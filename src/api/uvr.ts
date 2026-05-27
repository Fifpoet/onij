/** UVR 伴奏服务（经 Go 后端代理） */
import axios from 'axios'
import { apiClient, post } from '../util/http'

export type UvrSeparateResponse = {
  job_id: string
  instrumental_path: string
  instrumental_filename: string
  all_outputs: string[]
  elapsed_sec: number
}

const UVR_SEPARATE_TIMEOUT_MS = 30 * 60 * 1000
const UVR_DOWNLOAD_TIMEOUT_MS = 10 * 60 * 1000

export async function separateInstrumentalFromUrl(
  sourceUrl: string,
  options?: {
    id?: string | number
    architecture?: string
    model?: string
    use_gpu?: boolean
  },
): Promise<UvrSeparateResponse> {
  try {
    const payload: Record<string, unknown> = {
      source_url: sourceUrl,
      architecture: options?.architecture ?? 'mdx',
      model: options?.model ?? 'UVR-MDX-NET-Inst_HQ_3.onnx',
      use_gpu: options?.use_gpu ?? true,
      output_format: 'WAV',
      single_stem: 'Instrumental',
    }
    if (options?.id != null && String(options.id) !== '') {
      payload.id = String(options.id)
    }
    return (await post('/uvr/instrumental/from-url', payload, {
      timeout: UVR_SEPARATE_TIMEOUT_MS,
    })) as UvrSeparateResponse
  } catch (err: unknown) {
    if (axios.isAxiosError(err)) {
      const detail = err.response?.data
      if (typeof detail === 'string' && detail) {
        throw new Error(detail)
      }
      if (detail && typeof detail === 'object' && 'detail' in detail && typeof detail.detail === 'string') {
        throw new Error(detail.detail)
      }
      throw new Error(`UVR 请求失败 (${err.response?.status ?? 'network'})`)
    }
    throw err
  }
}

export async function downloadInstrumentalBlob(jobId: string): Promise<Blob> {
  const resp = await apiClient.get(`/uvr/jobs/${encodeURIComponent(jobId)}/instrumental`, {
    responseType: 'blob',
    timeout: UVR_DOWNLOAD_TIMEOUT_MS,
  })
  return resp.data as Blob
}
