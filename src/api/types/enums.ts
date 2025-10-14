export enum FileType {
    FT_Unknown = 0,
    FT_Pdf = 1,
    FT_Word = 2,
    FT_Excel = 3,
    FT_PowerPoint = 4,
    FT_Image = 5,
    FT_Video = 6,
    FT_Audio = 7,
    FT_Text = 8,
    FT_Archive = 9,
    FT_Code = 10,
    FT_CSV = 11,
    FT_Executable = 12,
    FT_Folder = 99,
}

export enum ArtistType {
    ATT_Unknown = 0,
    ATT_Artist = 1,
    ATT_Writer = 2,
    ATT_Composer = 3,
}

export enum AlbumType {
    ALT_Unknown = 0,
    ALT_Studio = 1,
    ALT_Live = 2,
    ALT_OST = 3,
    ALT_Cover = 4,
    ALT_Other = 99,
}

export enum AudioQuality {
    AQ_Unknown = 0,
    AQ_Low = 1,
}

export enum PerformType {
    PT_Unknown = 0,
}

export enum TagBiz {
    TB_Unknown = 0,
    TB_Music = 1,
}

export enum TagGroup {
    TG_Unknown = 0,
    TG_MusicTheme = 101,
    TG_MusicStyle = 102,
    TG_MusicEmotion = 103,
    TG_MusicHighlight = 104,
    TG_MusicCollection = 105,
    TG_ArtistStar = 120,
}

export enum TagType {
    TT_Unknown = 0,
    TT_NcBridge = 10401,
    TT_Star = 12001,
}

export enum ResourceType {
    RT_Unknown = 0,
    RT_Artist = 1,
    RT_Music = 2,
    RT_Album = 3,
}