const el = document.getElementById('app');
const cfg = window.__APP__ || { tenant: 'checkcall', apiBase: '/api' };

function tgInitData() {
  // Telegram WebApp
  if (window.Telegram?.WebApp?.initData) return window.Telegram.WebApp.initData;
  // Debug: можно подставить руками в браузере при тестах
  return window.__INIT_DATA__ || '';
}

async function api(path, opts = {}) {
  const res = await fetch(cfg.apiBase + path, {
    credentials: 'include',            // важно: чтобы cookie sid отправлялась
    headers: { 'Content-Type': 'application/json', ...(opts.headers || {}) },
    ...opts,
  });

  const ct = res.headers.get('content-type') || '';
  const body = ct.includes('application/json') ? await res.json() : await res.text();

  if (!res.ok) {
    const msg = typeof body === 'string' ? body : (body?.error || JSON.stringify(body));
    throw new Error(`${res.status} ${msg}`);
  }
  return body;
}

async function loginTelegram() {
  const initData = tgInitData();
  if (!initData) {
    throw new Error('Нет Telegram initData. Открой страницу из Telegram WebApp (или задай window.__INIT_DATA__ для теста).');
  }
  await api('/auth/telegram', {
    method: 'POST',
    body: JSON.stringify({ initData }),
  });
}

async function me() {
  return api('/me');
}

(async () => {
  try {
    // 1) логинимся через Telegram (ставится cookie sid)
    await loginTelegram();

    // 2) проверяем, что сессия работает
    const data = await me();

    el.textContent = `Tenant: ${cfg.tenant} | AUTH: OK | sid=${data.sid}`;
  } catch (e) {
    el.textContent = 'Ошибка: ' + (e?.message || e);
  }
})();
