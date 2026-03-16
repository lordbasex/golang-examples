/* global registerWasmHTTPListener */
importScripts("/wasm_exec.js");
importScripts("/db.js");
importScripts("https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@v1.1.0/sw.js");

// IMPORTANT: limit interception to /api/*
// This prevents the SW from handling "/" navigation requests.
registerWasmHTTPListener("/server.wasm", { base: "api" });

addEventListener("install", (event) => {
  event.waitUntil(self.skipWaiting());
});

addEventListener("activate", (event) => {
  event.waitUntil(self.clients.claim());
});
