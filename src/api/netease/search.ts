// 导入工具函数和类型
import { getNetease } from '@/util/http'
import type { SongDetailResponse } from './result'

// 歌曲搜索 - 艺术家
export interface SongSearchArtist {
    id: number
    name: string
    picUrl: string | null
    alias: string[]
    albumSize: number
    picId: number
    fansGroup: number | null
    img1v1Url: string
    img1v1: number
    trans: string | null
}

// 歌曲搜索 - 专辑
export interface SongSearchAlbum {
    id: number
    name: string
    artist: SongSearchArtist
    publishTime: number
    size: number
    copyrightId: number
    status: number
    picId: number
    mark: number
}

// 歌曲搜索 - 歌曲
export interface SongSearchItem {
    id: number
    name: string
    artists: SongSearchArtist[]
    album: SongSearchAlbum
    duration: number
    copyrightId: number
    status: number
    alias: string[]
    rtype: number
    ftype: number
    mvid: number
    fee: number
    rUrl: string | null
    mark: number
}

// 歌曲搜索结果
export interface SearchResult {
    hlWords: string[]
    hasMore: boolean
    songs: SongSearchItem[]
    songCount: number

    albums: AlbumSearchItem[]
    albumCount: number

    artists: ArtistSearchItem[]
    artistCount: number
    searchQcReminder: any | null
}


// 专辑搜索 - 艺术家
export interface AlbumSearchArtist {
    name: string
    id: number
    picId: number
    img1v1Id: number
    briefDesc: string
    picUrl: string
    img1v1Url: string
    albumSize: number
    alias: string[]
    trans: string
    musicSize: number
    topicPerson: number
    picId_str?: string
    img1v1Id_str?: string
    alia?: string[]
}

// 专辑搜索 - 专辑艺术家
export interface AlbumSearchAlbumArtist {
    name: string
    id: number
    picId: number
    img1v1Id: number
    briefDesc: string
    picUrl: string
    img1v1Url: string
    albumSize: number
    alias: string[]
    trans: string
    musicSize: number
    topicPerson: number
    img1v1Id_str?: string
}

// 专辑搜索 - 专辑
export interface AlbumSearchItem {
    name: string
    id: number
    type: string
    size: number
    picId: number
    blurPicUrl: string
    companyId: number
    pic: number
    picUrl: string
    publishTime: number
    description: string
    tags: string
    company: string
    briefDesc: string
    artist: AlbumSearchArtist
    songs: any[] | null
    alias: string[]
    status: number
    copyrightId: number
    commentThreadId: string
    artists: AlbumSearchAlbumArtist[]
    paid: boolean
    onSale: boolean
    picId_str: string
    alg: string
    mark: number
    containedSong: string
}


// 艺术家搜索 - 艺术家
export interface ArtistSearchItem {
    id: number
    name: string
    picUrl: string | null
    alias: string[]
    albumSize: number
    picId: number
    fansGroup: number | null
    img1v1Url: string
    accountId?: number
    img1v1: number
    identityIconUrl?: string
    mvSize: number
    followed: boolean
    alg: string
    alia?: string[]
    trans: string | null
}

// 获取歌曲详情
export async function getSongDetail(ids: number[]): Promise<SongDetailResponse> {
    const idsStr = ids.join(',')
    return await getNetease<SongDetailResponse>('/song/detail', {
        ids: idsStr
    })
}
