// src/api/types/file.ts

/* 枚举类型 */
export enum AppId {
  AppId_Unknown = 0,
  AppId_Music = 1,
}

export enum FileType {
  FT_Unknown = 0,
  FT_Mp3 = 1,
  FT_Lyrics = 2,
  FT_Png = 3,
}

/* 请求参数类型 */
export interface UploadFileReq {
  filename: string;       // 原始文件名（不含路径）
  app_id: AppId;          // 业务模块
  folders?: string[];     // 自定义文件夹
  file: File | Blob;      // 文件内容
}

/* 响应数据类型 */
export interface UploadFileResp {
  code: number;
  message: string;
  file_id: number;
  file_url: string;
}