# 前端苹果式工作台改造实施文档

## 1. 目标

在不破坏现有功能链路的前提下，把前端整体体验从“后台管理模板”调整为“简约、克制、产品化的协同工作台”，重点提升：

1. 全局导航层次
2. 首页首屏观感与重点路径
3. Mission 页的任务推进节奏
4. 运行日志与 transcript 的查看体验
5. 页面动效、留白、选中态与点击反馈

## 2. 设计方向

本轮采用“苹果官网式的简洁产品页语言”，但不做宣传页式的大段滚动叙事，核心原则如下：

1. 使用更轻的顶部悬浮导航，替代左侧后台式侧边栏
2. 使用大标题、低饱和浅色背景、柔和玻璃质感与克制阴影
3. 首屏明确主任务，不平均铺开所有对象卡片
4. 任务、Agent、文档与归档维持清晰主次关系
5. 默认先展示结构化事件，需要时再展开 transcript

## 3. 改动范围

本轮主要修改以下区域：

1. `frontend/src/app/layout.tsx`
2. `frontend/src/components/page-shell.tsx`
3. `frontend/src/pages/home.tsx`
4. `frontend/src/pages/mission.tsx`
5. `frontend/src/features/missions/*`
6. `frontend/src/features/tasks/*`
7. `frontend/src/features/runtime/*`
8. `frontend/src/styles.css`
9. `frontend/e2e/*.spec.ts`

## 4. 实施步骤

1. 先补首页与工作台新骨架的失败测试
2. 重构顶栏、页面头部和首页概览布局
3. 重构 Mission 页任务区、详情区与运行区的视觉关系
4. 统一按钮、卡片、输入框、选中态与微动效
5. 运行 `typecheck` 与 `playwright` 回归测试

## 5. 验证标准

完成后至少满足以下标准：

1. 首页打开后不再呈现后台侧栏布局
2. 顶部导航、首页首屏和 Mission 页首屏都有明确主次
3. `新增任务`、`选择任务`、`切换 Agent`、`展开 transcript` 交互仍然可用
4. `corepack pnpm typecheck` 通过
5. `corepack pnpm test:e2e` 通过
