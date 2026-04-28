// History 模块 — RAG 历史记录
(function () {
  const listEl = document.getElementById('history-list');

  window.loadHistory = loadHistory;

  async function loadHistory() {
    if (!listEl) return;
    listEl.innerHTML = `<div style="padding:24px;text-align:center;color:var(--text-tertiary)">${loadingEl()}</div>`;

    try {
      const data = await GET('/ai/rag-history');
      const histories = data.data || [];
      renderHistories(histories);
    } catch (err) {
      listEl.innerHTML = `<div class="empty-state"><div class="empty-state-title">加载失败</div><div class="empty-state-desc">${err.message}</div></div>`;
    }
  }

  function renderHistories(histories) {
    if (histories.length === 0) {
      listEl.innerHTML = `
        <div class="empty-state">
          <div class="empty-state-icon">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <circle cx="12" cy="12" r="10"/>
              <polyline points="12 6 12 12 16 14"/>
            </svg>
          </div>
          <div class="empty-state-title">暂无问答记录</div>
          <div class="empty-state-desc">在 RAG 问答中产生的对话会自动保存在这里</div>
        </div>`;
      return;
    }

    listEl.innerHTML = histories.map(h => `
      <div class="history-item" data-id="${h.id}">
        <div class="history-q">${escapeHtml(h.question)}</div>
        <div class="history-a">${escapeHtml(h.answer || '')}</div>
        <div class="history-meta">
          ${fmtTime(h.created_at)}
          <span style="margin-left:12px;cursor:pointer;color:var(--color-primary)" data-view="${h.id}">查看详情</span>
          <span style="margin-left:12px;cursor:pointer;color:var(--color-error)" data-del="${h.id}">删除</span>
        </div>
      </div>`).join('');

    listEl.querySelectorAll('[data-view]').forEach(btn => {
      btn.addEventListener('click', (e) => { e.stopPropagation(); viewHistory(btn.dataset.view); });
    });
    listEl.querySelectorAll('[data-del]').forEach(btn => {
      btn.addEventListener('click', (e) => { e.stopPropagation(); deleteHistory(btn.dataset.del); });
    });
  }

  async function viewHistory(id) {
    try {
      const data = await GET(`/ai/rag-history/${id}`);
      const h = data.data;
      showModal('问答详情', `
        <div style="display:flex;flex-direction:column;gap:16px;max-height:500px;overflow-y:auto">
          <div>
            <div class="form-label" style="margin-bottom:6px">问题</div>
            <div style="font-size:13.5px;line-height:1.7;padding:12px;background:var(--bg-base);border-radius:var(--radius-md)">${escapeHtml(h.question)}</div>
          </div>
          <div>
            <div class="form-label" style="margin-bottom:6px">回答</div>
            <div style="font-size:13.5px;line-height:1.7;white-space:pre-wrap">${escapeHtml(h.answer || '')}</div>
          </div>
          ${h.references && h.references.length > 0 ? `
          <div>
            <div class="form-label" style="margin-bottom:6px">参考资料</div>
            ${h.references.map((r,i) => `<div style="font-size:12.5px;line-height:1.6;color:var(--text-secondary);padding:8px 12px;background:var(--bg-base);border-radius:var(--radius-md);margin-bottom:6px">【${i+1}】${escapeHtml(r)}</div>`).join('')}
          </div>` : ''}
        </div>`, `<button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">关闭</button>`);
    } catch (err) {
      toast(err.message, 'error');
    }
  }

  async function deleteHistory(id) {
    if (!confirm('确定删除这条记录？')) return;
    try {
      await DEL(`/ai/rag-history/${id}`);
      toast('已删除', 'success');
      loadHistory();
    } catch (err) {
      toast(err.message, 'error');
    }
  }

  window.clearAllHistory = async function () {
    if (!confirm('确定清空所有问答记录？此操作不可恢复。')) return;
    try {
      await DEL('/ai/rag-history');
      toast('已清空', 'success');
      loadHistory();
    } catch (err) {
      toast(err.message, 'error');
    }
  };

  function escapeHtml(str) {
    const div = document.createElement('div');
    div.textContent = str || '';
    return div.innerHTML;
  }
})();
