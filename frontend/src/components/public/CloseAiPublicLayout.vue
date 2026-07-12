<template>
  <div class="min-h-screen overflow-x-clip bg-slate-50 text-slate-950 dark:bg-slate-950 dark:text-white">
    <header class="sticky top-0 z-40 border-b border-slate-200/80 bg-white/95 backdrop-blur dark:border-slate-800 dark:bg-slate-950/92">
      <nav class="mx-auto flex h-16 max-w-7xl items-center justify-between gap-2 px-3 sm:gap-3 sm:px-6 lg:px-8">
        <router-link to="/home" class="group flex min-w-0 items-center gap-2" :aria-label="brandName">
          <span class="flex h-9 w-9 shrink-0 items-center justify-center overflow-hidden rounded-lg bg-white shadow-sm shadow-primary-500/25 ring-1 ring-primary-100 transition-transform duration-200 group-hover:scale-95 dark:bg-slate-900 dark:ring-primary-900/50">
            <img :src="brandLogo" :alt="brandName" class="h-full w-full object-contain" />
          </span>
          <span class="hidden min-w-0 truncate text-sm font-semibold tracking-normal text-slate-950 dark:text-white min-[420px]:inline">{{ brandName }}</span>
        </router-link>

        <div class="hidden items-center gap-1 md:flex">
          <router-link
            v-for="item in navItems"
            :key="item.to"
            :to="item.to"
            class="rounded-full px-3 py-2 text-sm font-medium transition-all duration-200 ease-out"
            :class="isActive(item.to)
              ? 'bg-primary-50 text-primary-700 ring-1 ring-primary-100 dark:bg-primary-950/45 dark:text-primary-200 dark:ring-primary-900/60'
              : 'text-slate-600 hover:bg-slate-100 hover:text-slate-950 dark:text-slate-300 dark:hover:bg-slate-900 dark:hover:text-white'"
          >
            {{ item.label }}
          </router-link>
        </div>

        <div class="flex shrink-0 items-center gap-0.5 sm:gap-1">
          <LocaleSwitcher />
          <button
            type="button"
            class="hidden rounded-full p-2 text-slate-500 transition-all duration-200 hover:bg-slate-100 hover:text-slate-900 active:scale-95 dark:text-slate-300 dark:hover:bg-slate-900 dark:hover:text-white sm:inline-flex"
            :title="isDark ? copy.lightMode : copy.darkMode"
            @click="toggleTheme"
          >
            <Icon v-if="isDark" name="sun" size="md" />
            <Icon v-else name="moon" size="md" />
          </button>
          <router-link
            :to="consolePath"
            :title="copy.start"
            :aria-label="copy.start"
            class="inline-flex h-9 w-9 items-center justify-center rounded-full bg-primary-500 text-sm font-semibold text-white shadow-sm shadow-primary-500/25 transition-all duration-200 ease-out hover:-translate-y-px hover:bg-primary-600 hover:shadow-md hover:shadow-primary-500/30 active:translate-y-0 active:scale-95 focus:outline-none focus:ring-2 focus:ring-primary-500/35 focus:ring-offset-2 dark:focus:ring-offset-slate-950 sm:w-auto sm:px-4"
          >
            <Icon name="login" size="sm" class="sm:hidden" />
            <span class="hidden sm:inline">{{ copy.start }}</span>
          </router-link>
        </div>
      </nav>

      <div class="border-t border-slate-100 bg-white/95 px-3 py-2 dark:border-slate-900 dark:bg-slate-950/95 md:hidden">
        <div class="scrollbar-hide mx-auto flex max-w-7xl gap-1 overflow-x-auto">
          <router-link
            v-for="item in navItems"
            :key="item.to"
            :to="item.to"
            class="shrink-0 rounded-full px-3 py-1.5 text-sm font-medium transition-all duration-200"
            :class="isActive(item.to)
              ? 'bg-primary-50 text-primary-700 ring-1 ring-primary-100 dark:bg-primary-950/45 dark:text-primary-200 dark:ring-primary-900/60'
              : 'text-slate-600 dark:text-slate-300'"
          >
            {{ item.label }}
          </router-link>
        </div>
      </div>
    </header>

    <main>
      <slot />
    </main>

    <footer class="border-t border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-950">
      <div class="mx-auto grid max-w-7xl gap-8 px-4 py-10 sm:px-6 md:grid-cols-[1.2fr_0.8fr_0.8fr] lg:px-8">
        <div>
          <div class="mb-3 flex items-center gap-3">
            <span class="flex h-9 w-9 items-center justify-center overflow-hidden rounded-lg bg-white ring-1 ring-primary-100 dark:bg-slate-900 dark:ring-primary-900/50">
              <img :src="brandLogo" :alt="brandName" class="h-full w-full object-contain" />
            </span>
            <span class="font-semibold text-slate-950 dark:text-white">{{ brandName }}</span>
          </div>
          <p class="max-w-md text-sm leading-6 text-slate-600 dark:text-slate-300">
            {{ copy.footerIntro }}
          </p>
        </div>
        <div>
          <h2 class="mb-3 text-sm font-semibold text-slate-950 dark:text-white">{{ copy.product }}</h2>
          <div class="grid gap-2 text-sm text-slate-600 dark:text-slate-300">
            <router-link to="/home" class="transition-colors hover:text-primary-700 dark:hover:text-primary-300">{{ copy.home }}</router-link>
            <router-link to="/models" class="transition-colors hover:text-primary-700 dark:hover:text-primary-300">{{ copy.models }}</router-link>
            <router-link to="/docs" class="transition-colors hover:text-primary-700 dark:hover:text-primary-300">{{ copy.docs }}</router-link>
          </div>
        </div>
        <div>
          <h2 class="mb-3 text-sm font-semibold text-slate-950 dark:text-white">{{ copy.endpoint }}</h2>
          <code class="block overflow-x-auto rounded-lg border border-slate-200 bg-slate-50 px-3 py-2 text-xs text-slate-700 dark:border-slate-800 dark:bg-slate-900 dark:text-slate-200">
            {{ brandEndpoint }}
          </code>
          <p class="mt-3 text-xs leading-5 text-slate-500 dark:text-slate-400">
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
import { useAppStore, useAuthStore } from '@/stores'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'
import { activeBrand } from '@/brand'
import { sanitizeUrl } from '@/utils/url'

const props = defineProps<{
  pageTitle?: string
}>()

const route = useRoute()
const { locale } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()
const isDark = ref(document.documentElement.classList.contains('dark'))
const brandName = computed(() => appStore.siteName || activeBrand.siteName)
const brandLogo = computed(() => sanitizeUrl(appStore.siteLogo || activeBrand.logo, { allowRelative: true, allowDataUrl: true }) || activeBrand.logo)
const brandDomain = computed(() => activeBrand.domain)
const brandEndpoint = computed(() => activeBrand.apiBaseUrl)
const copy = computed(() => locale.value.startsWith('zh') ? activeBrand.publicLayout.zh : activeBrand.publicLayout.en)
const navItems = computed(() => [
  { to: '/home', label: copy.value.home },
  { to: '/models', label: copy.value.models },
  { to: '/docs', label: copy.value.docs },
])
const consolePath = computed(() => {
  if (!authStore.isAuthenticated) return '/login'
  return authStore.isAdmin ? '/admin/dashboard' : '/dashboard'
})

function isActive(path: string) {
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
}

watch([() => props.pageTitle, () => locale.value, () => route.fullPath, () => brandName.value], () => nextTick(syncDocumentTitle), {
  immediate: true
})

onMounted(() => {
  syncDocumentTitle()
})
</script>
