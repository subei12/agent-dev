import { FormEvent, useState } from "react";

type MissionRequestComposerProps = {
  isSubmitting: boolean;
  onCreateMission: (payload: { title: string; description: string }) => void;
};

/**
 * MissionRequestComposer 负责让用户提交一个完整需求，而不是手工拆设计/开发/测试任务。
 */
export function MissionRequestComposer({ isSubmitting, onCreateMission }: MissionRequestComposerProps) {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");

  /**
   * handleSubmit 提交单个需求并交由平台进入后续多 Agent 协作流程。
   */
  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!title.trim()) {
      return;
    }

    onCreateMission({
      title: title.trim(),
      description: description.trim()
    });
    setTitle("");
    setDescription("");
  };

  return (
    <section className="panel dashboard-request-card">
      <div className="panel-header">
        <p className="panel-kicker">新建需求</p>
        <span className="badge">用户主入口</span>
      </div>
      <h2>提交一个目标，剩下的交给平台里的 Agent 团队。</h2>
      <p className="panel-copy">
        你只需要描述要做什么。管理员 Agent 会组织讨论、收敛方案、拆出设计/开发/测试等内部任务，并在完成后判断是否结项。
      </p>
      <form className="stack" onSubmit={handleSubmit}>
        <label className="field">
          <span>需求标题</span>
          <input
            placeholder="例如：实现多 Agent 自动开发流程"
            value={title}
            onChange={(event) => setTitle(event.target.value)}
          />
        </label>
        <label className="field">
          <span>需求说明</span>
          <textarea
            className="field-textarea"
            placeholder="补充背景、目标、预期结果和限制条件。"
            value={description}
            onChange={(event) => setDescription(event.target.value)}
          />
        </label>
        <div className="action-row">
          <button className="action-button" disabled={isSubmitting} type="submit">
            提交需求
          </button>
        </div>
      </form>
    </section>
  );
}
