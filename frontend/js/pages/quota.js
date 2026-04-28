// Quota 模块 — 额度查询页面
(function () {
  window.loadQuota = loadQuota;

  async function loadQuota() {
    const el = document.getElementById('quota-content');
    if (!el) return;
    el.innerHTML = `<div style="padding:40px;text-align:center;color:var(--text-tertiary)">${loadingEl()}</div>`;

    try {
      const data = await GET('/quota/me');
      const q = data.data;
      const planClass = q.plan || 'free';
      const planLabel = { free:'免费版', basic:'基础版', pro:'专业版', unlimit:'无限制' }[planClass] || '免费版';

      const aiBar = q.daily_ai_quota < 0
        ? `<div style="color:var(--color-success);font-size:14px">无限次</div>`
        : `<div class="quota-bar"><div class="quota-bar-fill ${q.used_ai/q.daily_ai_quota>=0.9?'danger':q.used_ai/q.daily_ai_quota>=0.6?'warn':''}" style="width:${Math.min((q.used_ai/q.daily_ai_quota)*100,100)}%"></div></div><div class="quota-sub">${q.used_ai} / ${q.daily_ai_quota} 次</div>`;

      const ragBar = q.daily_rag_quota < 0
        ? `<div style="color:var(--color-success);font-size:14px">无限次</div>`
        : `<div class="quota-bar"><div class="quota-bar-fill ${q.used_rag/q.daily_rag_quota>=0.9?'danger':q.used_rag/q.daily_rag_quota>=0.6?'warn':''}" style="width:${Math.min((q.used_rag/q.daily_rag_quota)*100,100)}%"></div></div><div class="quota-sub">${q.used_rag} / ${q.daily_rag_quota} 次</div>`;

      el.innerHTML = `
        <div class="quota-grid">
          <div class="quota-card">
            <div class="quota-label">当前套餐</div>
            <div style="display:flex;align-items:center;gap:8px;margin-top:4px">
              <span class="plan-badge ${planClass}">${planLabel}</span>
            </div>
            ${q.expired_at ? `<div class="quota-sub" style="margin-top:8px">到期：${new Date(q.expired_at).toLocaleDateString('zh-CN')}</div>` : ''}
          </div>
          <div class="quota-card">
            <div class="quota-label">AI 今日剩余</div>
            <div class="quota-value">${q.remain_ai < 0 ? '∞' : q.remain_ai}<span class="quota-suffix">次</span></div>
            ${aiBar}
          </div>
          <div class="quota-card">
            <div class="quota-label">RAG 今日剩余</div>
            <div class="quota-value">${q.remain_rag < 0 ? '∞' : q.remain_rag}<span class="quota-suffix">次</span></div>
            ${ragBar}
          </div>
          <div class="quota-card">
            <div class="quota-label">额度重置</div>
            <div class="quota-value" style="font-size:18px;margin-top:4px">每日 00:00</div>
            <div class="quota-sub">北京时间自动重置</div>
          </div>
        </div>

        <div class="card">
          <div class="card-header"><span class="card-title">套餐详情</span></div>
          <div class="card-body">
            <table>
              <thead>
                <tr><th>套餐</th><th>每日 AI 次数</th><th>每日 RAG 次数</th></tr>
              </thead>
              <tbody>
                <tr><td><span class="plan-badge free">免费版</span></td><td>5 次</td><td>3 次</td></tr>
                <tr><td><span class="plan-badge basic">基础版</span></td><td>30 次</td><td>20 次</td></tr>
                <tr><td><span class="plan-badge pro">专业版</span></td><td>100 次</td><td>60 次</td></tr>
                <tr><td><span class="plan-badge unlimit">无限制</span></td><td>无限</td><td>无限</td></tr>
              </tbody>
            </table>
          </div>
          <div class="card-footer">
            <span style="font-size:12.5px;color:var(--text-secondary)">如需升级套餐，请联系管理员</span>
          </div>
        </div>`;
    } catch (err) {
      el.innerHTML = `<div class="empty-state"><div class="empty-state-title">加载失败</div><div class="empty-state-desc">${err.message}</div></div>`;
    }
  }
})();
