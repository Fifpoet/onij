// Popup 脚本 - 处理界面逻辑

class PopupManager {
  constructor() {
    this.init();
  }

  init() {
    this.bindEvents();
    this.loadRequests();
    this.setupMessageListener();
  }

  bindEvents() {
    // 清空按钮
    document.getElementById('clearBtn').addEventListener('click', () => {
      this.clearRequests();
    });

    // 刷新按钮
    document.getElementById('refreshBtn').addEventListener('click', () => {
      this.loadRequests();
    });
  }

  setupMessageListener() {
    // 监听来自background的消息
    chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
      if (message.type === 'REQUEST_UPDATE') {
        this.loadRequests();
      }
    });
  }

  async loadRequests() {
    try {
      const result = await chrome.storage.local.get(['musicRequests']);
      const requests = result.musicRequests || [];
      
      this.renderRequests(requests);
      this.updateStatus(requests.length > 0 ? 'active' : 'idle');
      
    } catch (error) {
      console.error('加载请求失败:', error);
      this.updateStatus('error');
    }
  }

  renderRequests(requests) {
    const listContainer = document.getElementById('requestList');
    const countElement = document.getElementById('requestCount');
    
    // 更新计数
    countElement.textContent = `${requests.length} 条`;
    
    if (requests.length === 0) {
      listContainer.innerHTML = `
        <div class="empty-state">
          <p>暂无捕获的请求</p>
          <p class="hint">访问 https://zm.i9mr.com/song/detail 开头的页面来测试</p>
        </div>
      `;
      return;
    }

    // 渲染请求列表
    listContainer.innerHTML = requests.map((request, index) => {
      if (request.musicInfo) {
        // 有音乐信息的情况
        return `
          <div class="request-item music-item" data-id="${request.id}">
            <div class="request-item-header">
              <div class="request-number">${index + 1}</div>
              <div class="request-time">${this.formatTime(request.timestamp)}</div>
            </div>
            <div class="music-info">
              <div class="music-title">
                <span class="song-name">${request.musicInfo.songName}</span>
                <span class="artists"> - ${request.musicInfo.artists}</span>
              </div>
              <div class="music-album">${request.musicInfo.album}</div>
            </div>
          </div>
        `;
      } else {
        // 没有音乐信息的情况（旧格式）
        return `
          <div class="request-item" data-id="${request.id}">
            <div class="request-item-header">
              <div class="request-number">${index + 1}</div>
              <div class="request-time">${this.formatTime(request.timestamp)}</div>
            </div>
            <div class="request-url">
              <span class="request-method">${request.method}</span>
              ${request.url}
            </div>
          </div>
        `;
      }
    }).join('');

    // 添加点击事件
    listContainer.querySelectorAll('.request-item').forEach(item => {
      item.addEventListener('click', () => {
        const requestId = item.dataset.id;
        this.handleRequestClick(requestId);
      });
    });
  }

  handleRequestClick(requestId) {
    // 可以在这里添加点击请求项的处理逻辑
    console.log('点击了请求:', requestId);
  }

  async clearRequests() {
    try {
      await chrome.storage.local.set({ musicRequests: [] });
      this.loadRequests();
      console.log('请求列表已清空');
    } catch (error) {
      console.error('清空请求失败:', error);
    }
  }

  updateStatus(status) {
    const statusIndicator = document.getElementById('statusIndicator');
    const statusText = document.getElementById('statusText');
    
    switch (status) {
      case 'active':
        statusIndicator.style.background = '#4ade80';
        statusText.textContent = '监听中...';
        break;
      case 'idle':
        statusIndicator.style.background = '#f59e0b';
        statusText.textContent = '等待中...';
        break;
      case 'error':
        statusIndicator.style.background = '#ef4444';
        statusText.textContent = '错误';
        break;
    }
  }

  formatTime(timestamp) {
    const date = new Date(timestamp);
    const now = new Date();
    const diff = now - date;
    
    if (diff < 60000) { // 1分钟内
      return '刚刚';
    } else if (diff < 3600000) { // 1小时内
      return `${Math.floor(diff / 60000)}分钟前`;
    } else if (diff < 86400000) { // 1天内
      return `${Math.floor(diff / 3600000)}小时前`;
    } else {
      return date.toLocaleDateString();
    }
  }
}

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', () => {
  new PopupManager();
});