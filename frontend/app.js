const API_BASE = localStorage.getItem("apiBase") || "http://localhost:3000/api";
const TOKEN_KEY = "currency_exchange_token";

const authStatusEl = document.getElementById("authStatus");
const toastEl = document.getElementById("toast");
const rateCountEl = document.getElementById("rateCount");
const articleCountEl = document.getElementById("articleCount");
const summaryAuthEl = document.getElementById("summaryAuth");

const registerForm = document.getElementById("registerForm");
const loginForm = document.getElementById("loginForm");
const logoutBtn = document.getElementById("logoutBtn");

const exchangeForm = document.getElementById("exchangeForm");
const refreshRatesBtn = document.getElementById("refreshRatesBtn");
const exchangeRatesEl = document.getElementById("exchangeRates");
const converterPairEl = document.getElementById("converterPair");
const converterAmountEl = document.getElementById("converterAmount");
const converterResultEl = document.getElementById("converterResult");
const reverseDirectionBtn = document.getElementById("reverseDirectionBtn");
const keepResultModeEl = document.getElementById("keepResultMode");

const articleForm = document.getElementById("articleForm");
const refreshArticlesBtn = document.getElementById("refreshArticlesBtn");
const articleListEl = document.getElementById("articleList");
const articleDetailEl = document.getElementById("articleDetail");

let latestRates = [];
let useInvertedRate = false;

function getCurrentConversionContext() {
    if (!Array.isArray(latestRates) || latestRates.length === 0 || !converterPairEl) {
        return null;
    }

    const index = Number(converterPairEl.value || 0);
    const selectedRate = latestRates[index];
    if (!selectedRate) {
        return null;
    }

    const rateValue = Number(selectedRate.rate);
    const effectiveRate = useInvertedRate ? 1 / rateValue : rateValue;
    if (!Number.isFinite(effectiveRate) || effectiveRate <= 0) {
        return null;
    }

    const amount = Number(converterAmountEl?.value);
    if (!Number.isFinite(amount)) {
        return null;
    }

    return {
        amount,
        converted: amount * effectiveRate,
    };
}

function getToken() {
    return localStorage.getItem(TOKEN_KEY) || "";
}

function setToken(token) {
    if (token) {
        localStorage.setItem(TOKEN_KEY, token);
    } else {
        localStorage.removeItem(TOKEN_KEY);
    }
    renderAuthStatus();
}

function renderAuthStatus() {
    const loggedIn = !!getToken();
    authStatusEl.textContent = loggedIn ? "已登录" : "未登录";
    if (summaryAuthEl) {
        summaryAuthEl.textContent = loggedIn ? "认证用户" : "游客";
    }
}

function showToast(message, isError = false) {
    toastEl.textContent = message;
    toastEl.style.background = isError ? "#b42318" : "#111827";
    toastEl.classList.add("show");
    setTimeout(() => toastEl.classList.remove("show"), 1800);
}

function getFormData(formElement) {
    return Object.fromEntries(new FormData(formElement).entries());
}

function formatAmount(value) {
    return Number(value).toLocaleString("zh-CN", {
        minimumFractionDigits: 2,
        maximumFractionDigits: 6,
    });
}

function updateConverterOptions(rates) {
    if (!converterPairEl || !converterResultEl) {
        return;
    }

    if (!Array.isArray(rates) || rates.length === 0) {
        converterPairEl.innerHTML = '<option value="">暂无可用汇率</option>';
        converterPairEl.disabled = true;
        converterResultEl.textContent = "请先新增或刷新汇率数据";
        return;
    }

    converterPairEl.disabled = false;
    useInvertedRate = false;
    converterPairEl.innerHTML = rates
        .map((item, index) => {
            const optionLabel = `${item.fromCurrency} → ${item.toCurrency} (1:${item.rate})`;
            return `<option value="${index}">${optionLabel}</option>`;
        })
        .join("");
}

function calculateConversion() {
    if (!converterPairEl || !converterAmountEl || !converterResultEl) {
        return;
    }

    if (!Array.isArray(latestRates) || latestRates.length === 0) {
        converterResultEl.textContent = "请先新增或刷新汇率数据";
        return;
    }

    const index = Number(converterPairEl.value || 0);
    const selectedRate = latestRates[index];
    const amount = Number(converterAmountEl.value);

    if (!selectedRate) {
        converterResultEl.textContent = "请选择可用币种对";
        return;
    }

    const fromCurrency = useInvertedRate ? selectedRate.toCurrency : selectedRate.fromCurrency;
    const toCurrency = useInvertedRate ? selectedRate.fromCurrency : selectedRate.toCurrency;
    const rateValue = Number(selectedRate.rate);
    const effectiveRate = useInvertedRate ? 1 / rateValue : rateValue;

    if (!Number.isFinite(effectiveRate) || effectiveRate <= 0) {
        converterResultEl.textContent = "当前汇率不可用于换算";
        return;
    }

    if (!Number.isFinite(amount)) {
        converterResultEl.textContent = `当前汇率：1 ${fromCurrency} = ${effectiveRate.toFixed(6)} ${toCurrency}`;
        return;
    }

    const converted = amount * effectiveRate;
    converterResultEl.textContent = `${formatAmount(amount)} ${fromCurrency} ≈ ${formatAmount(converted)} ${toCurrency}`;
}

function reverseConversionDirection() {
    if (!Array.isArray(latestRates) || latestRates.length === 0 || !converterPairEl) {
        showToast("请先加载汇率数据", true);
        return;
    }

    const currentIndex = Number(converterPairEl.value || 0);
    const currentRate = latestRates[currentIndex];

    if (!currentRate) {
        showToast("未找到当前币种对", true);
        return;
    }

    const shouldKeepResult = keepResultModeEl?.checked;
    const beforeReverse = shouldKeepResult ? getCurrentConversionContext() : null;

    const reverseIndex = latestRates.findIndex(
        (item) =>
            item.fromCurrency === currentRate.toCurrency &&
            item.toCurrency === currentRate.fromCurrency
    );

    if (reverseIndex >= 0) {
        converterPairEl.value = String(reverseIndex);
        useInvertedRate = false;
        if (beforeReverse && converterAmountEl) {
            converterAmountEl.value = beforeReverse.converted.toFixed(6);
        }
        calculateConversion();
        showToast("已切换为反向汇率");
        return;
    }

    const rateValue = Number(currentRate.rate);
    if (!Number.isFinite(rateValue) || rateValue <= 0) {
        showToast("当前汇率不可反转", true);
        return;
    }

    useInvertedRate = !useInvertedRate;
    if (beforeReverse && converterAmountEl) {
        converterAmountEl.value = beforeReverse.converted.toFixed(6);
    }
    calculateConversion();
    showToast("已使用倒数汇率反转方向");
}

async function request(path, options = {}) {
    const token = getToken();
    const headers = {
        "Content-Type": "application/json",
        ...(options.headers || {}),
    };

    if (token) {
        headers.Authorization = token;
    }

    const response = await fetch(`${API_BASE}${path}`, {
        ...options,
        headers,
    });

    const contentType = response.headers.get("content-type") || "";
    const data = contentType.includes("application/json")
        ? await response.json()
        : await response.text();

    if (!response.ok) {
        const message =
            (typeof data === "object" && data && data.error) ||
            `请求失败: ${response.status}`;
        throw new Error(message);
    }

    return data;
}

async function register(payload) {
    const data = await request("/auth/register", {
        method: "POST",
        body: JSON.stringify(payload),
    });
    setToken(data.token || "");
    showToast("注册成功");
}

async function login(payload) {
    const data = await request("/auth/login", {
        method: "POST",
        body: JSON.stringify(payload),
    });
    setToken(data.token || "");
    showToast("登录成功");
}

async function createExchangeRate(payload) {
    payload.rate = Number(payload.rate);
    await request("/exchangeRates", {
        method: "POST",
        body: JSON.stringify(payload),
    });
    showToast("汇率新增成功");
}

async function getExchangeRates() {
    return request("/exchangeRates");
}

async function createArticle(payload) {
    await request("/articles", {
        method: "POST",
        body: JSON.stringify(payload),
    });
    showToast("文章发布成功");
}

async function getArticles() {
    return request("/articles");
}

async function getArticleById(id) {
    return request(`/articles/${id}`);
}

async function likeArticle(id) {
    await request(`/articles/${id}/like`, { method: "POST" });
}

async function getArticleLikes(id) {
    return request(`/articles/${id}/like`);
}

function renderExchangeRates(rates) {
    latestRates = Array.isArray(rates) ? rates : [];

    if (rateCountEl) {
        rateCountEl.textContent = Array.isArray(rates) ? String(rates.length) : "0";
    }

    updateConverterOptions(latestRates);
    calculateConversion();

    if (!Array.isArray(rates) || rates.length === 0) {
        exchangeRatesEl.innerHTML = '<div class="item">暂无汇率数据</div>';
        return;
    }

    exchangeRatesEl.innerHTML = rates
        .map((item) => {
            const dateText = item.date ? new Date(item.date).toLocaleString() : "-";
            return `<div class="item">
          <h4>${item.fromCurrency} → ${item.toCurrency}</h4>
          <p>Rate: ${item.rate}</p>
          <p>更新时间: ${dateText}</p>
        </div>`;
        })
        .join("");
}

async function renderArticleDetail(article) {
    if (!article || !article.ID) {
        articleDetailEl.innerHTML = "未找到文章详情";
        return;
    }

    let likes = "0";
    try {
        const likeData = await getArticleLikes(article.ID);
        likes = likeData.likes || "0";
    } catch (error) {
        likes = "获取失败";
    }

    articleDetailEl.innerHTML = `
      <h4>${article.title}</h4>
      <p class="meta">点赞数：${likes}</p>
      <p><strong>导语：</strong>${article.preview}</p>
      <p>${article.content}</p>
    `;
}

function renderArticles(articles) {
    if (articleCountEl) {
        articleCountEl.textContent = Array.isArray(articles) ? String(articles.length) : "0";
    }

    if (!Array.isArray(articles) || articles.length === 0) {
        articleListEl.innerHTML = '<div class="item">暂无文章</div>';
        return;
    }

    articleListEl.innerHTML = articles
        .map(
            (item) => `<div class="item" data-id="${item.ID}">
        <h4>${item.title}</h4>
        <p>${item.preview}</p>
        <div class="row">
          <button class="secondary" data-action="view">查看详情</button>
          <button data-action="like">点赞</button>
        </div>
      </div>`
        )
        .join("");
}

async function refreshExchangeRates() {
    try {
        const data = await getExchangeRates();
        renderExchangeRates(data);
    } catch (error) {
        latestRates = [];
        updateConverterOptions([]);
        calculateConversion();
        showToast(error.message, true);
    }
}

async function refreshArticles() {
    try {
        const data = await getArticles();
        renderArticles(data);
    } catch (error) {
        articleListEl.innerHTML = '<div class="item">暂无文章</div>';
        showToast(error.message, true);
    }
}

registerForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    const payload = getFormData(registerForm);
    try {
        await register(payload);
        registerForm.reset();
    } catch (error) {
        showToast(error.message, true);
    }
});

loginForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    const payload = getFormData(loginForm);
    try {
        await login(payload);
        loginForm.reset();
    } catch (error) {
        showToast(error.message, true);
    }
});

logoutBtn.addEventListener("click", () => {
    setToken("");
    showToast("已退出登录");
});

exchangeForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    const payload = getFormData(exchangeForm);
    try {
        await createExchangeRate(payload);
        exchangeForm.reset();
        await refreshExchangeRates();
    } catch (error) {
        showToast(error.message, true);
    }
});

refreshRatesBtn.addEventListener("click", refreshExchangeRates);

converterPairEl.addEventListener("change", calculateConversion);
converterAmountEl.addEventListener("input", calculateConversion);
reverseDirectionBtn.addEventListener("click", reverseConversionDirection);

articleForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    const payload = getFormData(articleForm);
    try {
        await createArticle(payload);
        articleForm.reset();
        await refreshArticles();
    } catch (error) {
        showToast(error.message, true);
    }
});

refreshArticlesBtn.addEventListener("click", refreshArticles);

articleListEl.addEventListener("click", async (event) => {
    const target = event.target;
    if (!(target instanceof HTMLButtonElement)) {
        return;
    }

    const action = target.dataset.action;
    const item = target.closest(".item");
    const articleId = item?.dataset.id;

    if (!articleId) {
        return;
    }

    try {
        if (action === "view") {
            const article = await getArticleById(articleId);
            await renderArticleDetail(article);
            return;
        }

        if (action === "like") {
            await likeArticle(articleId);
            showToast("点赞成功");
            const article = await getArticleById(articleId);
            await renderArticleDetail(article);
        }
    } catch (error) {
        showToast(error.message, true);
    }
});

renderAuthStatus();
refreshExchangeRates();
refreshArticles();
