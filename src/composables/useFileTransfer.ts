import { UploadFile, DeleteFileById, DownloadFiles } from '@/api/file'
import type { UploadFileResp } from '@/api/types'
import { FileType } from '@/api/types/enums'
import { uploadToQiniu, getFileTypeFromFile, downloadByPrivateUrl } from '@/util/qiniu'

/** 与 FileUploadModal / FileList 一致的上传、删除、下载流程 */
export function useFileTransfer(parentId = 0) {
  async function uploadOne(
    file: File,
    onProgress?: (percent: number) => void,
  ): Promise<{
    fileId: number
    storeKey: string
    name: string
    size: number
    format: FileType
    previewUrl: string
  }> {
    const qiniu = await uploadToQiniu(file, ['tran'], onProgress)

    const resp = await UploadFile({
      parent_id: parentId,
      files: [
        {
          filename: file.name,
          store_key: qiniu.key,
          hash: qiniu.hash,
          size: file.size,
          format: getFileTypeFromFile(file),
          exif: '',
        },
      ],
    })

    const ids = resp.file_ids ?? (resp as UploadFileResp & { fileIds?: number[] }).fileIds
    const fileId = ids?.[0]
    if (!fileId) throw new Error(resp.message || '登记文件失败')

    return {
      fileId,
      storeKey: qiniu.key,
      name: file.name,
      size: file.size,
      format: getFileTypeFromFile(file),
      previewUrl: qiniu.url,
    }
  }

  async function removeOne(fileId: number) {
    const resp = await DeleteFileById({ file_id: fileId })
    if (resp.message !== 'ok') {
      throw new Error(resp.message || '删除失败')
    }
  }

  /** 从后端取签名链接触发下载（同 FileList 的 iframe 方式） */
  async function downloadOne(fileId: number) {
    const resp = await DownloadFiles({ file_ids: [fileId] })
    const url = resp.urls?.[0]
    if (!url) throw new Error(resp.message || '获取下载链接失败')
    downloadByPrivateUrl(url)
  }

  return { uploadOne, removeOne, downloadOne }
}
