// src/api/tag.ts
import { post } from '../util/http';
import { 
  UploadTagReq, 
  UploadTagResp,
  DeleteTagReq,
  DeleteTagResp
} from './types/tag';

// 上传标签
export const UploadTag = async (params: UploadTagReq): Promise<UploadTagResp> => {
  return (await post<UploadTagResp>('/tag/upload', params)).data;
};

// 删除标签
export const DeleteTag = async (params: DeleteTagReq): Promise<DeleteTagResp> => {
  return (await post<DeleteTagResp>('/tag/delete', params)).data;
};