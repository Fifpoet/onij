// src/api/artist.ts
import {post} from '@/util';
import {GetArtistListReq, GetArtistListResp, UploadArtistReq, UploadArtistResp} from '@/api/types';

// 上传艺术家
export const UploadArtist = async (params: UploadArtistReq): Promise<UploadArtistResp> => {
    return (await post<UploadArtistResp>('/artist/upload', params)).data;
};

// 搜索艺术家
export const GetArtist = async (params: GetArtistListReq): Promise<GetArtistListResp> => {
    return (await post<GetArtistListResp>('/artist/search', params)).data;
};