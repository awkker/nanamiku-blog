# Open Source Profile Release Execution Log

更新时间：2026-04-10

## Round 1

本轮按 `OPEN_SOURCE_PROFILE_RELEASE_PLAN.md` 的 `P0 / Phase 1 + Phase 4` 先处理开源前阻塞项，不扩展到新的 `about_page` 配置域。

已完成：

- 将 `frontend/src/content/copy/site.ts` 里的公开默认值替换为开源安全模板，移除真实邮箱、账号、域名与强个人化默认描述。
- 将 `frontend/src/content/copy/admin.ts` 的后台 placeholder 改为通用示例值，避免录屏、截图和源码继续暴露真实资料。
- 清理 About 页邮件 fallback，让未配置联系邮箱时不再回退到真实地址。
- 新增 `frontend/public/picture/diy/about/avatar-placeholder.svg` 作为默认头像，并统一替换相关 fallback 引用。
- 将后端测试中的真实域名、用户名与示例资料同步替换为通用测试数据，避免仓库级搜索继续命中旧值。
- 修正根 `README.md` 的 seed 命令，补充首次登录后应填写的设置项与 DIY 图片目录约定。
- 扩充根 `.gitignore`，补上 `.env`、`.env.*`、`*.local` 与常见本地构建产物忽略规则。

待继续：

- 做一轮仓库级敏感信息搜索复查，确认没有遗漏的真实账号、链接或素材引用。
- 评估 `about_page` 新配置域的字段结构、后端 key 与设置中心分组。
- 继续清理文档或素材里可能仍残留的个性化信息与图片资产。

## Round 2

本轮覆盖 `Phase 2 + Phase 3`，把 About 页中段高 DIY 内容正式接进现有 `site_settings` 与设置中心链路，并在 `/admin/settings` 内完成对应分组收口。

已完成：

- 新增后端 `site_settings` key `about_page`，并补齐公共读取、后台保存、审计日志与路由注册。
- 新增 `frontend/src/lib/about-page.ts` 与 `frontend/src/stores/aboutPage.ts`，为 About 页配置提供 normalize、payload 与本地缓存能力。
- 新增 `frontend/src/components/about/AboutDynamicSections.vue`，把 About 页中段的动态内容改为运行时消费 `about_page` 配置，同时保留 `copy/` 作为默认回退。
- 扩展 `frontend/src/components/admin/AdminSettingsCenter.vue`，新增 `?section=about-page` 分组，支持编辑补充卡片、时间线、能力分组、主推项目、本月目标、正在听与固定签名。
- 同步更新 `README.md` 与 `ADMIN_SETTINGS_CENTER_SETTINGS_DICTIONARY.md`，补齐新分组入口、fallback 规则与 `about_page` 字典说明。

阶段判断：

- `Phase 2` 已完成：`about_page` 配置域、后端 key、前端 lib/store、About 页消费链路均已打通。
- `Phase 3` 已完成：设置中心内已新增 `About 页面` 分组，并包含时间线、项目卡片、月度目标、正在听、固定签名与实时预览区域。

验证：

- `backend`: `go test ./...`
- `frontend`: `npm run build`

## Round 3

本轮继续细化 `Phase 3` 的设置中心体验，不新增配置域，也不调整前台 About 页结构，重点补强后台 `About 页面` 分组的编辑反馈与实时预览。

已完成：

- 为 `About 页面` 分组内的补充卡片、时间线、能力分组、主推项目、本月目标、正在听补上当前数量与上限显示，降低列表型配置的编辑盲区。
- 将右侧实时预览从“数量统计 + 部分卡片”扩展为更接近前台结构的小样板，新增时间线、能力分组、项目卡片、目标、正在听与固定签名的实时预览。
- 为各模块预览补上空态显示，便于区分“当前为空”与“保存后会回退到默认模板”的状态。

阶段判断：

- `Phase 3` 体验层补充已完成：`/admin/settings?section=about-page` 现在既能编辑，也能更直观地预览 About 页中段实际落点。

验证：

- `frontend`: `npm run build`

## Round 4

本轮补一个 About 页回归修复：GitHub 模块在默认值脱敏后，如果后台 `site_integrations.github_username` 未填写，就会被直接隐藏，容易造成“模块消失”的体感。

已完成：

- 保持 `site_integrations.github_username` 作为 About GitHub 模块的主配置来源。
- 为 About 页 GitHub 模块补上兼容回退：当后台未填写 GitHub 用户名时，尝试从作者公开社交链接中的 GitHub 地址解析用户名。
- 同步让 About 页底部联系区的 GitHub 按钮复用同一套解析逻辑，避免顶部 GitHub 模块和底部联系按钮行为不一致。

阶段判断：

- About 页 GitHub 展示链路已恢复到“后台显式配置优先，作者公开社交链接次级兜底”的状态。

验证：

- `frontend`: `npm run build`

## Round 5

本轮继续细化 About 页 GitHub 的默认回退策略，避免后台首次未填写时页面显得过空。

已完成：

- 将 About 页默认 GitHub 用户名改为公开测试账号 `octocat`，作为 `site_integrations.github_username` 未填写时的默认回退。
- 同步更新后台设置中心的 GitHub 缺省提示文案，使其与当前真实回退顺序保持一致。

阶段判断：

- About 页 GitHub 模块当前优先级为：后台显式配置 -> 默认测试账号 `octocat` -> 作者公开社交链接中的 GitHub 地址。

验证：

- `frontend`: `npm run build`

## Round 6

本轮处理两个回归反馈：默认测试 GitHub 用户会在本地环境持续触发 GitHub 接口报错，以及博客首页作者头像仍可能回退到历史 logo 资源。

已完成：

- 撤回 About 页的默认测试 GitHub 用户，恢复为“只有后台显式填写 `site_integrations.github_username` 时才展示 GitHub 资料卡”。
- 保留 About 页底部联系按钮的 GitHub 链接回退能力，但不再让它驱动顶部 GitHub 数据卡，避免未配置时主动请求 GitHub API。
- 将站点品牌默认文案恢复为 `NanaMiku Blog`，同步修正设置中心里的品牌 placeholder。
- 为作者公开资料新增遗留头像兜底：当后端仍返回旧的 `/picture/author.jpg` 时，前台自动替换为当前的通用默认头像，避免博客首页继续显示 logo 式头像。

阶段判断：

- About 页 GitHub 资料卡已恢复为“后台显式配置驱动”的行为。
- 博客首页作者卡片的历史头像 fallback 已完成兼容修复。

验证：

- `frontend`: `npm run build`

## Round 7

本轮继续处理设置中心使用反馈，重点解决“GitHub 用户名填写位置不清晰”和“作者头像在旧值或异常图片路径下持续停留在加载态”的问题。

已完成：

- 在设置中心 `author-profile` 分组补充提示，明确 About 页 GitHub 动态资料卡的用户名应在 `site-integrations` 分组填写。
- 为 Blog 侧栏与 About 首屏作者头像增加双层回退：先把历史 `/picture/author.jpg` 视为遗留值替换为通用占位图，再在图片加载失败时自动切回默认头像，避免持续停留在加载态。
- 将作者头像的遗留值兼容同步补到后端 `site_settings` 服务层，保证读取和保存 `author_profile` 时都能统一输出新的默认头像路径。
- 为后端作者头像兼容补上服务层测试，覆盖旧头像路径规范化场景。

阶段判断：

- `author-profile` 与 `site-integrations` 的职责边界现在在 UI 中更明确了。
- 作者头像 fallback 已形成“后端归一化 + 前端渲染兜底”的双保险。

验证：

- `backend`: `go test ./biz/service/...`
- `frontend`: `npm run build`

## Round 8

本轮继续处理设置中心移动端体验和作者头像输入兼容，重点解决“配置分组卡片过长”“作者头像仍显示异常”这两个反馈。

已完成：

- 将设置中心左侧导航改为可折叠菜单，并按“前台内容 / 作者与 About / 后台与合规”三组折叠展示，移动端默认收起，避免先出现一整张很长的配置分组卡片。
- 将作者头像默认展示收口为仓库内的 `/picture/author.jpg`，同时保留通用占位图作为第二层兜底，避免公开作者卡片继续停留在 loading 态。
- 前端作者头像地址归一化增强：新增对 `picture/author.jpg`、`./picture/author.jpg` 以及历史错误值 `https://picture/author.jpg` 的兼容，统一按站内资源处理。
- 后端 `site_settings` 的站内资源规范化同步增强，保存作者头像或其它站内图片路径时，不再把本地相对资源误判成外链。
- 为作者头像与站内资源规范化补充服务层测试，覆盖本地路径和历史错误 URL 场景。

阶段判断：

- 设置中心左侧导航现在更接近“左侧菜单 + 折叠分组”的结构，移动端表单区可用空间明显更稳定。
- 作者头像链路已形成“本地路径纠正 + 主默认图 + 通用占位图”的三级兜底。

验证：

- `backend`: `go test ./biz/service/...`
- `frontend`: `npm run build`

## Round 9

本轮纠正设置中心导航位置理解偏差，重点把配置分组从设置页内部导航，真正迁移到后台最左侧主侧边栏的 `设置中心` 节点下。

已完成：

- 新增设置中心当前分组的 Nano Store，同步 `?section=` 查询参数与后台主侧边栏高亮状态。
- 将后台最左侧主侧边栏的 `设置中心` 改成可折叠子菜单，直接展开 `前台内容 / 作者与 About / 后台与合规` 三组配置入口。
- 设置页内部原有那列分组导航已移除，页面主体改为“当前分组表单 + 右侧预览”两栏，避免出现两套重复分组入口。

阶段判断：

- 配置分组现在已经放到后台最左侧的主导航层级中，而不是继续停留在设置页内部。

验证：

- `frontend`: `npm run build`
