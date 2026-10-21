---
slug: presence-sync
title: 在线昵称与听歌同步
status: in_progress
---

# 在线昵称与听歌同步

## 目标

打开站点即有可改的随机中文昵称（顶栏可见）；能看到其他在线用户正在听的歌，点封面把对方的当前曲、进度和队列一次覆盖到自己这边并开始播放。

## 非目标

- 账号登录、跨设备同一身份
- 用户资料页、主页、私聊、好友关系
- 持续跟听（对方之后切歌/暂停/拖进度不再跟着走）
- 隐身 / 不广播正在听
- 同步 KTV、MV、云盘视频
- 新表 / 登录态；在线名单只活在服务端内存

## 大体方案

- **身份**：`sessionStorage` 存 `client_id`（每个标签页独立，刷新仍在）；`localStorage` 只记 `nickname`。首次进入用词库「形容词 + 的 + 名词」生成一次（如「安静的磁带」），改过后记住；允许重名。同一浏览器两个窗口必须是两个人，否则列表会把自己剔掉变成空的。
- **广播**：只要页面开着就心跳上报（id、昵称、在播、队列、进度秒、是否在播）。所有访客互相可见，无隐身。没在听也出现在列表，文案「未在播放」，封面点击无操作或 toast 即可。
- **通道**：Go 内存 map + TTL。`customizedRegister` 增加 `POST /presence/heartbeat`（上报并返回他人列表）。不改表。心跳约 10s，超时约 35s 视为离线。单次队列条数设上限（如 80），避免 payload 过大。
- **同步一次**：点对方**封面**→ 用对方 `nowPlaying` + `queue` **整表覆盖**自己的 `playQueue`（含清空己方原队列），按对方进度 seek 后播放。不要 `playNow` 默认的「当前曲插回队首」。点**歌手**→ `/artist`；点**歌名**→ `/album`。不能点进用户。
- **桌面顶栏**：左上在播不动。右侧（搜索旁）持久显示自己昵称（点击改名）+ 好友图标（打开在线列表）。不跟 `onij` 彩蛋绑定。
- **手机**：单独一版，不新路由。顶栏左侧仍是在播；右侧为**截断昵称**（约 5–7 字、可点改名）+ **好友图标** + 搜索（`onij` 菜单仍在最后）。列表用**全屏 overlay**（对齐现有搜索 overlay）：顶栏「在线」+ 关闭；列表行封面 | 歌手 / 歌名两行，未在播灰显。不要做资料卡。
- **`onij` 解锁**：改为持久（localStorage）。只控制 AI / 工具 / 语音 / 汉堡菜单，不影响昵称和好友列表。

### 拟改文件

- `src/util/presenceName.ts`：词库与首次昵称
- `src/store/presence.ts`（或等价 composable）：id/昵称持久、心跳、列表
- `src/api/presence.ts`：heartbeat
- `src/components/layout/`：桌面昵称+好友入口；手机昵称+好友 overlay（`DesktopHeader.vue` / `MobileHeader.vue` 最小接入，勿改搜索结果样式）
- `src/store/playQueue.ts` + `NeteaseAudioHost` / `neteasePlayerControl`：覆盖队列并在可播后 seek（不要拆第二套 audio）
- `src/composables/useNavToolsUnlock.ts`：`navToolsVisible` 持久
- `server/router.go` `customizedRegister` + `handler` / `logic`：内存 presence，无 DDL
- `App.vue`：仅在已有壳上挂心跳（若必须），禁止动全局布局/搜索页

### 与现有功能的关系

- 复用：`playQueue`、`ViewMusicListItem`、取链/播放宿主、顶栏左右分区、手机搜索 overlay 交互。
- 不要动：搜索结果列表样式、KTV、语音唤醒、七牛、歌词解析。

## 本轮任务

- [x] 本机匿名 id + 首次随机昵称 + 可改名并持久
- [x] `POST /presence/heartbeat` 内存在线表（全暴露，TTL 离线），返回他人列表
- [x] 页面开启即心跳；把自己的在播/队列/进度带上
- [x] 桌面顶栏右侧：昵称 + 好友列表（封面 / 歌手 / 歌名）
- [x] 手机顶栏截断昵称 + 好友图标；全屏 overlay 列表（单独排版）
- [x] 点封面：覆盖当前曲 + 队列 + 进度并播放；点歌手/歌名走现有详情路由
- [x] `onij` 解锁持久化（仍不控制好友入口）

## 相关未做

- 持续跟听（原因：本轮明确不做）
- 账号 / 跨设备同步身份
- 用户资料、隐身、私聊
- 同步 KTV / MV / 云盘视频
- 在线状态落库

## 验证

- 命令：`pnpm exec vue-tsc --noEmit`；`go test` 覆盖 presence logic（不要无故跑 handler 集成测）
- 浏览器：两个窗口（可隐私模式）打开**同一套环境**（本机都走 `18968`，不要一个本机一个线上），都能看到对方昵称、头像与在听的歌；点封面后本机当前曲、进度、队列与对方一致并开始播；点歌手/歌名进详情；改名刷新仍在；关掉一方约一分半内从列表消失
- 必须人确认：两台真实设备互相同步听歌是否顺（取链/seek）

## 风险

- `DesktopHeader` / `MobileHeader` / `App.vue` 必须最小 diff，禁止顺手改搜索页
- 心跳带整表队列，必须有条数上限
- seek 需等播放 URL 就绪，避免 `playNow` 把进度清零后播从头
- 工作流内不改表、不部署

## 实现备注

- 2026-10-13 后端：`POST /presence/heartbeat` 走 customizedRegister；进程内 map+mutex，TTL 35s，队列截 80，返回不含自己。handler 校验 client_id / 昵称非空。`wire ./inject` 已重生 `wire_gen.go`。`go build` 通过。
- 2026-10-13 前端：localStorage 持久 clientId（`randomUUID`）+ 随机「形容词的名词」昵称，可点改名。`App.vue` 仅 script 调 `startHeartbeat`（立即 + 10s，inflight 跳过）。桌面顶栏搜索前插昵称+好友面板；手机为截断昵称 | 好友图标 | 搜索 | 汉堡。点封面 `replaceFromPresence` 整表覆盖队列并记 `pendingSeekSec`（先写 seek 再改 `nowPlaying`，避免 watch 丢进度）；`NeteaseAudioHost` 等 loadedmetadata 后 seek 再播，不走 `playNow`。无在播封面 disabled（顶栏在 `n-message-provider` 外，未做 toast）。歌手 `/artist?ids=`、歌名 `/album?id=`。`navToolsVisible` 读写 localStorage，不控制好友入口。好友行额外显示对方昵称以便区分。
- 2026-10-13 review 后：标签页回到前台立刻心跳；改名立刻上报；同步后取链失败停在目标曲，不再 `advanceAfterEnded`。
- 2026-10-21：无痕页看不到原窗口，是原标签页进后台后 `setInterval` 被节流，35s TTL 把人清掉。心跳改到 Worker，打开列表立刻拉一次，TTL 90s；词库加大；头像由 client_id 哈希，不用另存。
- 2026-10-21：昵称不再截断；改名时头像留着；好友图标有人在线出红点人数；心跳 3s，切歌立刻上报。Playwright 双页已能看到对方昵称和红点。
- 2026-10-21：本地心跳改为 Vite 同源反代本机 Go（不再直连 8889）；Worker 改为模块脚本；播放中 `timeupdate` 也会续命，避免后台无痕被节流后从列表消失。
- 2026-10-21：线上接口本身是通的。同一浏览器两个窗口共用 localStorage `client_id`，服务端按 id 剔自己，两边都看到空列表。改为 `sessionStorage` 存 id（昵称仍 localStorage）；Worker 按经典脚本加载（构建产物是 IIFE，不能当 module）。

