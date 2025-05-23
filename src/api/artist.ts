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
  return post<UploadArtistResp>('/artist/upload', params);
};

// 搜索艺人
export const SearchArtist = async (params: SearchArtistReq): Promise<SearchArtistResp> => {
  return post<SearchArtistResp>('/artist/search', params);
};