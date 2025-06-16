// src/api/artist.ts
import { post } from '../util/http';
import { 
  UploadArtistReq, 
  UploadArtistResp,
  SearchArtistReq,
  SearchArtistResp 
} from './types/artist';

// 上传艺术家
export const UploadArtist = async (params: UploadArtistReq): Promise<UploadArtistResp> => {
  return (await post<UploadArtistResp>('/artist/upload', params)).data;
};

// 搜索艺术家
export const SearchArtist = async (params: SearchArtistReq): Promise<SearchArtistResp> => {
  return (await post<SearchArtistResp>('/artist/search', params)).data;
};