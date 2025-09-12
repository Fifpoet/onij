// 后台脚本 - 处理网络请求拦截
class RequestInterceptor {
  constructor() {
    this.init();
  }

  init() {
    // 监听网络请求
    chrome.webRequest.onBeforeRequest.addListener(
      (details) => {
        this.handleRequest(details);
      },
      { urls: ["<all_urls>"] },
      ["requestBody"]
    );

    // 监听请求完成
    chrome.webRequest.onCompleted.addListener(
      (details) => {
        this.handleResponse(details);
      },
      { urls: ["<all_urls>"] },
      ["responseHeaders"]
    );

    // 监听标签页更新
    chrome.tabs.onUpdated.addListener((tabId, changeInfo, tab) => {
      if (changeInfo.status === 'complete') {
        this.injectContentScript(tabId);
      }
    });
  }

  handleRequest(details) {
    // 检查是否是音乐相关的请求
    if (this.isMusicRequest(details.url)) {
      const requestData = {
        url: details.url,
        method: details.method,
        timestamp: details.timeStamp,
        tabId: details.tabId,
        status: 'pending',
        requestBody: details.requestBody
      };

      // 存储请求信息
      this.storeRequest(requestData);
    }
  }

  handleResponse(details) {
    // 检查是否是音乐相关的响应
    if (this.isMusicRequest(details.url)) {
      this.updateRequestStatus(details.tabId, details.url, 'completed');
    }
  }

  isMusicRequest(url) {
    // 定义音乐相关的URL模式
    const musicPatterns = [
      /\.mp3$/i,
      /\.wav$/i,
      /\.flac$/i,
      /\.aac$/i,
      /\.ogg$/i,
      /\.m4a$/i,
      /music/i,
      /audio/i,
      /song/i,
      /track/i,
      /album/i,
      /artist/i
    ];

    return musicPatterns.some(pattern => pattern.test(url));
  }

  storeRequest(requestData) {
    // 从存储中获取现有请求
    chrome.storage.local.get(['crawledRequests'], (result) => {
      const requests = result.crawledRequests || [];
      const requestId = this.generateRequestId(requestData);
      
      // 检查是否已存在
      const existingIndex = requests.findIndex(req => req.id === requestId);
      
      if (existingIndex >= 0) {
        // 更新现有请求
        requests[existingIndex] = { ...requests[existingIndex], ...requestData };
      } else {
        // 添加新请求
        requests.push({ ...requestData, id: requestId });
      }

      // 保存到存储
      chrome.storage.local.set({ crawledRequests: requests });
      
      // 通知popup更新
      this.notifyPopupUpdate();
    });
  }

  updateRequestStatus(tabId, url, status) {
    chrome.storage.local.get(['crawledRequests'], (result) => {
      const requests = result.crawledRequests || [];
      const request = requests.find(req => req.url === url && req.tabId === tabId);
      
      if (request) {
        request.status = status;
        chrome.storage.local.set({ crawledRequests: requests });
        this.notifyPopupUpdate();
      }
    });
  }

  generateRequestId(requestData) {
    return `${requestData.url}_${requestData.timestamp}`;
  }

  notifyPopupUpdate() {
    // 通知popup页面更新
    chrome.runtime.sendMessage({ type: 'REQUEST_UPDATE' });
  }

  injectContentScript(tabId) {
    // 注入内容脚本到页面
    chrome.scripting.executeScript({
      target: { tabId: tabId },
      files: ['content.js']
    }).catch(() => {
      // 忽略注入失败的错误（可能是特殊页面）
    });
  }
}

// 初始化请求拦截器
new RequestInterceptor();