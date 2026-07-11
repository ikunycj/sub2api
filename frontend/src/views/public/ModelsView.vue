<template>
  <CloseAiPublicLayout :page-title="page.metaTitle">
    <section class="border-b border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-950">
      <div class="mx-auto max-w-7xl px-4 py-14 sm:px-6 lg:px-8">
        <div class="max-w-4xl">
          <div class="mb-5 inline-flex items-center gap-2 rounded-md border border-blue-200 bg-blue-50 px-3 py-1.5 text-sm font-medium text-blue-700 dark:border-blue-900 dark:bg-blue-950/50 dark:text-blue-200">
            <Icon name="grid" size="sm" />
            <span>{{ page.badge }}</span>
          </div>
          <h1 class="break-words text-[2.125rem] font-semibold tracking-normal text-slate-950 max-sm:break-all dark:text-white sm:text-5xl">
            {{ page.title }}
          </h1>
          <p class="mt-5 max-w-3xl break-words text-lg leading-8 text-slate-600 dark:text-slate-300">
            {{ page.description }}
          </p>
        </div>
      </div>
    </section>

    <section class="bg-slate-50 py-8 dark:bg-slate-950">
      <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div class="rounded-lg border border-slate-200 bg-white p-4 shadow-sm dark:border-slate-800 dark:bg-slate-900">
          <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
            <label class="relative block flex-1">
              <Icon name="search" size="sm" class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
              <input
                v-model="searchQuery"
                type="search"
                :placeholder="page.searchPlaceholder"
                class="w-full rounded-xl border border-slate-200 bg-white py-2.5 pl-10 pr-3 text-sm text-slate-900 outline-none transition-all duration-200 placeholder:text-slate-400 hover:border-slate-300 focus:border-primary-500 focus:ring-4 focus:ring-primary-500/10 dark:border-slate-700 dark:bg-slate-950 dark:text-white dark:hover:border-slate-600"
              />
            </label>

            <div class="flex gap-2 overflow-x-auto">
              <button
                v-for="category in categories"
                :key="category.value"
                type="button"
                class="shrink-0 rounded-xl border px-3 py-2 text-sm font-medium transition-all duration-200 ease-out hover:-translate-y-px active:translate-y-0 active:scale-[0.98]"
                :class="activeCategory === category.value
                  ? 'border-primary-600 bg-primary-500 text-white shadow-sm shadow-primary-500/25'
                  : 'border-slate-200 bg-white text-slate-700 shadow-sm shadow-slate-950/5 hover:border-primary-200 hover:bg-primary-50/60 hover:text-primary-700 hover:shadow-md hover:shadow-primary-500/10 dark:border-slate-700 dark:bg-slate-950 dark:text-slate-200 dark:hover:border-primary-800 dark:hover:bg-primary-950/30 dark:hover:text-primary-200'"
                @click="activeCategory = category.value"
              >
                {{ category.label }}
              </button>
            </div>
          </div>
        </div>

        <div class="mt-6 grid gap-4 lg:grid-cols-[minmax(0,1fr)_320px]">
          <div class="overflow-hidden rounded-lg border border-slate-200 bg-white shadow-sm dark:border-slate-800 dark:bg-slate-900">
            <div class="hidden overflow-x-auto lg:block">
              <table class="min-w-full divide-y divide-slate-200 dark:divide-slate-800">
                <thead class="bg-slate-50 dark:bg-slate-950">
                  <tr>
                    <th v-for="column in page.columns" :key="column" class="px-5 py-3 text-left text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
                      {{ column }}
                    </th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-200 dark:divide-slate-800">
                  <tr v-for="model in filteredModels" :key="model.id" class="transition-colors duration-150 hover:bg-primary-50/50 dark:hover:bg-primary-950/20">
                    <td class="px-5 py-4">
                      <div class="font-semibold text-slate-950 dark:text-white">{{ model.name }}</div>
                      <div class="mt-1 font-mono text-xs text-slate-500 dark:text-slate-400">{{ model.id }}</div>
                    </td>
                    <td class="px-5 py-4 text-sm text-slate-700 dark:text-slate-200">{{ model.provider }}</td>
                    <td class="px-5 py-4">
                      <div class="flex flex-wrap gap-1.5">
                        <span
                          v-for="tag in model.tags"
                          :key="tag"
                          class="rounded-md bg-blue-50 px-2 py-1 text-xs font-medium text-blue-700 dark:bg-blue-950 dark:text-blue-200"
                        >
                          {{ tag }}
                        </span>
                      </div>
                    </td>
                    <td class="px-5 py-4 text-sm text-slate-700 dark:text-slate-200">{{ model.context }}</td>
                    <td class="px-5 py-4 text-sm text-slate-700 dark:text-slate-200">{{ model.billing }}</td>
                    <td class="px-5 py-4">
                      <span class="rounded-md bg-slate-100 px-2 py-1 text-xs font-semibold text-slate-700 dark:bg-slate-800 dark:text-slate-200">
                        {{ model.policy }}
                      </span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div class="grid gap-3 p-4 lg:hidden">
              <article
                v-for="model in filteredModels"
                :key="model.id"
                class="rounded-lg border border-slate-200 bg-slate-50 p-4 dark:border-slate-800 dark:bg-slate-950"
              >
                <div class="flex min-w-0 items-start justify-between gap-3">
                  <div class="min-w-0">
                    <h2 class="font-semibold text-slate-950 dark:text-white">{{ model.name }}</h2>
                    <p class="mt-1 font-mono text-xs text-slate-500 dark:text-slate-400">{{ model.id }}</p>
                  </div>
                  <span class="shrink-0 rounded-md bg-blue-50 px-2 py-1 text-xs font-semibold text-blue-700 dark:bg-blue-950 dark:text-blue-200">
                    {{ model.provider }}
                  </span>
                </div>
                <div class="mt-4 grid gap-2 text-sm text-slate-600 dark:text-slate-300">
                  <div>{{ page.mobile.context }}: {{ model.context }}</div>
                  <div>{{ page.mobile.billing }}: {{ model.billing }}</div>
                  <div>{{ page.mobile.policy }}: {{ model.policy }}</div>
                </div>
                <div class="mt-4 flex flex-wrap gap-1.5">
                  <span v-for="tag in model.tags" :key="tag" class="rounded-md bg-white px-2 py-1 text-xs dark:bg-slate-900">
                    {{ tag }}
                  </span>
                </div>
              </article>
            </div>
          </div>

          <aside class="grid gap-4 self-start">
            <div
              v-for="card in page.sidebar"
              :key="card.title"
              class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm dark:border-slate-800 dark:bg-slate-900"
              :class="card.highlight ? 'border-blue-200 bg-blue-50 dark:border-blue-900 dark:bg-blue-950/40' : ''"
            >
              <div class="mb-3 flex items-center gap-2 text-sm font-semibold text-slate-950 dark:text-white">
                <Icon :name="card.icon" size="sm" class="text-primary-600 dark:text-primary-300" />
                {{ card.title }}
              </div>
              <p class="text-sm leading-6 text-slate-600 dark:text-slate-300">{{ card.description }}</p>
              <router-link v-if="card.to" :to="card.to" class="mt-4 inline-flex items-center gap-2 rounded-full px-1 text-sm font-semibold text-primary-700 transition-colors hover:text-primary-800 focus:outline-none focus:ring-2 focus:ring-primary-500/25 dark:text-primary-300 dark:hover:text-primary-200">
                {{ card.cta }}
                <Icon name="arrowRight" size="sm" />
              </router-link>
            </div>
          </aside>
        </div>
      </div>
    </section>
  </CloseAiPublicLayout>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import CloseAiPublicLayout from '@/components/public/CloseAiPublicLayout.vue'
import Icon from '@/components/icons/Icon.vue'

type Category = 'all' | 'text' | 'reasoning' | 'vision' | 'code' | 'audio'
type SidebarIcon = 'filter' | 'sync' | 'terminal'

interface ModelItem {
  provider: string
  name: string
  id: string
  category: Category[]
  tags: string[]
  context: string
  billing: string
  policy: string
}

const { locale } = useI18n()
const searchQuery = ref('')
const activeCategory = ref<Category>('all')

const messages = {
  zh: {
    metaTitle: '模型广场',
    badge: '模型价格与 API 目录',
    title: '模型广场',
    description: '集中展示供应商、模型 ID、模态、上下文和计费口径。anytoken 额外提供企业可用的路由策略、预算控制和审计入口。',
    searchPlaceholder: '搜索模型、供应商或模型 ID',
    categories: {
      all: '全部',
      text: '文本',
      reasoning: '推理',
      vision: '视觉',
      code: '代码',
      audio: '音频',
    },
    columns: ['模型', '供应商', '能力', '上下文', '计费口径', '推荐策略'],
    mobile: {
      context: '上下文',
      billing: '计费口径',
      policy: '推荐策略',
    },
    sidebar: [
      { icon: 'filter' as SidebarIcon, title: '公开页负责发现', description: '模型广场面向售前和开发者解释有哪些模型、适合什么场景、如何接入。', cta: '', to: '', highlight: false },
      { icon: 'sync' as SidebarIcon, title: '企业版负责治理', description: '真实开通、额度、供应商 Key、路由和 fallback 仍在控制台完成。', cta: '', to: '', highlight: false },
      { icon: 'terminal' as SidebarIcon, title: 'OpenAI 兼容接入', description: '使用现有 SDK，把 base URL 切换为 anytoken.com 统一端点。', cta: '查看接入文档', to: '/docs', highlight: true },
    ],
    models: [
      { provider: 'OpenAI', name: 'GPT-4.1', id: 'openai/gpt-4.1', category: ['text', 'reasoning', 'code'], tags: ['文本', '推理', '代码'], context: '1M', billing: '实际费率', policy: '质量优先' },
      { provider: 'Anthropic', name: 'Claude Sonnet 4', id: 'anthropic/claude-sonnet-4', category: ['text', 'reasoning', 'code'], tags: ['文本', '长上下文', '代码'], context: '200K', billing: '实际费率', policy: '长文档' },
      { provider: 'Google', name: 'Gemini 2.5 Pro', id: 'google/gemini-2.5-pro', category: ['text', 'reasoning', 'vision'], tags: ['文本', '视觉', '长上下文'], context: '大上下文', billing: '实际费率', policy: '多模态' },
      { provider: 'DeepSeek', name: 'DeepSeek V3', id: 'deepseek/deepseek-chat', category: ['text', 'code'], tags: ['文本', '代码', '成本'], context: '128K', billing: '实际费率', policy: '成本优先' },
      { provider: 'Qwen', name: 'Qwen Max', id: 'qwen/qwen-max', category: ['text', 'reasoning'], tags: ['文本', '中文'], context: '中文友好', billing: '实际费率', policy: '中文业务' },
      { provider: 'OpenAI', name: 'GPT-4o mini', id: 'openai/gpt-4o-mini', category: ['text', 'vision'], tags: ['文本', '视觉', '高并发'], context: '通用', billing: '实际费率', policy: '高并发' },
      { provider: 'OpenAI', name: 'Embedding Large', id: 'openai/text-embedding-3-large', category: ['text'], tags: ['向量检索', '知识库'], context: '向量化', billing: '实际费率', policy: '知识库' },
      { provider: 'OpenAI', name: 'Audio Transcribe', id: 'openai/audio-transcribe', category: ['audio'], tags: ['音频', '转写'], context: '音频', billing: '实际费率', policy: '语音处理' },
    ] satisfies ModelItem[],
  },
  en: {
    metaTitle: 'Model Marketplace',
    badge: 'Model pricing and API catalog',
    title: 'Model Marketplace',
    description: 'Browse providers, model IDs, modalities, context windows, and billing conventions. anytoken adds enterprise routing, budget, and audit signals.',
    searchPlaceholder: 'Search model, provider, or model ID',
    categories: {
      all: 'All',
      text: 'Text',
      reasoning: 'Reasoning',
      vision: 'Vision',
      code: 'Code',
      audio: 'Audio',
    },
    columns: ['Model', 'Provider', 'Capabilities', 'Context', 'Billing', 'Policy'],
    mobile: {
      context: 'Context',
      billing: 'Billing',
      policy: 'Policy',
    },
    sidebar: [
      { icon: 'filter' as SidebarIcon, title: 'Public discovery', description: 'The marketplace explains which models exist, what they are good at, and how developers connect.', cta: '', to: '', highlight: false },
      { icon: 'sync' as SidebarIcon, title: 'Enterprise governance', description: 'Activation, quotas, provider keys, routing, and fallback remain in the console.', cta: '', to: '', highlight: false },
      { icon: 'terminal' as SidebarIcon, title: 'OpenAI-compatible access', description: 'Keep your SDK and switch the base URL to the anytoken.com unified endpoint.', cta: 'Read docs', to: '/docs', highlight: true },
    ],
    models: [
      { provider: 'OpenAI', name: 'GPT-4.1', id: 'openai/gpt-4.1', category: ['text', 'reasoning', 'code'], tags: ['Text', 'Reasoning', 'Code'], context: '1M', billing: 'Live rate', policy: 'Quality first' },
      { provider: 'Anthropic', name: 'Claude Sonnet 4', id: 'anthropic/claude-sonnet-4', category: ['text', 'reasoning', 'code'], tags: ['Text', 'Long', 'Code'], context: '200K', billing: 'Live rate', policy: 'Long docs' },
      { provider: 'Google', name: 'Gemini 2.5 Pro', id: 'google/gemini-2.5-pro', category: ['text', 'reasoning', 'vision'], tags: ['Text', 'Vision', 'Long'], context: 'Large', billing: 'Live rate', policy: 'Multimodal' },
      { provider: 'DeepSeek', name: 'DeepSeek V3', id: 'deepseek/deepseek-chat', category: ['text', 'code'], tags: ['Text', 'Code', 'Cost'], context: '128K', billing: 'Live rate', policy: 'Cost first' },
      { provider: 'Qwen', name: 'Qwen Max', id: 'qwen/qwen-max', category: ['text', 'reasoning'], tags: ['Text', 'Chinese'], context: 'Chinese', billing: 'Live rate', policy: 'Chinese business' },
      { provider: 'OpenAI', name: 'GPT-4o mini', id: 'openai/gpt-4o-mini', category: ['text', 'vision'], tags: ['Text', 'Vision', 'Fast'], context: 'General', billing: 'Live rate', policy: 'High throughput' },
      { provider: 'OpenAI', name: 'Embedding Large', id: 'openai/text-embedding-3-large', category: ['text'], tags: ['Embedding', 'RAG'], context: 'Vectors', billing: 'Live rate', policy: 'Knowledge base' },
      { provider: 'OpenAI', name: 'Audio Transcribe', id: 'openai/audio-transcribe', category: ['audio'], tags: ['Audio', 'Transcription'], context: 'Audio', billing: 'Live rate', policy: 'Speech workflows' },
    ] satisfies ModelItem[],
  },
} as const

const page = computed(() => locale.value.startsWith('zh') ? messages.zh : messages.en)
const categories = computed(() => ([
  { value: 'all' as Category, label: page.value.categories.all },
  { value: 'text' as Category, label: page.value.categories.text },
  { value: 'reasoning' as Category, label: page.value.categories.reasoning },
  { value: 'vision' as Category, label: page.value.categories.vision },
  { value: 'code' as Category, label: page.value.categories.code },
  { value: 'audio' as Category, label: page.value.categories.audio },
]))
const modelCatalog = computed<ModelItem[]>(() => page.value.models)
const filteredModels = computed(() => {
  const query = searchQuery.value.trim().toLowerCase()
  return modelCatalog.value.filter((model) => {
    const matchesCategory = activeCategory.value === 'all' || model.category.includes(activeCategory.value)
    const matchesQuery = !query
      || model.name.toLowerCase().includes(query)
      || model.id.toLowerCase().includes(query)
      || model.provider.toLowerCase().includes(query)
      || model.tags.some((tag) => tag.toLowerCase().includes(query))
    return matchesCategory && matchesQuery
  })
})
</script>
