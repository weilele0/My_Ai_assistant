// Chat 模块 — AI 对话 / RAG 问答
(function () {
  const messagesEl = document.getElementById('chat-messages');
  const textarea = document.getElementById('chat-textarea');
  const sendBtn = document.getElementById('chat-send');
  const modeTabs = document.getElementById('mode-tabs');
  const quotaMini = document.getElementById('quota-mini');

  let currentMode = 'generate';
  let conversation = [];
  let isLoading = false;

  // modeTabs 事件监听（如果元素存在才挂载）
  // 注意：主要通过 index.html 的事件委托来处理，这里仅作备用
  if (modeTabs) {
    modeTabs.addEventListener('click', (e) => {
      if (e.target.classList.contains('tab')) {
        modeTabs.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
        e.target.classList.add('active');
        currentMode = e.target.dataset.mode;
        clearChat();
      }
    });
  }

  // 外部切换模式（由 index.html 事件委托调用）
  window.setChatMode = function(mode) {
    currentMode = mode;
    clearChat();
    showEmpty(); // 清空后立即显示对应模式的空状态
  };

  textarea.addEventListener('input', () => {
    textarea.style.height = 'auto';
    textarea.style.height = Math.min(textarea.scrollHeight, 160) + 'px';
  });

  textarea.addEventListener('keydown', (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  });

  sendBtn.addEventListener('click', handleSend);

  async function handleSend() {
    const text = textarea.value.trim();
    if (!text || isLoading) return;

    if (!Token.has()) {
      toast('请先登录', 'error');
      router('auth');
      return;
    }

    textarea.value = '';
    textarea.style.height = 'auto';

    addMessage('user', text);
    conversation.push({ role: 'user', content: text });

    const placeholderId = addMessage('assistant', '', true);
    isLoading = true;
    sendBtn.disabled = true;

    try {
      let endpoint, body;
      if (currentMode === 'generate') {
        endpoint = '/ai/generate';
        body = { topic: text };
      } else {
        endpoint = '/ai/rag-generate';
        body = { question: text, top_k: 5 };
      }

      const data = await POST(endpoint, body);

      const ph = document.getElementById(placeholderId);
      if (ph) ph.remove();

      if (data.data && data.data.content) {
        addMessage('assistant', data.data.content);
        conversation.push({ role: 'assistant', content: data.data.content });
      } else if (data.data && data.data.answer) {
        addMessage('assistant', data.data.answer, false, data.data.references);
        conversation.push({ role: 'assistant', content: data.data.answer });
      }

      refreshQuotaMini();
    } catch (err) {
      const ph = document.getElementById(placeholderId);
      if (ph) {
        ph.querySelector('.msg-bubble').textContent = `请求失败：${err.message}`;
        ph.querySelector('.msg-bubble').style.color = 'var(--color-error)';
      }
    } finally {
      isLoading = false;
      sendBtn.disabled = false;
      messagesEl.scrollTop = messagesEl.scrollHeight;
    }
  }

  function addMessage(role, content, isPlaceholder = false, references) {
    const id = 'msg-' + Date.now() + Math.random();
    const time = fmtTime(new Date().toISOString());
    const initial = role === 'user' ? currentUser().charAt(0).toUpperCase() : 'A';
    const msgEl = document.createElement('div');
    msgEl.className = `message ${role}`;
    msgEl.id = id;

    let refsHTML = '';
    if (references && references.length > 0) {
      refsHTML = `
        <div class="msg-ref">
          <div class="msg-ref-label">参考资料</div>
          ${references.slice(0, 3).map((r, i) =>
            `<div style="margin-bottom:6px">【${i+1}】${r.length > 120 ? r.slice(0,120)+'…' : r}</div>`
          ).join('')}
        </div>`;
    }

    msgEl.innerHTML = `
      <div class="msg-avatar ${role}">${initial}</div>
      <div class="msg-content">
        <div class="msg-bubble">${isPlaceholder ? loadingEl() : escapeHtml(content)}</div>
        ${refsHTML}
        <div class="msg-time">${time}</div>
      </div>`;

    messagesEl.appendChild(msgEl);
    messagesEl.scrollTop = messagesEl.scrollHeight;
    return id;
  }

  function escapeHtml(str) {
    const div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
  }

  function clearChat() {
    messagesEl.innerHTML = '';
    conversation = [];
  }

  function loadChatState() {
    if (messagesEl.children.length === 0) {
      showEmpty();
    }
    refreshQuotaMini();
  }

  function showEmpty() {
    messagesEl.innerHTML = `
      <div class="chat-empty">
        <div class="chat-empty-icon">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
          </svg>
        </div>
        <div class="chat-empty-title">开始对话</div>
        <div class="chat-empty-desc">
          ${currentMode === 'generate'
            ? '输入任意主题，AI 将为你生成高质量文章'
            : '基于已上传的文档进行智能问答'}
        </div>
      </div>`;
  }

  async function refreshQuotaMini() {
    if (!Token.has()) return;
    try {
      const data = await GET('/quota/me');
      const q = data.data;
      quotaMini.innerHTML = `
        <span>AI: ${q.remain_ai < 0 ? '∞' : q.remain_ai} 次</span>
        <span style="color:var(--border-medium)">|</span>
        <span>RAG: ${q.remain_rag < 0 ? '∞' : q.remain_rag} 次</span>`;
    } catch {}
  }

  window.loadChatState = loadChatState;
})();
