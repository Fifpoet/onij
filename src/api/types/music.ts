import { AudioQuality, PerformType } from './enums'
import { Artist, Album, Tag } from './index'

export interface Music {
    id: number
    name: string
    artists: Artist[]
    composers: Artist[]
    writers: Artist[]
    issueTime: number
    mvUrl: string
    audioFileUrl: string
    lyricsFileUrl: string
    album: Album[]
    tags: Tag[]
}

// 请求和响应类型
export interface UploadMusicReq {
    // 基础信息
    name: string
    subName: string
    artistIds: number[]
    audioFileId: number
    lyricFileId: number
    timeLength: number

    // 表单信息
    composerIds: number[]
    writerIds: number[]
    albumId?: number
    mvUrl?: string
    rootMusicId?: number
    issueTime?: number
    priority?: number
    lyricContent?: string
    audioQuality: AudioQuality
    performType: PerformType
    thirdId?: number

    // update
    musicId?: number
}

export interface UploadMusicResp {
    code: number
    message: string
    musicId: number
}

export interface GetMusicDetailReq {
    musicId: number
}

export interface GetMusicDetailResp {
    code: number
    message: string
    music: Music
}

export interface UpdateMusicMvReq {
    music_id?: number
    third_id?: number
    mv_url: string
}

export interface UpdateMusicMvResp {
    code: number
    message: string
    music_id?: number
}

export interface GetMusicMvBatchReq {
    third_ids: number[]
}

export interface MusicMvItem {
    third_id: number
    music_id: number
    mv_url: string
}

export interface GetMusicMvBatchResp {
    code: number
    message: string
    items?: MusicMvItem[]
}

export interface GetMusicListReq {
    keyword: string
    artistIds: number[]
    tagTypes: number[]

    page: number
    limit: number
}

export interface GetMusicListResp {
    code: number
    message: string
    musics: Music[]
}