/** 从图片 EXIF 读取拍摄/创建时间（Unix 秒）；无 EXIF 或解析失败返回 undefined */
export async function readPhotoOriginAt(file: File): Promise<number | undefined> {
  if (!isLikelyPhoto(file)) return undefined

  try {
    const buf = await file.slice(0, 512 * 1024).arrayBuffer()
    const date = parseJpegExifDate(new DataView(buf))
    if (!date || Number.isNaN(date.getTime())) return undefined
    return Math.floor(date.getTime() / 1000)
  } catch {
    return undefined
  }
}

export function buildExifPayload(originAt?: number): string {
  if (!originAt || originAt <= 0) return ''
  return JSON.stringify({ origin_at: originAt })
}

function isLikelyPhoto(file: File): boolean {
  const mime = file.type.toLowerCase()
  if (mime.startsWith('image/')) return true
  return /\.(jpe?g|heic|tif{1,2})$/i.test(file.name)
}

function parseJpegExifDate(view: DataView): Date | undefined {
  if (view.byteLength < 4 || view.getUint16(0) !== 0xffd8) return undefined

  let offset = 2
  while (offset + 4 < view.byteLength) {
    if (view.getUint8(offset) !== 0xff) break
    const marker = view.getUint8(offset + 1)
    const len = view.getUint16(offset + 2)
    if (len < 2) break

    if (marker === 0xe1) {
      const exifBase = offset + 4
      if (exifBase + 8 > view.byteLength) break
      const sig = readAscii(view, exifBase, 4)
      if (sig === 'Exif') {
        const tiff = exifBase + 6
        const date =
          readTiffDate(view, tiff, 0x9003) ||
          readTiffDate(view, tiff, 0x0132) ||
          readTiffDate(view, tiff, 0x9004)
        if (date) return date
      }
    }

    offset += 2 + len
  }
  return undefined
}

function readTiffDate(view: DataView, tiff: number, tag: number): Date | undefined {
  if (tiff + 8 > view.byteLength) return undefined

  const le = view.getUint8(tiff) === 0x49 && view.getUint8(tiff + 1) === 0x49
  const be = view.getUint8(tiff) === 0x4d && view.getUint8(tiff + 1) === 0x4d
  if (!le && !be) return undefined

  const u16 = (o: number) => (le ? view.getUint16(o, true) : view.getUint16(o, false))
  const u32 = (o: number) => (le ? view.getUint32(o, true) : view.getUint32(o, false))

  const ifd0 = tiff + u32(tiff + 4)
  return walkIfd(view, tiff, ifd0, tag, u16, u32, 0)
}

function walkIfd(
  view: DataView,
  tiff: number,
  ifd: number,
  tag: number,
  u16: (o: number) => number,
  u32: (o: number) => number,
  depth: number,
): Date | undefined {
  if (depth > 4 || ifd + 2 > view.byteLength) return undefined

  const count = u16(ifd)
  let entry = ifd + 2
  for (let i = 0; i < count; i++) {
    if (entry + 12 > view.byteLength) break
    const entryTag = u16(entry)
    const type = u16(entry + 2)
    const countVal = u32(entry + 4)

    if (entryTag === tag && type === 2 && countVal >= 19) {
      const valueOffset = countVal > 4 ? tiff + u32(entry + 8) : entry + 8
      const raw = readAscii(view, valueOffset, Math.min(countVal - 1, 32))
      const parsed = parseExifDateString(raw)
      if (parsed) return parsed
    }

    if (entryTag === 0x8769 || entryTag === 0x8825) {
      const subIfd = tiff + u32(entry + 8)
      const sub = walkIfd(view, tiff, subIfd, tag, u16, u32, depth + 1)
      if (sub) return sub
    }

    entry += 12
  }
  return undefined
}

function parseExifDateString(raw: string): Date | undefined {
  const m = raw.match(/^(\d{4}):(\d{2}):(\d{2})[ T](\d{2}):(\d{2}):(\d{2})/)
  if (!m) return undefined
  const date = new Date(
    Number(m[1]),
    Number(m[2]) - 1,
    Number(m[3]),
    Number(m[4]),
    Number(m[5]),
    Number(m[6]),
  )
  return Number.isNaN(date.getTime()) ? undefined : date
}

function readAscii(view: DataView, offset: number, len: number): string {
  let s = ''
  for (let i = 0; i < len && offset + i < view.byteLength; i++) {
    const c = view.getUint8(offset + i)
    if (c === 0) break
    s += String.fromCharCode(c)
  }
  return s
}
