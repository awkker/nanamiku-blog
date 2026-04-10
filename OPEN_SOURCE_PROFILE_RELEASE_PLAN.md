# 个人资料脱敏与 About 页配置收口计划书

生成日期：2026-04-10

约束依据：
- 根目录 `agents.md`
- `ADMIN_SETTINGS_CENTER_UPGRADE_PLAN.md`
- `ADMIN_SETTINGS_CENTER_SETTINGS_DICTIONARY.md`
- `frontend/src/content/copy/site.ts`
- `frontend/src/content/copy/admin.ts`
- `frontend/src/lib/author-profile.ts`
- `frontend/src/lib/site-integrations.ts`
- `frontend/src/stores/authorProfile.ts`
- `frontend/src/stores/siteIntegrations.ts`
- `frontend/src/components/about/AboutIntroSection.vue`
- `frontend/src/components/about/AboutContactSection.vue`
- `frontend/src/pages/about.astro`
- `frontend/src/pages/friends.astro`
- 根目录 `README.md`

## 1. 计划目标

这份计划不追求“一次性把所有文案都继续搬后台”，而是分成两条线并行处理：

1. 先完成开源前必须做的个人信息脱敏与默认值清理，确保仓库公开后不会继续暴露真实邮箱、社交账号、头像或强个人化描述。
2. 再把 About 页里确实需要频繁 DIY、且已经明显超出 `copy/` 静态默认值边界的公开资料，升级为设置中心可维护配置。

核心判断：

- 当前最大问题已经不是“后台没有入口”，而是“源码默认值、组件兜底值、示例值、文档和素材里仍混有真实个人信息”。
- 设置中心已经具备 `author_profile` 和 `site_integrations`，说明基础设施足够，下一步应聚焦收口与治理，而不是再做一轮无边界扩张。

## 2. 当前现状判断

### 2.1 已经具备的能力

当前仓库实际上已经有一条可用的“公开资料运行时配置链路”：

| 业务域 | 当前 key | 现状 |
| --- | --- | --- |
| 首页首屏 | `home_hero` | 已后台化 |
| 首页视觉资源 | `home_assets` | 已后台化 |
| 作者公开资料 | `author_profile` | 已后台化 |
| GitHub / 天气 / 首页组件开关 | `site_integrations` | 已后台化 |
| 站点品牌 / 页尾 | `site_profile` / `footer` | 已后台化 |

说明：

- 这意味着 About 页里“展示名、角色、简介、技能、近况、联系邮箱、公开社交链接、GitHub 用户名”并不是完全写死的。
- 真正没收口的，是 `copy/site.ts` 中仍残留大量真实默认值，以及 About 页更深层的个性化内容块。

### 2.2 当前仍然存在的问题

| 类别 | 当前问题 | 风险 |
| --- | --- | --- |
| 公开默认值 | `site.ts` 中仍有真实姓名、邮箱、GitHub、X、QQ、Bilibili、域名、项目链接、个性化文案 | 仓库开源后直接暴露个人信息，且默认模板不通用 |
| 组件兜底值 | `author-profile.ts`、`AboutIntroSection.vue` 仍保留真实邮箱等硬编码 fallback | 即使后台未配置，也会回退到真实信息 |
| 后台示例值 | `admin.ts` 中 placeholder 使用了真实昵称、邮箱、GitHub 用户名 | 设置中心截图、录屏、源码都会继续泄露信息 |
| About 页内容边界 | 时间线、主推项目、月度目标、正在听、固定签名等仍全部留在 `copy/site.ts` | DIY 成本高，且个性化数据不适合继续跟源码强绑定 |
| 静态素材 | `frontend/public/picture/author.jpg` 等素材可能包含真实头像或强个人识别素材 | 开源仓库素材层仍可能暴露个人信息 |
| 文档与仓库卫生 | 根 `README.md` 的 seed 命令还是旧写法；根 `.gitignore` 过于精简 | 新用户跑不通，且后续更容易误提交本地敏感文件 |

## 3. About 页的治理边界

遵守 `agents.md` 的“文案集中管理”与“不是所有 copy 都该进后台”原则，About 页建议拆成三层。

### 3.1 保留在 `copy/` 的内容

这部分继续留在 `frontend/src/content/copy/site.ts`，但默认值必须改成开源安全版本：

- 标题标签、按钮文案、分区标题、导航文案
- `aria-label`
- 空态、提示态、说明性系统文案
- About 页的静态结构词，例如“写作地图”“主推项目”“本月目标”

### 3.2 继续使用现有设置中心的内容

以下内容继续放在现有 key，不新增业务域：

| 内容 | 配置 key |
| --- | --- |
| 展示名、头像、角色、Blog 简介、About 简介、所在地、起始时间 | `author_profile` |
| 技能标签、近况、引言、联系邮箱、公开社交链接、联系扩展入口 | `author_profile` |
| GitHub 用户名、首页天气/音乐/时钟开关 | `site_integrations` |

### 3.3 建议新增后台配置的 About 专属内容

About 页中以下内容，已经明显属于“高 DIY、强个人化、适合后台维护”的范围，建议新增一个单独 key：

建议 key：
- `about_page`

建议承载字段：

| 模块 | 建议字段 |
| --- | --- |
| 首屏说明补充 | `intro_cards` |
| 成长经历 | `milestones` |
| 能力分组 | `capability_groups` |
| 主推项目 | `featured_projects` |
| 月度目标 | `monthly_goals` |
| 正在听 | `listening_now` |
| 固定签名 | `signature` |

不建议在第一阶段继续迁移：

- `heroBadge`、`projectsSection.title`、`writingMapSection.title` 这类结构性标题
- `writingMap` 这种更像站内阅读引导而不是个人档案的数据

理由：

- 这些结构标题属于页面骨架，适合继续放在 `copy/`
- `writingMap` 更偏站点阅读路径，不属于强个人隐私，也不是最急需后台化的模块

### 3.4 DIY 图片目录约定

为避免图片资源继续散落在多个目录，建议把“可后台配置、可手工 DIY 的前台图片”统一收口到一个固定目录。

建议目录：
- `frontend/public/picture/diy/`

目录规则：
- 允许在该目录下继续按模块建子目录。
- 推荐子目录：`home/`、`about/`、`friends/`、`site/`、`projects/`
- 所有后台可编辑的图片字段，默认优先引用这个目录下的资源。

后台填写规则：
- 后台表单只填写相对路径，不填写本地绝对路径。
- 建议统一写成以 `/picture/diy/` 开头的公开路径。
- 示例：`/picture/diy/about/avatar.webp`
- 示例：`/picture/diy/home/hero/cover-01.webp`

适用范围：
- 首页背景图
- 作者头像与 About 页展示图
- 项目卡片封面
- 友链页站点展示图
- 其它后续准备交给设置中心维护的公开图片

约束建议：
- 新增可配置图片时，优先复用 `frontend/public/picture/diy/`，不要再把资源散放到 `picture/` 根目录各处。
- 后端保存的图片字段，优先保存公开相对路径；只有明确需要引用外部图片 CDN 时，才允许完整 URL。
- README 与后台帮助文案应同步说明这一规则，降低后续 DIY 成本。

## 4. 开源前必须先做的事

为了让仓库尽快具备公开条件，建议把工作拆成“上线阻塞项”和“后续增强项”。

### 4.1 上线阻塞项

这些建议在开源前完成，否则仓库仍然会暴露个人信息或给别人错误的默认体验。

1. 把 `frontend/src/content/copy/site.ts` 中真实个人资料替换为开源安全默认值。
2. 清理 `frontend/src/lib/author-profile.ts` 与 `frontend/src/components/about/AboutIntroSection.vue` 中的真实邮箱硬编码 fallback。
3. 将 `frontend/src/content/copy/admin.ts` 中所有真实示例值改为通用示例值。
4. 审核 `frontend/public/picture/author.jpg` 是否为真实头像；若是，替换为通用占位图或项目级默认头像。
5. 修正根 `README.md` 中已过期的 seed 命令，统一为当前后端要求的显式 `username / email / password` 写法。
6. 扩充根 `.gitignore`，至少补齐 `.env`、`.env.*`、`*.local` 等忽略规则，避免未来误提交。
7. 确认可公开图片资源的统一目录方案，优先收口到 `frontend/public/picture/diy/`。

### 4.2 后续增强项

这些不一定阻塞开源，但会直接影响后续 DIY 体验与长期维护成本。

1. 为 About 页新增 `about_page` 配置链路与设置中心分组。
2. 为 About 页复杂列表增加 normalize / 校验 / 默认回退逻辑。
3. 将 About 页的个性化内容从 `copy/site.ts` 收口到新的 store。
4. 为新增图片字段补上“相对路径优先”的前后端校验与帮助提示。
5. 补一轮 About 页与设置中心的前后端测试，避免后台保存后前台显示错位。

## 5. 推荐实施顺序

### Phase 0：脱敏盘点与边界冻结

目标：
- 明确哪些内容属于“真实个人资料”
- 明确哪些内容属于“公开但不该放进源码默认值”

交付：
- 一份敏感项清单
- 一份允许保留在仓库中的公开模板字段清单

建议覆盖对象：
- 邮箱
- GitHub / X / QQ / Bilibili 账号
- 真实头像
- 真实域名
- 强个人化句子、歌单、签名、项目路径

### Phase 1：源码默认值脱敏

目标：
- 不改变现有页面结构，不改视觉风格，只把默认值替换为开源安全版本

重点范围：
- `frontend/src/content/copy/site.ts`
- `frontend/src/content/copy/admin.ts`
- `frontend/src/lib/author-profile.ts`
- `frontend/src/lib/site-integrations.ts`
- `frontend/src/components/about/AboutIntroSection.vue`
- `frontend/src/components/about/AboutContactSection.vue`
- `frontend/src/pages/friends.astro`

交付：
- About / Blog / 首页 / 友链页仍可正常展示
- 后台未配置时只展示通用模板信息，不回退到真实个人信息
- 联系入口在未配置时可隐藏或回退到通用占位，不再出现真实邮箱

### Phase 2：About 页新增配置域

目标：
- 为 About 页专属个性化内容新增后台配置能力

建议新增：
- 后端 `site_settings` key：`about_page`
- 前端 `frontend/src/lib/about-page.ts`
- 前端 `frontend/src/stores/aboutPage.ts`
- 后台设置中心新增 `?section=about-page`

建议字段边界：
- 列表类内容使用数组结构，并提供最大长度限制
- 项目链接统一走公开 URL 规范化
- 空数组时回退到 `copy/` 的模板默认值
- 后端保存时写审计日志

### Phase 3：设置中心 UI 收口

目标：
- 在当前 `/admin/settings` 内继续收口，而不是再开新的零散页面

后台分组建议：

| 分组 | 内容 |
| --- | --- |
| 作者与社交 | 继续承载 `author_profile` |
| 集成与服务 | 继续承载 `site_integrations` |
| About 页面 | 新增承载 `about_page` |

About 页面分组内建议展示：
- 时间线编辑器
- 项目卡片编辑器
- 月度目标列表
- 正在听列表
- 固定签名编辑
- 实时预览区域

### Phase 4：文档、素材与仓库卫生修复

目标：
- 让仓库真正可公开、可克隆、可二次 DIY

交付：
- 根 `README.md` 修正启动流程
- 增加“首次部署后去后台填写公开资料”的说明
- 根 `.gitignore` 补全忽略规则
- 素材目录完成一轮版权与隐私审查
- 补充“DIY 图片统一放到 `frontend/public/picture/diy/`，后台填写相对路径”的说明

额外建议：

- `ADMIN_SETTINGS_CENTER_SETTINGS_DICTIONARY.md` 里的“回退顺序”表述，建议改成两套描述：
- 初始化顺序：`copy defaults -> local cache -> backend fetch`
- 运行时最终优先级：`backend real config -> local cache -> copy defaults`

这样更符合当前 store 的真实行为，避免文档继续误导。

### Phase 5：上线前最终检查

建议在准备发布前做一次统一检查。

### 5.1 搜索检查

建议对以下模式做仓库级搜索，要求结果只剩“通用示例值”或“明确允许保留的文档样例”：

- 真实邮箱
- 真实 GitHub 用户名
- 真实 X / QQ / Bilibili 链接
- 真实域名
- 真实昵称

### 5.2 运行检查

建议至少完成：

1. 前端 `npm run build`
2. 后端 `go test ./...`
3. 关键页面手工检查：`/`、`/blog`、`/about`、`/friends`、`/login`、`/admin/settings`
4. 未配置 `author_profile` / `site_integrations` / `about_page` 时的空态与回退检查

### 5.3 开源可用性检查

至少确认：

1. 新开发者按 README 能完成启动
2. 没有任何真实个人联系方式被默认展示
3. 没有真实头像或私人素材被默认暴露
4. 后台首次登录后可以把真实公开资料重新填回去

## 6. 建议的最终治理模型

建议最终固定为四层：

| 层级 | 存放位置 | 内容边界 |
| --- | --- | --- |
| 系统级静态文案 | `frontend/src/content/copy/` | 按钮、标题、提示语、结构性文案、开源安全默认值 |
| 公开运营配置 | `site_settings` | 作者公开资料、About 个性化内容、首页门面、GitHub 用户名、组件开关 |
| 敏感配置 | `.env` | Token、密码、JWT Secret、数据库与 Redis 凭证 |
| 私人初始化数据 | 本地未跟踪文件或后台手动填写 | 仓库不开源的真实个人资料 |

## 7. 建议优先级

如果目标是“尽快准备开源并上线”，推荐优先级如下：

### P0：必须先做

- 脱敏 `site.ts`
- 清理真实 fallback
- 清理后台示例值
- 检查头像素材
- 修正 README
- 补根 `.gitignore`

### P1：建议紧接着做

- 新增 `about_page` 配置域
- 设置中心新增 About 页面分组
- 补充 About 页的 store 与回退逻辑

### P2：上线后继续做

- 优化 About 页面编辑体验
- 增加更多审计与预览能力
- 继续检查内容正文与媒体素材的个人信息、授权和版权边界

## 8. 结论

当前仓库已经完成了“设置中心基础设施”这一步，真正缺的是“开源安全默认值”和“About 页专属配置域”的最后收口。

更稳妥的执行策略不是继续大面积重构，而是：

1. 先把所有真实个人信息从默认值、fallback、placeholder、素材、文档里移出去。
2. 再把 About 页剩余的高 DIY 内容单独接入 `about_page`，继续沿用现有 `/admin/settings` 架构。

这样既符合 `agents.md` 的边界控制，也最适合你当前“准备上线 + 计划开源”的阶段目标。
