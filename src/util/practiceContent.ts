export type PracticeContentBlock =
  | { type: 'text'; text: string }
  | { type: 'image'; fileId: number; previewUrl?: string }

/** 解析 content 为可编辑块（编辑态展示图片） */
export function parseContentToBlocks(
  content: string,
  urlMap: Record<number, string> = {},
): PracticeContentBlock[] {
  if (!content) return [{ type: 'text', text: '' }]

  const blocks: PracticeContentBlock[] = []
  const re = /!\[(\d+)\]/g
  let last = 0
  let m: RegExpExecArray | null

  while ((m = re.exec(content)) !== null) {
    if (m.index > last) {
      blocks.push({ type: 'text', text: content.slice(last, m.index) })
    }
    const fileId = Number(m[1])
    blocks.push({ type: 'image', fileId, previewUrl: urlMap[fileId] })
    last = m.index + m[0].length
  }
  if (last < content.length) {
    blocks.push({ type: 'text', text: content.slice(last) })
  }
  if (!blocks.length) blocks.push({ type: 'text', text: '' })
  return blocks
}

/** 块序列还原为 content 字符串 */
export function blocksToContent(blocks: PracticeContentBlock[]): string {
  return blocks
    .map((b) => (b.type === 'text' ? b.text : `![${b.fileId}]`))
    .join('')
}

/** 从练习 content 中提取 ![file_id] 图片引用 */
export function extractImageIds(content: string): number[] {
  const ids = new Set<number>()
  const re = /!\[(\d+)\]/g
  let m: RegExpExecArray | null
  while ((m = re.exec(content)) !== null) {
    ids.add(Number(m[1]))
  }
  return [...ids]
}

/** 将纯文本 + ![id] 转为可渲染 HTML（需传入 id -> url 映射） */
export function renderPracticeContentHtml(content: string, urlMap: Record<number, string>): string {
  if (!content) return ''

  const parts: string[] = []
  const re = /!\[(\d+)\]/g
  let last = 0
  let m: RegExpExecArray | null

  while ((m = re.exec(content)) !== null) {
    const before = content.slice(last, m.index)
    if (before) parts.push(escapeHtml(before).replace(/\n/g, '<br>'))

    const id = Number(m[1])
    const url = urlMap[id]
    if (url) {
      parts.push(`<img src="${escapeAttr(url)}" alt="" class="prac-content-img" loading="lazy" />`)
    } else {
      parts.push(escapeHtml(m[0]))
    }
    last = m.index + m[0].length
  }

  const tail = content.slice(last)
  if (tail) parts.push(escapeHtml(tail).replace(/\n/g, '<br>'))
  return parts.join('')
}

function escapeHtml(s: string) {
  return s
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}

function escapeAttr(s: string) {
  return s.replace(/&/g, '&amp;').replace(/"/g, '&quot;')
}

/** 选中日的 0 点 unix（本地时区） */
export function dateKeyToPracticeAt(dateKey: string): number {
  const [y, m, d] = dateKey.split('-').map(Number)
  return Math.floor(new Date(y, m - 1, d).getTime() / 1000)
}

export function todayDateKey() {
  const t = new Date()
  return toDateKey(t.getFullYear(), t.getMonth(), t.getDate())
}

export function toDateKey(year: number, month: number, day: number) {
  return `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`
}
