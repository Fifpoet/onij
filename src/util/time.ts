export function formatTime(seconds: number): string {
    const minutes = Math.floor(seconds / 60)
      .toString()
      .padStart(2, "0");
    const secondsPart = Math.floor(seconds % 60)
      .toString()
      .padStart(2, "0");
    return `${minutes}:${secondsPart}`;
  }