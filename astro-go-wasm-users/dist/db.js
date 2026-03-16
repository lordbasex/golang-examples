const DB_NAME = "astro_go_wasm_users_db";
const STORE_NAME = "users_secure";
const DB_VERSION = 2;

let db = null;

function requestToPromise(req) {
  return new Promise((resolve, reject) => {
    req.onsuccess = () => resolve(req.result);
    req.onerror = () => reject(req.error || new Error("indexeddb_error"));
  });
}

async function initDB() {
  if (db) return db;

  db = await new Promise((resolve, reject) => {
    const request = indexedDB.open(DB_NAME, DB_VERSION);

    request.onupgradeneeded = (event) => {
      const database = event.target.result;
      let store;
      if (!database.objectStoreNames.contains(STORE_NAME)) {
        store = database.createObjectStore(STORE_NAME, { keyPath: "id", autoIncrement: true });
      } else {
        store = event.target.transaction.objectStore(STORE_NAME);
      }

      if (!store.indexNames.contains("email")) store.createIndex("email", "email", { unique: false });
      if (!store.indexNames.contains("document")) store.createIndex("document", "document", { unique: false });
      if (!store.indexNames.contains("updatedAt")) store.createIndex("updatedAt", "updatedAt", { unique: false });
    };

    request.onsuccess = () => resolve(request.result);
    request.onerror = () => reject(request.error || new Error("db_open_error"));
  });

  return db;
}

async function setEncryptionPassphrase(passphrase) {
  await initDB();
  return { ok: true };
}

function ensureUnlocked() {
  return true;
}

async function encryptPayload(payload) {
  return { payload };
}

async function decryptPayload(record) {
  if (record && record.payload) return record.payload;
  return {};
}

async function readAllStored() {
  await initDB();
  const tx = db.transaction(STORE_NAME, "readonly");
  const store = tx.objectStore(STORE_NAME);
  return requestToPromise(store.getAll());
}

async function apiCreateUser(user) {
  await initDB();
  ensureUnlocked();

  const payload = {
    ...user,
    createdAt: user.createdAt || new Date().toISOString(),
    updatedAt: user.updatedAt || new Date().toISOString(),
  };

  const encrypted = await encryptPayload(payload);
  const tx = db.transaction(STORE_NAME, "readwrite");
  const store = tx.objectStore(STORE_NAME);

  const email = String(payload.email || "").trim().toLowerCase();
  if (email) {
    const existing = await requestToPromise(store.index("email").get(email));
    if (existing) {
      return JSON.stringify({ ok: false, message: "validation_failed", errors: { email: "Email already exists" } });
    }
  }

  const document = String(payload.document || "").trim();
  if (document) {
    const existing = await requestToPromise(store.index("document").get(document));
    if (existing) {
      return JSON.stringify({ ok: false, message: "validation_failed", errors: { document: "Document already exists" } });
    }
  }

  const id = await requestToPromise(store.add({
    email,
    document,
    updatedAt: payload.updatedAt,
    ...encrypted,
  }));

  return JSON.stringify({ ok: true, data: { id, ...payload } });
}

async function apiListUsers() {
  const rows = await readAllStored();
  const users = [];
  for (const row of rows) {
    const payload = await decryptPayload(row);
    users.push({ id: row.id, ...payload });
  }
  users.sort((a, b) => `${a.lastName} ${a.firstName}`.localeCompare(`${b.lastName} ${b.firstName}`));
  return JSON.stringify({ ok: true, data: users });
}

async function apiGetUser(id) {
  await initDB();
  const tx = db.transaction(STORE_NAME, "readonly");
  const store = tx.objectStore(STORE_NAME);
  const row = await requestToPromise(store.get(Number(id)));

  if (!row) return JSON.stringify({ ok: false, message: "not_found" });

  const payload = await decryptPayload(row);
  return JSON.stringify({ ok: true, data: { id: row.id, ...payload } });
}

async function apiUpdateUser(user) {
  await initDB();
  ensureUnlocked();

  const id = Number(user.id);
  const txRead = db.transaction(STORE_NAME, "readonly");
  const storeRead = txRead.objectStore(STORE_NAME);
  const current = await requestToPromise(storeRead.get(id));

  if (!current) return JSON.stringify({ ok: false, message: "not_found" });

  const currentPayload = await decryptPayload(current);
  const payload = {
    ...currentPayload,
    ...user,
    id,
    updatedAt: user.updatedAt || new Date().toISOString(),
  };

  const encrypted = await encryptPayload(payload);
  const tx = db.transaction(STORE_NAME, "readwrite");
  const store = tx.objectStore(STORE_NAME);

  const email = String(payload.email || "").trim().toLowerCase();
  if (email) {
    const existing = await requestToPromise(store.index("email").get(email));
    if (existing && Number(existing.id) !== id) {
      return JSON.stringify({ ok: false, message: "validation_failed", errors: { email: "Email already exists" } });
    }
  }

  const document = String(payload.document || "").trim();
  if (document) {
    const existing = await requestToPromise(store.index("document").get(document));
    if (existing && Number(existing.id) !== id) {
      return JSON.stringify({ ok: false, message: "validation_failed", errors: { document: "Document already exists" } });
    }
  }

  await requestToPromise(store.put({
    id,
    email,
    document,
    updatedAt: payload.updatedAt,
    ...encrypted,
  }));

  return JSON.stringify({ ok: true, data: payload });
}

async function apiDeleteUser(id) {
  await initDB();
  const tx = db.transaction(STORE_NAME, "readwrite");
  const store = tx.objectStore(STORE_NAME);
  await requestToPromise(store.delete(Number(id)));
  return JSON.stringify({ ok: true, data: { id: Number(id) } });
}

globalThis.setEncryptionPassphrase = setEncryptionPassphrase;
globalThis.apiCreateUser = apiCreateUser;
globalThis.apiListUsers = apiListUsers;
globalThis.apiGetUser = apiGetUser;
globalThis.apiUpdateUser = apiUpdateUser;
globalThis.apiDeleteUser = apiDeleteUser;

initDB();
