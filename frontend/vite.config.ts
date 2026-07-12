import { defineConfig, loadEnv, Plugin } from 'vite'
import vue from '@vitejs/plugin-vue'
import checker from 'vite-plugin-checker'
import { resolve } from 'path'

type ViteBrandKey = 'ikun' | 'anytoken'
type ViteBrandColorStep = '50' | '100' | '200' | '300' | '400' | '500' | '600' | '700' | '800' | '900' | '950'
type ViteBrandConfig = {
  siteName: string
  favicon: string
  primary: Record<ViteBrandColorStep, string>
}

const viteBrandConfigs = {
  ikun: {
    siteName: 'ikun.love',
    favicon: '/brands/ikun/logo.svg',
    primary: {
      50: '#fff8ed',
      100: '#ffefd0',
      200: '#ffdaa1',
      300: '#ffbe6b',
      400: '#ff9832',
      500: '#f97316',
      600: '#ea580c',
      700: '#c2410c',
      800: '#9a3412',
      900: '#7c2d12',
      950: '#431407',
    },
  },
  anytoken: {
    siteName: 'AnyToken',
    favicon: '/brands/anytoken/logo.svg',
    primary: {
      50: '#eff6ff',
      100: '#dbeafe',
      200: '#bfdbfe',
      300: '#93c5fd',
      400: '#60a5fa',
      500: '#2563eb',
      600: '#1d4ed8',
      700: '#1e40af',
      800: '#1e3a8a',
      900: '#172554',
      950: '#0f172a',
    },
  },
} satisfies Record<ViteBrandKey, ViteBrandConfig>

function resolveViteBrandConfig(value?: string): ViteBrandConfig {
  const normalized = (value || 'ikun').trim().toLowerCase()
  if (normalized in viteBrandConfigs) {
    return viteBrandConfigs[normalized as ViteBrandKey]
  }
  throw new Error(`Unsupported VITE_BRAND "${value}". Supported brands: ${Object.keys(viteBrandConfigs).join(', ')}`)
}

function hexToRgbChannels(hex: string): string {
  const normalized = hex.replace('#', '').trim()
  const value = normalized.length === 3
    ? normalized.split('').map((char) => `${char}${char}`).join('')
    : normalized
  return [
    Number.parseInt(value.slice(0, 2), 16),
    Number.parseInt(value.slice(2, 4), 16),
    Number.parseInt(value.slice(4, 6), 16),
  ].join(' ')
}

function getViteBrandCssText(brand: ViteBrandConfig): string {
  return Object.entries(brand.primary)
    .map(([step, color]) => `--color-primary-${step}: ${hexToRgbChannels(color)};`)
    .join(' ')
}

function iconMimeType(iconUrl: string): string {
  const cleanUrl = iconUrl.split('?')[0]?.toLowerCase() ?? ''
  if (cleanUrl.endsWith('.svg')) return 'image/svg+xml'
  if (cleanUrl.endsWith('.png')) return 'image/png'
  return 'image/x-icon'
}

function injectBrandDefaults(brand: ViteBrandConfig): Plugin {
  return {
    name: 'inject-brand-defaults',
    transformIndexHtml: {
      order: 'pre',
      handler(html) {
        const iconTag = `<link rel="icon" type="${iconMimeType(brand.favicon)}" href="${brand.favicon}" />`
        const title = `<title>${brand.siteName}</title>`
        const themeStyle = `<style id="brand-theme">:root{${getViteBrandCssText(brand)}}</style>`
        const withIcon = html.match(/<link\s+rel="icon"[^>]*>/i)
          ? html.replace(/<link\s+rel="icon"[^>]*>/i, iconTag)
          : html.replace('</head>', `${iconTag}\n</head>`)
        const withTitle = withIcon.match(/<title>.*?<\/title>/i)
          ? withIcon.replace(/<title>.*?<\/title>/i, title)
          : withIcon.replace('</head>', `${title}\n</head>`)
        return withTitle.includes('id="brand-theme"')
          ? withTitle
          : withTitle.replace('</head>', `${themeStyle}\n</head>`)
      }
    }
  }
}

/**
 * Vite 插件：开发模式下注入公开配置到 index.html
 * 与生产模式的后端注入行为保持一致，消除闪烁
 */
function injectPublicSettings(backendUrl: string): Plugin {
  return {
    name: 'inject-public-settings',
    apply: 'serve',
    transformIndexHtml: {
      order: 'pre',
      async handler(html) {
        try {
          const response = await fetch(`${backendUrl}/api/v1/settings/public`, {
            signal: AbortSignal.timeout(2000)
          })
          if (response.ok) {
            const data = await response.json()
            if (data.code === 0 && data.data) {
              const script = `<script>window.__APP_CONFIG__=${JSON.stringify(data.data)};</script>`
              return html.replace('</head>', `${script}\n</head>`)
            }
          }
        } catch (e) {
          console.warn('[vite] 无法获取公开配置，将回退到 API 调用:', (e as Error).message)
        }
        return html
      }
    }
  }
}

export default defineConfig(({ mode }) => {
  // 加载环境变量
  const env = loadEnv(mode, process.cwd(), '')
  const backendUrl = env.VITE_DEV_PROXY_TARGET || 'http://localhost:8080'
  const devPort = Number(env.VITE_DEV_PORT || 3000)
  const brand = resolveViteBrandConfig(env.VITE_BRAND)

  return {
    plugins: [
      injectBrandDefaults(brand),
      vue(),
      checker({
        vueTsc: true
      }),
      injectPublicSettings(backendUrl)
    ],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      // 使用 vue-i18n 运行时版本，避免 CSP unsafe-eval 问题
      'vue-i18n': 'vue-i18n/dist/vue-i18n.runtime.esm-bundler.js'
    }
  },
  define: {
    // 启用 vue-i18n JIT 编译，在 CSP 环境下处理消息插值
    // JIT 编译器生成 AST 对象而非 JS 代码，无需 unsafe-eval
    __INTLIFY_JIT_COMPILATION__: true
  },
  build: {
    outDir: '../backend/internal/web/dist',
    emptyOutDir: true,
    rollupOptions: {
      output: {
        /**
         * 手动分包配置
         * 分离第三方库并按功能合并应用代码，避免循环依赖
         */
        manualChunks(id: string) {
          if (id.includes('node_modules')) {
            // Vue 核心库
            if (
              id.includes('/vue/') ||
              id.includes('/vue-router/') ||
              id.includes('/pinia/') ||
              id.includes('/@vue/')
            ) {
              return 'vendor-vue'
            }

            // UI 工具库（较大，单独分离）
            if (id.includes('/@vueuse/') || id.includes('/xlsx/')) {
              return 'vendor-ui'
            }

            // 图表库
            if (id.includes('/chart.js/') || id.includes('/vue-chartjs/')) {
              return 'vendor-chart'
            }

            // 国际化
            if (id.includes('/vue-i18n/') || id.includes('/@intlify/')) {
              return 'vendor-i18n'
            }

            // 其他小型第三方库合并
            return 'vendor-misc'
          }

          // 应用代码：按入口点自动分包，不手动干预
          // 这样可以避免循环依赖，同时保持合理的 chunk 数量
        }
      }
    }
  },
    server: {
      host: '0.0.0.0',
      port: devPort,
      proxy: {
        '/api': {
          target: backendUrl,
          changeOrigin: true
        },
        '/v1': {
          target: backendUrl,
          changeOrigin: true
        },
        '/setup': {
          target: backendUrl,
          changeOrigin: true
        }
      }
    }
  }
})
