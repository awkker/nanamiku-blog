// @ts-check

import mdx from '@astrojs/mdx';
import sitemap from '@astrojs/sitemap';
import { defineConfig } from 'astro/config';

import vue from '@astrojs/vue';
import { siteCopy } from './src/content/copy/site.ts';

// 构建期站点 URL：优先读环境变量 SITE_URL（后台站点设置导出），回退到 siteCopy 静态默认值。
const buildSiteUrl = process.env.SITE_URL || siteCopy.seo.siteUrl;

// https://astro.build/config
export default defineConfig({
    site: buildSiteUrl,
    integrations: [mdx(), sitemap(), vue()],
    vite: {
        optimizeDeps: {
            exclude: ['marked'],
        },
        server: {
            proxy: {
                '/api': {
                    target: 'http://localhost:8080',
                    changeOrigin: true,
                },
            },
        },
    },
});
