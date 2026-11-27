import {post} from '@/util';
import {GetArtistListReq, GetArtistListResp, UploadArtistReq, UploadArtistResp} from '@/api/types';

export const UploadArtist = async (params: UploadArtistReq): Promise<UploadArtistResp> => {
    return post<UploadArtistResp>('/artist/upload', params);
};

export const GetArtist = async (params: GetArtistListReq): Promise<GetArtistListResp> => {
    return post<GetArtistListResp>('/artist/search', params);
};