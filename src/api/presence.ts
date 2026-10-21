import { post } from '../util/http'
import type { ViewMusicListItem } from '@/api/view/music'

export type PresenceTrack = ViewMusicListItem

export type PresencePeerDTO = {
  client_id: string
  nickname: string
  now_playing: PresenceTrack | null
  queue: PresenceTrack[]
  playback_current_sec: number
  is_playing: boolean
}

type HeartbeatResp = {
  code: number
  message: string
  peers: PresencePeerDTO[]
}

export const PostPresenceHeartbeat = async (
  params: {
    client_id: string
    nickname: string
    now_playing: PresenceTrack | null
    queue: PresenceTrack[]
    playback_current_sec: number
    is_playing: boolean
  },
  timeoutMs = 4000,
): Promise<HeartbeatResp> => {
  return post('/presence/heartbeat', params, { timeout: timeoutMs }) as Promise<HeartbeatResp>
}
