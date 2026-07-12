<template>
  <CloseAiPublicLayout :page-title="page.metaTitle">
    <section class="bg-white dark:bg-slate-950">
      <div class="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
        <div class="max-w-5xl">
          <h1 class="break-words text-[2.375rem] font-semibold tracking-normal text-slate-950 max-sm:break-all dark:text-white sm:text-5xl">
            {{ page.title }}
          </h1>
          <p class="mt-4 max-w-3xl break-words text-base leading-7 text-slate-600 dark:text-slate-300 sm:text-lg">
            {{ page.description }}
          </p>
        </div>

        <div class="mt-8 grid gap-3 lg:grid-cols-[minmax(280px,1fr)_minmax(160px,220px)_minmax(160px,220px)_auto_auto] lg:items-center">
          <label class="relative block">
            <Icon name="search" size="md" class="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-slate-400" />
            <input
              v-model="searchQuery"
              type="search"
              :placeholder="page.searchPlaceholder"
              class="h-12 w-full rounded-lg border border-slate-200 bg-white pl-12 pr-4 text-sm text-slate-900 outline-none transition-all duration-200 placeholder:text-slate-400 hover:border-slate-300 focus:border-primary-500 focus:ring-4 focus:ring-primary-500/10 dark:border-slate-800 dark:bg-slate-900 dark:text-white dark:hover:border-slate-700"
            />
          </label>

          <label class="relative block">
            <span class="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-sm font-medium text-slate-500 dark:text-slate-400">
              {{ page.platformLabel }}
            </span>
            <select v-model="selectedPlatform" class="h-12 w-full appearance-none rounded-lg border border-slate-200 bg-white pl-16 pr-10 text-sm font-medium text-slate-800 outline-none transition-all duration-200 hover:border-slate-300 focus:border-primary-500 focus:ring-4 focus:ring-primary-500/10 dark:border-slate-800 dark:bg-slate-900 dark:text-slate-100 dark:hover:border-slate-700">
              <option value="all">{{ page.all }}</option>
              <option v-for="platform in platforms" :key="platform" :value="platform">{{ platform }}</option>
            </select>
            <Icon name="chevronDown" size="sm" class="pointer-events-none absolute right-4 top-1/2 -translate-y-1/2 text-slate-400" />
          </label>

          <label class="relative block">
            <span class="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-sm font-medium text-slate-500 dark:text-slate-400">
              {{ page.groupLabel }}
            </span>
            <select v-model="selectedGroup" class="h-12 w-full appearance-none rounded-lg border border-slate-200 bg-white pl-14 pr-10 text-sm font-medium text-slate-800 outline-none transition-all duration-200 hover:border-slate-300 focus:border-primary-500 focus:ring-4 focus:ring-primary-500/10 dark:border-slate-800 dark:bg-slate-900 dark:text-slate-100 dark:hover:border-slate-700">
              <option v-for="group in groupOptions" :key="group.value" :value="group.value">{{ group.label }}</option>
            </select>
            <Icon name="chevronDown" size="sm" class="pointer-events-none absolute right-4 top-1/2 -translate-y-1/2 text-slate-400" />
          </label>

          <button
            type="button"
            class="inline-flex h-12 items-center justify-center gap-2 rounded-lg border border-slate-200 bg-white px-4 text-sm font-semibold text-slate-700 shadow-sm shadow-slate-950/5 transition-all duration-200 hover:border-primary-200 hover:bg-primary-50/60 hover:text-primary-700 dark:border-slate-800 dark:bg-slate-900 dark:text-slate-200 dark:hover:border-primary-800 dark:hover:bg-primary-950/30"
            @click="resetFilters"
          >
            <Icon name="refresh" size="sm" />
            <span>{{ page.reset }}</span>
          </button>

          <div class="grid h-12 grid-cols-2 rounded-lg border border-slate-200 bg-white p-1 shadow-sm shadow-slate-950/5 dark:border-slate-800 dark:bg-slate-900">
            <button
              type="button"
              class="flex min-w-10 items-center justify-center rounded-md transition-colors"
              :class="viewMode === 'grid' ? 'bg-primary-50 text-primary-700 dark:bg-primary-950/50 dark:text-primary-200' : 'text-slate-500 hover:text-slate-900 dark:text-slate-400 dark:hover:text-white'"
              :title="page.gridView"
              @click="viewMode = 'grid'"
            >
              <Icon name="grid" size="sm" />
            </button>
            <button
              type="button"
              class="flex min-w-10 items-center justify-center rounded-md transition-colors"
              :class="viewMode === 'list' ? 'bg-primary-50 text-primary-700 dark:bg-primary-950/50 dark:text-primary-200' : 'text-slate-500 hover:text-slate-900 dark:text-slate-400 dark:hover:text-white'"
              :title="page.listView"
              @click="viewMode = 'list'"
            >
              <Icon name="menu" size="sm" />
            </button>
          </div>
        </div>
      </div>
    </section>

    <section class="bg-slate-50 pb-14 pt-4 dark:bg-slate-950">
      <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div class="mb-5 flex flex-wrap items-center justify-between gap-3">
          <p class="text-sm text-slate-500 dark:text-slate-400">
            {{ page.resultPrefix }} <span class="font-semibold text-slate-900 dark:text-white">{{ filteredModels.length }}</span> {{ page.resultSuffix }}
          </p>
          <div class="flex flex-wrap gap-2">
            <span
              v-if="selectedPlatform !== 'all'"
              class="inline-flex items-center rounded-lg border border-slate-200 bg-white px-3 py-1.5 text-xs font-semibold text-slate-600 dark:border-slate-800 dark:bg-slate-900 dark:text-slate-300"
            >
              {{ selectedPlatform }}
            </span>
            <span
              v-if="selectedGroup !== 'all'"
              class="inline-flex items-center rounded-lg border border-slate-200 bg-white px-3 py-1.5 text-xs font-semibold text-slate-600 dark:border-slate-800 dark:bg-slate-900 dark:text-slate-300"
            >
              {{ page.categories[selectedGroup] }}
            </span>
          </div>
        </div>

        <div v-if="filteredModels.length === 0" class="rounded-lg border border-dashed border-slate-300 bg-white px-6 py-16 text-center dark:border-slate-800 dark:bg-slate-900">
          <Icon name="search" size="xl" class="mx-auto text-slate-300 dark:text-slate-600" />
          <h2 class="mt-4 text-lg font-semibold text-slate-950 dark:text-white">{{ page.emptyTitle }}</h2>
          <p class="mt-2 text-sm text-slate-500 dark:text-slate-400">{{ page.emptyDescription }}</p>
          <button type="button" class="mt-5 inline-flex rounded-lg bg-primary-500 px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-primary-600" @click="resetFilters">
            {{ page.reset }}
          </button>
        </div>

        <div v-else-if="viewMode === 'grid'" class="grid gap-5 md:grid-cols-2 xl:grid-cols-3">
          <article
            v-for="model in filteredModels"
            :key="model.id"
            class="model-card rounded-lg border border-slate-200 bg-white p-6 shadow-sm shadow-slate-950/[0.03] transition-all duration-200 hover:-translate-y-0.5 hover:border-primary-200 hover:shadow-md hover:shadow-primary-500/10 dark:border-slate-800 dark:bg-slate-900 dark:hover:border-primary-800"
          >
            <div class="flex items-start gap-4">
              <div class="flex h-11 w-11 shrink-0 items-center justify-center rounded-lg" :class="model.logoClass">
                <span class="text-lg font-semibold">{{ model.logo }}</span>
              </div>
              <div class="min-w-0 flex-1">
                <p class="text-sm font-semibold text-slate-500 dark:text-slate-400">{{ model.provider }}</p>
                <h2 class="mt-1 break-words text-xl font-semibold tracking-normal text-slate-950 dark:text-white">
                  {{ model.name }}
                </h2>
                <div class="mt-2 flex min-w-0 items-center gap-2">
                  <code class="truncate font-mono text-sm font-semibold text-slate-400 dark:text-slate-500">{{ model.id }}</code>
                  <button type="button" class="shrink-0 text-slate-400 transition-colors hover:text-primary-600" :title="page.copyId" @click="copyModelId(model.id)">
                    <Icon :name="copiedModelId === model.id ? 'check' : 'copy'" size="sm" />
                  </button>
                </div>
              </div>
            </div>

            <div class="mt-3 flex flex-wrap gap-1.5">
              <span
                v-for="group in model.category"
                :key="group"
                class="inline-flex items-center rounded-md bg-slate-100 px-2 py-1 text-xs font-medium text-slate-600 dark:bg-slate-800 dark:text-slate-300"
              >
                {{ page.categories[group] }}
              </span>
            </div>

            <div class="mt-4 flex flex-wrap items-center gap-2">
              <span
                v-for="capability in model.capabilities"
                :key="`${capability.from}-${capability.to}`"
                class="inline-flex items-center gap-2 rounded-md border border-slate-200 bg-white px-2 py-1 text-xs font-semibold text-slate-700 dark:border-slate-700 dark:bg-slate-950 dark:text-slate-200"
              >
                <span>{{ capability.from }}</span>
                <Icon name="arrowRight" size="xs" class="text-slate-400" />
                <span>{{ capability.to }}</span>
              </span>
            </div>

            <p class="model-description mt-6 text-sm leading-7 text-slate-600 dark:text-slate-300">
              {{ model.description }}
            </p>

            <div class="mt-6 grid grid-cols-3 gap-3 border-t border-slate-100 pt-4 dark:border-slate-800">
              <div>
                <p class="text-xs font-semibold text-slate-400 dark:text-slate-500">{{ page.maxContext }}</p>
                <p class="mt-1 text-sm font-semibold text-slate-950 dark:text-white">{{ model.context }}</p>
              </div>
              <div>
                <p class="text-xs font-semibold text-slate-400 dark:text-slate-500">{{ page.maxOutput }}</p>
                <p class="mt-1 text-sm font-semibold text-slate-950 dark:text-white">{{ model.output }}</p>
              </div>
              <div>
                <p class="text-xs font-semibold text-slate-400 dark:text-slate-500">{{ page.releaseDate }}</p>
                <p class="mt-1 text-sm font-semibold text-slate-950 dark:text-white">{{ model.releaseDate }}</p>
              </div>
            </div>

            <div class="mt-5 space-y-2">
              <div v-for="price in model.prices" :key="price.label" class="grid grid-cols-[6rem,1fr] gap-3 text-sm">
                <span class="text-slate-400 dark:text-slate-500">{{ price.label }}</span>
                <span class="text-right font-mono font-semibold text-emerald-700 dark:text-emerald-300">{{ price.value }}</span>
              </div>
            </div>
          </article>
        </div>

        <div v-else class="overflow-hidden rounded-lg border border-slate-200 bg-white shadow-sm dark:border-slate-800 dark:bg-slate-900">
          <div class="hidden grid-cols-[minmax(260px,1.4fr)_140px_160px_160px_220px] gap-4 border-b border-slate-200 bg-slate-50 px-5 py-3 text-xs font-semibold uppercase tracking-wide text-slate-500 dark:border-slate-800 dark:bg-slate-950 dark:text-slate-400 lg:grid">
            <span>{{ page.listColumns.model }}</span>
            <span>{{ page.listColumns.platform }}</span>
            <span>{{ page.listColumns.group }}</span>
            <span>{{ page.listColumns.context }}</span>
            <span>{{ page.listColumns.price }}</span>
          </div>
          <article
            v-for="model in filteredModels"
            :key="model.id"
            class="grid gap-4 border-b border-slate-100 px-5 py-5 last:border-b-0 dark:border-slate-800 lg:grid-cols-[minmax(260px,1.4fr)_140px_160px_160px_220px] lg:items-center"
          >
            <div class="flex min-w-0 items-center gap-4">
              <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg" :class="model.logoClass">
                <span class="font-semibold">{{ model.logo }}</span>
              </div>
              <div class="min-w-0">
                <h2 class="truncate text-base font-semibold text-slate-950 dark:text-white">{{ model.name }}</h2>
                <p class="mt-1 truncate font-mono text-xs text-slate-500 dark:text-slate-400">{{ model.id }}</p>
              </div>
            </div>
            <p class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ model.provider }}</p>
            <div class="flex flex-wrap gap-1.5">
              <span v-for="group in model.category.slice(0, 2)" :key="group" class="rounded-md bg-slate-100 px-2 py-1 text-xs font-medium text-slate-600 dark:bg-slate-800 dark:text-slate-300">{{ page.categories[group] }}</span>
            </div>
            <p class="text-sm text-slate-600 dark:text-slate-300">{{ model.context }} / {{ model.output }}</p>
            <div class="space-y-1 text-sm">
              <p v-for="price in model.prices.slice(0, 2)" :key="price.label" class="flex justify-between gap-3">
                <span class="text-slate-400 dark:text-slate-500">{{ price.label }}</span>
                <span class="font-mono font-semibold text-emerald-700 dark:text-emerald-300">{{ price.value }}</span>
              </p>
            </div>
          </article>
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
import { withBrandTokens } from '@/brand'

type GroupKey = 'all' | 'text' | 'reasoning' | 'vision' | 'code' | 'image' | 'audio'
type ViewMode = 'grid' | 'list'

interface ModelCapability {
  from: string
  to: string
}

interface ModelPrice {
  label: string
  value: string
}

interface ModelItem {
  provider: string
  logo: string
  logoClass: string
  name: string
  id: string
  category: GroupKey[]
  tags: string[]
  capabilities: ModelCapability[]
  description: string
  context: string
  output: string
  releaseDate: string
  prices: ModelPrice[]
}

const { locale } = useI18n()
const searchQuery = ref('')
const selectedGroup = ref<GroupKey>('all')
const selectedPlatform = ref('all')
const viewMode = ref<ViewMode>('grid')
const copiedModelId = ref('')

const messages = {
  zh: {
    metaTitle: '模型广场',
    title: '模型广场',
    description: '在一个生产级目录中查看 AI 模型价格、能力、平台和分组覆盖。筛选模型后可直接使用 OpenAI 兼容接口接入 anytoken。',
    searchPlaceholder: '搜索模型、平台、分组...',
    platformLabel: '平台',
    groupLabel: '分组选择',
    all: '全部',
    reset: '重置',
    gridView: '宫格视图',
    listView: '列表视图',
    copyId: '复制模型 ID',
    resultPrefix: '共找到',
    resultSuffix: '个模型',
    emptyTitle: '没有匹配的模型',
    emptyDescription: '换一个关键词，或重置筛选条件后再试。',
    maxContext: '最大上下文',
    maxOutput: '最大输出',
    releaseDate: '发布日期',
    categories: {
      all: '全部',
      text: '文本',
      reasoning: '推理',
      vision: '视觉',
      code: '代码',
      image: '图像',
      audio: '音频',
    },
    listColumns: {
      model: '模型',
      platform: '平台',
      group: '分组',
      context: '上下文 / 输出',
      price: '价格',
    },
    models: [
      {
        provider: 'DeepSeek',
        logo: 'D',
        logoClass: 'bg-blue-50 text-blue-600 dark:bg-blue-950/50 dark:text-blue-200',
        name: 'DeepSeek V4 Pro',
        id: 'deepseek-v4-pro',
        category: ['text', 'reasoning', 'code'],
        tags: ['文本', '推理', '代码'],
        capabilities: [{ from: 'T', to: 'T' }],
        description: 'DeepSeek V4 Pro 被描述为大规模 MoE 模型，拥有 16T 总参数与 49B 激活参数，并支持 1M token 上下文窗口，适合复杂推理与代码任务。',
        context: '1M',
        output: '384K',
        releaseDate: '2026年4月24日',
        prices: [
          { label: '输入', value: '¥12.6 / 百万 Token' },
          { label: '输出', value: '¥24.5 / 百万 Token' },
          { label: '缓存读取', value: '¥0.105 / 百万 Token' },
        ],
      },
      {
        provider: 'DeepSeek',
        logo: 'D',
        logoClass: 'bg-blue-50 text-blue-600 dark:bg-blue-950/50 dark:text-blue-200',
        name: 'DeepSeek V4 Flash',
        id: 'deepseek-v4-flash',
        category: ['text', 'reasoning', 'code'],
        tags: ['文本', '推理', '低延迟'],
        capabilities: [{ from: 'T', to: 'T' }],
        description: 'DeepSeek V4 Flash 保留 V4 系列的长上下文能力，但采用更轻量的 MoE 配置，适合高并发、低成本和常规生产流量。',
        context: '1M',
        output: '384K',
        releaseDate: '2026年4月24日',
        prices: [
          { label: '输入', value: '¥1.05 / 百万 Token' },
          { label: '输出', value: '¥2.1 / 百万 Token' },
          { label: '缓存读取', value: '¥0.021 / 百万 Token' },
        ],
      },
      {
        provider: 'Alibaba',
        logo: 'Q',
        logoClass: 'bg-violet-50 text-violet-600 dark:bg-violet-950/50 dark:text-violet-200',
        name: 'Qwen3.7 Plus',
        id: 'qwen3.7-plus',
        category: ['text', 'reasoning', 'vision'],
        tags: ['中文', '视觉', 'Agent'],
        capabilities: [{ from: 'T', to: 'T' }, { from: 'I', to: 'T' }],
        description: 'Qwen3.7 Plus 延续 Qwen 面向 Agent 的设计方向，强化文本和图片输入能力，适合中文业务、工具调用和多模态理解。',
        context: '1M',
        output: '64K',
        releaseDate: '2026年6月2日',
        prices: [
          { label: '输入', value: '¥2 / 百万 Token' },
          { label: '输出', value: '¥8 / 百万 Token' },
          { label: '缓存读取', value: '¥0.4 / 百万 Token' },
        ],
      },
      {
        provider: 'DeepSeek',
        logo: 'D',
        logoClass: 'bg-blue-50 text-blue-600 dark:bg-blue-950/50 dark:text-blue-200',
        name: 'DeepSeek R1',
        id: 'deepseek-r1',
        category: ['text', 'reasoning', 'code'],
        tags: ['推理', '数学', '代码'],
        capabilities: [{ from: 'T', to: 'T' }],
        description: 'DeepSeek R1 是 DeepSeek 的推理型模型，常被提到的特点是开放推理过程，以及在数学、逻辑和代码任务上的表现。',
        context: '128K',
        output: '32.8K',
        releaseDate: '2025年1月20日',
        prices: [
          { label: '输入', value: '¥4 / 百万 Token' },
          { label: '输出', value: '¥16 / 百万 Token' },
          { label: '缓存读取', value: '¥1 / 百万 Token' },
        ],
      },
      {
        provider: 'DeepSeek',
        logo: 'D',
        logoClass: 'bg-blue-50 text-blue-600 dark:bg-blue-950/50 dark:text-blue-200',
        name: 'DeepSeek V3',
        id: 'deepseek-v3',
        category: ['text', 'code'],
        tags: ['文本', '代码', '成本'],
        capabilities: [{ from: 'T', to: 'T' }],
        description: 'DeepSeek V3 是 V3 系列的通用 MoE 基座模型，适合聊天、代码生成、结构化输出和常规生产接口调用。',
        context: '128K',
        output: '16K',
        releaseDate: '2024年12月26日',
        prices: [
          { label: '输入', value: '¥2 / 百万 Token' },
          { label: '输出', value: '¥8 / 百万 Token' },
          { label: '缓存读取', value: '¥0.5 / 百万 Token' },
        ],
      },
      {
        provider: 'DeepSeek',
        logo: 'D',
        logoClass: 'bg-blue-50 text-blue-600 dark:bg-blue-950/50 dark:text-blue-200',
        name: 'DeepSeek V3.1 Terminus',
        id: 'deepseek-v3.1-terminus',
        category: ['text', 'reasoning', 'code'],
        tags: ['Agent', '稳定性', '代码'],
        capabilities: [{ from: 'T', to: 'T' }],
        description: 'DeepSeek V3.1 Terminus 是 V3.1 系列的升级版本，重点在语言一致性、更稳定的 Agent 行为和复杂任务执行。',
        context: '128K',
        output: '64K',
        releaseDate: '2025年9月22日',
        prices: [
          { label: '输入', value: '¥4 / 百万 Token' },
          { label: '输出', value: '¥12 / 百万 Token' },
          { label: '缓存读取', value: '¥1 / 百万 Token' },
        ],
      },
      {
        provider: 'OpenAI',
        logo: 'O',
        logoClass: 'bg-emerald-50 text-emerald-700 dark:bg-emerald-950/50 dark:text-emerald-200',
        name: 'GPT-4.1',
        id: 'gpt-4.1',
        category: ['text', 'reasoning', 'code', 'vision'],
        tags: ['文本', '视觉', '代码'],
        capabilities: [{ from: 'T', to: 'T' }, { from: 'I', to: 'T' }],
        description: 'GPT-4.1 面向通用生产任务，覆盖复杂文本、视觉理解、代码和工具调用，适合作为质量优先的默认模型。',
        context: '1M',
        output: '32K',
        releaseDate: '2025年4月14日',
        prices: [
          { label: '输入', value: '按实时渠道价' },
          { label: '输出', value: '按实时渠道价' },
          { label: '缓存读取', value: '按实时渠道价' },
        ],
      },
      {
        provider: 'Google',
        logo: 'G',
        logoClass: 'bg-amber-50 text-amber-700 dark:bg-amber-950/50 dark:text-amber-200',
        name: 'Gemini 2.5 Pro',
        id: 'gemini-2.5-pro',
        category: ['text', 'reasoning', 'vision', 'code'],
        tags: ['长上下文', '视觉', '推理'],
        capabilities: [{ from: 'T', to: 'T' }, { from: 'I', to: 'T' }],
        description: 'Gemini 2.5 Pro 适合长上下文、多模态理解和复杂推理场景，可用于文档分析、代码审查和跨模态业务流。',
        context: '1M',
        output: '64K',
        releaseDate: '2025年6月17日',
        prices: [
          { label: '输入', value: '按实时渠道价' },
          { label: '输出', value: '按实时渠道价' },
          { label: '缓存读取', value: '按实时渠道价' },
        ],
      },
      {
        provider: 'OpenAI',
        logo: 'I',
        logoClass: 'bg-cyan-50 text-cyan-700 dark:bg-cyan-950/50 dark:text-cyan-200',
        name: 'Image Generation',
        id: 'gpt-image-1',
        category: ['image', 'vision'],
        tags: ['图像', '多模态'],
        capabilities: [{ from: 'T', to: 'I' }, { from: 'I', to: 'I' }],
        description: '图像生成模型适合商品图、海报、设计素材和多轮图像编辑，可通过统一网关纳入同一套额度与审计。',
        context: '图像',
        output: '图像',
        releaseDate: '2025年4月',
        prices: [
          { label: '输入', value: '按图像规格' },
          { label: '输出', value: '按图像规格' },
          { label: '缓存读取', value: '-' },
        ],
      },
    ] satisfies ModelItem[],
  },
  en: {
    metaTitle: 'Model Marketplace',
    title: 'Model Marketplace',
    description: 'Browse AI model pricing, capabilities, platforms, and group coverage in one production catalog. Filter models and connect through the OpenAI-compatible AnyToken API.',
    searchPlaceholder: 'Search models, platforms, groups...',
    platformLabel: 'Platform',
    groupLabel: 'Group selection',
    all: 'All',
    reset: 'Reset',
    gridView: 'Grid view',
    listView: 'List view',
    copyId: 'Copy model ID',
    resultPrefix: 'Showing',
    resultSuffix: 'models',
    emptyTitle: 'No matching models',
    emptyDescription: 'Try another keyword or reset the filters.',
    maxContext: 'Max context',
    maxOutput: 'Max output',
    releaseDate: 'Released',
    categories: {
      all: 'All',
      text: 'Text',
      reasoning: 'Reasoning',
      vision: 'Vision',
      code: 'Code',
      image: 'Image',
      audio: 'Audio',
    },
    listColumns: {
      model: 'Model',
      platform: 'Platform',
      group: 'Group',
      context: 'Context / output',
      price: 'Price',
    },
    models: [
      { provider: 'DeepSeek', logo: 'D', logoClass: 'bg-blue-50 text-blue-600 dark:bg-blue-950/50 dark:text-blue-200', name: 'DeepSeek V4 Pro', id: 'deepseek-v4-pro', category: ['text', 'reasoning', 'code'], tags: ['Text', 'Reasoning', 'Code'], capabilities: [{ from: 'T', to: 'T' }], description: 'DeepSeek V4 Pro is positioned as a large-scale MoE model with a 1M token context window, suited for complex reasoning and code workflows.', context: '1M', output: '384K', releaseDate: 'Apr 24, 2026', prices: [{ label: 'Input', value: '¥12.6 / MTok' }, { label: 'Output', value: '¥24.5 / MTok' }, { label: 'Cache read', value: '¥0.105 / MTok' }] },
      { provider: 'DeepSeek', logo: 'D', logoClass: 'bg-blue-50 text-blue-600 dark:bg-blue-950/50 dark:text-blue-200', name: 'DeepSeek V4 Flash', id: 'deepseek-v4-flash', category: ['text', 'reasoning', 'code'], tags: ['Text', 'Reasoning', 'Fast'], capabilities: [{ from: 'T', to: 'T' }], description: 'DeepSeek V4 Flash keeps the long-context profile while using a lighter MoE configuration for high-throughput production traffic.', context: '1M', output: '384K', releaseDate: 'Apr 24, 2026', prices: [{ label: 'Input', value: '¥1.05 / MTok' }, { label: 'Output', value: '¥2.1 / MTok' }, { label: 'Cache read', value: '¥0.021 / MTok' }] },
      { provider: 'Alibaba', logo: 'Q', logoClass: 'bg-violet-50 text-violet-600 dark:bg-violet-950/50 dark:text-violet-200', name: 'Qwen3.7 Plus', id: 'qwen3.7-plus', category: ['text', 'reasoning', 'vision'], tags: ['Chinese', 'Vision', 'Agent'], capabilities: [{ from: 'T', to: 'T' }, { from: 'I', to: 'T' }], description: 'Qwen3.7 Plus is oriented toward agentic workflows and multimodal understanding, with strong Chinese-language business coverage.', context: '1M', output: '64K', releaseDate: 'Jun 2, 2026', prices: [{ label: 'Input', value: '¥2 / MTok' }, { label: 'Output', value: '¥8 / MTok' }, { label: 'Cache read', value: '¥0.4 / MTok' }] },
      { provider: 'DeepSeek', logo: 'D', logoClass: 'bg-blue-50 text-blue-600 dark:bg-blue-950/50 dark:text-blue-200', name: 'DeepSeek R1', id: 'deepseek-r1', category: ['text', 'reasoning', 'code'], tags: ['Reasoning', 'Math', 'Code'], capabilities: [{ from: 'T', to: 'T' }], description: 'DeepSeek R1 is a reasoning model often used for math, logic, coding, and tasks where transparent reasoning behavior matters.', context: '128K', output: '32.8K', releaseDate: 'Jan 20, 2025', prices: [{ label: 'Input', value: '¥4 / MTok' }, { label: 'Output', value: '¥16 / MTok' }, { label: 'Cache read', value: '¥1 / MTok' }] },
      { provider: 'DeepSeek', logo: 'D', logoClass: 'bg-blue-50 text-blue-600 dark:bg-blue-950/50 dark:text-blue-200', name: 'DeepSeek V3', id: 'deepseek-v3', category: ['text', 'code'], tags: ['Text', 'Code', 'Cost'], capabilities: [{ from: 'T', to: 'T' }], description: 'DeepSeek V3 is a general MoE base model suitable for chat, code generation, structured output, and cost-sensitive production calls.', context: '128K', output: '16K', releaseDate: 'Dec 26, 2024', prices: [{ label: 'Input', value: '¥2 / MTok' }, { label: 'Output', value: '¥8 / MTok' }, { label: 'Cache read', value: '¥0.5 / MTok' }] },
      { provider: 'DeepSeek', logo: 'D', logoClass: 'bg-blue-50 text-blue-600 dark:bg-blue-950/50 dark:text-blue-200', name: 'DeepSeek V3.1 Terminus', id: 'deepseek-v3.1-terminus', category: ['text', 'reasoning', 'code'], tags: ['Agent', 'Stable', 'Code'], capabilities: [{ from: 'T', to: 'T' }], description: 'DeepSeek V3.1 Terminus improves language consistency, stable agent behavior, and complex execution flows for V3.1-series workloads.', context: '128K', output: '64K', releaseDate: 'Sep 22, 2025', prices: [{ label: 'Input', value: '¥4 / MTok' }, { label: 'Output', value: '¥12 / MTok' }, { label: 'Cache read', value: '¥1 / MTok' }] },
      { provider: 'OpenAI', logo: 'O', logoClass: 'bg-emerald-50 text-emerald-700 dark:bg-emerald-950/50 dark:text-emerald-200', name: 'GPT-4.1', id: 'gpt-4.1', category: ['text', 'reasoning', 'code', 'vision'], tags: ['Text', 'Vision', 'Code'], capabilities: [{ from: 'T', to: 'T' }, { from: 'I', to: 'T' }], description: 'GPT-4.1 covers complex text, vision understanding, code, and tool calls, making it a quality-first default model choice.', context: '1M', output: '32K', releaseDate: 'Apr 14, 2025', prices: [{ label: 'Input', value: 'Live route rate' }, { label: 'Output', value: 'Live route rate' }, { label: 'Cache read', value: 'Live route rate' }] },
      { provider: 'Google', logo: 'G', logoClass: 'bg-amber-50 text-amber-700 dark:bg-amber-950/50 dark:text-amber-200', name: 'Gemini 2.5 Pro', id: 'gemini-2.5-pro', category: ['text', 'reasoning', 'vision', 'code'], tags: ['Long context', 'Vision', 'Reasoning'], capabilities: [{ from: 'T', to: 'T' }, { from: 'I', to: 'T' }], description: 'Gemini 2.5 Pro is suited for long-context, multimodal understanding, and complex reasoning across documents and code.', context: '1M', output: '64K', releaseDate: 'Jun 17, 2025', prices: [{ label: 'Input', value: 'Live route rate' }, { label: 'Output', value: 'Live route rate' }, { label: 'Cache read', value: 'Live route rate' }] },
      { provider: 'OpenAI', logo: 'I', logoClass: 'bg-cyan-50 text-cyan-700 dark:bg-cyan-950/50 dark:text-cyan-200', name: 'Image Generation', id: 'gpt-image-1', category: ['image', 'vision'], tags: ['Image', 'Multimodal'], capabilities: [{ from: 'T', to: 'I' }, { from: 'I', to: 'I' }], description: 'Image generation supports product visuals, posters, creative assets, and image editing through the same governed gateway.', context: 'Image', output: 'Image', releaseDate: 'Apr 2025', prices: [{ label: 'Input', value: 'By image size' }, { label: 'Output', value: 'By image size' }, { label: 'Cache read', value: '-' }] },
    ] satisfies ModelItem[],
  },
} as const

const page = computed(() => withBrandTokens(locale.value.startsWith('zh') ? messages.zh : messages.en))
const groupOptions = computed(() => ([
  { value: 'all' as GroupKey, label: page.value.categories.all },
  { value: 'text' as GroupKey, label: page.value.categories.text },
  { value: 'reasoning' as GroupKey, label: page.value.categories.reasoning },
  { value: 'vision' as GroupKey, label: page.value.categories.vision },
  { value: 'code' as GroupKey, label: page.value.categories.code },
  { value: 'image' as GroupKey, label: page.value.categories.image },
  { value: 'audio' as GroupKey, label: page.value.categories.audio },
]))
const modelCatalog = computed<ModelItem[]>(() => page.value.models)
const platforms = computed(() => Array.from(new Set(modelCatalog.value.map((model) => model.provider))))

const filteredModels = computed(() => {
  const query = searchQuery.value.trim().toLowerCase()
  return modelCatalog.value.filter((model) => {
    const matchesGroup = selectedGroup.value === 'all' || model.category.includes(selectedGroup.value)
    const matchesPlatform = selectedPlatform.value === 'all' || model.provider === selectedPlatform.value
    const matchesQuery = !query
      || model.name.toLowerCase().includes(query)
      || model.id.toLowerCase().includes(query)
      || model.provider.toLowerCase().includes(query)
      || model.tags.some((tag) => tag.toLowerCase().includes(query))
      || model.category.some((category) => page.value.categories[category].toLowerCase().includes(query))
      || model.description.toLowerCase().includes(query)
    return matchesGroup && matchesPlatform && matchesQuery
  })
})

function resetFilters() {
  searchQuery.value = ''
  selectedGroup.value = 'all'
  selectedPlatform.value = 'all'
}

async function copyModelId(id: string) {
  copiedModelId.value = id
  try {
    await navigator.clipboard?.writeText(id)
  } catch {
    // Clipboard access can be blocked in non-secure preview contexts.
  }
  window.setTimeout(() => {
    if (copiedModelId.value === id) {
      copiedModelId.value = ''
    }
  }, 1400)
}
</script>

<style scoped>
.model-card {
  min-height: 24.5rem;
}

.model-description {
  display: -webkit-box;
  min-height: 5.25rem;
  overflow: hidden;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 3;
}
</style>
