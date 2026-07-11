<template>
  <CloseAiPublicLayout :page-title="page.title">
    <div class="border-b border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-950">
      <div class="mx-auto max-w-[1440px] px-4 sm:px-6 lg:px-8">
        <div class="border-b border-slate-200 py-3 lg:hidden dark:border-slate-800">
          <button
            type="button"
            class="flex w-full items-center justify-between rounded-xl border border-slate-200 bg-white px-3 py-2.5 text-sm font-medium text-slate-800 shadow-sm transition hover:bg-slate-50 dark:border-slate-700 dark:bg-slate-900 dark:text-white dark:hover:bg-slate-800"
            @click="mobileNavigationOpen = !mobileNavigationOpen"
          >
            <span class="flex items-center gap-2"><Icon name="book" size="sm" />{{ copy.navigation }}</span>
            <Icon :name="mobileNavigationOpen ? 'chevronUp' : 'chevronDown'" size="sm" />
          </button>
          <nav v-if="mobileNavigationOpen" class="mt-2 max-h-[60vh] overflow-y-auto rounded-xl border border-slate-200 bg-white p-2 shadow-lg dark:border-slate-700 dark:bg-slate-900">
            <template v-for="group in navigation" :key="group.title">
              <p class="px-3 pb-1 pt-3 text-xs font-semibold text-slate-400">{{ group.title }}</p>
              <router-link
                v-for="item in group.items"
                :key="item.to"
                :to="item.to"
                class="flex rounded-lg px-3 py-2 text-sm transition"
                :class="isCurrentPage(item.to) ? activeNavigationClass : idleNavigationClass"
                @click="mobileNavigationOpen = false"
              >
                {{ item.label }}
              </router-link>
            </template>
          </nav>
        </div>

        <div class="grid grid-cols-1 gap-10 lg:grid-cols-[238px_minmax(0,760px)_180px] xl:gap-14">
          <aside class="hidden border-r border-slate-200 pr-5 lg:block dark:border-slate-800">
            <nav class="sticky top-20 max-h-[calc(100vh-5rem)] overflow-y-auto py-9">
              <div v-for="group in navigation" :key="group.title" class="mb-7">
                <p class="mb-2 px-3 text-xs font-semibold text-slate-400 dark:text-slate-500">{{ group.title }}</p>
                <div class="space-y-0.5">
                  <router-link
                    v-for="item in group.items"
                    :key="item.to"
                    :to="item.to"
                    class="group flex items-center justify-between rounded-xl px-3 py-2 text-sm transition-all duration-200"
                    :class="isCurrentPage(item.to) ? activeNavigationClass : idleNavigationClass"
                  >
                    <span>{{ item.label }}</span>
                    <Icon v-if="isCurrentPage(item.to)" name="chevronRight" size="xs" class="text-primary-600 dark:text-primary-300" />
                  </router-link>
                </div>
              </div>
            </nav>
          </aside>

          <article class="min-w-0 py-9 lg:py-12">
            <nav class="mb-8 flex flex-wrap items-center gap-2 text-sm text-slate-500 dark:text-slate-400" aria-label="Breadcrumb">
              <router-link to="/home" class="transition-colors hover:text-primary-700 dark:hover:text-primary-300">{{ copy.home }}</router-link>
              <Icon name="chevronRight" size="xs" class="text-slate-300 dark:text-slate-700" />
              <router-link to="/docs" class="transition-colors hover:text-primary-700 dark:hover:text-primary-300">{{ copy.docs }}</router-link>
              <template v-for="crumb in page.breadcrumbs" :key="crumb">
                <Icon name="chevronRight" size="xs" class="text-slate-300 dark:text-slate-700" />
                <span class="text-slate-700 dark:text-slate-200">{{ crumb }}</span>
              </template>
            </nav>

            <header class="border-b border-slate-200 pb-9 dark:border-slate-800">
              <p v-if="page.eyebrow" class="mb-3 text-sm font-semibold text-primary-600 dark:text-primary-300">{{ page.eyebrow }}</p>
              <h1 class="text-4xl font-semibold tracking-normal text-slate-950 sm:text-5xl dark:text-white">{{ page.title }}</h1>
              <p class="mt-5 max-w-3xl text-base leading-7 text-slate-600 sm:text-lg sm:leading-8 dark:text-slate-300">{{ page.description }}</p>
            </header>

            <div class="doc-prose pt-9">
              <p v-if="page.intro" class="doc-lead" v-html="page.intro" />

              <section v-for="section in page.sections" :id="section.id" :key="section.id" class="scroll-mt-28 pt-9 first:pt-0">
                <h2>{{ section.title }}</h2>

                <template v-for="(block, blockIndex) in section.blocks" :key="`${section.id}-${blockIndex}`">
                  <p v-if="block.type === 'paragraph'" v-html="block.content" />

                  <div v-else-if="block.type === 'code'" class="code-block">
                    <div class="code-toolbar">
                      <span>{{ block.label }}</span>
                      <button type="button" :aria-label="copy.copyCode" :title="copy.copyCode" @click="copyCode(block.content, `${section.id}-${blockIndex}`)">
                        <Icon :name="copiedCode === `${section.id}-${blockIndex}` ? 'check' : 'copy'" size="sm" />
                        <span>{{ copiedCode === `${section.id}-${blockIndex}` ? copy.copied : copy.copy }}</span>
                      </button>
                    </div>
                    <pre><code>{{ block.content }}</code></pre>
                  </div>

                  <ul v-else-if="block.type === 'list'" class="doc-list">
                    <li v-for="item in block.items" :key="item" v-html="item" />
                  </ul>

                  <ol v-else-if="block.type === 'steps'" class="doc-steps">
                    <li v-for="(item, itemIndex) in block.items" :key="item">
                      <span>{{ itemIndex + 1 }}</span>
                      <p v-html="item" />
                    </li>
                  </ol>

                  <div v-else-if="block.type === 'note'" class="doc-note">
                    <Icon name="infoCircle" size="md" />
                    <div>
                      <strong>{{ block.title }}</strong>
                      <p v-html="block.content" />
                    </div>
                  </div>

                  <div v-else-if="block.type === 'table'" class="doc-table-wrap">
                    <table>
                      <thead><tr><th v-for="heading in block.headings" :key="heading">{{ heading }}</th></tr></thead>
                      <tbody>
                        <tr v-for="row in block.rows" :key="row[0]">
                          <td v-for="cell in row" :key="cell" v-html="cell" />
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </template>
              </section>
            </div>

            <div class="mt-12 flex items-center gap-2 border-t border-slate-200 pt-6 text-sm text-slate-500 dark:border-slate-800 dark:text-slate-400">
              <Icon name="calendar" size="sm" />{{ copy.lastUpdated }}
            </div>

            <div class="mt-8 grid gap-3 sm:grid-cols-2">
              <router-link v-if="page.previous" :to="page.previous.to" class="doc-page-card">
                <Icon name="arrowLeft" size="sm" />
                <span><small>{{ copy.previous }}</small><strong>{{ page.previous.label }}</strong></span>
              </router-link>
              <div v-else />
              <router-link v-if="page.next" :to="page.next.to" class="doc-page-card justify-end text-right">
                <span><small>{{ copy.next }}</small><strong>{{ page.next.label }}</strong></span>
                <Icon name="arrowRight" size="sm" />
              </router-link>
            </div>
          </article>

          <aside class="hidden xl:block">
            <nav class="sticky top-24 py-12">
              <p class="mb-3 text-xs font-semibold text-slate-400 dark:text-slate-500">{{ copy.onThisPage }}</p>
              <div class="border-l border-slate-200 dark:border-slate-800">
                <a
                  v-for="section in page.sections"
                  :key="section.id"
                  :href="`#${section.id}`"
                  class="-ml-px block border-l px-4 py-1.5 text-sm transition-colors"
                  :class="activeSection === section.id ? 'border-primary-500 font-medium text-primary-700 dark:text-primary-300' : 'border-transparent text-slate-500 hover:text-slate-900 dark:text-slate-400 dark:hover:text-white'"
                >
                  {{ section.title }}
                </a>
              </div>
            </nav>
          </aside>
        </div>
      </div>
    </div>
  </CloseAiPublicLayout>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import CloseAiPublicLayout from '@/components/public/CloseAiPublicLayout.vue'
import Icon from '@/components/icons/Icon.vue'

type TextBlock = { type: 'paragraph'; content: string }
type CodeBlock = { type: 'code'; label: string; content: string }
type ListBlock = { type: 'list' | 'steps'; items: string[] }
type NoteBlock = { type: 'note'; title: string; content: string }
type TableBlock = { type: 'table'; headings: string[]; rows: string[][] }
type DocBlock = TextBlock | CodeBlock | ListBlock | NoteBlock | TableBlock
type DocSection = { id: string; title: string; blocks: DocBlock[] }
type DocLink = { label: string; to: string }
type DocPage = {
  title: string
  eyebrow?: string
  description: string
  intro?: string
  breadcrumbs: string[]
  sections: DocSection[]
  previous?: DocLink
  next?: DocLink
}

const route = useRoute()
const { locale } = useI18n()
const copiedCode = ref('')
const activeSection = ref('')
const mobileNavigationOpen = ref(false)
let observer: IntersectionObserver | undefined

const isZh = computed(() => locale.value.startsWith('zh'))
const copy = computed(() => isZh.value ? {
  home: '首页', docs: '文档', navigation: '文档导航', copied: '已复制', copy: '复制', copyCode: '复制代码', onThisPage: '本页目录', previous: '上一篇', next: '下一篇', lastUpdated: '最后更新于 2026 年 7 月 11 日',
} : {
  home: 'Home', docs: 'Docs', navigation: 'Documentation', copied: 'Copied', copy: 'Copy', copyCode: 'Copy code', onThisPage: 'On this page', previous: 'Previous', next: 'Next', lastUpdated: 'Last updated July 11, 2026',
})

const integrationItems = computed(() => isZh.value ? [
  ['集成指南', '/docs/integrations'], ['Claude Code', '/docs/integrations/claude-code'], ['Codex', '/docs/integrations/codex'], ['CC Switch', '/docs/integrations/cc-switch'], ['Cursor', '/docs/integrations/cursor'], ['Cline', '/docs/integrations/cline'], ['Roo Code', '/docs/integrations/roo-code'], ['Kilo Code', '/docs/integrations/kilo-code'], ['Aider', '/docs/integrations/aider'], ['Cherry Studio', '/docs/integrations/cherry-studio'], ['OpenCode', '/docs/integrations/opencode'], ['Qwen Code', '/docs/integrations/qwen-code'],
] : [
  ['Integrations', '/docs/integrations'], ['Claude Code', '/docs/integrations/claude-code'], ['Codex', '/docs/integrations/codex'], ['CC Switch', '/docs/integrations/cc-switch'], ['Cursor', '/docs/integrations/cursor'], ['Cline', '/docs/integrations/cline'], ['Roo Code', '/docs/integrations/roo-code'], ['Kilo Code', '/docs/integrations/kilo-code'], ['Aider', '/docs/integrations/aider'], ['Cherry Studio', '/docs/integrations/cherry-studio'], ['OpenCode', '/docs/integrations/opencode'], ['Qwen Code', '/docs/integrations/qwen-code'],
])

const navigation = computed(() => isZh.value ? [
  { title: 'API 参考', items: [{ label: 'API 参考', to: '/docs' }] },
  { title: 'AI 模型接口', items: [
    { label: '聊天', to: '/docs/api/chat' }, { label: 'Responses', to: '/docs/api/responses' }, { label: '嵌入', to: '/docs/api/embeddings' }, { label: '图像生成', to: '/docs/api/images' }, { label: '模型', to: '/docs/api/models' },
  ] },
  { title: '集成指南', items: integrationItems.value.map(([label, to]) => ({ label, to })) },
] : [
  { title: 'API Reference', items: [{ label: 'API Reference', to: '/docs' }] },
  { title: 'AI Model APIs', items: [
    { label: 'Chat', to: '/docs/api/chat' }, { label: 'Responses', to: '/docs/api/responses' }, { label: 'Embeddings', to: '/docs/api/embeddings' }, { label: 'Image generation', to: '/docs/api/images' }, { label: 'Models', to: '/docs/api/models' },
  ] },
  { title: 'Integrations', items: integrationItems.value.map(([label, to]) => ({ label, to })) },
])

const codexPage = computed<DocPage>(() => isZh.value ? {
  title: 'Codex',
  eyebrow: '集成指南',
  description: '安装 Codex CLI，并通过 anytoken 的 OpenAI Responses 兼容接口配置服务商、模型和验证命令。',
  intro: 'Codex 可以通过 anytoken 的 OpenAI Responses 兼容接口接入。Base URL 使用 <code>https://api.anytoken.com/v1</code>，请求由 Codex 发往 <code>/responses</code>，API Key 在 anytoken 控制台中创建。',
  breadcrumbs: ['Integrations', 'Codex'],
  sections: [
    { id: 'use-cases', title: '适合场景', blocks: [
      { type: 'paragraph', content: 'Codex 适合在终端中完成仓库理解、代码修改、命令执行、测试修复和日常开发任务。建议从 Git 工作区启动，这样每次变更都可以先通过 diff 审阅，再决定是否保留。' },
    ] },
    { id: 'install', title: '安装 Codex', blocks: [
      { type: 'paragraph', content: '使用 npm 安装 Codex CLI。安装前请确保本机已有受支持的 Node.js 与 npm 环境。' },
      { type: 'code', label: 'Terminal', content: 'npm install -g @openai/codex' },
      { type: 'paragraph', content: '安装完成后，先确认命令可以正常启动：' },
      { type: 'code', label: 'Terminal', content: 'codex --version' },
      { type: 'note', title: '团队环境', content: '在公司终端、构建机或 CI 中安装时，请把 CLI 版本固定到团队验证过的版本，并通过统一的软件分发流程升级。' },
    ] },
    { id: 'credentials', title: '准备 anytoken 凭证', blocks: [
      { type: 'paragraph', content: '进入 anytoken 控制台创建项目 API Key。密钥只应保存在环境变量或团队密钥管理系统中，不要提交到仓库。' },
      { type: 'code', label: 'macOS / Linux', content: 'export ANYTOKEN_API_KEY="ak-..."' },
      { type: 'code', label: 'PowerShell', content: '$env:ANYTOKEN_API_KEY="ak-..."' },
      { type: 'paragraph', content: '从模型广场或 <code>GET /v1/models</code> 返回值中选择完整 Model ID。本文以 <code>gpt-5.5</code> 为例；如果项目未授权该模型，请替换为当前 Key 可访问的 Codex 模型。' },
    ] },
    { id: 'provider', title: '推荐方式：配置 anytoken Provider', blocks: [
      { type: 'paragraph', content: '在 <code>~/.codex/config.toml</code> 中添加自定义 provider，让 Codex 每次启动都使用同一套网关配置：' },
      { type: 'code', label: '~/.codex/config.toml', content: 'model = "gpt-5.5"\nmodel_provider = "anytoken"\n\n[model_providers.anytoken]\nname = "anytoken"\nbase_url = "https://api.anytoken.com/v1"\nenv_key = "ANYTOKEN_API_KEY"\nwire_api = "responses"' },
      { type: 'note', title: 'Responses 接口', content: '<code>wire_api = "responses"</code> 对应 <code>https://api.anytoken.com/v1/responses</code>。请确认所选模型支持 Responses 请求格式。' },
      { type: 'paragraph', content: '<code>env_key</code> 只指定 Codex 应读取的环境变量名称，密钥本身不会写进配置文件。这样既能持久保存 provider，又能让个人电脑、CI 和服务器分别使用自己的凭证。' },
    ] },
    { id: 'environment-check', title: '环境变量验证', blocks: [
      { type: 'paragraph', content: '首次接入时，可以在当前终端临时设置密钥并启动 Codex。Base URL 与请求协议仍由 <code>config.toml</code> 决定。' },
      { type: 'code', label: 'macOS / Linux', content: 'export ANYTOKEN_API_KEY="ak-..."\ncodex --model gpt-5.5' },
      { type: 'code', label: 'PowerShell', content: '$env:ANYTOKEN_API_KEY="ak-..."\ncodex --model gpt-5.5' },
      { type: 'paragraph', content: '验证完成后，把密钥交给系统环境变量、密码管理器或企业密钥管理服务，不要长期保存在 Shell 历史中。' },
    ] },
    { id: 'commands', title: '启动和常用命令', blocks: [
      { type: 'paragraph', content: '在需要处理的项目目录中启动交互式 Codex：' },
      { type: 'code', label: 'Terminal', content: 'cd /path/to/your/repo\ncodex --model gpt-5.5' },
      { type: 'paragraph', content: '也可以直接提交一次性任务。首次验证建议使用只读问题：' },
      { type: 'code', label: 'Terminal', content: 'codex --model gpt-5.5 "说明这个仓库的主要技术栈，不要修改文件"' },
    ] },
    { id: 'verification', title: '验证', blocks: [
      { type: 'steps', items: [
        '先提交只读问题，确认 Codex 能读取项目并正常返回。',
        '在 anytoken 请求日志中确认请求命中了预期的模型与 <code>/v1/responses</code> 接口。',
        '再让 Codex 修改低风险文件，例如 README 或单元测试。',
        '使用 <code>git diff</code> 审阅变更，确认执行范围、代码质量和费用符合预期。',
      ] },
    ] },
    { id: 'troubleshooting', title: '常见问题', blocks: [
      { type: 'table', headings: ['现象', '处理方式'], rows: [
        ['401', '确认 <code>ANYTOKEN_API_KEY</code> 已设置，并且 provider 的 <code>env_key</code> 名称完全一致。'],
        ['404', '确认 Base URL 为 <code>https://api.anytoken.com/v1</code>，不要只填写根域名。'],
        ['Responses 接口错误', '确认 <code>wire_api = "responses"</code>，并选择支持 <code>/v1/responses</code> 的模型。'],
        ['模型名不识别', '使用 <code>GET /v1/models</code> 返回的完整 Model ID，并确认该模型支持 Responses API。'],
        ['命令风险过高', '先要求 Codex 输出计划和待执行命令，再逐项确认高风险操作。'],
        ['配置未生效', '检查 <code>~/.codex/config.toml</code> 的 provider 名称，并确认当前终端能够读取环境变量。'],
      ] },
    ] },
  ],
  previous: { label: 'Claude Code', to: '/docs/integrations/claude-code' },
  next: { label: 'CC Switch', to: '/docs/integrations/cc-switch' },
} : {
  title: 'Codex',
  eyebrow: 'Integration guide',
  description: 'Install Codex CLI and connect it to anytoken through an OpenAI Responses-compatible provider.',
  intro: 'Codex can use the anytoken OpenAI Responses-compatible endpoint. Set the Base URL to <code>https://api.anytoken.com/v1</code> and create the API key in the anytoken console.',
  breadcrumbs: ['Integrations', 'Codex'],
  sections: [
    { id: 'use-cases', title: 'Use cases', blocks: [{ type: 'paragraph', content: 'Codex works well for repository exploration, code changes, command execution, test repair, and everyday development from the terminal. Start it inside a Git worktree so every change can be reviewed before it is kept.' }] },
    { id: 'install', title: 'Install Codex', blocks: [{ type: 'paragraph', content: 'Install Codex CLI with npm, then verify the command.' }, { type: 'code', label: 'Terminal', content: 'npm install -g @openai/codex\ncodex --version' }] },
    { id: 'credentials', title: 'Prepare credentials', blocks: [{ type: 'paragraph', content: 'Create a project API key in the anytoken console and keep it outside your repository.' }, { type: 'code', label: 'Terminal', content: 'export ANYTOKEN_API_KEY="ak-..."' }] },
    { id: 'provider', title: 'Configure the anytoken provider', blocks: [{ type: 'paragraph', content: 'Add the provider to <code>~/.codex/config.toml</code>.' }, { type: 'code', label: '~/.codex/config.toml', content: 'model = "gpt-5.5"\nmodel_provider = "anytoken"\n\n[model_providers.anytoken]\nname = "anytoken"\nbase_url = "https://api.anytoken.com/v1"\nenv_key = "ANYTOKEN_API_KEY"\nwire_api = "responses"' }, { type: 'note', title: 'Responses API', content: 'The selected model must support requests to <code>/v1/responses</code>.' }] },
    { id: 'environment-check', title: 'Verify the environment', blocks: [{ type: 'code', label: 'Terminal', content: 'export ANYTOKEN_API_KEY="ak-..."\ncodex --model gpt-5.5' }] },
    { id: 'commands', title: 'Start Codex', blocks: [{ type: 'code', label: 'Terminal', content: 'cd /path/to/your/repo\ncodex --model gpt-5.5' }] },
    { id: 'verification', title: 'Verification', blocks: [{ type: 'steps', items: ['Ask a read-only repository question.', 'Confirm the request and model in anytoken logs.', 'Make one low-risk change.', 'Review the result with <code>git diff</code>.'] }] },
    { id: 'troubleshooting', title: 'Troubleshooting', blocks: [{ type: 'table', headings: ['Issue', 'Resolution'], rows: [['401', 'Check <code>ANYTOKEN_API_KEY</code> and <code>env_key</code>.'], ['404', 'Use <code>https://api.anytoken.com/v1</code> as the Base URL.'], ['Responses error', 'Use <code>wire_api = "responses"</code> and a compatible model.'], ['Unknown model', 'Use a complete Model ID from the model marketplace.']] }] },
  ],
  previous: { label: 'Claude Code', to: '/docs/integrations/claude-code' }, next: { label: 'CC Switch', to: '/docs/integrations/cc-switch' },
})

const ccSwitchPage = computed<DocPage>(() => isZh.value ? {
  title: 'CC Switch',
  eyebrow: '集成指南',
  description: '从 anytoken 控制台一键导入 Provider，并在 CC Switch 中统一切换 Claude Code、Codex 和 Gemini CLI 配置。',
  intro: 'CC Switch 是 Claude Code、Codex、Gemini CLI 等开发工具的本地 Provider 管理器。anytoken 控制台已经集成 <code>ccswitch://</code> 导入协议，可以把端点、API Key、应用类型和余额查询配置一次性写入 CC Switch。',
  breadcrumbs: ['Integrations', 'CC Switch'],
  sections: [
    { id: 'install', title: '安装 CC Switch', blocks: [
      { type: 'paragraph', content: '从 CC Switch 官方发布页下载适合当前系统的安装包。Windows 支持安装版和便携版，macOS 可以使用 Homebrew，Linux 提供 deb、rpm 与 AppImage。' },
      { type: 'code', label: 'macOS', content: 'brew install --cask cc-switch' },
      { type: 'paragraph', content: '首次启动后，确认操作系统已经注册 <code>ccswitch://</code> 协议。浏览器的一键导入依赖这个协议唤起桌面应用。' },
    ] },
    { id: 'create-key', title: '创建 API Key', blocks: [
      { type: 'steps', items: [
        '登录 anytoken 控制台，进入 API Key 页面。',
        '按用途选择对应平台的分组并创建 Key。Claude Code 使用 Anthropic 分组，Codex 使用 OpenAI 分组，Gemini CLI 使用 Gemini 分组。',
        '为不同项目创建独立 Key，并设置必要的模型权限、额度和速率限制。',
      ] },
      { type: 'note', title: '平台决定导入目标', content: '一键导入会读取 API Key 所属分组。OpenAI 分组导入为 Codex Provider，Anthropic 分组导入为 Claude Provider，Gemini 分组导入为 Gemini Provider。' },
    ] },
    { id: 'one-click-import', title: '一键导入', blocks: [
      { type: 'steps', items: [
        '保持 CC Switch 在本机已安装并至少启动过一次。',
        '在 anytoken 控制台的 API Key 列表中找到目标 Key。',
        '点击“导入到 CC Switch”，浏览器询问是否打开外部应用时选择允许。',
        '在 CC Switch 中检查 Provider 名称、应用类型和端点，然后确认导入。',
      ] },
      { type: 'paragraph', content: 'anytoken 会自动传递 Provider 名称、站点地址、请求端点、API Key 和余额查询脚本。导入 OpenAI 分组时还会写入 Codex 模型配置。' },
    ] },
    { id: 'provider-mapping', title: 'Provider 映射', blocks: [
      { type: 'table', headings: ['anytoken 分组', 'CC Switch 应用', '端点'], rows: [
        ['Anthropic', 'Claude', '<code>https://api.anytoken.com</code>'],
        ['OpenAI', 'Codex', '<code>https://api.anytoken.com</code>'],
        ['Gemini', 'Gemini', '<code>https://api.anytoken.com</code>'],
        ['Antigravity', 'Claude 或 Gemini', '<code>https://api.anytoken.com/antigravity</code>'],
      ] },
      { type: 'note', title: '不要重复添加 /v1', content: 'CC Switch 的 Provider 导入使用站点根地址，具体 CLI 配置和接口路径由应用类型处理。不要在自动导入结果后手工再拼接 <code>/v1</code>。' },
    ] },
    { id: 'usage', title: '余额与可用性', blocks: [
      { type: 'paragraph', content: '导入时会启用用量查询，并以 30 分钟为间隔请求 anytoken 的 <code>/v1/usage</code> 接口。CC Switch 可据此展示 Key 是否有效、剩余额度和计价单位。' },
      { type: 'code', label: 'Usage endpoint', content: 'GET https://api.anytoken.com/v1/usage\nAuthorization: Bearer <ANYTOKEN_API_KEY>' },
      { type: 'paragraph', content: '如果余额没有立即刷新，可以在 CC Switch 中手动触发更新，同时检查 Key 是否启用、是否存在余额或套餐，以及请求是否被本机代理拦截。' },
    ] },
    { id: 'switch-provider', title: '切换并生效', blocks: [
      { type: 'steps', items: [
        '在 CC Switch 中选中刚导入的 anytoken Provider。',
        '点击启用，或从系统托盘直接选择 Provider。',
        '重新打开 Codex 或 Gemini CLI 终端会话，使配置重新加载；Claude Code 通常可以直接读取切换结果。',
        '发送一个简短请求，并在 anytoken 请求日志中确认模型、API Key 和状态码。',
      ] },
    ] },
    { id: 'manual-configuration', title: '手动配置', blocks: [
      { type: 'paragraph', content: '如果浏览器无法唤起 CC Switch，可以在应用中手动新增自定义 Provider。不要复制 Deep Link 中的明文 API Key 到聊天、工单或截图。' },
      { type: 'code', label: 'Provider', content: 'Name: anytoken\nEndpoint: https://api.anytoken.com\nAPI Key: <ANYTOKEN_API_KEY>\nApplication: Claude / Codex / Gemini' },
    ] },
    { id: 'troubleshooting', title: '常见问题', blocks: [
      { type: 'table', headings: ['现象', '处理方式'], rows: [
        ['点击导入没有反应', '先启动一次 CC Switch，确认系统已注册 <code>ccswitch://</code>；仍失败时重装最新版或改用手动配置。'],
        ['导入到了错误应用', '检查 API Key 所属分组。应用类型由分组平台决定，不由按钮点击时临时选择。'],
        ['Codex 请求失败', '确认使用 OpenAI 分组，并检查所选模型支持 Responses API。'],
        ['Claude Code 请求失败', '确认使用 Anthropic 分组，且 Provider 端点没有重复附加 <code>/v1</code>。'],
        ['余额无法读取', '直接测试 <code>GET /v1/usage</code>，并确认 Key 有权限、未禁用且本机网络可访问 anytoken。'],
        ['切换后仍使用旧配置', '关闭旧终端并重新启动对应 CLI；同时确认 CC Switch 中当前 Provider 已启用。'],
      ] },
    ] },
  ],
  previous: { label: 'Codex', to: '/docs/integrations/codex' },
  next: { label: 'Cursor', to: '/docs/integrations/cursor' },
} : {
  title: 'CC Switch', eyebrow: 'Integration guide', description: 'Import anytoken providers and switch Claude Code, Codex, and Gemini CLI configurations from CC Switch.',
  intro: 'anytoken supports the <code>ccswitch://</code> import protocol and can send the endpoint, API key, application type, and usage configuration to CC Switch.', breadcrumbs: ['Integrations', 'CC Switch'],
  sections: [
    { id: 'install', title: 'Install CC Switch', blocks: [{ type: 'paragraph', content: 'Install a current CC Switch release and launch it once so the operating system registers the URL protocol.' }, { type: 'code', label: 'macOS', content: 'brew install --cask cc-switch' }] },
    { id: 'create-key', title: 'Create an API key', blocks: [{ type: 'list', items: ['Use an Anthropic group for Claude Code.', 'Use an OpenAI group for Codex.', 'Use a Gemini group for Gemini CLI.'] }] },
    { id: 'one-click-import', title: 'One-click import', blocks: [{ type: 'steps', items: ['Open the API Key page.', 'Click Import to CC Switch.', 'Allow the browser to open CC Switch.', 'Review and confirm the provider.'] }] },
    { id: 'provider-mapping', title: 'Provider mapping', blocks: [{ type: 'table', headings: ['anytoken group', 'CC Switch app', 'Endpoint'], rows: [['Anthropic', 'Claude', '<code>https://api.anytoken.com</code>'], ['OpenAI', 'Codex', '<code>https://api.anytoken.com</code>'], ['Gemini', 'Gemini', '<code>https://api.anytoken.com</code>']] }] },
    { id: 'usage', title: 'Usage status', blocks: [{ type: 'paragraph', content: 'The imported provider checks <code>/v1/usage</code> for validity and remaining balance.' }] },
    { id: 'switch-provider', title: 'Activate the provider', blocks: [{ type: 'steps', items: ['Enable the provider in CC Switch.', 'Restart the relevant CLI when needed.', 'Send a test request.', 'Confirm the request in anytoken logs.'] }] },
    { id: 'troubleshooting', title: 'Troubleshooting', blocks: [{ type: 'table', headings: ['Issue', 'Resolution'], rows: [['Import does nothing', 'Launch or reinstall CC Switch so <code>ccswitch://</code> is registered.'], ['Wrong application', 'Check the API key group platform.'], ['Old configuration remains active', 'Restart the terminal and re-enable the provider.']] }] },
  ], previous: { label: 'Codex', to: '/docs/integrations/codex' }, next: { label: 'Cursor', to: '/docs/integrations/cursor' },
})

const integrationDescriptions: Record<string, { zh: string; en: string; env: string; mode: 'anthropic' | 'openai-form' | 'shell' | 'json'; model?: string; note?: string }> = {
  'claude-code': { zh: '通过环境变量或用户级 settings.json 将 Claude Code 接入 anytoken。', en: 'Connect Claude Code to anytoken with environment variables or settings.json.', mode: 'anthropic', env: 'export ANTHROPIC_BASE_URL="https://api.anytoken.com"\nexport ANTHROPIC_AUTH_TOKEN="$ANYTOKEN_API_KEY"\nexport CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1', model: '<CLAUDE_MODEL_ID>', note: 'Claude Code 使用站点根地址，不要在 ANTHROPIC_BASE_URL 后追加 /v1。' },
  cursor: { zh: '在 Cursor 的自定义 OpenAI API 设置中配置 anytoken，用于聊天、编辑和代码理解。', en: 'Configure anytoken in Cursor custom OpenAI API settings.', mode: 'openai-form', env: 'Provider: OpenAI Compatible\nBase URL: https://api.anytoken.com/v1\nAPI Key: <ANYTOKEN_API_KEY>\nModel ID: <MODEL_ID>', note: 'Cursor 不同版本的自定义模型入口可能变化；若当前版本不允许自定义 Base URL，请使用其支持的 OpenAI Compatible Provider 入口。' },
  cline: { zh: '在 Cline 中选择 OpenAI Compatible，通过 anytoken 执行规划、编辑和命令任务。', en: 'Use the OpenAI Compatible provider in Cline through anytoken.', mode: 'openai-form', env: 'API Provider: OpenAI Compatible\nBase URL: https://api.anytoken.com/v1\nAPI Key: <ANYTOKEN_API_KEY>\nModel ID: <MODEL_ID>' },
  'roo-code': { zh: '为 Roo Code 的多模式 Agent 工作流配置统一模型网关和团队权限。', en: 'Configure a unified model gateway for Roo Code agent workflows.', mode: 'openai-form', env: 'API Provider: OpenAI Compatible\nBase URL: https://api.anytoken.com/v1\nAPI Key: <ANYTOKEN_API_KEY>\nModel ID: <MODEL_ID>' },
  'kilo-code': { zh: '将 Kilo Code 接入 anytoken，用项目级 Key 管理代码任务和模型费用。', en: 'Connect Kilo Code to anytoken with project-scoped access.', mode: 'openai-form', env: 'Provider: OpenAI Compatible\nBase URL: https://api.anytoken.com/v1\nAPI Key: <ANYTOKEN_API_KEY>\nModel ID: <MODEL_ID>' },
  aider: { zh: '在 Git 仓库中通过 anytoken 使用 Aider，统一模型选择、密钥和请求日志。', en: 'Use Aider through anytoken inside a Git repository.', mode: 'shell', env: 'export OPENAI_API_BASE="https://api.anytoken.com/v1"\nexport OPENAI_API_KEY="$ANYTOKEN_API_KEY"\naider --model openai/<MODEL_ID>', note: 'Aider 的模型前缀取决于所用 provider 适配器；如果当前版本要求 openai/ 前缀，请保留该前缀，实际请求中的模型 ID 仍应与 /v1/models 一致。' },
  'cherry-studio': { zh: '在 Cherry Studio 中添加 OpenAI 兼容服务，用于多模型聊天和知识库工作流。', en: 'Add anytoken as an OpenAI-compatible service in Cherry Studio.', mode: 'openai-form', env: 'Provider type: OpenAI Compatible\nAPI URL: https://api.anytoken.com/v1\nAPI Key: <ANYTOKEN_API_KEY>\nModel ID: <MODEL_ID>' },
  opencode: { zh: '在 OpenCode 中配置 anytoken Provider，让命令行代码任务使用统一模型入口。', en: 'Configure an anytoken provider in OpenCode.', mode: 'json', env: '{\n  "$schema": "https://opencode.ai/config.json",\n  "provider": {\n    "anytoken": {\n      "options": {\n        "baseURL": "https://api.anytoken.com/v1",\n        "apiKey": "{env:ANYTOKEN_API_KEY}"\n      }\n    }\n  }\n}', note: '完整模型清单应从控制台“使用密钥”中的 OpenCode 配置生成器复制，以保持上下文、输出上限和模型能力设置一致。' },
  'qwen-code': { zh: '为 Qwen Code 配置 OpenAI 兼容端点，并通过 anytoken 管理团队模型权限。', en: 'Configure an OpenAI-compatible anytoken endpoint for Qwen Code.', mode: 'openai-form', env: 'Provider: OpenAI Compatible\nBase URL: https://api.anytoken.com/v1\nAPI Key: <ANYTOKEN_API_KEY>\nModel ID: <MODEL_ID>', note: '仅在当前 Qwen Code 版本提供 OpenAI Compatible 或自定义 Provider 时使用此方式。' },
}

function genericIntegrationPage(slug: string): DocPage {
  const item = integrationDescriptions[slug] || integrationDescriptions.cursor
  const label = integrationItems.value.find(([, to]) => to.endsWith(`/${slug}`))?.[0] || slug
  const index = integrationItems.value.findIndex(([, to]) => to.endsWith(`/${slug}`))
  const previousItem = integrationItems.value[Math.max(0, index - 1)]
  const nextItem = integrationItems.value[Math.min(integrationItems.value.length - 1, index + 1)]
  return {
    title: label,
    eyebrow: isZh.value ? '集成指南' : 'Integration guide',
    description: isZh.value ? item.zh : item.en,
    intro: isZh.value ? `本页介绍如何将 ${label} 接入 anytoken。所有请求通过统一网关发送，密钥、模型权限和请求日志可由团队集中管理。` : `This page explains how to connect ${label} to anytoken for centralized credentials, model access, and request logs.`,
    breadcrumbs: ['Integrations', label],
    sections: isZh.value ? [
      { id: 'before-you-start', title: '准备工作', blocks: [{ type: 'list', items: ['在 anytoken 控制台创建项目 API Key。', '确认团队已授权需要使用的模型。', `安装并启动最新稳定版本的 ${label}。`] }] },
      { id: 'configuration', title: '配置', blocks: [
        { type: 'paragraph', content: item.mode === 'anthropic' ? `为 ${label} 设置 Anthropic 兼容环境变量，或把相同变量写入用户级配置。` : item.mode === 'shell' ? `在启动 ${label} 的终端中设置环境变量：` : item.mode === 'json' ? `在 OpenCode 配置文件中新增 anytoken Provider：` : `在 ${label} 的 Provider 或模型设置中选择 OpenAI Compatible，并填入以下参数：` },
        { type: 'code', label: item.mode === 'json' ? 'opencode.json' : 'Configuration', content: item.env },
        ...(item.note ? [{ type: 'note' as const, title: '配置说明', content: item.note }] : []),
      ] },
      { id: 'model', title: '选择模型', blocks: [{ type: 'paragraph', content: `先调用 <code>GET /v1/models</code> 或查看模型广场，把示例中的 <code>${item.model || '<MODEL_ID>'}</code> 替换为当前项目实际可用的完整 Model ID。生产环境建议固定模型，不要依赖未定义的别名。` }] },
      { id: 'verify', title: '验证连接', blocks: [{ type: 'steps', items: [`在 ${label} 中发送一个简短测试请求。`, '在 anytoken 请求日志中确认状态码、模型和费用。', '执行一次低风险真实任务，确认流式输出和工具调用符合预期。'] }] },
      { id: 'troubleshooting', title: '常见问题', blocks: [{ type: 'table', headings: ['现象', '处理方式'], rows: [['401', '检查 API Key 是否有效、所属分组是否正确，以及客户端是否正确传递认证头。'], ['404', item.mode === 'anthropic' ? 'Claude Code 使用站点根地址；不要额外添加 <code>/v1</code>。' : 'OpenAI Compatible 客户端通常使用包含 <code>/v1</code> 的 Base URL。'], ['模型不可用', '用 <code>GET /v1/models</code> 核对完整 Model ID，并检查项目模型权限。'], ['响应中断', '确认客户端支持流式响应，并在 anytoken 请求日志中查看网关或上游错误。']] }] },
    ] : [
      { id: 'before-you-start', title: 'Before you start', blocks: [{ type: 'list', items: ['Create a project API key.', 'Allow the required models.', `Install a stable version of ${label}.`] }] },
      { id: 'configuration', title: 'Configuration', blocks: [{ type: 'code', label: 'Configuration', content: item.env }] },
      { id: 'model', title: 'Choose a model', blocks: [{ type: 'paragraph', content: 'Use a complete Model ID returned by <code>GET /v1/models</code>, then pin it for production workflows.' }] },
      { id: 'verify', title: 'Verify', blocks: [{ type: 'steps', items: ['Send a short test request.', 'Check the anytoken request log.', 'Run one low-risk real task.'] }] },
      { id: 'troubleshooting', title: 'Troubleshooting', blocks: [{ type: 'table', headings: ['Issue', 'Resolution'], rows: [['401', 'Check the API key.'], ['404', 'Include <code>/v1</code> in the Base URL.'], ['Model unavailable', 'Check the Model ID and project policy.']] }] },
    ],
    previous: previousItem && previousItem[1] !== `/docs/integrations/${slug}` ? { label: previousItem[0], to: previousItem[1] } : undefined,
    next: nextItem && nextItem[1] !== `/docs/integrations/${slug}` ? { label: nextItem[0], to: nextItem[1] } : undefined,
  }
}

function overviewPage(): DocPage {
  return isZh.value ? {
    title: 'API 参考', description: 'anytoken 提供统一、稳定的 AI 模型 API。使用一个端点和一套凭证访问团队授权的模型。', breadcrumbs: [],
    sections: [
      { id: 'overview', title: '概览', blocks: [{ type: 'paragraph', content: 'anytoken API 兼容常见的 OpenAI 请求格式。现有应用通常只需替换 Base URL、API Key 和模型 ID 即可接入。' }, { type: 'code', label: 'Base URL', content: 'https://api.anytoken.com/v1' }] },
      { id: 'authentication', title: '身份验证', blocks: [{ type: 'paragraph', content: '所有 API 请求都通过 Bearer Token 鉴权。请为不同项目和环境创建独立 Key。' }, { type: 'code', label: 'HTTP Header', content: 'Authorization: Bearer $ANYTOKEN_API_KEY' }] },
      { id: 'quick-start', title: '快速开始', blocks: [{ type: 'paragraph', content: '先从 <code>GET /v1/models</code> 获取当前 Key 可访问的模型，再替换示例中的 <code>&lt;MODEL_ID&gt;</code>。' }, { type: 'code', label: 'cURL', content: 'curl https://api.anytoken.com/v1/chat/completions \\\n  -H "Authorization: Bearer $ANYTOKEN_API_KEY" \\\n  -H "Content-Type: application/json" \\\n  -d \'{"model":"<MODEL_ID>","messages":[{"role":"user","content":"你好"}]}\'' }] },
      { id: 'api-list', title: '接口列表', blocks: [{ type: 'table', headings: ['接口', '路径'], rows: [['聊天', '<code>POST /v1/chat/completions</code>'], ['Responses', '<code>POST /v1/responses</code>'], ['嵌入', '<code>POST /v1/embeddings</code>'], ['图像生成', '<code>POST /v1/images/generations</code>'], ['模型', '<code>GET /v1/models</code>']] }] },
    ], next: { label: '聊天', to: '/docs/api/chat' },
  } : {
    title: 'API Reference', description: 'Use one anytoken endpoint and one credential set to access approved AI models.', breadcrumbs: [],
    sections: [{ id: 'overview', title: 'Overview', blocks: [{ type: 'paragraph', content: 'anytoken supports common OpenAI-compatible request formats.' }, { type: 'code', label: 'Base URL', content: 'https://api.anytoken.com/v1' }] }, { id: 'authentication', title: 'Authentication', blocks: [{ type: 'code', label: 'HTTP Header', content: 'Authorization: Bearer $ANYTOKEN_API_KEY' }] }], next: { label: 'Chat', to: '/docs/api/chat' },
  }
}

function apiPage(slug: string): DocPage {
  const api: Record<string, { zh: string; en: string; method: string; path: string; summary: string; body: string }> = {
    chat: { zh: '聊天', en: 'Chat', method: 'POST', path: '/v1/chat/completions', summary: '发送多轮消息并获得模型回复，支持流式输出和工具调用。', body: '{\n  "model": "<MODEL_ID>",\n  "messages": [\n    { "role": "user", "content": "你好" }\n  ],\n  "stream": true\n}' },
    responses: { zh: 'Responses', en: 'Responses', method: 'POST', path: '/v1/responses', summary: '使用 Responses 格式执行文本生成、工具调用和 Codex 类 Agent 任务。', body: '{\n  "model": "<MODEL_ID>",\n  "input": "分析这个需求并给出实施步骤",\n  "stream": true\n}' },
    embeddings: { zh: '嵌入', en: 'Embeddings', method: 'POST', path: '/v1/embeddings', summary: '将文本转换为向量，用于检索、聚类和相似度计算。', body: '{\n  "model": "text-embedding-3-small",\n  "input": "需要向量化的文本"\n}' },
    images: { zh: '图像生成', en: 'Image generation', method: 'POST', path: '/v1/images/generations', summary: '通过兼容的图像模型生成图片。可用模型、尺寸和返回格式取决于项目分组与上游能力。', body: '{\n  "model": "gpt-image-1",\n  "prompt": "A clean product illustration",\n  "size": "1024x1024"\n}' },
    models: { zh: '模型', en: 'Models', method: 'GET', path: '/v1/models', summary: '查询当前 API Key 可以访问的模型及其 ID。', body: '' },
  }
  const item = api[slug] || api.chat
  const title = isZh.value ? item.zh : item.en
  const apiOrder = ['chat', 'responses', 'embeddings', 'images', 'models']
  const apiIndex = Math.max(0, apiOrder.indexOf(slug))
  const apiLabels = isZh.value ? { chat: '聊天', responses: 'Responses', embeddings: '嵌入', images: '图像生成', models: '模型' } : { chat: 'Chat', responses: 'Responses', embeddings: 'Embeddings', images: 'Image generation', models: 'Models' }
  const previousSlug = apiOrder[apiIndex - 1]
  const nextSlug = apiOrder[apiIndex + 1]
  const platformNotice = slug === 'embeddings' || slug === 'images'
    ? (isZh.value
        ? { type: 'note' as const, title: '平台限制', content: '此接口仅对 OpenAI 分组开放。使用其他平台分组的 API Key 会返回 404。' }
        : { type: 'note' as const, title: 'Platform limit', content: 'This endpoint is available only to OpenAI groups. API keys from other group platforms return 404.' })
    : undefined
  const code = item.body ? `curl https://api.anytoken.com${item.path} \\\n  -H "Authorization: Bearer $ANYTOKEN_API_KEY" \\\n  -H "Content-Type: application/json" \\\n  -d '${item.body.replace(/\n/g, '')}'` : `curl https://api.anytoken.com${item.path} \\\n  -H "Authorization: Bearer $ANYTOKEN_API_KEY"`
  return {
    title, eyebrow: isZh.value ? 'AI 模型接口' : 'AI Model API', description: isZh.value ? item.summary : `${item.method} ${item.path}`, breadcrumbs: [title],
    sections: isZh.value ? [
      { id: 'endpoint', title: '接口', blocks: [{ type: 'code', label: 'HTTP', content: `${item.method} ${item.path}` }, { type: 'paragraph', content: item.summary }, ...(platformNotice ? [platformNotice] : [])] },
      ...(item.body ? [{ id: 'request', title: '请求体', blocks: [{ type: 'code' as const, label: 'JSON', content: item.body }, { type: 'table' as const, headings: ['字段', '说明'], rows: [['<code>model</code>', '模型广场中的完整 Model ID。'], ['<code>stream</code>', '是否使用流式响应；仅适用于支持该参数的接口。']] }] }] : []),
      { id: 'example', title: '请求示例', blocks: [{ type: 'code', label: 'cURL', content: code }] },
      { id: 'errors', title: '错误处理', blocks: [{ type: 'table', headings: ['状态码', '说明'], rows: [['400', '请求参数或模型格式错误。'], ['401', 'API Key 缺失或无效。'], ['429', '触发项目速率或预算限制。'], ['5xx', '网关或上游模型暂时不可用。']] }] },
    ] : [{ id: 'endpoint', title: 'Endpoint', blocks: [{ type: 'code', label: 'HTTP', content: `${item.method} ${item.path}` }, ...(platformNotice ? [platformNotice] : [])] }, { id: 'example', title: 'Example', blocks: [{ type: 'code', label: 'cURL', content: code }] }],
    previous: previousSlug ? { label: apiLabels[previousSlug as keyof typeof apiLabels], to: `/docs/api/${previousSlug}` } : { label: isZh.value ? 'API 参考' : 'API Reference', to: '/docs' },
    next: nextSlug ? { label: apiLabels[nextSlug as keyof typeof apiLabels], to: `/docs/api/${nextSlug}` } : { label: isZh.value ? '集成指南' : 'Integrations', to: '/docs/integrations' },
  }
}

function integrationsPage(): DocPage {
  const cards = integrationItems.value.slice(1).map(([label, to]) => `<a href="${to}">${label}</a>`)
  return {
    title: isZh.value ? '集成指南' : 'Integrations', description: isZh.value ? '将常用开发工具、IDE 和桌面客户端接入 anytoken 的统一模型网关。' : 'Connect developer tools, IDEs, and desktop clients to anytoken.', breadcrumbs: ['Integrations'],
    sections: [{ id: 'overview', title: isZh.value ? '选择集成' : 'Choose an integration', blocks: [{ type: 'list', items: cards }] }, { id: 'common-settings', title: isZh.value ? '通用配置' : 'Common settings', blocks: [{ type: 'paragraph', content: isZh.value ? 'OpenAI Compatible 工具通常使用以下参数；Claude Code、Codex、Gemini CLI 等工具请以各自文章为准。' : 'OpenAI Compatible tools usually use these values. Follow the dedicated guide for Claude Code, Codex, and Gemini CLI.' }, { type: 'code', label: 'OpenAI Compatible', content: 'Base URL: https://api.anytoken.com/v1\nAPI Key:   <ANYTOKEN_API_KEY>\nModel:     <MODEL_ID>' }] }],
    previous: { label: isZh.value ? '模型' : 'Models', to: '/docs/api/models' }, next: { label: 'Claude Code', to: '/docs/integrations/claude-code' },
  }
}

const page = computed<DocPage>(() => {
  const path = route.path.replace(/\/$/, '') || '/docs'
  if (path === '/docs') return overviewPage()
  if (path === '/docs/integrations') return integrationsPage()
  if (path === '/docs/integrations/codex') return codexPage.value
  if (path === '/docs/integrations/cc-switch') return ccSwitchPage.value
  if (path.startsWith('/docs/integrations/')) return genericIntegrationPage(path.split('/').pop() || 'cursor')
  if (path.startsWith('/docs/api/')) return apiPage(path.split('/').pop() || 'chat')
  return overviewPage()
})

const activeNavigationClass = 'bg-primary-50 font-medium text-primary-700 dark:bg-primary-950/40 dark:text-primary-200'
const idleNavigationClass = 'text-slate-600 hover:bg-slate-100 hover:text-slate-950 dark:text-slate-300 dark:hover:bg-slate-900 dark:hover:text-white'

function isCurrentPage(path: string) {
  return route.path.replace(/\/$/, '') === path
}

async function copyCode(content: string, id: string) {
  await navigator.clipboard.writeText(content)
  copiedCode.value = id
  window.setTimeout(() => { copiedCode.value = '' }, 1500)
}

function observeSections() {
  observer?.disconnect()
  const sections = page.value.sections.map(section => document.getElementById(section.id)).filter(Boolean) as HTMLElement[]
  activeSection.value = page.value.sections[0]?.id || ''
  observer = new IntersectionObserver(entries => {
    const visible = entries.filter(entry => entry.isIntersecting).sort((a, b) => a.boundingClientRect.top - b.boundingClientRect.top)
    if (visible[0]) activeSection.value = visible[0].target.id
  }, { rootMargin: '-110px 0px -65% 0px', threshold: [0, 1] })
  sections.forEach(section => observer?.observe(section))
}

watch([() => route.path, () => locale.value], async () => {
  mobileNavigationOpen.value = false
  await nextTick()
  observeSections()
  window.scrollTo({ top: 0 })
})

onMounted(() => nextTick(observeSections))
onBeforeUnmount(() => observer?.disconnect())
</script>

<style scoped>
.doc-prose :deep(h2) {
  @apply mb-4 text-2xl font-semibold tracking-normal text-slate-950 sm:text-[28px] dark:text-white;
}

.doc-prose :deep(p) {
  @apply my-4 text-[15px] leading-7 text-slate-600 sm:text-base dark:text-slate-300;
}

.doc-prose :deep(code:not(pre code)) {
  @apply rounded-md bg-slate-100 px-1.5 py-0.5 font-mono text-[0.9em] font-medium text-slate-900 dark:bg-slate-800 dark:text-slate-100;
}

.doc-prose :deep(a) {
  @apply font-medium text-primary-700 underline decoration-primary-300 underline-offset-4 hover:text-primary-800 dark:text-primary-300 dark:hover:text-primary-200;
}

.doc-lead {
  @apply mt-0 rounded-2xl border border-slate-200 bg-slate-50 p-5 text-base leading-8 dark:border-slate-800 dark:bg-slate-900/60;
}

.code-block {
  @apply my-5 overflow-hidden rounded-2xl border border-slate-800 bg-slate-950 shadow-sm;
}

.code-toolbar {
  @apply flex h-11 items-center justify-between border-b border-white/10 px-4 text-xs font-medium text-slate-400;
}

.code-toolbar button {
  @apply inline-flex items-center gap-1.5 rounded-lg px-2 py-1.5 text-slate-400 transition hover:bg-white/10 hover:text-white;
}

.code-block pre {
  @apply max-w-full overflow-x-auto p-5 font-mono text-[13px] leading-6 text-blue-100;
}

.doc-list {
  @apply my-5 space-y-2.5 pl-6 text-[15px] leading-7 text-slate-600 marker:text-primary-500 dark:text-slate-300;
  list-style: disc;
}

.doc-steps {
  @apply my-6 space-y-4;
}

.doc-steps li {
  @apply flex gap-3;
}

.doc-steps li > span {
  @apply mt-1 flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-primary-50 text-xs font-semibold text-primary-700 ring-1 ring-primary-100 dark:bg-primary-950/50 dark:text-primary-200 dark:ring-primary-900;
}

.doc-steps li > p {
  @apply my-0;
}

.doc-note {
  @apply my-6 flex gap-3 rounded-2xl border border-primary-200 bg-primary-50/70 p-4 text-primary-700 dark:border-primary-900/70 dark:bg-primary-950/30 dark:text-primary-200;
}

.doc-note strong {
  @apply text-sm font-semibold;
}

.doc-note p {
  @apply mb-0 mt-1 text-sm leading-6 text-primary-800 dark:text-primary-200;
}

.doc-table-wrap {
  @apply my-6 overflow-x-auto rounded-2xl border border-slate-200 dark:border-slate-800;
}

.doc-table-wrap table {
  @apply w-full min-w-[540px] border-collapse text-left text-sm;
}

.doc-table-wrap th {
  @apply border-b border-slate-200 bg-slate-50 px-4 py-3 font-semibold text-slate-900 dark:border-slate-800 dark:bg-slate-900 dark:text-white;
}

.doc-table-wrap td {
  @apply border-b border-slate-100 px-4 py-3.5 leading-6 text-slate-600 last:border-b-0 dark:border-slate-800 dark:text-slate-300;
}

.doc-page-card {
  @apply flex min-h-24 items-center gap-3 rounded-2xl border border-slate-200 bg-white p-4 text-slate-600 shadow-sm shadow-slate-950/5 transition-all duration-300 hover:-translate-y-px hover:border-primary-200 hover:text-primary-700 hover:shadow-md hover:shadow-primary-500/10 dark:border-slate-800 dark:bg-slate-950 dark:text-slate-300 dark:hover:border-primary-800 dark:hover:text-primary-200;
}

.doc-page-card span {
  @apply flex min-w-0 flex-col;
}

.doc-page-card small {
  @apply mb-1 text-xs text-slate-400;
}

.doc-page-card strong {
  @apply truncate text-sm font-semibold;
}
</style>
