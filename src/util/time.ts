export function formatTime(seconds: number): string {
    const minutes = Math.floor(seconds / 60)
      .toString()
      .padStart(2, "0");
    const secondsPart = Math.floor(seconds % 60)
      .toString()
      .padStart(2, "0");
    return `${minutes}:${secondsPart}`;
}

/**
 * 格式化时间戳为日期字符串 (YYYY-MM-DD)
 * @param timestamp 时间戳（毫秒）
 * @returns 格式化后的日期字符串，如果时间戳无效则返回 '未知'
 */
export function formatDate(timestamp: number): string {
    if (!timestamp) return '未知'
    const date = new Date(timestamp)
    if (isNaN(date.getTime())) return '未知'
    return date.getFullYear() + '-' + 
           String(date.getMonth() + 1).padStart(2, '0') + '-' + 
           String(date.getDate()).padStart(2, '0')
}