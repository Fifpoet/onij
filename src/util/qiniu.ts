import * as qiniu from 'qiniu-js';
import { FileType } from '@/api/types/file';
import { GetUploadToken } from '@/api/file';

// 七牛云配置
const QINIU_CONFIG = {
  bucket: 'onij',
  domain: 'http://cloud.onij.fun',
  zone: 'huadong' as const
};

// 获取上传token
async function getUploadToken(): Promise<string> {
  try {
    const response = await GetUploadToken();
    if (response.upload_token) {
      return response.upload_token;
    } else {
      throw new Error(response.message || 'Failed to get upload token');
    }
  } catch (error) {
    console.error('Failed to get upload token:', error);
    throw error;
  }
}

// 计算文件hash
async function calculateFileHash(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = async (e) => {
      try {
        const arrayBuffer = e.target?.result as ArrayBuffer;
        const hashBuffer = await crypto.subtle.digest('SHA-256', arrayBuffer);
        const hashArray = Array.from(new Uint8Array(hashBuffer));
        const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
        resolve(hashHex);
      } catch (error) {
        reject(error);
      }
    };
    reader.onerror = reject;
    reader.readAsArrayBuffer(file);
  });
}

// 生成唯一的文件名
function generateUniqueFileName(originalName: string): string {
  const timestamp = Date.now();
  const randomStr = Math.random().toString(36).substring(2, 10);
  const extension = originalName.includes('.') ? originalName.split('.').pop() : '';
  const nameWithoutExt = originalName.includes('.') ? originalName.substring(0, originalName.lastIndexOf('.')) : originalName;
  
  return `${nameWithoutExt}_${timestamp}_${randomStr}${extension ? '.' + extension : ''}`;
}

// 上传文件到七牛云
export async function uploadToQiniu(
  file: File, 
  folderPath: string[] = []
): Promise<{ key: string; hash: string; url: string }> {
  return new Promise(async (resolve, reject) => {
    try {
      // 生成唯一文件名
      const uniqueFileName = generateUniqueFileName(file.name);
      const key = folderPath.length > 0 ? `${folderPath.join('/')}/${uniqueFileName}` : uniqueFileName;
      
      // 从后端获取上传token
      const token = await getUploadToken();
      
      // 创建上传配置
      const config = {
        useCdnDomain: false,
        region: qiniu.region[QINIU_CONFIG.zone]
      };
      
      // 创建上传对象
      const observable = qiniu.upload(file, key, token, {}, config);
      
      // 监听上传进度
      observable.subscribe({
        next: (res) => {
          console.log('Upload progress:', res.total.percent);
        },
        error: (err) => {
          console.error('Upload error:', err);
          reject(err);
        },
        complete: (res) => {
          console.log('Upload completed:', res);
          // 计算文件hash
          calculateFileHash(file).then(hash => {
            resolve({
              key: res.key,
              hash: hash,
              url: `${QINIU_CONFIG.domain}/${res.key}`
            });
          }).catch(reject);
        }
      });
    } catch (error) {
      reject(error);
    }
  });
}

// 获取文件类型
export function getFileTypeFromFile(file: File): FileType {
  const mimeType = file.type.toLowerCase();
  const extension = file.name.split('.').pop()?.toLowerCase() || '';
  
  // 根据MIME类型判断
  if (mimeType.startsWith('audio/') || extension === 'mp3') {
    return FileType.FT_Mp3;
  }
  if (mimeType.startsWith('image/') || ['png', 'jpg', 'jpeg', 'gif', 'bmp'].includes(extension)) {
    return FileType.FT_Png;
  }
  if (mimeType.startsWith('text/') || ['txt', 'lrc'].includes(extension)) {
    return FileType.FT_Lyrics;
  }
  
  return FileType.FT_Unknown;
}

// 批量上传文件
export async function batchUploadToQiniu(
  files: File[], 
  folderPath: string[] = []
): Promise<Array<{ key: string; hash: string; url: string; file: File }>> {
  const uploadPromises = files.map(file => 
    uploadToQiniu(file, folderPath).then(result => ({
      ...result,
      file
    }))
  );
  
  return Promise.all(uploadPromises);
} 