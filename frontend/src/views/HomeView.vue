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

  <CloseAiPublicLayout v-else :page-title="page.metaTitle">
    <section class="relative overflow-hidden border-b border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-950">
      <div class="closeai-grid absolute inset-0 opacity-80"></div>
      <div class="relative mx-auto grid min-w-0 max-w-7xl items-center gap-12 px-4 py-16 sm:px-6 lg:grid-cols-[minmax(0,1fr)_minmax(360px,0.92fr)] lg:px-8 lg:py-20">
        <div class="min-w-0">
          <div class="mb-6 inline-flex items-center rounded-md border border-primary-200 bg-primary-50 px-3 py-1.5 text-sm font-medium text-primary-700 dark:border-primary-900 dark:bg-primary-950/50 dark:text-primary-200">
            {{ page.hero.badge }}
          </div>
          <h1 class="max-w-3xl break-words text-[2.125rem] font-semibold leading-tight tracking-normal text-slate-950 max-sm:break-all dark:text-white sm:text-5xl lg:text-6xl">
            {{ page.hero.title }}
          </h1>
          <p class="mt-6 max-w-2xl break-words text-lg leading-8 text-slate-600 dark:text-slate-300">
            {{ page.hero.description }}
          </p>
          <div class="mt-8 flex flex-col gap-3 sm:flex-row">
            <router-link
              :to="startPath"
              class="inline-flex items-center justify-center gap-2 rounded-full bg-primary-500 px-5 py-3 text-sm font-semibold text-white shadow-sm shadow-primary-500/25 transition-all duration-200 ease-out hover:-translate-y-px hover:bg-primary-600 hover:shadow-md hover:shadow-primary-500/30 active:translate-y-0 active:scale-95 focus:outline-none focus:ring-2 focus:ring-primary-500/35 focus:ring-offset-2 dark:focus:ring-offset-slate-950"
            >
              <span>{{ page.hero.primaryCta }}</span>
              <span class="flex h-5 w-5 items-center justify-center rounded-full bg-white/90 text-primary-700">
                <Icon name="arrowRight" size="xs" />
              </span>
            </router-link>
            <router-link
              to="/models"
              class="inline-flex items-center justify-center rounded-full border border-slate-200 bg-white px-5 py-3 text-sm font-semibold text-slate-700 shadow-sm shadow-slate-950/5 transition-all duration-200 ease-out hover:-translate-y-px hover:border-primary-200 hover:bg-primary-50/60 hover:text-primary-700 hover:shadow-md hover:shadow-primary-500/10 active:translate-y-0 active:scale-95 focus:outline-none focus:ring-2 focus:ring-primary-500/25 focus:ring-offset-2 dark:border-slate-700 dark:bg-slate-950 dark:text-slate-100 dark:hover:border-primary-800 dark:hover:bg-primary-950/30 dark:hover:text-primary-200 dark:focus:ring-offset-slate-950"
            >
              {{ page.hero.secondaryCta }}
            </router-link>
          </div>
        </div>

        <div class="min-w-0 rounded-xl border border-slate-200 bg-white p-5 shadow-xl shadow-primary-950/8 dark:border-slate-800 dark:bg-slate-900">
          <div class="mb-5 flex items-center justify-between">
            <div>
              <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ page.heroPanel.title }}</h2>
              <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ activeBrand.domain }}</p>
            </div>
            <span class="rounded-md bg-emerald-50 px-2.5 py-1 text-xs font-medium text-emerald-700 dark:bg-emerald-950/50 dark:text-emerald-300">
              {{ page.heroPanel.status }}
            </span>
          </div>
          <div class="grid gap-3">
            <article
              v-for="model in featuredModels"
              :key="model.id"
              class="rounded-lg border border-slate-200 bg-slate-50 p-4 dark:border-slate-800 dark:bg-slate-950/70"
            >
              <div class="flex min-w-0 items-start justify-between gap-3">
                <div class="min-w-0">
                  <h3 class="font-semibold text-slate-950 dark:text-white">{{ model.name }}</h3>
                  <p class="mt-1 font-mono text-xs text-slate-500 dark:text-slate-400">{{ model.id }}</p>
                </div>
                <span class="shrink-0 rounded-md bg-primary-50 px-2 py-1 text-xs font-semibold text-primary-700 dark:bg-primary-950 dark:text-primary-200">
                  {{ model.provider }}
                </span>
              </div>
              <div class="mt-3 flex flex-wrap gap-2 text-xs text-slate-600 dark:text-slate-300">
                <span>{{ model.capability }}</span>
                <span class="text-primary-700 dark:text-primary-300">{{ model.policy }}</span>
              </div>
            </article>
          </div>
        </div>
      </div>
    </section>

    <section class="bg-slate-50 py-8 dark:bg-slate-950">
      <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <p class="mb-4 text-sm font-medium text-slate-500 dark:text-slate-400">{{ page.providersTitle }}</p>
        <div class="flex flex-wrap gap-3">
          <span
            v-for="provider in providers"
            :key="provider"
            class="inline-flex items-center gap-2 rounded-md border border-slate-200 bg-white px-3 py-2 text-sm font-medium text-slate-700 dark:border-slate-800 dark:bg-slate-900 dark:text-slate-200"
          >
            <span class="h-2 w-2 rounded-full bg-primary-600"></span>
            {{ provider }}
          </span>
        </div>
      </div>
    </section>

    <section class="bg-slate-50 py-16 dark:bg-slate-950">
      <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <SectionHeading :kicker="page.modalities.kicker" :title="page.modalities.title" :description="page.modalities.description" />
        <div class="mt-10 grid gap-4 md:grid-cols-3">
          <article
            v-for="item in page.modalities.items"
            :key="item.title"
            class="rounded-lg border border-slate-200 bg-white p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900"
          >
            <div class="mb-5 flex h-12 w-12 items-center justify-center rounded-lg bg-primary-50 text-sm font-semibold text-primary-700 dark:bg-primary-950 dark:text-primary-200">
              {{ item.mark }}
            </div>
            <h3 class="text-lg font-semibold text-slate-950 dark:text-white">{{ item.title }}</h3>
            <p class="mt-3 text-sm leading-6 text-slate-600 dark:text-slate-300">{{ item.description }}</p>
          </article>
        </div>
      </div>
    </section>

    <section class="bg-white py-16 dark:bg-slate-950">
      <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <SectionHeading :kicker="page.steps.kicker" :title="page.steps.title" :description="page.steps.description" />
        <div class="mt-10 grid gap-4 md:grid-cols-3">
          <article
            v-for="(step, index) in page.steps.items"
            :key="step.title"
            class="rounded-lg border border-slate-200 bg-slate-50 p-6 dark:border-slate-800 dark:bg-slate-900"
          >
            <div class="mb-5 flex h-9 w-9 items-center justify-center rounded-full bg-primary-500 text-sm font-semibold text-white">
              {{ index + 1 }}
            </div>
            <h3 class="text-lg font-semibold text-slate-950 dark:text-white">{{ step.title }}</h3>
            <p class="mt-3 text-sm leading-6 text-slate-600 dark:text-slate-300">{{ step.description }}</p>
            <div class="mt-5 rounded-md border border-slate-200 bg-white px-3 py-2 text-xs text-slate-500 dark:border-slate-700 dark:bg-slate-950 dark:text-slate-400">
              {{ step.meta }}
            </div>
          </article>
        </div>
      </div>
    </section>

    <section class="bg-slate-50 py-16 dark:bg-slate-950">
      <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <SectionHeading :kicker="page.ops.kicker" :title="page.ops.title" :description="page.ops.description" />
        <div class="mt-10 grid gap-4 md:grid-cols-3">
          <article
            v-for="item in page.ops.items"
            :key="item.number"
            class="rounded-lg border border-slate-200 bg-white p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900"
          >
            <div class="text-3xl font-semibold text-primary-600 dark:text-primary-300">{{ item.number }}</div>
            <h3 class="mt-5 text-lg font-semibold text-slate-950 dark:text-white">{{ item.title }}</h3>
            <p class="mt-3 text-sm leading-6 text-slate-600 dark:text-slate-300">{{ item.description }}</p>
            <router-link :to="item.to" class="mt-5 inline-flex rounded-full px-1 text-sm font-semibold text-primary-700 transition-colors hover:text-primary-800 focus:outline-none focus:ring-2 focus:ring-primary-500/25 dark:text-primary-300 dark:hover:text-primary-200">
              {{ item.cta }}
            </router-link>
          </article>
        </div>
      </div>
    </section>

    <section class="bg-slate-50 pb-16 dark:bg-slate-950">
      <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div class="grid items-center gap-8 rounded-xl border border-slate-200 bg-white p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900 lg:grid-cols-[0.9fr_1.1fr] lg:p-10">
          <div>
            <h2 class="text-3xl font-semibold tracking-normal text-slate-950 dark:text-white">{{ page.sdk.title }}</h2>
            <p class="mt-4 text-base leading-7 text-slate-600 dark:text-slate-300">{{ page.sdk.description }}</p>
            <div class="mt-6 flex flex-col gap-3 sm:flex-row">
              <router-link :to="startPath" class="inline-flex items-center justify-center gap-2 rounded-full bg-primary-500 px-4 py-2.5 text-sm font-semibold text-white shadow-sm shadow-primary-500/25 transition-all duration-200 ease-out hover:-translate-y-px hover:bg-primary-600 hover:shadow-md hover:shadow-primary-500/30 active:translate-y-0 active:scale-95 focus:outline-none focus:ring-2 focus:ring-primary-500/35 focus:ring-offset-2 dark:focus:ring-offset-slate-900">
                {{ page.sdk.primaryCta }}
                <Icon name="arrowRight" size="xs" />
              </router-link>
              <router-link to="/models" class="inline-flex items-center justify-center rounded-full border border-slate-200 bg-white px-4 py-2.5 text-sm font-semibold text-slate-700 shadow-sm shadow-slate-950/5 transition-all duration-200 ease-out hover:-translate-y-px hover:border-primary-200 hover:bg-primary-50/60 hover:text-primary-700 hover:shadow-md hover:shadow-primary-500/10 active:translate-y-0 active:scale-95 dark:border-slate-700 dark:bg-slate-950 dark:text-slate-100 dark:hover:border-primary-800 dark:hover:bg-primary-950/30 dark:hover:text-primary-200">
                {{ page.sdk.secondaryCta }}
              </router-link>
            </div>
          </div>
          <pre class="overflow-x-auto rounded-lg bg-slate-950 p-5 text-xs leading-6 text-primary-100"><code>{{ page.sdk.code }}</code></pre>
        </div>
      </div>
    </section>

    <section class="bg-white py-16 dark:bg-slate-950">
      <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div class="flex flex-col justify-between gap-4 md:flex-row md:items-end">
          <SectionHeading :kicker="page.pricing.kicker" :title="page.pricing.title" :description="page.pricing.description" />
          <router-link to="/models" class="inline-flex w-fit items-center gap-2 rounded-full bg-primary-500 px-4 py-2.5 text-sm font-semibold text-white shadow-sm shadow-primary-500/25 transition-all duration-200 ease-out hover:-translate-y-px hover:bg-primary-600 hover:shadow-md hover:shadow-primary-500/30 active:translate-y-0 active:scale-95">
            {{ page.pricing.cta }}
            <Icon name="arrowRight" size="xs" />
          </router-link>
        </div>
        <div class="mt-8 overflow-hidden rounded-lg border border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900">
          <div class="overflow-x-auto">
            <table class="min-w-full divide-y divide-slate-200 text-sm dark:divide-slate-800">
              <thead class="bg-slate-50 text-left text-xs font-semibold uppercase tracking-wide text-slate-500 dark:bg-slate-950 dark:text-slate-400">
                <tr>
                  <th v-for="column in page.pricing.columns" :key="column" class="px-5 py-3">{{ column }}</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-200 dark:divide-slate-800">
                <tr v-for="row in page.pricing.rows" :key="row.name" class="text-slate-700 dark:text-slate-200">
                  <td class="px-5 py-4 font-medium text-slate-950 dark:text-white">{{ row.name }}</td>
                  <td class="px-5 py-4">{{ row.modality }}</td>
                  <td class="px-5 py-4">{{ row.input }}</td>
                  <td class="px-5 py-4">{{ row.output }}</td>
                  <td class="px-5 py-4">{{ row.cache }}</td>
                  <td class="px-5 py-4">{{ row.context }}</td>
                  <td class="px-5 py-4 text-primary-700 dark:text-primary-300">{{ row.policy }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </section>

    <section class="bg-slate-50 py-16 dark:bg-slate-950">
      <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <SectionHeading :kicker="page.testimonials.kicker" :title="page.testimonials.title" />
        <div class="mt-10 grid gap-4 md:grid-cols-3">
          <article
            v-for="item in page.testimonials.items"
            :key="item.role"
            class="rounded-lg border border-slate-200 bg-white p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900"
          >
            <div class="text-4xl leading-none text-primary-600">“</div>
            <p class="mt-4 text-sm leading-6 text-slate-600 dark:text-slate-300">{{ item.quote }}</p>
            <p class="mt-6 text-sm font-semibold text-slate-950 dark:text-white">{{ item.role }}</p>
          </article>
        </div>
      </div>
    </section>

    <section class="bg-slate-50 pb-16 dark:bg-slate-950">
      <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <SectionHeading :kicker="page.faq.kicker" :title="page.faq.title" />
        <div class="mt-8 grid gap-3">
          <article
            v-for="item in page.faq.items"
            :key="item.question"
            class="rounded-lg border border-slate-200 bg-white p-5 dark:border-slate-800 dark:bg-slate-900"
          >
            <h3 class="text-base font-semibold text-slate-950 dark:text-white">{{ item.question }}</h3>
            <p class="mt-2 text-sm leading-6 text-slate-600 dark:text-slate-300">{{ item.answer }}</p>
          </article>
        </div>
      </div>
    </section>
  </CloseAiPublicLayout>
</template>

<script setup lang="ts">
import { computed, defineComponent, h, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import DOMPurify from 'dompurify'
import { marked } from 'marked'
import { useAuthStore, useAppStore } from '@/stores'
import CloseAiPublicLayout from '@/components/public/CloseAiPublicLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { activeBrand } from '@/brand'

const SectionHeading = defineComponent({
  props: {
    kicker: { type: String, required: true },
    title: { type: String, required: true },
    description: { type: String, default: '' },
  },
  setup(props) {
    return () => h('div', { class: 'max-w-3xl' }, [
      h('div', { class: 'mb-4 inline-flex rounded-md border border-primary-200 bg-primary-50 px-3 py-1.5 text-sm font-medium text-primary-700 dark:border-primary-900 dark:bg-primary-950/50 dark:text-primary-200' }, props.kicker),
      h('h2', { class: 'text-3xl font-semibold tracking-normal text-slate-950 dark:text-white sm:text-4xl' }, props.title),
      props.description ? h('p', { class: 'mt-4 text-base leading-7 text-slate-600 dark:text-slate-300' }, props.description) : null,
    ])
  }
})

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

const providers = ['OpenAI', 'Gemini', 'Claude', 'DeepSeek', 'ByteDance', 'ElevenLabs', 'Minimax', 'Kling', 'Vidu', 'Grok', 'Wan', 'Runway']
const featuredModels = computed(() => locale.value.startsWith('zh')
  ? [
      { provider: 'OpenAI', name: 'GPT-4.1', id: 'openai/gpt-4.1', capability: '文本 / 代码', policy: '质量优先' },
      { provider: 'Anthropic', name: 'Claude Sonnet 4', id: 'anthropic/claude-sonnet-4', capability: '长上下文', policy: '文档处理' },
      { provider: 'Google', name: 'Gemini 2.5 Pro', id: 'google/gemini-2.5-pro', capability: '视觉 / 推理', policy: '多模态' },
    ]
  : [
      { provider: 'OpenAI', name: 'GPT-4.1', id: 'openai/gpt-4.1', capability: 'Text / Code', policy: 'Quality first' },
      { provider: 'Anthropic', name: 'Claude Sonnet 4', id: 'anthropic/claude-sonnet-4', capability: 'Long context', policy: 'Document work' },
      { provider: 'Google', name: 'Gemini 2.5 Pro', id: 'google/gemini-2.5-pro', capability: 'Vision / Reasoning', policy: 'Multimodal' },
    ])

const page = computed(() => locale.value.startsWith('zh') ? activeBrand.home.zh : activeBrand.home.en)

onMounted(() => {
  authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>

<style scoped>
.closeai-grid {
  background-image:
    linear-gradient(rgb(var(--color-primary-600) / 0.08) 1px, transparent 1px),
    linear-gradient(90deg, rgb(var(--color-primary-600) / 0.08) 1px, transparent 1px),
    linear-gradient(180deg, rgb(var(--color-primary-50) / 0.98), rgba(255, 255, 255, 0.78));
  background-size: 48px 48px, 48px 48px, auto;
}
</style>
