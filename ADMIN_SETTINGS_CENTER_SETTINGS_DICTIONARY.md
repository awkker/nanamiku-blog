# 后台设置中心设置项字典

生成日期：2026-04-09

本文档用于收口 Phase 5，统一说明：

- `site_settings` 每个 key 的职责边界
- 默认值与 fallback 来源
- 前台消费位置
- 哪些文案禁止继续迁入后台
- 旧入口的清理结果

## 1. Fallback 规则

后台设置中心统一遵守以下回退顺序：

1. 前端 `copy/` 中的静态默认值
2. 对应 Nano Stores 的本地缓存值
3. 后端 `site_settings` 实时配置

具体约束：

- `copy/` 永远保留，负责系统文案与默认值。
- `site_settings` 只负责“可运营调整”的公开配置。
- 后端未配置、接口失败或客户端未完成 hydration 时，页面必须可退回到 `copy/` 默认值。
- 敏感信息只允许存在于 `.env`，禁止进 `site_settings`、前端源码或后台明文表单。

## 2. 设置项字典

| key | 作用域 | 管理入口 | 前台消费 | fallback 来源 |
| --- | --- | --- | --- | --- |
| `site_profile` | 品牌与默认 SEO | `/admin/settings?section=site-profile` | `BaseHead`、导航品牌、页尾品牌、RSS 标题 | `frontend/src/content/copy/site.ts` |
| `footer` | 页尾与备案信息 | `/admin/settings?section=site-footer` | `SiteFooter.vue` | `frontend/src/content/copy/site.ts` |
| `home_hero` | 首页首屏标题与副标题 | `/admin/settings?section=home-hero` | 首页 Hero 文案与打字机副标题 | `frontend/src/content/copy/site.ts` + 本地缓存 |
| `home_assets` | 首页视觉资源 | `/admin/settings?section=home-assets` | 首页背景图与换图逻辑 | `frontend/src/content/copy/site.ts` + 本地缓存 |
| `author_profile` | 作者展示资料 | `/admin/settings?section=author-profile` | Blog 侧栏、About 首屏、About 联系区 | `frontend/src/content/copy/site.ts` + Store fallback |
| `site_integrations` | GitHub、天气与首页组件开关 | `/admin/settings?section=site-integrations` | About GitHub、首页天气/音乐/时钟、后台依赖状态 | `frontend/src/content/copy/site.ts` + Store fallback |

## 3. 前端 Store 对照

| Store | 对应 key | 说明 |
| --- | --- | --- |
| `frontend/src/stores/siteProfile.ts` | `site_profile` | 负责品牌、SEO 与文档级同步 |
| `frontend/src/stores/siteFooter.ts` | `footer` | 负责页尾备案与说明 |
| `frontend/src/stores/homeHero.ts` | `home_hero` | 负责首页首屏文案与本地缓存 |
| `frontend/src/stores/homeAssets.ts` | `home_assets` | 负责首页背景资源与本地缓存 |
| `frontend/src/stores/authorProfile.ts` | `author_profile` | 负责作者展示身份资料 |
| `frontend/src/stores/siteIntegrations.ts` | `site_integrations` | 负责外部依赖配置与首页顶部组件开关 |

## 4. 旧入口清理结果

| 旧入口 | 当前状态 | 新入口 |
| --- | --- | --- |
| `/admin/footer` | 路由已移除 | `/admin/settings?section=site-footer` |
| `/admin/profile` | 路由已移除 | `/admin/settings?section=admin-profile` |

说明：

- 侧边栏已收口为“设置中心”统一入口。
- 不再保留兼容跳转页，历史收藏或旧文档需更新为设置中心地址。

## 5. 禁止继续迁入后台的内容

以下内容默认禁止继续迁入 `site_settings`：

- 后台系统按钮文案、toast、校验提示
- `aria-label` 与无障碍辅助说明
- 组件级 loading / empty / error 技术文案
- 状态枚举映射，例如“草稿 / 已发布 / 待审核”
- Token、API Secret、JWT Secret、数据库与 Redis 凭证
- 仅服务于开发和排障的技术常量

这些内容应继续留在：

- `frontend/src/content/copy/`
- 后端环境变量
- 业务表自身字段

## 6. 允许进入后台的内容

以下内容继续允许走设置中心：

- 首页门面内容
- 作者公开展示资料
- 页尾备案与运营说明
- GitHub 用户名、天气地点
- 首页顶部组件开关
- 可频繁调整、且直接影响对外展示的公开信息

## 7. 后续扩展规则

新增设置项时，必须同时满足：

1. 有明确业务域，不与现有 key 混写。
2. 有 `copy/` 默认值或明确空态回退。
3. 能说明前台消费位置。
4. 不包含敏感信息。
5. 能在设置中心标明影响范围。

如果不满足以上条件，优先继续留在 `copy/` 或业务表，而不是新增后台配置。
