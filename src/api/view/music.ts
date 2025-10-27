// 视图-音乐列表项

export interface ViewMusicListItem {
    id: number;
    name: string;
    artists: ({
        artist_id: number;
        artist_name: string;
    })[];
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
