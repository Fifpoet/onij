// 专辑详情 - 艺术家信息
export interface AlbumArtist {
    img1v1Id: number
    topicPerson: number
    picId: number
    musicSize: number
    albumSize: number
    briefDesc: string
    picUrl: string
    img1v1Url: string
    followed: boolean
    trans: string
    alias: string[]
    name: string
    id: number
    img1v1Id_str?: string
    picId_str?: string
}

// 专辑详情 - 评论信息
export interface AlbumCommentInfo {
    commentThread: {
        id: string
        resourceInfo: {
            id: number
            userId: number
            name: string
            imgUrl: string
            creator: any
            encodedId: any
            subTitle: any
            webUrl: any
        }
        resourceType: number
        commentCount: number
        likedCount: number
        shareCount: number
        hotCount: number
        latestLikedUsers: any
        resourceId: number
        resourceOwnerId: number
        resourceTitle: string
    }
    latestLikedUsers: any
    liked: boolean
    comments: any
    resourceType: number
    resourceId: number
    commentCount: number
    likedCount: number
    shareCount: number
    threadId: string
}

// 专辑详情 - 专辑信息
export interface AlbumDetail {
    songs: any[]
    paid: boolean
    onSale: boolean
    mark: number
    awardTags: any
    displayTags: any
    artists: AlbumArtist[]
    copyrightId: number
    picId: number
    artist: AlbumArtist
    briefDesc: string
    publishTime: number
    company: string
    picUrl: string
    commentThreadId: string
    blurPicUrl: string
    companyId: number
    pic: number
    status: number
    subType: string
    alias: string[]
    description: string
    tags: string
    name: string
    id: number
    type: string
    size: number
    info: AlbumCommentInfo
    picId_str?: string
}

// 专辑详情 - 歌曲信息（与 SongDetail 类似，但包含 privilege）
export interface AlbumSong {
    rtUrls: any[]
    ar: Array<{
        id: number
        name: string
        alia: string[]
    }>
    al: {
        id: number
        name: string
        pic_str: string
        pic: number
        alia: string[]
    }
    st: number
    noCopyrightRcmd: any
    songJumpInfo: any
    djId: number
    no: number
    fee: number
    mv: number
    cd: string
    t: number
    v: number
    rtype: number
    rurl: any
    pst: number
    alia: any[]
    pop: number
    rt: any
    mst: number
    cp: number
    crbt: any
    cf: string
    dt: number
    h: {
        br: number
        fid: number
        size: number
        vd: number
        sr: number
    }
    sq: {
        br: number
        fid: number
        size: number
        vd: number
        sr: number
    }
    hr: any
    l: {
        br: number
        fid: number
        size: number
        vd: number
        sr: number
    }
    rtUrl: any
    ftype: number
    a: any
    m: {
        br: number
        fid: number
        size: number
        vd: number
        sr: number
    }
    name: string
    id: number
    videoInfo?: {
        moreThanOne: boolean
        video: {
            vid: string
            type: number
            title: string
            playTime: number
            coverUrl: string
            publishTime: number
            artists: any
        }
    }
    privilege: {
        id: number
        fee: number
        payed: number
        st: number
        pl: number
        dl: number
        sp: number
        cp: number
        subp: number
        cs: boolean
        maxbr: number
        fl: number
        toast: boolean
        flag: number
        preSell: boolean
        playMaxbr: number
        downloadMaxbr: number
        maxBrLevel: string
        playMaxBrLevel: string
        downloadMaxBrLevel: string
        plLevel: string
        dlLevel: string
        flLevel: string
        rscl: any
        freeTrialPrivilege: {
            resConsumable: boolean
            userConsumable: boolean
            listenType: any
            cannotListenReason: any
            playReason: any
            freeLimitTagType: any
        }
        rightSource: number
        chargeInfoList: Array<{
            rate: number
            chargeUrl: any
            chargeMessage: any
            chargeType: number
        }>
        code: number
        message: any
        plLevels: any
        dlLevels: any
        ignoreCache: any
        bd: any
    }
}

// 专辑详情响应
export interface AlbumDetailResponse {
    resourceState: boolean
    songs: AlbumSong[]
    code: number
    album: AlbumDetail
}

// 获取专辑详情
import { getNetease } from '@/util/http'

export async function getAlbumDetail(id: string | number): Promise<AlbumDetailResponse> {
    return await getNetease<AlbumDetailResponse>('/album', {
        id: id
    })
}

