// 内容脚本 - 在页面中运行，用于提取音乐信息
class MusicInfoExtractor {
  constructor() {
    this.init();
  }

  init() {
    this.extractMusicInfo();
    this.setupMutationObserver();
  }

  extractMusicInfo() {
    // 提取页面中的音乐信息
    const musicInfo = {
      songName: this.extractSongName(),
      albumName: this.extractAlbumName(),
      artistName: this.extractArtistName(),
      coverUrl: this.extractCoverUrl(),
      lyrics: this.extractLyrics(),
      duration: this.extractDuration()
    };

    // 发送信息到background script
    if (this.hasValidMusicInfo(musicInfo)) {
      chrome.runtime.sendMessage({
        type: 'MUSIC_INFO_EXTRACTED',
        data: musicInfo
      });
    }
  }

  extractSongName() {
    // 尝试多种选择器提取歌曲名
    const selectors = [
      '.song-name',
      '.track-name',
      '.music-title',
      '.song-title',
      '[data-song-name]',
      '.title',
      'h1',
      'h2'
    ];

    for (const selector of selectors) {
      const element = document.querySelector(selector);
      if (element && element.textContent.trim()) {
        return element.textContent.trim();
      }
    }

    // 尝试从页面标题提取
    const title = document.title;
    if (title && !title.includes(' - ')) {
      return title;
    }

    return '';
  }

  extractAlbumName() {
    const selectors = [
      '.album-name',
      '.album-title',
      '.collection-name',
      '[data-album-name]',
      '.album'
    ];

    for (const selector of selectors) {
      const element = document.querySelector(selector);
      if (element && element.textContent.trim()) {
        return element.textContent.trim();
      }
    }

    return '';
  }

  extractArtistName() {
    const selectors = [
      '.artist-name',
      '.artist',
      '.singer',
      '.performer',
      '[data-artist-name]',
      '.author'
    ];

    for (const selector of selectors) {
      const element = document.querySelector(selector);
      if (element && element.textContent.trim()) {
        return element.textContent.trim();
      }
    }

    return '';
  }

  extractCoverUrl() {
    const selectors = [
      '.cover-image img',
      '.album-cover img',
      '.song-cover img',
      '.music-cover img',
      '[data-cover-image] img',
      '.thumbnail img'
    ];

    for (const selector of selectors) {
      const img = document.querySelector(selector);
      if (img && img.src) {
        return img.src;
      }
    }

    // 尝试从meta标签提取
    const metaImage = document.querySelector('meta[property="og:image"]');
    if (metaImage && metaImage.content) {
      return metaImage.content;
    }

    return '';
  }

  extractLyrics() {
    const selectors = [
      '.lyrics',
      '.song-lyrics',
      '.lyric-text',
      '[data-lyrics]',
      '.lyric'
    ];

    for (const selector of selectors) {
      const element = document.querySelector(selector);
      if (element && element.textContent.trim()) {
        return element.textContent.trim();
      }
    }

    return '';
  }

  extractDuration() {
    const selectors = [
      '.duration',
      '.song-duration',
      '.track-duration',
      '[data-duration]',
      '.time'
    ];

    for (const selector of selectors) {
      const element = document.querySelector(selector);
      if (element && element.textContent.trim()) {
        return this.parseDuration(element.textContent.trim());
      }
    }

    return 0;
  }

  parseDuration(durationText) {
    // 解析时长文本，如 "3:45" 或 "3分45秒"
    const timeRegex = /(\d+):(\d+)/;
    const match = durationText.match(timeRegex);
    
    if (match) {
      const minutes = parseInt(match[1]);
      const seconds = parseInt(match[2]);
      return minutes * 60 + seconds;
    }

    // 尝试解析中文格式
    const chineseRegex = /(\d+)分(\d+)秒/;
    const chineseMatch = durationText.match(chineseRegex);
    
    if (chineseMatch) {
      const minutes = parseInt(chineseMatch[1]);
      const seconds = parseInt(chineseMatch[2]);
      return minutes * 60 + seconds;
    }

    return 0;
  }

  hasValidMusicInfo(musicInfo) {
    return musicInfo.songName || musicInfo.albumName || musicInfo.artistName;
  }

  setupMutationObserver() {
    // 监听DOM变化，动态提取音乐信息
    const observer = new MutationObserver((mutations) => {
      let shouldExtract = false;
      
      mutations.forEach((mutation) => {
        if (mutation.type === 'childList' && mutation.addedNodes.length > 0) {
          // 检查是否有音乐相关的元素被添加
          mutation.addedNodes.forEach((node) => {
            if (node.nodeType === Node.ELEMENT_NODE) {
              const musicSelectors = [
                '.song-name', '.track-name', '.music-title',
                '.album-name', '.artist-name', '.lyrics'
              ];
              
              if (musicSelectors.some(selector => 
                node.matches && node.matches(selector) || 
                node.querySelector && node.querySelector(selector)
              )) {
                shouldExtract = true;
              }
            }
          });
        }
      });

      if (shouldExtract) {
        // 延迟提取，等待DOM稳定
        setTimeout(() => this.extractMusicInfo(), 500);
      }
    });

    observer.observe(document.body, {
      childList: true,
      subtree: true
    });
  }
}

// 页面加载完成后初始化
if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', () => {
    new MusicInfoExtractor();
  });
} else {
  new MusicInfoExtractor();
}