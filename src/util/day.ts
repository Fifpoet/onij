export function getWeekName(date: Date): string[] {
  const weekDays = ['星期天', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六'];
  const start = new Date(date);
  start.setDate(start.getDate() - start.getDay()); // 调整到本周的星期天

  const weekNames: string[] = [];
  for (let i = 0; i < 7; i++) {
    weekNames.push(weekDays[start.getDay()]);
    start.setDate(start.getDate() + 1);
  }

  return weekNames;
}

export function getMouthDay(date: Date): string[] {
  const start = new Date(date);
  start.setDate(start.getDate() - start.getDay()); // 调整到本周的星期天

  const mouthDays: string[] = [];
  for (let i = 0; i < 7; i++) {
    const month = start.getMonth() + 1; // getMonth() 返回的月份是从0开始的
    const day = start.getDate();
    mouthDays.push(`${month.toString().padStart(2, '0')}/${day.toString().padStart(2, '0')}`);
    start.setDate(start.getDate() + 1);
  }

  return mouthDays;
}

export function calculateTodayIndex(): number {
  return new Date().getDay();
}