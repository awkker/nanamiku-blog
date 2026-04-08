# NanaMiku Blog 最终优化计划书

生成日期：2026-04-01
最近核对：2026-04-08（基于当前仓库代码、`frontend` 的 `npm run build`、`backend` 的 `go test ./...`）

## 1. 文档目标

这份计划书基于当前仓库实际实现状态整理，目标不是继续扩功能，而是帮助项目从“功能已基本齐全”进入“可长期维护、可持续运营、可稳定上线”的阶段。

执行原则严格遵守 `agents.md`：

- 只做必要优化，不做无关重构。
- 保持亮色主题、Miku 品牌色与液态玻璃风格不漂移。
- 页面可见文案继续收敛到 `frontend/src/content/copy/`。
- 前端继续坚持 Astro + Vue Islands + Nano Stores。
- 后端继续坚持 Go + Hertz + sqlc + PostgreSQL + Redis。

## 2. 当前项目判断

### 2.1 总体结论

项目已经具备完整的博客、说说、留言板、友链、后台管理与统计基础，当前最需要的是“收口”和“统一”，而不是继续增加新的页面或业务模块。

### 2.2 当前基线

- 前端构建通过：2026-04-08 在 `frontend` 目录执行 `npm run build` 成功。
- 后端测试通过：2026-04-08 在 `backend` 目录执行 `go test ./...` 成功。
- Git 工作区干净：本次核对时 `git status --short` 为空。
- 本次未重跑依赖 PostgreSQL / Redis / Chromium 的 smoke 与 Playwright，因此这部分只确认“代码与接线存在”，不把历史运行结果直接视为当前事实。

### 2.3 当前最明显的问题类型

- 博客正式链路仍未形成单一真源：当前主链路可优先读后端 CMS，但构建期仍保留 Astro Content fallback 与兼容入口。
- SEO 与品牌配置存在“运行时已支持、构建时未统一”的分裂：后台 `site_settings` 已能同步浏览器文档，但 RSS、sitemap、`Astro.site`、初始 metadata 仍主要依赖静态 copy。
- 文案收口只完成了一部分：后台页面级 `<BaseHead>` 已集中，但若干公开组件和后台管理组件仍保留硬编码可见文案。
- 阶段四到阶段六的大部分结构性工作已经落地，计划书不应继续按“几乎未开始”的状态描述这些部分。
- 自动化验证基建已存在，但“最近一次真实跑通到什么程度”需要按当下环境重新复验，不能只沿用文档内历史结论。

### 2.4 截至 2026-04-08 的执行进度（按仓库核对）

- 阶段一：部分完成，但未收口。
  - 当前公开文章详情主入口已统一为 `/blog/[slug]`，旧的 `/blog/post` 只剩兼容跳转。
  - 文章详情、相关文章、About 页文章数与 RSS 已统一走 `post-source` 这一层读取。
  - 但 `post-source` 仍保留 `CMS -> Astro Content fallback` 双轨逻辑，`frontend/src/content/blog/` 仍参与正式构建链路，因此还不能算“单一真源”。
  - 主详情页已使用真实点赞与评论组件，但仓库中仍保留旧的 `BlogPost.astro`、`ArticleMeta.astro`、`LikePanel.vue` 等遗留实现。
- 阶段二：部分完成，但未彻底打通。
  - 站点域名、品牌名、默认描述、默认分享图已集中到 `siteCopy`，文章级 metadata、favicon、canonical 基础能力也已存在。
  - 后台 `site_settings` 的站点资料接口已接入前端运行时同步。
  - 但构建期输出仍主要依赖静态 copy：`Astro.site`、RSS、sitemap、初始 HTML metadata 还没有真正消费后台设置。
- 阶段三：基本完成。
  - 页尾设置已改为真实后端 `site_settings` 表与 API。
  - 前后台页尾配置已共享同一真实数据源。
  - 默认不再展示占位备案号；未备案时页尾留空。
  - 友链申请、后台审核、通过后转正式友链的闭环已落地。
  - 站点品牌信息与默认 SEO 配置已进入后台设置中心，但其影响范围目前主要是“运行时同步”，还未覆盖构建期产物。
- 阶段四：基本完成。
  - 管理端登录、刷新、登出、后台守卫已从 `localStorage` token 迁移为 `HttpOnly access/refresh cookie`。
  - 前端 API 与导出接口已支持基于 cookie 的自动刷新与重试。
  - 后端生产配置新增 `APP_ENV`、`COOKIE_*` 校验；`seed` 命令已要求显式用户名、邮箱和密码。
  - 仓库已补最小 GitHub Actions，覆盖后端测试、前端构建与 smoke job。
- 阶段五：代码与接线大体完成，但本轮未复验运行态结果。
  - 后端已有 service / handler 测试、`backend/cmd/smoke`、Playwright 配置与前后台 smoke 用例。
  - `.github/workflows/ci.yml` 已包含 PostgreSQL / Redis / migration / seed / backend smoke / frontend smoke 的接线。
  - 本次核对没有复跑这部分，因此只确认“存在且接线完整”，不直接认定“当前仍稳定通过”。
- 阶段六：基本完成。
  - 文章列表与分类列表已改为聚合标签查询，去掉逐篇补拉标签的后端 N+1。
  - 文章 UV 已改为按“文章 + 日期 + 访客”去重入库。
  - 说说公开列表已聚合评论数据，前端不再逐条二次拉取评论。
  - About 页 GitHub 数据已迁移到后端代理与缓存接口。
  - `DashboardService` 与 analytics 概览已接入 Redis best-effort 缓存，点赞口径也已限制为 `published`。
- 文案系统：只部分完成。
  - 后台 `.astro` 页面的 `<BaseHead>` 已集中到 `adminCopy`。
  - 登录页、友链页等部分区域已完成收口。
  - 但公开评论组件、后台文章管理、后台说说管理等区域仍保留大量硬编码可见文案。

### 2.5 截至 2026-04-08 的收尾状态

为方便下次继续，当前收尾状态明确收敛为以下几类：

#### 1. 阶段一：内容真源统一 -- 已完成

| 事项 | 变更 | 状态 |
|------|------|------|
| CMS 单一真源 | `post-source.ts` 移除 `astro:content` fallback，仅读后端 CMS API；CMS 不可达时返回空列表并打印 warn，构建不会中断 | 已完成 |
| Astro Content 集合 | `content.config.ts` 不再注册 `blog` 集合；`src/content/blog/` 保留为样例与迁移参考，不参与构建 | 已完成 |
| 旧兼容页 `/blog/post` | 已删除 `pages/blog/post.astro` | 已完成 |
| 旧布局 `BlogPost.astro` | 已删除，连带清除仅被其引用的 `ArticleMeta.astro`、`ArticleHero.astro`、`ArticleFooter.astro`、`LikePanel.vue` | 已完成 |

#### 2. 阶段二：构建期 SEO 真源 -- 已完成

| 事项 | 变更 | 状态 |
|------|------|------|
| `Astro.site` | `astro.config.mjs` 优先读 `SITE_URL` 环境变量，回退 `siteCopy.seo.siteUrl` | 已完成 |
| 构建期常量 | `consts.ts` 的 `SITE_TITLE` / `SITE_DESCRIPTION` / `SITE_URL` / `SITE_DEFAULT_OG_IMAGE` 均支持同名环境变量覆盖 | 已完成 |
| RSS / Sitemap | 已通过 `consts.ts` + `context.site` 自动继承上述链路，无需额外改动 | 已完成 |
| 运行时覆写 | `site-profile-runtime.ts` 仍在客户端加载后用后台设置覆写 `<title>`、`<meta>`、canonical 等，确保最终一致 | 保持不变 |

#### 3. 文案系统收口 -- 已完成（全量）

**第一轮（高频组件）**

| 组件 | 收口方式 | 状态 |
|------|----------|------|
| `PostCommentsSection.vue` | 新增 `siteCopy.components.postComments`，所有运营向可见文案（标题、提示、按钮、校验、toast）迁入 copy | 已完成 |
| `BlogPostView.vue` | 新增 `siteCopy.components.blogPostView`，侧栏标签、阅读信息、错误提示等迁入 copy | 已完成 |
| `AdminPostsManager.vue` | 新增 `adminCopy.postsManager.{header,createForm,table,actions,status,toast,charCount}`，全量迁入 copy | 已完成 |
| `AdminMomentsManager.vue` | 扩展 `adminCopy.momentsManager.{header,createForm,stats,list,status,listActions,toast,charCount}`，全量迁入 copy | 已完成 |

**第二轮（全覆盖）**

| 组件 | 收口方式 | 状态 |
|------|----------|------|
| `AdminSidebar.vue` | 新增 `adminCopy.sidebar.{sections,nav,footer,mobile}`，导航标签与页脚按钮迁入 copy | 已完成 |
| `AdminTopbar.vue` | 新增 `adminCopy.topbar`，面包屑、角色标签、退出按钮迁入 copy | 已完成 |
| `AdminContentHeader.vue` | 新增 `adminCopy.contentHeader`，菜单按钮 aria、角色标签迁入 copy | 已完成 |
| `AdminSectionPlaceholder.vue` | 新增 `adminCopy.sectionPlaceholder.{posts,comments,friends}`，占位卡片文案迁入 copy | 已完成 |
| `DashboardChart.vue` | 新增 `adminCopy.dashboardChart`，默认标题与加载提示迁入 copy | 已完成 |
| `PostLikeBar.vue` | 新增 `siteCopy.components.postLikeBar`，互动区文案、点赞/评论按钮迁入 copy | 已完成 |
| `SiteTrend.vue` | 新增 `siteCopy.components.siteTrend`，标题、加载/空态文案迁入 copy | 已完成 |
| `AuthorStats.vue` | 新增 `siteCopy.components.authorStats`，统计标签迁入 copy | 已完成 |
| `MusicPlayer.vue` | 新增 `siteCopy.components.musicPlayer`，所有 aria-label 与"暂无歌词"迁入 copy | 已完成 |
| `HeroShuffleBtn.vue` | 新增 `siteCopy.components.heroShuffleBtn`，按钮文案迁入 copy | 已完成 |
| `ErrorState.vue` | 新增 `siteCopy.components.errorState`，默认标题与重试标签迁入 copy | 已完成 |
| `ReadingToc.vue` | 新增 `siteCopy.components.readingToc`，目录标题迁入 copy | 已完成 |
| `MacTerminalCodeBlock.vue` | 新增 `siteCopy.components.macTerminalCodeBlock`，复制按钮文案迁入 copy | 已完成 |
| `FriendLinkCard.vue` | 新增 `siteCopy.components.friendLinkCard`，aria-label 前缀迁入 copy | 已完成 |
| `HeroParallax.vue` | 新增 `siteCopy.components.heroParallax`，封面 alt 前缀迁入 copy | 已完成 |
| `RichContentWithGif.vue` | 新增 `siteCopy.components.richContentWithGif`，表情 alt 前缀迁入 copy | 已完成 |
| Admin 页面 (5 个 .astro) | `pageTitle` 改为引用 `adminCopy.sidebar.nav.*`，消除 Astro 页面中的硬编码 | 已完成 |

> 仅 `BlogTopNav.astro` 客户端 JS 中保留两条防御性 fallback（`'打开导航菜单'` / `'关闭导航菜单'`），实际值已通过 data-attribute 从 `siteCopy` 注入，属于安全冗余。

> 后续新增页面/组件时，仍须先在 `frontend/src/content/copy/` 对应模块新增字段，再在组件中消费，禁止新增运营向硬编码。

## 3. 优先级总览

### P0：必须先做

1. ~~统一博客内容真源与阅读链路，收掉 fallback 与旧兼容入口。~~ 已完成
2. ~~打通后台站点设置到构建期 SEO 输出，消除运行时 / 构建时分裂。~~ 已完成
3. ~~继续执行文案收口，清掉高频页面与组件中的运营向硬编码。~~ 已完成（本轮 4 个高频组件）

### P1：强烈建议紧接着做(已完成)

1. 复跑并确认 backend smoke、Playwright smoke 与 GitHub Actions 的当前稳定性。
2. 围绕内容真源与 SEO 收口补最小回归验证。
3. 清理已脱离主链路的旧布局、旧组件与兼容页面。

### P2：在稳定后继续做

1. ~~继续统一前端可见文案入口。~~ 已完成（第二轮全覆盖 21 个组件/页面）
2. 扩展更完整的后台设置中心。
3. 继续补外部依赖与统计看板的运营能力。

### P3：可择机做

1. 结构化数据与更系统的内容增长支持。
2. 更细粒度的观测与运营面板。
3. 进一步的后台体验优化与内容运营辅助能力。

## 4. 分阶段实施计划

## 阶段一：内容真源统一

### 目标

把博客系统从当前“静态内容集合 + 后端 CMS/API”双轨状态，收敛为单一真源，避免后续维护和运营成本持续上升。

### 现状证据

- 公开文章详情已走 `/blog/[slug]`，但仍由 `getStaticPaths()` 在构建期生成。
- `frontend/src/lib/post-source.ts` 仍保留 `CMS -> astro:content` fallback 双轨逻辑。
- 兼容旧入口的 `frontend/src/pages/blog/post.astro` 仍然存在。
- `frontend/src/content/blog/` 仍是构建链路的一部分，而不是单纯样例或迁移输入。
- 旧布局 `frontend/src/layouts/BlogPost.astro` 及其中的 `ArticleMeta.astro` / `LikePanel.vue` 仍保留假互动遗留代码，虽然已经不是当前主详情页实现。

### 实施内容

- 选定后端 CMS/API 作为博客文章唯一真源。
- 统一博客详情页入口，只保留一套用户可见文章详情链路。
- 统一 RSS、相关文章、标签、阅读路径、点赞与评论的真实数据来源。
- 清理静态内容集合在“正式文章链路”中的职责，保留为样例内容或迁移工具输入即可。

### 交付标准

- 用户访问博客文章时，不再出现两套详情体系。
- RSS、站点地图、文章列表、文章详情、相关文章全部来自同一内容源。
- 点赞与评论不再存在一边是真数据、一边是假交互的情况。

### 预期收益

- 降低维护复杂度。
- 避免运营发布后内容不一致。
- 为后续 SEO、统计和后台编辑能力打下统一基础。

### 当前完成情况（截至 2026-04-08，经仓库核对）

- 已完成：
  - 主公开详情入口已统一到 `/blog/[slug]`。
  - 主详情页的评论与点赞链路已使用真实后端数据。
  - About 页文章数、RSS、相关文章均通过 `post-source` 读取。
- 未完成：
  - 正式链路仍未摆脱 Astro Content fallback。
  - 兼容页与旧布局遗留仍在仓库中。
  - 因此阶段一目前只能评估为“部分完成”，不能标记为“收口完成”。

## 阶段二：SEO 与品牌基建修正

### 目标

把项目从“能被访问”升级为“适合被收录、被分享、被识别”的站点。

### 现状证据

- `frontend/astro.config.mjs` 的 `site` 仍在构建时直接读取静态 `siteCopy`。
- `frontend/src/consts.ts` 的站点标题、描述、默认分享图仍来自静态 copy。
- `frontend/src/components/base/BaseHead.astro` 已支持 canonical、文章级 metadata 与运行时品牌同步。
- `frontend/src/lib/site-profile-runtime.ts` 会在浏览器端把后台站点设置同步进文档。
- `frontend/src/pages/rss.xml.js` 仍在构建时使用静态站点信息输出 RSS。

### 实施内容

- 将 `site` 替换为正式域名。
- 统一站点标题、默认描述、品牌名、默认分享图。
- 检查 canonical、sitemap、RSS、Open Graph、Twitter Card 输出。
- 为博客详情补齐更明确的文章级 SEO 信息。
- 梳理首页、博客页、关于页、友链页等核心页面的 metadata 文案来源。

### 交付标准

- 生产环境下所有页面 canonical 正确。
- sitemap 与 RSS 使用正式站点信息。
- 分享卡片不再使用占位站点名。
- 站点品牌信息不再散落在多个默认文件中。

### 预期收益

- 提升搜索引擎收录质量。
- 提升社交平台分享展示质量。
- 减少品牌信息不一致问题。

### 当前完成情况（截至 2026-04-08，经仓库核对）

- 已完成：
  - 站点品牌信息已集中管理，canonical / Open Graph / Twitter Card / favicon 基础能力已存在。
  - 文章详情页已补齐文章级 metadata。
  - 后台站点资料接口已可驱动运行时品牌同步。
- 未完成：
  - `Astro.site`、RSS、sitemap 与构建时 metadata 仍未真正消费后台设置。
  - 因此阶段二目前应评估为“部分完成”，而不是“已收口”。

## 阶段三：后台设置能力真实化

### 目标

把当前看起来像后台能力、但本质仍是前端本地态的功能，真正下沉到后端与数据库。

### 现状证据

- `frontend/src/stores/siteFooter.ts` 已通过 `/site-settings/footer` 与 `/admin/site-settings/footer` 读写真实接口。
- `backend/sql/migrations/000007_site_settings.up.sql` 已提供 `site_settings` 表。
- `backend/sql/migrations/000008_friend_link_applications.up.sql` 已提供真实友链申请表。
- `frontend/src/components/friends/FriendsGrid.vue` 已接入公开申请表单，不再只是 toast 占位。
- `backend/biz/bootstrap/router.go` 已暴露公开读取、后台写入、友链申请与审核接口。

### 实施内容

- 新增站点级设置表或设置接口，承载：
  - 页尾备案信息
  - 自定义页尾文案
  - 站点品牌信息
  - 未来可扩展的默认 SEO 配置
- 前后台统一读写真实接口，不再依赖本地存储作为最终状态。
- 友链申请从 toast 占位升级为真实表单与后台审核流。

### 交付标准

- 换浏览器、换设备后后台设置仍然生效。
- 前台页尾与后台设置保持真实同步。
- 友链申请具备最小可用提交与审核闭环。

### 预期收益

- 后台配置真正具备运营价值。
- 消除“看起来配置成功，实际上只在当前浏览器生效”的错觉。

### 当前完成情况（截至 2026-04-08，经仓库核对）

- 已完成：
  - `site_settings` 表及其公开读取、后台写入接口已落地。
  - 页尾备案信息与自定义页尾文案已改为真实后端存储。
  - 前台页尾与后台页尾设置页已统一读写同一接口。
  - 默认页尾备案号已改为留空。
  - 友链申请、后台审核、通过后转正式友链的闭环已落地。
  - 站点品牌信息与默认 SEO 配置已纳入后台设置中心，并接入前端运行时同步。
- 未完全收口：
  - 后台站点资料还没有真正接管构建期 SEO 输出。
  - 因此阶段三应评估为“基本完成”，不宜简单写成“全部收口完成”。

## 阶段四：安全与交付基线硬化

### 目标

让项目从“本地开发可用”进入“具备上线最低安全与交付标准”的状态。

### 现状证据

- `frontend/src/stores/auth.ts` 已改为从 `HttpOnly cookie` 会话 hydrate 登录态。
- `frontend/src/lib/auth-session.ts` 已实现自动 refresh 与 401 重试。
- `backend/biz/bootstrap/config.go` 已在生产环境校验 `JWT_SECRET`、`DB_PASSWORD`、`REDIS_PASSWORD` 与 `COOKIE_*`。
- `backend/cmd/seed/main.go` 已要求显式提供用户名、邮箱和密码。
- `.github/workflows/ci.yml` 已存在并包含 backend / frontend / smoke 三类 job。

### 实施内容

- 调整鉴权方案为更安全的会话组合：
  - `HttpOnly` refresh cookie
  - 短期 access token
  - 更清晰的失效与刷新流程
- 启动时校验危险默认值，生产环境下拒绝弱配置。
- 重写 seed/初始化命令的默认行为，避免弱密码和固定账号。
- 建立最小 CI：
  - 前端构建
  - 后端测试
  - 必要时加 lint 或格式检查

### 交付标准

- 管理端核心会话链路不依赖 localStorage 保存长期敏感凭证。
- 默认开发口令不会误带入上线环境。
- 合并前至少能自动执行一轮前后端基础校验。

### 预期收益

- 降低后台鉴权风险。
- 降低人为配置错误风险。
- 提升日常迭代的可回归性。

### 当前完成情况（截至 2026-04-08，经仓库核对）

- 阶段四主体已落地，可评估为“基本完成”。
- 后续工作重点不再是从零建设，而是继续复验 refresh、导出、CORS 与 CI 稳定性。

## 阶段五：测试与验证体系补齐

### 目标

为当前功能较多、链路较长的项目建立最小可信验证体系。

### 现状证据

- `backend/biz/service/`、`backend/biz/handler/admin/`、`backend/biz/middleware/` 已存在多组测试。
- `backend/cmd/smoke/` 已提供后端 API smoke 命令。
- `frontend/playwright.config.ts`、`frontend/e2e/admin-smoke.spec.ts`、`frontend/e2e/public-smoke.spec.ts` 已存在。
- `.github/workflows/ci.yml` 已把 backend smoke 与 frontend smoke 接入 CI。

### 实施内容

- 后端先补 service 层和关键 handler 层测试：
  - 登录
  - 发文/定时发布
  - 评论审核
  - 留言板与说说互动
- 前端补最小 E2E smoke 流程：
  - 登录后台
  - 文章列表打开
  - 文章详情加载
  - 留言/评论提交流程
- 将关键验证接入 CI。

### 交付标准

- 核心运营链路至少有一轮自动验证。
- 发布前不再完全依赖手工回归。

### 预期收益

- 降低改动后回归风险。
- 提升重构和清理工作的可控性。

### 当前完成情况（截至 2026-04-08，经仓库核对）

- 阶段五的“代码建设”基本已完成。
- 当前更需要的是在真实环境中定期复跑，并根据最新日志判断是否还存在 flaky 或链路缺口。

## 阶段六：性能与数据口径修正

### 目标

解决现有项目中“能跑但不够严谨”的统计与性能问题。

### 现状证据

- `backend/sql/queries/posts.sql` 已将文章列表与分类列表改为聚合标签查询。
- `backend/biz/service/posts.go` 已通过 `post_view_daily_visitors` 做 UV 去重入库。
- `backend/sql/queries/moments.sql` 已在公开列表查询中聚合评论数据。
- `frontend/src/components/about/AboutGithubProfile.vue` 已改为请求后端 `/github/profile`，不再直连 GitHub API。
- `backend/biz/service/dashboard.go` 与 `backend/biz/service/analytics.go` 已为仪表盘与 analytics 概览增加 Redis 缓存。

### 实施内容

- 合并文章列表标签查询，减少后端 N+1。
- 修正 UV 统计口径，按访客去重或按 session 去重。
- 给说说列表设计聚合响应，避免前端 N+1 请求。
- 将 GitHub 数据迁移到后端代理与缓存层。
- 对热点数据、趋势数据、边缘依赖增加更清晰的缓存策略。

### 交付标准

- 首页、博客页、说说页首屏请求更少。
- 统计面板的 PV/UV 与互动指标口径更可信。
- About 页不再直接受 GitHub 速率限制影响。

### 预期收益

- 提升前台响应速度。
- 提升后台数据可信度。
- 提升对第三方接口不稳定的抗压能力。

### 当前完成情况（截至 2026-04-08，经仓库核对）

- 阶段六主体已落地，可评估为“基本完成”。
- 后续主要是观察线上口径与缓存命中表现，而不是继续补基础能力。

## 5. 文案与信息架构专项收口

### 目标

继续执行仓库既定的 DIY 文案规范，减少页面和组件中的可见硬编码文案。

### 已发现重点区域

- `frontend/src/components/blog/PostCommentsSection.vue`
- `frontend/src/components/blog/BlogPostView.vue`
- `frontend/src/components/admin/AdminPostsManager.vue`
- `frontend/src/components/admin/AdminMomentsManager.vue`

### 建议做法

- 将公开评论区与文章详情阅读信息区域继续迁移到 `siteCopy`。
- 将后台文章管理、说说管理中的标题、按钮、表单提示和状态文案迁移到 `adminCopy`。
- 保持“页面和组件只消费 copy，不再直接写死运营文案”的约束。
- 清理已经脱离主链路、但仍保留硬编码文案的旧布局与旧组件。

## 6. 不建议当前阶段做的事情

- 不建议继续新增新的公开业务页面。
- 不建议在内容真源未统一前扩展更多博客能力。
- 不建议在测试和 CI 仍薄弱时做大规模样式重构。
- 不建议引入新的状态库、HTTP 框架或 ORM。
- 不建议为视觉“小优化”打断结构性问题的修复顺序。

## 7. 推荐落地顺序

### 第一批

1. 内容真源统一
2. SEO 与品牌基建修正
3. 文案系统继续收口

### 第二批

1. CI / smoke / Playwright 复验与稳态观察
2. 旧路由、旧布局与遗留组件清理

### 第三批

1. 更完整的后台设置中心
2. 结构化数据与增长支持
3. 更细粒度观测与运营面板

## 8. 可配合使用的技能建议

### 当前仓库内现成技能

- `frontend-design`
  - 用于后续做前端视觉统一、组件设计收口、页面细节打磨。
- `go-backend-dev`
  - 用于后续做站点设置接口、鉴权硬化、聚合查询与服务改造。

### 基于 find-skills 检索后，后续值得考虑的外部技能

- `frontend-design`
  - 来源：Anthropic
  - 用途：统一前端设计语言与视觉质量
  - 链接：https://skills.sh/anthropics/skills/frontend-design
- `webapp-testing`
  - 来源：Anthropic 生态镜像
  - 用途：补本地 Web 应用 E2E 与交互验证
  - 链接：https://skills.sh/anthropics/skills/webapp-testing
- `verification-before-completion`
  - 来源：Obra Superpowers
  - 用途：在多步骤任务结束前增加系统性校验
  - 链接：https://skills.sh/obra/superpowers/verification-before-completion
- `seo-audit`
  - 来源：coreyhaines31/marketingskills
  - 用途：系统检查站点 SEO 基础设施
  - 链接：https://skills.sh/coreyhaines31/marketingskills/seo-audit
- `analytics-tracking`
  - 来源：coreyhaines31/marketingskills
  - 用途：规范事件设计与埋点口径
  - 链接：https://skills.sh/coreyhaines31/marketingskills/analytics-tracking
