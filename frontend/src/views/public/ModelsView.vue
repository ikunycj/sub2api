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

        <div class="mt-8 grid gap-3 lg:grid-cols-[minmax(280px,1fr)_minmax(210px,260px)_auto_auto] lg:items-center">
          <label class="relative block">
            <Icon name="search" size="md" class="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-slate-400" />
            <input
              v-model="searchQuery"
              type="search"
              :placeholder="page.searchPlaceholder"
              class="h-12 w-full rounded-lg border border-slate-200 bg-white pl-12 pr-4 text-sm text-slate-900 outline-none transition-all duration-200 placeholder:text-slate-400 hover:border-slate-300 focus:border-primary-500 focus:ring-4 focus:ring-primary-500/10 dark:border-slate-800 dark:bg-slate-900 dark:text-white dark:hover:border-slate-700"
            />
          </label>

          <label class="relative flex h-12 items-center rounded-lg border border-slate-200 bg-white pl-4 pr-10 transition-all duration-200 hover:border-slate-300 focus-within:border-primary-500 focus-within:ring-4 focus-within:ring-primary-500/10 dark:border-slate-800 dark:bg-slate-900 dark:hover:border-slate-700">
            <span class="shrink-0 whitespace-nowrap text-sm font-medium text-slate-500 dark:text-slate-400">
              {{ page.groupLabel }}
            </span>
            <select
              v-model="selectedGroup"
              :disabled="catalogLoading"
              class="h-full min-w-0 flex-1 appearance-none bg-transparent pl-4 text-sm font-medium text-slate-800 outline-none disabled:cursor-wait disabled:text-slate-400 dark:text-slate-100 dark:disabled:text-slate-500"
            >
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
              v-if="hasOfficialPriceComparisons"
              class="inline-flex items-center rounded-lg border border-slate-200 bg-white px-3 py-1.5 text-xs text-slate-500 dark:border-slate-800 dark:bg-slate-900 dark:text-slate-400"
            >
              {{ page.officialPriceLegend }}
            </span>
            <span
              v-if="selectedGroup !== 'all'"
              class="inline-flex items-center rounded-lg border border-slate-200 bg-white px-3 py-1.5 text-xs font-semibold text-slate-600 dark:border-slate-800 dark:bg-slate-900 dark:text-slate-300"
            >
              {{ selectedGroupLabel }}
            </span>
          </div>
        </div>

        <div v-if="catalogLoading" class="overflow-hidden rounded-lg border border-slate-200 bg-white shadow-sm dark:border-slate-800 dark:bg-slate-900">
          <div class="overflow-x-auto">
            <div class="min-w-[1170px]">
              <div class="grid grid-cols-[minmax(220px,1.5fr)_90px_140px_76px_130px_130px_90px_130px] gap-4 border-b border-slate-200 bg-slate-50 px-5 py-3 text-xs font-semibold uppercase tracking-wide text-slate-500 dark:border-slate-800 dark:bg-slate-950 dark:text-slate-400">
                <span v-for="column in listColumnLabels" :key="column">{{ column }}</span>
              </div>
              <article
                v-for="index in 6"
                :key="index"
                class="grid grid-cols-[minmax(220px,1.5fr)_90px_140px_76px_130px_130px_90px_130px] gap-4 border-b border-slate-100 px-5 py-5 last:border-b-0 dark:border-slate-800"
              >
                <div class="flex animate-pulse items-center gap-4">
                  <div class="h-10 w-10 rounded-lg bg-slate-200 dark:bg-slate-800" />
                  <div class="min-w-0 flex-1 space-y-2">
                    <div class="h-4 w-2/3 rounded bg-slate-200 dark:bg-slate-800" />
                    <div class="h-3 w-1/2 rounded bg-slate-200 dark:bg-slate-800" />
                  </div>
                </div>
                <div v-for="cell in 7" :key="cell" class="animate-pulse">
                  <div class="h-4 w-3/4 rounded bg-slate-200 dark:bg-slate-800" />
                </div>
              </article>
            </div>
          </div>
        </div>

        <div v-else-if="catalogLoadFailed" class="rounded-lg border border-dashed border-slate-300 bg-white px-6 py-16 text-center dark:border-slate-800 dark:bg-slate-900">
          <Icon name="server" size="xl" class="mx-auto text-slate-300 dark:text-slate-600" />
          <h2 class="mt-4 text-lg font-semibold text-slate-950 dark:text-white">{{ page.loadFailedTitle }}</h2>
          <p class="mt-2 text-sm text-slate-500 dark:text-slate-400">{{ page.loadFailedDescription }}</p>
        </div>

        <div v-else-if="filteredModels.length === 0" class="rounded-lg border border-dashed border-slate-300 bg-white px-6 py-16 text-center dark:border-slate-800 dark:bg-slate-900">
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
            :key="model.key"
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
                v-for="group in model.groups"
                :key="group.id"
                class="inline-flex items-center rounded-md bg-slate-100 px-2 py-1 text-xs font-medium text-slate-600 dark:bg-slate-800 dark:text-slate-300"
              >
                {{ group.name }} · {{ formatRateMultiplier(group.rateMultiplier) }}
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

            <div v-if="model.prices.length > 0" class="mt-5 space-y-2 border-t border-slate-100 pt-4 dark:border-slate-800">
              <div v-for="price in model.prices" :key="price.label" class="grid grid-cols-[6rem,1fr] gap-3 text-sm">
                <span class="text-slate-400 dark:text-slate-500">{{ price.label }}</span>
                <div class="min-w-0 text-right">
                  <p class="font-mono font-semibold text-emerald-700 dark:text-emerald-300">{{ formatPriceDisplay(price) }}</p>
                  <p v-if="hasComparableOfficialPrice(price)" class="mt-1 flex flex-wrap items-center justify-end gap-x-1 text-[11px] leading-4 text-slate-400 dark:text-slate-500">
                    <span>{{ page.officialReference }}</span>
                    <span class="font-mono line-through">{{ formatOfficialPriceDisplay(price) }}</span>
                    <span v-if="hasPriceSavings(price)" class="font-medium text-emerald-600 dark:text-emerald-400">{{ formatSavingsDisplay(price) }}</span>
                  </p>
                </div>
              </div>
            </div>
          </article>
        </div>

        <div v-else class="overflow-hidden rounded-lg border border-slate-200 bg-white shadow-sm dark:border-slate-800 dark:bg-slate-900">
          <div class="overflow-x-auto">
            <div class="min-w-[1170px]">
              <div class="grid grid-cols-[minmax(220px,1.5fr)_90px_140px_76px_130px_130px_90px_130px] gap-4 border-b border-slate-200 bg-slate-50 px-5 py-3 text-xs font-semibold uppercase tracking-wide text-slate-500 dark:border-slate-800 dark:bg-slate-950 dark:text-slate-400">
                <span>{{ page.listColumns.model }}</span>
                <span>{{ page.listColumns.platform }}</span>
                <span>{{ page.listColumns.group }}</span>
                <span>{{ page.listColumns.multiplier }}</span>
                <span>{{ page.listColumns.input }}</span>
                <span>{{ page.listColumns.output }}</span>
                <span>{{ page.listColumns.modality }}</span>
                <span>{{ page.listColumns.cacheRead }}</span>
              </div>
              <article
                v-for="model in filteredModels"
                :key="model.key"
                class="grid grid-cols-[minmax(220px,1.5fr)_90px_140px_76px_130px_130px_90px_130px] gap-4 border-b border-slate-100 px-5 py-5 last:border-b-0 dark:border-slate-800"
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
                <p class="truncate text-sm font-medium text-slate-700 dark:text-slate-200">{{ model.provider }}</p>
                <div class="flex min-w-0 flex-wrap gap-1.5">
                  <span v-for="group in model.groups" :key="group.id" class="max-w-full truncate rounded-md bg-slate-100 px-2 py-1 text-xs font-medium text-slate-600 dark:bg-slate-800 dark:text-slate-300">{{ group.name }}</span>
                </div>
                <p class="font-mono text-sm font-semibold text-slate-700 dark:text-slate-200">{{ formatModelRateMultipliers(model) }}</p>
                <div class="min-w-0">
                  <p class="font-mono text-sm font-semibold text-emerald-700 dark:text-emerald-300">{{ formatListPriceCell(model, 'input') }}</p>
                  <p v-if="hasComparableListOfficialPrice(model, 'input')" class="mt-1 flex flex-wrap items-center gap-x-1 text-[11px] leading-4 text-slate-400 dark:text-slate-500">
                    <span>{{ page.officialReference }}</span>
                    <span class="font-mono line-through">{{ formatListOfficialPriceCell(model, 'input') }}</span>
                    <span v-if="hasListPriceSavings(model, 'input')" class="font-medium text-emerald-600 dark:text-emerald-400">{{ formatListSavingsCell(model, 'input') }}</span>
                  </p>
                </div>
                <div class="min-w-0">
                  <p class="font-mono text-sm font-semibold text-emerald-700 dark:text-emerald-300">{{ formatListPriceCell(model, 'output') }}</p>
                  <p v-if="hasComparableListOfficialPrice(model, 'output')" class="mt-1 flex flex-wrap items-center gap-x-1 text-[11px] leading-4 text-slate-400 dark:text-slate-500">
                    <span>{{ page.officialReference }}</span>
                    <span class="font-mono line-through">{{ formatListOfficialPriceCell(model, 'output') }}</span>
                    <span v-if="hasListPriceSavings(model, 'output')" class="font-medium text-emerald-600 dark:text-emerald-400">{{ formatListSavingsCell(model, 'output') }}</span>
                  </p>
                </div>
                <p class="text-sm text-slate-600 dark:text-slate-300">{{ formatModalityCell(model) }}</p>
                <div class="min-w-0">
                  <p class="font-mono text-sm font-semibold text-emerald-700 dark:text-emerald-300">{{ formatListPriceCell(model, 'cacheRead') }}</p>
                  <p v-if="hasComparableListOfficialPrice(model, 'cacheRead')" class="mt-1 flex flex-wrap items-center gap-x-1 text-[11px] leading-4 text-slate-400 dark:text-slate-500">
                    <span>{{ page.officialReference }}</span>
                    <span class="font-mono line-through">{{ formatListOfficialPriceCell(model, 'cacheRead') }}</span>
                    <span v-if="hasListPriceSavings(model, 'cacheRead')" class="font-medium text-emerald-600 dark:text-emerald-400">{{ formatListSavingsCell(model, 'cacheRead') }}</span>
                  </p>
                </div>
              </article>
            </div>
          </div>
        </div>
      </div>
    </section>
  </CloseAiPublicLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import CloseAiPublicLayout from '@/components/public/CloseAiPublicLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { withBrandTokens } from '@/brand'
import { publicModelsAPI, type PublicCatalogGroup, type PublicCatalogModel, type PublicModelCatalog, type PublicModelPricing } from '@/api/publicModels'
import { BILLING_MODE_IMAGE, BILLING_MODE_PER_REQUEST, BILLING_MODE_TOKEN } from '@/constants/channel'
import { createCatalogModelKey, selectCanonicalCatalogModels } from './modelCatalogNormalization'
import { calculateSavingsPercent, hasOfficialReferencePrice, isDiscountedPrice } from './modelPriceComparison'
import {
  BALANCE_USD_PER_CNY,
  calculateModelRechargePriceCny,
  calculateOfficialReferencePriceCny,
  formatModelCnyAmount,
  formatModelRateMultiplier,
  OFFICIAL_USD_TO_CNY_REFERENCE_RATE,
} from './modelPriceCurrency'

type CategoryKey = 'text' | 'reasoning' | 'vision' | 'code' | 'image' | 'audio'
type GroupFilterValue = 'all' | string
type ViewMode = 'grid' | 'list'
type CatalogSource = 'backend' | 'fallback'
type PriceCellKind = 'input' | 'output' | 'cacheRead'
type ModelPriceKind = PriceCellKind | 'perRequest'

interface ModelCapability {
  from: string
  to: string
}

interface ModelPrice {
  label: string
  value: string
  kind?: ModelPriceKind
  amountCny?: number | null
  officialAmountCny?: number | null
  unit?: string
}

interface ModelGroupTag {
  id: string
  name: string
  platform: string
  rateMultiplier: number
}

interface StaticModelItem {
  provider: string
  logo: string
  logoClass: string
  name: string
  id: string
  category: CategoryKey[]
  tags: string[]
  capabilities: ModelCapability[]
  description: string
  context: string
  output: string
  releaseDate: string
  prices: ModelPrice[]
}

interface ModelItem extends StaticModelItem {
  key: string
  groups: ModelGroupTag[]
  source: CatalogSource
}

interface CachedPublicModelCatalog {
  cachedAt: number
  catalog: PublicModelCatalog
}

type CachedPublicModelCatalogRead = PublicModelCatalog | null

const PUBLIC_MODEL_CATALOG_CACHE_KEY = 'public-model-catalog:v5'
const PUBLIC_MODEL_CATALOG_CACHE_MAX_STALE_MS = 7 * 24 * 60 * 60 * 1000
const EMPTY_CELL = '-'

const { locale } = useI18n()
const searchQuery = ref('')
const selectedGroup = ref<GroupFilterValue>('all')
const viewMode = ref<ViewMode>('list')
const copiedModelId = ref('')
const backendCatalogGroups = ref<PublicCatalogGroup[]>([])
const catalogLoaded = ref(false)
const catalogLoadFailed = ref(false)

const messages = {
  zh: {
    metaTitle: '模型广场',
    title: '模型广场',
    description: '在一个生产级目录中查看按人民币计算的 AI 模型价格、能力、平台和分组覆盖。筛选模型后可直接使用 OpenAI 兼容接口接入 anytoken。',
    searchPlaceholder: '搜索模型、分组...',
    groupLabel: '分组选择',
    reset: '重置',
    gridView: '宫格视图',
    listView: '列表视图',
    copyId: '复制模型 ID',
    officialReference: '官方',
    officialPriceLegend: `本站价按 1 元人民币 = ${BALANCE_USD_PER_CNY} 美元额度计算；删除线为按 ${OFFICIAL_USD_TO_CNY_REFERENCE_RATE} 参考汇率折算的官方人民币价`,
    savePrefix: '省',
    resultPrefix: '共找到',
    resultSuffix: '个模型',
    emptyTitle: '没有匹配的模型',
    emptyDescription: '换一个关键词，或重置筛选条件后再试。',
    loadFailedTitle: '模型目录暂时不可用',
    loadFailedDescription: '真实模型数据加载失败，请稍后刷新页面。',
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
      multiplier: '倍率',
      input: '输入',
      output: '输出',
      modality: '模态',
      cacheRead: '缓存读取',
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
    ] satisfies StaticModelItem[],
  },
  en: {
    metaTitle: 'Model Marketplace',
    title: 'Model Marketplace',
    description: 'Browse RMB-denominated AI model pricing, capabilities, platforms, and group coverage in one production catalog. Filter models and connect through the OpenAI-compatible AnyToken API.',
    searchPlaceholder: 'Search models and groups...',
    groupLabel: 'Group selection',
    reset: 'Reset',
    gridView: 'Grid view',
    listView: 'List view',
    copyId: 'Copy model ID',
    officialReference: 'Official',
    officialPriceLegend: `Our price uses CNY 1 = USD ${BALANCE_USD_PER_CNY} balance; strikethrough is the official RMB price converted at ${OFFICIAL_USD_TO_CNY_REFERENCE_RATE} CNY/USD`,
    savePrefix: 'Save',
    resultPrefix: 'Showing',
    resultSuffix: 'models',
    emptyTitle: 'No matching models',
    emptyDescription: 'Try another keyword or reset the filters.',
    loadFailedTitle: 'Model catalog unavailable',
    loadFailedDescription: 'Live model data could not be loaded. Refresh the page and try again.',
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
      multiplier: 'Multiplier',
      input: 'Input',
      output: 'Output',
      modality: 'Modality',
      cacheRead: 'Cache read',
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
    ] satisfies StaticModelItem[],
  },
} as const

const page = computed(() => withBrandTokens(locale.value.startsWith('zh') ? messages.zh : messages.en))
const isZh = computed(() => locale.value.startsWith('zh'))
const catalogLoading = computed(() => !catalogLoaded.value)
const listColumnLabels = computed(() => [
  page.value.listColumns.model,
  page.value.listColumns.platform,
  page.value.listColumns.group,
  page.value.listColumns.multiplier,
  page.value.listColumns.input,
  page.value.listColumns.output,
  page.value.listColumns.modality,
  page.value.listColumns.cacheRead,
])
const modelCatalog = computed<ModelItem[]>(() => {
  if (!catalogLoaded.value || catalogLoadFailed.value) {
    return []
  }
  return buildBackendModelCatalog(backendCatalogGroups.value)
})
const groupOptions = computed(() => {
  const seen = new Map<string, string>()
  for (const model of modelCatalog.value) {
    for (const group of model.groups) {
      if (!seen.has(group.id)) {
        seen.set(group.id, group.name)
      }
    }
  }
  return [
    { value: 'all', label: page.value.categories.all },
    ...Array.from(seen, ([value, label]) => ({ value, label })),
  ]
})
const selectedGroupLabel = computed(() => groupOptions.value.find((group) => group.value === selectedGroup.value)?.label || '')

const filteredModels = computed(() => {
  const query = searchQuery.value.trim().toLowerCase()
  return modelCatalog.value.filter((model) => {
    const matchesGroup = selectedGroup.value === 'all' || model.groups.some((group) => group.id === selectedGroup.value)
    const matchesQuery = !query
      || model.name.toLowerCase().includes(query)
      || model.id.toLowerCase().includes(query)
      || model.provider.toLowerCase().includes(query)
      || model.tags.some((tag) => tag.toLowerCase().includes(query))
      || model.category.some((category) => page.value.categories[category].toLowerCase().includes(query))
      || model.groups.some((group) => group.name.toLowerCase().includes(query))
      || model.description.toLowerCase().includes(query)
    return matchesGroup && matchesQuery
  })
})
const hasOfficialPriceComparisons = computed(() => filteredModels.value.some((model) => (
  model.prices.some((price) => hasComparableOfficialPrice(price))
)))

onMounted(() => {
  void initPublicModelCatalog()
})

async function initPublicModelCatalog() {
  const cached = readCachedPublicModelCatalog()
  if (cached) {
    applyPublicModelCatalog(cached)
    catalogLoadFailed.value = false
    catalogLoaded.value = true
    normalizeSelectedGroup()
  }

  await loadPublicModelCatalog({ preserveExistingCatalog: Boolean(cached) })
}

async function loadPublicModelCatalog(options: { preserveExistingCatalog?: boolean } = {}) {
  try {
    const catalog = await publicModelsAPI.getPublicModelCatalog()
    applyPublicModelCatalog(catalog)
    writeCachedPublicModelCatalog(catalog)
    catalogLoadFailed.value = false
  } catch {
    if (!options.preserveExistingCatalog) {
      catalogLoadFailed.value = true
    }
  } finally {
    catalogLoaded.value = true
    normalizeSelectedGroup()
  }
}

function applyPublicModelCatalog(catalog: PublicModelCatalog) {
  backendCatalogGroups.value = Array.isArray(catalog.groups) ? catalog.groups : []
}

function normalizeSelectedGroup() {
  if (selectedGroup.value !== 'all' && !groupOptions.value.some((group) => group.value === selectedGroup.value)) {
    selectedGroup.value = 'all'
  }
}

function readCachedPublicModelCatalog(): CachedPublicModelCatalogRead {
  try {
    const raw = window.localStorage.getItem(PUBLIC_MODEL_CATALOG_CACHE_KEY)
    if (!raw) {
      return null
    }
    const parsed = JSON.parse(raw) as Partial<CachedPublicModelCatalog>
    const cachedAt = Number(parsed.cachedAt)
    if (!Number.isFinite(cachedAt)) {
      return null
    }
    const age = Date.now() - cachedAt
    if (age < 0 || age > PUBLIC_MODEL_CATALOG_CACHE_MAX_STALE_MS || !parsed.catalog || !Array.isArray(parsed.catalog.groups)) {
      return null
    }
    const supportsOfficialPricing = parsed.catalog.groups.every((group) => (
      Array.isArray(group.models)
      && group.models.every((model) => Object.prototype.hasOwnProperty.call(model, 'official_pricing'))
    ))
    return supportsOfficialPricing ? parsed.catalog : null
  } catch {
    return null
  }
}

function writeCachedPublicModelCatalog(catalog: PublicModelCatalog) {
  try {
    window.localStorage.setItem(PUBLIC_MODEL_CATALOG_CACHE_KEY, JSON.stringify({
      cachedAt: Date.now(),
      catalog,
    } satisfies CachedPublicModelCatalog))
  } catch {
    // localStorage can be unavailable in private mode or full quota; the page still works without cache.
  }
}

function buildBackendModelCatalog(groups: readonly PublicCatalogGroup[]): ModelItem[] {
  const items: ModelItem[] = []
  const seen = new Set<string>()

  for (const group of groups) {
    const groupTag = toModelGroupTag(group)
    for (const model of selectCanonicalCatalogModels(group.models || [])) {
      const modelName = model.name.trim()
      if (!modelName) {
        continue
      }
      const key = `${group.id}:${(model.platform || group.platform).toLowerCase()}:${modelName.toLowerCase()}`
      if (seen.has(key)) {
        continue
      }
      seen.add(key)
      items.push(createBackendModelItem(model, groupTag))
    }
  }

  return items.sort((a, b) => a.name.localeCompare(b.name))
}

function toModelGroupTag(group: PublicCatalogGroup): ModelGroupTag {
  return {
    id: String(group.id),
    name: group.name,
    platform: group.platform,
    rateMultiplier: Number.isFinite(group.rate_multiplier) ? group.rate_multiplier : 1,
  }
}

function createBackendModelItem(model: PublicCatalogModel, group: ModelGroupTag): ModelItem {
  const platform = model.platform || group.platform
  const provider = formatPlatformLabel(platform)
  const category = inferModelCategories(model)
  const billingTag = formatBillingModeLabel(model.pricing?.billing_mode)
  return {
    key: createCatalogModelKey(group.id, platform, model.name),
    provider,
    logo: provider.slice(0, 1).toUpperCase(),
    logoClass: platformLogoClass(platform),
    name: model.name,
    id: model.name,
    category,
    groups: [group],
    tags: [provider, group.name, billingTag, ...category.map((item) => page.value.categories[item])],
    capabilities: inferModelCapabilities(category),
    description: backendModelDescription(model.name, provider, group.name),
    context: '',
    output: '',
    releaseDate: '',
    prices: formatModelPrices(model.pricing, model.official_pricing, group.rateMultiplier),
    source: 'backend',
  }
}

function inferModelCategories(model: PublicCatalogModel): CategoryKey[] {
  const value = `${model.name} ${model.platform}`.toLowerCase()
  const categories: CategoryKey[] = []
  if (/(image|dall-e|gpt-image|stable|flux|midjourney|图像)/.test(value) || model.pricing?.billing_mode === BILLING_MODE_IMAGE) {
    categories.push('image', 'vision')
  } else {
    categories.push('text')
  }
  if (/(reason|r1|o1|o3|o4|thinking|推理)/.test(value)) {
    categories.push('reasoning')
  }
  if (/(vision|gpt-4o|gemini|claude|multimodal|视觉)/.test(value)) {
    categories.push('vision')
  }
  if (/(code|coder|codex|代码)/.test(value)) {
    categories.push('code')
  }
  return Array.from(new Set(categories))
}

function inferModelCapabilities(categories: CategoryKey[]): ModelCapability[] {
  if (categories.includes('image')) {
    return [{ from: 'T', to: 'I' }, { from: 'I', to: 'I' }]
  }
  if (categories.includes('vision')) {
    return [{ from: 'T', to: 'T' }, { from: 'I', to: 'T' }]
  }
  return [{ from: 'T', to: 'T' }]
}

function backendModelDescription(modelName: string, provider: string, groupName: string): string {
  if (isZh.value) {
    return `${modelName} 来自 ${provider} 平台，属于「${groupName}」标准余额分组。可用性和价格会跟随后台渠道配置更新。`
  }
  return `${modelName} is available on ${provider} through the ${groupName} standard balance group. Availability and pricing follow the live channel configuration.`
}

function formatModelPrices(
  pricing: PublicModelPricing | null,
  officialPricing: PublicModelPricing | null | undefined,
  rateMultiplier: number,
): ModelPrice[] {
  if (!pricing) {
    return []
  }
  const comparableOfficialPricing = pricing.billing_mode === officialPricing?.billing_mode
    ? officialPricing
    : null
  if (pricing.billing_mode === BILLING_MODE_PER_REQUEST) {
    return [createPrice(
      isZh.value ? '按次' : 'Per request',
      pricing.per_request_price,
      comparableOfficialPricing?.per_request_price ?? null,
      rateMultiplier,
      1,
      'request',
      'perRequest',
    )]
  }
  if (pricing.billing_mode === BILLING_MODE_IMAGE) {
    return [
      createPrice(isZh.value ? '输入' : 'Input', pricing.input_price, comparableOfficialPricing?.input_price ?? null, rateMultiplier, 1_000_000, 'MTok', 'input'),
      createPrice(
        isZh.value ? '图像输出' : 'Image output',
        pricing.image_output_price ?? pricing.per_request_price,
        comparableOfficialPricing?.image_output_price ?? comparableOfficialPricing?.per_request_price ?? null,
        rateMultiplier,
        1,
        'request',
        'output',
      ),
    ]
  }
  return [
    createPrice(isZh.value ? '输入' : 'Input', pricing.input_price, comparableOfficialPricing?.input_price ?? null, rateMultiplier, 1_000_000, 'MTok', 'input'),
    createPrice(isZh.value ? '输出' : 'Output', pricing.output_price, comparableOfficialPricing?.output_price ?? null, rateMultiplier, 1_000_000, 'MTok', 'output'),
    createPrice(isZh.value ? '缓存读取' : 'Cache read', pricing.cache_read_price, comparableOfficialPricing?.cache_read_price ?? null, rateMultiplier, 1_000_000, 'MTok', 'cacheRead'),
  ]
}

function formatListPriceCell(model: ModelItem, kind: PriceCellKind): string {
  const price = findModelPrice(model, kind)
  return price ? formatPriceDisplay(price) : EMPTY_CELL
}

function hasComparableListOfficialPrice(model: ModelItem, kind: PriceCellKind): boolean {
  const price = findModelPrice(model, kind)
  return price ? hasComparableOfficialPrice(price) : false
}

function hasListPriceSavings(model: ModelItem, kind: PriceCellKind): boolean {
  const price = findModelPrice(model, kind)
  return price ? hasPriceSavings(price) : false
}

function formatListOfficialPriceCell(model: ModelItem, kind: PriceCellKind): string {
  const price = findModelPrice(model, kind)
  return price ? formatOfficialPriceDisplay(price) : EMPTY_CELL
}

function formatListSavingsCell(model: ModelItem, kind: PriceCellKind): string {
  const price = findModelPrice(model, kind)
  return price ? formatSavingsDisplay(price) : ''
}

function findModelPrice(model: ModelItem, kind: PriceCellKind): ModelPrice | undefined {
  return model.prices.find((item) => item.kind === kind || (!item.kind && priceLabelMatchesKind(item.label, kind)))
}

function priceLabelMatchesKind(label: string, kind: PriceCellKind): boolean {
  const normalized = label.trim().toLowerCase()
  const aliases: Record<PriceCellKind, string[]> = {
    input: ['input', '输入'],
    output: ['output', 'image output', '输出', '图像输出'],
    cacheRead: ['cache read', '缓存读取'],
  }
  return aliases[kind].includes(normalized)
}

function formatModalityCell(model: ModelItem): string {
  return formatOptionalCell(model.category.map((category) => page.value.categories[category]).join(' / '))
}

function formatOptionalCell(value?: string | null): string {
  const text = String(value || '').trim()
  return text || EMPTY_CELL
}

function formatRateMultiplier(rateMultiplier: number): string {
  return formatModelRateMultiplier(rateMultiplier)
}

function formatModelRateMultipliers(model: ModelItem): string {
  const values = model.groups.map((group) => formatRateMultiplier(group.rateMultiplier))
  return values.length > 0 ? values.join(' / ') : EMPTY_CELL
}

function createPrice(
  label: string,
  value: number | null,
  officialValue: number | null,
  rateMultiplier: number,
  scale: number,
  unit: string,
  kind: ModelPriceKind,
): ModelPrice {
  return {
    label,
    value: EMPTY_CELL,
    kind,
    amountCny: value == null ? null : calculateModelRechargePriceCny(value, rateMultiplier, scale),
    officialAmountCny: officialValue == null ? null : calculateOfficialReferencePriceCny(officialValue, scale),
    unit,
  }
}

function formatPriceDisplay(price: ModelPrice): string {
  if (price.amountCny == null || !price.unit) {
    return formatOptionalCell(price.value)
  }
  return `${formatModelCnyAmount(price.amountCny)} / ${price.unit}`
}

function hasComparableOfficialPrice(price: ModelPrice): boolean {
  return hasOfficialReferencePrice(price.officialAmountCny)
}

function hasPriceSavings(price: ModelPrice): boolean {
  return isDiscountedPrice(price.amountCny, price.officialAmountCny)
}

function formatOfficialPriceDisplay(price: ModelPrice): string {
  if (!hasComparableOfficialPrice(price) || price.officialAmountCny == null) {
    return EMPTY_CELL
  }
  return formatModelCnyAmount(price.officialAmountCny)
}

function formatSavingsDisplay(price: ModelPrice): string {
  const percentage = calculateSavingsPercent(price.amountCny, price.officialAmountCny)
  if (percentage == null) {
    return ''
  }
  return `${page.value.savePrefix} ${percentage}%`
}

function formatBillingModeLabel(mode?: string): string {
  switch (mode) {
    case BILLING_MODE_IMAGE:
      return isZh.value ? '图像' : 'Image'
    case BILLING_MODE_PER_REQUEST:
      return isZh.value ? '按次' : 'Per request'
    case BILLING_MODE_TOKEN:
    default:
      return isZh.value ? 'Token 计费' : 'Token billing'
  }
}

function formatPlatformLabel(platform: string): string {
  const normalized = platform.toLowerCase()
  const labels: Record<string, string> = {
    openai: 'OpenAI',
    anthropic: 'Anthropic',
    gemini: 'Gemini',
    antigravity: 'Antigravity',
  }
  return labels[normalized] || platform || 'AI'
}

function platformLogoClass(platform: string): string {
  const normalized = platform.toLowerCase()
  if (normalized === 'openai') return 'bg-emerald-50 text-emerald-700 dark:bg-emerald-950/50 dark:text-emerald-200'
  if (normalized === 'anthropic') return 'bg-orange-50 text-orange-700 dark:bg-orange-950/50 dark:text-orange-200'
  if (normalized === 'gemini') return 'bg-blue-50 text-blue-700 dark:bg-blue-950/50 dark:text-blue-200'
  if (normalized === 'antigravity') return 'bg-violet-50 text-violet-700 dark:bg-violet-950/50 dark:text-violet-200'
  return 'bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-200'
}

function resetFilters() {
  searchQuery.value = ''
  selectedGroup.value = 'all'
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
  min-height: 20rem;
}

.model-description {
  display: -webkit-box;
  min-height: 5.25rem;
  overflow: hidden;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 3;
}
</style>
