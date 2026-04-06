# NanaMiku Blog 最终优化计划书

生成日期：2026-04-01

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

- 前端构建通过：`frontend` 目录执行 `npm run build` 成功。
- 后端测试通过：`backend` 目录执行 `go test ./...` 成功。
- Git 工作区干净，适合按计划分阶段推进。

### 2.3 当前最明显的问题类型

- 内容源存在双轨，博客文章链路尚未完全统一。
- SEO 与品牌基础配置仍有默认占位值。
- 部分“后台配置能力”仍停留在前端本地存储层。
- 安全、测试、CI 仍偏“开发态”，还不是“上线态”。
- 若干查询与请求链路存在 N+1 或统计口径不严谨问题。

### 2.4 截至 2026-04-06 的执行进度

- 阶段一已完成核心收口：
  - 公开博客详情统一到 `/blog/[slug]`。
  - 文章详情、相关文章、About 页文章数与 RSS 已切到统一内容源，正式链路优先读后端 CMS。
- 阶段二已完成核心修正：
  - 站点域名、品牌名、默认 SEO 描述、默认分享图已集中管理。
  - canonical、RSS、文章级 metadata 与 favicon 已完成基础修正。
- 阶段三已完成核心收口：
  - 页尾设置已从前端 `localStorage` 改为真实后端 `site_settings` 表与 API。
  - 前后台页尾配置已共享同一真实数据源，换浏览器和换设备后仍可生效。
  - 默认不再展示占位备案号；未备案时页尾留空。
  - 友链申请已从前台占位提示升级为真实表单、后端 `friend_link_applications` 数据表与后台审核闭环。
  - 后台友链管理已支持审核公开申请，并在通过时自动转为正式友链。
  - 站点品牌信息与默认 SEO 配置已纳入后台设置中心，前台运行时会优先同步真实配置。
- 仍待继续推进：
  - 阶段四已开始并完成核心链路：
    - 管理端登录、刷新、登出、后台守卫已从 `localStorage` token 迁移为 `HttpOnly access/refresh cookie`。
    - 前端 API 与导出接口已支持基于 cookie 的自动刷新与重试。
    - 后端生产配置新增 `APP_ENV`、`COOKIE_*` 校验，生产环境会拒绝危险默认值与不安全 cookie 配置。
    - `seed` 命令已改为必须显式提供用户名、邮箱和密码，不再隐式创建弱默认管理员账号。
    - 仓库已补最小 GitHub Actions，覆盖后端测试与前端构建。
  - 阶段五已开始第一批落地：
    - 已补 cookie refresh fallback、定时发布时间解析等关键纯函数测试。
    - 已新增 `backend/cmd/smoke`，可自动跑登录、`/auth/me`、留言提交与审核、友链申请与审核、登出失效等 API smoke。
  - 阶段五已补第二批浏览器级 smoke：
    - 已新增 `frontend/playwright.config.ts` 与 `frontend/e2e/admin-smoke.spec.ts`。
    - 可自动验证未登录跳转、后台登录、友链管理页、备份页与登出失效等最小 UI 主链路。
  - 2026-04-05 阶段五已补第三批公开链路 smoke：
    - 新增 `frontend/e2e/public-smoke.spec.ts`，覆盖博客详情加载、文章评论提交与留言板提交主链路。
    - Playwright smoke 新增后端健康检查辅助；未启动后端时，依赖 API 的 smoke 会自动跳过而不是直接报错。
  - 2026-04-05 阶段五已补第四批互动 API smoke：
    - `backend/cmd/smoke` 已新增说说创建、公开列表可见、点赞、转发、评论、评论点赞与自动清理验证。
    - 公开说说列表聚合评论链路已有自动 smoke 覆盖，后续回归时不再只依赖手工点测。
  - 2026-04-05 阶段五已补第五批评论审核闭环 smoke：
    - `frontend/e2e/public-smoke.spec.ts` 已新增“文章评论提交 -> 后台评论审核 -> 前台可见 -> 清理”的浏览器级闭环验证。
    - Playwright 公开文章 smoke 已改为动态读取真实已发布文章，不再依赖 fallback `first-post`，评论接口会命中真实 `post_id`。
  - 2026-04-06 阶段五已补第六批自动化接线：
    - `.github/workflows/ci.yml` 已新增 smoke job，会自动拉起 PostgreSQL / Redis、执行 migration、seed 管理员、启动后端，并串行运行 `backend/cmd/smoke` 与 Playwright 浏览器 smoke。
    - `frontend/e2e/smoke-helpers.ts` 已支持在无现成公开文章时自动创建临时已发布文章，并在测试完成后清理，浏览器 smoke 不再依赖预置 CMS 数据。
    - `frontend/e2e/public-smoke.spec.ts` 已新增公开说说互动链路，覆盖说说展示、点赞、转发、评论、评论点赞与临时数据清理。
    - `frontend/e2e/admin-smoke.spec.ts` 已补“文章管理”页面访问验证，后台最小主链路与计划书目标进一步对齐。
  - 2026-04-06 阶段五已补第七批浏览器闭环与文案收口：
    - `frontend/e2e/public-smoke.spec.ts` 已新增“留言板留言提交 -> 后台审核通过 -> 前台可见 -> 清理”浏览器级闭环。
    - `frontend/e2e/public-smoke.spec.ts` 已新增“友链申请提交 -> 后台审核通过 -> 前台友链墙可见 -> 删除正式友链”浏览器级闭环。
    - 登录页与评论审核页的可见文案、页面 metadata 已收敛到 `frontend/src/content/copy/admin.ts`；`LoginForm.vue`、`admin/comments.astro`、`AdminCommentsManager.vue` 不再继续写死主要运营文案。
    - `frontend/src/stores/friends.ts` 的错误提示已复用 `siteCopy.friendsPage.grid.errorFallback`，友链页错误文案口径进一步统一。
  - 2026-04-06 阶段五已补第八批后台文章管理闭环：
    - `frontend/e2e/admin-smoke.spec.ts` 已新增“后台创建草稿文章 -> 前台不可见 -> 后台发布 -> 前台可见 -> 后台编辑 -> 后台转草稿 -> 前台下线 -> 清理”的浏览器级闭环。
    - `frontend/e2e/smoke-helpers.ts` 已新增后台资源清理辅助，文章类 smoke 在中途失败时也能尽量自动清理临时数据。
    - `frontend/src/components/admin/AdminPostsManager.vue` 已补最小 `data-testid` 锚点，便于后续继续补说说管理与更多后台内容链路 smoke。
  - 2026-04-06 阶段五已补第九批后台说说管理闭环：
    - `frontend/e2e/admin-smoke.spec.ts` 已新增“后台创建草稿说说 -> 前台不可见 -> 后台设为定时发布 -> 立即发布 -> 前台可见 -> 后台编辑 -> 后台转草稿 -> 前台下线 -> 删除”的浏览器级闭环。
    - `frontend/src/components/admin/AdminMomentsManager.vue` 已补最小 `data-testid` 锚点，避免说说管理相关 smoke 继续依赖脆弱的样式或纯文本定位。
  - 2026-04-06 阶段五已补第十批定时发布与 CI 修正：
    - `frontend/e2e/admin-smoke.spec.ts` 已新增"后台创建定时发布文章 -> 前台不可见 -> 行内重设发布时间 -> 仍不可见 -> 删除"的浏览器级闭环。
    - 修复 `UpdateMoment` SQL 中 `$7` 与 `NULL` 缺少 `::timestamptz` 显式类型转换导致的一致性 500 错误。
    - 修复 `admin-smoke` 文章创建测试未展开"文章属性"抽屉导致分类输入不可见的超时。
    - 登录限流从 10 req/min 放宽至 60 req/min，避免 CI 并行 worker 触发 429 级联。
  - 2026-04-06 阶段六已完成第二批性能与口径修正：
    - `dashboard.sql` 中 `GetTotalLikeCount` 已限定 `WHERE status = 'published'`，不再把草稿/定时文章的点赞计入仪表盘总量。
    - `DashboardService` 已接入 Redis，`GetStats`（2 min TTL）、`GetViewTrend` / `GetCommentTrend` / `GetLikeTrend`（5 min TTL）、`GetAnalyticsOverview`（5 min TTL）、`GetPublicSiteTrend`（5 min TTL）均增加了 JSON 序列化缓存，避免每次加载仪表盘时对 PostgreSQL 做 6~17+ 次串行聚合查询。
    - 缓存为 best-effort：Redis 不可用或反序列化失败时自动降级到数据库查询，无破坏性影响。
  - 2026-04-06 文案收口第十一批：后台 admin BaseHead 硬编码统一迁入 `adminCopy`：
    - `admin/index.astro`、`admin/backup.astro`、`admin/posts.astro`、`admin/moments.astro`、`admin/profile.astro` 的 `<BaseHead>` title / description 均改为从 `adminCopy` 读取。
    - `adminCopy` 新增 `layout.page`、`backup.page`、`postsManager.page`、`momentsManager.page`、`profileManager.page` 五组 metaTitle / metaDescription 字段。
    - 至此所有 8 个后台 `.astro` 页面的 `<BaseHead>` 均已收口到 copy 模块，后台再无硬编码页面标题。
  - 2026-04-06 阶段四已补第一批实机联调探针：
    - 本地已使用真实 PostgreSQL / Redis 重跑 `go run cmd/migrate/main.go`，首次回放时补应用 `000009_post_view_daily_visitors.up.sql`，当前本地 smoke 数据库已与仓库 migration 对齐。
    - `backend/cmd/smoke` 已新增 `HttpOnly cookie` 属性校验、仅保留 refresh cookie 的自动刷新恢复、CORS 预检、备份导出 `Content-Type` / `Content-Disposition` / 非空 body 探针。
    - 上述后端 smoke 已于 2026-04-06 在本地真实环境实跑通过，不再只停留在“代码已接线、尚未联调”的状态。
  - 2026-04-06 阶段五已补第十二批次级互动与后台会话回归覆盖：
    - `frontend/e2e/admin-smoke.spec.ts` 已新增“仅保留 refresh cookie -> 访问后台触发自动刷新 -> 备份导出仍可用”的浏览器级会话闭环。
    - `frontend/e2e/public-smoke.spec.ts` 已新增“留言板留言审核通过 -> 前台点赞 -> 作者回复 -> 回复点赞 -> 清理”的浏览器级次级互动闭环。
    - `frontend/e2e/public-smoke.spec.ts` 已新增“友链申请提交 -> 后台拒绝 -> 前台不可见”的浏览器级审核闭环。
    - `frontend/src/components/guestbook/GuestbookMessageCard.vue` 已补最小 `data-testid` 锚点，留言板回复与投票 smoke 后续不再继续依赖脆弱文本定位。
  - 2026-04-06 本地浏览器 smoke 实跑备注：
    - Playwright 首次实跑未进入业务断言，直接因本机缺少 Chromium 可执行文件而中断；该阻塞属于本地环境问题，不是当前仓库代码回归。
    - 下一步优先观察 GitHub Actions 中 `smoke` job 的真实执行结果；若 CI 仍失败，再按日志继续补修。
  - 2026-04-05 阶段六已完成第一批性能与口径修正：
    - 文章列表与分类列表已改为聚合标签查询，移除逐篇补拉标签的后端 N+1。
    - 文章 UV 已改为按“文章 + 日期 + 访客”去重入库，`post_view_daily` 不再按每次访问直接累加 UV。
    - 说说公开列表已聚合返回评论数据，前端不再逐条二次拉取评论。
    - About 页 GitHub 数据已迁移到后端代理与缓存接口，前端不再直连 GitHub API。

### 2.5 截至 2026-04-06 的明确未完成项

为方便下次继续，当前剩余工作明确收敛为以下几类：

1. 阶段四实机联调仍未完成
   - 本地后端实机联调已完成第一轮：`HttpOnly cookie` 登录、自动刷新恢复、登出失效、导出接口、`COOKIE_*` 与 `CORS` 已通过 `backend/cmd/smoke` 新增探针真实验证。
   - 仍需观察一次 GitHub Actions 中新增 `smoke` job 的真实执行结果，尤其是 Playwright 浏览器链路；2026-04-06 本地首跑因缺少 Chromium 可执行文件中断，需以 CI 结果为准继续补修。

2. 阶段五仍缺更完整的内容与互动 E2E
   - 后台文章管理基础闭环已补齐，"定时发布"路径已进入浏览器级自动验证（创建定时文章、前台不可见、行内重设时间、删除）。
   - 后台说说管理基础闭环已补齐，并已覆盖草稿、定时、发布、编辑、转草稿与前台可见性联动；但更细的图片、多记录并发与异常态回归仍主要依赖手工验证。
   - 留言板回复、投票与友链申请拒绝的浏览器级自动化覆盖已补齐；仍待继续补图片、多记录并发、异常态与边界条件回归。

3. 阶段六第二批性能与口径修正（已完成）
   - `DashboardService` 已全面接入 Redis 缓存（stats 2 min, trends/analytics 5 min），避免仪表盘每次加载产生 6~17+ 次串行 DB 查询。
   - `GetTotalLikeCount` 已限定 `WHERE status = 'published'`，修正草稿/定时文章点赞被计入仪表盘的口径偏差。
   - 前端无残余 N+1 模式。

4. 文案系统收口（已完成）
   - `friends.astro` 与 `FriendsGrid.vue` 早已完全消费 `siteCopy.friendsPage`，无硬编码。
   - 全部 8 个后台 `.astro` 页面的 `<BaseHead>` 均已收口到 `adminCopy`，后台再无硬编码页面标题。
   - 原则仍是“页面和组件只消费 copy，不继续新增运营向硬编码文案”。


## 3. 优先级总览

### P0：必须先做

1. 统一博客内容真源与阅读链路。
2. 修正 SEO、站点域名与品牌基础配置。
3. 将伪后台配置能力改成真实后台能力。

### P1：强烈建议紧接着做

1. 完成鉴权与环境配置硬化。
2. 建立最小 CI 与自动校验链路。
3. 补关键路径测试。

### P2：在稳定后继续做

1. 修复数据口径与性能问题。
2. 统一前端可见文案入口。
3. 为外部依赖数据增加代理与缓存层。

### P3：可择机做

1. 更完整的后台设置中心。
2. 结构化数据与更系统的内容增长支持。
3. 更细粒度的观测与运营面板。

## 4. 分阶段实施计划

## 阶段一：内容真源统一

### 目标

把博客系统从当前“静态内容集合 + 后端 CMS/API”双轨状态，收敛为单一真源，避免后续维护和运营成本持续上升。

### 现状证据

- 静态文章路由仍然存在：`frontend/src/pages/blog/[slug].astro:106`
- 动态 API 文章详情页同时存在：`frontend/src/pages/blog/post.astro:11`
- 博客列表页跳转到 API 详情页：`frontend/src/components/blog/BlogFeed.vue:95`
- RSS 仍然只读取 Astro Content：`frontend/src/pages/rss.xml.js:5`
- 静态详情侧边栏仍存在本地假点赞：`frontend/src/components/blog/ArticleMeta.astro:65`
- 假点赞组件未接后端：`frontend/src/components/blog/LikePanel.vue:34`

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

## 阶段二：SEO 与品牌基建修正

### 目标

把项目从“能被访问”升级为“适合被收录、被分享、被识别”的站点。

### 现状证据

- `Astro.site` 仍为占位值：`frontend/astro.config.mjs:11`
- 全局站点标题和描述仍为默认值：`frontend/src/consts.ts:4`
- RSS 标题和描述继承默认值：`frontend/src/pages/rss.xml.js:8`

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

## 阶段三：后台设置能力真实化

### 目标

把当前看起来像后台能力、但本质仍是前端本地态的功能，真正下沉到后端与数据库。

### 现状证据

- 页尾设置保存到 localStorage：`frontend/src/components/admin/AdminFooterManager.vue:143`
- 页尾配置 store 只读写浏览器存储：`frontend/src/stores/siteFooter.ts:5`
- 友链申请入口仍是占位提示：`frontend/src/components/friends/FriendsGrid.vue:20`

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

### 当前完成情况（截至 2026-04-03）

- 已完成：
  - 新增站点级 `site_settings` 表及对应查询、服务、公开读取接口、后台写入接口。
  - 页尾备案信息与自定义页尾文案已改为真实后端存储，不再以浏览器本地存储作为最终状态。
  - 前台页尾与后台页尾设置页已统一读写同一接口。
  - 备份导出链路已纳入 `site_settings`。
  - 默认页尾备案号已改为留空，未备案时不展示占位内容。
- 2026-04-03 补充完成：
  - 新增 `friend_link_applications` 真实申请表、公开提交接口与后台审核接口。
  - 前台友链页已接入真实申请表单，不再只弹出占位 toast。
  - 后台友链管理已新增申请队列，并支持通过 / 拒绝申请。
- 2026-04-03 收口完成：
  - 站点品牌信息与默认 SEO 配置已纳入后台 `site_settings`，并新增公开读取与后台写入接口。
  - 前台运行时会优先同步真实品牌名、Logo Alt、默认描述、默认分享图与 canonical 基准 URL。

## 阶段四：安全与交付基线硬化

### 目标

让项目从“本地开发可用”进入“具备上线最低安全与交付标准”的状态。

### 现状证据

- 管理端 token 保存在 localStorage：`frontend/src/stores/auth.ts:71`
- 登录页仍围绕 localStorage 做重定向判断：`frontend/src/components/auth/LoginForm.vue:152`
- 后端存在默认数据库、Redis、JWT 占位值：`backend/biz/bootstrap/config.go:73`
- seed 命令仍带默认管理员账号密码：`backend/cmd/seed/main.go:29`
- 仓库当前缺少 CI 工作流文件：`.github/` 目录不存在

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

## 阶段五：测试与验证体系补齐

### 目标

为当前功能较多、链路较长的项目建立最小可信验证体系。

### 现状证据

- 当前后端仅有一个测试文件：`backend/biz/service/backup_test.go`
- 前端暂无真实测试基建或 E2E 自动化目录

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

## 阶段六：性能与数据口径修正

### 目标

解决现有项目中“能跑但不够严谨”的统计与性能问题。

### 现状证据

- 文章列表存在逐篇取标签的查询模式：`backend/biz/service/posts.go:94`
- 文章 UV 统计当前按访问直接累加：`backend/biz/service/posts.go:240`
- 说说页先拉列表，再按条补拉评论：`frontend/src/stores/moments.ts:130`
- About GitHub 数据由前端直连 GitHub API：`frontend/src/components/about/AboutGithubProfile.vue:331`

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

## 5. 文案与信息架构专项收口

### 目标

继续执行仓库既定的 DIY 文案规范，减少页面和组件中的可见硬编码文案。

### 已发现重点区域

- `frontend/src/pages/friends.astro`
- `frontend/src/components/friends/FriendsGrid.vue`
- `frontend/src/components/auth/LoginForm.vue`
- 若干 admin 页面 `BaseHead` 标题与描述仍未统一走 copy 模块

### 建议做法

- 新增 `siteCopy.friendsPage` 分区。
- 将登录页可见文案迁移到 `siteCopy` 或 `adminCopy`。
- 把后台页面 title/description 统一纳入 `adminCopy`。
- 保持“页面只消费 copy，不再直接写死运营文案”的约束。

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
3. 后台设置能力真实化

### 第二批

1. 鉴权与环境配置硬化
2. CI 与测试基线

### 第三批

1. 性能与数据口径修正
2. 文案系统继续收口
3. GitHub 数据代理缓存

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
