// 视图-音乐列表项

export interface ViewMusicListItem {
    id: number;
    name: string;
    artists: ({
        artist_id: number;
        artist_name: string;
    })[];
    /** 完整表演者（可选；专辑页列表仍用 artists 表示「非专辑歌手」，队列/搜索样式可用此字段展示全员） */
    display_artists?: {
        artist_id: number;
        artist_name: string;
    }[];
    time_long: number;
    composer_id?: number;
    composer_name?: string;
    writer_id?: number;
    writer_name?: string;
    issue_time?: number;
    mv_url?: string;

    // 专辑信息
    album_id: number;
    album_name: string;
    cover_file_url?: string;
    // 标签信息
    tags?: string[];
}
