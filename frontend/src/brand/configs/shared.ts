import type { SharedBrandContent } from '../types'

export const sharedBrandContent: SharedBrandContent = {
  publicLayout: {
    zh: {
      home: '首页',
      models: '模型广场',
      docs: '文档',
      start: '开始使用',
      product: 'OpenAI 与 Claude',
      endpoint: '统一端点',
      footerIntro: 'AnyToken 通过一套 OpenAI 兼容 API 提供 OpenAI 与 Claude 模型调用，部分模型价格最低仅为官方 API 的 1/50。',
      lightMode: '切换到浅色模式',
      darkMode: '切换到深色模式',
    },
    en: {
      home: 'Home',
      models: 'Models',
      docs: 'Docs',
      start: 'Start',
      product: 'OpenAI and Claude',
      endpoint: 'Endpoint',
      footerIntro: 'AnyToken provides OpenAI-compatible access to OpenAI and Claude models, with selected models priced from as low as 1/50 of official API rates.',
      lightMode: 'Switch to light mode',
      darkMode: 'Switch to dark mode',
    },
  },
  home: {
    zh: {
      metaTitle: 'AnyToken OpenAI 与 Claude API｜低至官方价 1/50',
      hero: {
        badge: 'OpenAI + Claude 模型 API',
        title: '用更低成本调用 OpenAI 与 Claude',
        description: '通过一套 OpenAI 兼容 API 接入 OpenAI 与 Claude 模型。部分模型价格最低仅为官方 API 的 1/50，一个 API Key 即可开始调用。',
        primaryCta: '开始使用',
        secondaryCta: '探索模型',
      },
      heroPanel: {
        title: 'OpenAI 与 Claude 统一 API',
        status: '当前可用',
      },
      providersTitle: '目前专注支持 OpenAI 与 Claude 两大模型生态',
      modalities: {
        kicker: '核心优势',
        title: '两家主流厂商，一套更省成本的 API',
        description: 'AnyToken 当前专注 OpenAI 与 Claude 模型，减少重复接入；部分模型价格最低仅为官方 API 的 1/50。',
        items: [
          { mark: 'O', title: 'OpenAI 模型', description: '接入当前已上架的 OpenAI 模型，具体模型与实时价格以模型广场为准。' },
          { mark: 'C', title: 'Claude 模型', description: '通过同一套 API 调用当前已上架的 Claude 模型，不需要维护另一套接入配置。' },
          { mark: '1/50', title: '更低调用成本', description: '部分模型最低仅为官方 API 价格的 1/50，实际费率以模型广场实时展示为准。' },
        ],
      },
      steps: {
        kicker: '三步接入',
        title: '从注册到第一条请求',
        description: '注册账号、充值余额、创建 API Key，然后从模型广场复制实际 Model ID 即可调用。',
        items: [
          { title: '注册并登录', description: '创建 AnyToken 账号，进入控制台完成基础设置。', meta: '账号就绪' },
          { title: '充值统一余额', description: '余额可用于当前支持的 OpenAI 与 Claude 模型。', meta: '统一余额' },
          { title: '创建 API Key', description: '复制 Base URL、API Key 与模型 ID 接入现有应用。', meta: 'OpenAI 兼容' },
        ],
      },
      ops: {
        kicker: '产品能力',
        title: '把调用成本与接入复杂度一起降下来',
        description: 'AnyToken 的价值不在于堆叠厂商数量，而是用更低价格和统一协议，把 OpenAI 与 Claude 调用做得简单、透明。',
        items: [
          { number: '01', title: '部分模型低至官方 1/50', description: '先在模型广场查看实时费率，再根据预算选择模型；不同模型与分组价格以实际展示为准。', cta: '查看价格', to: '/models' },
          { number: '02', title: '明确支持 OpenAI 与 Claude', description: '当前厂商范围清晰，不用在大量尚未支持的模型中筛选；可用模型以模型广场为准。', cta: '浏览模型', to: '/models' },
          { number: '03', title: '兼容现有 OpenAI SDK', description: '多数现有项目只需替换 Base URL、API Key 和模型 ID，即可接入 AnyToken。', cta: '阅读文档', to: '/docs' },
        ],
      },
      sdk: {
        title: '保留 OpenAI SDK，只替换接入配置',
        description: '使用现有 OpenAI SDK 时，通常只需替换 Base URL、API Key 和模型 ID；OpenAI 与 Claude 模型都从同一入口调用。',
        primaryCta: '生成 API Key',
        secondaryCta: '比较价格',
        code: `import OpenAI from "openai";

const client = new OpenAI({
  baseURL: "https://api.anytoken.com/v1",
  apiKey: process.env.ANYTOKEN_API_KEY,
});

await client.chat.completions.create({
  model: "YOUR_MODEL_ID",
  messages: [{ role: "user", content: "Optimize this logic." }],
});`,
      },
      pricing: {
        kicker: '价格优势',
        title: '部分模型低至官方 API 价格的 1/50',
        description: '低价是 AnyToken 当前最直接的优势。不同模型、分组以及输入输出费率可能不同，请以模型广场和控制台实际价格为准。',
        cta: '查看全部模型',
        columns: ['模型厂商', '状态', '接入方式', '价格', '计费', '模型范围', '使用建议'],
        rows: [
          { name: 'OpenAI', modality: '已支持', input: '统一 API', output: '实时费率', cache: '按实际用量', context: '模型广场更新', policy: '按需选择' },
          { name: 'Claude', modality: '已支持', input: '统一 API', output: '实时费率', cache: '按实际用量', context: '模型广场更新', policy: '按需选择' },
        ],
      },
      testimonials: {
        kicker: '为什么选择',
        title: '围绕真实优势，不做无效堆叠',
        items: [
          { role: '价格优势', quote: '部分模型最低仅为官方 API 价格的 1/50，让测试和生产调用的预算更可控。' },
          { role: '厂商范围', quote: '当前明确支持 OpenAI 与 Claude，首页、模型广场和文档保持同一口径。' },
          { role: '接入成本', quote: '统一使用 OpenAI 兼容 API，现有项目通常只需修改接入配置。' },
        ],
      },
      faq: {
        kicker: 'FAQ',
        title: '开始使用前最常见的问题',
        items: [
          { question: 'AnyToken 当前支持哪些模型厂商？', answer: '目前支持 OpenAI 与 Claude 两家。具体可用模型会随后台配置变化，请以模型广场实时列表为准。' },
          { question: '价格真的可以低到官方的 1/50 吗？', answer: '部分模型最低可达到官方 API 价格的 1/50，并非所有模型都采用同一折扣。最终价格以模型广场和控制台实际费率为准。' },
          { question: '现有 OpenAI SDK 代码还能继续使用吗？', answer: '可以。通常替换 Base URL、API Key 和模型 ID 后，就能继续使用熟悉的请求方式。' },
          { question: '费用如何计算和管理？', answer: '先在工作台充值余额，再按实际模型和请求用量扣费。不同模型的输入、输出等费率以控制台实际配置为准。' },
        ],
      },
    },
    en: {
      metaTitle: 'AnyToken OpenAI and Claude API | From 1/50 of Official Rates',
      hero: {
        badge: 'OpenAI + Claude Model API',
        title: 'Use OpenAI and Claude at a lower cost',
        description: 'Access OpenAI and Claude models through one OpenAI-compatible API. Selected models start from as low as 1/50 of official API pricing, with one API key for integration.',
        primaryCta: 'Start',
        secondaryCta: 'Explore models',
      },
      heroPanel: {
        title: 'Unified OpenAI and Claude API',
        status: 'Available now',
      },
      providersTitle: 'Currently focused on the OpenAI and Claude model ecosystems',
      modalities: {
        kicker: 'Core advantages',
        title: 'Two leading providers through one lower-cost API',
        description: 'AnyToken currently focuses on OpenAI and Claude models, reducing duplicate integration work while offering selected models from as low as 1/50 of official API rates.',
        items: [
          { mark: 'O', title: 'OpenAI models', description: 'Access the OpenAI models currently listed on AnyToken. Availability and live pricing follow the model marketplace.' },
          { mark: 'C', title: 'Claude models', description: 'Call currently listed Claude models through the same API without maintaining a separate integration path.' },
          { mark: '1/50', title: 'Lower API cost', description: 'Selected models start from as low as 1/50 of official API pricing. Live marketplace rates are authoritative.' },
        ],
      },
      steps: {
        kicker: 'Three steps',
        title: 'From signup to your first request',
        description: 'Create an account, top up your balance, issue an API key, and copy an available model ID from the marketplace.',
        items: [
          { title: 'Create an account', description: 'Sign in to AnyToken and finish the basic console setup.', meta: 'Account ready' },
          { title: 'Top up one balance', description: 'Use one balance for the currently supported OpenAI and Claude models.', meta: 'Unified balance' },
          { title: 'Create an API key', description: 'Use the Base URL, API key, and model ID in your existing application.', meta: 'OpenAI compatible' },
        ],
      },
      ops: {
        kicker: 'Product capabilities',
        title: 'Reduce both API cost and integration complexity',
        description: 'AnyToken is focused on clear value rather than provider count: lower pricing and one consistent protocol for OpenAI and Claude access.',
        items: [
          { number: '01', title: 'Selected models from 1/50 of official rates', description: 'Review live marketplace rates before choosing a model. Pricing varies by model and group.', cta: 'View pricing', to: '/models' },
          { number: '02', title: 'Clear OpenAI and Claude support', description: 'The current provider scope is explicit, and actual model availability follows the marketplace.', cta: 'Browse models', to: '/models' },
          { number: '03', title: 'Compatible with the OpenAI SDK', description: 'Most existing projects only need a new Base URL, API key, and model ID to connect.', cta: 'Read docs', to: '/docs' },
        ],
      },
      sdk: {
        title: 'Keep the OpenAI SDK and replace only connection settings',
        description: 'Existing OpenAI SDK integrations usually only need a new Base URL, API key, and model ID. OpenAI and Claude models share the same AnyToken entry point.',
        primaryCta: 'Generate API key',
        secondaryCta: 'Compare pricing',
        code: `import OpenAI from "openai";

const client = new OpenAI({
  baseURL: "https://api.anytoken.com/v1",
  apiKey: process.env.ANYTOKEN_API_KEY,
});

await client.chat.completions.create({
  model: "YOUR_MODEL_ID",
  messages: [{ role: "user", content: "Optimize this logic." }],
});`,
      },
      pricing: {
        kicker: 'Pricing advantage',
        title: 'Selected models from 1/50 of official API pricing',
        description: 'Lower pricing is AnyToken\'s clearest current advantage. Rates vary by model, group, input, and output, so the marketplace and console remain authoritative.',
        cta: 'View all models',
        columns: ['Provider', 'Status', 'Access', 'Price', 'Billing', 'Model range', 'Guidance'],
        rows: [
          { name: 'OpenAI', modality: 'Supported', input: 'Unified API', output: 'Live rates', cache: 'Usage based', context: 'Marketplace updates', policy: 'Choose as needed' },
          { name: 'Claude', modality: 'Supported', input: 'Unified API', output: 'Live rates', cache: 'Usage based', context: 'Marketplace updates', policy: 'Choose as needed' },
        ],
      },
      testimonials: {
        kicker: 'Why AnyToken',
        title: 'Focused on real advantages, not inflated provider counts',
        items: [
          { role: 'Pricing', quote: 'Selected models start from as low as 1/50 of official API rates, making test and production budgets easier to control.' },
          { role: 'Provider scope', quote: 'Current support is explicit: OpenAI and Claude, with the homepage, marketplace, and docs using the same scope.' },
          { role: 'Integration effort', quote: 'One OpenAI-compatible API means existing projects usually only need connection-setting changes.' },
        ],
      },
      faq: {
        kicker: 'FAQ',
        title: 'Common questions before getting started',
        items: [
          { question: 'Which model providers does AnyToken support today?', answer: 'AnyToken currently supports OpenAI and Claude. The live model marketplace is the source of truth for specific available models.' },
          { question: 'Can pricing really be as low as 1/50 of official rates?', answer: 'Selected models can reach 1/50 of official API pricing, but the same discount does not apply to every model. Final marketplace and console rates are authoritative.' },
          { question: 'Can we keep existing OpenAI SDK code?', answer: 'Yes. In most cases, replace the Base URL, API key, and model ID while keeping the familiar request flow.' },
          { question: 'How are costs calculated and managed?', answer: 'Top up your workspace balance, then usage is deducted according to the actual model and request rates configured in the console.' },
        ],
      },
    },
  },
}
