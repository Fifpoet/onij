---
slug: cloud-video
title: 云盘视频本页播放与打点
status: in_progress
---

# 云盘视频本页播放与打点

## 目标

云盘（及同样展示 `FT_Video` 的中转文件）上传 MP4 后，在**当前页拉起**播放器：短片边下边播，演唱会级长片（最长约 3h）用单码率切片可拖进度；进度条可打密集关键点，悬停时局部放大以便点选。清晰度固定一档。

## 非目标

- 多码率 / ABR、自适应清晰度
- 非 MP4 源转码（AVI/MKV/MOV/HEVC 本轮不收）
- 弹幕、直播、独立 `/video/:id` 路由
- 悬停预览缩略图 / 雪碧图
- 视频文件经 Go 把 body 转一遍
- 改音乐 `playQueue` / `NeteaseAudioHost` 的播放核心（只做互斥暂停）

## 大体方案

- **存储**：仍走现有七牛直传 + `file` 表（`FT_Video` 已有）。大文件改为**分片断点续传**（不要整文件 FormData）。签名域保持 `http://cloud.onij.fun`，浏览器经 `https://onij.fun/cdn/...`。
- **可播形态**：片源已是 MP4。短片直出原对象（依赖 Range + `moov` 在前；必要时七牛 `avthumb` 做 faststart）。**时长 > 20 分钟或体积 > 512MB**（演唱会 3h 必走这条）生成**单码率 HLS**（m3u8 + ts，固定清晰度）；播放用原生 `<video>` + `hls.js`（Safari 可 native HLS）。
- **HLS 私有空间**：m3u8 里的切片必须带与现网一致的下载签名（后端改写清单，或七牛 `pm3u8`），切片仍走 `/cdn`，不走 Go body。
- **本页拉起**：不新开路由。全局挂一个视频宿主（类比 `NeteaseAudioHost`，但是 `<video>` + 遮罩/底栏进度）。云盘 `FileList`、中转站视频文件点击即打开。开视频时**暂停音乐**；关播放器不自动恢复音乐。
- **打点**：不要塞 `file.extra`（VARCHAR(1024)，已有 EXIF）。新表 `file_video_marker`（`file_id + time_ms + label`）。进度条画标记；悬停出现**局部放大尺**（放大光标附近一段时间窗），密集节点靠放大后点选，不靠缩略图。
- **nginx `/cdn`**：转发 `Range`；`proxy_buffering off`；读超时加长。否则长片 206 会被缓冲成整段或 502。

### 拟改文件

- `deploy/nginx/web-vue.conf`：`/cdn` Range + 关闭缓冲
- `src/util/qiniu.ts`：视频/大文件分片续传（本地 HTTP 避开 `crypto.subtle` 问题）
- `src/components/file/FileList.vue`、`src/views/TranPage.vue`：点击视频拉起播放器，不要当图片灯箱
- `src/components/player/` 新视频宿主 + 进度条（标记 + 悬停放大）；`App.vue` 最小 diff 挂宿主
- `src/store/` 视频会话状态（当前 file、播放中、勿并进 `playQueue`）
- `server/model/ddl/file_video_marker.sql`：新表（**只写 DDL，执行须用户确认**）
- `server/router.go` `customizedRegister` + logic：打点 CRUD；长片触发/查询 HLS 状态；签发或改写播放 URL
- `server/util/qiniu.go`：pfop（faststart / 单码率 m3u8）、私有 m3u8 处理
- 依赖：按需加 `hls.js`（不要新 UI 框架）

### 与现有功能的关系

- 复用：`/file/upload_token`、`/file/upload` 落库、`normalizeFileCdnUrl`、云盘文件夹与列表。
- 不要动：搜索页、音乐队列与歌词全屏、KTV、图片灯箱逻辑（视频走另一条打开路径）。
- 图片仍走现有 preview；视频禁止 `loading="lazy"` 式整文件预拉。

## 本轮任务

- [x] `/cdn` 支持视频 Range（206），`proxy_buffering off`，超时足够切片/拖动
- [x] 大 MP4 分片上传 + 现有 `file` 落库（`FT_Video`）；小文件可维持现上传
- [x] 本页拉起视频播放器（云盘 + 中转视频）；暂停音乐；不跳路由
- [x] 短片：签名 MP4 边下边播、可拖进度（faststart 不够则 pfop）
- [x] 长片（>20min 或 >512MB）：单码率 HLS，切片 URL 可签名经 `/cdn` 播放
- [x] `file_video_marker` DDL + 列表/批量保存/删除 API（改表等用户执行）
- [x] 进度条展示密集打点；悬停局部放大；点击标记 seek
- [x] 播放中可新增/改名/删除当前时间点的标记并写回后端

## 相关未做

- 多码率 ABR、非 MP4 转码（原因：范围外）
- 云盘上传大文件进度条（中转已接 onProgress；FileUploadModal 未接）
- 删除原片时清理 HLS m3u8/ts 产物
- 悬停缩略图 / 雪碧图（原因：用局部放大代替）
- 独立视频页、弹幕、直播
- HLS 转码完成前的「处理中」精细队列 UI 以外的运营工具
- 自动按歌声/镜头切点（本轮仅人工打点）

## 验证

- 命令：`pnpm exec vue-tsc --noEmit`；相关 `go test`（logic，不要无故跑 handler 集成测）
- 浏览器：`/cloud` 上传短 MP4 → 本页打开能未下完就播、能拖；再测长片（或 >512MB）HLS 可拖；进度条打多个点，悬停放大后能点中；中转站视频同样拉起；开视频时音乐应停
- 必须人确认：七牛上传、签名 `/cdn`、长片 HLS 是否出片（异步 pfop）

## 风险

- `App.vue` 只挂宿主，禁止改全局样式/搜索布局
- 七牛签名必须仍是桶绑定域，禁止签成 `onij.fun` 当对象路径
- `file.extra` 不够装打点，必须新表；**工作流内不执行 DDL**
- nginx 缓冲不关则长片 Range 失败；HLS 清单不改写则 ts 无 token 全 401
- 本地 Vite HTTP 下勿用依赖 `crypto.subtle` 的七牛 SDK 摘要
- 视频与音乐两套宿主并存，避免两个声音同时出

## 实现备注

- 2026-09-04：本轮任务已实现。`deploy/nginx/web-vue.conf` 的 `/cdn` 转发 Range、关缓冲、超时 3600s、gzip off。后端新增 `file_video_marker` DDL（未执行）、Dal/Logic、`/file/video/play` 与三个 marker 接口；短片签 MP4（moov 不在前则 pfop faststart），长片（时长>20min 或 DB size>512MB）异步 pfop 单码率 HLS 并改写 m3u8 切片为签名 URL。删除文件时顺带软删标记。前端 `hls.js`；≥8MB 七牛 mkblk/bput/mkfile；`VideoPlayerHost` + `VideoMarkerBar` 本页播放与打点；开视频 `pause` 停音乐，关播放器不恢复。`wire ./inject` 已生成，`go build` 可通过。打点表需用户在库上执行 `server/model/ddl/file_video_marker.sql` 后 marker API 才可用；线上 `/cdn` 需同步 nginx 配置。
- 2026-09-05：标记点改绿色；回车+失焦不再重复保存；去掉关闭按钮；全屏工具条贴底，鼠标下移唤出、3s 无悬停收起；中央 ‹ › 跳转相邻标记。

