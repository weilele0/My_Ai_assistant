// Auth 模块 — 登录 / 注册
(function () {
  const form = document.getElementById('auth-form');
  const tabLogin = document.getElementById('tab-login');
  const tabRegister = document.getElementById('tab-register');
  const submitBtn = document.getElementById('auth-submit');
  const btnText = document.getElementById('auth-btn-text');
  const errorEl = document.getElementById('auth-error');

  let mode = 'login';

  function switchMode(m) {
    mode = m;
    tabLogin.classList.toggle('active', m === 'login');
    tabRegister.classList.toggle('active', m === 'register');
    submitBtn.querySelector('.btn-text').textContent = m === 'login' ? '登录' : '注册';
    document.getElementById('auth-link-text').textContent =
      m === 'login' ? '还没有账号？' : '已有账号？';
    document.getElementById('auth-link-btn').textContent = m === 'login' ? '立即注册' : '去登录';
    errorEl.textContent = '';
    form.reset();
  }

  tabLogin.addEventListener('click', () => switchMode('login'));
  tabRegister.addEventListener('click', () => switchMode('register'));
  document.getElementById('auth-link-btn').addEventListener('click', () => {
    switchMode(mode === 'login' ? 'register' : 'login');
  });

  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    errorEl.textContent = '';
    submitBtn.disabled = true;
    btnText.innerHTML = loadingEl();

    const username = document.getElementById('auth-username').value.trim();
    const password = document.getElementById('auth-password').value;
    const email = document.getElementById('auth-email').value.trim();

    try {
      if (mode === 'login') {
        const data = await POST('/user/login', { username, password });
        const { token, id, username: uname, is_admin } = data.data;
        Token.set(token);
        localStorage.setItem('uid', id);
        // 优先用服务器返回的用户名，没有则用表单填写的
        localStorage.setItem('username', uname || username);
        localStorage.setItem('is_admin', String(is_admin));

        // 从 quota 接口获取真实套餐
        try {
          const quotaRes = await GET('/quota');
          localStorage.setItem('plan', quotaRes.data?.plan || 'free');
        } catch (e) {
          localStorage.setItem('plan', 'free');
        }

        toast('登录成功', 'success');
        // 登录成功后：显示主界面 + 切换 URL 到 /chat
        showApp();
        navigateTo('chat');
      } else {
        await POST('/user/register', { username, email, password });
        toast('注册成功，请登录', 'success');
        switchMode('login');
        document.getElementById('auth-username').value = username;
        document.getElementById('auth-password').value = '';
        document.getElementById('auth-email').value = '';
      }
    } catch (err) {
      errorEl.textContent = err.message;
    } finally {
      submitBtn.disabled = false;
      btnText.textContent = mode === 'login' ? '登录' : '注册';
    }
  });
})();
