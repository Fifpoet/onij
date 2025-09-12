// 音乐爬取助手 - 主界面逻辑
class MusicCrawler {
  constructor() {
    this.requests = new Map();
    this.urlFilter = '';
    this.init();
  }

  init() {
    this.bindEvents();
    this.loadStoredRequests();
    this.setupUrlFilter();
  }

  bindEvents() {
    // URL过滤
    document.getElementById('urlFilter').addEventListener('input', (e) => {
      this.urlFilter = e.target.value;
      this.filterRequests();
    });

    // 清除过滤
    document.getElementById('clearFilter').addEventListener('click', () => {
      document.getElementById('urlFilter').value = '';
      this.urlFilter = '';
      this.filterRequests();
    });

    // 模态框事件
    document.getElementById('closeModal').addEventListener('click', () => {
      this.closeModal();
    });

    document.getElementById('cancelEntry').addEventListener('click', () => {
      this.closeModal();
    });

    // 表单提交
    document.getElementById('entryForm').addEventListener('submit', (e) => {
      e.preventDefault();
      this.handleFormSubmit();
    });

    // 点击模态框背景关闭
    document.getElementById('entryModal').addEventListener('click', (e) => {
      if (e.target.id === 'entryModal') {
        this.closeModal();
      }
    });
  }

  setupUrlFilter() {
    // 从存储中恢复过滤条件
    chrome.storage.local.get(['urlFilter'], (result) => {
      if (result.urlFilter) {
        this.urlFilter = result.urlFilter;
        document.getElementById('urlFilter').value = this.urlFilter;
        this.filterRequests();
      }
    });
  }

  loadStoredRequests() {
    chrome.storage.local.get(['crawledRequests'], (result) => {
      if (result.crawledRequests) {
        this.requests = new Map(result.crawledRequests);
        this.renderRequests();
      }
    });
  }

  saveRequests() {
    const requestsArray = Array.from(this.requests.entries());
    chrome.storage.local.set({ crawledRequests: requestsArray });
  }

  addRequest(requestData) {
    const requestId = this.generateRequestId(requestData);
    this.requests.set(requestId, {
      ...requestData,
      id: requestId,
      timestamp: Date.now()
    });
    this.saveRequests();
    this.renderRequests();
  }

  generateRequestId(requestData) {
    return `${requestData.url}_${requestData.timestamp || Date.now()}`;
  }

  filterRequests() {
    // 保存过滤条件
    chrome.storage.local.set({ urlFilter: this.urlFilter });
    
    const filteredRequests = Array.from(this.requests.values()).filter(request => {
      if (!this.urlFilter) return true;
      return request.url.toLowerCase().includes(this.urlFilter.toLowerCase());
    });

    this.renderFilteredRequests(filteredRequests);
  }

  renderRequests() {
    const allRequests = Array.from(this.requests.values());
    this.renderFilteredRequests(allRequests);
  }

  renderFilteredRequests(requests) {
    const requestList = document.getElementById('requestList');
    const emptyState = document.getElementById('emptyState');

    if (requests.length === 0) {
      requestList.style.display = 'none';
      emptyState.style.display = 'block';
      return;
    }

    requestList.style.display = 'block';
    emptyState.style.display = 'none';

    requestList.innerHTML = requests.map(request => this.createRequestItem(request)).join('');
    
    // 绑定事件
    this.bindRequestItemEvents();
  }

  createRequestItem(request) {
    const statusClass = this.getStatusClass(request.status);
    const statusText = this.getStatusText(request.status);
    
    return `
      <div class="request-item" data-request-id="${request.id}">
        <div class="request-item-header">
          <div class="request-url" title="${request.url}">${this.truncateUrl(request.url)}</div>
          <div class="request-status ${statusClass}">${statusText}</div>
        </div>
        
        <div class="music-form">
          <div class="form-field">
            <label>歌曲名称</label>
            <div class="search-dropdown">
              <input type="text" 
                     class="song-search" 
                     placeholder="搜索或输入歌曲名" 
                     value="${request.songName || ''}"
                     data-field="songName" />
              <div class="dropdown-menu" data-field="songName"></div>
            </div>
          </div>
          
          <div class="form-field">
            <label>专辑名称</label>
            <div class="search-dropdown">
              <input type="text" 
                     class="album-search" 
                     placeholder="搜索或输入专辑名" 
                     value="${request.albumName || ''}"
                     data-field="albumName" />
              <div class="dropdown-menu" data-field="albumName"></div>
            </div>
          </div>
        </div>
        
        <div class="form-actions">
          <button class="btn btn-primary btn-sm entry-btn" data-request-id="${request.id}">
            录入
          </button>
        </div>
      </div>
    `;
  }

  bindRequestItemEvents() {
    // 搜索下拉框事件
    document.querySelectorAll('.search-dropdown input').forEach(input => {
      input.addEventListener('input', (e) => {
        this.handleSearchInput(e);
      });
      
      input.addEventListener('focus', (e) => {
        this.showDropdown(e.target);
      });
      
      input.addEventListener('blur', (e) => {
        // 延迟隐藏，让点击事件能够触发
        setTimeout(() => this.hideDropdown(e.target), 200);
      });
    });

    // 录入按钮事件
    document.querySelectorAll('.entry-btn').forEach(btn => {
      btn.addEventListener('click', (e) => {
        const requestId = e.target.getAttribute('data-request-id');
        this.openEntryModal(requestId);
      });
    });
  }

  handleSearchInput(e) {
    const input = e.target;
    const field = input.getAttribute('data-field');
    const query = input.value.trim();
    
    if (query.length < 2) {
      this.hideDropdown(input);
      return;
    }

    this.searchData(field, query).then(results => {
      this.showSearchResults(input, results);
    });
  }

  async searchData(field, query) {
    // 模拟搜索API调用
    // 这里后续会替换为实际的API调用
    return new Promise((resolve) => {
      setTimeout(() => {
        const mockData = this.getMockSearchData(field, query);
        resolve(mockData);
      }, 300);
    });
  }

  getMockSearchData(field, query) {
    // 模拟数据，后续替换为真实API
    const mockSongs = [
      { id: 1, name: '夜曲', artist: '周杰伦' },
      { id: 2, name: '青花瓷', artist: '周杰伦' },
      { id: 3, name: '稻香', artist: '周杰伦' }
    ];
    
    const mockAlbums = [
      { id: 1, name: '十一月的萧邦', artist: '周杰伦' },
      { id: 2, name: '魔杰座', artist: '周杰伦' },
      { id: 3, name: '叶惠美', artist: '周杰伦' }
    ];

    const data = field === 'songName' ? mockSongs : mockAlbums;
    return data.filter(item => 
      item.name.toLowerCase().includes(query.toLowerCase())
    );
  }

  showSearchResults(input, results) {
    const dropdown = input.parentElement.querySelector('.dropdown-menu');
    
    if (results.length === 0) {
      dropdown.innerHTML = '<div class="dropdown-item">未找到匹配项</div>';
    } else {
      dropdown.innerHTML = results.map(item => 
        `<div class="dropdown-item" data-id="${item.id}" data-name="${item.name}">
          ${item.name} ${item.artist ? `- ${item.artist}` : ''}
        </div>`
      ).join('');
      
      // 绑定选择事件
      dropdown.querySelectorAll('.dropdown-item').forEach(item => {
        item.addEventListener('click', (e) => {
          this.selectSearchResult(input, e.target);
        });
      });
    }
    
    this.showDropdown(input);
  }

  selectSearchResult(input, item) {
    const name = item.getAttribute('data-name');
    input.value = name;
    this.hideDropdown(input);
    
    // 更新请求数据
    const requestId = input.closest('.request-item').getAttribute('data-request-id');
    const field = input.getAttribute('data-field');
    this.updateRequestField(requestId, field, name);
  }

  showDropdown(input) {
    const dropdown = input.parentElement.querySelector('.dropdown-menu');
    dropdown.classList.add('show');
  }

  hideDropdown(input) {
    const dropdown = input.parentElement.querySelector('.dropdown-menu');
    dropdown.classList.remove('show');
  }

  updateRequestField(requestId, field, value) {
    const request = this.requests.get(requestId);
    if (request) {
      request[field] = value;
      this.saveRequests();
    }
  }

  openEntryModal(requestId) {
    const request = this.requests.get(requestId);
    if (!request) return;

    // 填充表单数据
    document.getElementById('songName').value = request.songName || '';
    document.getElementById('albumName').value = request.albumName || '';
    document.getElementById('artistName').value = request.artistName || '';
    document.getElementById('duration').value = request.duration || '';
    document.getElementById('fileUrl').value = request.url || '';
    document.getElementById('coverUrl').value = request.coverUrl || '';
    document.getElementById('lyrics').value = request.lyrics || '';

    // 存储当前编辑的请求ID
    document.getElementById('entryForm').setAttribute('data-request-id', requestId);

    // 显示模态框
    document.getElementById('entryModal').classList.add('show');
  }

  closeModal() {
    document.getElementById('entryModal').classList.remove('show');
    document.getElementById('entryForm').reset();
  }

  handleFormSubmit() {
    const form = document.getElementById('entryForm');
    const requestId = form.getAttribute('data-request-id');
    
    const formData = {
      songName: document.getElementById('songName').value,
      albumName: document.getElementById('albumName').value,
      artistName: document.getElementById('artistName').value,
      duration: parseInt(document.getElementById('duration').value) || 0,
      fileUrl: document.getElementById('fileUrl').value,
      coverUrl: document.getElementById('coverUrl').value,
      lyrics: document.getElementById('lyrics').value
    };

    // 验证必填字段
    if (!formData.songName || !formData.albumName) {
      alert('歌曲名称和专辑名称为必填项');
      return;
    }

    // 更新请求数据
    const request = this.requests.get(requestId);
    if (request) {
      Object.assign(request, formData);
      request.status = 'completed';
      this.saveRequests();
      this.renderRequests();
    }

    this.closeModal();
    
    // 这里可以添加API调用逻辑
    this.submitToAPI(formData);
  }

  async submitToAPI(formData) {
    try {
      // 模拟API调用
      console.log('提交数据到API:', formData);
      
      // 这里后续会替换为实际的API调用
      // const response = await fetch('/api/music', {
      //   method: 'POST',
      //   headers: { 'Content-Type': 'application/json' },
      //   body: JSON.stringify(formData)
      // });
      
      alert('音乐信息已保存！');
    } catch (error) {
      console.error('保存失败:', error);
      alert('保存失败，请重试');
    }
  }

  getStatusClass(status) {
    switch (status) {
      case 'completed': return 'status-success';
      case 'error': return 'status-error';
      default: return 'status-pending';
    }
  }

  getStatusText(status) {
    switch (status) {
      case 'completed': return '已完成';
      case 'error': return '错误';
      default: return '待处理';
    }
  }

  truncateUrl(url) {
    if (url.length <= 50) return url;
    return url.substring(0, 47) + '...';
  }
}

// 初始化应用
document.addEventListener('DOMContentLoaded', () => {
  new MusicCrawler();
});