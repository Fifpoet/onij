# 音乐信息提取器 Chrome 插件

一个用于监听和提取音乐网站请求信息的Chrome浏览器插件。

## 功能特性

- 🎵 监听 `https://zm.i9mr.com/song/detail` 开头的所有请求
- 📋 在弹窗中展示捕获的请求列表
- 🔢 显示请求序号、时间戳和URL
- 🗑️ 支持清空请求列表
- 🔄 实时更新请求状态

## 安装方法

1. 打开Chrome浏览器
2. 访问 `chrome://extensions/`
3. 开启"开发者模式"
4. 点击"加载已解压的扩展程序"
5. 选择包含 `manifest.json` 的文件夹
6. 插件安装完成

## 使用方法

1. 安装插件后，点击浏览器工具栏中的插件图标
2. 访问任何以 `https://zm.i9mr.com/song/detail` 开头的页面
3. 插件会自动捕获这些请求并显示在列表中
4. 在弹窗中查看捕获的请求详情

## 测试方法

1. 打开 `test.html` 文件进行功能测试
2. 点击测试链接或使用动态测试按钮
3. 查看插件弹窗中的请求列表

## 文件结构

```
chrome/
├── manifest.json      # 插件配置文件
├── background.js      # 后台脚本（监听请求）
├── popup.html         # 弹窗界面
├── popup.js          # 弹窗逻辑
├── styles.css        # 样式文件
├── test.html         # 测试页面
└── README.md         # 说明文档
```

## 技术实现

- **Manifest V3**: 使用最新的Chrome扩展API
- **WebRequest API**: 监听网络请求
- **Chrome Storage**: 存储请求数据
- **Service Worker**: 后台脚本运行环境

## 权限说明

- `webRequest`: 监听网络请求
- `storage`: 存储请求数据
- `activeTab`: 访问当前标签页
- `https://zm.i9mr.com/*`: 访问目标网站

## 开发说明

插件使用纯JavaScript开发，无需额外依赖。主要功能包括：

1. **请求监听**: 使用 `chrome.webRequest.onBeforeRequest` 监听特定URL模式的请求
2. **数据存储**: 使用 `chrome.storage.local` 存储请求信息
3. **界面更新**: 通过消息传递机制实现实时更新
4. **错误处理**: 完善的错误处理和用户反馈

## 注意事项

- 插件仅在目标网站有请求时才会捕获数据
- 请求数据存储在本地，不会上传到服务器
- 建议定期清空请求列表以节省存储空间