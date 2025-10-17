// src/api/types/index.ts
export * from './music';
export * from './file';
export * from './artist';
export * from './album';
export * from './tag';
export * from './enums.ts';

export enum MidShowWhat {
    ShowMusicDetail = 1,
    ShowMusicLike = 2,
    ShowMemoList = 3,
    ShowFileList = 4,
    ShowChatWindow = 5,
}
