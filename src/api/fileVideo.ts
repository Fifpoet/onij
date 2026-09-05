import { post } from '../util/http'
import { normalizeFileCdnUrl } from '@/util/qiniuCdn'

const PLAY_TIMEOUT_MS = 60000

export type VideoPlayMode = 'mp4' | 'hls' | 'processing' | 'unsupported'

export type PlayFileVideoResp = {
  code: number
  message: string
  mode: VideoPlayMode
  url?: string
  playlist?: string
  duration_ms: number
}

export type FileVideoMarker = {
  id: number
  file_id: number
  time_ms: number
  label: string
}

type MarkerListResp = {
  code: number
  message: string
  markers: FileVideoMarker[]
}

type MarkerOkResp = {
  code: number
  message: string
}

function normalizePlaylist(text: string): string {
  return text.replace(/https?:\/\/cloud\.onij\.fun\/[^\s"']+/gi, (u) => normalizeFileCdnUrl(u))
}

export const PlayFileVideo = async (fileId: number): Promise<PlayFileVideoResp> => {
  const resp = (await post(
    '/file/video/play',
    { file_id: fileId },
    { timeout: PLAY_TIMEOUT_MS },
  )) as PlayFileVideoResp
  if (resp.url) resp.url = normalizeFileCdnUrl(resp.url)
  if (resp.playlist) resp.playlist = normalizePlaylist(resp.playlist)
  return resp
}

export const ListFileVideoMarkers = async (fileId: number): Promise<MarkerListResp> => {
  const resp = (await post('/file/video/marker/list', { file_id: fileId })) as MarkerListResp
  resp.markers = resp.markers ?? []
  return resp
}

export const SaveFileVideoMarkers = async (
  fileId: number,
  markers: Array<Pick<FileVideoMarker, 'time_ms' | 'label'> & { id?: number }>,
): Promise<MarkerListResp> => {
  const resp = (await post('/file/video/marker/save', {
    file_id: fileId,
    markers,
  })) as MarkerListResp
  resp.markers = resp.markers ?? []
  return resp
}

export const DeleteFileVideoMarker = async (id: number): Promise<MarkerOkResp> => {
  return post('/file/video/marker/delete', { id }) as Promise<MarkerOkResp>
}
