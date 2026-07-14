<template>
  <div class="tokenhub-public-shell min-h-screen overflow-x-clip bg-[var(--background)] text-[color:var(--foreground)]">
    <header
      :class="[
        'z-50 w-full bg-transparent px-4 md:px-10',
        floatingNav ? 'fixed top-0 pt-5' : 'sticky top-0 py-4'
      ]"
    >
      <nav class="tokenhub-nav-shell relative mx-auto flex max-w-[1280px] items-center justify-between gap-3 rounded-[1.7rem] py-1.5 pl-1.5 pr-4">
        <router-link
          to="/"
          class="inline-flex min-w-0 items-center gap-2 rounded-full px-2.5 py-1.5 transition-opacity hover:opacity-85"
          :aria-label="brandName"
        >
          <span class="brand-icon flex size-9 shrink-0 items-center justify-center overflow-hidden rounded-full">
            <img :src="brandLogo" :alt="brandName" class="h-full w-full object-contain" />
          </span>
          <span class="brand-wordmark hidden min-w-0 truncate text-[15px] font-semibold tracking-normal min-[420px]:inline">
            {{ brandName }}
          </span>
        </router-link>

        <div class="absolute left-1/2 top-1/2 z-[80] hidden -translate-x-1/2 -translate-y-1/2 items-center gap-1 lg:flex">
          <router-link
            v-for="item in navItems"
            :key="item.to"
            :to="item.to"
            class="tokenhub-nav-link inline-flex h-9 items-center rounded-full px-3 text-[14px] transition-all duration-700 ease-[cubic-bezier(0.32,0.72,0,1)] sm:px-3.5"
            :class="isActive(item.to) ? 'is-active' : ''"
          >
            <span>{{ item.label }}</span>
          </router-link>
        </div>

        <div class="relative z-[80] ml-auto flex shrink-0 items-center gap-1">
          <div class="hidden sm:block">
            <LocaleSwitcher variant="circle" />
          </div>
          <button
            type="button"
            class="tokenhub-icon-button inline-flex size-9 items-center justify-center rounded-full transition-all duration-300 ease-out"
            :title="isDark ? copy.lightMode : copy.darkMode"
            @click="toggleTheme"
          >
            <Icon v-if="isDark" name="sun" size="sm" />
            <Icon v-else name="moon" size="sm" />
          </button>
          <router-link
            :to="consolePath"
            class="nav-primary-cta inline-flex h-9 items-center rounded-full px-4 text-sm font-medium"
          >
            {{ copy.start }}
          </router-link>
        </div>
      </nav>

      <div class="mx-auto mt-2 max-w-[1280px] lg:hidden">
        <div class="tokenhub-mobile-nav scrollbar-hide flex gap-1 overflow-x-auto rounded-full p-1">
          <router-link
            v-for="item in navItems"
            :key="item.to"
            :to="item.to"
            class="tokenhub-nav-link shrink-0 rounded-full px-3 py-1.5 text-sm transition-all duration-300"
            :class="isActive(item.to) ? 'is-active' : ''"
          >
            {{ item.label }}
          </router-link>
        </div>
      </div>
    </header>

    <main>
      <slot />
    </main>

    <footer class="footer-shell border-t border-[color:var(--surface-line)]">
      <div class="mx-auto grid max-w-[1280px] gap-8 px-6 py-12 md:grid-cols-[1.2fr_0.8fr_0.8fr] md:px-10">
        <div>
          <div class="mb-4 flex items-center gap-3">
            <span class="brand-icon flex size-9 items-center justify-center overflow-hidden rounded-full">
              <img :src="brandLogo" :alt="brandName" class="h-full w-full object-contain" />
            </span>
            <span class="home-title font-semibold">{{ brandName }}</span>
          </div>
          <p class="home-muted max-w-md text-sm leading-6">
            {{ copy.footerIntro }}
          </p>
        </div>
        <div>
          <h2 class="home-title mb-3 text-sm font-semibold">{{ copy.product }}</h2>
          <div class="grid gap-2 text-sm">
            <router-link to="/" class="home-muted transition-colors hover:text-[color:var(--home-heading)]">{{ copy.home }}</router-link>
            <router-link to="/models" class="home-muted transition-colors hover:text-[color:var(--home-heading)]">{{ copy.models }}</router-link>
            <router-link to="/docs" class="home-muted transition-colors hover:text-[color:var(--home-heading)]">{{ copy.docs }}</router-link>
          </div>
        </div>
        <div>
          <h2 class="home-title mb-3 text-sm font-semibold">{{ copy.endpoint }}</h2>
          <code class="home-code block overflow-x-auto rounded-lg px-3 py-2 text-xs">
            {{ brandEndpoint }}
          </code>
          <p class="home-faint mt-3 text-xs leading-5">
            © 2026 {{ brandName }}. {{ brandDomain }}
          </p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'
import {
  activeBrand,
  getCurrentApiBaseUrl,
  getCurrentHost,
  sharedBrandContent,
  withBrandTokens,
} from '@/brand'

const props = defineProps<{
  pageTitle?: string
  pageDescription?: string
  canonicalPath?: string
  ogImage?: string
}>()

const route = useRoute()
const { locale } = useI18n()
const authStore = useAuthStore()
const isDark = ref(document.documentElement.classList.contains('dark'))
const brandName = computed(() => activeBrand.siteName)
const brandLogo = computed(() => activeBrand.logo)
const brandDomain = computed(() => getCurrentHost())
const brandEndpoint = computed(() => getCurrentApiBaseUrl())
const copy = computed(() => withBrandTokens(
  locale.value.startsWith('zh') ? sharedBrandContent.publicLayout.zh : sharedBrandContent.publicLayout.en,
))
const floatingNav = computed(() => route.name === 'Home')
const navItems = computed(() => [
  { to: '/', label: copy.value.home },
  { to: '/models', label: copy.value.models },
  { to: '/docs', label: copy.value.docs },
])
const consolePath = computed(() => {
  if (!authStore.isAuthenticated) return '/login'
  return authStore.isAdmin ? '/admin/dashboard' : '/dashboard'
})

function isActive(path: string) {
  if (path === '/') return route.name === 'Home'
  return route.path === path
}

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

function syncDocumentTitle() {
  const title = props.pageTitle?.trim()
  document.title = title ? `${title} - ${brandName.value}` : brandName.value
  syncSeoMeta()
}

function absoluteBrandUrl(path: string): string {
  if (/^https?:\/\//i.test(path)) return path
  const base = `${window.location.origin}/`
  return new URL(path.replace(/^\//, ''), base).toString()
}

function upsertMeta(selector: string, attrs: Record<string, string>) {
  let element = document.head.querySelector<HTMLMetaElement>(selector)
  if (!element) {
    element = document.createElement('meta')
    document.head.appendChild(element)
  }
  Object.entries(attrs).forEach(([key, value]) => element?.setAttribute(key, value))
}

function upsertCanonical(href: string) {
  let element = document.head.querySelector<HTMLLinkElement>('link[rel="canonical"]')
  if (!element) {
    element = document.createElement('link')
    element.rel = 'canonical'
    document.head.appendChild(element)
  }
  element.href = href
}

function syncSeoMeta() {
  const description = props.pageDescription?.trim()
  const canonicalUrl = absoluteBrandUrl(props.canonicalPath || route.path || '/')
  const imageUrl = props.ogImage ? absoluteBrandUrl(props.ogImage) : absoluteBrandUrl(activeBrand.logo)

  upsertCanonical(canonicalUrl)
  if (!description) return

  upsertMeta('meta[name="description"]', { name: 'description', content: description })
  upsertMeta('meta[property="og:type"]', { property: 'og:type', content: 'website' })
  upsertMeta('meta[property="og:title"]', { property: 'og:title', content: document.title })
  upsertMeta('meta[property="og:description"]', { property: 'og:description', content: description })
  upsertMeta('meta[property="og:url"]', { property: 'og:url', content: canonicalUrl })
  upsertMeta('meta[property="og:image"]', { property: 'og:image', content: imageUrl })
  upsertMeta('meta[name="twitter:card"]', { name: 'twitter:card', content: 'summary_large_image' })
  upsertMeta('meta[name="twitter:title"]', { name: 'twitter:title', content: document.title })
  upsertMeta('meta[name="twitter:description"]', { name: 'twitter:description', content: description })
  upsertMeta('meta[name="twitter:image"]', { name: 'twitter:image', content: imageUrl })
}

watch([() => props.pageTitle, () => props.pageDescription, () => props.canonicalPath, () => props.ogImage, () => locale.value, () => route.fullPath, () => brandName.value], () => nextTick(syncDocumentTitle), {
  immediate: true
})

onMounted(() => {
  syncDocumentTitle()
})
</script>

<style scoped>
.tokenhub-public-shell {
  --background: #f8fafc;
  --foreground: #0f172a;
  --primary: rgb(var(--color-primary-600));
  --primary-foreground: #ffffff;
  --surface-shell: transparent;
  --surface-core: #ffffffd9;
  --surface-line: rgb(var(--color-primary-700) / 0.12);
  --surface-highlight: #ffffff;
  --surface-control: #ffffffd9;
  --surface-control-hover: rgb(var(--color-primary-50));
  --surface-control-shadow: 0 10px 28px rgb(var(--color-primary-600) / 0.09);
  --home-heading: #0f172a;
  --home-muted: #64748b;
  --home-faint: #94a3b8;
  --home-code: #0f172a;
  --brand-ring: rgb(var(--color-primary-600) / 0.22);
  --nav-muted: #334155cc;
  --nav-strong: #0f172a;
  --nav-shell: #ffffffd9;
  --nav-shadow: 0 18px 54px rgb(var(--color-primary-600) / 0.09);
}

.dark .tokenhub-public-shell {
  --background: #07111f;
  --foreground: #f8fbff;
  --primary: rgb(var(--color-primary-300));
  --primary-foreground: #07111f;
  --surface-shell: transparent;
  --surface-core: #0f1f35d9;
  --surface-line: rgb(var(--color-primary-300) / 0.15);
  --surface-highlight: #172a46;
  --surface-control: #10223ab8;
  --surface-control-hover: #19365c;
  --surface-control-shadow: 0 18px 58px #0000002e;
  --home-heading: #f8fbff;
  --home-muted: #cbd5e1;
  --home-faint: #94a3b8;
  --home-code: #061121;
  --brand-ring: rgb(var(--color-primary-300) / 0.24);
  --nav-muted: #dbeafecc;
  --nav-strong: #f8fbff;
  --nav-shell: #0f1f35d9;
  --nav-shadow: 0 18px 58px #00000038;
}

.tokenhub-nav-shell {
  background: var(--nav-shell);
  box-shadow: var(--nav-shadow);
  backdrop-filter: blur(26px) saturate(1.05);
}

.tokenhub-mobile-nav {
  background: var(--nav-shell);
  box-shadow: var(--nav-shadow);
  backdrop-filter: blur(26px) saturate(1.05);
}

.brand-icon {
  background: #ffffff;
  box-shadow: inset 0 0 0 1px rgb(255 255 255 / 0.12);
}

.brand-wordmark,
.home-title {
  color: var(--home-heading);
}

.home-muted {
  color: var(--home-muted);
}

.home-faint {
  color: var(--home-faint);
}

.home-code {
  background: var(--home-code);
  color: color-mix(in oklab, var(--home-heading) 72%, transparent);
}

.tokenhub-nav-link {
  color: var(--nav-muted);
}

.tokenhub-nav-link:hover,
.tokenhub-nav-link.is-active {
  background: var(--surface-control-hover);
  color: var(--nav-strong);
  box-shadow: var(--surface-control-shadow);
}

.tokenhub-icon-button {
  background: var(--surface-control);
  color: var(--nav-muted);
  box-shadow: var(--surface-control-shadow);
}

.tokenhub-icon-button:hover {
  background: var(--surface-control-hover);
  color: var(--nav-strong);
}

.nav-primary-cta {
  background: var(--primary);
  color: var(--primary-foreground);
  box-shadow: 0 10px 28px var(--brand-ring), inset 0 1px 0 rgb(255 255 255 / 0.2);
  transition: background-color 0.18s ease-out, box-shadow 0.18s ease-out, transform 0.18s ease-out;
}

.nav-primary-cta:hover {
  filter: brightness(1.06);
  box-shadow: 0 13px 32px var(--brand-ring), inset 0 1px 0 rgb(255 255 255 / 0.24);
  transform: translateY(-1px);
}

.nav-primary-cta:active {
  transform: translateY(0) scale(0.98);
}

.footer-shell {
  background:
    linear-gradient(180deg, var(--surface-highlight), transparent),
    var(--surface-core);
  backdrop-filter: blur(26px) saturate(1.05);
  box-shadow: 0 24px 90px #15172013;
}
</style>
