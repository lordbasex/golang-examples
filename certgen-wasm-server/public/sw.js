/* global registerWasmHTTPListener */
importScripts("/wasm_exec.js");
importScripts("https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@v1.1.0/sw.js");

// This Service Worker bootstraps the Go WASM HTTP handlers and only intercepts "/api/*".
console.log("certgen-wasm-server: sw.js loaded");

// IMPORTANT: limit interception to /api/*
// This prevents the SW from handling "/" navigation requests.
registerWasmHTTPListener("/server.wasm", { base: "api" });

addEventListener("install", (event) => {
  console.log("certgen-wasm-server: install");
  event.waitUntil(self.skipWaiting());
});

addEventListener("activate", (event) => {
  console.log("certgen-wasm-server: activate");
  event.waitUntil(self.clients.claim());
});

