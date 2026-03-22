import { FormEvent, useState } from "react";

import { AgentConfig } from "../../lib/api";

type TaskComposerProps = {
  agents: AgentConfig[];
  onCreateTask: (payload: { title: string; type: string; assignedAgentId?: string }) => void;
  isSubmitting: boolean;
};

/**
 * TaskComposer 提供新增任务的快捷表单。
 */
export function TaskComposer({ agents, onCreateTask, isSubmitting }: TaskComposerProps) {
  const [title, setTitle] = useState("");
  const [type, setType] = useState("code");
  const [assignedAgentId, setAssignedAgentId] = useState("");

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!title.trim()) return;
    onCreateTask({
      title: title.trim(),
      type,
      assignedAgentId: assignedAgentId || undefined
    });
    setTitle("");
  };

  return (
    <section className="panel panel--feature task-composer-card">
      <div className="panel-header panel-header--compact">
        <p className="panel-kicker">新增任务</p>
        <span className="badge">快速创建</span>
      </div>
      <h2>把新任务直接放进当前推进节奏里。</h2>
      <p className="panel-copy">创建后会直接回到任务看板，方便立刻分派 Agent、领取或进入评审流转。</p>
      <form className="inline-form inline-form--wide task-composer-form" onSubmit={handleSubmit}>
        <label className="field">
          <span>任务标题</span>
          <input placeholder="例如：补充管理员复核流程" value={title} onChange={(event) => setTitle(event.target.value)} />
        </label>
        <label className="field field--compact">
          <span>任务类型</span>
          <select value={type} onChange={(event) => setType(event.target.value)}>
            <option value="analysis">分析</option>
            <option value="design">设计</option>
            <option value="code">开发</option>
            <option value="test">测试</option>
            <option value="review">评审</option>
          </select>
        </label>
        <label className="field field--compact">
          <span>负责 Agent</span>
          <select value={assignedAgentId} onChange={(event) => setAssignedAgentId(event.target.value)}>
            <option value="">未分配</option>
            {agents.map((agent) => (
              <option key={agent.id} value={agent.id}>
                {agent.name}
              </option>
            ))}
          </select>
        </label>
        <button className="action-button" disabled={isSubmitting} type="submit">
          新增任务
        </button>
      </form>
      <div className="task-composer__tips">
        <span className="badge">主操作常驻首屏</span>
        <span className="badge">支持立即指定 Agent</span>
      </div>
    </section>
  );
}
