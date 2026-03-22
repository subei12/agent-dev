const statusMap: Record<string, string> = {
  active: "进行中",
  adopted: "已采纳",
  claimed: "已领取",
  completed: "已完成",
  discussion: "讨论中",
  done: "已完成",
  failed: "失败",
  handoff_pending: "等待交接",
  implementation: "开发中",
  in_progress: "进行中",
  open: "开放",
  pending: "待处理",
  proposed: "待采纳",
  review: "评审中",
  sealed: "已封存",
  todo: "待开始",
  uploaded: "已上传",
  waiting_admin: "等待管理员",
  working: "执行中"
};

const typeMap: Record<string, string> = {
  architecture: "架构文档",
  code: "代码任务",
  design: "设计任务",
  implementation_plan: "实施计划",
  review: "评审任务",
  review_report: "评审报告",
  runtime: "运行态",
  session: "会话",
  transcript: "转录",
  doc_review: "文档评审"
};

const accessModeMap: Record<string, string> = {
  redacted_transcript: "脱敏转录",
  summary: "摘要",
  export: "导出"
};

/**
 * formatStatusLabel 把后端状态码映射成适合界面显示的中文标签。
 */
export function formatStatusLabel(value: string | undefined) {
  if (!value) return "未知";
  return statusMap[value] ?? value;
}

/**
 * formatTypeLabel 把类型值映射成适合界面显示的中文标签。
 */
export function formatTypeLabel(value: string | undefined) {
  if (!value) return "未分类";
  return typeMap[value] ?? value;
}

/**
 * formatAccessModeLabel 把访问模式映射成中文标签。
 */
export function formatAccessModeLabel(value: string | undefined) {
  if (!value) return "未知";
  return accessModeMap[value] ?? value;
}
