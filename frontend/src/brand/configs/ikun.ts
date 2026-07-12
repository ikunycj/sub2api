import type { BrandConfig } from '../types'

export const ikunBrand: BrandConfig = {
  key: 'ikun',
  siteName: 'ikun.love',
  shortName: 'ikun',
  domain: 'ikun.love',
  apiBaseUrl: 'https://ikun.love/v1',
  apiKeyEnvVar: 'IKUN_API_KEY',
  logo: '/brands/ikun/logo.svg',
  favicon: '/brands/ikun/logo.svg',
  logoInitial: 'i',
  titleSuffix: 'AI API Gateway',
  theme: {
    primary: {
      50: '#fff8ed',
      100: '#ffefd0',
      200: '#ffdaa1',
      300: '#ffbe6b',
      400: '#ff9832',
      500: '#f97316',
      600: '#ea580c',
      700: '#c2410c',
      800: '#9a3412',
      900: '#7c2d12',
      950: '#431407',
    },
  },
  publicLayout: {
    zh: {
      home: '首页',
      models: '模型广场',
      docs: '文档',
      start: '开始使用',
      product: '服务开发者',
      endpoint: '统一端点',
      footerIntro: 'ikun.love 基于 sub2api 提供统一模型网关、账号池调度、API Key 分发、用量统计和成本控制能力。',
      lightMode: '切换到浅色模式',
      darkMode: '切换到深色模式',
    },
    en: {
      home: 'Home',
      models: 'Models',
      docs: 'Docs',
      start: 'Start',
      product: 'For developers',
      endpoint: 'Endpoint',
      footerIntro: 'ikun.love is powered by sub2api, combining a unified model gateway, account-pool routing, API key distribution, usage tracking, and cost controls.',
      lightMode: 'Switch to light mode',
      darkMode: 'Switch to dark mode',
    },
  },
  home: {
    zh: {
      metaTitle: 'AI API 中转与账号池网关',
      hero: {
        badge: 'ikun.love sub2api 服务',
        title: '一个入口管理多模型 API 调用',
        description: '通过兼容 OpenAI 的统一网关接入 GPT、Claude、Gemini、DeepSeek 等模型，把账号池、路由、额度和用量统计放进同一个控制台。',
        primaryCta: '进入控制台',
        secondaryCta: '查看模型',
      },
      heroPanel: {
        title: '可用模型与路由预览',
        status: '在线服务',
      },
      providersTitle: '统一浏览和接入常用模型供应商',
      modalities: {
        kicker: '模型接入',
        title: '把多模型能力收进一个网关',
        description: '在 ikun.love 的模型目录里查看文本、视觉、图像、视频和语音能力，用同一套 API Key 管理不同模型路径。',
        items: [
          { mark: '文', title: '文本与代码模型', description: '把聊天、代码、推理和长上下文模型统一到 OpenAI 兼容调用方式里。' },
          { mark: '图', title: '视觉与图像能力', description: '为图像理解、图像生成和多模态任务保留统一的接入和计费入口。' },
          { mark: '音', title: '语音与音频能力', description: '把语音转写、合成和音频理解能力纳入同一套运营和审计体系。' },
        ],
      },
      steps: {
        kicker: '快速接入',
        title: '从创建 Key 到发出第一条请求',
        description: '在控制台创建分组和 API Key，选择可用模型，再把 base URL 切到统一端点即可开始调用。',
        items: [
          { title: '创建 API Key', description: '为不同用途生成独立 Key，便于限额、审计和成本归因。', meta: '用户 / 分组 / Key' },
          { title: '选择模型和路由', description: '按质量、成本、上下文和可用性选择模型路径。', meta: 'GPT  Claude  Gemini' },
          { title: '接入现有 SDK', description: '替换 base URL 和密钥，继续使用熟悉的 OpenAI SDK 调用方式。', meta: '兼容 /v1' },
        ],
      },
      ops: {
        kicker: '运营能力',
        title: '让 API 服务更容易维护',
        description: 'sub2api 的账号池、分组、额度、监控和错误记录能力，可以帮助你把模型调用从零散配置变成可运营的服务。',
        items: [
          { number: '01', title: '账号池调度', description: '把多个上游账号纳入统一池子，根据分组、平台和可用性进行调度。', cta: '查看模型', to: '/models' },
          { number: '02', title: 'Key 分发与额度', description: '为用户或项目分配独立 API Key，并通过余额、并发和订阅控制使用范围。', cta: '进入控制台', to: '/login' },
          { number: '03', title: '日志、监控与排错', description: '集中查看请求、失败原因、渠道状态和成本归因，减少人工排查成本。', cta: '阅读文档', to: '/docs' },
        ],
      },
      sdk: {
        title: '保留现有 SDK，只替换统一端点',
        description: 'ikun.love 提供 OpenAI 兼容接口。大多数已有代码只需要替换 base URL 和 API Key，即可切到统一网关。',
        primaryCta: '生成 API Key',
        secondaryCta: '查看模型',
        code: `import OpenAI from "openai";

const client = new OpenAI({
  baseURL: "https://ikun.love/v1",
  apiKey: process.env.IKUN_API_KEY,
});

await client.chat.completions.create({
  model: "gpt-4.1",
  messages: [{ role: "user", content: "Optimize this logic." }],
});`,
      },
      pricing: {
        kicker: '成本与额度',
        title: '调用前先看清模型和额度',
        description: '不同模型、账号池和分组可能对应不同费率与限制。最终价格、余额和可用模型以控制台实际配置为准。',
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
        kicker: '使用场景',
        title: '少维护几套模型接入配置',
        items: [
          { role: '个人开发者', quote: '已有 OpenAI SDK 的项目，可以直接切换到统一端点做多模型测试。' },
          { role: '服务维护者', quote: '账号池、额度和日志集中管理后，排查失败请求更直接。' },
          { role: '团队管理员', quote: '不同用户用不同 Key 和分组，成本归因和权限控制都更清楚。' },
        ],
      },
      faq: {
        kicker: 'FAQ',
        title: '接入统一网关前常见问题',
        items: [
          { question: 'ikun.love 和直接调用单个模型供应商有什么不同？', answer: 'ikun.love 把账号池、模型路由、Key 分发、额度和日志放在一起，减少为每家供应商单独维护接入的成本。' },
          { question: '现有 OpenAI SDK 代码还能继续使用吗？', answer: '可以。替换 base URL 和 API Key 后，就能继续使用熟悉的聊天补全请求方式。' },
          { question: '模型和价格在哪里看？', answer: '公开模型页提供基础说明，最终可用模型、费率和额度以登录后的控制台配置为准。' },
          { question: '后台自定义首页还能用吗？', answer: '可以。管理员设置 home_content 后会覆盖默认品牌首页，适合临时公告页或完全自定义落地页。' },
        ],
      },
    },
    en: {
      metaTitle: 'AI API Proxy and Account-Pool Gateway',
      hero: {
        badge: 'ikun.love sub2api service',
        title: 'One gateway for multi-model API calls',
        description: 'Use an OpenAI-compatible gateway for GPT, Claude, Gemini, DeepSeek, and more while managing account pools, routing, quotas, and usage in one console.',
        primaryCta: 'Open console',
        secondaryCta: 'View models',
      },
      heroPanel: {
        title: 'Model and routing preview',
        status: 'Live service',
      },
      providersTitle: 'Browse common model providers through one gateway',
      modalities: {
        kicker: 'Model access',
        title: 'Bring multi-model capabilities into one gateway',
        description: 'Use one API key and one operational console for text, vision, image, video, and audio model paths.',
        items: [
          { mark: 'T', title: 'Text and code models', description: 'Route chat, code, reasoning, and long-context models through a familiar OpenAI-compatible interface.' },
          { mark: 'I', title: 'Vision and image capabilities', description: 'Keep image understanding, generation, and multimodal workflows behind one billing and access layer.' },
          { mark: 'A', title: 'Speech and audio capabilities', description: 'Manage transcription, speech synthesis, and audio understanding with the same gateway controls.' },
        ],
      },
      steps: {
        kicker: 'Quick start',
        title: 'From API key to first request',
        description: 'Create groups and API keys in the console, choose available models, then point your SDK at the unified endpoint.',
        items: [
          { title: 'Create an API key', description: 'Issue separate keys for different users, projects, or environments.', meta: 'Users / groups / keys' },
          { title: 'Choose models and routing', description: 'Pick model paths by quality, cost, context, and availability.', meta: 'GPT  Claude  Gemini' },
          { title: 'Keep your SDK', description: 'Replace the base URL and key while keeping OpenAI-compatible request code.', meta: '/v1 compatible' },
        ],
      },
      ops: {
        kicker: 'Operations',
        title: 'Easier API service maintenance',
        description: 'sub2api account pools, groups, quotas, monitoring, and error logs turn scattered model calls into an operable service.',
        items: [
          { number: '01', title: 'Account-pool routing', description: 'Manage multiple upstream accounts and route requests by group, platform, and availability.', cta: 'View models', to: '/models' },
          { number: '02', title: 'Key distribution and quotas', description: 'Give users or projects their own API keys with balance, concurrency, and subscription controls.', cta: 'Open console', to: '/login' },
          { number: '03', title: 'Logs, monitoring, and debugging', description: 'Inspect requests, failures, channel status, and cost attribution in one place.', cta: 'Read docs', to: '/docs' },
        ],
      },
      sdk: {
        title: 'Keep your SDK and replace the endpoint',
        description: 'ikun.love exposes an OpenAI-compatible API. Most existing code only needs a new base URL and API key.',
        primaryCta: 'Generate API key',
        secondaryCta: 'View models',
        code: `import OpenAI from "openai";

const client = new OpenAI({
  baseURL: "https://ikun.love/v1",
  apiKey: process.env.IKUN_API_KEY,
});

await client.chat.completions.create({
  model: "gpt-4.1",
  messages: [{ role: "user", content: "Optimize this logic." }],
});`,
      },
      pricing: {
        kicker: 'Cost and quota',
        title: 'Understand models and limits before calling',
        description: 'Models, account pools, and groups can carry different rates and limits. Final availability and pricing follow the console configuration.',
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
        kicker: 'Use cases',
        title: 'Maintain fewer model integrations',
        items: [
          { role: 'Solo developer', quote: 'Projects that already use the OpenAI SDK can switch endpoints and test multiple models quickly.' },
          { role: 'Service operator', quote: 'Central account pools, quotas, and logs make failed request diagnosis much more direct.' },
          { role: 'Team admin', quote: 'Separate keys and groups make cost attribution and access control clearer.' },
        ],
      },
      faq: {
        kicker: 'FAQ',
        title: 'Questions before using a unified gateway',
        items: [
          { question: 'How is ikun.love different from calling one provider directly?', answer: 'ikun.love combines account pools, model routing, key distribution, quotas, and logs so you do not maintain a separate integration for every provider.' },
          { question: 'Can we keep existing OpenAI SDK code?', answer: 'Yes. Replace the base URL and API key, then keep using familiar chat completion requests.' },
          { question: 'Where do models and prices come from?', answer: 'The public model page explains the basics. Final availability, rates, and quotas follow the authenticated console configuration.' },
          { question: 'Can the admin custom home page still override this?', answer: 'Yes. When home_content is configured, it replaces the default branded home page for temporary notices or a fully custom landing page.' },
        ],
      },
    },
  },
}
