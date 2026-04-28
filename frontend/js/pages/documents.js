// Documents 模块 — 文档管理
(function () {
  const listEl = document.getElementById('doc-list');

  window.loadDocuments = loadDocuments;

  function loadDocuments() {
    listEl.innerHTML = `<div style="padding:40px;text-align:center;color:var(--text-tertiary)">${loadingEl()}</div>`;
    GET('/documents')
      .then(data => renderDocs(data.data || []))
      .catch(err => {
        listEl.innerHTML = `<div class="empty-state"><div class="empty-state-title">加载失败</div><div class="empty-state-desc">${err.message}</div></div>`;
      });
  }

  function renderDocs(docs) {
    if (docs.length === 0) {
      listEl.innerHTML = `
        <div class="empty-state">
          <div class="empty-state-icon">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
              <polyline points="14 2 14 8 20 8"/>
            </svg>
          </div>
          <div class="empty-state-title">暂无文档</div>
          <div class="empty-state-desc">上传文档后，即可基于文档进行 RAG 问答</div>
        </div>`;
      return;
    }

    listEl.innerHTML = docs.map(doc => `
      <div class="doc-item" data-id="${doc.id}">
        <div class="doc-icon">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
            <polyline points="14 2 14 8 20 8"/>
          </svg>
        </div>
        <div class="doc-info">
          <div class="doc-title">${escapeHtml(doc.title)}</div>
          <div class="doc-meta">${fmtTime(doc.created_at)} · ${doc.content ? doc.content.length + ' 字' : '无内容'}</div>
        </div>
        <div class="doc-actions">
          <button class="btn btn-ghost btn-sm" data-view="${doc.id}" title="查看">${viewIcon()}</button>
          <button class="btn btn-danger btn-sm" data-del="${doc.id}" title="删除">${delIcon()}</button>
        </div>
      </div>`).join('');

    listEl.querySelectorAll('[data-view]').forEach(btn => {
      btn.addEventListener('click', (e) => { e.stopPropagation(); viewDoc(btn.dataset.view); });
    });
    listEl.querySelectorAll('[data-del]').forEach(btn => {
      btn.addEventListener('click', (e) => { e.stopPropagation(); deleteDoc(btn.dataset.del); });
    });
  }

  function viewIcon() {
    return `<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>`;
  }

  function delIcon() {
    return `<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14H6L5 6"/><path d="M10 11v6M14 11v6"/><path d="M9 6V4h6v2"/></svg>`;
  }

  async function viewDoc(id) {
    try {
      const data = await GET(`/documents/${id}`);
      const doc = data.data;
      showModal(doc.title, `
        <div style="max-height:400px;overflow-y:auto;line-height:1.8;font-size:13.5px;white-space:pre-wrap">${escapeHtml(doc.content || '（无内容）')}</div>
      `, `<button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">关闭</button>`);
    } catch (err) {
      toast(err.message, 'error');
    }
  }

  async function deleteDoc(id) {
    if (!confirm('确定要删除这篇文档吗？')) return;
    try {
      await DEL(`/documents/${id}`);
      toast('删除成功', 'success');
      loadDocuments();
    } catch (err) {
      toast(err.message, 'error');
    }
  }

  window.openUploadModal = function () {
    showModal('上传文档', `
      <div style="display:flex;flex-direction:column;gap:16px">
        <div class="form-group">
          <label class="form-label">文档标题</label>
          <input class="form-input" id="up-title" placeholder="输入文档标题" />
        </div>
        <div class="form-group">
          <label class="form-label">文档内容</label>
          <textarea class="form-textarea" id="up-content" placeholder="粘贴文档正文内容…" style="min-height:200px"></textarea>
        </div>
      </div>`, `
      <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">取消</button>
      <button class="btn btn-primary" id="up-submit">上传</button>`);

    document.getElementById('up-submit').addEventListener('click', async () => {
      const title = document.getElementById('up-title').value.trim();
      const content = document.getElementById('up-content').value.trim();
      if (!title) { toast('请输入标题', 'error'); return; }
      if (!content) { toast('请输入内容', 'error'); return; }

      const btn = document.getElementById('up-submit');
      btn.disabled = true; btn.textContent = loadingEl();

      try {
        await POST('/documents', { title, content });
        toast('上传成功，正在向量化…', 'success');
        document.querySelector('.modal-overlay').remove();
        loadDocuments();
      } catch (err) {
        toast(err.message, 'error');
        btn.disabled = false; btn.textContent = '上传';
      }
    });
  };

  function escapeHtml(str) {
    const div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
  }
})();
