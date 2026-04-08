import { siteCopy } from './content/copy'

// 构建期 SEO 常量：优先读环境变量（后台站点设置导出），回退到 siteCopy 静态默认值。
export const SITE_TITLE = import.meta.env.SITE_TITLE || siteCopy.seo.siteTitle
export const SITE_DESCRIPTION = import.meta.env.SITE_DESCRIPTION || siteCopy.seo.defaultDescription
export const SITE_URL = import.meta.env.SITE_URL || siteCopy.seo.siteUrl
export const SITE_DEFAULT_OG_IMAGE = import.meta.env.SITE_DEFAULT_OG_IMAGE || siteCopy.seo.defaultSocialImage
