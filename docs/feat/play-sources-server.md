---
slug: play-sources-server
title: 播放源服务端配置与 GD 签名同步
status: done
---

# 播放源服务端配置与 GD 签名同步

## 目标

全站共用一套播放源（存在 Go 服务端，监控页改完所有人立刻用）。GD 播放器签名从他们站点自动拉，crc32 轮换不用再改代码、重新发前端。Meting 去掉。

## 非目标

- 登录鉴权 / 谁能改源（本站公开，监控页谁都能改；本轮不加锁）
- 新表 / migration（JSON 文件即可）
- 自建 lx-music、UnblockNetease、liveling 内建
- 回落网易 `/song/url`
- 把 crc32 源码做成文本框让人粘贴
- 改搜索 / 歌词 / 网易代理本身

## 大体方案

- **列表落服务端**：`customizedRegister` 增加 `POST /play-source/list`、`POST /play-source/save`。JSON 文件存在进程工作目录 `play_sources.json`（生产即 `release/bin/play_sources.json`）。无文件时用内建默认。不改表。
- **内建默认（顺序即取链顺序）**：
  1. GD 公开 API（`kind=gdapi`，`https://music-api.gdstudio.xyz`）——无 crc32，`types=url` 直接 JSON。
  2. GD 播放器（`kind=gdstudio`，生产同源 `/gdstudio`）——自动同步签名；网易空再酷我搜+取链。
- **删 Meting**：去掉内建、添加表单选项、`requestMetingPlayUrl`、qijieya 默认。旧浏览器 `localStorage.playSources` 不再读取。
- **自动同步签名（仅 `gdstudio`）**：运行时从该源 origin 拉 `js/player.js`（读 `version`）和 `js/crc32.min.js`，现有沙箱跑出 md5。内存缓存；`401 Invalid request` 清缓存再拉一次。监控卡片「同步签名」。打包 vendor 只作拉失败时的离线兜底。
- **前端**：`playSources` 不再 persist；启动/进监控拉 list；增删改排序/开关都 `save` 整表。取链读服务端那份，不要各人一份。

### 拟改文件

- `server/router.go`：注册 `/play-source/list`、`/play-source/save`
- `server/handler/` + `server/logic/` + `server/biz/prm/` + `inject/wire.go`：读写 JSON、缺省合并内建
- `src/store/playSources.ts`：去 persist；kind 改为 `gdstudio | gdapi`；从 API 灌入
- `src/api/`：play-source 客户端；`gdstudioSign.ts` 运行时拉脚本；`songUrl.ts` / `playSourceHealth.ts` 去掉 Meting，gdapi 不签名
- `src/views/MonitorPage.vue`：去 Meting；GD 播放器卡「同步签名」
- `.cursor/rules/onij-features.mdc`：更新取链约定

### 与现有功能的关系

- 复用：监控页音乐源卡片、GD 取链（url → autosource → kuwo）、`/gdstudio` 反代、`NeteaseAudioHost` / `playQueue`。
- 不要动：搜索页、`App.vue` 布局、云盘、KTV、presence。

## 本轮任务

- [x] `POST /play-source/list` / `save`，JSON 文件，无文件用内建两源；不改表
- [x] 前端去 pinia persist 与 Meting；启动拉服务端列表，监控改完即 save；`App.vue` 只加 script 加载；Vite 反代加上 `play-source`
- [x] `gdstudio` 运行时拉 player.js + crc32.min.js；401 再拉；监控「同步签名」；vendor 仅兜底
- [x] `gdapi`：无签名 `types=url`；内建默认排在播放器源前面
- [x] 探测：gdapi / Meting 不再用空曲 `185811`；gdapi 用能吐链的网易 id（如 1357375695）；gdstudio 仍走真实播放地址（空则酷我）
- [x] 更新 `onij-features.mdc` 取链描述

## 相关未做

- 监控页改源鉴权（原因：本轮明确公开可改）
- liveling / 其它 Meting 公益源内建（原因：Meting 类不稳定，已删）
- 自建 lx-music / UNM
- ~~crc32 同步结果写回服务端文件~~（已做：监控页「同步签名」写入 `gdstudio_sign.json`）
- 内建源相对顺序可改（原因：前后端 merge 钉死 gdapi → gdstudio，监控上移刷新会打回）
- `loadFromServer` 失败保留当前列表（原因：现在会重置成内建两项，再保存会盖掉自定义源）
- ~~「同步签名」失败要提示，不要静默 vendor~~（已做：失败写在卡片上，不退回 vendor）
- `play_sources.json` 进 `.gitignore`（原因：本机 `go run` 会在工作目录生成）

## 验证

- 命令：`pnpm exec vue-tsc --noEmit`；`go test` 覆盖 play-source logic（缺省合并、非法 kind 拒绝、保存后再读）
- 浏览器：`/monitor` 音乐源只剩 GD 公开 API + GD 播放器；改开关/顺序刷新后仍在（换浏览器/无痕也应一样）；点「同步签名」后 GD 播放器 ping 正常；搜一首歌能播
- 必须人确认：线上换浏览器后源列表是否同一套

## 风险

- 监控页公开可写，任何人能改全站取链（本站可接受；不要顺手加登录）
- `App.vue` / 搜索页最小 diff
- `play_sources.json` 在 `release/bin`，部署脚本不要 `--delete` 掉它
- 工作流内不部署、不改表
- 2026-11-24：监控「同步签名」把 player.js 版本 + crc32 写到服务端 `gdstudio_sign.json`；启动拉这份给全站用。按钮失败不再静默 vendor。沙箱同时认 `this.md5` / `window.md5`，避免浏览器 UMD 挂到真实 window 后看起来像同步过仍 401。

## 实现备注

- 2026-11-16：Go 增加 `POST /play-source/list` / `save`，工作目录 `play_sources.json`（缺/坏/空则内建 gdapi + gdstudio；非法 kind / 空 id / 重复 id 返回 400 且不写文件）。前端去掉 pinia persist 与 Meting；启动与监控页拉服务端列表，增删改排序/开关立刻 save，源地址输入 debounce 400ms 或失焦再写。`gdstudio` 运行时拉 `player.js` + `crc32.min.js`，401 清缓存再拉一次，监控卡片「同步签名」，vendor 仅兜底。`gdapi` 无签名 `types=url`，探测改用 1357375695。改动文件：`server/biz/prm/play_source.go`、`server/logic/play_source_logic.go`、`server/logic/play_source_logic_test.go`、`server/handler/play_source.go`、`server/router.go`、`server/logic/0_init.go`、`server/inject/wire.go`、`server/inject/wire_gen.go`、`src/api/playSource.ts`、`src/store/playSources.ts`、`src/api/netease/gdstudioSign.ts`、`src/api/netease/songUrl.ts`、`src/api/netease/playSourceHealth.ts`、`src/views/MonitorPage.vue`、`src/App.vue`（仅 script 加载）、`vite.config.ts`、`.cursor/rules/onij-features.mdc`。未做「相关未做」。
- 2026-11-16 验收：`go test ./logic/ -run PlaySource` 通过；本轮文件 ReadLints 无诊断。本机 `/monitor` 音乐源仅 GD 公开 API + GD 播放器，两源探测拿到网易 mp3；停用刷新后仍停用（服务端 JSON）；「同步签名」后再探测正常；搜索「海阔天空」双击可播。`vue-tsc` 全仓仍有既有报错（专辑/标签等无关 WIP），本轮文件未出现。未部署。线上换浏览器是否同一套源须人确认。
