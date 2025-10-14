// 歌手详情 - 歌手信息
export interface ArtistDetail {
    img1v1Id: number
    topicPerson: number
    picId: number
    briefDesc: string
    albumSize: number
    musicSize: number
    picUrl: string
    img1v1Url: string
    followed: boolean
    trans: string
    alias: string[]
    name: string
    id: number
    publishTime: number
    picId_str: string
    img1v1Id_str: string
    mvSize: number
}

// 歌手详情 - 热门歌曲
export interface ArtistHotSong {
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
    rtype: number
    rurl: any
    pst: number
    alia: any[]
    pop: number
    rt: string
    mst: number
    cp: number
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
}

// 歌手详情 - 权限信息
export interface SongPrivilege {
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

