const http = require('http');
const httpProxy = require('http-proxy-middleware');
const express = require('express');
const path = require('path');

const app = express();

// Middleware to disable caching
app.use((req, res, next) => {
  res.setHeader('Cache-Control', 'no-cache, no-store, must-revalidate');
  res.setHeader('Pragma', 'no-cache');
  res.setHeader('Expires', '0');
  next();
});

// Statische Dateien servieren (unsere HTML-Frontends)
app.use(express.static(__dirname));

// Proxy für Auth-Service
app.use('/api/auth', httpProxy.createProxyMiddleware({
  target: 'http://localhost:8081',
  changeOrigin: true,
  pathRewrite: {
    '^/api/auth': '', // entfernt /api/auth vom Path
  },
}));

// Proxy für Shopping-Service
app.use('/api/shopping', httpProxy.createProxyMiddleware({
  target: 'http://localhost:8080',
  changeOrigin: true,
  pathRewrite: {
    '^/api/shopping': '', // entfernt /api/shopping vom Path
  },
}));

// Proxy für Checkout-Service (falls vorhanden)
app.use('/api/checkout', httpProxy.createProxyMiddleware({
  target: 'http://localhost:8082',
  changeOrigin: true,
  pathRewrite: {
    '^/api/checkout': '',
  },
}));

const PORT = 3000;
app.listen(PORT, () => {
  console.log(`🌐 Proxy Server läuft auf http://localhost:${PORT}`);
  console.log(`📄 Frontend URLs:`);
  console.log(`   Main Dashboard: http://localhost:${PORT}/index.html`);
  console.log(`   Auth Service: http://localhost:${PORT}/auth-service/index.html`);
  console.log(`   Shopping Service: http://localhost:${PORT}/shopping-service/index.html`);
  console.log(`   Checkout Service: http://localhost:${PORT}/checkout-service/index.html`);
});
