// Tasks 模块 — 异步任务（列表 + 进度轮询）
(function () {
  const form = document.getElementById('task-form');
  const topicInput = document.getElementById('task-topic');
  const submitBtn = document.getElementById('task-submit');
  const taskListEl = document.getElementById('task-list');

  let pollingTimer = null;

  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    const topic = topicInput.value.trim();
    if (!topic) return;

    submitBtn.disabled = true;
    submitBtn.querySelector('.btn-text').innerHTML = loadingEl();

    try {
      const res = await POST('/tasks/generate', { topic });
      toast('任务已提交，正在生成…', 'success');
      topicInput.value = '';
      await loadTasks();
      startPolling();
    } catch (err) {
      toast(err.message, 'error');
    } finally {
      submitBtn.disabled = false;
      submitBtn.querySelector('.btn-text').textContent = '提交任务';
    }
  });

  window.loadTasks = loadTasks;

  async function loadTasks() {
    if (!Token.has()) return;
    try {
      const res = await GET('/tasks');
      const tasks = res.data || [];
      renderTasks(tasks);
    } catch (err) {
      taskListEl.innerHTML = `
        <div class="empty-state">
          <div class="empty-state-icon">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
          </div>
          <div class="empty-state-title">加载失败</div>
          <div class="empty-state-desc">${err.message}</div>
        </div>`;
    }
  }

  function renderTasks(tasks) {
    if (tasks.length === 0) {
      taskListEl.innerHTML = `
        <div class="empty-state">
          <div class="empty-state-icon">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
          </div>
          <div class="empty-state-title">暂无任务</div>
          <div class="empty-state-desc">提交主题后，任务将在后台处理，完成后可在此查看结果</div>
        </div>`;
      return;
    }

    taskListEl.innerHTML = tasks.map(task => renderTaskCard(task)).join('');

    // 绑定展开/收起事件
    taskListEl.querySelectorAll('.task-expand-btn').forEach(btn => {
      btn.addEventListener('click', (e) => {
        const card = e.target.closest('.task-item');
        const content = card.querySelector('.task-result');
        const isExpanded = card.classList.contains('expanded');

        if (isExpanded) {
          card.classList.remove('expanded');
          content.style.maxHeight = '0';
          e.target.innerHTML = expandIcon() + ' 展开查看';
        } else {
          card.classList.add('expanded');
          content.style.maxHeight = content.scrollHeight + 'px';
          e.target.innerHTML = collapseIcon() + ' 收起';
        }
      });
    });

    // 绑定保存到文档事件
    taskListEl.querySelectorAll('.task-save-btn').forEach(btn => {
      btn.addEventListener('click', (e) => {
        const taskId = btn.dataset.id;
        const taskTopic = btn.dataset.topic;
        const taskContent = btn.dataset.content;
        openSaveDocModal(taskId, taskTopic, taskContent);
      });
    });
  }

  function renderTaskCard(task) {
    const statusMap = {
      pending:    { label: '等待中',  cls: 'status-pending',   dot: '⏳' },
      processing: { label: '生成中…', cls: 'status-processing', dot: '⚙️' },
      completed:  { label: '已完成',  cls: 'status-completed',  dot: '✅' },
      failed:     { label: '失败',    cls: 'status-failed',     dot: '❌' },
    };
    const s = statusMap[task.status] || statusMap.pending;
    const time = fmtTime(task.created_at);
    const hasResult = task.status === 'completed' && task.output_text;

    // 截取预览内容
    const previewText = hasResult && task.output_text.length > 150
      ? task.output_text.substring(0, 150) + '…'
      : '';

    // 错误信息
    const errorHtml = task.status === 'failed'
      ? `<div style="margin-top:var(--space-3);font-size:12px;color:var(--color-error)">${escapeHtml(task.error_msg || '生成失败')}</div>`
      : '';

    // 操作按钮
    let actionBtns = '';
    if (task.status === 'completed' && hasResult) {
      actionBtns = `
        <button class="btn btn-primary btn-sm task-save-btn"
          data-id="${task.id}"
          data-topic="${escapeHtmlAttr(task.input_text)}"
          data-content="${escapeHtmlAttr(task.output_text)}">
          ${saveIcon()} 保存到文档
        </button>
        <button class="btn btn-ghost btn-sm task-expand-btn">
          ${expandIcon()} 展开查看
        </button>`;
    }

    return `
      <div class="task-item ${s.cls}">
        <div class="task-header">
          <div class="task-meta">
            <span class="task-dot">${s.dot}</span>
            <span class="task-label">${s.label}</span>
            <span class="task-time">${time}</span>
          </div>
          <div class="task-actions">
            ${actionBtns}
          </div>
        </div>
        <div class="task-topic">${escapeHtml(task.input_text)}</div>
        ${hasResult ? `
          <div class="task-result" style="max-height:0;overflow:hidden;transition:max-height 0.3s ease">
            <div class="task-result-content">${escapeHtml(task.output_text)}</div>
          </div>
        ` : ''}
        ${errorHtml}
      </div>`;
  }

  // 保存到文档弹窗
  window.openSaveDocModal = function(taskId, topic, content) {
    showModal('保存到文档', `
      <div style="display:flex;flex-direction:column;gap:16px">
        <div class="form-group">
          <label class="form-label">文档标题</label>
          <input class="form-input" id="save-doc-title" value="${escapeHtml(topic)}" placeholder="输入文档标题" />
        </div>
        <div class="form-group">
          <label class="form-label">文档内容</label>
          <textarea class="form-textarea" id="save-doc-content" placeholder="文档内容" style="min-height:200px">${escapeHtml(content)}</textarea>
        </div>
      </div>`, `
      <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">取消</button>
      <button class="btn btn-primary" id="save-doc-submit">保存</button>`);

    document.getElementById('save-doc-submit').addEventListener('click', async () => {
      const title = document.getElementById('save-doc-title').value.trim();
      const docContent = document.getElementById('save-doc-content').value.trim();

      if (!title) { toast('请输入标题', 'error'); return; }
      if (!docContent) { toast('请输入内容', 'error'); return; }

      const btn = document.getElementById('save-doc-submit');
      btn.disabled = true;
      btn.innerHTML = loadingEl();

      try {
        await POST('/documents', { title, content: docContent });
        toast('保存成功！', 'success');
        document.querySelector('.modal-overlay').remove();
      } catch (err) {
        toast(err.message, 'error');
        btn.disabled = false;
        btn.textContent = '保存';
      }
    });
  };

  // 图标
  function expandIcon() {
    return `<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="6 9 12 15 18 9"/></svg>`;
  }
  function collapseIcon() {
    return `<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="18 15 12 9 6 15"/></svg>`;
  }
  function saveIcon() {
    return `<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/><polyline points="17 21 17 13 7 13 7 21"/><polyline points="7 3 7 8 15 8"/></svg>`;
  }

  function escapeHtml(str) {
    if (!str) return '';
    const div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
  }

  function escapeHtmlAttr(str) {
    if (!str) return '';
    return str.replace(/"/g, '&quot;').replace(/'/g, '&#39;');
  }

  // 轮询
  function startPolling() {
    if (pollingTimer) return;
    pollingTimer = setInterval(async () => {
      try {
        const res = await GET('/tasks');
        const tasks = res.data || [];
        renderTasks(tasks);
        const hasActive = tasks.some(t => t.status === 'pending' || t.status === 'processing');
        if (!hasActive) {
          clearInterval(pollingTimer);
          pollingTimer = null;
        }
      } catch {
        clearInterval(pollingTimer);
        pollingTimer = null;
      }
    }, 2000);
  }

  window.addEventListener('beforeunload', () => {
    if (pollingTimer) clearInterval(pollingTimer);
  });
})();
