import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";

import { PageShell } from "../components/page-shell";
import { AgentConfig, getAgentConfigs, updateAgentConfig } from "../lib/api";

const projectId = "proj_1";

/**
 * AgentsPage 渲染 Agent 配置页。
 */
export function AgentsPage() {
  const queryClient = useQueryClient();
  const [selectedAgent, setSelectedAgent] = useState<AgentConfig | null>(null);

  const agentsQuery = useQuery({
    queryKey: ["agent-configs"],
    queryFn: () => getAgentConfigs(projectId)
  });

  const updateMutation = useMutation({
    mutationFn: (payload: Partial<AgentConfig> & { id: string }) => updateAgentConfig(projectId, payload.id, payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["agent-configs"] });
    }
  });

  const agents = agentsQuery.data ?? [];
  const current = selectedAgent ?? agents[0] ?? null;

  return (
    <PageShell
      eyebrow="Agent 配置"
      title="本地 CLI 执行器配置"
      description="当前阶段不启用沙盒，直接配置本地命令、参数和执行权限。"
      aside={
        <div className="hero-stat">
          <span>Agent 数量</span>
          <strong>{agents.length}</strong>
        </div>
      }
    >
      <section className="panel">
        <div className="panel-header">
          <p className="panel-kicker">Agent 列表</p>
        </div>
        <div className="stack">
          {agents.map((agent) => (
            <button key={agent.id} className="list-card" onClick={() => setSelectedAgent(agent)} type="button">
              <div className="list-card__meta">
                <span>{agent.executorType}</span>
                <span>{agent.enabled ? "启用" : "禁用"}</span>
              </div>
              <h3>{agent.name}</h3>
              <p>{agent.command}</p>
            </button>
          ))}
        </div>
      </section>
      <section className="panel panel--feature">
        <div className="panel-header">
          <p className="panel-kicker">配置详情</p>
        </div>
        {current ? (
          <AgentConfigForm
            agent={current}
            isSaving={updateMutation.isPending}
            onSave={(payload) => updateMutation.mutate({ id: current.id, ...payload })}
          />
        ) : (
          <div className="empty-card">
            <strong>当前没有 Agent 配置。</strong>
          </div>
        )}
      </section>
    </PageShell>
  );
}

function AgentConfigForm({
  agent,
  isSaving,
  onSave
}: {
  agent: AgentConfig;
  isSaving: boolean;
  onSave: (payload: Partial<AgentConfig>) => void;
}) {
  const [draft, setDraft] = useState(agent);

  return (
    <div className="stack">
      <label className="field">
        <span>Agent 名称</span>
        <input value={draft.name} onChange={(event) => setDraft({ ...draft, name: event.target.value })} />
      </label>
      <label className="field">
        <span>执行器类型</span>
        <select value={draft.executorType} onChange={(event) => setDraft({ ...draft, executorType: event.target.value })}>
          <option value="codex_cli">codex_cli</option>
          <option value="claude_code">claude_code</option>
          <option value="gemini_cli">gemini_cli</option>
        </select>
      </label>
      <label className="field">
        <span>本地命令</span>
        <input value={draft.command} onChange={(event) => setDraft({ ...draft, command: event.target.value })} />
      </label>
      <label className="field">
        <span>参数（逗号分隔）</span>
        <input
          value={draft.args.join(",")}
          onChange={(event) => setDraft({ ...draft, args: event.target.value.split(",").filter(Boolean) })}
        />
      </label>
      <button className="action-button" disabled={isSaving} onClick={() => onSave(draft)} type="button">
        保存 Agent 配置
      </button>
    </div>
  );
}
