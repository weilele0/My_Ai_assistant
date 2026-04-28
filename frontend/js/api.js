// API 配置文件
const API_BASE = 'http://localhost:8080/api/v1';

// Token 管理
const Token = {
  get() {
    return localStorage.getItem('token') || '';
  },
  set(token) {
    localStorage.setItem('token', token);
  },
  clear() {
    localStorage.removeItem('token');
  },
  has() {
    return !!this.get();
  }
};

// 通用 HTTP 工具
async function api(endpoint, options = {}) {
  const url = `${API_BASE}${endpoint}`;
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers
  };

  const token = Token.get();
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const res = await fetch(url, {
    ...options,
    headers
  });

  const data = await res.json().catch(() => ({}));

  if (!res.ok) {
    throw new Error(data.error || `请求失败 (${res.status})`);
  }

  return data;
}

// 快捷方法
const GET  = (url, opts) => api(url, { ...opts, method: 'GET' });
const POST = (url, body, opts) => api(url, { ...opts, method: 'POST', body: JSON.stringify(body) });
const DEL  = (url, opts) => api(url, { ...opts, method: 'DELETE' });
