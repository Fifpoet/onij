// 背景脚本 - 监听网络请求和处理音乐信息

class RequestListener {
  constructor() {
    this.init();
  }

  init() {
    // 监听所有网络请求
    chrome.webRequest.onBeforeRequest.addListener(
      (details) => {
        this.handleRequest(details);
      },
      {
        urls: ["https://zm.i9mr.com/song/detail*"]
      },
      ["requestBody"]
    );

    // 监听来自content script的消息
    chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
      if (message.type === 'MUSIC_INFO_EXTRACTED') {
        this.handleMusicInfo(message.musicInfo);
      }
    });

    console.log('音乐信息提取器已启动，正在监听请求...');
  }

  handleRequest(details) {
    console.log('捕获到音乐详情请求:', details.url);
    
    // 提取请求信息
    const requestInfo = {
      id: Date.now(),
      url: details.url,
      timestamp: new Date().toISOString(),
      method: details.method,
      tabId: details.tabId
    };

    // 保存到存储
    this.saveRequest(requestInfo);
  }

  handleMusicInfo(musicInfo) {
    console.log('🎵 Background收到音乐信息:', musicInfo);
    
    // 保存音乐信息到存储
    this.saveMusicInfo(musicInfo);
  }

  async saveMusicInfo(musicInfo) {
    try {
      // 获取现有数据
      const result = await chrome.storage.local.get(['musicRequests']);
      const requests = result.musicRequests || [];
      
      // 查找对应的请求并更新音乐信息
      const requestIndex = requests.findIndex(req => req.url === musicInfo.url);
      if (requestIndex !== -1) {
        requests[requestIndex].musicInfo = musicInfo;
      } else {
        // 如果没有找到对应请求，创建新的
        const newRequest = {
          id: Date.now(),
          url: musicInfo.url,
          timestamp: musicInfo.timestamp,
          method: 'GET',
          musicInfo: musicInfo
        };
        requests.unshift(newRequest);
      }
      
      // 限制保存数量（最多100条）
      if (requests.length > 100) {
        requests.splice(100);
      }
      
      // 保存到存储
      await chrome.storage.local.set({ musicRequests: requests });
      
      console.log('💾 音乐信息已保存到存储:', musicInfo);
      console.log('📊 当前存储的请求数量:', requests.length);
      
      // 通知popup更新
      this.notifyPopupUpdate();
      
    } catch (error) {
      console.error('保存音乐信息失败:', error);
    }
  }

  async saveRequest(requestInfo) {
    try {
      // 获取现有数据
      const result = await chrome.storage.local.get(['musicRequests']);
      const requests = result.musicRequests || [];
      
      // 添加新请求
      requests.unshift(requestInfo);
      
      // 限制保存数量（最多100条）
      if (requests.length > 100) {
        requests.splice(100);
      }
      
      // 保存到存储
      await chrome.storage.local.set({ musicRequests: requests });
      
      console.log('请求已保存:', requestInfo);
      
      // 通知popup更新
      this.notifyPopupUpdate();
      
    } catch (error) {
      console.error('保存请求失败:', error);
    }
  }

  notifyPopupUpdate() {
    // 通知popup页面更新
    chrome.runtime.sendMessage({ 
      type: 'REQUEST_UPDATE',
      timestamp: Date.now()
    }).catch(() => {
      // 忽略popup未打开时的错误
    });
  }
}

// 初始化请求监听器
new RequestListener();