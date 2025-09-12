// Content Script - 拦截API响应并提取音乐信息

class MusicInfoExtractor {
  constructor() {
    this.init();
  }

  init() {
    // 拦截fetch请求
    this.interceptFetch();
    
    // 拦截XMLHttpRequest
    this.interceptXHR();
    
    console.log('音乐信息提取器 Content Script 已启动');
  }

  interceptFetch() {
    const originalFetch = window.fetch;
    const self = this;
    
    window.fetch = async function(...args) {
      const response = await originalFetch.apply(this, args);
      
      // 记录所有请求以便调试
      if (args[0] && typeof args[0] === 'string') {
        console.log('🌐 拦截到请求:', args[0]);
        
        // 检查是否是音乐详情API
        if (args[0].includes('/song/detail') || args[0].includes('song/detail') || 
            args[0].includes('detail') || args[0].includes('song')) {
          console.log('🎵 拦截到可能的音乐API请求:', args[0]);
          self.handleApiResponse(response.clone(), args[0]);
        }
      }
      
      return response;
    };
  }

  interceptXHR() {
    const originalOpen = XMLHttpRequest.prototype.open;
    const originalSend = XMLHttpRequest.prototype.send;
    const self = this;

    XMLHttpRequest.prototype.open = function(method, url, ...args) {
      this._url = url;
      return originalOpen.apply(this, [method, url, ...args]);
    };

    XMLHttpRequest.prototype.send = function(...args) {
      const xhr = this;
      
      // 监听响应
      xhr.addEventListener('load', function() {
        if (xhr._url) {
          console.log('🌐 XHR请求完成:', xhr._url);
          
          if (xhr._url.includes('/song/detail') || xhr._url.includes('song/detail')) {
            console.log('🎵 拦截到音乐详情XHR请求:', xhr._url);
            self.handleXHRResponse(xhr);
          }
        }
      });
      
      return originalSend.apply(this, args);
    };
  }

  async handleApiResponse(response, url) {
    try {
      console.log('🔍 开始解析API响应:', url);
      console.log('🔍 响应状态:', response.status);
      console.log('🔍 响应头:', response.headers);
      
      const data = await response.json();
      console.log('📄 API响应数据:', data);
      console.log('📄 数据类型:', typeof data);
      console.log('📄 数据键:', Object.keys(data || {}));
      
      this.extractMusicInfo(data, url);
    } catch (error) {
      console.error('❌ 解析API响应失败:', error);
      console.error('❌ 错误详情:', error.message);
    }
  }

  handleXHRResponse(xhr) {
    try {
      console.log('🔍 开始解析XHR响应:', xhr._url);
      console.log('🔍 XHR状态:', xhr.status);
      console.log('🔍 XHR响应文本长度:', xhr.responseText.length);
      
      const data = JSON.parse(xhr.responseText);
      console.log('📄 XHR响应数据:', data);
      console.log('📄 数据类型:', typeof data);
      console.log('📄 数据键:', Object.keys(data || {}));
      
      this.extractMusicInfo(data, xhr._url);
    } catch (error) {
      console.error('❌ 解析XHR响应失败:', error);
      console.error('❌ 错误详情:', error.message);
      console.error('❌ 响应文本:', xhr.responseText.substring(0, 500));
    }
  }

  extractMusicInfo(data, url) {
    try {
      console.log('🎼 开始提取音乐信息...');
      console.log('🎼 完整数据结构:', JSON.stringify(data, null, 2));
      
      let song = null;
      
      // 尝试不同的数据结构
      if (data && data.songs && Array.isArray(data.songs) && data.songs.length > 0) {
        song = data.songs[0];
        console.log('🎵 从songs数组中找到歌曲:', song);
      } else if (data && data.data && data.data.songs && Array.isArray(data.data.songs) && data.data.songs.length > 0) {
        song = data.data.songs[0];
        console.log('🎵 从data.songs数组中找到歌曲:', song);
      } else if (data && data.result && data.result.songs && Array.isArray(data.result.songs) && data.result.songs.length > 0) {
        song = data.result.songs[0];
        console.log('🎵 从result.songs数组中找到歌曲:', song);
      } else if (data && data.song) {
        song = data.song;
        console.log('🎵 从song对象中找到歌曲:', song);
      } else {
        console.log('❌ 未找到歌曲信息，尝试直接使用data作为歌曲信息');
        console.log('❌ 数据结构分析:');
        console.log('  - data存在:', !!data);
        console.log('  - data.songs存在:', !!(data && data.songs));
        console.log('  - data.songs是数组:', !!(data && data.songs && Array.isArray(data.songs)));
        console.log('  - data.songs长度:', data && data.songs ? data.songs.length : 'N/A');
        console.log('  - data键:', data ? Object.keys(data) : 'N/A');
        
        // 如果数据结构完全不同，尝试直接使用data
        if (data && typeof data === 'object') {
          song = data;
        } else {
          return;
        }
      }

      if (!song) {
        console.log('❌ 无法找到歌曲信息');
        return;
      }

      console.log('🎵 最终歌曲数据:', song);
      
      // 提取歌曲信息 - 支持多种字段名
      const musicInfo = {
        songName: song.name || song.title || song.songName || '未知歌曲',
        artists: this.extractArtists(song),
        album: song.al ? song.al.name : (song.album ? song.album.name : song.albumName) || '未知专辑',
        songId: song.id || song.songId || song.song_id,
        albumCover: this.extractAlbumCover(song),
        url: url,
        timestamp: new Date().toISOString(),
        rawData: song // 保存原始数据用于调试
      };

      console.log('✅ 提取到音乐信息:', musicInfo);

      // 发送到background script
      chrome.runtime.sendMessage({
        type: 'MUSIC_INFO_EXTRACTED',
        musicInfo: musicInfo
      }, (response) => {
        if (chrome.runtime.lastError) {
          console.error('❌ 发送消息失败:', chrome.runtime.lastError);
        } else {
          console.log('✅ 音乐信息已发送到background script');
        }
      });

    } catch (error) {
      console.error('❌ 提取音乐信息失败:', error);
      console.error('❌ 错误堆栈:', error.stack);
    }
  }

  extractArtists(song) {
    if (song.ar && Array.isArray(song.ar)) {
      return song.ar.map(artist => artist.name || artist).join(', ');
    } else if (song.artists && Array.isArray(song.artists)) {
      return song.artists.map(artist => artist.name || artist).join(', ');
    } else if (song.artist) {
      return song.artist;
    } else if (song.artistName) {
      return song.artistName;
    }
    return '未知歌手';
  }

  extractAlbumCover(song) {
    if (song.al && song.al.picUrl) {
      return song.al.picUrl;
    } else if (song.album && song.album.picUrl) {
      return song.album.picUrl;
    } else if (song.albumCover) {
      return song.albumCover;
    } else if (song.cover) {
      return song.cover;
    }
    return null;
  }
}

// 初始化音乐信息提取器
new MusicInfoExtractor();