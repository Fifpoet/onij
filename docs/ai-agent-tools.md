# onij AI Agent：Tool 清单与运行约定

> 状态：已接入 Agent loop + HITL 软暂停（合集写操作需确认）  
> 约定：语音助手走同一套 `/ai/chat/stream`（独立 `kind=voice` 会话，仅喂最近 3 轮 user）；共用本机环境变量 `dskey` / DeepSeek API。  
> 数据：`ai_session` + `ai_message` + `ai_confirm_pending`，**暂无 turn 表**。  
> 网易云代理：后端默认 `http://onij.fun:3000`（可用 `NETEASE_API_BASE` 覆盖）。  
> 步数上限：`AI_AGENT_MAX_STEPS`（默认 6）。

---

## 0. 关键代码位置（先看这里）

| 职责 | 路径 |
|------|------|
| Tool 清单文档（本文） | `docs/ai-agent-tools.md` |
| DeepSeek Chat / tools 协议 | `server/util/deepseek.go` |
| 网易云搜索/详情 HTTP | `server/util/netease.go` |
| Tool schema + 执行器 | `server/logic/ai_tools.go` |
| Agent 循环（LLM↔tool） | `server/logic/ai_logic.go` → `runAgentLoop` / `executeToolBatch` |
| HITL 确认 / 续跑 | `server/logic/ai_confirm.go` |
| Pending 表 DDL | `server/model/ddl/ai_confirm_pending.sql` |
| SSE 入口 | `server/handler/ai.go` → `AiChatStream` / `AiConfirmStream` |
| 路由 | `server/router.go` → `/ai/chat/stream`、`/ai/confirm/stream` |
| 前端 SSE 客户端 | `src/api/ai.ts` → `AiChatStream` / `AiConfirmStream` |
| 前端确认卡片 | `src/views/AiPage.vue` |
| 前端副作用执行 | `src/composables/useAiClientTools.ts` |
| 合集业务（tool 复用） | `server/logic/collection_logic.go` |
| 语音助手（唤醒词「小雪」→ 主流程） | `src/composables/useVoicePipeline.ts`；会话 `kind=1` 顶置 |

### SSE 事件

| event | 含义 |
|-------|------|
| `session` | 会话（含懒创建） |
| `user` | 用户消息已落库 |
| `assistant_tools` | 模型请求调用 tools |
| `tool_call` / `tool_result` | 单次 tool 调用与结果 |
| `client_action` | 需浏览器执行的前端 tool |
| `needs_confirm` | HITL：下发确认 UI（`confirm_id` + `ui`） |
| `paused` | 本轮软暂停结束（等用户确认） |
| `delta` | 最终回复文本（当前整段推送） |
| `done` / `error` | 正常结束 / 错误 |

### HITL 软暂停（要点）

1. 闸门在 `executeToolBatch`：**校验后、执行前**；当前强制确认：`collection_add_songs`、`collection_remove_song`。  
2. 落库 `ai_confirm_pending`（结构化 `ui_json` + `draft_args`），SSE `needs_confirm` 后 **结束本次请求**（不是 LLM 原生 pause）。  
3. 用户选择 → `POST /ai/confirm/stream` → 写 `role=tool`（`tool_call_id` 对齐）→ `runAgentLoop` 续跑。  
4. UI 形态：`yes_no` / `multi_select`（多首歌加入合集）等，由后端按 tool 组装，不喂给模型正文。  
5. 同批后续 tool 若遇暂停，写 `skipped` tool 结果以闭合协议。

---

## 1. Tool 清单

| name | 侧 | 实现状态 | 说明 | 主要参数 | 危险级 |
|------|----|----------|------|----------|--------|
| `search_music` | 后端 | ✅ | 网易云关键词搜索 | `query`；可选 `limit` | 低 |
| `get_song_detail` | 后端 | ✅ | 按三方 id 拉详情 | `song_ids: number[]` | 低 |
| `queue_add` | 前端* | ✅ | 加入播放队列 | `song_id` | 低 |
| `queue_play_now` | 前端* | ✅ | 立即播放 | `song_id` | 低 |
| `player_toggle` | 前端* | ✅ | 播放/暂停 | — | 低 |
| `player_next` | 前端* | ✅ | 下一首 | — | 低 |
| `open_page` | 前端* | ✅ | 路由跳转 | `path`, `query?` | 低 |
| `collection_list` | 后端 | ✅ | 合集列表/搜索 | `keyword?` | 低 |
| `collection_create` | 后端 | ✅ | 新建合集 | `name` | 低 |
| `collection_detail` | 后端 | ✅ | 合集详情+曲目 id | `id` | 低 |
| `collection_add_songs` | 后端 | ✅ + HITL | 合集加歌（确认后执行） | `id`, `songs[...]` | 中 |
| `collection_remove_song` | 后端 | ✅ + HITL | 合集删歌（确认后执行） | `id`, `song_id` | 中 |

\*「前端」：后端校验后 SSE `client_action`，由 `useAiClientTools` 执行；tool 结果写 `role=tool` 回灌模型。

### 后续可补充

- `queue_clear`（高危 HITL）
- 合集写操作 HITL 确认令牌
- `practice_log` / `uvr_separate` / `tran_*` / `cloud_search`
- 最终回复真正 token 级流式（当前 tool 轮非流式，终态整段 `delta`）

---

## 2. 测试用调用场景

| # | 用户说法（示例） | 期望 tool 轨迹（示意） |
|---|------------------|------------------------|
| 1 | 搜索林俊杰的江南 | `search_music` → 文字列出结果；（可选）`open_page` `/search?q=` |
| 2 | 把第一首加入队列并播放 | `search_music` → `queue_play_now` |
| 3 | 暂停 / 下一首 | `player_toggle` / `player_next` |
| 4 | 新建合集叫夜跑 | `collection_create` |
| 5 | 把江南加进夜跑合集 | `search_music` → `collection_list` → `collection_add_songs` → **HITL 确认** → 续跑 |
| 6 | 夜跑合集里有哪些歌 | `collection_list`/`detail` → 答复 |
| 7 | 打开合集页 | `collection_detail` → `open_page` `/collection?id=` |
| 8 | 搜周杰伦（空结果或错字） | `search_music`；失败则错误回灌再纠 query |

语音助手：先说「小雪」再给口令（如「小雪，播放江南」）。无唤醒词不调 LLM。会话在 AI 页列表顶置，标题固定「语音助手」。

---

## 3. 统一校验 · 错误回灌 · 步数上限

### 3.1 统一校验

1. **Schema / 参数校验（执行前）**  
   - 必填缺失、类型不对 → **不执行业务**，`tool` 回灌 `{ ok:false, error:"validation", data.detail }`。
2. **业务校验（executor 内）**  
   - id > 0、合集存在、`open_page` 路径白名单等。  
   - HITL：写操作暂未强制确认（后续补）。
3. **幂等**  
   - 合集已含该歌等由 collection 逻辑处理；队列 duplicate 由前端 enqueue 返回。

### 3.2 错误回灌

- 执行结果一律 `role=tool` 写入 `ai_message` 并 SSE `tool_result`。  
- 下一轮 LLM 可见失败原因；禁止静默吞掉。

### 3.3 失败重试

| 类型 | 策略 |
|------|------|
| 网络超时 / 5xx | 暂未自动重试网易请求；由模型决定是否再调 |
| Schema / 业务非法 | 不重试，回灌 |
| 搜索空结果 | 不重试同一 query；由模型决定是否换词 |

### 3.4 步数上限

- 单次用户发送：`AI_AGENT_MAX_STEPS`（默认 **6**，一次 LLM 调用算一步）。  
- 触及上限：最终 assistant「步骤过多已停止」。

### 3.5 Token

- DeepSeek `usage` 写入对应 assistant `ai_message`；session `token_total` 累加。  
- 不做 turn 表聚合。

---

## 4. 与代码的对应（当前进度）

| 能力 | 状态 |
|------|------|
| Tool 注册表 | ✅ `ai_tools.go` + 本文 §1 |
| `/ai` 页 + session/message + 流式骨架 | ✅ |
| Agent loop + 首批 tools | ✅ |
| HITL 软暂停（合集写） | ✅ `ai_confirm_pending` + `/ai/confirm/stream` |
| 前端副作用 | ✅ `useAiClientTools.ts` |
| 语音隔离 | ✅ 独立 `/ai/*` |
| HITL / 真正流式终态 | 后续 |
