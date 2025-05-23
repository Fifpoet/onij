// src/api/tag.ts
import { post } from '../util/http';
import { UploadTagReq, UploadTagResp } from './types/tag';

// 上传标签
export const UploadTag = async (params: UploadTagReq): Promise<UploadTagResp> => {
  return post<UploadTagResp>('/tag/upload', params);
};