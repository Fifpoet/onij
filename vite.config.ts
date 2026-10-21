import path from 'node:path'
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import unocss from 'unocss/vite'
import viteCompression from 'vite-plugin-compression'
import { visualizer } from 'rollup-plugin-visualizer'

export default defineConfig((configEnv) => {
  const env = loadEnv(configEnv.mode, process.cwd())

  return {
    base: env.VITE_PUBLIC_PATH || '/',
    resolve: {
      alias: {
        '@': path.resolve(path.resolve(process.cwd()), 'src'),
        '~': path.resolve(process.cwd()),
      },
    },
    plugins: [
      vue(),
      unocss(),
      viteCompression({ algorithm: 'gzip' }),
      visualizer({ open: false, gzipSize: true, brotliSize: true }),
    ],
    server: {
      host: '0.0.0.0',
      port: 18968,
      open: false,
      proxy: {
        // 避免 HSTS 把 http://onij.fun:3000 升成 https://onij.fun:3000（无证书）
        '/netease': {
          target: 'http://onij.fun:3000',
          changeOrigin: true,
          rewrite: (p) => p.replace(/^\/netease/, ''),
        },
        // 取链走播放器同源 api.php（music.gdstudio.org），带 crc32 签名
        '/gdstudio-music': {
          target: 'https://music.gdstudio.org',
          changeOrigin: true,
          rewrite: (p) => p.replace(/^\/gdstudio-music/, ''),
          configure: (proxy) => {
            proxy.on('proxyReq', (proxyReq) => {
              proxyReq.setHeader('Origin', 'https://music.gdstudio.org')
              proxyReq.setHeader('Referer', 'https://music.gdstudio.org/')
            })
          },
        },
        // 主站 Go：页面同源 /presence 等，避免浏览器直连 127.0.0.1:8889
        '^/(ping|album|artist|file|music|practice|tag|uvr|voice|ai|collection|tran|presence|api)(/|$)': {
          target: 'http://127.0.0.1:8889',
          changeOrigin: true,
        },
      },
    },
    // https://cn.vitejs.dev/guide/api-javascript.html#build
    build: {
      chunkSizeWarningLimit: 1024, // chunk 大小警告的限制 (单位 kb)
    },
    esbuild: {
      drop: ['debugger'], // console
    },
  }
})
