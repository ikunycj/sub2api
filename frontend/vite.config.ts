import { defineConfig, loadEnv, Plugin } from 'vite'
import vue from '@vitejs/plugin-vue'
import checker from 'vite-plugin-checker'
import { resolve } from 'path'
import { DEFAULT_BRAND_ENV, PRIMARY_COLOR_SCHEMES } from './src/brand/palettes'
import type { BrandColorScale, BrandPrimaryScheme } from './src/brand/palettes'

type ViteSiteConfig = {
  siteName: string
  logo: string
  primaryScheme: BrandPrimaryScheme
  primary: BrandColorScale
}

function resolveViteSiteConfig(env: Record<string, string>): ViteSiteConfig {
  const siteName = env.VITE_SITE_NAME?.trim() || DEFAULT_BRAND_ENV.siteName
  const logo = env.VITE_SITE_LOGO?.trim() || DEFAULT_BRAND_ENV.logo
  const primaryScheme = (env.VITE_PRIMARY_SCHEME?.trim().toLowerCase() || DEFAULT_BRAND_ENV.primaryScheme) as BrandPrimaryScheme
  const primary = PRIMARY_COLOR_SCHEMES[primaryScheme]

  if (!primary) {
    throw new Error(
      `Unsupported VITE_PRIMARY_SCHEME "${env.VITE_PRIMARY_SCHEME}". Supported schemes: ${Object.keys(PRIMARY_COLOR_SCHEMES).join(', ')}`,
    )
  }

  return { siteName, logo, primaryScheme, primary }
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

function getPrimaryCssText(primary: BrandColorScale): string {
  return Object.entries(primary)
    .map(([step, color]) => `--color-primary-${step}: ${hexToRgbChannels(color)};`)
    .join(' ')
}

function iconMimeType(iconUrl: string): string {
  const cleanUrl = iconUrl.split('?')[0]?.toLowerCase() ?? ''
  if (cleanUrl.endsWith('.svg')) return 'image/svg+xml'
  if (cleanUrl.endsWith('.png')) return 'image/png'
  return 'image/x-icon'
}

function escapeHtml(value: string): string {
  return value
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}

function injectSiteDefaults(site: ViteSiteConfig): Plugin {
  return {
    name: 'inject-site-defaults',
    transformIndexHtml: {
      order: 'pre',
      handler(html) {
        const siteName = escapeHtml(site.siteName)
        const logo = escapeHtml(site.logo)
        const description = escapeHtml(`${site.siteName} provides OpenAI-compatible access to OpenAI and Claude models, with selected models priced from as low as 1/50 of official API rates.`)
        const iconTag = `<link rel="icon" type="${iconMimeType(site.logo)}" href="${logo}" />`
        const title = `<title>${siteName}</title>`
        const themeStyle = `<style id="site-theme">:root{${getPrimaryCssText(site.primary)}}</style>`
        const seoTags = [
          '<meta name="robots" content="index,follow" />',
          `<meta name="description" content="${description}" />`,
          '<meta property="og:type" content="website" />',
          `<meta property="og:title" content="${siteName}" />`,
          `<meta property="og:description" content="${description}" />`,
          `<meta property="og:image" content="${logo}" />`,
          '<meta name="twitter:card" content="summary" />',
          `<meta name="twitter:title" content="${siteName}" />`,
          `<meta name="twitter:description" content="${description}" />`,
          `<meta name="twitter:image" content="${logo}" />`,
        ].join('\n        ')
        const withScheme = /<html\b[^>]*\bdata-primary-scheme=/i.test(html)
          ? html.replace(/(<html\b[^>]*\bdata-primary-scheme=)["'][^"']*["']/i, `$1"${site.primaryScheme}"`)
          : html.replace(/<html\b([^>]*)>/i, `<html$1 data-primary-scheme="${site.primaryScheme}">`)
        const withIcon = withScheme.match(/<link\s+rel="icon"[^>]*>/i)
          ? withScheme.replace(/<link\s+rel="icon"[^>]*>/i, iconTag)
          : withScheme.replace('</head>', `${iconTag}\n</head>`)
        const withTitle = withIcon.match(/<title>.*?<\/title>/i)
          ? withIcon.replace(/<title>.*?<\/title>/i, title)
          : withIcon.replace('</head>', `${title}\n</head>`)
        const withTheme = withTitle.includes('id="site-theme"')
          ? withTitle
          : withTitle.replace('</head>', `${themeStyle}\n</head>`)
        return withTheme.includes('name="description"')
          ? withTheme
          : withTheme.replace('</head>', `        ${seoTags}\n</head>`)
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
  const site = resolveViteSiteConfig(env)

  return {
    plugins: [
      injectSiteDefaults(site),
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
