# CloseAI Sub2API

CloseAI Sub2API 是基于 [Wei-Shaw/sub2api](https://github.com/Wei-Shaw/sub2api) 的个人定制 fork，用来搭建自用的 AI API 分发、中转、计费和运营管理平台。

这个仓库不是上游项目的纯镜像。它保留 Sub2API 的核心网关能力，并围绕 Ikun/CloseAI 的实际使用场景做了前端品牌、套餐支付、Codex/OpenAI 兼容、账号运营、部署脚本和 Git 管理方式的定制。

## 项目定位

- 面向个人或小团队的 AI API 网关与账号池管理平台。
- 通过平台 API Key 对外提供 OpenAI、Claude、Gemini、Antigravity、Bedrock 等能力。
- 管理上游账号、代理、分组、用量、限流、计费、订阅套餐和用户余额。
- 重点保障 Codex CLI、Claude Code、OpenAI Responses API、Chat Completions、WebSocket 和图片生成等调用链路的兼容性。
- 作为 `closeai` 总项目的子模块维护，代码魔改集中在 `custom` 分支。

## 技术栈

| 模块 | 技术 |
| --- | --- |
| 后端 | Go 1.26+, Gin, Ent, PostgreSQL, Redis |
| 前端 | Vue 3, Vite, TypeScript, Pinia, TailwindCSS |
| 支付 | EasyPay, 支付宝, 微信支付, Stripe, Airwallex |
| 部署 | Docker Compose, systemd, Caddy/OpenResty, 嵌入式前端二进制 |

## 主要魔改功能

### 品牌与前端

- 默认站点名改为 `Ikun`，替换默认 Logo 和主题样式。
- 清理上游赞助商、版本提示、在线更新入口等不适合自用部署的展示。
- 优化首页、控制台、支付页、套餐卡片、用户中心和管理后台的视觉与交互。
- 为登录注册增加协议确认、邮箱验证、OAuth 补邮箱等流程。

### 账号与网关兼容

- 强化 OpenAI/Codex OAuth 账号的导入、刷新、用量探测和调度逻辑。
- 支持 OpenAI Responses API、Chat Completions、Embeddings、Images、WebSocket 等路径的转发和兼容转换。
- 针对 Codex CLI / Claude Code 做工具调用、上下文续链、SSE 终止事件、图片桥接、模型映射和客户端识别适配。
- 支持账号模型白名单、分组模型列表、自定义 `/v1/models` 输出和上游模型同步。
- 增加代理有效期、代理失败回退、账号冷却、模型级不可用冷却和容量错误 failover。

### 计费、套餐与支付

- 优化用户订阅套餐、余额充值、支付结果页和管理端订单展示。
- 支持套餐外部订阅链接、外部订阅弹窗文案、支付目标校验和充值金额处理。
- 调整人民币展示、二维码支付、支付成功通知和订单恢复流程。
- 扩展多币种与 Airwallex 支付，并保留 EasyPay、支付宝、微信、Stripe 等支付方式。
- 支持用户按平台维度的 USD 配额、订阅日卡额度、平台用量拆分和失败请求记录。

### 管理后台与运营

- 用户列表支持按 API Key 所在分组过滤，支持查看已删除用户历史用量。
- 用量页增加按日明细、缓存命中/创建 token、图片 token、失败请求和请求类型筛选。
- 运维面板增强错误日志归因、SLA 口径、TTFT 统计、告警指标和本地业务限制分类。
- 账号管理增加批量编辑、创建时间、可用账号统计、上游模型同步和 5h/7d 用量窗口提示。
- 兑换码支持批量修改、有效期设置和更完整的后台管理。

### 登录、通知与风控

- 增加钉钉 OAuth 登录和 internal-only 用户属性同步。
- 扩展 OIDC、微信、GitHub、Google 等第三方登录的邮箱绑定和自动创建逻辑。
- 增加邮件模板管理、支付成功通知、余额提醒、订阅提醒和退订入口。
- 内容审计支持关键词拦截、风险阈值、按模型生效和管理员豁免。
- 强化响应头、CSP、URL 白名单、敏感凭证脱敏和 API Key ACL。

### 部署与工具

- 增加 `deploy/build.sh`，一键构建前端并交叉编译 Linux AMD64 / macOS ARM64 部署包。
- 提供 `deploy/Caddyfile`、Docker Compose、本地目录版部署、systemd 服务模板和示例配置。
- 保留 `skills/sub2api-admin`，用于通过 Codex Skill 管理 Sub2API 后台接口。
- 部署时重点保留 `database.*`、`redis.*`、`jwt.secret`、`totp.encryption_key`、服务端口和 `.installed` 状态文件。

## 项目结构

```text
sub2api/
├── backend/                  # Go 后端、Ent 模型、网关、计费、调度和后台 API
│   ├── cmd/server/           # 服务入口和 Wire 注入
│   ├── ent/                  # Ent schema 与生成代码
│   ├── internal/             # 核心业务代码
│   ├── migrations/           # 数据库迁移
│   └── resources/            # 价格、模型等资源文件
├── frontend/                 # Vue 3 管理端和用户端
│   └── src/
├── deploy/                   # Docker、systemd、Caddy、构建脚本和配置模板
├── docs/                     # 功能说明文档
├── skills/                   # Codex Skill 辅助工具
└── tools/                    # 项目脚本工具
```

## 本地开发

### 前置依赖

- Node.js 18+
- pnpm 9+
- Go 1.26+
- PostgreSQL 15+
- Redis 7+

### 安装依赖

```bash
cd frontend
pnpm install --frozen-lockfile

cd ../backend
go mod download
```

Windows 本地如果遇到 Go 缓存权限问题，可以把缓存放到仓库内：

```powershell
cd backend
$env:GOMODCACHE = "$PWD\..\.cache\go-mod"
$env:GOCACHE = "$PWD\..\.cache\go-build"
go mod download
```

### 启动开发环境

先准备后端配置：

```bash
cp deploy/config.example.yaml backend/config.yaml
```

启动后端：

```bash
cd backend
go run ./cmd/server
```

启动前端：

```bash
cd frontend
pnpm dev
```

### 常用检查

```bash
# 前端
cd frontend
pnpm typecheck
pnpm test:run

# 后端
cd backend
go test ./...
```

修改 `backend/ent/schema` 后需要重新生成：

```bash
cd backend
go generate ./ent
go generate ./cmd/server
```

## 构建与部署

### 一键构建部署包

```bash
cd deploy
bash build.sh
```

构建产物：

```text
backend/bin/linux-amd64/
backend/bin/darwin-arm64/
```

Linux 部署包内包含：

- `sub2api` 二进制
- `config.yaml`
- `data/`
- `start.sh`
- `sub2api.service`
- `deploy.sh`

### Docker Compose

```bash
cd deploy
cp .env.example .env
docker compose -f docker-compose.local.yml up -d
```

本地目录版 `docker-compose.local.yml` 更适合生产备份和迁移。

### systemd

生产环境推荐把服务放到 `/opt/sub2api`，通过 systemd 管理：

```bash
sudo systemctl status sub2api
sudo journalctl -u sub2api -f
sudo systemctl restart sub2api
```

反向代理如果需要粘性会话，请确认代理不会丢弃带下划线的请求头。例如 Nginx/OpenResty 需要：

```nginx
underscores_in_headers on;
```

## Git 管理方式

这个项目有两层 Git：

| 层级 | 路径 | 作用 |
| --- | --- | --- |
| 外层仓库 | `closeai/` | 管理总项目、文档、子模块指针 |
| 内层仓库 | `closeai/sub2api/` | 管理 Sub2API 实际代码魔改 |

外层仓库的 `.gitmodules` 指向：

```text
path = sub2api
url = https://github.com/ikunycj/sub2api.git
branch = custom
```

### 克隆项目

```bash
git clone --recurse-submodules git@github.com:ikunycj/closeai.git
cd closeai
git submodule update --init --recursive
```

如果已经克隆但子模块为空：

```bash
git submodule update --init --recursive
```

### 分支约定

```text
upstream/main  ->  Wei-Shaw/sub2api 开源上游
origin/main    ->  ikunycj/sub2api 的上游同步分支，原则上只同步不直接魔改
origin/custom  ->  个人魔改分支，日常开发都在这里
```

核心原则：

- `main` 用来跟上游，尽量保持干净。
- `custom` 用来放自己的功能、品牌、部署和业务改造。
- 从上游同步新功能时，把 `custom` rebase 到最新 `main` 上。
- rebase 后推送 `custom` 必须用 `--force-with-lease`，不要用裸 `--force`。

### 日常魔改流程

```bash
cd sub2api
git checkout custom
git pull --ff-only origin custom

# 修改代码
git status
git add .
git commit -m "feat: describe your change"
git push origin custom
```

如果这次提交改变了子模块指针，还要回到外层仓库提交指针：

```bash
cd ..
git status
git add sub2api
git commit -m "chore: update sub2api custom pointer"
git push origin master
```

顺序一定是：先提交并推送 `sub2api` 内层仓库，再提交外层 `closeai` 的子模块指针。

### 同步上游 Sub2API

```bash
cd sub2api
git fetch upstream

# 更新本地 main
git checkout main
git merge --ff-only upstream/main
git push origin main

# 把魔改搬到最新 main 上
git checkout custom
git rebase main
```

如果有冲突：

```bash
# 手动解决冲突后
git add <resolved-files>
git rebase --continue
```

完成后推送：

```bash
git push --force-with-lease origin custom
```

最后更新外层子模块指针：

```bash
cd ..
git add sub2api
git commit -m "chore: sync sub2api custom with upstream"
git push origin master
```

### 查看当前状态

```bash
# 外层仓库
git status --short --branch

# 内层 sub2api 仓库
git -C sub2api status --short --branch
git -C sub2api remote -v
git -C sub2api log --oneline --decorate -n 10
```

### 注意事项

- 不要把生产 `config.yaml`、数据库 dump、Redis dump、`.env`、真实密钥、`backend/bin/`、`node_modules/` 提交进仓库。
- 处理冲突时优先保留本 fork 的品牌、支付、Codex 兼容和部署脚本改造，再按需吸收上游实现。
- 外层仓库只记录子模块指针，不会自动包含 `sub2api` 内部提交内容。
- 如果外层显示 `sub2api` 有改动，先进入 `sub2api` 确认内层是否已经提交。
- 如果 rebase 后 CI 或本地测试失败，优先检查 Ent 生成代码、前端类型、支付字段和网关 DTO 是否同步。

## 上游关系

- 上游项目：[Wei-Shaw/sub2api](https://github.com/Wei-Shaw/sub2api)
- 当前 fork：[ikunycj/sub2api](https://github.com/ikunycj/sub2api)
- 总项目：[ikunycj/closeai](https://github.com/ikunycj/closeai)

需要查上游文档时，可以切到 `main` 或直接访问上游仓库；需要改自己的生产功能时，只在 `custom` 上开发。
