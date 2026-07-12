import type { BrandConfig } from '../types'

export const anytokenBrand: BrandConfig = {
  key: 'anytoken',
  siteName: 'AnyToken',
  shortName: 'anytoken',
  domain: 'anytoken.com',
  apiBaseUrl: 'https://api.anytoken.com/v1',
  apiKeyEnvVar: 'ANYTOKEN_API_KEY',
  logo: '/brands/anytoken/logo.svg',
  favicon: '/brands/anytoken/logo.svg',
  logoInitial: 'a',
  titleSuffix: 'AI API Gateway',
  theme: {
    primary: {
      50: '#eff6ff',
      100: '#dbeafe',
      200: '#bfdbfe',
      300: '#93c5fd',
      400: '#60a5fa',
      500: '#2563eb',
      600: '#1d4ed8',
      700: '#1e40af',
      800: '#1e3a8a',
      900: '#172554',
      950: '#0f172a',
    },
  },
  publicLayout: {
    zh: {
      home: '首页',
      models: '模型广场',
      docs: '文档',
      start: '开始使用',
      product: '服务 AI 团队',
      endpoint: '统一端点',
      footerIntro: 'AnyToken 基础设施为企业团队提供统一模型目录、OpenAI 兼容网关、路由、预算和审计能力。',
      lightMode: '切换到浅色模式',
      darkMode: '切换到深色模式',
    },
    en: {
      home: 'Home',
      models: 'Models',
      docs: 'Docs',
      start: 'Start',
      product: 'For AI teams',
      endpoint: 'Endpoint',
      footerIntro: 'AnyToken infrastructure gives enterprise teams a unified model catalog, OpenAI-compatible gateway, routing, budgets, and auditability.',
      lightMode: 'Switch to light mode',
      darkMode: 'Switch to dark mode',
    },
  },
  home: {
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
        description: '在 AnyToken 模型中心评估聊天、视觉、图像、视频和语音模型，不必在多个供应商控制台之间来回切换。',
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
          { title: '创建工作台', description: '为你的团队设置 AnyToken 工作台、组织和环境。', meta: '组织 / 环境' },
          { title: '配置预算和供应商', description: '一套预算规则覆盖供应商、模型和客户环境。', meta: '$20  $100  $500' },
          { title: '签发网关 Key', description: '创建统一 API Key，让请求通过 AnyToken 自动路由。', meta: 'API_KEY 可开始集成' },
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
        description: 'AnyToken 提供由模型网关支撑的 OpenAI 兼容 API。切换模型家族时，通常只需要改配置，不必重写调用代码。',
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
          { question: 'AnyToken 和单一供应商 API 有什么不同？', answer: 'AnyToken 把模型目录、路由网关、统一预算和审计放在一起。团队可以跨供应商调用模型，不用为每家供应商单独维护集成。' },
          { question: '现有 OpenAI SDK 代码还能继续使用吗？', answer: '可以。替换 base URL 和 API Key 后，就能继续使用熟悉的聊天补全请求方式。' },
          { question: 'AnyToken 如何帮助 B2B 产品团队？', answer: '统一 Key、环境、预算和供应商访问，减少客户指定模型时的重复交付。' },
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
        description: 'Evaluate chat, vision, image, video, and audio models in the AnyToken model center without jumping across provider consoles.',
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
          { title: 'Create workspace', description: 'Set up your AnyToken workspace, organization, and environments.', meta: 'Org / environments' },
          { title: 'Configure budgets and providers', description: 'Use one budget policy across providers, models, and customer environments.', meta: '$20  $100  $500' },
          { title: 'Issue gateway key', description: 'Create one API key and route requests automatically through AnyToken.', meta: 'API_KEY ready' },
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
        description: 'AnyToken provides an OpenAI-compatible API backed by a model gateway. Switching model families usually means changing config, not application code.',
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
          { question: 'How is AnyToken different from a single provider API?', answer: 'AnyToken brings model catalog, gateway routing, unified budgets, and audit logs together so teams can call models across providers without maintaining separate integrations.' },
          { question: 'Can we keep existing OpenAI SDK code?', answer: 'Yes. Replace the base URL and API key, then keep using familiar chat completion requests.' },
          { question: 'How does AnyToken help B2B product teams?', answer: 'It unifies keys, environments, budgets, and provider access, reducing duplicate delivery work when customers require specific models.' },
          { question: 'Where is spend managed?', answer: 'The console shows model spend, budget alerts, and team cost attribution.' },
        ],
      },
    },
  },
}
