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
  const today = new Date();
  // getDay() 返回 0 表示星期天，1 表示星期一，以此类推，直到 6 表示星期六
  // 我们通过取余数确保索引在 0 到 6 的范围内
  const dayOfWeek = today.getDay();
  
  // 根据你的应用设置，你可能需要调整 dayOfWeek 的值，例如：
  // 如果你想将 0 作为星期天的索引，而 6 作为星期六的索引，你不需要做任何改变
  // 如果你想将 1 作为星期一的索引，0 作为星期天的索引，你可以使用以下调整：
  // const adjustedDayOfWeek = (dayOfWeek + 5) % 7 + 1;
  
  return dayOfWeek; // 或者使用 adjustedDayOfWeek
}