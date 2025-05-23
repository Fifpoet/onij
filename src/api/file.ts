// src/api/file.ts
import { post } from '../util/http';
import { UploadFileReq, UploadFileResp } from './types/file';

// 上传文件
export const UploadFile = async (params: UploadFileReq): Promise<UploadFileResp> => {
  const formData = new FormData();
  formData.append('filename', params.filename);
  formData.append('app_id', params.app_id.toString());
  
  if (params.folders && params.folders.length > 0) {
    params.folders.forEach((folder, index) => {
      formData.append(`folders[${index}]`, folder);
    });
  }
  
  formData.append('file', params.file);
  
  return (await post<UploadFileResp>('/file/upload', formData)).data;
};