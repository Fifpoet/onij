// src/api/artist.ts
import { post } from '../util/http';
import { UploadArtistReq, UploadArtistResp } from './types/artist';

// 上传艺术家
export const UploadArtist = async (params: UploadArtistReq): Promise<UploadArtistResp> => {
  return post<UploadArtistResp>('/artist/upload', params);
};