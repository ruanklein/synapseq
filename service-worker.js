const CACHE_VERSION = "synapseq-pwa-v1.0.1";

const APP_SHELL = [
  "/",
  "/index.html",
  "/script.js",
  "/styles.css",
  "/wasm/synapseq.wasm",
  "/wasm/wasm_exec.js",
  "/wasm/synapseq.js",

  // icons
  "/assets/icon-192.png",
  "/assets/icon-256.png",
  "/assets/icon-384.png",
  "/assets/icon-512.png",
  "/assets/icon-maskable.png",
];

self.addEventListener("install", (event) => {
  event.waitUntil(
    caches.open(CACHE_VERSION).then((cache) => {
      return cache.addAll(APP_SHELL);
    })
  );
  self.skipWaiting();
});

self.addEventListener("activate", (event) => {
  event.waitUntil(
    caches.keys().then((keys) => {
      return Promise.all(
        keys.filter((k) => k !== CACHE_VERSION).map((k) => caches.delete(k))
      );
    })
  );
  self.clients.claim();
});

self.addEventListener("fetch", (event) => {
  const req = event.request;

  if (req.mode === "navigate") {
    event.respondWith(fetch(req).catch(() => caches.match("/index.html")));
    return;
  }

  event.respondWith(
    caches.match(req).then((cached) => {
      if (cached) return cached;

      return fetch(req).then((res) => {
        if (req.method === "GET") {
          const clone = res.clone();
          caches.open(CACHE_VERSION).then((cache) => {
            cache.put(req, clone);
          });
        }
        return res;
      });
    })
  );
});
