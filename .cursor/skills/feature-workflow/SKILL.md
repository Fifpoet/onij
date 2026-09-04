---
name: feature-workflow
description: Orchestrates onij feature work from docs/feat/<slug>.md as planner → implementer → tester → reviewer → verifier. Use when the user says 做 feat, 做feat, 按工作流, 走工作流, /feature-workflow, or asks to implement a named feat.
---

# 功能工作流

本 Skill 只编排实现阶段。方案正文在 `docs/feat/<slug>.md`，格式见 `docs/feat/README.md`。

项目约束以 `onij-project.mdc` 为准；动手前 Read `onij-frontend.mdc` / `onij-backend.mdc` / `onij-features.mdc`（按改动面）。

## 触发

| 用户说法 | 做什么 |
|----------|--------|
| 「做 feat \<slug\>」「做feat \<slug\>」 | 读 `docs/feat/<slug>.md`，走下方实现 checklist |
| 「按工作流 / 走工作流」且已有对应 feat 文件 | 同上 |
| 方案讨论收口（就这样 / 定了 / 收口） | **不走实现**。按 `docs/feat/README.md` 写 md，并用其中「固定输出」回复用户 |
| 闲聊 / 八股 / 方案还在问 | 不写 feat、不开实现 |
| 「部署 / 上线」 | 走部署，不走本流程 |

没有 `docs/feat/<slug>.md` 就说「做 feat slug」：能从本对话提炼出目标/范围/非目标则先补写 md 再实现；否则停，继续讨论。

子 Agent **看不到**主对话。委派时必须带上：feat 文件路径、本轮任务原文、文件范围、约束。

## 角色

| 阶段 | 子 Agent | 说明 |
|------|----------|------|
| 拆任务 | `planner` | 只读 feat md + 仓库。产出文件级计划。任务含糊则停 |
| 实现 | `implementer` | 改代码，并**勾选** feat 里「本轮任务」、写「实现备注」 |
| 写测 | `tester` | 补与改动匹配的测试；不改 feat md |
| Review | `reviewer` | 只读。对照 feat 本轮任务与 diff。禁止自审 |
| 验收 | `verifier` | 只读。编译/单测/lint/（UI 则浏览器） |

## 实现 Checklist

```
- [ ] 0 已读 docs/feat/<slug>.md；本轮任务可验收。不够 → 停，提问
- [ ] 1 planner（对照 feat；不要另起一套需求）
- [ ] 2 主 Agent 如有必要：把 planner 的文件名单补进 feat「本轮任务」
- [ ] 3 已 Read 对应 mdc
- [ ] 4 implementer（改代码 + 勾选 feat）
- [ ] 5 tester（无测试床则注明跳过）
- [ ] 6 reviewer。有 Critical → 回 4，再 Review
- [ ] 7 verifier
- [ ] 8 主 Agent 收口 feat：补勾、更新 status、追加新的「相关未做」。不 commit / push
```

## 委派规则

- 实现与写测 **串行**，同一工作区；不要并行两个会改文件的子 Agent。
- 文案/样式一行级且不改行为：可跳过 planner 与 tester，仍要 reviewer + verifier；仍要勾 feat（若有对应项）。
- 涉及逻辑或 UI 行为：planner、reviewer、verifier 都不可省。
- 改表 / 写库 / SSH / 部署：工作流内 **不做**。
- 「相关未做」不是本轮范围，implementer 不准顺手做，也不准勾成完成。

## 验收命令（verifier 用）

- 改了 `src/**`：`pnpm exec vue-tsc --noEmit`；ReadLints 看改过的文件。
- 改了 `server/**`：相关包 `go test`，不要无故跑 handler 集成测。
- 改了 UI：浏览器走 feat「验证」里写的路径；麦克风/网易云/七牛标给人确认。

## 工作流收口（对用户）

```
feat <slug>：本轮 …
实现：改动文件 …
feat 任务：已勾 … / 未勾 …
相关未做：…
测试 / Review / 验收：…
未验证：…
```
