// 显示主界面（登录成功后调用）
function showApp() {
  document.getElementById('page-auth').classList.remove('active');
  document.getElementById('app-layout').style.display = '';
  updateUserCard();
}

// 路由到指定页面，并更新 URL
function navigateTo(page) {
  if (!Token.has()) {
    navigateTo('auth');
    return;
  }
  router(page);
  // 更新 URL（pushState 不刷新页面）
  const path = page === 'chat' ? '/' : '/' + page;
  history.pushState({ page }, '', path);
}

// 路由与页面切换
function router(page) {
  // 未登录不允许进入聊天等页面
  if (page !== 'auth' && !Token.has()) {
    router('auth');
    return;
  }

  document.querySelectorAll('.page').forEach(p => p.classList.remove('active'));
  document.querySelectorAll('.nav-item').forEach(n => n.classList.remove('active'));
  const nav = document.querySelector(`.nav-item[data-page="${page}"]`);
  if (nav) nav.classList.add('active');

  const target = document.getElementById(`page-${page}`);
  if (target) target.classList.add('active');

  if (page !== 'auth') updateUserCard();
  if (page === 'chat') loadChatState();
  if (page === 'documents') loadDocuments();
  if (page === 'tasks') loadTasks();
  if (page === 'quota') loadQuota();
  if (page === 'admin') loadAdmin();
  if (page === 'history') loadHistory();
}

// 更新侧边栏用户信息
function updateUserCard() {
  const username = localStorage.getItem('username') || '用户';
  const plan = localStorage.getItem('plan') || 'free';
  const el = document.getElementById('user-card');
  if (!el) return;
  el.querySelector('.user-name').textContent = username;
  el.querySelector('.user-plan').textContent = planName(plan);
  el.querySelector('.user-avatar').textContent = username.charAt(0).toUpperCase();

  const adminNav = document.querySelector('[data-page="admin"]');
  const isAdmin = localStorage.getItem('is_admin') === 'true';
  if (adminNav) adminNav.style.display = isAdmin ? '' : 'none';
}

// 套餐名映射
function planName(plan) {
  const map = { free: '免费版', basic: '基础版', pro: '专业版', unlimit: '无限制' };
  return map[plan] || '免费版';
}

// Toast 提示
function toast(msg, type = 'info') {
  let container = document.querySelector('.toast-container');
  if (!container) {
    container = document.createElement('div');
    container.className = 'toast-container';
    document.body.appendChild(container);
  }
  const t = document.createElement('div');
  t.className = `toast ${type}`;
  t.innerHTML = `<span>${msg}</span>`;
  container.appendChild(t);
  setTimeout(() => t.remove(), 3000);
}

// 检查登录状态，未登录则跳转
function requireAuth() {
  if (!Token.has()) {
    router('auth');
    return false;
  }
  return true;
}

// 登出
function logout() {
  Token.clear();
  localStorage.removeItem('uid');
  localStorage.removeItem('username');
  localStorage.removeItem('is_admin');
  localStorage.removeItem('plan');
  document.getElementById('app-layout').style.display = 'none';
  document.getElementById('page-auth').classList.add('active');
  history.pushState({ page: 'auth' }, '', '/login');
}

// 当前用户名
function currentUser() {
  return localStorage.getItem('username') || 'U';
}

// 格式化时间
function fmtTime(iso) {
  if (!iso) return '';
  const d = new Date(iso);
  const now = new Date();
  const diff = now - d;
  if (diff < 60000) return '刚刚';
  if (diff < 3600000) return `${Math.floor(diff/60000)} 分钟前`;
  if (diff < 86400000) return `${Math.floor(diff/3600000)} 小时前`;
  return d.toLocaleDateString('zh-CN', { month:'short', day:'numeric' });
}

// 渲染额度条
function renderQuotaBar(used, total) {
  if (total < 0) return '<span style="color:var(--color-success)">无限制</span>';
  const pct = total > 0 ? Math.min((used / total) * 100, 100) : 0;
  const cls = pct >= 90 ? 'danger' : pct >= 60 ? 'warn' : '';
  return `
    <div class="quota-bar"><div class="quota-bar-fill ${cls}" style="width:${pct}%"></div></div>
    <span class="quota-sub">${used} / ${total} 次</span>`;
}

// 显示模态框
function showModal(title, bodyHTML, footerHTML, onClose) {
  const overlay = document.createElement('div');
  overlay.className = 'modal-overlay';
  overlay.innerHTML = `
    <div class="modal">
      <div class="modal-header">
        <span class="modal-title">${title}</span>
        <div class="modal-close" data-close>${closeIcon()}</div>
      </div>
      <div class="modal-body">${bodyHTML}</div>
      <div class="modal-footer">${footerHTML}</div>
    </div>`;
  document.body.appendChild(overlay);

  overlay.querySelector('[data-close]').addEventListener('click', () => {
    overlay.remove();
    if (onClose) onClose();
  });
  overlay.addEventListener('click', (e) => {
    if (e.target === overlay) { overlay.remove(); if (onClose) onClose(); }
  });
  return overlay;
}

function closeIcon() {
  return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 6L6 18M6 6l12 12"/></svg>`;
}

function loadingEl() {
  return `<div class="loading-dots"><span></span><span></span><span></span></div>`;
}
