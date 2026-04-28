// Admin 模块 — 管理员后台
(function () {
  const tabs = document.getElementById('admin-tabs');
  const contentEl = document.getElementById('admin-content');
  let currentTab = 'users';

  tabs.addEventListener('click', (e) => {
    if (e.target.classList.contains('tab')) {
      tabs.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
      e.target.classList.add('active');
      currentTab = e.target.dataset.tab;
      loadAdminContent();
    }
  });

  window.loadAdmin = function () {
    if (localStorage.getItem('is_admin') !== 'true') {
      toast('无管理员权限', 'error');
      router('chat');
      return;
    }
    loadAdminContent();
  };

  function loadAdminContent() {
    if (currentTab === 'users') loadUsers();
    else if (currentTab === 'stats') loadStats();
    else if (currentTab === 'quotas') loadQuotas();
  }

  async function loadUsers() {
    contentEl.innerHTML = `<div style="padding:40px;text-align:center;color:var(--text-tertiary)">${loadingEl()}</div>`;
    try {
      const data = await GET('/admin/users');
      const users = data.data || [];

      contentEl.innerHTML = `
        <div class="section-label">用户列表（${users.length} 人）</div>
        <div class="table-wrap">
          <table>
            <thead>
              <tr><th>ID</th><th>用户名</th><th>邮箱</th><th>管理员</th><th>注册时间</th><th>操作</th></tr>
            </thead>
            <tbody>
              ${users.map(u => `
                <tr>
                  <td style="color:var(--text-tertiary)">#${u.id}</td>
                  <td>${escapeHtml(u.username)}</td>
                  <td style="color:var(--text-secondary)">${escapeHtml(u.email || '—')}</td>
                  <td>${u.is_admin ? '<span style="color:var(--color-success);font-weight:500">是</span>' : '否'}</td>
                  <td style="color:var(--text-tertiary);font-size:12px">${fmtTime(u.created_at)}</td>
                  <td><button class="btn btn-ghost btn-sm" data-plan="${u.id}">设置套餐</button></td>
                </tr>`).join('')}
            </tbody>
          </table>
        </div>`;

      contentEl.querySelectorAll('[data-plan]').forEach(btn => {
        btn.addEventListener('click', () => openSetPlanModal(btn.dataset.plan));
      });
    } catch (err) {
      contentEl.innerHTML = `<div class="empty-state"><div class="empty-state-title">加载失败</div><div class="empty-state-desc">${err.message}</div></div>`;
    }
  }

  async function loadStats() {
    contentEl.innerHTML = `<div style="padding:40px;text-align:center;color:var(--text-tertiary)">${loadingEl()}</div>`;
    try {
      const data = await GET('/admin/stats');
      const s = data.data;
      const dist = s.plan_distribution || {};

      contentEl.innerHTML = `
        <div class="stats-grid">
          <div class="stat-card"><div class="stat-value">${s.total_users || 0}</div><div class="stat-label">总用户数</div></div>
          <div class="stat-card"><div class="stat-value">${s.today_active_users || 0}</div><div class="stat-label">今日活跃用户</div></div>
          <div class="stat-card"><div class="stat-value">${s.today_ai_calls || 0}</div><div class="stat-label">今日 AI 调用次数</div></div>
          <div class="stat-card"><div class="stat-value">${s.today_rag_calls || 0}</div><div class="stat-label">今日 RAG 调用次数</div></div>
        </div>
        <div class="card">
          <div class="card-header"><span class="card-title">套餐分布</span></div>
          <div class="card-body">
            <div style="display:flex;gap:16px;flex-wrap:wrap">
              <div style="display:flex;align-items:center;gap:6px"><span class="plan-badge free">免费版</span><span style="font-size:14px;font-weight:600">${dist['free'] || 0}</span> 人</div>
              <div style="display:flex;align-items:center;gap:6px"><span class="plan-badge basic">基础版</span><span style="font-size:14px;font-weight:600">${dist['basic'] || 0}</span> 人</div>
              <div style="display:flex;align-items:center;gap:6px"><span class="plan-badge pro">专业版</span><span style="font-size:14px;font-weight:600">${dist['pro'] || 0}</span> 人</div>
              <div style="display:flex;align-items:center;gap:6px"><span class="plan-badge unlimit">无限制</span><span style="font-size:14px;font-weight:600">${dist['unlimit'] || 0}</span> 人</div>
            </div>
          </div>
        </div>`;
    } catch (err) {
      contentEl.innerHTML = `<div class="empty-state"><div class="empty-state-title">加载失败</div><div class="empty-state-desc">${err.message}</div></div>`;
    }
  }

  async function loadQuotas() {
    contentEl.innerHTML = `<div style="padding:40px;text-align:center;color:var(--text-tertiary)">${loadingEl()}</div>`;
    try {
      const [quotaData, subData] = await Promise.all([
        GET('/admin/quotas/today'),
        GET('/admin/subscriptions')
      ]);

      const quotas = quotaData.data || [];
      const subs = subData.data || [];

      contentEl.innerHTML = `
        <div class="section-label">今日额度使用（${quotas.length} 人有记录）</div>
        <div class="table-wrap" style="margin-bottom:32px">
          <table>
            <thead><tr><th>用户ID</th><th>AI 已用</th><th>RAG 已用</th><th>日期</th></tr></thead>
            <tbody>
              ${quotas.length === 0 ? '<tr><td colspan="4" style="text-align:center;color:var(--text-tertiary);padding:24px">暂无数据</td></tr>' :
                quotas.map(q => `<tr><td>#${q.user_id}</td><td>${q.used_ai} 次</td><td>${q.used_rag} 次</td><td style="color:var(--text-tertiary)">${q.date}</td></tr>`).join('')}
            </tbody>
          </table>
        </div>
        <div class="section-label">用户订阅（${subs.length} 条）</div>
        <div class="table-wrap">
          <table>
            <thead><tr><th>用户ID</th><th>套餐</th><th>到期时间</th><th>更新时间</th></tr></thead>
            <tbody>
              ${subs.length === 0 ? '<tr><td colspan="4" style="text-align:center;color:var(--text-tertiary);padding:24px">暂无数据</td></tr>' :
                subs.map(s => `<tr><td>#${s.user_id}</td><td><span class="plan-badge ${s.plan}">${planLabel(s.plan)}</span></td><td style="color:var(--text-secondary)">${s.expired_at ? new Date(s.expired_at).toLocaleDateString('zh-CN') : '永久'}</td><td style="color:var(--text-tertiary);font-size:12px">${fmtTime(s.updated_at)}</td></tr>`).join('')}
            </tbody>
          </table>
        </div>`;
    } catch (err) {
      contentEl.innerHTML = `<div class="empty-state"><div class="empty-state-title">加载失败</div><div class="empty-state-desc">${err.message}</div></div>`;
    }
  }

  function openSetPlanModal(userId) {
    showModal('设置用户套餐', `
      <div style="display:flex;flex-direction:column;gap:12px">
        <div class="form-label">选择套餐</div>
        <div class="plan-grid">
          <div class="plan-card selected" data-plan="free"><div class="plan-card-name">免费版</div><div class="plan-card-desc">AI 5次/日，RAG 3次/日</div></div>
          <div class="plan-card" data-plan="basic"><div class="plan-card-name">基础版</div><div class="plan-card-desc">AI 30次/日，RAG 20次/日</div></div>
          <div class="plan-card" data-plan="pro"><div class="plan-card-name">专业版</div><div class="plan-card-desc">AI 100次/日，RAG 60次/日</div></div>
          <div class="plan-card" data-plan="unlimit"><div class="plan-card-name">无限制</div><div class="plan-card-desc">不限次数（谨慎分配）</div></div>
        </div>
      </div>`, `
      <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">取消</button>
      <button class="btn btn-primary" id="set-plan-btn">确认设置</button>`);

    let selectedPlan = 'free';
    document.querySelectorAll('.plan-card').forEach(card => {
      card.addEventListener('click', () => {
        document.querySelectorAll('.plan-card').forEach(c => c.classList.remove('selected'));
        card.classList.add('selected');
        selectedPlan = card.dataset.plan;
      });
    });

    document.getElementById('set-plan-btn').addEventListener('click', async () => {
      const btn = document.getElementById('set-plan-btn');
      btn.disabled = true; btn.textContent = '设置中…';
      try {
        await POST(`/admin/users/${userId}/plan`, { plan: selectedPlan });
        toast('套餐设置成功', 'success');
        document.querySelector('.modal-overlay').remove();
        // 如果是当前用户，更新 localStorage
        if (userId === localStorage.getItem('uid')) {
          localStorage.setItem('plan', selectedPlan);
          updateUserCard();
        }
        loadUsers();
      } catch (err) {
        toast(err.message, 'error');
        btn.disabled = false; btn.textContent = '确认设置';
      }
    });
  }

  function planLabel(p) {
    return { free:'免费版', basic:'基础版', pro:'专业版', unlimit:'无限制' }[p] || p;
  }

  function escapeHtml(str) {
    const div = document.createElement('div');
    div.textContent = str || '';
    return div.innerHTML;
  }
})();
