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
        '/file': {
          target: env.VITE_SERVER_URL || 'http://127.0.0.1:8889',
          changeOrigin: true,
        },
        '/api': {
          target: env.VITE_SERVER_URL || 'http://127.0.0.1:8889',
          changeOrigin: true,
        },
        // 备用代理；播放直链默认由浏览器直连 music.gdstudio.org（该站已放行 localhost CORS）
        '/gdstudio-music': {
          target: 'https://music.gdstudio.org',
          changeOrigin: true,
          rewrite: (p) => p.replace(/^\/gdstudio-music/, ''),
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
