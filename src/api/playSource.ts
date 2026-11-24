import { post } from '../util/http'

export type PlaySourceKind = 'gdstudio' | 'gdapi'

export type PlaySourceDTO = {
  id: string
  title: string
  kind: PlaySourceKind
  origin: string
  enabled: boolean
  builtin?: boolean
}

type PlaySourceListResp = {
  code: number
  message: string
  sources: PlaySourceDTO[]
}

export const ListPlaySources = async (): Promise<PlaySourceListResp> => {
  return post('/play-source/list', {}) as Promise<PlaySourceListResp>
}

export const SavePlaySources = async (sources: PlaySourceDTO[]): Promise<PlaySourceListResp> => {
  return post('/play-source/save', { sources }) as Promise<PlaySourceListResp>
}

export type GdstudioSignResp = {
  code: number
  message: string
  version: string
  crc32: string
  updatedAt: number
}

const SIGN_TIMEOUT = 20_000

export const GetGdstudioSign = async (): Promise<GdstudioSignResp> => {
  return post('/play-source/gdstudio-sign', {}, { timeout: SIGN_TIMEOUT }) as Promise<GdstudioSignResp>
}

export const SaveGdstudioSign = async (payload: {
  version: string
  crc32: string
}): Promise<GdstudioSignResp> => {
  return post('/play-source/gdstudio-sign/save', payload, { timeout: SIGN_TIMEOUT }) as Promise<GdstudioSignResp>
}
