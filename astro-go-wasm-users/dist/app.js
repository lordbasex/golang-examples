let wasmReady = false;
let allUsers = [];
let editingUserId = null;
let pendingDeleteId = null;

const fields = [
  "firstName", "lastName", "email", "phone", "document", "address",
  "city", "state", "zipCode", "country", "company", "role", "notes"
];

const constraints = {
  firstName: { required: true, max: 80 },
  lastName: { required: true, max: 80 },
  email: { required: true, max: 120, type: "email" },
  phone: { max: 30, type: "phone" },
  document: { max: 30 },
  address: { max: 180 },
  city: { max: 80 },
  state: { max: 80 },
  zipCode: { max: 20 },
  country: { max: 80 },
  company: { max: 120 },
  role: { max: 80 },
  notes: { max: 1000 }
};

const el = (id) => document.getElementById(id);

function getTheme() {
  return localStorage.getItem("astro_go_wasm_theme") || "light";
}

function applyTheme(theme) {
  document.documentElement.setAttribute("data-theme", theme);
  localStorage.setItem("astro_go_wasm_theme", theme);
}

function toggleTheme() {
  applyTheme(getTheme() === "dark" ? "light" : "dark");
}

function showToast(message) {
  const toast = el("toast");
  toast.textContent = message;
  toast.classList.add("is-visible");
  clearTimeout(showToast._timer);
  showToast._timer = setTimeout(() => toast.classList.remove("is-visible"), 2600);
}

function clearFieldError(field) {
  const input = el(field);
  const errorEl = document.querySelector(`[data-error="${field}"]`);
  if (input) input.classList.remove("is-invalid");
  if (errorEl) errorEl.textContent = "";
}

function clearAllErrors() {
  fields.forEach(clearFieldError);
}

function setFieldError(field, message) {
  const input = el(field);
  const errorEl = document.querySelector(`[data-error="${field}"]`);
  if (input) input.classList.add("is-invalid");
  if (errorEl) errorEl.textContent = message || "";
}

function paintErrors(errors = {}) {
  clearAllErrors();
  Object.entries(errors).forEach(([field, message]) => setFieldError(field, message));
}

function validateEmail(value) {
  if (!value) return false;
  if (value.includes(" ")) return false;
  const parts = value.split("@");
  if (parts.length !== 2) return false;
  const [local, domain] = parts;
  return !!local && !!domain && domain.includes(".") && !domain.startsWith(".") && !domain.endsWith(".");
}

function validatePhone(value) {
  if (!value) return true;
  const re = /^[0-9+\-\s()]+$/;
  const digits = value.replace(/\D/g, "");
  return re.test(value) && digits.length >= 6 && digits.length <= 20;
}

function validateField(name, value) {
  const rules = constraints[name];
  if (!rules) return "";

  const trimmed = String(value || "").trim();

  if (rules.required && !trimmed) {
    const labels = {
      firstName: "First name is required",
      lastName: "Last name is required",
      email: "Email is required"
    };
    return labels[name] || "Required field";
  }

  if (trimmed && rules.max && [...trimmed].length > rules.max) {
    return `Maximum ${rules.max} characters`;
  }

  if (trimmed && rules.type === "email" && !validateEmail(trimmed)) {
    return "Invalid email";
  }

  if (trimmed && rules.type === "phone" && !validatePhone(trimmed)) {
    return "Invalid phone";
  }

  return "";
}

function validateFrontend(payload) {
  const errors = {};
  for (const name of fields) {
    const error = validateField(name, payload[name]);
    if (error) errors[name] = error;
  }
  return errors;
}

function setFormMode(editing) {
  el("formTitle").textContent = editing ? "Edit user" : "Create user";
}

function openFormModal() {
  el("formModal").classList.remove("is-hidden");
}

function closeFormModal() {
  el("formModal").classList.add("is-hidden");
}

function resetForm() {
  editingUserId = null;
  el("userId").value = "";
  fields.forEach((name) => { el(name).value = ""; clearFieldError(name); });
  setFormMode(false);
}

function startCreateUser() {
  resetForm();
  openFormModal();
}

function formToPayload() {
  const payload = {};
  for (const name of fields) {
    payload[name] = el(name).value.trim();
  }
  return payload;
}

function fillForm(user) {
  editingUserId = user.id;
  el("userId").value = user.id;
  fields.forEach((name) => {
    el(name).value = user[name] || "";
    clearFieldError(name);
  });
  setFormMode(true);
  openFormModal();
}

function escapeHTML(value = "") {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}

function formatDate(value) {
  if (!value) return "-";
  try { return new Date(value).toLocaleString(); } catch { return value; }
}

function renderUsers(list = allUsers) {
  const container = el("usersList");
  if (!list.length) {
    container.innerHTML = `<tr><td colspan="7" class="empty-row">No users to display.</td></tr>`;
    return;
  }

  container.innerHTML = list.map((user) => `
    <tr>
      <td>${user.id}</td>
      <td>${escapeHTML(user.firstName)} ${escapeHTML(user.lastName)}</td>
      <td>${escapeHTML(user.email || "-")}</td>
      <td>${escapeHTML(user.phone || "-")}</td>
      <td>${escapeHTML(user.company || "-")}</td>
      <td>${escapeHTML(formatDate(user.updatedAt))}</td>
      <td>
        <div class="table-actions">
          <button class="btn btn--ghost btn--icon" title="View" aria-label="View" data-action="view" data-id="${user.id}">👁️</button>
          <button class="btn btn--ghost btn--icon" title="Edit" aria-label="Edit" data-action="edit" data-id="${user.id}">✏️</button>
          <button class="btn btn--danger btn--icon" title="Delete" aria-label="Delete" data-action="delete" data-id="${user.id}">🗑️</button>
        </div>
      </td>
    </tr>
  `).join("");
}

function renderDetail(user) {
  const body = el("modalBody");
  const items = [
    ["ID", user.id],
    ["First name", user.firstName],
    ["Last name", user.lastName],
    ["Email", user.email],
    ["Phone", user.phone],
    ["Document", user.document],
    ["Company", user.company],
    ["Role", user.role],
    ["Address", user.address],
    ["City", user.city],
    ["State / Province", user.state],
    ["Zip code", user.zipCode],
    ["Country", user.country],
    ["Created", formatDate(user.createdAt)],
    ["Updated", formatDate(user.updatedAt)],
    ["Notes", user.notes],
  ];

  body.innerHTML = items.map(([label, value]) => `
    <div class="detail__item">
      <span>${escapeHTML(label)}</span>
      <strong>${escapeHTML(value || "-")}</strong>
    </div>
  `).join("");

  el("modal").classList.remove("is-hidden");
}

function closeModal() {
  el("modal").classList.add("is-hidden");
}

function openDeleteModal(user) {
  pendingDeleteId = Number(user.id);
  const name = `${user.firstName || ""} ${user.lastName || ""}`.trim() || `#${pendingDeleteId}`;
  el("deleteModalText").textContent = `Are you sure you want to delete ${name}?`;
  el("deleteModal").classList.remove("is-hidden");
}

function closeDeleteModal() {
  pendingDeleteId = null;
  el("deleteModal").classList.add("is-hidden");
}

function openSideMenu() {
  el("sideMenu").classList.add("is-open");
  el("sideMenuBackdrop").classList.remove("is-hidden");
}

function closeSideMenu() {
  el("sideMenu").classList.remove("is-open");
  el("sideMenuBackdrop").classList.add("is-hidden");
}

function normalizeEmail(value) {
  return String(value || "").trim().toLowerCase();
}

function normalizeDocument(value) {
  return String(value || "").trim();
}

async function apiRequest(url, options = {}) {
  if (!wasmReady) throw new Error("WASM is not ready yet");
  const canDirectFallback = /^\/api\/users(\/\d+)?$/.test(url);
  if (!navigator.serviceWorker?.controller && canDirectFallback) {
    return apiRequestDirect(url, options);
  }

  const response = await fetch(url, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });

  const raw = await response.text();
  let data = {};
  try {
    data = raw ? JSON.parse(raw) : {};
  } catch {
    // Keep empty object when response is not JSON (e.g. HTML 404 page).
    data = {};
  }

  if (!response.ok || data.ok === false) {
    // If SW responds 404 for API routes (usually stale worker/base mismatch),
    // fallback to direct local API to keep app usable.
    if (response.status === 404 && canDirectFallback) {
      return apiRequestDirect(url, options);
    }
    const error = new Error(data.message || "request_failed");
    error.payload = data;
    throw error;
  }

  return data;
}

async function apiRequestDirect(url, options = {}) {
  const method = (options.method || "GET").toUpperCase();
  const body = options.body ? JSON.parse(options.body) : null;
  const match = /^\/api\/users\/(\d+)$/.exec(url);

  let raw;
  if (url === "/api/users" && method === "GET") {
    raw = await globalThis.apiListUsers();
  } else if (url === "/api/users" && method === "POST") {
    raw = await globalThis.apiCreateUser(body || {});
  } else if (match && method === "GET") {
    raw = await globalThis.apiGetUser(Number(match[1]));
  } else if (match && method === "PUT") {
    raw = await globalThis.apiUpdateUser({ ...(body || {}), id: Number(match[1]) });
  } else if (match && method === "DELETE") {
    raw = await globalThis.apiDeleteUser(Number(match[1]));
  } else {
    throw new Error("endpoint_not_supported");
  }

  const data = JSON.parse(raw || "{}");
  if (data.ok === false) {
    const error = new Error(data.message || "request_failed");
    error.payload = data;
    throw error;
  }
  return data;
}

async function refreshUsers() {
  const data = await apiRequest("/api/users");
  allUsers = data.data || [];
  filterUsers();
}

function filterUsers() {
  const q = el("searchInput").value.trim().toLowerCase();
  if (!q) return renderUsers(allUsers);

  const filtered = allUsers.filter((u) => {
    const haystack = [
      u.firstName, u.lastName, u.email, u.phone, u.company, u.role, u.city, u.country
    ].join(" ").toLowerCase();
    return haystack.includes(q);
  });

  renderUsers(filtered);
}

async function onSubmitForm(event) {
  event.preventDefault();

  const payload = formToPayload();
  const frontendErrors = validateFrontend(payload);

  if (!editingUserId) {
    const email = normalizeEmail(payload.email);
    if (email && allUsers.some((u) => normalizeEmail(u.email) === email)) {
      frontendErrors.email = "Email already exists";
    }

    const document = normalizeDocument(payload.document);
    if (document && allUsers.some((u) => normalizeDocument(u.document) === document)) {
      frontendErrors.document = "Document already exists";
    }
  }

  if (Object.keys(frontendErrors).length) {
    paintErrors(frontendErrors);
    return;
  }

  clearAllErrors();

  try {
    if (editingUserId) {
      const current = allUsers.find((u) => Number(u.id) === Number(editingUserId));
      await apiRequest(`/api/users/${editingUserId}`, {
        method: "PUT",
        body: JSON.stringify({ ...current, ...payload }),
      });
      showToast("User updated");
    } else {
      await apiRequest("/api/users", {
        method: "POST",
        body: JSON.stringify(payload),
      });
      showToast("User created");
    }

    resetForm();
    closeFormModal();
    await refreshUsers();
  } catch (err) {
    if (err.payload?.errors) {
      paintErrors(err.payload.errors);
      return;
    }
    showToast(err.message || "Could not save user");
  }
}

async function deleteUser(id, silent = false) {
  await apiRequest(`/api/users/${id}`, { method: "DELETE" });
  if (!silent) showToast("User deleted");
  if (Number(editingUserId) === Number(id)) resetForm();
  await refreshUsers();
}

async function onUsersListClick(event) {
  const button = event.target.closest("[data-action]");
  if (!button) return;

  const id = Number(button.dataset.id);
  const action = button.dataset.action;
  const user = allUsers.find((u) => Number(u.id) === id);
  if (!user) return;

  if (action === "view") {
    renderDetail(user);
  } else if (action === "edit") {
    fillForm(user);
  } else if (action === "delete") {
    openDeleteModal(user);
  }
}

function bindRealtimeValidation() {
  for (const name of fields) {
    const input = el(name);
    if (!input) continue;

    const handler = () => {
      const error = validateField(name, input.value);
      if (error) setFieldError(name, error);
      else clearFieldError(name);
    };

    input.addEventListener("input", handler);
    input.addEventListener("blur", handler);
  }
}

async function initWasm() {
  if (!("serviceWorker" in navigator)) {
    throw new Error("Service Worker is not available in this browser");
  }

  try {
    const registration = await navigator.serviceWorker.register("/sw.js?v=2", { scope: "/" });
    await registration.update();
    await navigator.serviceWorker.ready;
    wasmReady = true;
  } catch (err) {
    throw err;
  }
}

async function boot() {
  applyTheme(getTheme());
  bindRealtimeValidation();

  el("themeToggle").addEventListener("click", toggleTheme);
  el("newFromListBtn").addEventListener("click", startCreateUser);
  el("menuToggleBtn").addEventListener("click", () => {
    if (el("sideMenu").classList.contains("is-open")) closeSideMenu();
    else openSideMenu();
  });
  el("sideMenuBackdrop").addEventListener("click", (e) => {
    if (e.target.dataset.closeSideMenu === "true") closeSideMenu();
  });

  el("userForm").addEventListener("submit", onSubmitForm);
  el("refreshBtn").addEventListener("click", async () => {
    try { await refreshUsers(); showToast("List updated"); }
    catch (err) { showToast(err.message || "Could not refresh list"); }
  });

  el("searchInput").addEventListener("input", filterUsers);
  el("usersList").addEventListener("click", onUsersListClick);
  el("closeModalBtn").addEventListener("click", closeModal);
  el("modal").addEventListener("click", (e) => {
    if (e.target.dataset.closeModal === "true") closeModal();
  });
  el("closeFormModalBtn").addEventListener("click", closeFormModal);
  el("formModal").addEventListener("click", (e) => {
    if (e.target.dataset.closeFormModal === "true") closeFormModal();
  });

  el("closeDeleteModalBtn").addEventListener("click", closeDeleteModal);
  el("cancelDeleteBtn").addEventListener("click", closeDeleteModal);
  el("deleteModal").addEventListener("click", (e) => {
    if (e.target.dataset.closeDeleteModal === "true") closeDeleteModal();
  });
  el("confirmDeleteBtn").addEventListener("click", async () => {
    if (!pendingDeleteId) return;
    const id = pendingDeleteId;
    closeDeleteModal();
    try {
      await deleteUser(id, true);
      showToast("User deleted");
    } catch (err) {
      showToast(err.message || "Could not delete user");
    }
  });

  try {
    await initWasm();
    showToast("WASM loaded successfully");
  } catch (err) {
    showToast(`WASM initialization error: ${err.message}`);
    return;
  }

  try {
    await refreshUsers();
  } catch (err) {
    renderUsers([]);
  }
}

boot();
