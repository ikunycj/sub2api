<template>
  <div v-if="homeContent.trim()" class="min-h-screen">
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <div v-else v-html="customHomeHtml"></div>
  </div>

  <CloseAiPublicLayout
    v-else
    :page-title="page.metaTitle"
    :page-description="page.hero.description"
    canonical-path="/"
  >
    <div class="home-surface relative isolate overflow-hidden text-[color:var(--foreground)]">
      <section class="mx-auto flex min-h-[92svh] w-full max-w-[1280px] flex-col justify-center overflow-hidden px-6 pb-8 pt-28 md:px-10 md:pb-10 md:pt-32">
        <section class="relative min-w-0 text-center">
          <div class="mx-auto w-full max-w-[1200px]">
            <div class="reveal-up flex flex-col items-center justify-center">
              <p class="hero-badge mb-7 inline-flex max-w-full items-center gap-2 rounded-full px-4 py-2 text-[12px] font-medium uppercase tracking-[0.14em]">
                <span class="inline-flex size-1.5 rounded-full bg-[var(--primary)]"></span>
                <span class="truncate">{{ page.hero.badge }}</span>
              </p>
              <h1 class="hero-title home-title mx-auto w-full min-w-0 max-w-[1200px] break-words text-4xl font-semibold leading-[1.04] tracking-normal [font-family:var(--font-display)] sm:text-5xl lg:text-7xl lg:leading-[0.98] xl:text-[5.125rem]">
                {{ page.hero.title }}
              </h1>
              <p class="hero-copy home-copy mx-auto mt-8 w-full min-w-0 max-w-[620px] text-[17px] leading-[1.75]">
                {{ page.hero.description }}
              </p>
              <div class="mt-10 flex flex-wrap items-center justify-center gap-4">
                <router-link
                  :to="startPath"
                  class="home-primary-action group inline-flex min-h-14 items-center justify-center gap-3 rounded-full py-1.5 pl-7 pr-2 text-[15px] font-semibold transition-all duration-700 ease-[cubic-bezier(0.32,0.72,0,1)] active:scale-[0.98]"
                >
                  <span>{{ page.hero.primaryCta }}</span>
                  <span class="flex size-10 items-center justify-center rounded-full bg-[var(--background)] text-[color:var(--foreground)] transition-transform duration-700 ease-[cubic-bezier(0.32,0.72,0,1)] group-hover:-translate-y-0.5 group-hover:translate-x-0.5 group-hover:scale-105">
                    <Icon name="arrowRight" size="sm" />
                  </span>
                </router-link>
                <router-link
                  to="/models"
                  class="home-soft-button inline-flex h-14 items-center rounded-full px-7 text-[15px] font-medium transition-all duration-700 ease-[cubic-bezier(0.32,0.72,0,1)] active:scale-[0.98]"
                >
                  {{ page.hero.secondaryCta }}
                </router-link>
              </div>
            </div>
          </div>
        </section>

        <section class="mx-auto mt-12 w-full max-w-[1040px] md:mt-14">
          <div class="hero-product-visual overflow-hidden rounded-lg border border-[color:var(--surface-line)]">
            <div class="hero-product-header flex flex-wrap items-center justify-between gap-4 border-b border-[color:var(--surface-line)] px-5 py-4 sm:px-6">
              <div class="flex min-w-0 items-center gap-3">
                <span class="brand-icon flex size-9 shrink-0 items-center justify-center overflow-hidden rounded-full">
                  <img :src="activeBrand.logo" :alt="activeBrand.siteName" class="h-full w-full object-contain" />
                </span>
                <div class="min-w-0">
                  <p class="home-title truncate text-sm font-semibold">{{ page.heroPanel.title }}</p>
                  <p class="home-faint mt-0.5 truncate text-xs">{{ heroVisualCopy.gatewayProtocol }}</p>
                </div>
              </div>
              <span class="home-accent-panel inline-flex items-center gap-2 rounded-full px-3 py-1.5 text-xs font-medium">
                <span class="inline-flex size-1.5 rounded-full bg-[var(--primary)]"></span>
                {{ page.heroPanel.status }}
              </span>
            </div>

            <div class="hero-product-grid grid lg:grid-cols-[1fr_0.78fr_1fr]">
              <section class="hero-product-section p-5 sm:p-6">
                <p class="home-faint text-[11px] font-medium uppercase tracking-[0.14em]">{{ heroVisualCopy.providerLabel }}</p>
                <div class="mt-4">
                  <div
                    v-for="provider in providers"
                    :key="provider.name"
                    class="hero-provider-row flex items-center justify-between gap-4 border-b border-[color:var(--surface-line-soft)] py-4 last:border-b-0"
                  >
                    <div class="flex min-w-0 items-center gap-3">
                      <span class="provider-mark inline-flex size-9 shrink-0 items-center justify-center rounded-full text-xs font-semibold" :style="{ '--provider-color': provider.color }">
                        {{ provider.mark }}
                      </span>
                      <span class="home-title truncate text-sm font-semibold">{{ provider.name }}</span>
                    </div>
                    <span class="home-muted shrink-0 text-xs">{{ heroVisualCopy.providerStatus }}</span>
                  </div>
                </div>
              </section>

              <section class="hero-gateway-section flex min-h-52 flex-col items-center justify-center border-y border-[color:var(--surface-line)] p-6 text-center lg:border-x lg:border-y-0">
                <span class="hero-gateway-node flex size-20 items-center justify-center rounded-lg">
                  <img :src="activeBrand.logo" alt="" class="size-12 object-contain" />
                </span>
                <p class="home-title mt-5 text-base font-semibold">{{ activeBrand.siteName }} API</p>
                <p class="home-muted mt-2 max-w-40 text-xs leading-5">{{ heroVisualCopy.gatewayLabel }}</p>
              </section>

              <section class="hero-product-section p-5 sm:p-6">
                <p class="home-faint text-[11px] font-medium uppercase tracking-[0.14em]">{{ heroVisualCopy.metricsLabel }}</p>
                <div class="mt-4">
                  <div
                    v-for="metric in heroMetrics"
                    :key="metric.label"
                    class="flex items-start justify-between gap-5 border-b border-[color:var(--surface-line-soft)] py-3.5 last:border-b-0"
                  >
                    <span class="home-muted text-xs">{{ metric.label }}</span>
                    <span class="home-title max-w-[65%] break-words text-right font-mono text-xs font-semibold">{{ metric.value }}</span>
                  </div>
                </div>
              </section>
            </div>

            <div class="hero-price-strip flex flex-col gap-4 border-t border-[color:var(--surface-line)] px-5 py-5 sm:flex-row sm:items-center sm:justify-between sm:px-6">
              <div>
                <p class="home-muted text-xs">{{ heroVisualCopy.priceLabel }}</p>
                <p class="home-title mt-1 text-2xl font-semibold">1/50</p>
              </div>
              <div class="min-w-0 flex-1 sm:max-w-[56%]">
                <div class="hero-price-track h-2 overflow-hidden rounded-full">
                  <div class="hero-price-fill h-full rounded-full"></div>
                </div>
                <p class="home-faint mt-2 text-right text-[11px]">{{ heroVisualCopy.priceNote }}</p>
              </div>
            </div>
          </div>
        </section>

        <section class="w-full min-w-0 overflow-hidden pt-8">
          <p class="home-faint mb-5 text-center text-[10px] uppercase tracking-[0.24em] md:mb-7">
            {{ page.providersTitle }}
          </p>
          <div class="relative w-full min-w-0 overflow-hidden">
            <div class="pointer-events-none absolute inset-y-0 left-0 z-10 w-20 bg-gradient-to-r from-[var(--background)] to-transparent"></div>
            <div class="pointer-events-none absolute inset-y-0 right-0 z-10 w-20 bg-gradient-to-l from-[var(--background)] to-transparent"></div>
            <div class="home-marquee-track flex">
              <div
                v-for="(provider, index) in marqueeProviders"
                :key="`${provider.name}-${index}`"
                class="home-subpanel mx-2.5 flex shrink-0 items-center gap-3 rounded-full px-5 py-3"
              >
                <span class="provider-mark inline-flex size-[22px] shrink-0 items-center justify-center rounded-full text-[11px] font-semibold" :style="{ '--provider-color': provider.color }">
                  {{ provider.mark }}
                </span>
                <span class="home-title text-[14px]">{{ provider.name }}</span>
              </div>
            </div>
          </div>
        </section>
      </section>

      <section class="home-section px-6 md:px-10">
        <div class="mx-auto max-w-[1280px]">
          <div class="grid gap-6 lg:grid-cols-[0.92fr_1.08fr] lg:items-stretch">
            <div class="home-shell p-1.5">
              <div class="home-core flex h-full flex-col p-7 md:p-9">
                <p class="home-eyebrow text-[12px] font-medium uppercase tracking-[0.18em]">{{ page.modalities.kicker }}</p>
                <h2 class="home-title mt-5 text-[clamp(2rem,4.2vw,3.75rem)] font-semibold leading-[1.02] tracking-normal">
                  {{ page.modalities.title }}
                </h2>
                <p class="home-muted mt-5 max-w-[560px] text-[16px] leading-[1.8]">
                  {{ page.modalities.description }}
                </p>
                <div class="mt-auto grid gap-3 pt-8">
                  <div
                    v-for="metric in heroMetrics"
                    :key="metric.label"
                    class="home-subpanel flex items-center justify-between rounded-lg px-4 py-3"
                  >
                    <span class="home-muted text-[13px]">{{ metric.label }}</span>
                    <span class="home-title font-mono text-[14px] font-semibold">{{ metric.value }}</span>
                  </div>
                </div>
              </div>
            </div>

            <div class="grid gap-4 md:grid-cols-3">
              <article
                v-for="item in page.modalities.items"
                :key="item.title"
                class="home-shell p-1.5"
              >
                <div class="home-core flex h-full flex-col p-6">
                  <div class="home-accent-panel mb-10 inline-flex size-14 items-center justify-center text-[16px] font-semibold">
                    {{ item.mark }}
                  </div>
                  <h3 class="home-title text-[22px] font-semibold leading-snug tracking-normal">{{ item.title }}</h3>
                  <p class="home-muted mt-4 text-[14px] leading-[1.75]">{{ item.description }}</p>
                </div>
              </article>
            </div>
          </div>
        </div>
      </section>

      <section class="home-section px-6 pt-0 md:px-10">
        <div class="mx-auto max-w-[1280px]">
          <div class="mb-12 max-w-[760px]">
            <p class="home-eyebrow text-[12px] font-medium uppercase tracking-[0.18em]">{{ page.steps.kicker }}</p>
            <h2 class="home-title mt-5 text-[clamp(2rem,4vw,3.375rem)] font-semibold leading-[1.04] tracking-normal">
              {{ page.steps.title }}
            </h2>
            <p class="home-muted mt-5 text-[16px] leading-[1.8]">{{ page.steps.description }}</p>
          </div>
          <div class="grid gap-4 md:grid-cols-3">
            <article
              v-for="(step, index) in page.steps.items"
              :key="step.title"
              class="home-shell p-1.5"
            >
              <div class="home-core h-full p-6">
                <div class="mb-8 flex items-center justify-between">
                  <span class="home-accent-dot inline-flex size-3 rounded-full"></span>
                  <span class="home-faint font-mono text-[12px]">0{{ index + 1 }}</span>
                </div>
                <h3 class="home-title text-[21px] font-semibold leading-snug tracking-normal">{{ step.title }}</h3>
                <p class="home-muted mt-4 min-h-[5.25rem] text-[14px] leading-[1.75]">{{ step.description }}</p>
                <div class="home-subpanel mt-6 rounded-lg px-4 py-3 font-mono text-[12px] text-[color:var(--home-muted)]">
                  {{ step.meta }}
                </div>
              </div>
            </article>
          </div>
        </div>
      </section>

      <section class="home-section px-6 pt-0 md:px-10">
        <div class="mx-auto grid max-w-[1280px] gap-8 lg:grid-cols-[0.95fr_1.05fr] lg:items-center">
          <div>
            <p class="home-eyebrow text-[12px] font-medium uppercase tracking-[0.18em]">{{ page.ops.kicker }}</p>
            <h2 class="home-title mt-5 text-[clamp(2rem,4vw,3.375rem)] font-semibold leading-[1.04] tracking-normal">
              {{ page.ops.title }}
            </h2>
            <p class="home-muted mt-5 max-w-[620px] text-[16px] leading-[1.8]">{{ page.ops.description }}</p>
          </div>
          <div class="home-shell p-1.5">
            <div class="home-core p-6">
              <div class="grid gap-3">
                <article
                  v-for="item in page.ops.items"
                  :key="item.number"
                  class="home-subpanel p-5"
                >
                  <div class="flex min-w-0 items-start gap-4">
                    <span class="home-faint w-10 shrink-0 font-mono text-[13px]">{{ item.number }}</span>
                    <div class="min-w-0">
                      <h3 class="home-title text-[19px] font-semibold leading-snug tracking-normal">{{ item.title }}</h3>
                      <p class="home-muted mt-2 text-[14px] leading-[1.7]">{{ item.description }}</p>
                      <router-link :to="item.to" class="home-inline-link mt-4 inline-flex items-center gap-2 text-[13px] font-medium">
                        {{ item.cta }}
                        <Icon name="arrowRight" size="xs" />
                      </router-link>
                    </div>
                  </div>
                </article>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section class="home-section px-6 pt-0 md:px-10">
        <div class="mx-auto max-w-[1280px]">
          <div class="home-shell p-1.5">
            <div class="home-core grid gap-8 p-6 md:p-9 lg:grid-cols-[0.88fr_1.12fr] lg:items-center">
              <div>
                <h2 class="home-title text-[clamp(2rem,3.6vw,3.125rem)] font-semibold leading-[1.04] tracking-normal">{{ page.sdk.title }}</h2>
                <p class="home-muted mt-5 text-[16px] leading-[1.8]">{{ page.sdk.description }}</p>
                <div class="mt-10 flex flex-wrap gap-4">
                  <router-link :to="startPath" class="home-primary-action inline-flex h-11 items-center justify-center rounded-full px-5 text-[13px] font-medium">
                    {{ page.sdk.primaryCta }}
                  </router-link>
                  <router-link to="/models" class="home-soft-button inline-flex h-11 items-center justify-center rounded-full px-5 text-[13px] font-medium">
                    {{ page.sdk.secondaryCta }}
                  </router-link>
                </div>
              </div>
              <pre class="home-code overflow-x-auto rounded-lg p-5 text-[12px] leading-6"><code>{{ page.sdk.code }}</code></pre>
            </div>
          </div>
        </div>
      </section>

      <section id="pricing" class="home-section px-6 pt-0 md:px-10">
        <div class="mx-auto max-w-[1280px]">
          <div class="mb-10 flex flex-col justify-between gap-5 md:flex-row md:items-end">
            <div class="max-w-[760px]">
              <p class="home-eyebrow text-[12px] font-medium uppercase tracking-[0.18em]">{{ page.pricing.kicker }}</p>
              <h2 class="home-title mt-5 text-[clamp(2rem,4vw,3.375rem)] font-semibold leading-[1.04] tracking-normal">
                {{ page.pricing.title }}
              </h2>
              <p class="home-muted mt-5 text-[16px] leading-[1.8]">{{ page.pricing.description }}</p>
            </div>
            <router-link to="/models" class="home-soft-button inline-flex h-11 w-fit items-center justify-center rounded-full px-5 text-[13px] font-medium">
              {{ page.pricing.cta }}
            </router-link>
          </div>
          <div class="home-shell p-1.5">
            <div class="home-core overflow-hidden">
              <div class="overflow-x-auto">
                <table class="min-w-full text-left text-[14px]">
                  <thead class="home-faint text-[11px] uppercase tracking-[0.14em]">
                    <tr>
                      <th v-for="column in page.pricing.columns" :key="column" class="whitespace-nowrap px-5 py-4 font-medium">{{ column }}</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-[color:var(--surface-line)]">
                    <tr v-for="row in page.pricing.rows" :key="row.name">
                      <td class="home-title whitespace-nowrap px-5 py-5 font-semibold">{{ row.name }}</td>
                      <td class="home-muted whitespace-nowrap px-5 py-5">{{ row.modality }}</td>
                      <td class="home-muted whitespace-nowrap px-5 py-5">{{ row.input }}</td>
                      <td class="home-muted whitespace-nowrap px-5 py-5">{{ row.output }}</td>
                      <td class="home-muted whitespace-nowrap px-5 py-5">{{ row.cache }}</td>
                      <td class="home-muted whitespace-nowrap px-5 py-5">{{ row.context }}</td>
                      <td class="whitespace-nowrap px-5 py-5">
                        <span class="home-accent-panel rounded-full px-3 py-1 text-[12px]">{{ row.policy }}</span>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section class="home-section px-6 pt-0 md:px-10">
        <div class="mx-auto max-w-[1280px]">
          <h2 class="home-title max-w-[780px] text-[clamp(2rem,4vw,3.375rem)] font-semibold leading-[1.04] tracking-normal">
            {{ page.testimonials.title }}
          </h2>
          <div class="mt-12 grid gap-4 md:grid-cols-3">
            <article
              v-for="(item, index) in page.testimonials.items"
              :key="item.role"
              class="home-shell p-1.5"
            >
              <div class="home-core h-full p-6">
                <p class="home-copy text-[16px] leading-[1.7]">“{{ item.quote }}”</p>
                <div class="mt-6 flex items-center gap-3">
                  <span class="testimonial-avatar flex size-10 items-center justify-center rounded-full text-[13px] font-semibold" :style="{ '--avatar-color': testimonialColors[index] }">
                    {{ item.role.slice(0, 1) }}
                  </span>
                  <div>
                    <p class="home-title text-[14px] font-medium">{{ testimonialNames[index] }}</p>
                    <p class="home-faint text-[11px] uppercase tracking-[0.1em]">{{ item.role }}</p>
                  </div>
                </div>
              </div>
            </article>
          </div>
        </div>
      </section>

      <section class="home-section px-6 pt-0 md:px-10">
        <div class="mx-auto max-w-[1280px]">
          <div id="faq" class="mx-auto max-w-[980px]">
            <h2 class="home-title text-[clamp(2rem,4vw,3.375rem)] font-semibold leading-[1.04] tracking-normal">
              {{ page.faq.title }}
            </h2>
            <div class="mt-12">
              <article
                v-for="item in page.faq.items"
                :key="item.question"
                class="home-divider border-t py-8 first:border-t-0"
              >
                <h3 class="home-title text-[18px] font-semibold leading-snug tracking-normal">{{ item.question }}</h3>
                <p class="home-muted mt-3 text-[15px] leading-relaxed">{{ item.answer }}</p>
              </article>
            </div>
          </div>
        </div>
      </section>
    </div>
  </CloseAiPublicLayout>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import DOMPurify from 'dompurify'
import { marked } from 'marked'
import { useAuthStore, useAppStore } from '@/stores'
import CloseAiPublicLayout from '@/components/public/CloseAiPublicLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import {
  activeBrand,
  getCurrentApiBaseUrl,
  sharedBrandContent,
  withBrandTokens,
} from '@/brand'

const { locale } = useI18n()
const authStore = useAuthStore()
const appStore = useAppStore()

const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://') || content.startsWith('/')
})
const customHomeHtml = computed(() => {
  const content = homeContent.value.trim()
  if (!content || isHomeContentUrl.value) return ''
  return DOMPurify.sanitize(marked.parse(content, { async: false }) as string)
})
const startPath = computed(() => {
  if (!authStore.isAuthenticated) return '/login'
  return authStore.isAdmin ? '/admin/dashboard' : '/dashboard'
})

const providers = [
  { name: 'OpenAI', mark: 'O', color: '#2563eb' },
  { name: 'Claude', mark: 'C', color: '#0ea5e9' },
]
const marqueeProviders = computed(() => Array.from({ length: 6 }, () => providers).flat())
const heroMetrics = computed(() => locale.value.startsWith('zh')
  ? [
      { label: '统一端点', value: getCurrentApiBaseUrl() },
      { label: '支持厂商', value: 'OpenAI + Claude' },
      { label: '价格优势', value: '低至官方 1/50' },
    ]
  : [
      { label: 'Unified endpoint', value: getCurrentApiBaseUrl() },
      { label: 'Providers', value: 'OpenAI + Claude' },
      { label: 'Price advantage', value: 'From 1/50 official' },
    ])
const testimonialNames = computed(() => locale.value.startsWith('zh')
  ? ['低至官方 1/50', 'OpenAI + Claude', 'OpenAI 兼容']
  : ['From 1/50 official', 'OpenAI + Claude', 'OpenAI compatible'])
const testimonialColors = ['#2563eb', '#0ea5e9', '#1d4ed8']
const heroVisualCopy = computed(() => locale.value.startsWith('zh')
  ? {
      gatewayProtocol: 'OpenAI 兼容协议',
      providerLabel: '当前支持',
      providerStatus: '可用',
      gatewayLabel: '一个入口调用两家模型',
      metricsLabel: '接入信息',
      priceLabel: '部分模型最低为官方 API 价格的',
      priceNote: '具体价格以模型广场实时费率为准',
    }
  : {
      gatewayProtocol: 'OpenAI-compatible protocol',
      providerLabel: 'Currently supported',
      providerStatus: 'Available',
      gatewayLabel: 'One entry point for both providers',
      metricsLabel: 'Connection details',
      priceLabel: 'Selected models start from',
      priceNote: 'See the model marketplace for live rates',
    })

const page = computed(() => withBrandTokens(
  locale.value.startsWith('zh') ? sharedBrandContent.home.zh : sharedBrandContent.home.en,
))

onMounted(() => {
  authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>

<style scoped>
.home-surface {
  --font-display: ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  --background: #f8fafc;
  --foreground: #0f172a;
  --primary: rgb(var(--color-primary-600));
  --primary-foreground: #ffffff;
  --surface-shell: transparent;
  --surface-core: #ffffffd9;
  --surface-line: rgb(var(--color-primary-700) / 0.12);
  --surface-line-soft: rgb(var(--color-primary-700) / 0.08);
  --surface-highlight: #ffffff;
  --surface-control: #ffffffd9;
  --surface-control-hover: rgb(var(--color-primary-50));
  --surface-control-shadow: 0 10px 28px rgb(var(--color-primary-600) / 0.09);
  --home-heading: #0f172a;
  --home-copy: #334155;
  --home-muted: #64748b;
  --home-faint: #94a3b8;
  --home-layer: rgb(var(--color-primary-50) / 0.72);
  --home-code: #0f172a;
  --brand-ring: rgb(var(--color-primary-600) / 0.2);
  --brand-soft: rgb(var(--color-primary-100));
  --accent: #0d9488;
  --accent-strong: #0f766e;
  --accent-soft: #ccfbf1;
  --accent-ring: #14b8a633;
  background:
    linear-gradient(180deg, rgb(var(--color-primary-50)) 0%, #f8fafc 34%, #ffffff 100%),
    linear-gradient(90deg, rgb(var(--color-primary-200) / 0.28) 1px, transparent 1px),
    linear-gradient(180deg, rgb(var(--color-primary-200) / 0.28) 1px, transparent 1px),
    var(--background);
  background-size: auto, 80px 80px, 80px 80px;
}

.dark .home-surface {
  --background: #07111f;
  --foreground: #f8fbff;
  --primary: rgb(var(--color-primary-300));
  --primary-foreground: #07111f;
  --surface-shell: transparent;
  --surface-core: #0f1f35d9;
  --surface-line: rgb(var(--color-primary-300) / 0.15);
  --surface-line-soft: rgb(var(--color-primary-300) / 0.08);
  --surface-highlight: #172a46;
  --surface-control: #10223ab8;
  --surface-control-hover: #19365c;
  --surface-control-shadow: 0 18px 58px #0000002e;
  --home-heading: #f8fbff;
  --home-copy: #e2e8f0;
  --home-muted: #cbd5e1;
  --home-faint: #94a3b8;
  --home-layer: #10223ab8;
  --home-code: #061121;
  --brand-ring: rgb(var(--color-primary-300) / 0.24);
  --brand-soft: rgb(var(--color-primary-900));
  --accent: #2dd4bf;
  --accent-strong: #5eead4;
  --accent-soft: #134e4a;
  --accent-ring: #2dd4bf40;
}

.home-surface::before {
  content: "";
  position: fixed;
  inset: 0;
  z-index: -1;
  pointer-events: none;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noiseFilter'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.85' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noiseFilter)'/%3E%3C/svg%3E");
  background-size: 128px 128px;
  opacity: 0.03;
}

.home-section {
  padding-block: 7rem;
}

.home-title {
  color: var(--home-heading);
}

.hero-title,
.hero-copy {
  overflow-wrap: anywhere;
}

.home-copy {
  color: var(--home-copy);
}

.home-muted {
  color: var(--home-muted);
}

.home-faint {
  color: var(--home-faint);
}

.home-eyebrow {
  color: var(--accent-strong);
}

.hero-badge {
  background: var(--surface-control);
  border: 1px solid var(--surface-line);
  color: var(--accent-strong);
  box-shadow: var(--surface-control-shadow);
}

.home-shell {
  background: var(--surface-shell);
  border-radius: 0.5rem;
}

.home-core {
  border: 1px solid var(--surface-line);
  border-radius: 0.5rem;
  background:
    linear-gradient(180deg, var(--surface-highlight), transparent 58%),
    var(--surface-core);
  box-shadow: 0 20px 56px rgb(var(--color-primary-600) / 0.07);
}

.home-subpanel {
  border: 1px solid var(--surface-line-soft);
  border-radius: 0.5rem;
  background: var(--home-layer);
  box-shadow: 0 10px 28px rgb(var(--color-primary-600) / 0.05);
}

.home-accent-panel {
  border-radius: 0.5rem;
  background: var(--accent-soft);
  color: var(--accent-strong);
  box-shadow: 0 10px 28px rgb(var(--color-primary-600) / 0.04);
}

.home-accent-dot {
  background: var(--accent);
  box-shadow: 0 0 16px var(--accent-ring);
}

.home-primary-action {
  background: var(--primary);
  color: var(--primary-foreground);
  box-shadow: 0 18px 48px var(--brand-ring);
}

.home-primary-action:hover {
  filter: brightness(1.08);
  box-shadow: 0 22px 58px var(--brand-ring);
}

.home-soft-button {
  background: var(--surface-control);
  box-shadow: var(--surface-control-shadow);
  color: var(--home-muted);
}

.home-soft-button:hover {
  background: var(--surface-control-hover);
  color: var(--home-heading);
}

.home-inline-link {
  color: var(--accent-strong);
}

.home-code {
  background: var(--home-code);
  color: color-mix(in oklab, var(--home-heading) 72%, transparent);
}

.home-divider {
  border-color: var(--surface-line);
}

.hero-product-visual {
  background: var(--surface-core);
  box-shadow: 0 24px 70px rgb(var(--color-primary-600) / 0.16);
}

.hero-product-header,
.hero-price-strip {
  background: var(--surface-highlight);
}

.hero-product-section {
  background: color-mix(in oklab, var(--surface-core) 88%, var(--home-layer));
}

.hero-gateway-section {
  background: var(--home-layer);
}

.hero-gateway-node {
  background: var(--primary);
  box-shadow: 0 16px 42px var(--brand-ring);
}

.hero-price-track {
  background: color-mix(in oklab, var(--home-muted) 18%, transparent);
}

.hero-price-fill {
  width: 2%;
  min-width: 0.5rem;
  background: var(--primary);
}

.provider-mark {
  --provider-color: #2f647d;
  background: color-mix(in oklab, var(--provider-color) 16%, transparent);
  color: var(--provider-color);
}

.testimonial-avatar {
  --avatar-color: #2f647d;
  background: color-mix(in oklab, var(--avatar-color) 18%, transparent);
  color: var(--avatar-color);
}

.reveal-up {
  animation: reveal-up 0.72s cubic-bezier(0.32, 0.72, 0, 1) both;
}

.home-marquee-track {
  width: max-content;
  animation: home-marquee 35s linear infinite;
}

.home-marquee-track:hover {
  animation-play-state: paused;
}

@keyframes home-marquee {
  from {
    transform: translateX(0);
  }

  to {
    transform: translateX(-50%);
  }
}

@keyframes reveal-up {
  from {
    opacity: 0;
    filter: blur(10px);
    transform: translateY(28px);
  }

  to {
    opacity: 1;
    filter: blur(0);
    transform: translateY(0);
  }
}

@media (prefers-reduced-motion: reduce) {
  .home-marquee-track,
  .reveal-up {
    animation: none;
  }
}

@media (max-width: 640px) {
  .hero-title {
    max-width: calc(100vw - 3rem);
    font-size: 2.25rem;
    line-height: 1.04;
  }

  .hero-copy {
    max-width: calc(100vw - 3rem);
    font-size: 1rem;
  }

  .home-section {
    padding-block: 4.5rem;
  }
}
</style>
