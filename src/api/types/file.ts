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
  FT_Folder = 99,
}

/* 文件信息 */
export interface FileInfo {
  filename?: string;       // 文件名
  file: File;             // 文件内容
  origin_at?: number;     // 原始文件时间
}

/* 请求参数类型 */
export interface UploadFileReq {
  parent_id: number;      // 父级文件夹id
  files: FileInfo[];      // 文件列表
}

/* 响应数据类型 */
export interface UploadFileResp {
  code: number;
  message: string;
  file_ids: number[];
}

/* 获取文件列表请求参数 */
export interface GetFileListReq {
  parent_id: number;
  page: number;
  limit: number;
}

/* 获取文件列表响应数据 */
export interface GetFileListResp {
  code: number;
  message: string;
  files: FileDetail[];
  total: number;
}

/* 文件详情 */
export interface FileDetail {
  id: number;
  name: string;
  format: FileType;
  parent_id: number;
  origin_at: number;
  created_at: number;
  updated_at: number;
}

/* 下载文件请求参数 */
export interface DownloadFileReq {
  file_ids: number[];
}

/* 下载文件响应数据 */
export interface DownloadFileResp {
  code: number;
  message: string;
  urls: string[];
}