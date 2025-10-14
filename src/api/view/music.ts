// 视图-音乐列表项

import {Music} from "@/api/types";

export interface ViewMusicListItem {
    id: number;
    name: string;
    artists: ({
        artist_id: number;
        artist_name: string;
    })[];
    composer_id: number;
    composer_name: string;
    writer_id: number;
    writer_name: string;
    issue_time: number;
    mv_url: string;

    mp3_file_url: string;
    lyrics_file_url: string;
    lyrics_content: string;

    // 专辑信息
    album_id: number;
    album_name: string;
    cover_file_url: string;
    // 标签信息
    tags: string[];
}

export function MusicDetailToViewItem(music: Music) : ViewMusicListItem {
    return {
        id: music.id,
        name: music.name,
    }
}