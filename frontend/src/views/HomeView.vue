<template>
  <div v-if="homeContent.trim()" class="min-h-screen">
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <div v-else v-html="homeContent"></div>
  </div>

  <CloseAiPublicLayout v-else :page-title="page.metaTitle">
    <section class="relative overflow-hidden border-b border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-950">
      <div class="closeai-grid absolute inset-0 opacity-80"></div>
      <div class="relative mx-auto grid min-w-0 max-w-7xl items-center gap-12 px-4 py-16 sm:px-6 lg:grid-cols-[minmax(0,1fr)_minmax(360px,0.92fr)] lg:px-8 lg:py-20">
        <div class="min-w-0">
          <div class="mb-6 inline-flex items-center rounded-md border border-blue-200 bg-blue-50 px-3 py-1.5 text-sm font-medium text-blue-700 dark:border-blue-900 dark:bg-blue-950/50 dark:text-blue-200">
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

        <div class="min-w-0 rounded-xl border border-slate-200 bg-white p-5 shadow-xl shadow-blue-950/8 dark:border-slate-800 dark:bg-slate-900">
          <div class="mb-5 flex items-center justify-between">
            <div>
              <h2 class="text-sm font-semibold text-slate-950 dark:text-white">{{ page.heroPanel.title }}</h2>
              <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">anytoken.com</p>
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
                <span class="shrink-0 rounded-md bg-blue-50 px-2 py-1 text-xs font-semibold text-blue-700 dark:bg-blue-950 dark:text-blue-200">
                  {{ model.provider }}
                </span>
              </div>
              <div class="mt-3 flex flex-wrap gap-2 text-xs text-slate-600 dark:text-slate-300">
                <span>{{ model.capability }}</span>
                <span class="text-blue-700 dark:text-blue-300">{{ model.policy }}</span>
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
            <span class="h-2 w-2 rounded-full bg-blue-600"></span>
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
            <div class="mb-5 flex h-12 w-12 items-center justify-center rounded-lg bg-blue-50 text-sm font-semibold text-blue-700 dark:bg-blue-950 dark:text-blue-200">
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
            <div class="text-3xl font-semibold text-blue-600 dark:text-blue-300">{{ item.number }}</div>
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
          <pre class="overflow-x-auto rounded-lg bg-slate-950 p-5 text-xs leading-6 text-blue-100"><code>{{ page.sdk.code }}</code></pre>
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
                  <td class="px-5 py-4 text-blue-700 dark:text-blue-300">{{ row.policy }}</td>
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
            <div class="text-4xl leading-none text-blue-600">“</div>
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
import { useAuthStore, useAppStore } from '@/stores'
import CloseAiPublicLayout from '@/components/public/CloseAiPublicLayout.vue'
import Icon from '@/components/icons/Icon.vue'

const SectionHeading = defineComponent({
  props: {
    kicker: { type: String, required: true },
    title: { type: String, required: true },
    description: { type: String, default: '' },
  },
  setup(props) {
    return () => h('div', { class: 'max-w-3xl' }, [
      h('div', { class: 'mb-4 inline-flex rounded-md border border-blue-200 bg-blue-50 px-3 py-1.5 text-sm font-medium text-blue-700 dark:border-blue-900 dark:bg-blue-950/50 dark:text-blue-200' }, props.kicker),
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

const messages = {
  zh: {
    metaTitle: 'AI 模型 API 网关与市场',
    hero: {
      badge: 'anytoken.com 企业版',
      title: '一个 API Key 接入主流 AI 模型',
      description: '通过兼容 OpenAI 的统一网关调用 GPT、Claude、Gemini、DeepSeek 等模型，计费、路由和模型发现都在一个工作台完成。',
      primaryCta: '开始使用',
      secondaryCta: '探索模型',
    },
    heroPanel: {
      title: '模型市场实时预览',
      status: '企业可用',
    },
    providersTitle: '在一个模型聚合平台浏览主流供应商',
    modalities: {
      kicker: '模型目录',
      title: '用一个目录接入文本以外的能力',
      description: '在 anytoken 模型中心评估聊天、视觉、图像、视频和语音模型，不必在多个供应商控制台之间来回切换。',
      items: [
        { mark: '图', title: '视觉与图像模型', description: '在同一个 API 市场选择图像生成和视觉理解模型，和语言模型共用一套接入方式。' },
        { mark: '视', title: '视频生成任务', description: '把创意视频和产品视频任务发送到网关托管端点，开发接口保持一致。' },
        { mark: '音', title: '语音与音频能力', description: '用和 LLM 工作流相同的运营层接入语音合成、转写和音频理解模型。' },
      ],
    },
    steps: {
      kicker: '快速接入',
      title: '从选模型到发出第一条请求',
      description: '在模型广场挑选模型，配置企业预算，再通过统一 API 接入应用。少写供应商分支，把时间留给产品本身。',
      items: [
        { title: '创建工作台', description: '为你的团队设置 anytoken 工作台、组织和环境。', meta: '组织 / 环境' },
        { title: '配置预算和供应商', description: '一套预算规则覆盖供应商、模型和客户环境。', meta: '$20  $100  $500' },
        { title: '签发网关 Key', description: '创建统一 API Key，让请求通过 anytoken 自动路由。', meta: 'API_KEY 可开始集成' },
      ],
    },
    ops: {
      kicker: '运营层',
      title: '统一模型接入的运营层',
      description: '当团队需要比较供应商、统一接入规范、快速上线 AI 功能时，统一 API 平台能减少重复集成和后续维护成本。',
      items: [
        { number: '01', title: '网关级路由', description: '把流量送到合适的供应商路径，同时让文本、视觉和 Agent 场景保持一致请求格式。', cta: '探索路由', to: '/models' },
        { number: '02', title: '选模型不用堆适配器', description: '在同一目录里比较成本、能力和可用性，再决定哪些模型接入生产流量。', cta: '比较模型', to: '/models' },
        { number: '03', title: '面向 B2B 团队的集成控制', description: '统一 Key、环境、计费和供应商访问，加快客户交付。', cta: '阅读集成文档', to: '/docs' },
      ],
    },
    sdk: {
      title: '保留现有 SDK，灵活切换模型路径',
      description: 'anytoken 提供由模型网关支撑的 OpenAI 兼容 API。切换模型家族时，通常只需要改配置，不必重写调用代码。',
      primaryCta: '生成 API Key',
      secondaryCta: '比较价格',
      code: `import OpenAI from "openai";

const client = new OpenAI({
  baseURL: "https://api.anytoken.com/v1",
  apiKey: process.env.ANYTOKEN_API_KEY,
});

await client.chat.completions.create({
  model: "anytoken/auto",
  messages: [{ role: "user", content: "Optimize this logic." }],
});`,
    },
    pricing: {
      kicker: '价格与成本',
      title: '路由流量前先看清成本',
      description: '模型进入生产环境前，先比较不同供应商的 Token 成本、上下文窗口和适用场景。最终价格以控制台实际配置为准。',
      cta: '查看全部模型',
      columns: ['模型名称', '模态', '输入', '输出', '缓存读取', '上下文窗口', '推荐策略'],
      rows: [
        { name: 'DeepSeek V3', modality: '文本', input: '实际费率', output: '实际费率', cache: '实际费率', context: '128K', policy: '成本优先' },
        { name: 'GPT-4.1', modality: '文本 / 图像', input: '实际费率', output: '实际费率', cache: '实际费率', context: '1M', policy: '质量优先' },
        { name: 'Claude Sonnet 4', modality: '文本', input: '实际费率', output: '实际费率', cache: '实际费率', context: '200K', policy: '长文本' },
        { name: 'Qwen Max', modality: '文本', input: '实际费率', output: '实际费率', cache: '实际费率', context: '中文友好', policy: '中文业务' },
      ],
    },
    testimonials: {
      kicker: '客户视角',
      title: '让团队少做模型运营杂活',
      items: [
        { role: 'AI 平台负责人', quote: '我们的评测框架继续保持 OpenAI 兼容，产品团队也能通过一个网关测试多个供应商。' },
        { role: 'AI 基础设施工程师', quote: '上线前，我们能先看清成本、端点和模型适配度。' },
        { role: 'SaaS 创始人', quote: '客户指定不同模型供应商时，统一 API 能明显缩短交付周期。' },
      ],
    },
    faq: {
      kicker: 'FAQ',
      title: '标准化模型接入前团队常问的问题',
      items: [
        { question: 'anytoken 和单一供应商 API 有什么不同？', answer: 'anytoken 把模型目录、路由网关、统一预算和审计放在一起。团队可以跨供应商调用模型，不用为每家供应商单独维护集成。' },
        { question: '现有 OpenAI SDK 代码还能继续使用吗？', answer: '可以。替换 base URL 和 API Key 后，就能继续使用熟悉的聊天补全请求方式。' },
        { question: 'anytoken 如何帮助 B2B 产品团队？', answer: '统一 Key、环境、预算和供应商访问，减少客户指定模型时的重复交付。' },
        { question: '费用在哪里管理？', answer: '在控制台查看模型支出、预算告警和团队成本归因。' },
      ],
    },
  },
  en: {
    metaTitle: 'AI Model API Gateway and Marketplace',
    hero: {
      badge: 'anytoken.com enterprise edition',
      title: 'One API key for leading AI models',
      description: 'Call GPT, Claude, Gemini, DeepSeek, and more through an OpenAI-compatible gateway, with billing, routing, and model discovery in one workspace.',
      primaryCta: 'Start',
      secondaryCta: 'Explore models',
    },
    heroPanel: {
      title: 'Model marketplace preview',
      status: 'Enterprise ready',
    },
    providersTitle: 'Browse leading providers in one model aggregation platform',
    modalities: {
      kicker: 'Model catalog',
      title: 'Use one catalog for capabilities beyond text',
      description: 'Evaluate chat, vision, image, video, and audio models in the anytoken model center without jumping across provider consoles.',
      items: [
        { mark: 'I', title: 'Vision and image models', description: 'Select image generation and vision models in one API marketplace with the same integration pattern as language models.' },
        { mark: 'V', title: 'Video generation tasks', description: 'Send product and creative video tasks through hosted gateway endpoints while keeping the developer interface consistent.' },
        { mark: 'A', title: 'Speech and audio capabilities', description: 'Connect TTS, transcription, and audio understanding through the same operational layer as LLM workflows.' },
      ],
    },
    steps: {
      kicker: 'Quick start',
      title: 'From model selection to first request',
      description: 'Pick a model, configure enterprise budgets, and connect your app through one API. Spend less time on provider branches and more time on product.',
      items: [
        { title: 'Create workspace', description: 'Set up your anytoken workspace, organization, and environments.', meta: 'Org / environments' },
        { title: 'Configure budgets and providers', description: 'Use one budget policy across providers, models, and customer environments.', meta: '$20  $100  $500' },
        { title: 'Issue gateway key', description: 'Create one API key and route requests automatically through anytoken.', meta: 'API_KEY ready' },
      ],
    },
    ops: {
      kicker: 'Operations',
      title: 'The operations layer for unified model access',
      description: 'When teams compare providers, standardize integration, and ship AI features quickly, a unified API platform reduces duplicate work and maintenance.',
      items: [
        { number: '01', title: 'Gateway routing', description: 'Route traffic to the right provider path while keeping text, vision, and agent requests consistent.', cta: 'Explore routing', to: '/models' },
        { number: '02', title: 'Choose models without adapters', description: 'Compare cost, capability, and availability in one catalog before moving traffic to production.', cta: 'Compare models', to: '/models' },
        { number: '03', title: 'B2B integration control', description: 'Unify keys, environments, billing, and provider access to accelerate customer delivery.', cta: 'Read docs', to: '/docs' },
      ],
    },
    sdk: {
      title: 'Keep your SDK and switch model paths flexibly',
      description: 'anytoken provides an OpenAI-compatible API backed by a model gateway. Switching model families usually means changing config, not application code.',
      primaryCta: 'Generate API key',
      secondaryCta: 'Compare pricing',
      code: `import OpenAI from "openai";

const client = new OpenAI({
  baseURL: "https://api.anytoken.com/v1",
  apiKey: process.env.ANYTOKEN_API_KEY,
});

await client.chat.completions.create({
  model: "anytoken/auto",
  messages: [{ role: "user", content: "Optimize this logic." }],
});`,
    },
    pricing: {
      kicker: 'Pricing and cost',
      title: 'Understand cost before routing traffic',
      description: 'Compare token cost, context windows, and scenarios before models enter production. Final rates follow the actual console configuration.',
      cta: 'View all models',
      columns: ['Model', 'Modality', 'Input', 'Output', 'Cache read', 'Context', 'Policy'],
      rows: [
        { name: 'DeepSeek V3', modality: 'Text', input: 'Live rate', output: 'Live rate', cache: 'Live rate', context: '128K', policy: 'Cost first' },
        { name: 'GPT-4.1', modality: 'Text / Image', input: 'Live rate', output: 'Live rate', cache: 'Live rate', context: '1M', policy: 'Quality first' },
        { name: 'Claude Sonnet 4', modality: 'Text', input: 'Live rate', output: 'Live rate', cache: 'Live rate', context: '200K', policy: 'Long text' },
        { name: 'Qwen Max', modality: 'Text', input: 'Live rate', output: 'Live rate', cache: 'Live rate', context: 'Chinese', policy: 'Chinese business' },
      ],
    },
    testimonials: {
      kicker: 'Customer view',
      title: 'Less model operations work for the team',
      items: [
        { role: 'AI platform lead', quote: 'Our evaluation framework stays OpenAI compatible, and product teams can test multiple providers through one gateway.' },
        { role: 'AI infrastructure engineer', quote: 'Before launch, we can see cost, endpoints, and model fit clearly.' },
        { role: 'SaaS founder', quote: 'When customers request different model providers, one API shortens the delivery cycle.' },
      ],
    },
    faq: {
      kicker: 'FAQ',
      title: 'Questions teams ask before standardizing model access',
      items: [
        { question: 'How is anytoken different from a single provider API?', answer: 'anytoken brings model catalog, gateway routing, unified budgets, and audit logs together so teams can call models across providers without maintaining separate integrations.' },
        { question: 'Can we keep existing OpenAI SDK code?', answer: 'Yes. Replace the base URL and API key, then keep using familiar chat completion requests.' },
        { question: 'How does anytoken help B2B product teams?', answer: 'It unifies keys, environments, budgets, and provider access, reducing duplicate delivery work when customers require specific models.' },
        { question: 'Where is spend managed?', answer: 'The console shows model spend, budget alerts, and team cost attribution.' },
      ],
    },
  },
} as const

const page = computed(() => locale.value.startsWith('zh') ? messages.zh : messages.en)

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
    linear-gradient(rgba(37, 99, 235, 0.08) 1px, transparent 1px),
    linear-gradient(90deg, rgba(37, 99, 235, 0.08) 1px, transparent 1px),
    linear-gradient(180deg, rgba(239, 246, 255, 0.98), rgba(255, 255, 255, 0.78));
  background-size: 48px 48px, 48px 48px, auto;
}
</style>
